package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const conColName = "consistency"
const conID = "consistency-status"

// HandleSubrouter for api access to health status
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {
	respond.PrepAppRoutes(s, confhandler, appRoutes)
}

var appRoutes = []respond.AppRoutes{
	{"health", "GET", "/health", GetStatus},
	{"health", "OPTIONS", "/health", Options},
}

// GetStatus gets health status
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

	var data DataMongo
	var result Result
	filter := bson.M{"_id": conID}

	err = conCol.FindOne(context.TODO(), filter).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			code = 404
			output, _ = createMessage("Health status not yet available", code, contentType)
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

	// check if ack exists and applies still
	if data.AckStatus != "" && data.AckMsg != "" && data.AckTimeoutHours > 0 {
		ackTime, err := time.Parse(time.RFC3339, data.AckTimestamp)
		if err == nil {
			ackDur := time.Duration(data.AckTimeoutHours) * time.Hour
			if now.Sub(ackTime.UTC()) < ackDur {
				// make status get the ack result
				result.Status = data.AckStatus
				result.Message = data.AckMsg
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
