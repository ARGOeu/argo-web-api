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

package trends

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

const flapEndpointGroups = "flipflop_trends_endpoint_groups"
const flapEndpoints = "flipflop_trends_endpoints"
const flapServices = "flipflop_trends_services"
const flapMetrics = "flipflop_trends_metrics"

func getDateRange(urlValues url.Values) (int, int, error) {
	dateStr := urlValues.Get("date")
	startDateStr := urlValues.Get("start_date")
	endDateStr := urlValues.Get("end_date")

	// if date declared as url parameter use it exclusively as a start and end date
	if dateStr != "" {
		date, _, err := utils.ParseZuluDate(dateStr)
		if err != nil {
			return -1, -1, err
		}
		return date, date, nil
	}

	if startDateStr != "" && endDateStr != "" {
		startDate, _, err := utils.ParseZuluDate(startDateStr)
		if err != nil {
			return -1, -1, err
		}
		endDate, _, err := utils.ParseZuluDate(endDateStr)
		if err != nil {
			return -1, -1, err
		}
		return startDate, endDate, nil
	}

	if (startDateStr == "" && endDateStr != "") || (startDateStr != "" && endDateStr == "") {
		return -1, -1, errors.New("Please use either a date url parameter or a combination of start_date " +
			"and end_date parameters to declare range")
	}

	date, _, err := utils.ParseZuluDate(dateStr)
	return date, date, err

}

// FlatFlappingMetrics returns a list of top flapping metrics for the day
func ListFlappingMetrics(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("List Flapping Metrics")
	err := error(nil)
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Mongo Session
	results := []MetricData{}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	metricCollection := session.DB(tenantDbConfig.Db).C(flapMetrics)

	// Query the detailed metric results
	reportID, err := mongo.GetReportID(session, tenantDbConfig.Db, vars["report_name"])

	startDate, endDate, err := getDateRange(urlValues)

	limit := -1
	limStr := urlValues.Get("top")
	if limStr != "" {
		limit, err = strconv.Atoi(limStr)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	// query for metrics
	filter := bson.M{"report": reportID, "date": bson.M{"$gte": startDate, "$lte": endDate}}

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"group":    "$group",
				"service":  "$service",
				"endpoint": "$endpoint",
				"metric":   "$metric"},
			"flipflop": bson.M{"$sum": "$flipflop"},
		}},

		{"$project": bson.M{
			"group":    "$_id.group",
			"service":  "$_id.service",
			"endpoint": "$_id.endpoint",
			"metric":   "$_id.metric",
			"flipflop": "$flipflop"}},
		{"$sort": bson.D{
			{"flipflop", -1}}}}

	if limit > 0 {
		query = append(query, bson.M{"$limit": limit})

	}

	err = metricCollection.Pipe(query).All(&results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createMetricListView(results, "Success", 200)

	return code, h, output, err
}

// FlatFlappingEndpoints returns a list of top flapping endpoints for the day
func ListFlappingEndpoints(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("List Flapping Endpoints")
	err := error(nil)
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Mongo Session
	results := []EndpointData{}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	endpointCollection := session.DB(tenantDbConfig.Db).C(flapEndpoints)

	// Query the detailed endpoint results
	reportID, err := mongo.GetReportID(session, tenantDbConfig.Db, vars["report_name"])

	startDate, endDate, err := getDateRange(urlValues)

	limit := -1
	limStr := urlValues.Get("top")
	if limStr != "" {
		limit, err = strconv.Atoi(limStr)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	// query for metrics
	filter := bson.M{"report": reportID, "date": bson.M{"$gte": startDate, "$lte": endDate}}

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"group":    "$group",
				"service":  "$service",
				"endpoint": "$endpoint"},
			"flipflop": bson.M{"$sum": "$flipflop"},
		}},

		{"$project": bson.M{
			"group":    "$_id.group",
			"service":  "$_id.service",
			"endpoint": "$_id.endpoint",
			"flipflop": "$flipflop"}},
		{"$sort": bson.D{
			{"flipflop", -1}}}}

	if limit > 0 {
		query = append(query, bson.M{"$limit": limit})

	}

	err = endpointCollection.Pipe(query).All(&results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createEndpointListView(results, "Success", 200)

	return code, h, output, err
}

// FlatFlappingServices returns a list of top flapping services for the day
func ListFlappingServices(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("List Flapping Services")
	err := error(nil)
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Mongo Session
	results := []ServiceData{}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	servicesCollection := session.DB(tenantDbConfig.Db).C(flapServices)

	// Query the detailed service results
	reportID, err := mongo.GetReportID(session, tenantDbConfig.Db, vars["report_name"])

	startDate, endDate, err := getDateRange(urlValues)

	limit := -1
	limStr := urlValues.Get("top")
	if limStr != "" {
		limit, err = strconv.Atoi(limStr)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	// query for metrics
	filter := bson.M{"report": reportID, "date": bson.M{"$gte": startDate, "$lte": endDate}}

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"group":   "$group",
				"service": "$service"},
			"flipflop": bson.M{"$sum": "$flipflop"},
		}},
		{"$project": bson.M{
			"group":    "$_id.group",
			"service":  "$_id.service",
			"flipflop": "$flipflop"}},
		{"$sort": bson.D{
			{"flipflop", -1}}}}

	if limit > 0 {
		query = append(query, bson.M{"$limit": limit})

	}

	err = servicesCollection.Pipe(query).All(&results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createServiceListView(results, "Success", 200)

	return code, h, output, err
}

// FlatFlappingEndpointGroups returns a list of top flapping endpoint groups for the day
func ListFlappingEndpointGroups(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("List Flapping Endpoint Groups")
	err := error(nil)
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Mongo Session
	results := []EndpointGroupData{}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	eGroupsCollection := session.DB(tenantDbConfig.Db).C(flapEndpointGroups)

	// Query the detailed endpoint group results
	reportID, err := mongo.GetReportID(session, tenantDbConfig.Db, vars["report_name"])

	startDate, endDate, err := getDateRange(urlValues)

	limit := -1
	limStr := urlValues.Get("top")
	if limStr != "" {
		limit, err = strconv.Atoi(limStr)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	// query for metrics
	filter := bson.M{"report": reportID, "date": bson.M{"$gte": startDate, "$lte": endDate}}

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"group": "$group"},
			"flipflop": bson.M{"$sum": "$flipflop"},
		}},
		{"$project": bson.M{
			"group":    "$_id.group",
			"flipflop": "$flipflop"}},
		{"$sort": bson.D{
			{"flipflop", -1}}}}

	if limit > 0 {
		query = append(query, bson.M{"$limit": limit})

	}

	err = eGroupsCollection.Pipe(query).All(&results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createEndpointGroupListView(results, "Success", 200)

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
