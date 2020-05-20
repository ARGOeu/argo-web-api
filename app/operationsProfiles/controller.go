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

package operationsProfiles

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/context"
	"gopkg.in/mgo.v2/bson"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/mux"
)

// datastore collection name that contains operations profile records
const opsColName = "operations_profiles"

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
				// operations_profiles collection is meant to have an index with date_integer:-1 and id:1 so
				// when searching by date the documents are sorted with the recent timestamp first
				// so we need the recent item available to our query timepoint which is specific date
				"id":               bson.M{"$first": "$id"},
				"date":             bson.M{"$first": "$date"},
				"name":             bson.M{"$first": "$name"},
				"available_states": bson.M{"$first": "$available_states"},
				"defaults":         bson.M{"$first": "$defaults"},
				"operations":       bson.M{"$first": "$operations"},
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
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Open session to tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Retrieve Results from database
	result := OpsProfile{}

	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}
	opsQuery := prepQuery(dt, vars["ID"])

	opsCol := session.DB(tenantDbConfig.Db).C(opsColName)
	err = opsCol.Find(opsQuery).One(&result)
	if err != nil {
		if err.Error() == "not found" {
			output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
			code = 404
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createListView([]OpsProfile{result}, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// List the existing metric profiles for the tenant making the request
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
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Open session to tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}
	opsQuery := prepMultiQuery(dt, name)

	opsCol := session.DB(tenantDbConfig.Db).C(opsColName)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	results := []OpsProfile{}
	err = opsCol.Pipe(opsQuery).All(&results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createListView(results, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

//Create a new metric profile
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
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)
	urlValues := r.URL.Query()
	dateStr := urlValues.Get("date")
	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	incoming := OpsProfile{}
	incoming.DateInt = dt
	incoming.Date = dateStr

	// Try ingest request body
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
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

	// Validate States
	var errList []string
	errList = append(errList, incoming.validateDuplicates()...)
	errList = append(errList, incoming.validateStates()...)
	errList = append(errList, incoming.validateMentions()...)

	if len(errList) > 0 {
		output, err = createErrView("Validation Error", 422, errList)
		code = 422
		return code, h, output, err
	}

	// check if the operations profile's name is unique
	results := []OpsProfile{}
	query := bson.M{"name": incoming.Name}

	err = mongo.Find(session, tenantDbConfig.Db, opsColName, query, "", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// If results are returned for the specific name
	// then we already have an existing report and we must
	// abort creation notifying the user
	if len(results) > 0 {
		output, _ = respond.MarshalContent(respond.ErrConflict("Operations profile with the same name already exists"), contentType, "", " ")
		code = http.StatusConflict
		return code, h, output, err
	}

	// Generate new id
	incoming.ID = mongo.NewUUID()
	err = mongo.Insert(session, tenantDbConfig.Db, opsColName, incoming)

	if err != nil {
		panic(err)
	}

	// Create view of the results
	output, err = createRefView(incoming, "Operations Profile successfully created", 201, r) //Render the results into JSON
	code = 201
	return code, h, output, err
}

//Update function to update contents of an existing metric profile
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
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	incoming := OpsProfile{}
	incoming.DateInt = dt
	incoming.Date = dateStr

	// ingest body data
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
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

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// create filter to retrieve specific profile with id
	filter := bson.M{"id": vars["ID"], "date_integer": bson.M{"$lte": dt}}

	incoming.ID = vars["ID"]

	// Retrieve Results from database
	results := []OpsProfile{}
	err = mongo.Find(session, tenantDbConfig.Db, opsColName, filter, "name", &results)

	if err != nil {
		panic(err)
	}

	// Check if nothing found
	if len(results) < 1 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	// Validate States
	var errList []string
	errList = append(errList, incoming.validateDuplicates()...)
	errList = append(errList, incoming.validateStates()...)
	errList = append(errList, incoming.validateMentions()...)

	if len(errList) > 0 {
		output, err = createErrView("Validation Error", 422, errList)
		code = 422
		return code, h, output, err
	}

	// check if the operations profile's name is unique
	if incoming.Name != results[0].Name {

		results = []OpsProfile{}
		query := bson.M{"name": incoming.Name, "id": bson.M{"$ne": vars["ID"]}}

		err = mongo.Find(session, tenantDbConfig.Db, opsColName, query, "", &results)

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		// If results are returned for the specific name
		// then we already have an existing operations profile and we must
		// abort creation notifying the user
		if len(results) > 0 {
			output, _ = respond.MarshalContent(respond.ErrConflict("Operations profile with the same name already exists"), contentType, "", " ")
			code = http.StatusConflict
			return code, h, output, err
		}
	}

	// run the update query

	opsCol := session.DB(tenantDbConfig.Db).C(opsColName)
	info, err := opsCol.Upsert(bson.M{"id": vars["ID"], "date_integer": dt}, incoming)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	updMsg := "Operations Profile successfully updated"

	if info.Updated <= 0 {
		updMsg = "Operations Profile successfully updated (new history snapshot)"
	}

	// Create view for response message
	output, err = createMsgView(updMsg, 200) //Render the results into JSON
	code = 200
	return code, h, output, err
}

//Delete metric profile based on id
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
	results := []OpsProfile{}
	err = mongo.Find(session, tenantDbConfig.Db, opsColName, filter, "name", &results)

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

	// Remove also removes multiple documents with specific id, thus erasing all history
	mongo.Remove(session, tenantDbConfig.Db, opsColName, filter)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createMsgView("Operations Profile Successfully Deleted", 200) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
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
