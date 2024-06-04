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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/ARGOeu/argo-web-api/app/feeds"
	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/app/tenants"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/store"
	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// datastore collection name that contains aggregations profile records

const endpointColName = "topology_endpoints"
const groupColName = "topology_groups"
const serviceTypeColName = "topology_service_types"
const feedsDataCol = "feeds_data"

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)
	groupType := gcontext.Get(r, "group_type").(string)
	egroupType := gcontext.Get(r, "endpoint_group_type").(string)

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	reportName := vars["report_name"]
	//Time Related
	dateStr := urlValues.Get("date")

	dt, _, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	// find the report id first
	reportCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	reportID, err := store.GetReportID(reportCol, reportName)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	var serviceResults []string
	var egroupResults []string
	var groupResults []string

	serviceCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("service_ar")
	eGroupCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("endpoint_group_ar")

	distinctServices, err := serviceCol.Distinct(context.TODO(), "name", bson.M{"report": reportID, "date": dt})
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	serviceResults, err = utils.DistinctCast(distinctServices)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	distinctEgroup, err := eGroupCol.Distinct(context.TODO(), "name", bson.M{"report": reportID, "date": dt})
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	egroupResults, err = utils.DistinctCast(distinctEgroup)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	distinctGroup, err := eGroupCol.Distinct(context.TODO(), "supergroup", bson.M{"report": reportID, "date": dt})
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	groupResults, err = utils.DistinctCast(distinctGroup)
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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	colEndpoint := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(endpointColName)
	colGroup := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(groupColName)

	dt, _, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	expDate := store.GetCloseDate(colEndpoint, dt)

	resTagsEndpoint := []TagValues{}
	resTagsGroup := []TagValues{}

	if expDate < 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	cursor, err := colEndpoint.Aggregate(context.TODO(), prepTagAggr(expDate))

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &resTagsEndpoint)

	expDate = store.GetCloseDate(colGroup, dt)

	if expDate < 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	cursor, err = colGroup.Aggregate(context.TODO(), prepTagAggr(expDate))

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &resTagsGroup)

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
	mode := urlValues.Get("mode")

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	dt, _, err := utils.ParseZuluDate(dateStr)
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

	results, expDate, err := getEndpointResults(cfg.MongoClient, tenantDbConfig, dt, fltr)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if mode == "combined" {

		// check for feeds
		dbConfigs := (getComboDBConfigs(cfg.MongoClient, tenantDbConfig, cfg))
		for _, dbConfig := range dbConfigs {
			// append subresults to list of combined results

			subResults, _, err := getEndpointResults(cfg.MongoClient, dbConfig.Config, dt, fltr)

			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}
			// tag results with tenant name
			for i := range subResults {
				subResults[i].Tenant = dbConfig.Tenant
			}
			results = append(results, subResults...)
		}
	}

	// check if nothing found
	if expDate < 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
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

