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

package aggregationProfiles

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	gcontext "github.com/gorilla/context"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// datastore collection name that contains aggregations profile records
const aggColName = "aggregation_profiles"

func prepMultiQuery(dt int, name string) interface{} {

	matchQuery := bson.M{"date_integer": bson.M{"$lte": dt}}

	if name != "" {
		matchQuery["name"] = name
	}

	return []bson.M{
		{
			"$match": matchQuery,
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"id": "$id",
				},
				// aggregation_profiles collection is meant to have an index with date_integer:-1 and id:1 so
				// when searching by date the documents are sorted with the recent timestamp first
				// so we need the recent item available to our query timepoint which is specific date
				"id":                bson.M{"$first": "$id"},
				"date":              bson.M{"$first": "$date"},
				"name":              bson.M{"$first": "$name"},
				"namespace":         bson.M{"$first": "$namespace"},
				"endpoint_group":    bson.M{"$first": "$endpoint_group"},
				"metric_operation":  bson.M{"$first": "$metric_operation"},
				"profile_operation": bson.M{"$first": "$profile_operation"},
				"metric_profile":    bson.M{"$first": "$metric_profile"},
				"groups":            bson.M{"$first": "$groups"},
			},
		},
		{
			"$sort": bson.M{"id": 1},
		},
	}

}

func prepQuery(dt int, id string) interface{} {

	return bson.M{"date_integer": bson.M{"$lte": dt}, "id": id}

}

// ListOne handles the listing of one specific profile based on its given id
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

	vars := mux.Vars(r)
	urlValues := r.URL.Query()
	dateStr := urlValues.Get("date")

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	// Retrieve Results from database
	result := AggProfile{}
	dt, _, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}
	aggQuery := prepQuery(dt, vars["ID"])

	aggCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(aggColName)
	err = aggCol.FindOne(context.TODO(), aggQuery).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
			code = 404
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	// Create view of the results
	output, err = createListView([]AggProfile{result}, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

// List the existing aggregation profiles for the tenant making the request
// Also there is an optional url param "name" to filter results by
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
	dateStr := urlValues.Get("date")
	name := urlValues.Get("name")

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	dt, _, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}
	aggQuery := prepMultiQuery(dt, name)

	aggCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(aggColName)

	results := []AggProfile{}
	cursor, err := aggCol.Aggregate(context.TODO(), aggQuery)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &results)
	// Create view of the results
	output, err = createListView(results, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// Create a new metric profile
func Create(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	urlValues := r.URL.Query()
	dateStr := urlValues.Get("date")
	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	incoming := AggProfile{}
	incoming.DateInt = dt
	incoming.Date = dateStr

	// Try ingest request body
	body, err := io.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// Parse body json
	if err := json.Unmarshal(body, &incoming); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
		code = 400
		return code, h, output, err
	}

	aggCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(aggColName)
	metricCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("metric_profiles")
	// Validate
	err = incoming.MetricProf.validateID(metricCol)
	// Respond 422 unprocessabe entity
	if err != nil {
		output, _ = respond.MarshalContent(respond.ErrUnprocessableEntity("Referenced metric profile ID is not found"), contentType, "", " ")
		code = 422
		return code, h, output, err
	}

	// check if the aggregation profile's name is unique
	query := bson.M{"name": incoming.Name}

	queryResult := aggCol.FindOne(context.TODO(), query)

	if queryResult.Err() == nil {
		output, _ = respond.MarshalContent(respond.ErrConflict("Aggregation profile with the same name already exists"), contentType, "", " ")
		code = http.StatusConflict
		return code, h, output, err
	}

	if queryResult.Err() != mongo.ErrNoDocuments {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Generate new id
	incoming.ID = utils.NewUUID()
	_, err = aggCol.InsertOne(context.TODO(), incoming)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createRefView(incoming, "Aggregation Profile successfully created", 201, r) //Render the results into JSON
	code = 201
	return code, h, output, err
}

// Update function to update contents of an existing aggregation profile
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
	urlValues := r.URL.Query()
	dateStr := urlValues.Get("date")
	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	incoming := AggProfile{}
	incoming.DateInt = dt
	incoming.Date = dateStr

	// ingest body data
	body, err := io.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	// parse body json
	if err := json.Unmarshal(body, &incoming); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
		code = 400
		return code, h, output, err
	}

	aggCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(aggColName)
	metricCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("metric_profiles")
	// create filter to retrieve specific profile with id
	query := bson.M{"id": vars["ID"]}

	incoming.ID = vars["ID"]

	// Retrieve Results from database

	queryResult := aggCol.FindOne(context.TODO(), query)

	if queryResult.Err() != nil {
		if queryResult.Err() == mongo.ErrNoDocuments {
			output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
			code = 404
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	queryCheck := bson.M{"name": incoming.Name, "id": bson.M{"$ne": vars["ID"]}}

	queryResult = aggCol.FindOne(context.TODO(), queryCheck)

	if queryResult.Err() == nil {
		output, _ = respond.MarshalContent(respond.ErrConflict("Aggregation profile with the same name already exists"), contentType, "", " ")
		code = http.StatusConflict
		return code, h, output, err
	}

	if queryResult.Err() != mongo.ErrNoDocuments {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Validate
	err = incoming.MetricProf.validateID(metricCol)
	// Respond 422 unprocessabe entity
	if err != nil {
		output, _ = respond.MarshalContent(respond.ErrUnprocessableEntity("Referenced metric profile ID is not found"), contentType, "", " ")
		code = 422
		return code, h, output, err
	}

	replaceResult, err := aggCol.ReplaceOne(context.TODO(), bson.M{"id": vars["ID"], "date_integer": dt}, incoming, options.Replace().SetUpsert(true))
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if replaceResult.MatchedCount == 0 && replaceResult.UpsertedCount == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = http.StatusNotFound
		return code, h, output, err
	}

	updMsg := "Aggregations Profile successfully updated"

	if replaceResult.UpsertedCount > 0 {
		updMsg = "Aggregations Profile successfully updated (new history snapshot)"
	}

	// Create view for response message
	output, err = createMsgView(updMsg, 200) //Render the results into JSON
	code = 200
	return code, h, output, err
}

// Delete metric profile based on id
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

	aggCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(aggColName)

	query := bson.M{"id": vars["ID"]}

	deleteResult, err := aggCol.DeleteMany(context.TODO(), query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if deleteResult.DeletedCount == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = http.StatusNotFound
		return code, h, output, err
	}

	// Create view of the results
	output, err = createMsgView("Aggregation Profile Successfully Deleted", 200) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// Options implements the http option request
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
