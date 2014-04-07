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

package availabilityProfiles

import (
	"api/utils/authentication"
	"api/utils/config"
	"api/utils/mongo"
	"net/http"
	"strings"
)

func DeleteProfiles(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte {
	
	answer:=""
	
	if authentication.Authenticate(r.Header, cfg) {
		
		urlValues := r.URL.Path
		
		id := strings.Split(urlValues,"/")[4]
		
		session := mongo.OpenSession(cfg)
	
		err := mongo.IdRemove(session,"AR","aps",id)
			
		if err!=nil{
			answer="No profile matching the requested id"
		} else {
			answer = "Delete successful"
		}
	} else {
		answer = http.StatusText(403)
	}
	output, err := messageXML(answer)
	if err != nil {
		panic(err)
	}
 	return output
	
}
