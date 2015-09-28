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

package recomputations2

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ARGOeu/argo-web-api/utils/authentication"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
)

var recomputationsColl = "recomputations"

func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	// contentType := "application/json"
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	urlValues := r.URL.Query()
	contentType := r.Header.Get("Accept")

	filter := IncomingRecomputation{
		StartTime: urlValues.Get("start_time"),
		EndTime:   urlValues.Get("end_time"),
		Reason:    urlValues.Get("reason"),
		Report:    urlValues.Get("report"),
	}

	tenantDbConfig, err := authentication.AuthenticateTenant(r.Header, cfg)

	if err != nil {
		output = []byte(http.StatusText(http.StatusUnauthorized))
		code = http.StatusUnauthorized //If wrong api key is passed we return UNAUTHORIZED http status
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	session, err := mongo.OpenSession(tenantDbConfig)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	results := []MongoInterface{}
	err = mongo.Find(session, tenantDbConfig.Db, recomputationsColl, filter, "timestamp", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createListView(results, contentType)

	return code, h, output, err

}

func SubmitRecomputation(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	contentType := r.Header.Get("Accept")

	tenantDbConfig, err := authentication.AuthenticateTenant(r.Header, cfg)

	if err != nil {
		output = []byte(http.StatusText(http.StatusUnauthorized))
		code = http.StatusUnauthorized //If wrong api key is passed we return UNAUTHORIZED http status
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	session, err := mongo.OpenSession(tenantDbConfig)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	var recompSubmission IncomingRequest
	// urlValues := r.URL.Query()

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &recompSubmission); err != nil {
		code = 422 // unprocessable entity
		output = []byte("Unprocessable JSON")
		return code, h, output, err
	}
	now := time.Now()
	recomputations := []MongoInterface{}
	for _, r := range recompSubmission.Data {
		recomputations = append(recomputations,
			MongoInterface{
				RequesterName:  tenantDbConfig.User,
				RequesterEmail: tenantDbConfig.Email,
				StartTime:      r.StartTime,
				EndTime:        r.EndTime,
				Reason:         r.Reason,
				Report:         r.Report,
				Exclude:        r.Exclude,
				Timestamp:      now.Format("2006-01-02 15:04:05"),
				Status:         "pending",
			})
	}

	// err = session.DB(cfg.MongoDB.Db).C(recomputationsColl).Insert(recomputations)

	err = mongo.MultiInsert(session, tenantDbConfig.Db, recomputationsColl, recomputations)

	if err != nil {
		panic(err)
	}

	output, err = createSubmitView(recomputations, contentType)

	return code, h, output, err
}
