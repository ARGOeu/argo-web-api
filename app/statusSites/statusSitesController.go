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

package statusSites

import (
	//"bytes"
	"fmt"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	// URL PATH_VALUES
	urlPath := r.URL.Path
	group := strings.Split(urlPath, "/")[6]

	urlValues := r.URL.Query()

	input := StatusSitesInput{
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues.Get("vo"),
		urlValues.Get("profile"),
		urlValues.Get("group_type"),
		group,
	}

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
	results := []StatusSitesOutput{}

	session, err := mongo.OpenSession(cfg)

	c := session.DB("AR").C("status_sites")
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

func prepQuery(input StatusSitesInput) bson.M {

	//Time Related
	const zuluForm = "2006-01-02T15:04:05Z"
	const ymdForm = "20060102"

	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	//teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	// parse time as integer
	ts_int := (ts.Hour() * 10000) + (ts.Minute() * 100) + ts.Second()
	te_int := (te.Hour() * 10000) + (te.Minute() * 100) + te.Second()

	if input.group_type == "site" {

		query := bson.M{
			"di":   tsYMD,
			"site": input.group,
			"ti":   bson.M{"$gte": ts_int, "$lte": te_int},
		}

		return query

	} else if input.group_type == "ngi" {
		query := bson.M{
			"di":  tsYMD,
			"roc": input.group,
			"ti":  bson.M{"$gte": ts_int, "$lte": te_int},
		}

		return query

	}

	return bson.M{"di": 0}

}
