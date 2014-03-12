/*
 * Copyright (c) 2013 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
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

package vos

import (
	//"api/utils/caches"
	"api/utils/config"
	"api/utils/mongo"
	"net/http"
	"strings"
)

func VoAvailabilityInProfile(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte {

	// This is the input we will receive from the API
	urlValues := r.URL.Query()

	input := ApiVoAvailabilityInProfileInput{
		urlValues.Get("vo_name"),
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues["profile_name"],
		urlValues.Get("type"),
	}
	
	output := []byte("")

	// found, output := caches.HitCache("ngis", input, cfg)
// 	if found {
// 		return output
// 	}
	session := mongo.OpenSession(cfg)

	results := []ApiVoAvailabilityInProfileOutput{}

	err := error(nil)
	if len(input.availabilityperiod) == 0 || strings.ToLower(input.availabilityperiod) == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"

		query := Daily(input)
		
		
		err := mongo.Find(session, "AR", "voreports", query, "d", &results)
		if err != nil{
			panic(err)
		}

		//err = mongo.Pipe(session, "AR", "voreports", query, &results)

	} else if strings.ToLower(input.availabilityperiod) == "monthly" {
		
		return []byte("<root>WAIT!</root>")
		// customForm[0] = "200601"
// 		customForm[1] = "2006-01"
// 
// 		query := Monthly(input)
// 
// 		err = mongo.Pipe(session, "AR", "sites", query, &results)
	}

	if err != nil {
		return []byte("<root><error>" + err.Error() + "</error></root>")
	}
	
	output, err = CreateXMLResponse(results)
	
	// if len(results) > 0 {
// 		caches.WriteCache("ngis", input, output, cfg)
// 	}

	mongo.CloseSession(session)

	return output
}