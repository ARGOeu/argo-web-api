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
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
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
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	vars := mux.Vars(r)
	urlValues := r.URL.Query()

	input := metricResultQuery{
		EndpointName: vars["endpoint_name"],
		ExecTime:     urlValues.Get("exec_time"),
	}

	filter := urlValues.Get("filter")

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	results := []metricResultOutput{}

	// Query the detailed metric results
	err = mongo.Pipe(session, tenantDbConfig.Db, "status_metrics", prepMultipleQuery(input, filter), &results)

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
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	input := metricResultQuery{
		EndpointName: vars["endpoint_name"],
		MetricName:   vars["metric_name"],
		ExecTime:     urlValues.Get("exec_time"),
	}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	result := metricResultOutput{}

	metricCol := session.DB(tenantDbConfig.Db).C("status_metrics")

	// Query the detailed metric results
	err = metricCol.Find(prepQuery(input)).One(&result)

	output, err = createMetricResultView(result, contentType)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

func prepQuery(input metricResultQuery) bson.M {

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
			{"service", 1},
			{"metric", 1},
			{"timestamp", 1},
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
	h.Set("Allow", fmt.Sprintf("GET, OPTIONS"))
	return code, h, output, err

}
