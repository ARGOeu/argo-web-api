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
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/hbase"
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

	doInfo := urlValues.Get("info")

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

	dataSrc := urlValues.Get("datasource")
	// If hbase bypass mongo session
	if dataSrc == "hbase" {
		// Get hbase configuration
		hbCfg := context.Get(r, "hbase_conf").(config.HbaseConfig)
		// Get tenant name
		tenantName := context.Get(r, "tenant_name").(string)

		// Query Results from hbase
		hbResults, errHb := hbase.QueryStatusEndpoints(hbCfg, tenantName, input.report, strconv.Itoa(input.startTime), input.group, input.service, input.hostname)
		if errHb != nil {
			code = http.StatusInternalServerError
			return code, h, output, errHb
		}
		// Convert hbase results to data output format
		doResults := hbaseToDataOutput(hbResults)
		// Render the reults into xml
		output, err = createView(doResults, input, urlValues.Get("end_time")) //Render the results into JSON/XML format

		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, errHb
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

	if doInfo != "false" {
		err = metricCollection.Pipe(infoAggr(prepareQuery(input, reportID))).All(&results)
	} else {
		err = metricCollection.Find(prepareQuery(input, reportID)).All(&results)

	}

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

	output, err = createView(results, input, urlValues.Get("end_time")) //Render the results into JSON/XML format

	return code, h, output, err
}

func prepareQuery(input InputParams, reportID string) bson.M {

	// prepare the match filter
	filter := bson.M{
		"date_integer":   bson.M{"$gte": input.startTime, "$lte": input.endTime},
		"report":         reportID,
		"endpoint_group": input.group,
		"service":        input.service,
	}

	if len(input.hostname) > 0 {
		filter["host"] = input.hostname
	}

	return filter
}

type list []interface{}

func infoAggr(filter bson.M) []bson.M {

	query := []bson.M{
		{"$match": filter},
		{"$lookup": bson.M{
			"from": "topology_endpoints",
			"let": bson.M{
				"endpoint_date":    "$date_integer",
				"endpoint_name":    "$host",
				"endpoint_service": "$service",
			},
			"pipeline": []bson.M{
				{"$match": bson.M{
					"$expr": bson.M{
						"$and": []bson.M{
							{"$eq": []string{"$$endpoint_date", "$date_integer"}},
							{"$eq": []string{"$$endpoint_name", "$hostname"}},
							{"$eq": []string{"$$endpoint_service", "$service"}},
						},
					},
				}},
				{"$project": bson.M{
					"_id":  0,
					"tags": "$tags"},
				},
			},
			"as": "extra"},
		},
		{"$project": bson.M{
			"date_integer":   1,
			"report":         1,
			"endpoint_group": 1,
			"service":        1,
			"host":           1,
			"status":         1,
			"timestamp":      1,
			"info":           bson.M{"$arrayElemAt": list{"$extra.tags", 0}},
		}},
		// {"$sort": bson.D{
		// 	{"supergroup", 1},
		// 	{"service", 1},
		// 	{"host", 1},
		// 	{"date", 1}}}}
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
