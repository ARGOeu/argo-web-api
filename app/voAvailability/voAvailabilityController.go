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

package voAvailability

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

	input := ApiVoAvailabilityInProfileInput{
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues.Get("availability_profile"),
		urlValues.Get("granularity"),
		//urlValues.Get("format"),
		urlValues["group_name"],
	}

	output := []byte("")

	found, output := caches.HitCache("vos", input, cfg)
	if found {
		return output
	}
	session := mongo.OpenSession(cfg)

	results := []ApiVoAvailabilityInProfileOutput{}

	err := error(nil)
	if len(input.granularity) == 0 || strings.ToLower(input.granularity) == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"

		query := Daily(input)

		err = mongo.Pipe(session, "AR", "voreports", query, &results)
		if err != nil {
			panic(err)
		}

		//err = mongo.Pipe(session, "AR", "voreports", query, &results)

	} else if strings.ToLower(input.granularity) == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"

		query := Monthly(input)

		err = mongo.Pipe(session, "AR", "voreports", query, &results)

		if err != nil {
			panic(err)
		}

	}
	output, err = createResponse(results)

	if len(results) > 0 {
		caches.WriteCache("vos", input, output, cfg)
	}

	mongo.CloseSession(session)

	return output
}

func createResponse(results []ApiVoAvailabilityInProfileOutput) ([]byte, error) {

	docRoot := &Root{}

	prevProfile := ""
	prevVo := ""
	vo := &Vo{}
	profile := &Profile{}
	// we iterate through the results struct array
	// keeping only the value of each row
	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], row.Date)
		//if new profile value does not match the previous profile value
		//we create a new profile in the xml
		if prevProfile != row.Profile {
			prevProfile = row.Profile
			profile = &Profile{
				Name: row.Profile,
			}
			docRoot.Profile = append(docRoot.Profile, profile)
			prevVo = ""
		}
		//if new ngi does not match the previous ngi value
		//we create a new ngi entry in the xml
		if prevVo != row.Vo {
			prevVo = row.Vo
			vo = &Vo{
				Vo: row.Vo,
			}
			profile.Vo = append(profile.Vo, vo)
		}
		//we append the new availability values
		vo.Availability = append(vo.Availability,
			&Availability{
				Timestamp:    timestamp.Format(customForm[1]),
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability)})
	}
	//we create the xml response and record the output and any possible errors
	//in the appropriate variables
	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	//we return the output
	return output, err
}
