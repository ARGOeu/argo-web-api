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

// List all the current recomputation request from a specific tenantdb
func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	//TODO: change this to the actual tenantdb using a call
	// tenantdb := get_tenant_db(r , cfg) where r is the http request
	// that has the header x-api-key which needs to be checked
	tenantDbConfig, err := authentication.AuthenticateTenant(r.Header, cfg)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	session, err := mongo.OpenSession(tenantDbConfig)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	results := []InputOutput{}
	err = mongo.Find(session, tenantDbConfig.Db, "recomputations", nil, "timestamp", &results)

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

//Create handles new recomputation requests
func Create(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	tenantDbConfig, err := authentication.AuthenticateTenant(r.Header, cfg)

	//only authenticated requests triger the handling code
	if err != nil {
		output = []byte(http.StatusText(http.StatusUnauthorized))
		code = http.StatusUnauthorized //If wrong api key is passed we return UNAUTHORIZED http status
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	session, err := mongo.OpenSession(tenantDbConfig)
	err = r.ParseForm()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	urlValues := r.Form
	fmt.Println(urlValues)
	now := time.Now()
	input := InputOutput{
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues.Get("reason"),
		urlValues.Get("group"),
		urlValues["exclude_site"],
		"pending",
		now.Format("2006-01-02 15:04:05"),
	}

	query := insertQuery(input)
	err = mongo.Insert(session, tenantDbConfig.Db, "recomputations", query)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	message := "A recalculation request has been filed"
	output, err = messageXML(message) //Render the response into XML

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err

}
