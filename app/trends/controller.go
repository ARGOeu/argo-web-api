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
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/store"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

const flapEndpointGroups = "flipflop_trends_endpoint_groups"
const flapEndpoints = "flipflop_trends_endpoints"
const flapServices = "flipflop_trends_services"
const flapMetrics = "flipflop_trends_metrics"
const statusMetrics = "status_trends_metrics"
const statusEndpoints = "status_trends_endpoints"
const statusServices = "status_trends_services"
const statusEgroups = "status_trends_groups"

type list []interface{}

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
		return -1, -1, errors.New("please use either a date url parameter or a combination of start_date " +
			"and end_date parameters to declare range")
	}

	date, _, err := utils.ParseZuluDate(dateStr)
	return date, date, err

}

// ListStatusMetrics returns a list of top status metrics for the day
func ListStatusMetrics(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	metricCollection := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(statusMetrics)
	// Query the detailed metric results
	reportID, err := store.GetReportID(rCol, vars["report_name"])
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	startDate, endDate, err := getDateRange(urlValues)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	limit := -1
	limStr := urlValues.Get("top")
	if limStr != "" {
		limit, err = strconv.Atoi(limStr)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	granularity := urlValues.Get("granularity")

	// query for metrics
	filter := bson.M{"report": reportID, "date": bson.M{"$gte": startDate, "$lte": endDate}}

	// apply query for bucketed monthly results if granularity is set to monthly
	if granularity == "monthly" {

		results := []StatusMonthMetricData{}

		query := []bson.M{
			{"$match": filter},
			{"$group": bson.M{
				"_id": bson.M{
					"month":    bson.M{"$substr": list{"$date", 0, 6}},
					"group":    "$group",
					"service":  "$service",
					"endpoint": "$endpoint",
					"metric":   "$metric",
					"status":   "$status"},
				"events": bson.M{"$sum": "$trends"},
			}},
			{"$sort": bson.D{{Key: "_id.month", Value: 1}, {Key: "_id.status", Value: 1}, {Key: "events", Value: -1}}},
			{
				"$group": bson.M{
					"_id": bson.M{"month": "$_id.month", "status": "$_id.status"},
					"top": bson.M{"$push": bson.M{"group": "$_id.group", "service": "$_id.service", "endpoint": "$_id.endpoint", "metric": "$_id.metric", "status": "$_id.status", "events": "$events"}}}},
		}

		// trim down the list in each month-bucket according to the limit parameter
		if limit > 0 {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"status": "$_id.status",
				"top":    bson.M{"$slice": list{"$top", limit}}}})
		} else {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"status": "$_id.status",
				"top":    "$top"}})
		}

		// sort end results by month bucket ascending
		query = append(query, bson.M{"$sort": bson.D{{Key: "date", Value: 1}, {Key: "status", Value: 1}}})

		cursor, err := metricCollection.Aggregate(context.TODO(), query)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &results)

		output, err = createStatusMonthMetricListView(results, "Success", 200)

		return code, h, output, err

	}

	// continue by calculating non monthly bucketed results
	results := []StatusGroupMetricData{}

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"group":    "$group",
				"service":  "$service",
				"endpoint": "$endpoint",
				"metric":   "$metric",
				"status":   "$status"},
			"events": bson.M{"$sum": "$trends"},
		}},
		{"$sort": bson.D{{Key: "_id.status", Value: 1}, {Key: "events", Value: -1}}},
		{
			"$group": bson.M{
				"_id": bson.M{"status": "$_id.status"},
				"top": bson.M{"$push": bson.M{"group": "$_id.group", "service": "$_id.service", "endpoint": "$_id.endpoint", "metric": "$_id.metric", "status": "$_id.status", "events": "$events"}}}},
	}

	// trim down the list in each month-bucket according to the limit parameter
	if limit > 0 {
		query = append(query, bson.M{"$project": bson.M{"status": "$_id.status",
			"top": bson.M{"$slice": list{"$top", limit}}}})
	} else {
		query = append(query, bson.M{"$project": bson.M{"status": "$_id.status",
			"top": "$top"}})
	}

	// sort end results by month bucket ascending
	query = append(query, bson.M{"$sort": bson.D{{Key: "status", Value: 1}}})

	cursor, err := metricCollection.Aggregate(context.TODO(), query)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	output, err = createStatusMetricListView(results, "Success", 200)

	return code, h, output, err
}

