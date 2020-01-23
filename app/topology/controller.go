/*
 * Copyright (c) 2018 GRNET S.A.
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

package topology

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// datastore collection name that contains aggregations profile records
const topoColName = "topology"
const endpointColName = "topology_endpoints"
const groupColName = "topology_groups"

func getCloseDate(c *mgo.Collection, dt int) int {
	dateQuery := bson.M{"date_integer": bson.M{"$lte": dt}}
	result := Endpoint{}
	err := c.Find(dateQuery).One(&result)
	if err != nil {
		return -1
	}
	return result.DateInt
}

// ListTopoStats list statistics about the topology used in the report
func ListTopoStats(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	groupType := context.Get(r, "group_type").(string)
	egroupType := context.Get(r, "endpoint_group_type").(string)

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	reportName := vars["report_name"]
	//Time Related
	dateStr := urlValues.Get("date")

	const zuluForm = "2006-01-02"
	const ymdForm = "20060102"

	ts := time.Now().UTC()
	if dateStr != "" {
		ts, _ = time.Parse(zuluForm, dateStr)
	}

	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// find the report id first
	reportID, err := mongo.GetReportID(session, tenantDbConfig.Db, reportName)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	var serviceResults []string
	var egroupResults []string
	var groupResults []string

	serviceCol := session.DB(tenantDbConfig.Db).C("service_ar")
	eGroupCol := session.DB(tenantDbConfig.Db).C("endpoint_group_ar")

	err = serviceCol.Find(bson.M{"report": reportID, "date": tsYMD}).Distinct("name", &serviceResults)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	err = eGroupCol.Find(bson.M{"report": reportID, "date": tsYMD}).Distinct("name", &egroupResults)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	err = eGroupCol.Find(bson.M{"report": reportID, "date": tsYMD}).Distinct("supergroup", &groupResults)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// fill the topology JSON struct
	result := Topology{}
	result.GroupCount = len(groupResults)
	result.GroupType = groupType
	result.GroupList = groupResults
	result.EndGroupCount = len(egroupResults)
	result.EndGroupType = egroupType
	result.EndGroupList = egroupResults
	result.ServiceCount = len(serviceResults)
	result.ServiceList = serviceResults

	output, err = createTopoView(result, contentType, http.StatusOK)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

// ListEndpoints by date
func ListEndpoints(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Open session to tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	colEndpoint := session.DB(tenantDbConfig.Db).C(endpointColName)

	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		return code, h, output, err
	}

	expDate := getCloseDate(colEndpoint, dt)

	results := []Endpoint{}

	if expDate < 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	err = colEndpoint.Find(bson.M{"date_integer": expDate}).All(&results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createListEndpoint(results, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// CreateEndpoints Creates a list of Endpoint Groups fon an item
func CreateEndpoints(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
		return code, h, output, err
	}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	incoming := []Endpoint{}
	// Try ingest request body
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))

	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// check if topology already exists for current day

	existing := Endpoint{}
	endpointCol := session.DB(tenantDbConfig.Db).C(endpointColName)
	err = endpointCol.Find(bson.M{"date_integer": dt}).One(&existing)
	if err != nil {
		// Stop at any error except not found. We want to have not found
		if err.Error() != "not found" {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
		// else continue correctly -
	} else {
		// If found we need to inform user that the topology is already created for this date
		output, err = createMessageOUT(fmt.Sprintf("Topology already exists for date: %s, please either update it or delete it first!", dateStr), 409, "json")
		code = 409
		return code, h, output, err
	}

	// Parse body json
	if err := json.Unmarshal(body, &incoming); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
		code = 400
		return code, h, output, err
	}

	for i := range incoming {
		incoming[i].Date = dateStr
		incoming[i].DateInt = dt
	}

	err = mongo.MultiInsert(session, tenantDbConfig.Db, endpointColName, incoming)

	if err != nil {
		panic(err)
	}

	// Create view of the results
	output, err = createMessageOUT(fmt.Sprintf("Topology of %d endpoints created for date: %s", len(incoming), dateStr), 201, "json") //Render the results into JSON
	code = 201
	return code, h, output, err
}

// Options implements the option request on resource
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
	h.Set("Allow", fmt.Sprintf("GET, OPTIONS"))
	return code, h, output, err

}
