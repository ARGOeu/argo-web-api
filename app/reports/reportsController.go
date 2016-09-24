/*
 * Copyright (c) 2014 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
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
 * or implied, of either GRNET S.A., SRCE or IN2P3 CNRS Computing
 * Centre
 *
 * The work represented by this source file is partially funded by
 * the EGI-InSPIRE project through the European Commission's 7th
 * Framework Programme (contract # INFSO-RI-261323)
 */

package reports

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"gopkg.in/mgo.v2/bson"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/mux"
)

var reportsColl = "reports"

// Create function is used to implement the create report request.
// The request is an http POST request with the report description
// provided as json structure in the request body
func Create(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusCreated
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

	//Reading the json input from the request body
	reqBody, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))

	if err != nil {
		return code, h, output, err
	}
	input := MongoInterface{}
	//Unmarshalling the json input into byte form

	err = json.Unmarshal(reqBody, &input)

	// Check if json body is malformed
	if err != nil {
		output, _ := respond.MarshalContent(respond.MalformedJSONInput, contentType, "", " ")
		code = http.StatusBadRequest
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	// Try to open the mongo session
	session, err := mongo.OpenSession(tenantDbConfig)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Validate profiles given in report
	validationErrors := input.ValidateProfiles(session.DB(tenantDbConfig.Db))

	if len(validationErrors) > 0 {
		code = 422
		out := respond.UnprocessableEntity
		out.Errors = validationErrors
		output = out.MarshalTo(contentType)
		return code, h, output, err
	}

	// Prepare structure for storing query results
	results := []MongoInterface{}

	// Check if report with the same name exists in datastore
	query := searchName(input.Info.Name)
	err = mongo.Find(session, tenantDbConfig.Db, reportsColl, query, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// If results are returned for the specific name
	// then we already have an existing report and we must
	// abort creation notifing the user
	if len(results) > 0 {
		// Name was found so print the error message in xml
		out := respond.ResponseMessage{
			Status: respond.StatusResponse{
				Message: "Report with the same name already exists",
				Code:    strconv.Itoa(http.StatusConflict),
			}}

		output, _ = respond.MarshalContent(out, contentType, "", " ")

		code = http.StatusConflict
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err

	}
	input.Info.Created = time.Now().Format("2006-01-02 15:04:05")
	input.Info.Updated = input.Info.Created
	input.ID = mongo.NewUUID()
	// If no report exists with this name create a new one

	err = mongo.Insert(session, tenantDbConfig.Db, reportsColl, input)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Notify user that the report has been created. In xml style
	selfLink := "https://" + r.Host + r.URL.Path + "/" + input.ID
	output, err = SubmitSuccesful(input, contentType, selfLink)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err

}

// List function that implements the http GET request that retrieves
// all avaiable report information
func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	urlValues := r.URL.Query()

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Try to open the mongo session
	session, err := mongo.OpenSession(cfg.MongoDB)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	query := bson.M{}

	if urlValues.Get("name") != "" {
		query["info.name"] = urlValues["name"]
	}

	// Create structure for storing query results
	results := []MongoInterface{}
	// Query tenant collection for all available documents.
	// nil query param == match everything
	err = mongo.Find(session, tenantDbConfig.Db, reportsColl, nil, "id", &results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// After successfully retrieving the db results
	// call the createView function to render them into idented xml
	output, err = createView(results, contentType)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// ListOne function that implements the http GET request that retrieves
// all avaiable report information
func ListOne(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	//Extracting urlvar "name" from url path

	id := mux.Vars(r)["id"]
	// Try to open the mongo session
	session, err := mongo.OpenSession(tenantDbConfig)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create structure for storing query results
	result := MongoInterface{}
	// Create a simple query object to query by name
	query := bson.M{"id": id}
	// Query collection tenants for the specific tenant name
	err = mongo.FindOne(session, tenantDbConfig.Db, reportsColl, query, &result)

	// If query returned zero result then no tenant matched this name,
	// abort and notify user accordingly
	if err != nil {
		code = http.StatusNotFound
		output, _ := ReportNotFound(contentType)
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	// After successfully retrieving the db results
	// call the createView function to render them into idented xml
	output, err = createView([]MongoInterface{result}, contentType)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	code = http.StatusOK
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// Update function used to implement update report request.
// This is an http PUT request that gets a specific report's name
// as a urlvar parameter input and a json structure in the request
// body in order to update the datastore document for the specific
// report
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

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	//Extracting report name from url
	id := mux.Vars(r)["id"]

	//Reading the json input
	reqBody, err := ioutil.ReadAll(r.Body)

	input := MongoInterface{}
	//Unmarshalling the json input into byte form
	err = json.Unmarshal(reqBody, &input)

	if err != nil {

		// User provided malformed json input data
		output, _ := respond.MarshalContent(respond.MalformedJSONInput, contentType, "", " ")
		code = http.StatusBadRequest
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	sanitizedInput := bson.M{
		"$set": bson.M{
			// "info": bson.M{
			"info.name":        input.Info.Name,
			"info.description": input.Info.Description,
			"info.updated":     time.Now().Format("2006-01-02 15:04:05"),
			// },
			"profiles":        input.Profiles,
			"filter_tags":     input.Tags,
			"topology_schema": input.Topology,
		}}

	// Try to open the mongo session
	session, err := mongo.OpenSession(tenantDbConfig)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Validate profiles given in report
	validationErrors := input.ValidateProfiles(session.DB(tenantDbConfig.Db))

	if len(validationErrors) > 0 {
		code = 422
		out := respond.UnprocessableEntity
		out.Errors = validationErrors
		output = out.MarshalTo(contentType)
		return code, h, output, err
	}

	// We search by name and update
	query := bson.M{"id": id}
	err = mongo.Update(session, tenantDbConfig.Db, reportsColl, query, sanitizedInput)
	if err != nil {
		if err.Error() != "not found" {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
		//Render the response into XML
		code = http.StatusNotFound
		output, err = ReportNotFound(contentType)
		return code, h, output, err
	}

	//Render the response into XML
	output, err = respond.CreateResponseMessage("Report was successfully updated", "200", contentType)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	code = http.StatusOK
	return code, h, output, err

}

// Delete function used to implement remove report request
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

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	//Extracting record id from url
	id := mux.Vars(r)["id"]

	// Try to open the mongo session
	session, err := mongo.OpenSession(tenantDbConfig)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// We search by name and delete the document in db
	query := bson.M{"id": id}
	info, err := mongo.Remove(session, tenantDbConfig.Db, reportsColl, query)

	if err != nil {
		if err.Error() != "not found" {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
		//Render the response into XML
		code = http.StatusNotFound
		output, err = ReportNotFound(contentType)
		return code, h, output, err
	}

	// info.Removed > 0 means that many documents have been removed
	// If deletion took place we notify user accordingly.
	// Else we notify that no tenant matched the specific name
	if info.Removed > 0 {
		code = http.StatusOK
		output, err = respond.CreateResponseMessage("Report was successfully deleted", "200", contentType)
	} else {
		code = http.StatusNotFound
		output, err = ReportNotFound(contentType)
	}
	//Render the response into XML
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
	h.Set("Allow", fmt.Sprintf("GET, POST, DELETE, PUT, OPTIONS"))
	return code, h, output, err

}