// ListStatusMetricsTags returns a list of top status metrics for the day sorted by tags
func ListStatusMetricsTags(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	metricCollection := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(statusMetrics)
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	// Query the detailed metric results
	reportID, err := store.GetReportID(rCol, vars["report_name"])
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	startDate, endDate, err := getDateRange(urlValues)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	limit := -1
	limStr := urlValues.Get("top")
	if limStr != "" {
		limit, err = strconv.Atoi(limStr)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	granularity := urlValues.Get("granularity")

	// query for metrics
	filter := bson.M{"report": reportID, "tags": bson.M{"$exists": true}, "date": bson.M{"$gte": startDate, "$lte": endDate}}

	// apply query for bucketed monthly results if granularity is set to monthly
	if granularity == "monthly" {

		results := []StatusMonthMetricData{}

		query := []bson.M{
			{"$match": filter},
			{"$unwind": "$tags"},
			{"$group": bson.M{
				"_id": bson.M{
					"month":    bson.M{"$substr": list{"$date", 0, 6}},
					"tag":      "$tags",
					"group":    "$group",
					"service":  "$service",
					"endpoint": "$endpoint",
					"metric":   "$metric",
					"status":   "$status"},
				"events": bson.M{"$sum": "$trends"},
			}},
			{"$sort": bson.D{{Key: "_id.month", Value: 1}, {Key: "_id.status", Value: 1}, {Key: "_id.tag", Value: 1}, {Key: "events", Value: -1}}},
			{
				"$group": bson.M{
					"_id": bson.M{"month": "$_id.month", "status": "$_id.status", "tag": "$_id.tag"},
					"top": bson.M{"$push": bson.M{"group": "$_id.group", "service": "$_id.service", "endpoint": "$_id.endpoint", "metric": "$_id.metric", "status": "$_id.status", "events": "$events"}}}},
		}

		// trim down the list in each month-bucket according to the limit parameter
		if limit > 0 {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"status": "$_id.status",
				"tag":    "$_id.tag",
				"top":    bson.M{"$slice": list{"$top", limit}}}})
		} else {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"status": "$_id.status",
				"tag":    "$_id.tag",
				"top":    "$top"}})
		}

		// sort end results by month bucket ascending
		query = append(query, bson.M{"$sort": bson.D{{Key: "date", Value: 1}, {Key: "status", Value: 1}, {Key: "tag", Value: 1}}})

		cursor, err := metricCollection.Aggregate(context.TODO(), query)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &results)

		output, err = createStatusMonthMetricListView(results, "Success", 200)

		return code, h, output, err

	}

	// continue by calculating non monthly bucketed results
	results := []StatusGroupMetricData{}

	query := []bson.M{
		{"$match": filter},
		{"$unwind": "$tags"},
		{"$group": bson.M{
			"_id": bson.M{
				"tag":      "$tags",
				"group":    "$group",
				"service":  "$service",
				"endpoint": "$endpoint",
				"metric":   "$metric",
				"status":   "$status"},
			"events": bson.M{"$sum": "$trends"},
		}},
		{"$sort": bson.D{{Key: "_id.status", Value: 1}, {Key: "_id.tag", Value: 1}, {Key: "events", Value: -1}}},
		{
			"$group": bson.M{
				"_id": bson.M{"status": "$_id.status", "tag": "$_id.tag"},
				"top": bson.M{"$push": bson.M{"group": "$_id.group", "service": "$_id.service", "endpoint": "$_id.endpoint", "metric": "$_id.metric", "status": "$_id.status", "events": "$events"}}}},
	}

	// trim down the list in each month-bucket according to the limit parameter
	if limit > 0 {
		query = append(query, bson.M{"$project": bson.M{"status": "$_id.status", "tag": "$_id.tag",
			"top": bson.M{"$slice": list{"$top", limit}}}})
	} else {
		query = append(query, bson.M{"$project": bson.M{"status": "$_id.status", "tag": "$_id.tag",
			"top": "$top"}})
	}

	// sort end results by month bucket ascending
	query = append(query, bson.M{"$sort": bson.D{{Key: "status", Value: 1}}})

	cursor, err := metricCollection.Aggregate(context.TODO(), query)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	output, err = createStatusMetricListView(results, "Success", 200)

	return code, h, output, err
}

