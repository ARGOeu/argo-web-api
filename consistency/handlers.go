package consistency

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const conColName = "consistency"
const conID = "consistency-status"

// HandleSubrouter for api access to consistency information
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {
	respond.PrepAppRoutes(s, confhandler, appRoutes)
}

var appRoutes = []respond.AppRoutes{
	{"consistency.result", "GET", "/consistency", GetStatus},
	{"consistency.options", "OPTIONS", "/consistency", Options},
	{"consistency.ack", "POST", "/consistency/ack", PostAck},
	{"consistency.ack.options", "OPTIONS", "/consistency/ack", Options},
	{"consistency.auto-check", "POST", "/consistency/auto-check", PostAutoCheck},
	{"consistency.auto-check.options", "OPTIONS", "/consistency/auto-check", Options},
}

// GetStatus gets consistency status
func GetStatus(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	var output []byte
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Try to get mongo client and target consistency collection
	conCol := cfg.MongoClient.Database(cfg.MongoDB.Db).Collection(conColName)

	urlValues := r.URL.Query()
	verbose := urlValues.Has("verbose")

	var data DataMongo
	var result Result
	filter := bson.M{"_id": conID}

	err = conCol.FindOne(context.TODO(), filter).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			code = 404
			output, _ = createMessage("Constistency information is not yet available", code, contentType)
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	now := time.Now().UTC()

	// by default result equals to the result of the auto check
	result.Status = data.AutoCheckStatus
	result.Message = data.AutoCheckMsg
	result.Timestamp = now.Format(time.RFC3339)

	// add auto-check details
	if verbose {
		result.AutoCheckMsg = data.AutoCheckMsg
		result.AutoCheckStatus = data.AutoCheckStatus
		result.AutoCheckTimestamp = data.AutoCheckTimestamp
	}
	// check if ack exists and applies still
	if data.AckStatus != "" && data.AckMsg != "" && data.AckTimeoutHours > 0 {
		ackTime, err := time.Parse(time.RFC3339, data.AckTimestamp)
		if err == nil {
			ackDur := time.Duration(data.AckTimeoutHours) * time.Hour
			if now.Sub(ackTime.UTC()) < ackDur {
				// make status get the ack result
				result.Status = data.AckStatus
				result.Message = data.AckMsg
				if verbose {
					result.AckStatus = data.AckStatus
					result.AckMsg = data.AckMsg
					result.AckTimeoutHours = data.AckTimeoutHours
					result.AckTimestamp = data.AckTimestamp
				}
			}

		}

	}

	output, err = respond.MarshalContent(result, contentType, "", " ")

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

// PostAutoCheck posts results about an auto check
func PostAutoCheck(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	var output []byte
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Try to get mongo client and target consistency collection
	conCol := cfg.MongoClient.Database(cfg.MongoDB.Db).Collection(conColName)

	var autoCheckInput AutoCheckMongo

	body, err := io.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &autoCheckInput); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
		code = http.StatusBadRequest
		return code, h, output, err
	}

	autoCheckInput.AutoCheckTimestamp = time.Now().UTC().Format(time.RFC3339)

	filter := bson.M{"_id": conID}
	update := bson.M{"$set": autoCheckInput}
	opts := options.Update().SetUpsert(true)

	_, err = conCol.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Fatal("MongoDB update failed:", err)
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	message := "The Auto Check event was posted succesfully"
	output, err = createMessage(message, code, contentType)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

// PostAck posts an acknowledgement
func PostAck(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	var output []byte
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Try to get mongo client and target consistency collection
	conCol := cfg.MongoClient.Database(cfg.MongoDB.Db).Collection(conColName)

	var ackInput AckMongo

	body, err := io.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &ackInput); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
		code = http.StatusBadRequest
		return code, h, output, err
	}

	ackInput.AckTimestamp = time.Now().UTC().Format(time.RFC3339)
	if ackInput.AckTimeoutHours == 0 {
		// default is 6 hours - TODO: to be configurable
		ackInput.AckTimeoutHours = 6
	}

	filter := bson.M{"_id": conID}
	update := bson.M{"$set": ackInput}
	opts := options.Update().SetUpsert(false)

	updateResult, err := conCol.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Fatal("MongoDB update failed:", err)
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if updateResult.MatchedCount == 0 {
		message := "There is no auto check event yet to acknowledge"
		output, err = createMessage(message, http.StatusNotFound, contentType)
		return http.StatusNotFound, h, output, err
	}

	message := "The Ack event was posted succesfully"
	output, err = createMessage(message, code, contentType)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

func Options(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/plain"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	h.Set("Allow", "GET, OPTIONS")
	return code, h, output, err

}
