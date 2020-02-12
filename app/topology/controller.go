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
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/app/reports"

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

	fltr := fltrEndpoint{}
	fltr.Group = urlValues.Get("group")
	fltr.GroupType = urlValues.Get("type")
	fltr.Hostname = urlValues.Get("hostname")
	fltr.Service = urlValues.Get("service")
	fltr.Tags = urlValues.Get("tags")

	expDate := getCloseDate(colEndpoint, dt)

	results := []Endpoint{}

	if expDate < 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	err = colEndpoint.Find(prepEndpointQuery(expDate, fltr)).All(&results)
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

// CreateEndpoints Creates a list of endpoints for a specific date
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

// DeleteEndpoints deletes a list of endpoints topology for a specific date
func DeleteEndpoints(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	endpointCol := session.DB(tenantDbConfig.Db).C(endpointColName)
	change, err := endpointCol.RemoveAll(bson.M{"date_integer": dt})
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if change.Removed == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	// Create view of the results
	output, err = createMessageOUT(fmt.Sprintf("Topology of %d endpoints deleted for date: %s", change.Removed, dateStr), 200, "json")
	return code, h, output, err
}

// CreateGroups creates a list of groups as a group topology for a specific date
func CreateGroups(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	incoming := []Group{}
	// Try ingest request body
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))

	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// check if topology already exists for current day

	existing := Group{}
	groupCol := session.DB(tenantDbConfig.Db).C(groupColName)
	err = groupCol.Find(bson.M{"date_integer": dt}).One(&existing)
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

	err = mongo.MultiInsert(session, tenantDbConfig.Db, groupColName, incoming)

	if err != nil {
		panic(err)
	}

	// Create view of the results
	output, err = createMessageOUT(fmt.Sprintf("Topology of %d groups created for date: %s", len(incoming), dateStr), 201, "json") //Render the results into JSON
	code = 201
	return code, h, output, err
}

func handleWildcard(item string) (string, bool) {
	if strings.Contains(item, "*") {
		return "^" + strings.Replace(item, "*", ".*", -1) + "$", true
	}
	return item, false
}

func appendTags(query bson.M, tags string) bson.M {
	// split tags in kv paris
	kvs := strings.Split(tags, ",")
	for _, kv := range kvs {
		// split kv in key and value
		kvsplit := strings.Split(kv, ":")
		// if not properly split ignore
		if len(kvsplit) != 2 {
			continue
		}

		if value, reg := handleWildcard(kvsplit[1]); reg == true {
			query["tags."+kvsplit[0]] = bson.RegEx{Pattern: value}
		} else {
			query["tags."+kvsplit[0]] = value
		}

	}

	return query

}

func prepEndpointQuery(date int, filter fltrEndpoint) bson.M {

	query := bson.M{"date_integer": date}
	// if filter struct not empty begin adding filters

	if (fltrEndpoint{} != filter) {
		if filter.Group != "" {
			if value, reg := handleWildcard(filter.Group); reg == true {
				query["group"] = bson.RegEx{Pattern: value}
			} else {
				query["group"] = value
			}

		}
		if filter.GroupType != "" {
			if value, reg := handleWildcard(filter.GroupType); reg == true {
				query["type"] = bson.RegEx{Pattern: value}
			} else {
				query["type"] = value
			}
		}
		if filter.Service != "" {
			if value, reg := handleWildcard(filter.Service); reg == true {
				query["service"] = bson.RegEx{Pattern: value}
			} else {
				query["service"] = value
			}
		}
		if filter.Hostname != "" {
			if value, reg := handleWildcard(filter.Hostname); reg == true {
				query["hostname"] = bson.RegEx{Pattern: value}
			} else {
				query["hostname"] = value
			}
		}
	}

	// check if tags exist to append them to query
	if filter.Tags != "" {
		appendTags(query, filter.Tags)
	}

	return query
}

