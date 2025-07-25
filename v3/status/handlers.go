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

package status

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// ListStatus lists group and endpoint status timelines
func ListStatus(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("List Metric Timelines")
	err := error(nil)
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	// This is going to be used to determine a detailed/latest view of the results
	view := urlValues.Get("view")
	latest := false
	details := false
	if view == "details" {
		details = true
	} else if view == "latest" {
		latest = true
	}

	var parsedStart, parsedEnd, parsedPrev int

	// check if user provided start and end time correctly
	urlStartTime := urlValues.Get("start_time")
	urlEndTime := urlValues.Get("end_time")

	endDate := ""

	if urlStartTime == "" && urlEndTime == "" {
		isoTimeNow := time.Now().UTC().Format(time.RFC3339)
		startDate := strings.Split(isoTimeNow, "T")[0]
		startTime := startDate + "T00:00:00Z"
		endDate = startDate
		parsedStart, _ = parseZuluDate(startTime)
		parsedEnd, _ = parseZuluDate(isoTimeNow)
		parsedPrev, _ = getPrevDay(startTime)
	} else {
		if parsedStart, err = parseZuluDate(urlStartTime); err != nil {
			code = http.StatusBadRequest
			message := fmt.Sprintf("Error parsing start_time=%s - please use zulu format like %s", urlStartTime, zuluForm)
			output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(message), contentType, "", " ")
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}

		if parsedEnd, err = parseZuluDate(urlEndTime); err != nil {
			code = http.StatusBadRequest
			message := fmt.Sprintf("Error parsing end_time=%s - please use zulu format like %s", urlEndTime, zuluForm)
			output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(message), contentType, "", " ")
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}

		endDate = strings.Split(urlEndTime, "T")[0]

		if latest {
			code = http.StatusBadRequest
			message := "Parameter view=latest should not be used when specifing start_time and end_time period"
			output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(message), contentType, "", " ")
			return code, h, output, err
		}

		parsedPrev, _ = getPrevDay(urlStartTime)

	}

	input := InputParams{
		parsedStart,
		parsedEnd,
		vars["report_name"],
		vars["group_type"],
		vars["group_name"],
		contentType,
		"",
	}

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	// Query the detailed results

	// prepare the match filter
	report := reports.MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	err = rCol.FindOne(context.TODO(), bson.M{"info.name": vars["report_name"]}).Decode(&report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report_name"] + " does not exist"
		output, err := createErrorMessage(message, code) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	input.groupType = report.Topology.Group.Group.Type

	groupCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(statusGroupColName)
	groupResults := []GroupData{}

	cursor, err := groupCol.Find(context.TODO(), queryGroups(input, report.ID))
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &groupResults)

	//if no status results yet show previous days results
	if len(groupResults) == 0 {
		// Zero query results
		input.startTime = parsedPrev
		cursor, err := groupCol.Find(context.TODO(), queryGroups(input, report.ID))
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &groupResults)
	}

	endpointCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(statusEndpointColName)
	endpointResults := []EndpointData{}

	cursor, err = endpointCol.Find(context.TODO(), queryEndpoints(input, report.ID))
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &endpointResults)
	output, err = createCombinedView(groupResults, endpointResults, input, endDate, details, latest) //Render the results into JSON

	return code, h, output, err
}

// ListStatusByID lists endpoint status timeline for a specific resource id
func ListStatusByID(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("List Metric Timelines")
	err := error(nil)
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	urlStartTime := urlValues.Get("start_time")
	urlEndTime := urlValues.Get("end_time")

	// This is going to be used to determine a detailed/latest view of the results
	view := urlValues.Get("view")
	latest := false
	if urlStartTime == "" && urlEndTime == "" {
		latest = true
	}
	details := false
	if view == "details" {
		details = true
		latest = false
	} else if view == "latest" {
		latest = true
	}

	var parsedStart, parsedEnd, parsedPrev int

	// check if user provided start and end time correctly

	endDate := ""

	if urlStartTime == "" && urlEndTime == "" {
		isoTimeNow := time.Now().UTC().Format(time.RFC3339)
		startDate := strings.Split(isoTimeNow, "T")[0]
		startTime := startDate + "T00:00:00Z"
		endDate = startDate
		parsedStart, _ = parseZuluDate(startTime)
		parsedEnd, _ = parseZuluDate(isoTimeNow)
		parsedPrev, _ = getPrevDay(startTime)
	} else {
		if parsedStart, err = parseZuluDate(urlStartTime); err != nil {
			code = http.StatusBadRequest
			message := fmt.Sprintf("Error parsing start_time=%s - please use zulu format like %s", urlStartTime, zuluForm)
			output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(message), contentType, "", " ")
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}

		if parsedEnd, err = parseZuluDate(urlEndTime); err != nil {
			code = http.StatusBadRequest
			message := fmt.Sprintf("Error parsing end_time=%s - please use zulu format like %s", urlEndTime, zuluForm)
			output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(message), contentType, "", " ")
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}

		endDate = strings.Split(urlEndTime, "T")[0]

		if latest {
			code = http.StatusBadRequest
			message := "Parameter view=latest should not be used when specifing start_time and end_time period"
			output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(message), contentType, "", " ")
			return code, h, output, err
		}

		parsedPrev, _ = getPrevDay(urlStartTime)

	}

	input := InputParams{
		parsedStart,
		parsedEnd,
		vars["report_name"],
		vars["group_type"],
		vars["group_name"],
		contentType,
		vars["id"],
	}

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	// Query the detailed results

	// prepare the match filter
	report := reports.MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	err = rCol.FindOne(context.TODO(), bson.M{"info.name": vars["report_name"]}).Decode(&report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report_name"] + " does not exist"
		output, err := createErrorMessage(message, code) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	input.groupType = report.Topology.Group.Group.Type

	groupCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(statusGroupColName)
	groupResults := []GroupData{}

	cursor, err := groupCol.Find(context.TODO(), queryGroups(input, report.ID))
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &groupResults)

	//if no status results yet show previous days results
	if len(groupResults) == 0 {
		// Zero query results
		input.startTime = parsedPrev
		cursor, err := groupCol.Find(context.TODO(), queryGroups(input, report.ID))
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &groupResults)
	}

	endpointCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(statusEndpointColName)
	endpointResults := []EndpointData{}

	cursor, err = endpointCol.Find(context.TODO(), queryEndpoints(input, report.ID))
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &endpointResults)

	if len(endpointResults) == 0 {
		code = http.StatusNotFound
		message := "No endpoints found with resource-id: " + vars["id"]
		output, err := createErrorMessage(message, code) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	output, err = createViewByID(endpointResults, input, endDate, details, latest) //Render the results into JSON

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
