/*
 * Copyright (c) 2014 GRNET S.A.
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

package ngiAvailability

import (
	"encoding/xml"
	"fmt"
	"github.com/argoeu/ar-web-api/utils/caches"
	"github.com/argoeu/ar-web-api/utils/config"
	"github.com/argoeu/ar-web-api/utils/mongo"
	"net/http"
	"strings"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte {

	// This is the input we will receive from the API
	urlValues := r.URL.Query()

	input := ApiNgiAvailabilityInProfileInput{
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues.Get("availability_profile"),
		urlValues.Get("granularity"),
		urlValues.Get("infrastructure"),
		urlValues.Get("production"),
		urlValues.Get("monitored"),
		urlValues.Get("certification"),
		//urlValues.Get("format"),
		urlValues["group_name"],
	}

	if len(input.Infrastructure) == 0 {
		input.Infrastructure = "Production"
	}

	if len(input.Production) == 0 || input.Production == "true" {
		input.Production = "Y"
	} else {
		input.Production = "N"
	}

	if len(input.Monitored) == 0 || input.Monitored == "true" {
		input.Monitored = "Y"
	} else {
		input.Monitored = "N"
	}

	if len(input.Certification) == 0 {
		input.Certification = "Certified"
	}

	found, output := caches.HitCache("ngis", input, cfg)
	if found {
		return output
	}
	session := mongo.OpenSession(cfg)

	results := []ApiNgiAvailabilityInProfileOutput{}

	err := error(nil)
	if len(input.Granularity) == 0 || strings.ToLower(input.Granularity) == "daily" {
		CustomForm[0] = "20060102"
		CustomForm[1] = "2006-01-02"

		query := Daily(input)

		err = mongo.Pipe(session, "AR", "sites", query, &results)

	} else if strings.ToLower(input.Granularity) == "monthly" {
		CustomForm[0] = "200601"
		CustomForm[1] = "2006-01"

		query := Monthly(input)

		err = mongo.Pipe(session, "AR", "sites", query, &results)
	}

	if err != nil {
		return []byte("<root><error>" + err.Error() + "</error></root>")
	}

	output, err = createResponse(results)
	if len(results) > 0 {
		caches.WriteCache("ngis", input, output, cfg)
	}

	mongo.CloseSession(session)

	return output
}

func createResponse(results []ApiNgiAvailabilityInProfileOutput) ([]byte, error) {

	docRoot := &Root{}

	prevProfile := ""
	prevNgi := ""
	ngi := &Ngi{}
	profile := &Profile{}
	// we iterate through the results struct array
	// keeping only the value of each row
	for _, row := range results {
		timestamp, _ := time.Parse(CustomForm[0], row.Date)
		//if new profile value does not match the previous profile value
		//we create a new profile in the xml
		if prevProfile != row.Profile {
			prevProfile = row.Profile
			profile = &Profile{
				Name:      row.Profile,
				Namespace: row.Namespace}
			docRoot.Profile = append(docRoot.Profile, profile)
			prevNgi = ""
		}
		//if new ngi does not match the previous ngi value
		//we create a new ngi entry in the xml
		if prevNgi != row.Ngi {
			prevNgi = row.Ngi
			ngi = &Ngi{
				Ngi: row.Ngi,
			}
			profile.Ngi = append(profile.Ngi, ngi)
		}
		//we append the new availability values
		ngi.Availability = append(ngi.Availability,
			&Availability{
				Timestamp:    timestamp.Format(CustomForm[1]),
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability)})
	}
	//we create the xml response and record the output and any possible errors
	//in the appropriate variables
	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	//we return the output
	return output, err
}
