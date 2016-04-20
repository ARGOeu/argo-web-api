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

package tenants

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/mux"
)

// Create function is used to implement the create tenant request.
// The request is an http POST request with the tenant description
// provided as json structure in the request body
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

	// Try ingest request body
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	incoming := Tenant{}

	// Parse body json
	if err := json.Unmarshal(body, &incoming); err != nil {

		output, _ = respond.MarshalContent(respond.BadRequestBadJSON, contentType, "", " ")
		code = http.StatusBadRequest
		return code, h, output, err
	}

	// Try to open the mongo session
	session, err := mongo.OpenSession(cfg.MongoDB)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Check if name exists
	sameName := []Tenant{}
	filter := bson.M{"info.name": incoming.Info.Name}
	err = mongo.Find(session, cfg.MongoDB.Db, "tenants", filter, "name", &sameName)

	if len(sameName) > 0 {
		code = http.StatusConflict
		output, err = createMsgView("Tenant with same name already exists", code)
		return code, h, output, err
	}

	// Generate new id
	incoming.ID = mongo.NewUUID()
	incoming.Info.Created = time.Now().Format("2006-01-02 15:04:05")
	incoming.Info.Updated = incoming.Info.Created
	err = mongo.Insert(session, cfg.MongoDB.Db, "tenants", incoming)

	if err != nil {
		panic(err)
	}

	// Create view of the results
	output, err = createRefView(incoming, "Tenant was succesfully created", 201, r) //Render the results into JSON
	code = http.StatusCreated
	return code, h, output, err
}

// List function that implements the http GET request that retrieves
// all avaiable tenant information
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

	// Try to open the mongo session
	session, err := mongo.OpenSession(cfg.MongoDB)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create structure for storing query results
	results := []Tenant{}
	// Query tenant collection for all available documents.
	// nil query param == match everything
	err = mongo.Find(session, cfg.MongoDB.Db, "tenants", nil, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	// After successfully retrieving the db results
	// call the createView function to render them into idented xml
	output, err = createListView(results, "Success", code)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// ListOne function implement an http GET request that accepts
// a name parameter urlvar and retrieves information only for the
// specific tenant
func ListOne(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	vars := mux.Vars(r)

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Try to open the mongo session
	session, err := mongo.OpenSession(cfg.MongoDB)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create structure to hold query results
	results := []Tenant{}

	// Create a simple query object to query by id
	query := bson.M{"id": vars["ID"]}
	// Query collection tenants for the specific tenant id
	err = mongo.Find(session, cfg.MongoDB.Db, "tenants", query, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Check if nothing found
	if len(results) < 1 {
		output, _ = respond.MarshalContent(respond.NotFound, contentType, "", " ")
		code = http.StatusNotFound
		return code, h, output, err
	}

	// After successfully retrieving the db results
	// call the createView function to render them into idented xml
	output, err = createListView(results, "Success", code)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// Update function used to implement update tenant request.
// This is an http PUT request that gets a specific tenant's name
// as a urlvar parameter input and a json structure in the request
// body in order to update the datastore document for the specific
// tenant
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

	incoming := Tenant{}

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
		output, _ = respond.MarshalContent(respond.BadRequestBadJSON, contentType, "", " ")
		code = http.StatusBadRequest
		return code, h, output, err
	}

	// Try to open the mongo session
	session, err := mongo.OpenSession(cfg.MongoDB)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// create filter to retrieve specific profile with id
	filter := bson.M{"id": vars["ID"]}

	incoming.ID = vars["ID"]

	// Retrieve Results from database
	results := []Tenant{}
	err = mongo.Find(session, cfg.MongoDB.Db, "tenants", filter, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Check if nothing found
	if len(results) < 1 {

		output, _ = respond.MarshalContent(respond.NotFound, contentType, "", " ")
		code = http.StatusNotFound
		return code, h, output, err
	}

	// If user chose to change name - check if name already exists
	if results[0].Info.Name != incoming.Info.Name {
		sameName := []Tenant{}
		filter = bson.M{"info.name": incoming.Info.Name}

		err = mongo.Find(session, cfg.MongoDB.Db, "tenants", filter, "name", &sameName)

		if len(sameName) > 1 {
			code = http.StatusConflict
			output, err = createMsgView("Tenant with same name already exists", code)
			return code, h, output, err
		}
	}

	// run the update query
	incoming.Info.Created = results[0].Info.Created

	incoming.Info.Updated = time.Now().Format("2006-01-02 15:04:05")
	filter = bson.M{"id": vars["ID"]}
	err = mongo.Update(session, cfg.MongoDB.Db, "tenants", filter, incoming)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view for response message
	output, err = createMsgView("Tenant successfully updated", 200) //Render the results into JSON
	code = http.StatusOK
	return code, h, output, err

}

// Delete function used to implement remove tenant request
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

	// Try to open the mongo session
	session, err := mongo.OpenSession(cfg.MongoDB)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	filter := bson.M{"id": vars["ID"]}

	// Retrieve Results from database
	results := []Tenant{}
	err = mongo.Find(session, cfg.MongoDB.Db, "tenants", filter, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Check if nothing found
	if len(results) < 1 {
		output, _ = respond.MarshalContent(respond.NotFound, contentType, "", " ")
		code = http.StatusNotFound
		return code, h, output, err
	}

	mongo.Remove(session, cfg.MongoDB.Db, "tenants", filter)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createMsgView("Tenant Successfully Deleted", 200) //Render the results into JSON

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
	h.Set("Allow", fmt.Sprintf("GET,POST,PUT,DELETE,OPTIONS"))
	return code, h, output, err

}
