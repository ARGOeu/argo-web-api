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
	"fmt"
	"net/http"

	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
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

	parsedStart, _ := parseZuluDate(urlValues.Get("start_time"))

	parsedEnd, _ := parseZuluDate(urlValues.Get("end_time"))

	input := InputParams{
		parsedStart,
		parsedEnd,
		vars["report_name"],
		vars["group_type"],
		vars["group_name"],
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

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	// Query the detailed results

	// prepare the match filter
	report := reports.MongoInterface{}
	err = mongo.FindOne(session, tenantDbConfig.Db, "reports", bson.M{"info.name": vars["report_name"]}, &report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report_name"] + " does not exist"
		output, err := createErrorMessage(message, code, contentType) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	input.groupType = report.Topology.Group.Group.Type

	groupCol := session.DB(tenantDbConfig.Db).C(statusGroupColName)
	groupResults := []GroupData{}

	err = groupCol.Find(queryGroups(input, report.ID)).All(&groupResults)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	parsedPrev, _ := getPrevDay(urlValues.Get("start_time"))

	//if no status results yet show previous days results
	if len(groupResults) == 0 {
		// Zero query results
		input.startTime = parsedPrev
		err = groupCol.Find(queryGroups(input, report.ID)).All(&groupResults)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	endpointCol := session.DB(tenantDbConfig.Db).C(statusEndpointColName)
	endpointResults := []EndpointData{}

	err = endpointCol.Find(queryEndpoints(input, report.ID)).All(&endpointResults)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createCombinedView(groupResults, endpointResults, input, urlValues.Get("end_time"), details) //Render the results into JSON

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
	h.Set("Allow", fmt.Sprintf("GET, OPTIONS"))
	return code, h, output, err

}