func prepGroupQuery(date int, filter fltrGroup) bson.M {

	query := bson.M{"date_integer": date}
	// if filter struct not empty begin adding filters
	if (fltrGroup{} != filter) {
		if filter.Group != "" {
			if value, reg := handleWildcard(filter.Group); reg == true {
				query["group"] = bson.RegEx{Pattern: value}
			} else {
				query["group"] = value
			}

		}
		if filter.GroupType != "" {
			if value, reg := handleWildcard(filter.GroupType); reg == true {
				query["type"] = bson.RegEx{Pattern: value}
			} else {
				query["type"] = value
			}
		}
		if filter.Subgroup != "" {
			if value, reg := handleWildcard(filter.Subgroup); reg == true {
				query["subgroup"] = bson.RegEx{Pattern: value}
			} else {
				query["subgroup"] = value
			}
		}
	}

	// check if tags exist to append them to query
	if filter.Tags != "" {
		appendTags(query, filter.Tags)
	}

	return query
}

func getReportEndpointGroupType(r reports.MongoInterface) string {

	if r.Topology.Group != nil {
		if r.Topology.Group.Group != nil {
			return r.Topology.Group.Group.Type
		}
	}

	return ""
}

func getReportTags(r reports.MongoInterface) string {
	tagStr := ""
	first := false
	for _, tag := range r.Tags {
		if tag.Context == "argo.group.filter.tags" {
			if !first {
				tagStr = tagStr + ","
			} else {
				first = false
			}

			tagStr = tagStr + tag.Name + ":" + tag.Value
		}
	}
	return tagStr
}

func getReportGroupType(r reports.MongoInterface) string {

	if r.Topology.Group != nil {
		return r.Topology.Group.Type
	}

	return ""
}

//ListGroupsByReport lists group topology by report
func ListGroupsByReport(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	reportName := vars["report"]

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Open session to tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	colGroup := session.DB(tenantDbConfig.Db).C(groupColName)
	colReports := session.DB(tenantDbConfig.Db).C("reports")
	//get the report

	report := reports.MongoInterface{}
	err = colReports.Find(bson.M{"info.name": reportName}).One(&report)

	if err != nil {
		if err.Error() == "not found" {
			output, err = createMessageOUT(fmt.Sprintf("No report with name: %s exists!", reportName), 404, "json")
			code = 404
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	groupType := getReportGroupType(report)
	groupTags := getReportTags(report)

	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		return code, h, output, err
	}

	fltr := fltrGroup{}
	if groupType != "" {
		fltr.GroupType = groupType
	}
	if groupTags != "" {
		fltr.Tags = groupTags
	}

	expDate := getCloseDate(colGroup, dt)

	results := []Group{}

	if expDate < 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	err = colGroup.Find(prepGroupQuery(expDate, fltr)).All(&results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createListGroup(results, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// ListGroups by date
func ListGroups(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	colGroup := session.DB(tenantDbConfig.Db).C(groupColName)

	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		return code, h, output, err
	}

	fltr := fltrGroup{}
	fltr.Group = urlValues.Get("group")
	fltr.GroupType = urlValues.Get("type")
	fltr.Subgroup = urlValues.Get("subgroup")
	fltr.Tags = urlValues.Get("tags")

	expDate := getCloseDate(colGroup, dt)

	results := []Group{}

	if expDate < 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	err = colGroup.Find(prepGroupQuery(expDate, fltr)).All(&results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createListGroup(results, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// DeleteGroups deletes a list of groups topology for a specific date
func DeleteGroups(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	groupCol := session.DB(tenantDbConfig.Db).C(groupColName)
	change, err := groupCol.RemoveAll(bson.M{"date_integer": dt})
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if change.Removed == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	// Create view of the results
	output, err = createMessageOUT(fmt.Sprintf("Topology of %d groups deleted for date: %s", change.Removed, dateStr), 200, "json")
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