// getEndpointResults accepts an date in integer format YYYYMMDD, a tenand db configuration,
// a filter endpoint object and returns relevant topology endpoints
func getEndpointResults(client *mongo.Client, dbConfig config.MongoConfig, dateInt int, filterE fltrEndpoint) ([]Endpoint, int, error) {
	subResults := []Endpoint{}

	colEndpoint := client.Database(dbConfig.Db).Collection(endpointColName)
	expDate := store.GetCloseDate(colEndpoint, dateInt)
	if expDate < 0 {
		return subResults, expDate, nil
	}

	cursor, err := colEndpoint.Find(context.TODO(), prepEndpointQuery(expDate, filterE))

	if err != nil {
		return subResults, expDate, err
	}

	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &subResults)

	return subResults, expDate, err
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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)
	urlValues := r.URL.Query()
	dateStr := urlValues.Get("date")

	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	incoming := []Endpoint{}
	// Try ingest request body
	body, err := io.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))

	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// check if topology already exists for current day

	existing := Endpoint{}
	endpointCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(endpointColName)
	err = endpointCol.FindOne(context.TODO(), bson.M{"date_integer": dt}).Decode(&existing)
	if err != nil {
		// Stop at any error except not found. We want to have not found
		if err != mongo.ErrNoDocuments {
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

	incomingInf := make([]interface{}, len(incoming))
	for i, value := range incoming {
		incomingInf[i] = value
	}

	_, err = endpointCol.InsertMany(context.TODO(), incomingInf)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)
	urlValues := r.URL.Query()

	dateStr := urlValues.Get("date")
	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	endpointCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(endpointColName)
	deleteResult, err := endpointCol.DeleteMany(context.TODO(), bson.M{"date_integer": dt})
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if deleteResult.DeletedCount == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	// Create view of the results
	output, err = createMessageOUT(fmt.Sprintf("Topology of %d endpoints deleted for date: %s", deleteResult.DeletedCount, dateStr), 200, "json")
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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)
	urlValues := r.URL.Query()
	dateStr := urlValues.Get("date")
	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	incoming := []Group{}
	// Try ingest request body
	body, err := io.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))

	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// check if topology already exists for current day

	existing := Group{}
	groupCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(groupColName)
	err = groupCol.FindOne(context.TODO(), bson.M{"date_integer": dt}).Decode(&existing)
	if err != nil {
		// Stop at any error except not found. We want to have not found
		if err != mongo.ErrNoDocuments {
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

	incomingInf := make([]interface{}, len(incoming))
	for i, value := range incoming {
		incomingInf[i] = value
	}

	_, err = groupCol.InsertMany(context.TODO(), incomingInf)

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

		if _, exists := tagmap[kvsplit[0]]; !exists {
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
				if value, reg := handleWildcard(tagVal); reg {
					tagRegStr = tagRegStr + value
				} else {
					tagRegStr = tagRegStr + tagVal
				}
			}
			tagRegStr = tagRegStr + ")$"
			query["tags."+tag] = primitive.Regex{Pattern: tagRegStr}

		} else {

			// check if tag filter value has exclude operator and trim it
			trimValue, exclude := handleExclude(tagmap[tag][0])

			if value, reg := handleWildcard(trimValue); reg {
				if exclude {
					query["tags."+tag] = bson.M{"$not": primitive.Regex{Pattern: value}}
				} else {
					query["tags."+tag] = primitive.Regex{Pattern: value}
				}

			} else {
				if exclude {
					query["tags."+tag] = bson.M{"$ne": value}
				} else {
					query["tags."+tag] = value
				}
			}
		}
	}

	return query

}

