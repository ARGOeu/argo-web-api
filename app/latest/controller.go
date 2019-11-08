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

package latest

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

// GetMetricResult returns the detailed message from a probe
func ListLatestResults(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	endpointGroup := vars["group_name"]
	reportName := vars["report_name"]
	limit := urlValues.Get("limit")
	dateStr := urlValues.Get("date")
	filter := urlValues.Get("filter")

	strict := false
	if urlValues.Get("strict") == "true" {
		strict = true
	}

	lim, err := strconv.Atoi(limit)
	if err != nil {
		lim = 500
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

	resultItems := []MetricData{}
	result := MetricList{}

	metricCol := session.DB(tenantDbConfig.Db).C("status_metrics")

	// Query the detailed metric results
	metricResultsQuery := prepQuery(dateStr, reportID, endpointGroup, filter, strict, lim)
	if strict {
		err = metricCol.Pipe(metricResultsQuery).All(&resultItems)
	} else {
		if lim == -1 {
			err = metricCol.Find(metricResultsQuery).Sort("-time_integer").All(&resultItems)
		} else {
			err = metricCol.Find(metricResultsQuery).Sort("-time_integer").Limit(lim).All(&resultItems)
		}
	}

	result.MetricData = resultItems

	output, err = createLatestView(result, contentType, http.StatusOK)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

func prepQuery(dateStr string, report string, group string, filter string, strict bool, limit int) interface{} {

	//Time Related
	const zuluForm = "2006-01-02T15:04:05Z"
	const ymdForm = "20060102"

	ts := time.Now().UTC()
	if dateStr != "" {
		ts, _ = time.Parse(zuluForm, dateStr)
	}

	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))

	filter = strings.ToUpper(filter)
	query := bson.M{
		"report":       report,
		"date_integer": tsYMD,
	}

	if group != "" {
		query["endpoint_group"] = group
	}

	if filter == "NON-OK" {

		query["status"] = bson.M{"$ne": "OK"}

	} else if filter == "CRITICAL" ||
		filter == "WARNING" ||
		filter == "OK" ||
		filter == "MISSING" ||
		filter == "UNKNOWN" {

		query["status"] = filter

	}

	if strict {

		// start building the aggregation pipeline
		// the first step in aggregation is the match query from above
		// then we will request a reverse sort of results in time but
		// with a limit following so to optimize the sorting in the last n elements of the day

		pipelineQuery := []bson.M{
			{
				"$match": query,
			},
			{
				"$sort": bson.M{
					"time_integer": -1,
				},
			},
			{
				"$limit": limit,
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"endpoint_group": "$endpoint_group",
						"host":           "$host",
						"service":        "$service",
						"metric":         "$metric",
					},
					// since now we are applying reverse sorting we need to select the first element (last in reverse = first)
					"endpoint_group": bson.M{"$first": "$endpoint_group"},
					"service":        bson.M{"$first": "$service"},
					"host":           bson.M{"$first": "$host"},
					"metric":         bson.M{"$first": "$metric"},
					"timestamp":      bson.M{"$first": "$timestamp"},
					"status":         bson.M{"$first": "$status"},
					"message":        bson.M{"$first": "$message"},
					"summary":        bson.M{"$first": "$summary"},
					"time_integer":   bson.M{"$first": "$time_integer"},
				},
			},
			{
				"$sort": bson.M{
					"time_integer": -1,
				},
			},
		}

		return pipelineQuery
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
