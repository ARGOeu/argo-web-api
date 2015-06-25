/*
 * Copyright (c) 2014 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
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
 * or implied, of either GRNET S.A., SRCE or IN2P3 CNRS Computing
 * Centre
 *
 * The work represented by this source file is partially funded by
 * the EGI-InSPIRE project through the European Commission's 7th
 * Framework Programme (contract # INFSO-RI-261323)
 */

package statusDetail

import (
	//"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/argoeu/argo-web-api/app/jobs"
	"github.com/argoeu/argo-web-api/app/metricProfiles"
	"github.com/argoeu/argo-web-api/utils/authentication"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"

	"labix.org/v2/mgo/bson"
)

// List the status metric details
func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("Hello there!")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//var buffer bytes.Buffer

	//STANDARD DECLARATIONS END

	// URL PATH_VALUES
	urlPath := r.URL.Path
	group := strings.Split(urlPath, "/")[6]

	urlValues := r.URL.Query()

	input := InputParams{
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues.Get("job"),
		urlValues.Get("group_type"),
		group,
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

	// Structure to hold job information
	jobResult := jobs.Job{}
	metricProfileResult := metricProfiles.MongoInterface{}

	// Mongo Session
	results := []DataOutput{}

	selectedGroupType := ""
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	metricCol := session.DB(tenantDbConfig.Db).C("status_metric")
	jobCol := session.DB(tenantDbConfig.Db).C("jobs")
	profileCol := session.DB(tenantDbConfig.Db).C("metric_profiles")

	// Get Job details
	err = jobCol.Find(bson.M{"name": input.job}).One(&jobResult)
	if err != nil {
		output = []byte("Error on retrieving job information")
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Search job for used metric profile
	metricProfileName := ""
	metricProfileName, err = jobs.GetMetricProfile(jobResult)

	// Query details for the metric profile used
	err = profileCol.Find(bson.M{"name": metricProfileName}).One(&metricProfileResult)
	if err != nil {

		output = []byte("Error on retrieving metric profile information")
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	// Find if the selected group type is endpoint group
	// or group of groups. If is not found in the job
	if jobResult.GroupOfGroups == input.groupType {
		selectedGroupType = "group"
	} else if jobResult.EndpointGroup == input.groupType {
		selectedGroupType = "endpoint"
	} else {
		message := "the specific group type is not supported in this job"
		output, err := messageXML(message) //Render the response into XML
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err

	}

	// Query the detailed metric results
	err = metricCol.Find(prepQuery(input, selectedGroupType)).All(&results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create the view
	output, err = createView(results, input, metricProfileResult) //Render the results into XML format
	//if strings.ToLower(input.format) == "json" {
	//	contentType = "application/json"
	//}
	//buffer.WriteString(strconv.Itoa(len(results)))
	//output = []byte(buffer.String())
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	return code, h, output, err
}

func prepQuery(input InputParams, selectedGroupType string) bson.M {

	//Time Related
	const zuluForm = "2006-01-02T15:04:05Z"
	const ymdForm = "20060102"

	ts, _ := time.Parse(zuluForm, input.startTime)
	te, _ := time.Parse(zuluForm, input.endTime)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	//teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	// parse time as integer
	tsInt := (ts.Hour() * 10000) + (ts.Minute() * 100) + ts.Second()
	teInt := (te.Hour() * 10000) + (te.Minute() * 100) + te.Second()

	if selectedGroupType == "endpoint" {

		query := bson.M{
			"job":            input.job,
			"date_int":       tsYMD,
			"endpoint_group": input.group,
			"time_int":       bson.M{"$gte": tsInt, "$lte": teInt},
		}

		return query

	} else if selectedGroupType == "group" {

		query := bson.M{
			"job":        input.job,
			"date_int":   tsYMD,
			"supergroup": input.group,
			"time_int":   bson.M{"$gte": tsInt, "$lte": teInt},
		}

		return query

	} else if input.groupType == "host" {
		query := bson.M{
			"job":      input.job,
			"date_int": tsYMD,
			"host":     input.group,
			"time_int": bson.M{"$gte": tsInt, "$lte": teInt},
		}

		return query
	}

	return bson.M{"date_int": 0}

}
