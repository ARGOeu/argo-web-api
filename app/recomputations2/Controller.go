/*
 * Copyright (c) 2015 GRNET S.A.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the
 * License. You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an "AS
 * IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language
 * governing permissions and limitations under the License.
 *
 * The views and conclusions contained in the software and
 * documentation are those of the authors and should not be
 * interpreted as representing official policies, either expressed
 * or implied, of GRNET S.A.
 *
 */

package recomputations2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils"
	"github.com/ARGOeu/argo-web-api/utils/config"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var recomputationsColl = "recomputations"

// List existing recomputations
func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	urlValues := r.URL.Query()
	// date argument expected in the format YYYY-MM-DD
	qDate := urlValues.Get("date")
	qReport := urlValues.Get("report")

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	query := bson.M{}

	// check if there are relevant recomputations for given date
	if qDate != "" {
		if respond.ValidateDateOnly(qDate) != nil {
			output, _ = respond.MarshalContent(respond.ErrBadRequestDetails("date argument should be in the YYYY-MM-DD format"), contentType, "", " ")
			code = http.StatusBadRequest
			return code, h, output, err
		}

		query["$where"] = fmt.Sprintf("'%s' >= this.start_time.split('T')[0] && '%s' <= this.end_time.split('T')[0]", qDate, qDate)
	}
	if qReport != "" {
		query["report"] = qReport
	}

	results := []MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(recomputationsColl)
	findOptions := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})
	cursor, err := rCol.Find(context.TODO(), query, findOptions)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)
	output, err = createListView(results, contentType)

	return code, h, output, err

}

// ListOne lists a single recomputation according to the given id
func ListOne(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	// contentType := "application/json"
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	query := bson.M{"id": vars["ID"]}

	result := MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(recomputationsColl)
	err = rCol.FindOne(context.TODO(), query).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
			code = http.StatusNotFound
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createListView(result, contentType)

	return code, h, output, err

}

// SubmitRecomputation insert a new pending recomputation in the tenants database
func SubmitRecomputation(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
	//STANDARD DECLARATIONS START
	code := http.StatusAccepted
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	var recompSubmission IncomingRecomputation
	// urlValues := r.URL.Query()

	body, err := io.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &recompSubmission); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
		code = http.StatusBadRequest
		return code, h, output, err
	}
	now := time.Now()

	statusItem := HistoryItem{Status: "pending", Timestamp: now.Format("2006-01-02T15:04:05Z")}
	history := []HistoryItem{statusItem}

	recomputation := MongoInterface{
		ID:               utils.NewUUID(),
		RequesterName:    recompSubmission.RequesterName,
		RequesterEmail:   recompSubmission.RequesterEmail,
		StartTime:        recompSubmission.StartTime,
		EndTime:          recompSubmission.EndTime,
		Reason:           recompSubmission.Reason,
		Report:           recompSubmission.Report,
		Exclude:          recompSubmission.Exclude,
		ExcludeMetrics:   recompSubmission.ExcludeMetrics,
		ExcludeMonSource: recompSubmission.ExcludeMonSource,
		Timestamp:        now.Format("2006-01-02T15:04:05Z"),
		Status:           "pending",
		History:          history,
	}

	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(recomputationsColl)

	_, err = rCol.InsertOne(context.TODO(), recomputation)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createSubmitView(recomputation, contentType, r)
	return code, h, output, err
}

// ResetStatus resets status changes back to pending when recomputation was created
func ResetStatus(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	query := bson.M{"id": vars["ID"]}

	// Retrieve Results from database
	result := MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(recomputationsColl)
	err = rCol.FindOne(context.TODO(), query).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
			code = 404
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	statusItem := HistoryItem{Status: "pending", Timestamp: result.Timestamp}
	history := []HistoryItem{statusItem}

	recomputation := MongoInterface{
		ID:             vars["ID"],
		RequesterName:  result.RequesterName,
		RequesterEmail: result.RequesterEmail,
		StartTime:      result.StartTime,
		EndTime:        result.EndTime,
		Reason:         result.Reason,
		Report:         result.Report,
		Exclude:        result.Exclude,
		Status:         "pending",
		Timestamp:      result.Timestamp,
		History:        history,
	}

	replaceResult, err := rCol.ReplaceOne(context.TODO(), query, recomputation)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if replaceResult.MatchedCount == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	output, err = createMsgView("Recomputation status reset successfully to: pending", 200)

	return code, h, output, err
}

