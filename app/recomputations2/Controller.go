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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
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

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	filter := IncomingRecomputation{
		StartTime: urlValues.Get("start_time"),
		EndTime:   urlValues.Get("end_time"),
		Reason:    urlValues.Get("reason"),
		Report:    urlValues.Get("report"),
	}

	session, err := mongo.OpenSession(tenantDbConfig)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	results := []MongoInterface{}
	err = mongo.Find(session, tenantDbConfig.Db, recomputationsColl, filter, "timestamp", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

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
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	filter := IncomingRecomputation{
		ID: vars["ID"],
	}
	session, err := mongo.OpenSession(tenantDbConfig)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	result := MongoInterface{}
	err = mongo.FindOne(session, tenantDbConfig.Db, recomputationsColl, filter, &result)

	if err != nil {
		if err.Error() == "not found" {
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
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	session, err := mongo.OpenSession(tenantDbConfig)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	var recompSubmission IncomingRecomputation
	// urlValues := r.URL.Query()

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
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
		ID:             mongo.NewUUID(),
		RequesterName:  recompSubmission.RequesterName,
		RequesterEmail: recompSubmission.RequesterEmail,
		StartTime:      recompSubmission.StartTime,
		EndTime:        recompSubmission.EndTime,
		Reason:         recompSubmission.Reason,
		Report:         recompSubmission.Report,
		Exclude:        recompSubmission.Exclude,
		Timestamp:      now.Format("2006-01-02T15:04:05Z"),
		Status:         "pending",
		History:        history,
	}

	err = mongo.Insert(session, tenantDbConfig.Db, recomputationsColl, recomputation)

	if err != nil {
		panic(err)
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
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	session, err := mongo.OpenSession(tenantDbConfig)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	query := bson.M{"id": vars["ID"]}

	// Retrieve Results from database
	results := []MongoInterface{}
	err = mongo.Find(session, tenantDbConfig.Db, recomputationsColl, query, "", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Check if nothing found
	if len(results) < 1 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	statusItem := HistoryItem{Status: "pending", Timestamp: results[0].Timestamp}
	history := []HistoryItem{statusItem}

	recomputation := MongoInterface{
		ID:             vars["ID"],
		RequesterName:  results[0].RequesterName,
		RequesterEmail: results[0].RequesterEmail,
		StartTime:      results[0].StartTime,
		EndTime:        results[0].EndTime,
		Reason:         results[0].Reason,
		Report:         results[0].Report,
		Exclude:        results[0].Exclude,
		Status:         "pending",
		Timestamp:      results[0].Timestamp,
		History:        history,
	}

	if err = mongo.Update(session, tenantDbConfig.Db, recomputationsColl, query, recomputation); err != nil {
		code = http.StatusInternalServerError
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
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	session, err := mongo.OpenSession(tenantDbConfig)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	var statusSubmit IncomingStatus
	// urlValues := r.URL.Query()

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	query := bson.M{"id": vars["ID"]}

	// Retrieve Results from database
	results := []MongoInterface{}
	err = mongo.Find(session, tenantDbConfig.Db, recomputationsColl, query, "", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Check if nothing found
	if len(results) < 1 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	if err := json.Unmarshal(body, &statusSubmit); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
		code = http.StatusBadRequest
		return code, h, output, err
	}

	if isValidStatus(statusSubmit.Status) == false {
		code = http.StatusConflict
		output, _ = respond.MarshalContent(
			respond.ErrConflict("status should be among values: \"pending\",\"approved\",\"rejected\",\"running\",\"done\""),
			contentType, "", " ")
		return code, h, output, err

	}

	now := time.Now()
	history := results[0].History
	statusItem := HistoryItem{Status: statusSubmit.Status, Timestamp: now.Format("2006-01-02T15:04:05Z")}
	history = append(history, statusItem)

	recomputation := MongoInterface{
		ID:             vars["ID"],
		RequesterName:  results[0].RequesterName,
		RequesterEmail: results[0].RequesterEmail,
		StartTime:      results[0].StartTime,
		EndTime:        results[0].EndTime,
		Reason:         results[0].Reason,
		Report:         results[0].Report,
		Exclude:        results[0].Exclude,
		Status:         statusSubmit.Status,
		Timestamp:      results[0].Timestamp,
		History:        history,
	}

	if err = mongo.Update(session, tenantDbConfig.Db, recomputationsColl, query, recomputation); err != nil {
		code = http.StatusInternalServerError
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
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	session, err := mongo.OpenSession(tenantDbConfig)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	var recompSubmission IncomingRecomputation
	// urlValues := r.URL.Query()

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	query := bson.M{"id": vars["ID"]}

	// Retrieve Results from database
	results := []MongoInterface{}
	err = mongo.Find(session, tenantDbConfig.Db, recomputationsColl, query, "", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Check if nothing found
	if len(results) < 1 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = 404
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
		Status:         results[0].Status,
		Timestamp:      results[0].Timestamp,
	}

	if err = mongo.Update(session, tenantDbConfig.Db, recomputationsColl, query, recomputation); err != nil {
		code = http.StatusInternalServerError
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
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Open session to tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	filter := bson.M{"id": vars["ID"]}

	// Retrieve Results from database
	results := []MongoInterface{}
	err = mongo.Find(session, tenantDbConfig.Db, recomputationsColl, filter, "", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Check if nothing found
	if len(results) < 1 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	mongo.Remove(session, tenantDbConfig.Db, recomputationsColl, filter)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createMsgView("Recomputation Successfully Deleted", 200)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}
