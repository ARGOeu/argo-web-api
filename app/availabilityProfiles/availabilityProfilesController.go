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

package availabilityProfiles

import (
	"encoding/json"
	"github.com/argoeu/ar-web-api/utils/authentication"
	"github.com/argoeu/ar-web-api/utils/config"
	"github.com/argoeu/ar-web-api/utils/mongo"
	"io/ioutil"
	"net/http"
	"strings"
	"fmt"
)

func List(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte {

	err := error(nil)

	recordId := []string{}

	//Read the search values
	urlValues := r.URL.Query()

	//Searchig is based on name and namespace
	input := AvailabilityProfileSearch{
		urlValues["name"],
		urlValues["namespace"],
	}

	results := []AvailabilityProfileOutput{}

	session := mongo.OpenSession(cfg)

	query := readOne(input)

	if len(input.Name) == 0 {
		query = nil //If no name and namespace is provided then we have to retrieve all profiles thus we send nil into db query
	}

	err = mongo.Find(session, "AR", "aps", query, "name", &results)

	recordId, err = mongo.GetId(session, "AR", "aps", query)

	for i := range results {
		results[i].ID = recordId[i] //We add a record id value to the records we retrieved
	}

	output, err := createResponse(results) //Render the results into XML format

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=%s", "text/xml", "utf-8"))

	return output
}

func Create(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte {

	answer := ""

	//Authentication procedure
	if authentication.Authenticate(r.Header, cfg) {

		var name []string
		var namespace []string

		//Reading the json input
		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			panic(err)
		}

		input := AvailabilityProfileInput{}

		results := []AvailabilityProfileOutput{}

		//Unmarshalling the json input into byte form
		err = json.Unmarshal(reqBody, &input)

		//Making sure that no profile with the requested name and namespace combination already exists in the DB
		name = append(name, input.Name)

		namespace = append(namespace, input.Namespace)

		search := AvailabilityProfileSearch{
			name,
			namespace,
		}

		session := mongo.OpenSession(cfg)

		query := readOne(search)

		err = mongo.Find(session, "AR", "aps", query, "name", &results)

		if len(results) <= 0 {
			//If name-namespace combination is unique we insert the new record into mongo
			query := createOne(input)

			err = mongo.Insert(session, "AR", "aps", query)

			if err != nil {
				panic(err)
			}
			//Providing with the appropriate user response
			answer = "Availability Profile record successfully created"

		} else {
			answer = "An availability profile with that name already exists"
		}

	} else {
		answer = http.StatusText(403) //If wrong api key is passed we return FORBIDDEN http status
	}

	output, err := messageXML(answer) //Render the response into XML

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=%s", "text/xml", "utf-8"))

	return output
}

func Update(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte {

	answer := ""

	//Authentication procedure
	if authentication.Authenticate(r.Header, cfg) {

		//Extracting record id from url
		urlValues := r.URL.Path

		id := strings.Split(urlValues, "/")[4]

		//Reading the json input
		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			panic(err)
		}

		input := AvailabilityProfileInput{}

		//Unmarshalling the json input into byte form
		err = json.Unmarshal(reqBody, &input)

		session := mongo.OpenSession(cfg)

		//We update the record bassed on its unique id
		err = mongo.IdUpdate(session, "AR", "aps", id, input)

		if err != nil {
			answer = "No profile matching the requested id" //If not found we inform the user
		} else {
			answer = "Update successful" //We provide with the appropriate user response
		}
	} else {
		answer = http.StatusText(403) //If wrong api key is passed we return FORBIDDEN http status
	}

	output, err := messageXML(answer) //Render the response into XML

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=%s", "text/xml", "utf-8"))
	
	return output

}

func Delete(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte {

	answer := ""

	//Authentication procedure
	if authentication.Authenticate(r.Header, cfg) {

		//Extracting record id from url
		urlValues := r.URL.Path

		id := strings.Split(urlValues, "/")[4]

		session := mongo.OpenSession(cfg)

		//We remove the record bassed on its unique id
		err := mongo.IdRemove(session, "AR", "aps", id)

		if err != nil {
			answer = "No profile matching the requested id" //If not found we inform the user
		} else {
			answer = "Delete successful" //We provide with the appropriate user response
		}
	} else {
		answer = http.StatusText(403) //If wrong api key is passed we return FORBIDDEN http status
	}
	output, err := messageXML(answer)

	if err != nil {
		panic(err)
	}
	
	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=%s", "text/xml", "utf-8"))
	
	return output

}

func createResponse(results []AvailabilityProfileOutput) ([]byte, error) {
	
	output, err := CreateView(results)
	
	return output, err
	
}


