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

package siteAvailability

import (
	"fmt"
	"github.com/argoeu/ar-web-api/utils/caches"
	"github.com/argoeu/ar-web-api/utils/config"
	"github.com/argoeu/ar-web-api/utils/mongo"
	"net/http"
	"strings"
)

// func CreateResponse(w http.ResponseWriter, r *http.Request, cfg config.Config) ([]byte, error,int){
//
// }

func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK

	h := http.Header{}

	output := []byte("")

	err := error(nil)

	contentType := "text/xml"

	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Parse the request into the input
	urlValues := r.URL.Query()

	input := SiteAvailabilityInput{
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues.Get("availability_profile"),
		urlValues.Get("granularity"),
		urlValues.Get("infrastructure"),
		urlValues.Get("production"),
		urlValues.Get("monitored"),
		urlValues.Get("certification"),
		urlValues.Get("format"),
		urlValues["group_name"],
	}

	if len(input.infrastructure) == 0 {
		input.infrastructure = "Production"
	}

	if len(input.production) == 0 || input.production == "true" {
		input.production = "Y"
	} else {
		input.production = "N"
	}

	if len(input.monitored) == 0 || input.monitored == "true" {
		input.monitored = "Y"
	} else {
		input.monitored = "N"
	}

	if len(input.certification) == 0 {
		input.certification = "Certified"
	}

	if strings.ToLower(input.format) == "json" {

		contentType = "application/json"

	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	found, output := caches.HitCache("sites", input, cfg)

	if found {

		return code, h, output, err

	}

	session, err := mongo.OpenSession(cfg)

	if err != nil {

		code = http.StatusInternalServerError

		return code, h, output, err
	}

	results := []SiteAvailabilityOutput{}

	// Select the granularity of the search daily/monthly
	if len(input.granularity) == 0 || strings.ToLower(input.granularity) == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"

		query := Daily(input)

		err = mongo.Pipe(session, "AR", "sites", query, &results)

	} else if strings.ToLower(input.granularity) == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"

		query := Monthly(input)

		err = mongo.Pipe(session, "AR", "sites", query, &results)
	}

	if err != nil {

		code = http.StatusInternalServerError

		return code, h, output, err
	}

	output, err = createView(results, input.format)

	if err != nil {

		code = http.StatusInternalServerError

		return code, h, output, err
	}

	if len(results) > 0 {
		caches.WriteCache("sites", input, output, cfg)
	}

	mongo.CloseSession(session)

	return code, h, output, err
}
