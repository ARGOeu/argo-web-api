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
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

func isAdminRestricted(roles []string) bool {
	return len(roles) > 0 && roles[0] == "super_admin_restricted"
}

func isAdminUI(roles []string) bool {
	return len(roles) > 0 && roles[0] == "super_admin_ui"
}

// Provides a global tenant status (true/false) based on tenant's status details
func calcTotalStatus(details StatusDetail) bool {
	// Check first tenant configuration in argo-engine
	if details.EngineConfig == false {
		return false
	}
	// Check tenant configuration statuses regarding AMS service
	if details.AMS.MetricData.Ingestion == false {
		return false
	}
	if details.AMS.MetricData.Publishing == false {
		return false
	}
	if details.AMS.MetricData.StatusStreaming == false {
		return false
	}
	if details.AMS.SyncData.Ingestion == false {
		return false
	}
	if details.AMS.SyncData.Publishing == false {
		return false
	}
	if details.AMS.SyncData.StatusStreaming == false {
		return false
	}
	// Check tenant configuration statuses regarding HDFS
	if details.HDFS.MetricData == false {
		return false
	}
	for _, item := range details.HDFS.SyncData {
		if item.AggregationProf == false {
			return false
		}
		if item.ConfigProf == false {
			return false
		}
		if item.Donwtimes == false {
			return false
		}
		if item.GroupEndpoints == false {
			return false
		}
		if item.GroupGroups == false {
			return false
		}
		if item.MetricProf == false {
			return false
		}
		if item.OpsProf == false {
			return false
		}
		if item.Recomp == false {
			return false
		}
		if item.Weight == false {
			return false
		}
	}
	return true
}

func restrictTenantOutput(results []Tenant) []Tenant {
	restricted := []Tenant{}
	for _, tenant := range results {
		rItem := Tenant{}
		rItem.ID = tenant.ID
		rItem.Info = tenant.Info
		rItem.Topology = tenant.Topology
		restricted = append(restricted, rItem)
	}
	return restricted
}