// ListStatusEndpoints returns a list of top status endpoints (in duration) for the day
func ListStatusEndpoints(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	endpointCollection := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(statusEndpoints)
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")

	// Query the detailed status endpoints trend results
	reportID, err := store.GetReportID(rCol, vars["report_name"])
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	startDate, endDate, err := getDateRange(urlValues)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	limit := -1
	limStr := urlValues.Get("top")
	if limStr != "" {
		limit, err = strconv.Atoi(limStr)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	granularity := urlValues.Get("granularity")

	// query for endpoints
	filter := bson.M{"report": reportID, "date": bson.M{"$gte": startDate, "$lte": endDate}}

	// apply query for bucketed monthly results if granularity is set to monthly
	if granularity == "monthly" {

		results := []StatusMonthEndpointData{}

		query := []bson.M{
			{"$match": filter},
			{"$group": bson.M{
				"_id": bson.M{
					"month":    bson.M{"$substr": list{"$date", 0, 6}},
					"group":    "$group",
					"service":  "$service",
					"endpoint": "$endpoint",
					"status":   "$status"},
				"duration": bson.M{"$sum": "$duration"},
			}},
			{"$sort": bson.D{{Key: "_id.month", Value: 1}, {Key: "_id.status", Value: 1}, {Key: "duration", Value: -1}}},
			{
				"$group": bson.M{
					"_id": bson.M{"month": "$_id.month", "status": "$_id.status"},
					"top": bson.M{"$push": bson.M{"group": "$_id.group", "service": "$_id.service", "endpoint": "$_id.endpoint", "status": "$_id.status", "duration": "$duration"}}}},
		}

		// trim down the list in each month-bucket according to the limit parameter
		if limit > 0 {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"status": "$_id.status",
				"top":    bson.M{"$slice": list{"$top", limit}}}})
		} else {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"status": "$_id.status",
				"top":    "$top"}})
		}

		// sort end results by month bucket ascending
		query = append(query, bson.M{"$sort": bson.D{{Key: "date", Value: 1}, {Key: "status", Value: 1}}})

		cursor, err := endpointCollection.Aggregate(context.TODO(), query)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &results)

		output, err = createStatusMonthEndpointListView(results, "Success", 200)

		return code, h, output, err

	}

	// continue by calculating non monthly bucketed results
	results := []StatusGroupEndpointData{}

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"group":    "$group",
				"service":  "$service",
				"endpoint": "$endpoint",
				"status":   "$status"},
			"duration": bson.M{"$sum": "$duration"},
		}},
		{"$sort": bson.D{{Key: "_id.status", Value: 1}, {Key: "duration", Value: -1}}},
		{
			"$group": bson.M{
				"_id": bson.M{"status": "$_id.status"},
				"top": bson.M{"$push": bson.M{"group": "$_id.group", "service": "$_id.service", "endpoint": "$_id.endpoint", "status": "$_id.status", "duration": "$duration"}}}},
	}

	// trim down the list in each month-bucket according to the limit parameter
	if limit > 0 {
		query = append(query, bson.M{"$project": bson.M{"status": "$_id.status",
			"top": bson.M{"$slice": list{"$top", limit}}}})
	} else {
		query = append(query, bson.M{"$project": bson.M{"status": "$_id.status",
			"top": "$top"}})
	}

	// sort end results by month bucket ascending
	query = append(query, bson.M{"$sort": bson.D{{Key: "status", Value: 1}}})

	cursor, err := endpointCollection.Aggregate(context.TODO(), query)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	output, err = createStatusEndpointListView(results, "Success", 200)

	return code, h, output, err
}