func isValidStatus(status string) bool {

	switch status {
	case
		"pending",
		"approved",
		"rejected",
		"running",
		"done":
		return true
	}

	return false
}

// ChangeStatus updates the status of an existing recomputation
func ChangeStatus(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	var statusSubmit IncomingStatus
	// urlValues := r.URL.Query()

	body, err := io.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	query := bson.M{"id": vars["ID"]}

	// Retrieve Results from database
	result := MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(recomputationsColl)
	err = rCol.FindOne(context.TODO(), query).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
			code = 404
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if err := json.Unmarshal(body, &statusSubmit); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
		code = http.StatusBadRequest
		return code, h, output, err
	}

	if !isValidStatus(statusSubmit.Status) {
		code = http.StatusConflict
		output, _ = respond.MarshalContent(
			respond.ErrConflict("status should be among values: \"pending\",\"approved\",\"rejected\",\"running\",\"done\""),
			contentType, "", " ")
		return code, h, output, err

	}

	now := time.Now()
	history := result.History
	statusItem := HistoryItem{Status: statusSubmit.Status, Timestamp: now.Format("2006-01-02T15:04:05Z")}
	history = append(history, statusItem)

	recomputation := MongoInterface{
		ID:             vars["ID"],
		RequesterName:  result.RequesterName,
		RequesterEmail: result.RequesterEmail,
		StartTime:      result.StartTime,
		EndTime:        result.EndTime,
		Reason:         result.Reason,
		Report:         result.Report,
		Exclude:        result.Exclude,
		Status:         statusSubmit.Status,
		Timestamp:      result.Timestamp,
		History:        history,
	}

	replaceResult, err := rCol.ReplaceOne(context.TODO(), query, recomputation)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if replaceResult.MatchedCount == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	output, err = createMsgView("Recomputation status updated successfully to: "+statusSubmit.Status, 200)
	return code, h, output, err
}

// Update updates an already existing recomputation
func Update(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	var recompSubmission IncomingRecomputation
	// urlValues := r.URL.Query()

	body, err := io.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	query := bson.M{"id": vars["ID"]}

	// Retrieve Result from database
	result := MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(recomputationsColl)
	err = rCol.FindOne(context.TODO(), query).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
			code = 404
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if err := json.Unmarshal(body, &recompSubmission); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
		code = http.StatusBadRequest
		return code, h, output, err
	}

	recomputation := MongoInterface{
		ID:             vars["ID"],
		RequesterName:  recompSubmission.RequesterName,
		RequesterEmail: recompSubmission.RequesterEmail,
		StartTime:      recompSubmission.StartTime,
		EndTime:        recompSubmission.EndTime,
		Reason:         recompSubmission.Reason,
		Report:         recompSubmission.Report,
		Exclude:        recompSubmission.Exclude,
		Status:         result.Status,
		Timestamp:      result.Timestamp,
	}

	replaceResult, err := rCol.ReplaceOne(context.TODO(), query, recomputation)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if replaceResult.MatchedCount == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	output, err = createMsgView("Recomputation updated successfully", 200)
	return code, h, output, err
}

// Delete recomputation
func Delete(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(recomputationsColl)

	query := bson.M{"id": vars["ID"]}

	deleteResult, err := rCol.DeleteOne(context.TODO(), query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if deleteResult.DeletedCount == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = http.StatusNotFound
		return code, h, output, err
	}
	output, err = createMsgView("Recomputation Successfully Deleted", 200)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}
