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

package results

import (
	"fmt"
	"net/http"

	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// ListServiceFlavorResults is responsible for handling request to list service flavor results
func ListServiceFlavorResults(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
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

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	report := reports.MongoInterface{}
	err = mongo.FindOne(session, tenantDbConfig.Db, "reports", bson.M{"info.name": vars["report_name"]}, &report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report_name"] + " does not exist"
		output, err := createErrorMessage(message, code, contentType) //Render the response into XML or JSON
		return code, h, output, err
	}

	input := serviceFlavorResultQuery{
		basicQuery: basicQuery{
			Name:        vars["service_type"],
			Granularity: urlValues.Get("granularity"),
			Format:      contentType,
			StartTime:   urlValues.Get("start_time"),
			EndTime:     urlValues.Get("end_time"),
			Report:      report,
			Vars:        vars,
		},
		EndpointGroup: vars["lgroup_name"],
	}

	tenantDB := session.DB(tenantDbConfig.Db)
	errs := input.Validate(tenantDB)
	if len(errs) > 0 {
		out := respond.BadRequestSimple
		out.Errors = errs
		output = out.MarshalTo(contentType)
		code = 400
		return code, h, output, err
	}

	if vars["lgroup_type"] != report.GetEndpointGroupType() {
		code = http.StatusNotFound
		message := "The report " + vars["report_name"] + " does not define endpoint group type: " + vars["lgroup_type"] + ". Try using " + report.GetEndpointGroupType() + " instead."
		output, err := createErrorMessage(message, code, contentType) //Render the response into XML or JSON
		return code, h, output, err
	}

	results := []ServiceFlavorInterface{}

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Construct the query to mongodb based on the input
	filter := bson.M{
		"date":   bson.M{"$gte": input.StartTimeInt, "$lte": input.EndTimeInt},
		"report": report.ID,
	}

	if input.Name != "" {
		filter["name"] = input.Name
	}

	if input.EndpointGroup != "" {
		filter["supergroup"] = input.EndpointGroup
	}

	// Select the granularity of the search daily/monthly
	if input.Granularity == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		query := DailyServiceFlavor(filter)
		err = mongo.Pipe(session, tenantDbConfig.Db, "service_ar", query, &results)
	} else if input.Granularity == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query := MonthlyServiceFlavor(filter)
		err = mongo.Pipe(session, tenantDbConfig.Db, "service_ar", query, &results)
	}

	// mongo.Find(session, tenantDbConfig.Db, "endpoint_group_ar", bson.M{}, "_id", &results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if len(results) == 0 {
		code = http.StatusNotFound
		message := "No results found for given query"
		output, err = createErrorMessage(message, code, contentType)
		return code, h, output, err
	}

	output, err = createServiceFlavorResultView(results, report, input.Format)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err

}

// ListEndpointGroupResults endpoint group availabilities according to the http request
func ListEndpointGroupResults(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	report := reports.MongoInterface{}
	err = mongo.FindOne(session, tenantDbConfig.Db, "reports", bson.M{"info.name": vars["report_name"]}, &report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report_name"] + " does not exist"
		output, err := createErrorMessage(message, code, contentType) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	input := endpointGroupResultQuery{
		basicQuery{
			Name:        vars["lgroup_name"],
			Granularity: urlValues.Get("granularity"),
			Format:      contentType,
			StartTime:   urlValues.Get("start_time"),
			EndTime:     urlValues.Get("end_time"),
			Report:      report,
			Vars:        vars,
		}, "",
	}

	tenantDB := session.DB(tenantDbConfig.Db)
	errs := input.Validate(tenantDB)
	if len(errs) > 0 {
		out := respond.BadRequestSimple
		out.Errors = errs
		output = out.MarshalTo(contentType)
		code = 400
		return code, h, output, err
	}

	if vars["lgroup_type"] != report.GetEndpointGroupType() {
		code = http.StatusNotFound
		message := "The report " + vars["report_name"] + " does not define endpoint group type: " + vars["lgroup_type"] + ". Try using " + report.GetEndpointGroupType() + " instead."
		output, err := createErrorMessage(message, code, contentType) //Render the response into XML or JSON
		return code, h, output, err
	}

	results := []EndpointGroupInterface{}

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Construct the query to mongodb based on the input
	filter := bson.M{
		"date":   bson.M{"$gte": input.StartTimeInt, "$lte": input.EndTimeInt},
		"report": report.ID,
	}

	if input.Name != "" {
		filter["name"] = input.Name
	}

	// Select the granularity of the search daily/monthly
	if input.Granularity == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		query := DailyEndpointGroup(filter)
		err = mongo.Pipe(session, tenantDbConfig.Db, "endpoint_group_ar", query, &results)
	} else if input.Granularity == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query := MonthlyEndpointGroup(filter)
		err = mongo.Pipe(session, tenantDbConfig.Db, "endpoint_group_ar", query, &results)
	}

	// mongo.Find(session, tenantDbConfig.Db, "endpoint_group_ar", bson.M{}, "_id", &results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if len(results) == 0 {
		code = http.StatusNotFound
		message := "No results found for given query"
		output, err = createErrorMessage(message, code, contentType)
		return code, h, output, err
	}

	output, err = createEndpointGroupResultView(results, report, input.Format)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

