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
	"sort"
	"strings"

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

	err = serviceCol.Find(bson.M{"report": reportID, "date": dt}).Distinct("name", &serviceResults)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	err = eGroupCol.Find(bson.M{"report": reportID, "date": dt}).Distinct("name", &egroupResults)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	err = eGroupCol.Find(bson.M{"report": reportID, "date": dt}).Distinct("supergroup", &groupResults)
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

// ListTopoTags list statistics the distinct values appearing for each tag in topology items
func ListTopoTags(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

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
	colGroup := session.DB(tenantDbConfig.Db).C(groupColName)

	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	expDate := getCloseDate(colEndpoint, dt)

	resTagsEndpoint := []TagValues{}
	resTagsGroup := []TagValues{}

	if expDate < 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	err = colEndpoint.Pipe(prepTagAggr(expDate)).All(&resTagsEndpoint)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	expDate = getCloseDate(colGroup, dt)

	if expDate < 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	err = colGroup.Pipe(prepTagAggr(expDate)).All(&resTagsGroup)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// sort unique tag values in each tag occurence
	for _, tg := range resTagsEndpoint {
		sort.Strings(tg.Values)
	}

	for _, tg := range resTagsGroup {
		sort.Strings(tg.Values)
	}

	resTags := []TagInfo{{Name: "endpoints", Values: resTagsEndpoint}, {Name: "groups", Values: resTagsGroup}}

	output, err = createListTags(resTags, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
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
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	fltr := fltrEndpoint{}
	fltr.Group = urlValues["group"]
	fltr.GroupType = urlValues["type"]
	fltr.Hostname = urlValues["hostname"]
	fltr.Service = urlValues["service"]
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
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
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
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
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
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
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
	// split tags in kv parts
	tagmap := map[string][]string{}
	kvs := strings.Split(tags, ",")
	for _, kv := range kvs {
		// split kv in key and value
		kvsplit := strings.Split(kv, ":")
		// if not properly split ignore
		if len(kvsplit) != 2 {
			continue
		}

		if _, exists := tagmap[kvsplit[0]]; exists == false {
			tagmap[kvsplit[0]] = []string{}
		}

		tagmap[kvsplit[0]] = append(tagmap[kvsplit[0]], kvsplit[1])

	}

	// iterate over map
	for tag := range tagmap {
		if len(tagmap[tag]) > 1 {
			tagRegStr := "^("
			for i, tagVal := range tagmap[tag] {
				if i != 0 {
					tagRegStr = tagRegStr + "|"
				}
				if value, reg := handleWildcard(tagVal); reg == true {
					tagRegStr = tagRegStr + value
				} else {
					tagRegStr = tagRegStr + tagVal
				}
			}
			tagRegStr = tagRegStr + ")$"
			query["tags."+tag] = bson.RegEx{Pattern: tagRegStr}

		} else {
			if value, reg := handleWildcard(tagmap[tag][0]); reg == true {
				query["tags."+tag] = bson.RegEx{Pattern: value}
			} else {
				query["tags."+tag] = value
			}
		}
	}

	return query

}

func prepEndpointQuery(date int, filter fltrEndpoint) bson.M {

	query := bson.M{"date_integer": date}
	// if filter group has values
	if len(filter.Group) > 0 {
		if len(filter.Group) > 1 {
			groupRegStr := "^("
			for i, groupVal := range filter.Group {
				if i != 0 {
					groupRegStr = groupRegStr + "|"
				}
				if value, reg := handleWildcard(groupVal); reg == true {
					groupRegStr = groupRegStr + value
				} else {
					groupRegStr = groupRegStr + groupVal
				}
			}
			groupRegStr = groupRegStr + ")$"
			query["group"] = bson.RegEx{Pattern: groupRegStr}
		} else {
			if value, reg := handleWildcard(filter.Group[0]); reg == true {
				query["group"] = bson.RegEx{Pattern: value}
			} else {
				query["group"] = value
			}
		}
	}

	// if filter group type has values
	if len(filter.GroupType) > 0 {
		if len(filter.GroupType) > 1 {
			groupRegStr := "^("
			for i, groupVal := range filter.GroupType {
				if i != 0 {
					groupRegStr = groupRegStr + "|"
				}
				if value, reg := handleWildcard(groupVal); reg == true {
					groupRegStr = groupRegStr + value
				} else {
					groupRegStr = groupRegStr + groupVal
				}
			}
			groupRegStr = groupRegStr + ")$"
			query["type"] = bson.RegEx{Pattern: groupRegStr}
		} else {
			if value, reg := handleWildcard(filter.GroupType[0]); reg == true {
				query["type"] = bson.RegEx{Pattern: value}
			} else {
				query["type"] = value
			}
		}
	}

	// if filter service has values
	if len(filter.Service) > 0 {
		if len(filter.Service) > 1 {
			groupRegStr := "^("
			for i, groupVal := range filter.Service {
				if i != 0 {
					groupRegStr = groupRegStr + "|"
				}
				if value, reg := handleWildcard(groupVal); reg == true {
					groupRegStr = groupRegStr + value
				} else {
					groupRegStr = groupRegStr + groupVal
				}
			}
			groupRegStr = groupRegStr + ")$"
			query["service"] = bson.RegEx{Pattern: groupRegStr}
		} else {
			if value, reg := handleWildcard(filter.Service[0]); reg == true {
				query["service"] = bson.RegEx{Pattern: value}
			} else {
				query["service"] = value
			}
		}
	}

	// if filter hostname has values
	if len(filter.Hostname) > 0 {
		if len(filter.Hostname) > 1 {
			groupRegStr := "^("
			for i, groupVal := range filter.Hostname {
				if i != 0 {
					groupRegStr = groupRegStr + "|"
				}
				if value, reg := handleWildcard(groupVal); reg == true {
					groupRegStr = groupRegStr + value
				} else {
					groupRegStr = groupRegStr + groupVal
				}
			}
			groupRegStr = groupRegStr + ")$"
			query["hostname"] = bson.RegEx{Pattern: groupRegStr}
		} else {
			if value, reg := handleWildcard(filter.Hostname[0]); reg == true {
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

func prepTagAggr(date int) []bson.M {

	aggr := []bson.M{
		{"$match": bson.M{"date_integer": date}},
		{"$project": bson.M{"_id": 0, "tags": bson.M{"$objectToArray": "$tags"}}},
		{"$unwind": "$tags"},
		{"$replaceRoot": bson.M{"newRoot": "$tags"}},
		{"$project": bson.M{"k": "$k", "v": bson.M{"$split": []string{"$v", ", "}}}},
		{"$unwind": "$v"},
		{"$group": bson.M{"_id": "$k", "v": bson.M{"$addToSet": "$v"}}},
		{"$sort": bson.M{"_id": 1}},
	}

	return aggr

}

func prepGroupEndpointAggr(date int, fgroup fltrGroup, fendpoint fltrEndpoint) []bson.M {

	// prepare match query for applying filtering on the groups based on report
	groupMatch := prepGroupQuery(date, fgroup)
	// prepare match query for applying further filtering on already filtered by groups endpoints based on report
	endpointMatch := prepEndpointQuery(date, fendpoint)

	// Endpoints must be filtered based on the filtered groups they belong to. So we begin an aggregation firstly
	// by filtering groups and then trying to find using lookup which endpoints belong to those filtered groups.
	// The lookup aggregation step does that and creates to reference variables group_date = each group's date_integer value
	// and group_subgroup = each group's subgroup value. For each group those reference variables are used in a matching
	// rule to get only endpoints that their date_integer and group values match group_date and group_subgroup var values
	// respectively. For each group iterated we store the matched endpoint records in an array named endp. Since
	// the pipeline produces a list of groups with nested arrays of matched endpoints we unwind the arrays to get
	// a flat association of each group with each matched endpoint record (again stored in an nested endp field).
	// Lastly we replace the root of each record with the contents of the endp nested field which is the actual
	// endpoint record we want.
	aggr := []bson.M{
		{"$match": groupMatch},
		{"$lookup": bson.M{"from": endpointColName,
			"let": bson.M{
				"group_date":     "$date_integer",
				"group_subgroup": "$subgroup",
			},
			"pipeline": []bson.M{
				{"$match": bson.M{
					"$expr": bson.M{
						"$and": []bson.M{
							{"$eq": []string{"$$group_subgroup", "$group"}},
							{"$eq": []string{"$$group_date", "$date_integer"}},
						},
					},
				},
				},
			},
			"as": "endp"}},
		{"$unwind": "$endp"},
		{"$replaceRoot": bson.M{"newRoot": "$endp"}},
		{"$match": endpointMatch},
	}

	return aggr

}

func prepGroupQuery(date int, filter fltrGroup) bson.M {

	query := bson.M{"date_integer": date}
	// if filter group has values
	if len(filter.Group) > 0 {
		if len(filter.Group) > 1 {
			groupRegStr := "^("
			for i, groupVal := range filter.Group {
				if i != 0 {
					groupRegStr = groupRegStr + "|"
				}
				if value, reg := handleWildcard(groupVal); reg == true {
					groupRegStr = groupRegStr + value
				} else {
					groupRegStr = groupRegStr + groupVal
				}
			}
			groupRegStr = groupRegStr + ")$"
			query["group"] = bson.RegEx{Pattern: groupRegStr}
		} else {
			if value, reg := handleWildcard(filter.Group[0]); reg == true {
				query["group"] = bson.RegEx{Pattern: value}
			} else {
				query["group"] = value
			}
		}
	}

	// if filter group type has values
	if len(filter.GroupType) > 0 {
		if len(filter.GroupType) > 1 {
			groupRegStr := "^("
			for i, groupVal := range filter.GroupType {
				if i != 0 {
					groupRegStr = groupRegStr + "|"
				}
				if value, reg := handleWildcard(groupVal); reg == true {
					groupRegStr = groupRegStr + value
				} else {
					groupRegStr = groupRegStr + groupVal
				}
			}
			groupRegStr = groupRegStr + ")$"
			query["type"] = bson.RegEx{Pattern: groupRegStr}
		} else {
			if value, reg := handleWildcard(filter.GroupType[0]); reg == true {
				query["type"] = bson.RegEx{Pattern: value}
			} else {
				query["type"] = value
			}
		}
	}

	// if filter subgroup has values
	if len(filter.Subgroup) > 0 {
		if len(filter.Subgroup) > 1 {
			groupRegStr := "^("
			for i, groupVal := range filter.Subgroup {
				if i != 0 {
					groupRegStr = groupRegStr + "|"
				}
				if value, reg := handleWildcard(groupVal); reg == true {
					groupRegStr = groupRegStr + value
				} else {
					groupRegStr = groupRegStr + groupVal
				}
			}
			groupRegStr = groupRegStr + ")$"
			query["subgroup"] = bson.RegEx{Pattern: groupRegStr}
		} else {
			if value, reg := handleWildcard(filter.Subgroup[0]); reg == true {
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

// getReportFilters reads a report definition and extracs all related group and endpoint filters
func getReportFilters(r reports.MongoInterface) (fltrGroup, fltrEndpoint) {
	// prepare filter structs for groups and endpoints
	fGroup := fltrGroup{}
	fEndpoint := fltrEndpoint{}
	// prepare empty tag strings for groups and endpoints
	groupTags := ""
	endpointTags := ""
	// help vars to decide when to put commas when building tag strings
	groupFirst := false
	endpointFirst := false

	// Iterate over report filter tags
	for _, tag := range r.Tags {

		if tag.Context == "argo.group.filter.tags" {
			// check if it has special context refering to group filtering based on tags
			// then construct / update group tag string

			if !groupFirst {
				groupTags = groupTags + ","
			} else {
				groupFirst = false
			}
			groupTags = groupTags + tag.Name + ":" + tag.Value

		} else if tag.Context == "argo.group.filter.fields" {
			// check if it has special context refering to group filtering based on basic fields
			// then based on tag name update corresponding field filter

			if tag.Name == "group" {
				fGroup.Group = append(fGroup.Group, tag.Value)
			} else if tag.Name == "type" {
				fGroup.GroupType = append(fGroup.GroupType, tag.Value)
			} else if tag.Name == "subgroup" {
				fGroup.Subgroup = append(fGroup.Subgroup, tag.Value)
			}
		} else if tag.Context == "argo.endpoint.filter.tags" {
			// check if it has special context refering to endpoint filtering based on tags
			// then construct / update endpoint tag string

			if !endpointFirst {
				endpointTags = endpointTags + ","
			} else {
				endpointFirst = false
			}
			endpointTags = endpointTags + tag.Name + ":" + tag.Value
		} else if tag.Context == "argo.endpoint.filter.fields" {
			// check if it has special context refering to endpoint filtering based on basic fields
			// then based on tag name update corresponding field filter

			if tag.Name == "group" {
				fEndpoint.Group = append(fEndpoint.Group, tag.Value)
			} else if tag.Name == "type" {
				fEndpoint.GroupType = append(fEndpoint.GroupType, tag.Value)
			} else if tag.Name == "service" {
				fEndpoint.Service = append(fEndpoint.Service, tag.Value)
			} else if tag.Name == "hostname" {
				fEndpoint.Hostname = append(fEndpoint.Hostname, tag.Value)
			}
		}
	}

	// lastly add constructued tag string to filters
	fGroup.Tags = groupTags
	fEndpoint.Tags = endpointTags

	return fGroup, fEndpoint
}

func getReportGroupType(r reports.MongoInterface) string {

	if r.Topology.Group != nil {
		return r.Topology.Group.Type
	}

	return ""
}

//ListEndpointsByReport lists endpoint topology by report
func ListEndpointsByReport(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	fGroup := fltrGroup{}
	fEndpoint := fltrEndpoint{}

	fGroup, fEndpoint = getReportFilters(report)

	if groupType != "" {
		fGroup.GroupType = append(fGroup.GroupType)
	}

	expDate := getCloseDate(colGroup, dt)

	results := []Endpoint{}

	if expDate < 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	query := prepGroupEndpointAggr(expDate, fGroup, fEndpoint)

	colGroup.Pipe(query).All(&results)

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

	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}
	fGroup := fltrGroup{}
	fGroup, _ = getReportFilters(report)

	if groupType != "" {
		fGroup.GroupType = append(fGroup.GroupType, groupType)
	}

	expDate := getCloseDate(colGroup, dt)

	results := []Group{}

	if expDate < 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	err = colGroup.Find(prepGroupQuery(expDate, fGroup)).All(&results)
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
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	fltr := fltrGroup{}
	fltr.Group = urlValues["group"]
	fltr.GroupType = urlValues["type"]
	fltr.Subgroup = urlValues["subgroup"]
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
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
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