// check if a filter value has a negative operator for exclusion
func handleExclude(value string) (string, bool) {

	if strings.HasPrefix(value, "~") {
		return value[1:], true
	} else {
		return value, false
	}
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
				if value, reg := handleWildcard(groupVal); reg {
					groupRegStr = groupRegStr + value
				} else {
					groupRegStr = groupRegStr + groupVal
				}
			}
			groupRegStr = groupRegStr + ")$"
			query["group"] = primitive.Regex{Pattern: groupRegStr}
		} else {

			// check if value has exclude prefix and trim it
			trimValue, exclude := handleExclude(filter.Group[0])

			if value, reg := handleWildcard(trimValue); reg {
				if exclude {
					query["group"] = bson.M{"$not": primitive.Regex{Pattern: value}}
				} else {
					query["group"] = primitive.Regex{Pattern: value}
				}

			} else {
				if exclude {
					query["group"] = bson.M{"$ne": value}
				} else {
					query["group"] = value
				}

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
				if value, reg := handleWildcard(groupVal); reg {
					groupRegStr = groupRegStr + value
				} else {
					groupRegStr = groupRegStr + groupVal
				}
			}
			groupRegStr = groupRegStr + ")$"
			query["type"] = primitive.Regex{Pattern: groupRegStr}
		} else {

			// check if value has exclude prefix and trim it
			trimValue, exclude := handleExclude(filter.GroupType[0])

			if value, reg := handleWildcard(trimValue); reg {
				if exclude {
					query["type"] = bson.M{"$not": primitive.Regex{Pattern: value}}
				} else {
					query["type"] = primitive.Regex{Pattern: value}
				}
			} else {
				if exclude {
					query["type"] = bson.M{"$ne": value}
				} else {
					query["type"] = value
				}
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
				if value, reg := handleWildcard(groupVal); reg {
					groupRegStr = groupRegStr + value
				} else {
					groupRegStr = groupRegStr + groupVal
				}
			}
			groupRegStr = groupRegStr + ")$"
			query["service"] = primitive.Regex{Pattern: groupRegStr}
		} else {

			// check if value has exclude prefix and trim it
			trimValue, exclude := handleExclude(filter.Service[0])

			if value, reg := handleWildcard(trimValue); reg {
				if exclude {
					query["service"] = bson.M{"$not": primitive.Regex{Pattern: value}}
				} else {
					query["service"] = primitive.Regex{Pattern: value}
				}
			} else {
				if exclude {
					query["service"] = bson.M{"$ne": value}
				} else {
					query["service"] = value
				}
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
				if value, reg := handleWildcard(groupVal); reg {
					groupRegStr = groupRegStr + value
				} else {
					groupRegStr = groupRegStr + groupVal
				}
			}
			groupRegStr = groupRegStr + ")$"
			query["hostname"] = primitive.Regex{Pattern: groupRegStr}
		} else {

			// check if value has exclude prefix and trim it
			trimValue, exclude := handleExclude(filter.Hostname[0])

			if value, reg := handleWildcard(trimValue); reg {
				if exclude {
					query["hostname"] = bson.M{"$not": primitive.Regex{Pattern: value}}
				} else {
					query["hostname"] = primitive.Regex{Pattern: value}
				}
			} else {
				if exclude {
					query["hostname"] = bson.M{"$ne": value}
				} else {
					query["hostname"] = value
				}
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
				if value, reg := handleWildcard(groupVal); reg {
					groupRegStr = groupRegStr + value
				} else {
					groupRegStr = groupRegStr + groupVal
				}
			}
			groupRegStr = groupRegStr + ")$"
			query["group"] = primitive.Regex{Pattern: groupRegStr}
		} else {

			// check if value has exclude prefix and trim it
			trimValue, exclude := handleExclude(filter.Group[0])

			if value, reg := handleWildcard(trimValue); reg {
				if exclude {
					query["group"] = bson.M{"$not": primitive.Regex{Pattern: value}}
				} else {
					query["group"] = primitive.Regex{Pattern: value}
				}
			} else {
				if exclude {
					query["group"] = bson.M{"$ne": value}
				} else {
					query["group"] = value
				}
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
				if value, reg := handleWildcard(groupVal); reg {
					groupRegStr = groupRegStr + value
				} else {
					groupRegStr = groupRegStr + groupVal
				}
			}
			groupRegStr = groupRegStr + ")$"
			query["type"] = primitive.Regex{Pattern: groupRegStr}
		} else {

			// check if value has exclude prefix and trim it
			trimValue, exclude := handleExclude(filter.GroupType[0])

			if value, reg := handleWildcard(trimValue); reg {
				if exclude {
					query["type"] = bson.M{"$not": primitive.Regex{Pattern: value}}
				} else {
					query["type"] = primitive.Regex{Pattern: value}
				}
			} else {
				if exclude {
					query["type"] = bson.M{"$ne": value}
				} else {
					query["type"] = value
				}
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
				if value, reg := handleWildcard(groupVal); reg {
					groupRegStr = groupRegStr + value
				} else {
					groupRegStr = groupRegStr + groupVal
				}
			}
			groupRegStr = groupRegStr + ")$"
			query["subgroup"] = primitive.Regex{Pattern: groupRegStr}
		} else {

			// check if value has exclude prefix and trim it
			trimValue, exclude := handleExclude(filter.Subgroup[0])

			if value, reg := handleWildcard(trimValue); reg {
				if exclude {
					query["subgroup"] = bson.M{"$not": primitive.Regex{Pattern: value}}
				} else {
					query["subgroup"] = primitive.Regex{Pattern: value}
				}
			} else {
				if exclude {
					query["subgroup"] = bson.M{"$ne": value}
				} else {
					query["subgroup"] = value
				}
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

func getReportTagsArray(tag reports.Tag) string {

	// split multiple values using commas as delimeter
	vals := strings.Split(tag.Value, ",")
	// trim values from space
	result := ""
	for i := range vals {
		val := strings.TrimSpace(vals[i])
		result = result + tag.Name + ":*" + val + "*"
		if i < len(vals)-1 {
			result = result + ","
		}
	}

	return result

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

		if tag.Context == "argo.group.filter.tags.array" {
			// check if it has special context refering to group filtering based on tags
			// then construct / update group tag string

			if !groupFirst {
				groupTags = groupTags + ","
			} else {
				groupFirst = false
			}
			groupTags = groupTags + getReportTagsArray(tag)

		}
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
		} else if tag.Context == "argo.endpoint.filter.tags.array" {
			// check if it has special context refering to endpoint filtering based on tags
			// then construct / update endpoint tag string

			if !endpointFirst {
				endpointTags = endpointTags + ","
			} else {
				endpointFirst = false
			}
			endpointTags = endpointTags + getReportTagsArray(tag)
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

// ListEndpointsByReport lists endpoint topology by report
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
	mode := urlValues.Get("mode")
	reportName := vars["report"]

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	colReports := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	//get the report

	report := reports.MongoInterface{}
	err = colReports.FindOne(context.TODO(), bson.M{"info.name": reportName}).Decode(&report)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			output, err = createMessageOUT(fmt.Sprintf("No report with name: %s exists!", reportName), 404, "json")
			code = 404
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	groupType := getReportGroupType(report)
	egroupType := getReportEndpointGroupType(report)

	dt, _, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	fGroup, fEndpoint := getReportFilters(report)

	if groupType != "" && (len(fGroup.GroupType) == 0) {
		fGroup.GroupType = append(fGroup.GroupType, groupType)
	}

	if egroupType != "" && (len(fEndpoint.GroupType) == 0) {
		fEndpoint.GroupType = append(fEndpoint.GroupType, egroupType)
	}

	results, _, err := getGroupEndpointResults(cfg.MongoClient, tenantDbConfig, dt, fGroup, fEndpoint)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if mode == "combined" {
		// check for feeds
		dbConfigs := (getComboDBConfigs(cfg.MongoClient, tenantDbConfig, cfg))
		for _, dbConfig := range dbConfigs {
			// append subresults to list of combined results

			subResults, _, err := getGroupEndpointResults(cfg.MongoClient, dbConfig.Config, dt, fGroup, fEndpoint)
			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}

			// tag results with tenant name
			for i := range subResults {
				subResults[i].Tenant = dbConfig.Tenant
			}
			results = append(results, subResults...)
		}
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

// getGroupEndpointResults accepts an date in integer format YYYYMMDD, a tenand db configuration,
// a filter group object and a filter endpoint object and returns relevant topology endpoints
func getGroupEndpointResults(client *mongo.Client, dbConfig config.MongoConfig, dateInt int, filterG fltrGroup, filterE fltrEndpoint) ([]Endpoint, int, error) {
	subResults := []Endpoint{}

	colGroup := client.Database(dbConfig.Db).Collection(groupColName)
	expDate := store.GetCloseDate(colGroup, dateInt)
	if expDate < 0 {
		return subResults, expDate, nil
	}

	query := prepGroupEndpointAggr(expDate, filterG, filterE)

	cursor, err := colGroup.Aggregate(context.TODO(), query)

	if err != nil {
		return subResults, expDate, err
	}

	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &subResults)

	return subResults, expDate, err
}

// ListGroupsByReport lists group topology by report
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
	mode := urlValues.Get("mode")
	reportName := vars["report"]

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	colReports := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	//get the report

	report := reports.MongoInterface{}
	err = colReports.FindOne(context.TODO(), bson.M{"info.name": reportName}).Decode(&report)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			output, err = createMessageOUT(fmt.Sprintf("No report with name: %s exists!", reportName), 404, "json")
			code = 404
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	groupType := getReportGroupType(report)

	dt, _, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	fGroup := fltrGroup{}
	fGroup, _ = getReportFilters(report)

	if groupType != "" && (len(fGroup.GroupType) == 0) {
		fGroup.GroupType = append(fGroup.GroupType, groupType)
	}

	results, _, err := getGroupResults(cfg.MongoClient, tenantDbConfig, dt, fGroup)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if mode == "combined" {
		// check for feeds
		dbConfigs := (getComboDBConfigs(cfg.MongoClient, tenantDbConfig, cfg))
		for _, dbConfig := range dbConfigs {
			// append subresults to list of combined results

			subResults, _, err := getGroupResults(cfg.MongoClient, dbConfig.Config, dt, fGroup)
			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}
			// tag results with tenant name
			for i := range subResults {
				subResults[i].Tenant = dbConfig.Tenant
			}
			results = append(results, subResults...)
		}
	}

	// check if nothing found
	if len(results) == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
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

// CreateServiceTypes creates a list of service types for a specific date
func CreateServiceTypes(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	incoming := []ServiceType{}
	// Try ingest request body
	body, err := io.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))

	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// check if topology already exists for current day

	existing := ServiceType{}
	serviceTypeCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(serviceTypeColName)
	err = serviceTypeCol.FindOne(context.TODO(), bson.M{"date_integer": dt}).Decode(&existing)
	if err != nil {
		// Stop at any error except not found. We want to have not found
		if err != mongo.ErrNoDocuments {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
		// else continue correctly -
	} else {
		// If found we need to inform user that the topology is already created for this date
		output, err = createMessageOUT(fmt.Sprintf("Topology list of service types already exists for date: %s, please either update it or delete it first!", dateStr), 409, "json")
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

	incomingInf := make([]interface{}, len(incoming))
	for i, value := range incoming {
		incomingInf[i] = value
	}

	_, err = serviceTypeCol.InsertMany(context.TODO(), incomingInf)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createMessageOUT(fmt.Sprintf("Topology of %d service types created for date: %s", len(incoming), dateStr), 201, "json") //Render the results into JSON
	code = 201
	return code, h, output, err
}

// ListServiceTypes by date
func ListServiceTypes(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	mode := urlValues.Get("mode")

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	dt, _, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	results, err := getServiceTypeResults(cfg.MongoClient, tenantDbConfig, dt)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if mode == "combined" {

		// check for feeds

		dbConfigs := (getComboDBConfigs(cfg.MongoClient, tenantDbConfig, cfg))
		for _, dbConfig := range dbConfigs {
			// append subresults to list of combined results
			subResults, err := getServiceTypeResults(cfg.MongoClient, dbConfig.Config, dt)
			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}
			// tag results with tenant name
			for i := range subResults {
				subResults[i].Tenant = dbConfig.Tenant
			}

			results = append(results, subResults...)
		}
	}

	// check if nothing found
	if len(results) == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	// Create view of the results
	output, err = createListService(results, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// getServiceTypeResults accepts an date in integer format YYYYMMDD and a tenand db configuration
// and returns relevant service types
func getServiceTypeResults(client *mongo.Client, dbConfig config.MongoConfig, dateInt int) ([]ServiceType, error) {
	subResults := []ServiceType{}

	subServiceTypeCol := client.Database(dbConfig.Db).Collection(serviceTypeColName)
	expDate := store.GetCloseDate(subServiceTypeCol, dateInt)
	if expDate < 0 {
		return subResults, nil
	}

	cursor, err := subServiceTypeCol.Find(context.TODO(), bson.M{"date_integer": expDate})

	if err != nil {
		return subResults, err
	}

	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &subResults)

	return subResults, err
}

// getComboDBConfigs returns a list of database configuration objects based on a list of tenant names
func getComboDBConfigs(client *mongo.Client, comboTenantCfg config.MongoConfig, cfg config.Config) []TenantDB {
	// empty result
	result := []TenantDB{}

	feeds := feeds.Data{}
	ftCol := client.Database(comboTenantCfg.Db).Collection(feedsDataCol)
	err := ftCol.FindOne(context.TODO(), bson.M{}).Decode(&feeds)
	// if feeds exist

	if err == nil && feeds.Tenants != nil {

		tenantsCol := client.Database(cfg.MongoDB.Db).Collection("tenants")
		tenantResults := []tenants.Tenant{}
		// grab tenants involved in feeds
		cursor, err := tenantsCol.Find(context.TODO(), bson.M{"id": bson.M{"$in": feeds.Tenants}})
		if err != nil {
			// emtpy result
			return result
		}

		defer cursor.Close(context.TODO())

		cursor.All(context.TODO(), &tenantResults)

		for _, item := range tenantResults {

			if len(item.DbConf) == 0 {
				continue
			}
			itemConf := item.DbConf[0]

			result = append(result, TenantDB{
				Tenant: item.Info.Name,
				Config: config.MongoConfig{
					Db:       itemConf.Database,
					Store:    itemConf.Store,
					Host:     itemConf.Server,
					Port:     itemConf.Port,
					Username: itemConf.Username,
					Password: itemConf.Password,
				}})
		}

	}

	return result
}

// DeleteServiceTypes deletes a list of service types in topology for a specific date
func DeleteServiceTypes(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	serviceCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(serviceTypeColName)
	deleteResult, err := serviceCol.DeleteMany(context.TODO(), bson.M{"date_integer": dt})
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if deleteResult.DeletedCount == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	// Create view of the results
	output, err = createMessageOUT(fmt.Sprintf("Topology of %d service types deleted for date: %s", deleteResult.DeletedCount, dateStr), 200, "json")
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
	mode := urlValues.Get("mode")

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	dt, _, err := utils.ParseZuluDate(dateStr)
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

	results, expDate, err := getGroupResults(cfg.MongoClient, tenantDbConfig, dt, fltr)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if mode == "combined" {
		// check for feeds
		dbConfigs := (getComboDBConfigs(cfg.MongoClient, tenantDbConfig, cfg))
		for _, dbConfig := range dbConfigs {
			// append subresults to list of combined results
			subResults, _, err := getGroupResults(cfg.MongoClient, dbConfig.Config, dt, fltr)

			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}
			// tag results with tenant name
			for i := range subResults {
				subResults[i].Tenant = dbConfig.Tenant
			}
			results = append(results, subResults...)
		}
	}

	// check if nothing found
	if expDate < 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
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

// getGroupResults accepts an date in integer format YYYYMMDD, a tenand db configuration
// and a filter group object and returns relevant topology groups
func getGroupResults(client *mongo.Client, dbConfig config.MongoConfig, dateInt int, fltr fltrGroup) ([]Group, int, error) {
	subResults := []Group{}

	subGroupsCol := client.Database(dbConfig.Db).Collection(groupColName)
	expDate := store.GetCloseDate(subGroupsCol, dateInt)
	if expDate < 0 {
		return subResults, expDate, nil
	}
	cursor, err := subGroupsCol.Find(context.TODO(), prepGroupQuery(expDate, fltr))

	if err != nil {
		return subResults, expDate, nil
	}

	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &subResults)

	return subResults, expDate, err
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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)
	urlValues := r.URL.Query()
	dateStr := urlValues.Get("date")
	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	groupCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(groupColName)
	deleteResult, err := groupCol.DeleteMany(context.TODO(), bson.M{"date_integer": dt})
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if deleteResult.DeletedCount == 0 {
		output, _ = respond.MarshalContent(respond.ErrNotFoundQuery, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	// Create view of the results
	output, err = createMessageOUT(fmt.Sprintf("Topology of %d groups deleted for date: %s", deleteResult.DeletedCount, dateStr), 200, "json")
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
	h.Set("Allow", "GET, OPTIONS")
	return code, h, output, err

}
