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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	gcontext "github.com/gorilla/context"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	//Reading the json input from the request body
	reqBody, err := io.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))

	if err != nil {
		return code, h, output, err
	}
	input := MongoInterface{}
	//Unmarshalling the json input into byte form

	err = json.Unmarshal(reqBody, &input)

	// check if user declared any thresholds or else provide defaults
	if input.Thresholds == nil {
		t := defaultThresholds()
		input.Thresholds = &t
	}

	// check if user declared what needs to be computed or else provide defaults
	if input.Computations == nil {
		c := genDefaultComp()
		input.Computations = c
	}

	// Check if json body is malformed
	if err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
		code = 400
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	db := cfg.MongoClient.Database(tenantDbConfig.Db)

	// Validate profiles given in report

	validationErrors := input.ValidateProfiles(db)
	validationErrors = append(validationErrors, input.ValidateTrends()...)

	if len(validationErrors) > 0 {
		code = 422
		out := respond.UnprocessableEntity
		out.Errors = validationErrors
		output = out.MarshalTo(contentType)
		return code, h, output, err
	}

	// Check if report with the same name exists in datastore
	query := searchName(input.Info.Name)

	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(reportsColl)
	queryResult := rCol.FindOne(context.TODO(), query)

	if queryResult.Err() == nil {
		output, _ = respond.MarshalContent(respond.ErrConflict("Report with the same name already exists"), contentType, "", " ")
		code = http.StatusConflict
		return code, h, output, err
	}

	if queryResult.Err() != mongo.ErrNoDocuments {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// If results are returned for the specific name
	// then we already have an existing report and we must
	// abort creation notifing the user

	input.Info.Created = time.Now().Format("2006-01-02 15:04:05")
	input.Info.Updated = input.Info.Created
	input.ID = utils.NewUUID()
	// If no report exists with this name create a new one

	_, err = rCol.InsertOne(context.TODO(), input)

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
// all available report information
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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)
	tenantName := gcontext.Get(r, "tenant_name").(string)

	query := bson.M{}

	if urlValues.Get("name") != "" {
		query["info.name"] = urlValues["name"]
	}

	// Create structure for storing query results
	results := []MongoInterface{}
	// Query tenant collection for all available documents.
	// nil query param == match everything
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(reportsColl)
	findOptions := options.Find().SetSort(bson.D{{Key: "id", Value: 1}})
	cursor, err := rCol.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	for indx := range results {
		results[indx].Tenant = tenantName

		// check if computations field is not set and return the default value
		if results[indx].Computations == nil {
			results[indx].Computations = genDefaultComp()
		}
	}

	// After successfully retrieving the db results
	// call the createView function to render them into indented xml
	output, err = createView(results, contentType)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// ListOne function that implements the http GET request that retrieves
// the specified report's information
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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)
	tenantName := gcontext.Get(r, "tenant_name").(string)

	//Extracting urlvar "name" from url path

	id := mux.Vars(r)["id"]

	// Create structure for storing query results
	result := MongoInterface{}
	// Create a simple query object to query by name
	query := bson.M{"id": id}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(reportsColl)
	// Query collection tenants for the specific tenant name
	err = rCol.FindOne(context.TODO(), query).Decode(&result)

	// If query returned zero result then no tenant matched this name,
	// abort and notify user accordingly
	if err != nil {
		if err == mongo.ErrNoDocuments {
			output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
			code = 404
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Enrich report with tenant name -- used in argo engine
	result.Tenant = tenantName

	if result.Computations == nil {
		result.Computations = genDefaultComp()
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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	//Extracting report name from url
	id := mux.Vars(r)["id"]

	//Reading the json input
	reqBody, err := io.ReadAll(r.Body)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	input := MongoInterface{}
	//Unmarshalling the json input into byte form
	err = json.Unmarshal(reqBody, &input)

	if err != nil {

		// User provided malformed json input data
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
		code = http.StatusBadRequest
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	if input.Computations == nil {
		input.Computations = genDefaultComp()
	}

	sanitizedInput := bson.M{
		"$set": bson.M{
			// "info": bson.M{
			"info.name":        input.Info.Name,
			"info.description": input.Info.Description,
			"info.updated":     time.Now().Format("2006-01-02 15:04:05"),
			"weight":           input.Weight,
			"disabled":         input.Disabled,
			"thresholds":       input.Thresholds,
			"computations":     input.Computations,
			// },
			"profiles":        input.Profiles,
			"filter_tags":     input.Tags,
			"topology_schema": input.Topology,
		}}

	db := cfg.MongoClient.Database(tenantDbConfig.Db)
	// Validate profiles given in report
	validationErrors := input.ValidateProfiles(db)
	validationErrors = append(validationErrors, input.ValidateTrends()...)

	if len(validationErrors) > 0 {
		code = 422
		out := respond.UnprocessableEntity
		out.Errors = validationErrors
		output = out.MarshalTo(contentType)
		return code, h, output, err
	}

	queryById := bson.M{"id": id}

	// before updating, check if the report exists and the name is unique
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(reportsColl)
	result := MongoInterface{}

	err = rCol.FindOne(context.TODO(), queryById).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
			code = http.StatusNotFound
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if result.Info.Name != input.Info.Name {

		queryByName := bson.M{"info.name": input.Info.Name}

		queryResult := rCol.FindOne(context.TODO(), queryByName)

		if queryResult.Err() == nil {
			output, _ = respond.MarshalContent(respond.ErrConflict("Report with the same name already exists"), contentType, "", " ")
			code = http.StatusConflict
			return code, h, output, err

		}

		if queryResult.Err() != mongo.ErrNoDocuments {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

	}

	updateChange, err := rCol.UpdateOne(context.TODO(), queryById, sanitizedInput)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if updateChange.MatchedCount == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = 404
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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(reportsColl)

	//Extracting record id from url
	id := mux.Vars(r)["id"]

	// We search by name and delete the document in db
	query := bson.M{"id": id}
	deleteResult, err := rCol.DeleteOne(context.TODO(), query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if deleteResult.DeletedCount == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	code = http.StatusOK
	output, err = respond.CreateResponseMessage("Report was successfully deleted", "200", contentType)
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
	h.Set("Allow", "GET, POST, DELETE, PUT, OPTIONS")
	return code, h, output, err

}
