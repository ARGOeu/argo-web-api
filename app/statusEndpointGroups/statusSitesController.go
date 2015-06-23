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

package statusEndpointGroup

import (
	//"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/utils/authentication"
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
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues.Get("vo"),
		urlValues.Get("profile"),
		urlValues.Get("group_type"),
		group,
	}

	fmt.Println(group)
	// Set default values
	if len(input.profile) == 0 {
		input.profile = "ch.cern.sam.ROC_CRITICAL"
	}

	if len(input.group_type) == 0 {
		input.group_type = "site"
	}

	if len(input.vo) == 0 {
		input.vo = "ops"
	}

	// Mongo Session
	results := []StatusEndpointGroupOutput{}

	session, err = mongo.OpenSession(tenantDbConfig)

	c := session.DB(tenantDbConfig.Db).C("status_sites")
	err = c.Find(prepQuery(input)).All(&results)

	mongo.CloseSession(session)

	output, err = createView(results, input) //Render the results into XML format
	//if strings.ToLower(input.format) == "json" {
	//	contentType = "application/json"
	//}
	//buffer.WriteString(strconv.Itoa(len(results)))
	//output = []byte(buffer.String())
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	return code, h, output, err
}

func prepQuery(input StatusEndpointGroupInput) bson.M {

	//Time Related
	const zuluForm = "2006-01-02T15:04:05Z"
	const ymdForm = "20060102"

	// timeStart, _ := time.Parse(zuluForm, input.start_time)
	// timeEnd, _ := time.Parse(zuluForm, input.end_time)
	// timeStartYMD, _ := strconv.Atoi(timeStart.Format(ymdForm))
	// //timeEndYMD, _ := strconv.Atoi(timeEnd.Format(ymdForm))
	// // parse time as integer
	// timeStart_int := (timeStart.Hour() * 10000) + (timeStart.Minute() * 100) + timeStart.Second()
	// timeEnd_int := (timeEnd.Hour() * 10000) + (timeEnd.Minute() * 100) + timeEnd.Second()
	//
	// fmt.Println(timeStart, timeEnd, timeStartYMD, timeStart_int, timeEnd_int)

	query := bson.M{}
	timeStartInt := 0
	timeEndInt := 0
	timeStartYMD := 0

	if input.start_time != "" {
		timeStart, _ := time.Parse(zuluForm, input.start_time)
		timeStartYMD, _ = strconv.Atoi(timeStart.Format(ymdForm))
		timeStartInt = (timeStart.Hour() * 10000) + (timeStart.Minute() * 100) + timeStart.Second()
		query["date_integer"] = timeStartYMD
	}

	if input.end_time != "" {
		timeEnd, _ := time.Parse(zuluForm, input.end_time)
		timeEndInt = (timeEnd.Hour() * 10000) + (timeEnd.Minute() * 100) + timeEnd.Second()

	}

	if input.end_time != "" && input.start_time != "" {
		query["time_integer"] = bson.M{"$gte": timeStartInt, "$lte": timeEndInt}

	}

	if input.endpointGroup != "" {
		query["endpoint_group"] = input.endpointGroup
		query["endpoint_group_type"] = input.endpointGroupType
	}

	if input.group != "" {
		query["group"] = input.Group
	}

	return query
	// if input.group_type == "site" {
	//
	// 	query := bson.M{
	// 		"date_integer":   timeStartYMD,
	// 		"endpoint_group": input.group,
	// 		"time_integer":   bson.M{"$gte": timeStart_int, "$lte": timeEnd_int},
	// 	}
	// 	return query
	//
	// } else if input.group_type == "ngi" {
	// 	query := bson.M{
	// 		"date_integer": timeStartYMD,
	// 		"group":        input.group,
	// 		"time_integer": bson.M{"$gte": timeStart_int, "$lte": timeEnd_int},
	// 	}
	// 	return query
	// }
	// return bson.M{"di": 0}

}
