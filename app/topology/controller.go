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

package topology

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// ListTopoStats list statistics about the topology used in the report
func ListTopoStats(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	groupType := context.Get(r, "group_type").(string)
	egroupType := context.Get(r, "endpoint_group_type").(string)

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	report_name := vars["report_name"]
	//Time Related
	dateStr := urlValues.Get("date")

	const zuluForm = "2006-01-02"
	const ymdForm = "20060102"

	ts := time.Now().UTC()
	if dateStr != "" {
		ts, _ = time.Parse(zuluForm, dateStr)
	}

	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// find the report id first
	reportID, err := mongo.GetReportID(session, tenantDbConfig.Db, report_name)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	var serviceResults []string
	var egroupResults []string
	var groupResults []string

	serviceCol := session.DB(tenantDbConfig.Db).C("service_ar")
	eGroupCol := session.DB(tenantDbConfig.Db).C("endpoint_group_ar")

	err = serviceCol.Find(bson.M{"report": reportID, "date": tsYMD}).Distinct("name", &serviceResults)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	err = eGroupCol.Find(bson.M{"report": reportID, "date": tsYMD}).Distinct("name", &egroupResults)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	err = eGroupCol.Find(bson.M{"report": reportID, "date": tsYMD}).Distinct("supergroup", &groupResults)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// fill the topology JSON struct
	result := Topology{}
	result.GroupCount = len(groupResults)
	result.GroupType = groupType
	result.GroupList = groupResults
	result.EndGroupCount = len(egroupResults)
	result.EndGroupType = egroupType
	result.EndGroupList = egroupResults
	result.ServiceCount = len(serviceResults)
	result.ServiceList = serviceResults

	output, err = createTopoView(result, contentType, http.StatusOK)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
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