// ListSuperGroupResults supergroup availabilities according to the http request
func ListSuperGroupResults(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	report := reports.MongoInterface{}
	err = mongo.FindOne(session, tenantDbConfig.Db, "reports", bson.M{"info.name": vars["report_name"]}, &report)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	input := endpointGroupResultQuery{
		basicQuery{
			Name:        vars["group_name"],
			Granularity: urlValues.Get("granularity"),
			Format:      contentType,
			StartTime:   urlValues.Get("start_time"),
			EndTime:     urlValues.Get("end_time"),
			Report:      report,
			Vars:        vars,
		}, "",
	}

	tenantDB := session.DB(tenantDbConfig.Db)
	errs := input.Validate(tenantDB)
	if len(errs) > 0 {
		out := respond.BadRequestSimple
		out.Errors = errs
		output = out.MarshalTo(contentType)
		code = 400
		return code, h, output, err
	}

	results := []SuperGroupInterface{}

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Construct the query to mongodb based on the input
	filter := bson.M{
		"date":   bson.M{"$gte": input.StartTimeInt, "$lte": input.EndTimeInt},
		"report": report.ID,
	}

	if input.Name != "" {
		filter["supergroup"] = input.Name
	}

	// Select the granularity of the search daily/monthly
	if input.Granularity == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		query := DailySuperGroup(filter)
		err = mongo.Pipe(session, tenantDbConfig.Db, "endpoint_group_ar", query, &results)
	} else if input.Granularity == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query := MonthlySuperGroup(filter)
		err = mongo.Pipe(session, tenantDbConfig.Db, "endpoint_group_ar", query, &results)
	}
	// mongo.Find(session, tenantDbConfig.Db, "endpoint_group_ar", bson.M{}, "_id", &results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if len(results) == 0 {
		code = http.StatusNotFound
		message := "No results found for given query"
		output, err = createErrorMessage(message, code, contentType)
		return code, h, output, err
	}

	output, err = createSuperGroupView(results, report, input.Format)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

// DailyServiceFlavor query to aggregate daily SF results from mongoDB
func DailyServiceFlavor(filter bson.M) []bson.M {
	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"date":         bson.D{{"$substr", list{"$date", 0, 8}}},
				"name":         "$name",
				"supergroup":   "$supergroup",
				"availability": "$availability",
				"reliability":  "$reliability",
				"unknown":      "$unknown",
				"up":           "$up",
				"down":         "$down",
				"report":       "$report"}}},
		{"$project": bson.M{
			"date":         "$_id.date",
			"name":         "$_id.name",
			"availability": "$_id.availability",
			"reliability":  "$_id.reliability",
			"unknown":      "$_id.unknown",
			"up":           "$_id.up",
			"down":         "$_id.down",
			"supergroup":   "$_id.supergroup",
			"report":       "$_id.report"}},
		{"$sort": bson.D{
			{"supergroup", 1},
			{"name", 1},
			{"date", 1}}}}

	return query
}

