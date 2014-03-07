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

package services

import (
	"api/utils/config"
	"api/utils/mongo"
	"api/utils/caches"
	"net/http"
)

//Reply to requests about service_availability_in_profile
func ServiceAvailabilityInProfile(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte {

	// Parse the request into the input
	urlValues := r.URL.Query()

	input := ApiServiceAvailabilityInProfileInput{
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues["vo_name"],
		urlValues["profile_name"],
		urlValues["group_type"],
		urlValues.Get("type"),
		urlValues.Get("output"),
		urlValues["namespace"],
		urlValues["group_name"],
		urlValues["service_flavour"],
		urlValues["service_hostname"],
	}

	found, output := caches.HitCache("service_endpoint ", input, cfg)
	if found {
		return output
	}
	
	err := error(nil)
	// Create a mongodb session
	session := mongo.OpenSession(cfg)

	results := []ApiServiceAvailabilityInProfileOutput{}

	query := Timeline(input)

	err = mongo.Pipe(session, "AR", "sites", query, &results)

	//err = c.Find(q).Sort("p", "h", "sf").All(&results)
	if err != nil {
		return []byte("<root><error>" + err.Error() + "</error></root>")
	}

	//rootfmt.Println(results)
	output, err = CreateXMLResponse(results)

	if len(results) > 0 {
		caches.WriteCache("service_endpoint ", input, output, cfg)
	}

	mongo.CloseSession(session)

	return output
}