// ListStatusServices returns a list of top status services (in duration) for the day
func ListStatusServices(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	serviceCollection := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(statusServices)
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	// Query the detailed status services trend results
	reportID, err := store.GetReportID(rCol, vars["report_name"])
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	startDate, endDate, err := getDateRange(urlValues)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	limit := -1
	limStr := urlValues.Get("top")
	if limStr != "" {
		limit, err = strconv.Atoi(limStr)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	granularity := urlValues.Get("granularity")

	// query for services
	filter := bson.M{"report": reportID, "date": bson.M{"$gte": startDate, "$lte": endDate}}

	// apply query for bucketed monthly results if granularity is set to monthly
	if granularity == "monthly" {

		results := []StatusMonthServiceData{}

		query := []bson.M{
			{"$match": filter},
			{"$group": bson.M{
				"_id": bson.M{
					"month":   bson.M{"$substr": list{"$date", 0, 6}},
					"group":   "$group",
					"service": "$service",
					"status":  "$status"},
				"duration": bson.M{"$sum": "$duration"},
			}},
			{"$sort": bson.D{{Key: "_id.month", Value: 1}, {Key: "_id.status", Value: 1}, {Key: "duration", Value: -1}}},
			{
				"$group": bson.M{
					"_id": bson.M{"month": "$_id.month", "status": "$_id.status"},
					"top": bson.M{"$push": bson.M{"group": "$_id.group", "service": "$_id.service", "status": "$_id.status", "duration": "$duration"}}}},
		}

		// trim down the list in each month-bucket according to the limit parameter
		if limit > 0 {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"status": "$_id.status",
				"top":    bson.M{"$slice": list{"$top", limit}}}})
		} else {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"status": "$_id.status",
				"top":    "$top"}})
		}

		// sort end results by month bucket ascending
		query = append(query, bson.M{"$sort": bson.D{{Key: "date", Value: 1}, {Key: "status", Value: 1}}})

		cursor, err := serviceCollection.Aggregate(context.TODO(), query)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &results)

		output, err = createStatusMonthServiceListView(results, "Success", 200)

		return code, h, output, err

	}

	// continue by calculating non monthly bucketed results
	results := []StatusGroupServiceData{}

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"group":   "$group",
				"service": "$service",
				"status":  "$status"},
			"duration": bson.M{"$sum": "$duration"},
		}},
		{"$sort": bson.D{{Key: "_id.status", Value: 1}, {Key: "duration", Value: -1}}},
		{
			"$group": bson.M{
				"_id": bson.M{"status": "$_id.status"},
				"top": bson.M{"$push": bson.M{"group": "$_id.group", "service": "$_id.service", "status": "$_id.status", "duration": "$duration"}}}},
	}

	// trim down the list in each month-bucket according to the limit parameter
	if limit > 0 {
		query = append(query, bson.M{"$project": bson.M{"status": "$_id.status",
			"top": bson.M{"$slice": list{"$top", limit}}}})
	} else {
		query = append(query, bson.M{"$project": bson.M{"status": "$_id.status",
			"top": "$top"}})
	}

	// sort end results by month bucket ascending
	query = append(query, bson.M{"$sort": bson.D{{Key: "status", Value: 1}}})

	cursor, err := serviceCollection.Aggregate(context.TODO(), query)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	output, err = createStatusServiceListView(results, "Success", 200)

	return code, h, output, err
}

