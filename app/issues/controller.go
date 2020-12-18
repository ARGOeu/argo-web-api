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

package issues

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// FlatListEndpointTimelines returns a list of metric timelines
func FlatListEndpointTimelines(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Mongo Session
	results := []EndpointData{}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	metricCollection := session.DB(tenantDbConfig.Db).C("status_endpoints")

	// Query the detailed metric results
	reportID, err := mongo.GetReportID(session, tenantDbConfig.Db, vars["report_name"])

	dt, _, err := utils.ParseZuluDate(urlValues.Get("date"))

	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	err = metricCollection.Pipe(prepareIssueQuery(reportID, dt, urlValues.Get("filter"))).All(&results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createEndpointListView(results, "Success", 200)

	return code, h, output, err
}

func prepareIssueQuery(reportID string, dt int, filter string) []bson.M {
	// prepare the match filter
	match1 := bson.M{"report": reportID, "date_integer": dt}

	group := bson.M{
		"_id": bson.M{
			"host":           "$host",
			"service":        "$service",
			"endpoint_group": "$endpoint_group",
		},
		"status":         bson.M{"$last": "$status"},
		"timestamp":      bson.M{"$last": "$timestamp"},
		"host":           bson.M{"$last": "$host"},
		"service":        bson.M{"$last": "$service"},
		"endpoint_group": bson.M{"$last": "$endpoint_group"},
		"info":           bson.M{"$last": "$info"},
	}

	match2 := bson.M{"status": bson.M{"$ne": "OK"}}

	if strings.ToLower(filter) != "" {
		match2 = bson.M{"status": strings.ToUpper(filter)}
	}

	sorted := bson.M{"host": 1}

	agg := []bson.M{
		{"$match": match1},
		{"$group": group},
		{"$match": match2},
		{"$sort": sorted},
	}

	return agg

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