// MonthlyServiceFlavor query to aggregate daily SF results from mongoDB
func MonthlyServiceFlavor(filter bson.M) []bson.M {
	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"date":       bson.D{{"$substr", list{"$date", 0, 6}}},
				"name":       "$name",
				"supergroup": "$supergroup",
				"report":     "$report"},
			"avgup":      bson.M{"$avg": "$up"},
			"avgunknown": bson.M{"$avg": "$unknown"},
			"avgdown":    bson.M{"$avg": "$down"}}},
		{"$project": bson.M{
			"date":       "$_id.date",
			"name":       "$_id.name",
			"supergroup": "$_id.supergroup",
			"report":     "$_id.report",
			"unknown":    "$avgunknown",
			"up":         "$avgup",
			"down":       "$avgdown",
			"availability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{1.00000001, "$avgunknown"}}}},
					100}},
			"reliability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgunknown"}}, "$avgdown"}}}},
					100}}}},
		{"$sort": bson.D{
			{"supergroup", 1},
			{"name", 1},
			{"date", 1}}}}
	return query
}

// DailyEndpointGroup query to aggregate daily results from mongodb
func DailyEndpointGroup(filter bson.M) []bson.M {
	// Mongo aggregation pipeline
	// Select all the records that match q
	// Project to select just the first 8 digits of the date YYYYMMDD
	// Sort by profile->supergroup->endpointGroup->datetime
	query := []bson.M{
		{"$match": filter},
		{"$project": bson.M{
			"date":         bson.M{"$substr": list{"$date", 0, 8}},
			"availability": 1,
			"reliability":  1,
			"unknown":      1,
			"up":           1,
			"down":         1,
			"report":       1,
			"supergroup":   1,
			"name":         1}},
		{"$sort": bson.D{
			{"report", 1},
			{"supergroup", 1},
			{"name", 1},
			{"date", 1}}}}

	return query
}

// MonthlyEndpointGroup query to aggregate monthly results from mongodb
func MonthlyEndpointGroup(filter bson.M) []bson.M {

	// Mongo aggregation pipeline
	// Select all the records that match q
	// Group them by the first six digits of their date (YYYYMM), their supergroup, their endpointGroup, their profile, etc...
	// from that group find the average of the uptime, u, downtime
	// Project the result to a better format and do this computation
	// availability = (avgup/(1.00000001 - avgu))*100
	// reliability = (avgup/((1.00000001 - avgu)-avgd))*100
	// Sort the results by namespace->profile->supergroup->endpointGroup->datetime

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"date":       bson.M{"$substr": list{"$date", 0, 6}},
				"name":       "$name",
				"supergroup": "$supergroup",
				"report":     "$report"},
			"avgup":      bson.M{"$avg": "$up"},
			"avgunknown": bson.M{"$avg": "$unknown"},
			"avgdown":    bson.M{"$avg": "$down"}}},
		{"$project": bson.M{
			"date":       "$_id.date",
			"name":       "$_id.name",
			"report":     "$_id.report",
			"supergroup": "$_id.supergroup",
			"unknown":    "$avgunknown",
			"up":         "$avgup",
			"down":       "$avgdown",
			"avgup":      1,
			"avgunknown": 1,
			"avgdown":    1,
			"availability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{1.00000001, "$avgunknown"}}}},
					100}},
			"reliability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgunknown"}}, "$avgdown"}}}},
					100}}}},
		{"$sort": bson.D{
			{"report", 1},
			{"supergroup", 1},
			{"name", 1},
			{"date", 1}}}}

	return query
}