// ListStatusEgroups returns a list of top status endpoint groups (in duration) for the day
func ListStatusEgroups(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	egroupCollection := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(statusEgroups)
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	// Query the detailed status services trend results
	reportID, err := store.GetReportID(rCol, vars["report_name"])
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	startDate, endDate, err := getDateRange(urlValues)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	limit := -1
	limStr := urlValues.Get("top")
	if limStr != "" {
		limit, err = strconv.Atoi(limStr)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	granularity := urlValues.Get("granularity")

	// query for endpoint groups
	filter := bson.M{"report": reportID, "date": bson.M{"$gte": startDate, "$lte": endDate}}

	// apply query for bucketed monthly results if granularity is set to monthly
	if granularity == "monthly" {

		results := []StatusMonthEgroupData{}

		query := []bson.M{
			{"$match": filter},
			{"$group": bson.M{
				"_id": bson.M{
					"month":  bson.M{"$substr": list{"$date", 0, 6}},
					"group":  "$group",
					"status": "$status"},
				"duration": bson.M{"$sum": "$duration"},
			}},
			{"$sort": bson.D{{Key: "_id.month", Value: 1}, {Key: "_id.status", Value: 1}, {Key: "duration", Value: -1}}},
			{
				"$group": bson.M{
					"_id": bson.M{"month": "$_id.month", "status": "$_id.status"},
					"top": bson.M{"$push": bson.M{"group": "$_id.group", "status": "$_id.status", "duration": "$duration"}}}},
		}

		// trim down the list in each month-bucket according to the limit parameter
		if limit > 0 {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"status": "$_id.status",
				"top":    bson.M{"$slice": list{"$top", limit}}}})
		} else {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"status": "$_id.status",
				"top":    "$top"}})
		}

		// sort end results by month bucket ascending
		query = append(query, bson.M{"$sort": bson.D{{Key: "date", Value: 1}, {Key: "status", Value: 1}}})

		cursor, err := egroupCollection.Aggregate(context.TODO(), query)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &results)

		output, err = createStatusMonthEgroupListView(results, "Success", 200)

		return code, h, output, err

	}

	// continue by calculating non monthly bucketed results
	results := []StatusGroupEgroupData{}

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"group":  "$group",
				"status": "$status"},
			"duration": bson.M{"$sum": "$duration"},
		}},
		{"$sort": bson.D{{Key: "_id.status", Value: 1}, {Key: "duration", Value: -1}}},
		{
			"$group": bson.M{
				"_id": bson.M{"status": "$_id.status"},
				"top": bson.M{"$push": bson.M{"group": "$_id.group", "status": "$_id.status", "duration": "$duration"}}}},
	}

	// trim down the list in each month-bucket according to the limit parameter
	if limit > 0 {
		query = append(query, bson.M{"$project": bson.M{"status": "$_id.status",
			"top": bson.M{"$slice": list{"$top", limit}}}})
	} else {
		query = append(query, bson.M{"$project": bson.M{"status": "$_id.status",
			"top": "$top"}})
	}

	// sort end results by month bucket ascending
	query = append(query, bson.M{"$sort": bson.D{{Key: "status", Value: 1}}})

	cursor, err := egroupCollection.Aggregate(context.TODO(), query)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	output, err = createStatusEgroupListView(results, "Success", 200)

	return code, h, output, err
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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	metricCollection := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(flapMetrics)

	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	// Query the detailed metric results
	reportID, err := store.GetReportID(rCol, vars["report_name"])
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	startDate, endDate, err := getDateRange(urlValues)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	limit := -1
	limStr := urlValues.Get("top")
	if limStr != "" {
		limit, err = strconv.Atoi(limStr)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	granularity := urlValues.Get("granularity")

	// query for metrics
	filter := bson.M{"report": reportID, "date": bson.M{"$gte": startDate, "$lte": endDate}}

	// apply query for bucketed monthly results if granularity is set to monthly
	if granularity == "monthly" {

		results := []MonthMetricData{}

		query := []bson.M{
			{"$match": filter},
			{"$group": bson.M{
				"_id": bson.M{
					"month":    bson.M{"$substr": list{"$date", 0, 6}},
					"group":    "$group",
					"service":  "$service",
					"endpoint": "$endpoint",
					"metric":   "$metric"},
				"flipflop": bson.M{"$sum": "$flipflop"},
			}},
			{"$sort": bson.D{{Key: "_id.month", Value: 1}, {Key: "flipflop", Value: -1}}},
			{
				"$group": bson.M{
					"_id": bson.M{"month": "$_id.month"},
					"top": bson.M{"$push": bson.M{"group": "$_id.group", "service": "$_id.service", "endpoint": "$_id.endpoint", "metric": "$_id.metric", "flipflop": "$flipflop"}}}},
		}

		// trim down the list in each month-bucket according to the limit parameter
		if limit > 0 {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"top": bson.M{"$slice": list{"$top", limit}}}})
		} else {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"top": "$top"}})
		}

		// sort end results by month bucket ascending
		query = append(query, bson.M{"$sort": bson.D{{Key: "date", Value: 1}}})

		cursor, err := metricCollection.Aggregate(context.TODO(), query)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &results)

		output, err = createMonthMetricListView(results, "Success", 200)

		return code, h, output, err

	}

	// continue by calculating non monthly bucketed results
	results := []MetricData{}

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
			{Key: "flipflop", Value: -1}}}}

	if limit > 0 {
		query = append(query, bson.M{"$limit": limit})

	}

	cursor, err := metricCollection.Aggregate(context.TODO(), query)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	output, err = createMetricListView(results, "Success", 200)

	return code, h, output, err
}

