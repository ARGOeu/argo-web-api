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

package statusServices

import (
	"fmt"
	"net/http"

	"github.com/ARGOeu/argo-web-api/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/ARGOeu/argo-web-api/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/authentication"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
)

// ListMetricTimelines returns a list of metric timelines
func ListServiceTimelines(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("List Metric Timelines")
	err := error(nil)
	contentType := "application/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	contentType, err = respond.ParseAcceptHeader(r)
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	if err != nil {
		code = http.StatusNotAcceptable
		output, _ = respond.MarshalContent(respond.NotAcceptableContentType, contentType, "", " ")
		return code, h, output, err
	}

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	parsedStart, parsedEnd, errs := respond.ValidateDateRange(urlValues.Get("start_time"), urlValues.Get("end_time"))
	if len(errs) > 0 {
		code = http.StatusBadRequest
		output = respond.CreateFailureResponseMessage("Bad Request", "400", errs).MarshalTo(contentType)
	}

	input := InputParams{
		parsedStart,
		parsedEnd,
		vars["report_name"],
		vars["group_type"],
		vars["group_name"],
		vars["service_name"],
		contentType,
	}

	// Call authenticateTenant to check the api key and retrieve
	// the correct tenant db conf
	tenantDbConfig, err := authentication.AuthenticateTenant(r.Header, cfg)

	if err != nil {
		output = []byte(http.StatusText(http.StatusUnauthorized))
		code = http.StatusUnauthorized //If wrong api key is passed we return UNAUTHORIZED http status
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	// Mongo Session
	results := []DataOutput{}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	metricCollection := session.DB(tenantDbConfig.Db).C("status_services")

	// Query the detailed metric results
	reportID, err := mongo.GetReportID(session, tenantDbConfig.Db, input.report)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	err = metricCollection.Find(prepareQuery(input, reportID)).All(&results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createView(results, input) //Render the results into XML format

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

func prepareQuery(input InputParams, reportID string) bson.M {

	// prepare the match filter
	filter := bson.M{
		"date_integer":   bson.M{"$gte": input.startTime, "$lte": input.endTime},
		"report":         reportID,
		"endpoint_group": input.group,
	}

	if len(input.service) > 0 {
		filter["service"] = input.service
	}

	return filter
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
