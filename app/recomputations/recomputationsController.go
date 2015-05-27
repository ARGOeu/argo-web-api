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

package recomputations

import (
	"fmt"
	"net/http"
	"time"

	"github.com/argoeu/argo-web-api/utils/authentication"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"
)

func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//TODO: change this to the actual tenantdb
	tenantdb := "AR_test"
	//STANDARD DECLARATIONS END

	session, err := mongo.OpenSession(cfg.MongoDB)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	results := []RecomputationsInputOutput{}
	err = mongo.Find(session, tenantdb, "recalculations", nil, "t", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createView(results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	mongo.CloseSession(session)
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

func Create(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	message := ""
	//only authenticated requests triger the handling code
	if authentication.Authenticate(r.Header, cfg) {

		session, err := mongo.OpenSession(cfg.MongoDB)

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		err = r.ParseForm()

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		urlValues := r.Form
		now := time.Now()

		input := RecomputationsInputOutput{
			urlValues.Get("start_time"),
			urlValues.Get("end_time"),
			urlValues.Get("reason"),
			urlValues.Get("ngi_name"),
			urlValues["exclude_site"],
			"pending",
			now.Format("2006-01-02 15:04:05"),
			//urlValues["exclude_sf"],
			//urlValues["exclude_end_point"],
		}

		query := insertQuery(input)
		err = mongo.Insert(session, "AR", "recalculations", query)

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		mongo.CloseSession(session)
		message = "A recalculation request has been filed"
		output, err := messageXML(message) //Render the response into XML

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err

	} else {
		output = []byte(http.StatusText(http.StatusUnauthorized))
		code = http.StatusUnauthorized //If wrong api key is passed we return UNAUTHORIZED http status
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}
}
