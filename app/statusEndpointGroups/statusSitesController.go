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

package statusEndpointGroups

import (
	//"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/argoeu/argo-web-api/utils/authentication"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"
	"labix.org/v2/mgo/bson"
)

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

	tenantDbConfig, err := authentication.AuthenticateTenant(r.Header, cfg)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// URL PATH_VALUES
	urlPath := r.URL.Path
	group := strings.Split(urlPath, "/")[6]

	urlValues := r.URL.Query()

	input := StatusEndpointGroupInput{
		Start:      urlValues.Get("start_time"),
		End:        urlValues.Get("end_time"),
		Job:        urlValues.Get("job"),
		SuperGroup: urlValues.Get("supergroup_name"),
		Name:       group,
	}

	// Mongo Session
	results := []Job{}

	session, err = mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	c := session.DB(tenantDbConfig.Db).C("status_endpointgroups")
	// err = c.Find(prepQuery(input)).All(&results)

	query := aggregateQuery(input)
	err = c.Pipe(query).All(&results)

	output, err = createView(results, input) //Render the results into XML format
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	return code, h, output, err
}

func aggregateQuery(input StatusEndpointGroupInput) []bson.M {

	filter := prepQuery(input)

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"job":        "$job",
				"supergroup": "$supergroup",
				"name":       "$name"},
			"name": bson.M{
				"$first": "$name",
			},
			"statuses": bson.M{
				"$push": "$$ROOT",
			}}},
		{"$group": bson.M{
			"_id": bson.M{
				"job": "$_id.job",
			},
			"endpointgroup": bson.M{
				"$push": "$$ROOT",
			}}},
		{"$project": bson.M{
			"job":           "$_id.job",
			"endpointgroup": "$endpointgroup",
			// "name":            "$_id.name",
			// "supergroup":      "$statuses.supergroup",
			// "status":          "$statuses.status",
			// "previous_status": "$statuses.previous_status",
			// "timestamp":       "$statuses.timestamp",
		}},
		// {"$sort": bson.D{
		// 	{"job", 1},
		// 	{"supergroup", 1},
		// 	{"name", 1},
		// 	{"timestamp", 1},
		// 	// {"date_integer", 1},
		// 	// {"time_integer", 1},
		// }},
	}

	return query

}

func prepQuery(input StatusEndpointGroupInput) bson.M {

	//Time Related
	const zuluForm = "2006-01-02T15:04:05Z"
	const ymdForm = "20060102"

	// timeStart, _ := time.Parse(zuluForm, input.Start)
	// timeEnd, _ := time.Parse(zuluForm, input.End)
	// timeStartYMD, _ := strconv.Atoi(timeStart.Format(ymdForm))
	// //timeEndYMD, _ := strconv.Atoi(timeEnd.Format(ymdForm))
	// // parse time as integer
	// timeStart_int := (timeStart.Hour() * 10000) + (timeStart.Minute() * 100) + timeStart.Second()
	// timeEnd_int := (timeEnd.Hour() * 10000) + (timeEnd.Minute() * 100) + timeEnd.Second()

	query := bson.M{}
	timeStartInt := 0
	timeEndInt := 0
	timeStartYMD := 0

	if input.Start != "" {
		timeStart, _ := time.Parse(zuluForm, input.Start)
		timeStartYMD, _ = strconv.Atoi(timeStart.Format(ymdForm))
		timeStartInt = (timeStart.Hour() * 10000) + (timeStart.Minute() * 100) + timeStart.Second()
		query["date_integer"] = timeStartYMD
	}

	if input.End != "" {
		timeEnd, _ := time.Parse(zuluForm, input.End)
		timeEndInt = (timeEnd.Hour() * 10000) + (timeEnd.Minute() * 100) + timeEnd.Second()
	}
	if input.End != "" && input.Start != "" {
		query["time_integer"] = bson.M{"$gte": timeStartInt, "$lte": timeEndInt}

	}

	if input.Name != "" {
		query["name"] = input.Name
	}

	if input.SuperGroup != "" {
		query["supergroup"] = input.SuperGroup
	}

	return query

}
