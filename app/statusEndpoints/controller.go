/*
 * Copyright (c) 2015 GRNET S.A.
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

package statusEndpoints

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// parseZuluDate is used to parse a zulu formatted date to integer
func parseZuluDate(dateStr string) (int, error) {
	parsedTime, _ := time.Parse(zuluForm, dateStr)
	return strconv.Atoi(parsedTime.Format(ymdForm))
}

// getPrevDay returns the previous day
func getPrevDay(dateStr string) (int, error) {
	parsedTime, _ := time.Parse(zuluForm, dateStr)
	prevTime := parsedTime.AddDate(0, 0, -1)
	return strconv.Atoi(prevTime.Format(ymdForm))
}

// FlatListEndpointTimelines returns a list of metric timelines
func FlatListEndpointTimelines(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("List Endpoint Timelines")
	err := error(nil)
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	skip := 0
	tkStr := urlValues.Get("nextPageToken")
	if tkStr != "" {
		tk, err := base64.StdEncoding.DecodeString(tkStr)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
		skip, err = strconv.Atoi(string(tk))
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	limit := -1
	limStr := urlValues.Get("pageSize")
	if limStr != "" {
		limit, err = strconv.Atoi(limStr)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	parsedStart, _ := parseZuluDate(urlValues.Get("start_time"))
	parsedEnd, _ := parseZuluDate(urlValues.Get("end_time"))

	input := InputParams{
		parsedStart,
		parsedEnd,
		vars["report_name"],
		vars["group_type"],
		vars["group_name"],
		vars["service_name"],
		vars["endpoint_name"],
		contentType,
	}

	// This is going to be used to determine a detailed view or not of the results
	view := urlValues.Get("view")
	details := false
	if view == "details" {
		details = true
	}

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Mongo Session
	results := []DataOutput{}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	metricCollection := session.DB(tenantDbConfig.Db).C("status_endpoints")

	// Query the detailed metric results
	reportID, err := mongo.GetReportID(session, tenantDbConfig.Db, input.report)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	err = metricCollection.Pipe(prepareFlatQuery(input, reportID, limit, skip)).All(&results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	parsedPrev, _ := getPrevDay(urlValues.Get("start_time"))

	//if no status results yet show previous days results
	if len(results) == 0 {
		// Zero query results
		input.startTime = parsedPrev
		err = metricCollection.Find(prepareQuery(input, reportID)).All(&results)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	output, err = createView(results, input, urlValues.Get("end_time"), details) //Render the results into JSON/XML format

	return code, h, output, err
}

// ListEndpointTimelines returns a list of metric timelines
func ListEndpointTimelines(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("List Endpoint Timelines")
	err := error(nil)
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	parsedStart, _ := parseZuluDate(urlValues.Get("start_time"))
	parsedEnd, _ := parseZuluDate(urlValues.Get("end_time"))

	input := InputParams{
		parsedStart,
		parsedEnd,
		vars["report_name"],
		vars["group_type"],
		vars["group_name"],
		vars["service_name"],
		vars["endpoint_name"],
		contentType,
	}

	// This is going to be used to determine a detailed view or not of the results
	view := urlValues.Get("view")
	details := false
	if view == "details" {
		details = true
	}

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Mongo Session
	results := []DataOutput{}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	metricCollection := session.DB(tenantDbConfig.Db).C("status_endpoints")

	// Query the detailed metric results
	reportID, err := mongo.GetReportID(session, tenantDbConfig.Db, input.report)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	err = metricCollection.Find(prepareQuery(input, reportID)).All(&results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	parsedPrev, _ := getPrevDay(urlValues.Get("start_time"))

	//if no status results yet show previous days results
	if len(results) == 0 {
		// Zero query results
		input.startTime = parsedPrev
		err = metricCollection.Find(prepareQuery(input, reportID)).All(&results)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	output, err = createView(results, input, urlValues.Get("end_time"), details) //Render the results into JSON/XML format

	return code, h, output, err
}

func prepareQuery(input InputParams, reportID string) bson.M {

	// prepare the match filter
	filter := bson.M{
		"date_integer":   bson.M{"$gte": input.startTime, "$lte": input.endTime},
		"report":         reportID,
		"endpoint_group": input.group,
	}

	if len(input.service) > 0 {
		filter["service"] = input.service
	}

	if len(input.hostname) > 0 {
		filter["host"] = input.hostname
	}

	return filter
}

func prepareFlatQuery(input InputParams, reportID string, limit int, skip int) []bson.M {

	// prepare the match filter
	match := bson.M{
		"date_integer":   bson.M{"$gte": input.startTime, "$lte": input.endTime},
		"report":         reportID,
		"endpoint_group": input.group,
		"service":        input.service,
	}

	if len(input.hostname) > 0 {
		match["host"] = input.hostname
	}

	query := []bson.M{
		bson.M{"$match": match},
	}

	if limit > 0 {
		query = append(query, bson.M{"$skip": skip})
		query = append(query, bson.M{"$limit": limit + 1})

	}

	return query
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