// DailySuperGroup function to build the MongoDB aggregation query for daily calculations
func DailySuperGroup(filter bson.M) []bson.M {
	// Mongo aggregation pipeline
	// Select all the records that match q
	// Project the results to add 1 to every weights to avoid having 0 as a weights
	// Group them by the first 8 digits of datetime (YYYYMMDD) and each group find
	// availability = sum(availability*weights)
	// reliability = sum(reliability*weights)
	// weights = sum(weights)
	// Project to a better format and do these computations
	// availability = availability/weights
	// reliability = reliability/weights
	// Sort by report->supergroup->name->datetime
	query := []bson.M{
		{"$match": filter},
		{"$project": bson.M{
			"date":         1,
			"availability": 1,
			"reliability":  1,
			"report":       1,
			"supergroup":   1,
			"weight": bson.M{
				"$add": list{"$weight", 1}}},
		},
		{"$group": bson.M{
			"_id": bson.M{
				"date":       bson.D{{"$substr", list{"$date", 0, 8}}},
				"supergroup": "$supergroup",
				"report":     "$report"},
			"availability": bson.M{
				"$sum": bson.M{
					"$multiply": list{"$availability", "$weight"}}},
			"reliability": bson.M{
				"$sum": bson.M{
					"$multiply": list{"$reliability", "$weight"}}},
			"weight": bson.M{"$sum": "$weight"}},
		},
		{"$project": bson.M{
			"date":       "$_id.date",
			"supergroup": "$_id.supergroup",
			"report":     "$_id.report",
			"availability": bson.M{
				"$divide": list{"$availability", "$weight"}},
			"reliability": bson.M{
				"$divide": list{"$reliability", "$weight"}}},
		},
		{"$sort": bson.D{
			{"report", 1},
			{"supergroup", 1},
			{"date", 1}},
		}}

	return query
}

// MonthlySuperGroup function to build the MongoDB aggregation query for monthly calculations
func MonthlySuperGroup(filter bson.M) []bson.M {
	filter["availability"] = bson.M{"$gte": 0}
	filter["reliability"] = bson.M{"$gte": 0}
	// Mongo aggregation pipeline
	// Select all the records that match q
	// Project the results to add 1 to every weights to avoid having 0 as a weights
	// Group them by the first 8 digits of datetime (YYYYMMDD) and each group find
	// availability = sum(availability*weights)
	// reliability = sum(reliability*weights)
	// weights = sum(weights)
	// Project to a better format and do these computations
	// availability = availability/weights
	// reliability = reliability/weights
	// Group by the first 6 digits of the datetime (YYYYMM) and by ngi,site,profile and for each group find
	// availability = average(availability)
	// reliability = average(reliability)
	// Project the results to a better format
	// Sort by namespace->report->supergroup->datetime

	query := []bson.M{
		{"$match": filter},
		{"$project": bson.M{
			"date":         1,
			"availability": 1,
			"reliability":  1,
			"report":       1,
			"supergroup":   1,
			"weight": bson.M{
				"$add": list{"$weight", 1}}},
		},
		{"$group": bson.M{
			"_id": bson.M{
				"date":       bson.D{{"$substr", list{"$date", 0, 8}}},
				"supergroup": "$supergroup",
				"report":     "$report"},
			"availability": bson.M{"$sum": bson.M{"$multiply": list{"$availability", "$weight"}}},
			"reliability":  bson.M{"$sum": bson.M{"$multiply": list{"$reliability", "$weight"}}},
			"weight":       bson.M{"$sum": "$weight"}},
		},
		{"$match": bson.M{
			"weight": bson.M{"$gt": 0}},
		},
		{"$project": bson.M{
			"date":         "$_id.date",
			"supergroup":   "$_id.supergroup",
			"report":       "$_id.report",
			"availability": bson.M{"$divide": list{"$availability", "$weight"}},
			"reliability":  bson.M{"$divide": list{"$reliability", "$weight"}}},
		},
		{"$group": bson.M{
			"_id": bson.M{
				"date":       bson.D{{"$substr", list{"$date", 0, 6}}},
				"supergroup": "$supergroup", "report": "$report"},
			"availability": bson.M{"$avg": "$availability"},
			"reliability":  bson.M{"$avg": "$reliability"}},
		},
		{"$project": bson.M{
			"date":         "$_id.date",
			"supergroup":   "$_id.supergroup",
			"report":       "$_id.report",
			"availability": 1,
			"reliability":  1},
		},
		{"$sort": bson.D{
			{"report", 1},
			{"supergroup", 1},
			{"date", 1}},
		}}

	return query
}

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
