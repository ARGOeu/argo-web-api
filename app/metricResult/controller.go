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

package metricResult

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/utils/config"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetMultipleMetricResults returns the detailed message from a probe
func GetMultipleMetricResults(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	vars := mux.Vars(r)
	urlValues := r.URL.Query()

	input := metricResultQuery{
		EndpointName: vars["endpoint_name"],
		ExecTime:     urlValues.Get("exec_time"),
		Service:      urlValues.Get("service"),
	}

	filter := urlValues.Get("filter")

	results := []metricResultOutput{}

	col := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("status_metrics")

	// Query the detailed metric results
	cursor, err := col.Aggregate(context.TODO(), prepMultipleQuery(input, filter))

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	output, err = createMultipleMetricResultsView(results, contentType)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

// GetMetricResult returns the detailed message from a probe
func GetMetricResult(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	reportName := urlValues.Get("report")
	reportID := ""

	reportCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")

	if reportName != "" {

		queryResult := reportCol.FindOne(context.TODO(), bson.M{"info.name": reportName})

		if queryResult.Err() != nil {
			if queryResult.Err() == mongo.ErrNoDocuments {
				code = http.StatusNotFound
				message := "The report with the name " + reportName + " does not exist"
				output, err := createErrorMessage(message, code, contentType) //Render the response into XML or JSON
				h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
				return code, h, output, err
			}
			code = http.StatusInternalServerError
			return code, h, output, err
		}

	}

	input := metricResultQuery{
		EndpointName: vars["endpoint_name"],
		MetricName:   vars["metric_name"],
		ExecTime:     urlValues.Get("exec_time"),
		Service:      urlValues.Get("service"),
	}

	results := []metricResultOutput{}

	metricCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("status_metrics")

	// Query the detailed metric results
	cursor, err := metricCol.Find(context.TODO(), prepQuery(input, reportID))

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	if len(results) == 0 {
		code = http.StatusNotFound
		message := "Metric not found!"
		output, err := createErrorMessage(message, code, contentType)
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	output, err = createMultipleMetricResultsView(results, contentType)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

func prepQuery(input metricResultQuery, reportID string) bson.M {

	//Time Related
	const zuluForm = "2006-01-02T15:04:05Z"
	const ymdForm = "20060102"

	ts, _ := time.Parse(zuluForm, input.ExecTime)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))

	// parse time as integer
	tsInt := (ts.Hour() * 10000) + (ts.Minute() * 100) + ts.Second()

	query := bson.M{
		"date_integer": tsYMD,
		"host":         input.EndpointName,
		"metric":       input.MetricName,
		"time_integer": tsInt,
	}

	// filter by service type
	if input.Service != "" {
		query["service"] = input.Service
	}

	if reportID != "" {
		query["report"] = reportID
	}

	return query

}

func prepMultipleQuery(input metricResultQuery, filter string) []bson.M {

	//Time Related
	const zuluForm = "2006-01-02T15:04:05Z"
	const ymdForm = "20060102"

	ts, _ := time.Parse(zuluForm, input.ExecTime)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))

	matchQuery := bson.M{
		"date_integer": tsYMD,
		"host":         input.EndpointName,
	}

	// filter by service type
	if input.Service != "" {
		matchQuery["service"] = input.Service
	}

	// convert to lower case for agility in checks
	filter = strings.ToUpper(filter)

	if filter == "NON-OK" {

		matchQuery["status"] = bson.M{"$ne": "OK"}

	} else if filter == "CRITICAL" ||
		filter == "WARNING" ||
		filter == "OK" ||
		filter == "MISSING" ||
		filter == "UNKNOWN" {

		matchQuery["status"] = filter

	}

	aggrQuery := []bson.M{
		{"$match": matchQuery},
		{"$group": bson.M{
			"_id": bson.M{
				"host":                   "$host",
				"service":                "$service",
				"metric":                 "$metric",
				"timestamp":              "$timestamp",
				"message":                "$message",
				"summary":                "$summary",
				"status":                 "$status",
				"info":                   "$info",
				"actual_data":            "$actual_data",
				"threshold_rule_applied": "$threshold_rule_applied",
				"original_status":        "$original_status"},
		},
		},
		{"$project": bson.M{
			"_id":                    0,
			"host":                   "$_id.host",
			"info":                   "$_id.info",
			"metric":                 "$_id.metric",
			"service":                "$_id.service",
			"timestamp":              "$_id.timestamp",
			"status":                 "$_id.status",
			"summary":                "$_id.summary",
			"message":                "$_id.message",
			"actual_data":            "$_id.actual_data",
			"threshold_rule_applied": "$_id.threshold_rule_applied",
			"original_status":        "$_id.original_status"},
		},
		{"$sort": bson.D{
			{Key: "service", Value: 1},
			{Key: "metric", Value: 1},
			{Key: "timestamp", Value: 1},
		},
		},
	}

	return aggrQuery

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