// FlatFlappingMetricsTags returns a list of top flapping metrics for the day per tag
func ListFlappingMetricsTags(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	metricCollection := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(flapMetrics)

	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	// Query the detailed metric results
	reportID, err := store.GetReportID(rCol, vars["report_name"])
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	startDate, endDate, err := getDateRange(urlValues)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	limit := -1
	limStr := urlValues.Get("top")
	if limStr != "" {
		limit, err = strconv.Atoi(limStr)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	granularity := urlValues.Get("granularity")

	// query only for metrics that have tags
	filter := bson.M{"report": reportID, "tags": bson.M{"$exists": true}, "date": bson.M{"$gte": startDate, "$lte": endDate}}

	// apply query for bucketed monthly results if granularity is set to monthly
	if granularity == "monthly" {

		results := []MonthMetricData{}

		query := []bson.M{
			{"$match": filter},
			{"$unwind": "$tags"},
			{"$group": bson.M{
				"_id": bson.M{
					"month":    bson.M{"$substr": list{"$date", 0, 6}},
					"tag":      "$tags",
					"group":    "$group",
					"service":  "$service",
					"endpoint": "$endpoint",
					"metric":   "$metric"},
				"flipflop": bson.M{"$sum": "$flipflop"},
			}},
			{"$sort": bson.D{{Key: "_id.month", Value: 1}, {Key: "_id.tag", Value: 1}, {Key: "flipflop", Value: -1}}},
			{
				"$group": bson.M{
					"_id": bson.M{"month": "$_id.month", "tag": "$_id.tag"},
					"top": bson.M{"$push": bson.M{"group": "$_id.group", "service": "$_id.service", "endpoint": "$_id.endpoint", "metric": "$_id.metric", "flipflop": "$flipflop"}}}},
		}

		// trim down the list in each month-bucket according to the limit parameter
		if limit > 0 {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"tag": "$_id.tag",
				"top": bson.M{"$slice": list{"$top", limit}}}})
		} else {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"tag": "$_id.tag",
				"top": "$top"}})
		}

		// sort end results by month bucket ascending
		query = append(query, bson.M{"$sort": bson.D{{Key: "date", Value: 1}, {Key: "tag", Value: 1}}})

		cursor, err := metricCollection.Aggregate(context.TODO(), query)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &results)
		output, err = createMonthMetricListView(results, "Success", 200)

		return code, h, output, err

	}

	// continue by calculating non monthly bucketed results
	results := []TagMetricData{}

	query := []bson.M{
		{"$match": filter},
		{"$unwind": "$tags"},
		{"$group": bson.M{
			"_id": bson.M{
				"tag":      "$tags",
				"group":    "$group",
				"service":  "$service",
				"endpoint": "$endpoint",
				"metric":   "$metric"},
			"flipflop": bson.M{"$sum": "$flipflop"},
		}},
		{"$sort": bson.D{{Key: "_id.tag", Value: 1}, {Key: "flipflop", Value: -1}}},
		{
			"$group": bson.M{
				"_id": bson.M{"tag": "$_id.tag"},
				"top": bson.M{"$push": bson.M{"group": "$_id.group", "service": "$_id.service", "endpoint": "$_id.endpoint", "metric": "$_id.metric", "flipflop": "$flipflop"}}}},
	}

	// trim down the list in each month-bucket according to the limit parameter
	if limit > 0 {
		query = append(query, bson.M{"$project": bson.M{
			"tag": "$_id.tag",
			"top": bson.M{"$slice": list{"$top", limit}}}})
	} else {
		query = append(query, bson.M{"$project": bson.M{
			"tag": "$_id.tag",
			"top": "$top"}})
	}

	// sort end results by tag
	query = append(query, bson.M{"$sort": bson.D{{Key: "tag", Value: 1}}})

	cursor, err := metricCollection.Aggregate(context.TODO(), query)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	output, err = createTagMetricDataListView(results, "Success", 200)

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	// Mongo Session
	results := []EndpointData{}

	endpointCollection := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(flapEndpoints)
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	// Query the detailed endpoint results
	reportID, err := store.GetReportID(rCol, vars["report_name"])
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	startDate, endDate, err := getDateRange(urlValues)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	limit := -1
	limStr := urlValues.Get("top")
	if limStr != "" {
		limit, err = strconv.Atoi(limStr)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	granularity := urlValues.Get("granularity")

	// query for metrics
	filter := bson.M{"report": reportID, "date": bson.M{"$gte": startDate, "$lte": endDate}}

	// apply query for bucketed monthly results if granularity is set to monthly
	if granularity == "monthly" {

		results := []MonthEndpointData{}

		query := []bson.M{
			{"$match": filter},
			{"$group": bson.M{
				"_id": bson.M{
					"month":    bson.M{"$substr": list{"$date", 0, 6}},
					"group":    "$group",
					"service":  "$service",
					"endpoint": "$endpoint"},
				"flipflop": bson.M{"$sum": "$flipflop"},
			}},
			{"$sort": bson.D{{Key: "_id.month", Value: 1}, {Key: "flipflop", Value: -1}}},
			{
				"$group": bson.M{
					"_id": bson.M{"month": "$_id.month"},
					"top": bson.M{"$push": bson.M{"group": "$_id.group", "service": "$_id.service", "endpoint": "$_id.endpoint", "flipflop": "$flipflop"}}}},
		}

		// trim down the list in each month-bucket according to the limit parameter
		if limit > 0 {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"top": bson.M{"$slice": list{"$top", limit}}}})
		} else {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"top": "$top"}})
		}

		// sort end results by month bucket ascending
		query = append(query, bson.M{"$sort": bson.D{{Key: "date", Value: 1}}})

		cursor, err := endpointCollection.Aggregate(context.TODO(), query)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &results)
		output, err = createMonthEndpointListView(results, "Success", 200)

		return code, h, output, err

	}

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
			{Key: "flipflop", Value: -1}}}}

	if limit > 0 {
		query = append(query, bson.M{"$limit": limit})

	}

	cursor, err := endpointCollection.Aggregate(context.TODO(), query)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	// Mongo Session
	results := []ServiceData{}

	servicesCollection := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(flapServices)
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	// Query the detailed service results
	reportID, err := store.GetReportID(rCol, vars["report_name"])
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	startDate, endDate, err := getDateRange(urlValues)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	limit := -1
	limStr := urlValues.Get("top")
	if limStr != "" {
		limit, err = strconv.Atoi(limStr)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	granularity := urlValues.Get("granularity")

	// query for metrics
	filter := bson.M{"report": reportID, "date": bson.M{"$gte": startDate, "$lte": endDate}}

	// apply query for bucketed monthly results if granularity is set to monthly
	if granularity == "monthly" {

		results := []MonthServiceData{}

		query := []bson.M{
			{"$match": filter},
			{"$group": bson.M{
				"_id": bson.M{
					"month":   bson.M{"$substr": list{"$date", 0, 6}},
					"group":   "$group",
					"service": "$service"},
				"flipflop": bson.M{"$sum": "$flipflop"},
			}},
			{"$sort": bson.D{{Key: "_id.month", Value: 1}, {Key: "flipflop", Value: -1}}},
			{
				"$group": bson.M{
					"_id": bson.M{"month": "$_id.month"},
					"top": bson.M{"$push": bson.M{"group": "$_id.group", "service": "$_id.service", "flipflop": "$flipflop"}}}},
		}

		// trim down the list in each month-bucket according to the limit parameter
		if limit > 0 {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"top": bson.M{"$slice": list{"$top", limit}}}})
		} else {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"top": "$top"}})
		}

		// sort end results by month bucket ascending
		query = append(query, bson.M{"$sort": bson.D{{"date", 1}}})

		cursor, err := servicesCollection.Aggregate(context.TODO(), query)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &results)

		output, err = createMonthServiceListView(results, "Success", 200)

		return code, h, output, err

	}

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
			{Key: "flipflop", Value: -1}}}}

	if limit > 0 {
		query = append(query, bson.M{"$limit": limit})

	}

	cursor, err := servicesCollection.Aggregate(context.TODO(), query)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	// Mongo Session
	results := []EndpointGroupData{}

	eGroupsCollection := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(flapEndpointGroups)
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	// Query the detailed endpoint group results
	reportID, err := store.GetReportID(rCol, vars["report_name"])
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

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

	granularity := urlValues.Get("granularity")

	// query for metrics
	filter := bson.M{"report": reportID, "date": bson.M{"$gte": startDate, "$lte": endDate}}

	// apply query for bucketed monthly results if granularity is set to monthly
	if granularity == "monthly" {

		results := []MonthEndpointGroupData{}

		query := []bson.M{
			{"$match": filter},
			{"$group": bson.M{
				"_id": bson.M{
					"month": bson.M{"$substr": list{"$date", 0, 6}},
					"group": "$group"},
				"flipflop": bson.M{"$sum": "$flipflop"},
			}},
			{"$sort": bson.D{{Key: "_id.month", Value: 1}, {Key: "flipflop", Value: -1}}},
			{
				"$group": bson.M{
					"_id": bson.M{"month": "$_id.month"},
					"top": bson.M{"$push": bson.M{"group": "$_id.group", "flipflop": "$flipflop"}}}},
		}

		// trim down the list in each month-bucket according to the limit parameter
		if limit > 0 {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"top": bson.M{"$slice": list{"$top", limit}}}})
		} else {
			query = append(query, bson.M{"$project": bson.M{"date": bson.M{"$concat": list{bson.M{"$substr": list{"$_id.month", 0, 4}},
				"-", bson.M{"$substr": list{"$_id.month", 4, 6}}}},
				"top": "$top"}})
		}

		// sort end results by month bucket ascending
		query = append(query, bson.M{"$sort": bson.D{{Key: "date", Value: 1}}})

		cursor, err := eGroupsCollection.Aggregate(context.TODO(), query)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &results)

		output, err = createMonthEndpointGroupListView(results, "Success", 200)

		return code, h, output, err

	}

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
			{Key: "flipflop", Value: -1}}}}

	if limit > 0 {
		query = append(query, bson.M{"$limit": limit})

	}

	cursor, err := eGroupsCollection.Aggregate(context.TODO(), query)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

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
	h.Set("Allow", "GET, OPTIONS")
	return code, h, output, err

}
