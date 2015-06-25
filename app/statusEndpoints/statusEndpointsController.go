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

package statusEndpoints

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	//STANDARD DECLARATIONS END

	tenantDbConfig, err := authentication.AuthenticateTenant(r.Header, cfg)

	if err != nil {
		if err.Error() == "Unauthorized" {
			code = http.StatusUnauthorized
			return code, h, output, err
		}
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

	input := StatusEndpointsInput{
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues.Get("job"),
		urlValues.Get("group_type"),
		group,
	}

	if len(input.group_type) == 0 {
		input.group_type = "site"
	}

	// Mongo Session
	results := []StatusEndpointsOutput{}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	c := session.DB(tenantDbConfig.Db).C("status_endpoints")
	err = c.Find(prepQuery(input)).All(&results)

	output, err = createView(results, input) //Render the results into XML format
	//if strings.ToLower(input.format) == "json" {
	//	contentType = "application/json"
	//}
	//buffer.WriteString(strconv.Itoa(len(results)))
	//output = []byte(buffer.String())
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	return code, h, output, err
}

func prepQuery(input StatusEndpointsInput) bson.M {

	//Time Related
	const zuluForm = "2006-01-02T15:04:05Z"
	const ymdForm = "20060102"

	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))

	tsInt := (ts.Hour() * 10000) + (ts.Minute() * 100) + ts.Second()
	teInt := (te.Hour() * 10000) + (te.Minute() * 100) + te.Second()

	if input.group_type == "endpoint" {
		query := bson.M{
			"job":            input.job,
			"date_int":       tsYMD,
			"endpoint_group": input.group,
			"time_int":       bson.M{"$gte": tsInt, "$lte": teInt},
		}

		return query

	} else if input.group_type == "group" {
		query := bson.M{
			"job":            input.job,
			"date_int":       tsYMD,
			"supergroup":     input.group,
			"time_int":       bson.M{"$gte": tsInt, "$lte": teInt},
		}

		return query

	} else if input.group_type == "host" {
		query := bson.M{
			"job":            input.job,
			"date_int":       tsYMD,
			"host":           input.group,
			"time_int":       bson.M{"$gte": tsInt, "$lte": teInt},
		}

		return query

	}

	return bson.M{"date_int": 0}

}
