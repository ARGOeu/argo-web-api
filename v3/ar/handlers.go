/*
 * Copyright (c) 2022 GRNET S.A.
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

package ar

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ListGroupAR lists supergroup and group availabilities according to the http request
func ListGroupAR(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	report := reports.MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	err = rCol.FindOne(context.TODO(), bson.M{"info.name": vars["report_name"]}).Decode(&report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report_name"] + " does not exist"
		output, err := createErrorMessage(message, code, contentType) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	input := GroupResultQuery{
		basicQuery{
			Name:        "",
			Granularity: urlValues.Get("granularity"),
			Format:      contentType,
			StartTime:   urlValues.Get("start_time"),
			EndTime:     urlValues.Get("end_time"),
			Report:      report,
			Vars:        vars,
		}, "",
	}

	errs := input.Validate()
	if len(errs) > 0 {
		out := respond.BadRequestSimple
		out.Errors = errs
		output = out.MarshalTo(contentType)
		code = 400
		return code, h, output, err
	}

	resultsGroups := []GroupInterface{}
	resultsSuperGroups := []SuperGroupInterface{}

	// Construct the query to mongodb based on the input
	filter := bson.M{
		"date":   bson.M{"$gte": input.StartTimeInt, "$lte": input.EndTimeInt},
		"report": report.ID,
	}

	// Prepare the supergroup results
	// Select the granularity of the search daily/monthly
	arCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(groupColName)
	var query []primitive.M
	custom := false
	if input.Granularity == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		query = DailySuperGroup(filter)

	} else if input.Granularity == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = MonthlySuperGroup(filter)

	} else if input.Granularity == "custom" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = CustomSuperGroup(filter)
		custom = true
	}

	cursor, err := arCol.Aggregate(context.TODO(), query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &resultsSuperGroups)

	// Prepare the group results
	// Select the granularity of the search daily/monthly
	if input.Granularity == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		query = DailyGroup(filter)

	} else if input.Granularity == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = MonthlyGroup(filter)

	} else if input.Granularity == "custom" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = CustomGroup(filter)
		custom = true
	}

	cursor, err = arCol.Aggregate(context.TODO(), query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &resultsGroups)

	output, err = createResultView(resultsSuperGroups, resultsGroups, report, custom)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

// ListEndpointARByID lists endpoints a/r results for a specific resource id
func ListEndpointARByID(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	report := reports.MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	err = rCol.FindOne(context.TODO(), bson.M{"info.name": vars["report_name"]}).Decode(&report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report_name"] + " does not exist"
		output, err := createErrorMessage(message, code, contentType) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	input := basicQuery{
		Name:        "",
		Granularity: urlValues.Get("granularity"),
		Format:      contentType,
		StartTime:   urlValues.Get("start_time"),
		EndTime:     urlValues.Get("end_time"),
		Report:      report,
		Vars:        vars,
	}

	// if user has not defined a start/end period construct by default the a/r period including the days of this month
	if input.StartTime == "" && input.EndTime == "" {
		input.StartTime = time.Now().UTC().Format("2006-01") + "-01T00:00:00Z"
		input.EndTime = time.Now().UTC().Format("2006-01-02") + "T23:59:59Z"
	}

	errs := input.Validate()
	if len(errs) > 0 {
		out := respond.BadRequestSimple
		out.Errors = errs
		output = out.MarshalTo(contentType)
		code = 400
		return code, h, output, err
	}

	results := []EndpointInterface{}

	// Construct the query to mongodb based on the input
	filter := bson.M{
		"date":    bson.M{"$gte": input.StartTimeInt, "$lte": input.EndTimeInt},
		"report":  report.ID,
		"info.ID": input.Vars["id"],
	}

	// Prepare the group results
	// Select the granularity of the search daily/monthly
	arCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(endpointColName)
	var query []primitive.M
	custom := false
	if input.Granularity == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		query = DailyEndpoint(filter)

	} else if input.Granularity == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = MonthlyEndpoint(filter)

	} else if input.Granularity == "custom" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = CustomEndpoint(filter)
		custom = true
	}

	cursor, err := arCol.Aggregate(context.TODO(), query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	// if number of returned results is 0 respond with not found
	if len(results) == 0 {
		code = http.StatusNotFound
		message := "No endpoints found with resource-id: " + vars["id"]
		output, err := createErrorMessage(message, code, contentType) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	output, err = createEndpointResult(results, input.Vars["id"], custom)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
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
	h.Set("Allow", "GET, OPTIONS")
	return code, h, output, err

}