func removeNonUIUsers(results []Tenant) []Tenant {
	restricted := []Tenant{}
	for _, tenant := range results {
		uiUsers := []TenantUser{}
		rItem := Tenant{}
		rItem.ID = tenant.ID
		rItem.Info = tenant.Info
		rItem.Topology = tenant.Topology
		for _, user := range tenant.Users {
			for _, role := range user.Roles {
				if role == "admin_ui" {
					uiUsers = append(uiUsers, user)
				}
			}
		}
		rItem.Users = uiUsers

		restricted = append(restricted, rItem)
	}
	return restricted
}

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
	incoming.Topology = TopologyInfo{TopoType: "", Feed: ""}

	// Parse body json
	if err := json.Unmarshal(body, &incoming); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
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

	if errMsg, errCode := validateTenantUsers(incoming, session, cfg); errMsg != "" && errCode != 0 {
		output, _ = respond.MarshalContent(respond.ErrConflict(errMsg), contentType, "", " ")
		code = errCode
		return code, h, output, err
	}

	// generate a unique id for each of the tenant users
	for idx := range incoming.Users {
		incoming.Users[idx].ID = mongo.NewUUID()
	}

	// Check if name exists
	sameName := []Tenant{}
	filter := bson.M{"info.name": incoming.Info.Name}
	err = mongo.Find(session, cfg.MongoDB.Db, "tenants", filter, "name", &sameName)

	if len(sameName) > 0 {
		output, _ = respond.MarshalContent(respond.ErrConflict("Tenant with same name already exists"), contentType, "", " ")
		code = http.StatusConflict
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
	output, err = createRefView(incoming, "Tenant was successfully created", 201, r) //Render the results into JSON
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

	roles := context.Get(r, "roles").([]string)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create structure for storing query results
	results := []Tenant{}
	// Query tenant collection for all available documents.
	// nil query param == match everything
	if isAdminUI(roles) {
		err = mongo.Find(session, cfg.MongoDB.Db, "tenants", bson.M{"users.roles": "admin_ui"}, "name", &results)
	} else {
		err = mongo.Find(session, cfg.MongoDB.Db, "tenants", nil, "name", &results)
	}

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Quicky check if super admin is restricted to remove restricted info
	if isAdminRestricted(roles) {
		results = restrictTenantOutput(results)
	}

	// remove non ui users from results
	if isAdminUI(roles) {
		results = removeNonUIUsers(results)
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

// ListStatus show tenant status
func ListStatus(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	results := []TenantStatus{}

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
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = http.StatusNotFound
		return code, h, output, err
	}

	results[0].Status.TotalStatus = calcTotalStatus(results[0].Status)

	// After successfully retrieving the db results
	// call the createView function to render them into idented xml
	output, err = createStatusView(results, "Success", code)

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

	roles := context.Get(r, "roles").([]string)

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

	if isAdminUI(roles) {
		query = bson.M{"id": vars["ID"], "users.roles": "admin_ui"}
	}

	// Query collection tenants for the specific tenant id
	err = mongo.Find(session, cfg.MongoDB.Db, "tenants", query, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Check if nothing found
	if len(results) < 1 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = http.StatusNotFound
		return code, h, output, err
	}

	if isAdminUI(roles) {
		results = removeNonUIUsers(results)
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

func UpdateStatus(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	incomingStatus := StatusDetail{}

	// ingest body data
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	// parse body json
	if err := json.Unmarshal(body, &incomingStatus); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
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

	// Retrieve Results from database
	results := []Tenant{}
	err = mongo.Find(session, cfg.MongoDB.Db, "tenants", filter, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Check if nothing found
	if len(results) < 1 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = http.StatusNotFound
		return code, h, output, err
	}

	filter = bson.M{"id": vars["ID"]}

	// update and set only the status field
	setIncoming := bson.M{"$set": bson.M{"status": incomingStatus}}
	err = mongo.Update(session, cfg.MongoDB.Db, "tenants", filter, setIncoming)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view for response message
	output, err = createMsgView("Tenant successfully updated", 200) //Render the results into JSON
	code = http.StatusOK
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
	incoming.Topology = TopologyInfo{TopoType: "", Feed: ""}

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
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = http.StatusNotFound
		return code, h, output, err
	}

	if errMsg, errCode := validateTenantUsers(incoming, session, cfg); errMsg != "" && errCode != 0 {
		output, _ = respond.MarshalContent(respond.ErrConflict(errMsg), contentType, "", " ")
		code = errCode
		return code, h, output, err
	}

	// If user chose to change name - check if name already exists
	if results[0].Info.Name != incoming.Info.Name {
		sameName := []Tenant{}
		filter = bson.M{"info.name": incoming.Info.Name}

		err = mongo.Find(session, cfg.MongoDB.Db, "tenants", filter, "name", &sameName)

		if len(sameName) > 0 {
			code = http.StatusConflict
			output, _ = respond.MarshalContent(respond.ErrConflict("Tenant with same name already exists"), contentType, "", " ")
			return code, h, output, err
		}
	}

	// save all the previous users' ids
	// use the apikey since it is a unique field
	ids := map[string]string{}
	for _, u := range results[0].Users {
		ids[u.APIkey] = u.ID
	}

	// for the old users, reuse their ids
	// for the new ones, generate new ids
	for idx, u := range incoming.Users {
		// check if the user was already present
		id, found := ids[u.APIkey]
		if found {
			incoming.Users[idx].ID = id
			continue
		}
		// generate new uuid
		incoming.Users[idx].ID = mongo.NewUUID()
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
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
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

// Delete function used to implement remove tenant request
func GetUserByID(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	exportFilter := r.URL.Query().Get("export")

	// Try to open the mongo session
	session, err := mongo.OpenSession(cfg.MongoDB)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	query := bson.M{"users.id": vars["ID"]}
	results := []Tenant{}

	err = mongo.Find(session, cfg.MongoDB.Db, "tenants", query, "", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if len(results) == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = http.StatusNotFound
		return code, h, output, err
	}

	for _, user := range results[0].Users {
		if user.ID == vars["ID"] {
			output, err = createUserView(user, "User was successfully retrieved", 200, exportFilter)
			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}
		}
	}

	return code, h, output, err
}

// validateTenantUsers validates the uniqueness of the tenant's users' keys
func validateTenantUsers(tenant Tenant, session *mgo.Session, cfg config.Config) (string, int) {

	usersKeys := make(map[string]bool)
	errMsg := ""
	errCode := 0

	for _, tUser := range tenant.Users {

		results := []Tenant{}

		// for each of the tenant's users make sure there is no other that holds the same key across all users from all tenants
		// check that there are no others users registered under different tenants that might have the same api key

		query := bson.M{
			"$and": []bson.M{
				{
					"id": bson.M{
						"$ne": tenant.ID}},
				{
					"users": bson.M{
						"$elemMatch": bson.M{
							"api_key": tUser.APIkey}}}}}

		if err := mongo.Find(session, cfg.MongoDB.Db, "tenants", query, "", &results); err != nil {
			return err.Error(), http.StatusInternalServerError
		}

		if len(results) > 0 {
			errMsg = fmt.Sprintf("More than one users found using the key: %v", tUser.APIkey)
			errCode = http.StatusConflict
			return errMsg, errCode
		}

		// use a map with all the keys that we have evaluated to check whether or not users inside the same tenant have the same key declared
		// when we evaluate a key, we try to see if we have seen that key again in a previous user

		if _, ok := usersKeys[tUser.APIkey]; ok {
			errMsg = fmt.Sprintf("More than one users found using the key: %v", tUser.APIkey)
			errCode = http.StatusConflict
			return errMsg, errCode
		}

		// if the current key isn't present inside the map, add it
		usersKeys[tUser.APIkey] = true
	}

	return errMsg, errCode

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
