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
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/argoeu/argo-web-api/utils/authentication"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"
)

// List an availability profile
func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	//Read the search values
	urlValues := r.URL.Query()

	//Searchig is based on name and namespace
	input := AvailabilityProfileSearch{
		urlValues["name"],
		urlValues["namespace"],
	}

	results := []AvailabilityProfileOutput{}
	session, err := mongo.OpenSession(cfg.MongoDB)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	query := readOne(input)

	if len(input.Name) == 0 {
		query = nil //If no name and namespace is provided then we have to retrieve all profiles thus we send nil into db query
	}

	err = mongo.Find(session, cfg.MongoDB.Db, "aps", query, "_id", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	mongo.CloseSession(session)

	output, err = createView(results) //Render the results into XML format

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// Create an availability profile
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

	//Authentication procedure
	if authentication.Authenticate(r.Header, cfg) {

		session, err := mongo.OpenSession(cfg.MongoDB)

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		name := []string{}
		namespace := []string{}

		//Reading the json input
		reqBody, err := ioutil.ReadAll(r.Body)

		input := AvailabilityProfileInput{}
		results := []AvailabilityProfileOutput{}
		//Unmarshalling the json input into byte form
		err = json.Unmarshal(reqBody, &input)

		if err != nil {
			if err != nil {
				message = "Malformated json input data" // User provided malformed json input data
				output, err := messageXML(message)

				if err != nil {
					code = http.StatusInternalServerError
					return code, h, output, err
				}

				code = http.StatusBadRequest
				h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
				return code, h, output, err
			}
		}

		//Making sure that no profile with the requested name and namespace combination already exists in the DB
		name = append(name, input.Name)
		namespace = append(namespace, input.Namespace)

		search := AvailabilityProfileSearch{
			name,
			namespace,
		}

		query := readOne(search)
		err = mongo.Find(session, cfg.MongoDB.Db, "aps", query, "name", &results)

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		if len(results) <= 0 {
			//If name-namespace combination is unique we insert the new record into mongo
			query := createOne(input)
			err = mongo.Insert(session, cfg.MongoDB.Db, "aps", query)

			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}

			mongo.CloseSession(session)

			//Providing with the appropriate user response
			message = "Availability Profile record successfully created"
			output, err := messageXML(message) //Render the response into XML

			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}

			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err

		} else {
			message = "An availability profile with that name already exists"
			output, err := messageXML(message) //Render the response into XML

			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}

			code = http.StatusBadRequest
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}

	} else {
		output = []byte(http.StatusText(http.StatusUnauthorized))
		code = http.StatusUnauthorized //If wrong api key is passed we return UNAUTHORIZED http status
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

}

// Update an availability profile
func Update(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	message := ""

	//Authentication procedure
	if authentication.Authenticate(r.Header, cfg) {

		//Extracting record id from url
		urlValues := r.URL.Path
		id := strings.Split(urlValues, "/")[4]

		//Reading the json input
		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		input := AvailabilityProfileInput{}
		//Unmarshalling the json input into byte form
		err = json.Unmarshal(reqBody, &input)

		if err != nil {
			message = "Malformated json input data" // User provided malformed json input data
			output, err := messageXML(message)

			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}

			code = http.StatusBadRequest
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}

		session, err := mongo.OpenSession(cfg.MongoDB)

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		//We update the record bassed on its unique id
		err = mongo.IdUpdate(session, cfg.MongoDB.Db, "aps", id, input)

		mongo.CloseSession(session)

		if err != nil {
			message = "No profile matching the requested id" //If not found we inform the user
			output, err := messageXML(message)

			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}

			code = http.StatusBadRequest
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err

		} else {
			// Everything went fine and profile was deleted
			message = "Availability Profile was successfully updated"
			output, err := messageXML(message) //Render the response into XML

			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}

			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}
	} else {
		code = http.StatusUnauthorized
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		output = []byte(http.StatusText(http.StatusUnauthorized)) //If wrong api key is passed we return UNAUTHORIZED http status
		return code, h, output, err
	}

}

// Delete an availability profile
func Delete(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END
	message := ""

	//Authentication procedure
	if authentication.Authenticate(r.Header, cfg) {

		//Extracting record id from url
		urlValues := r.URL.Path
		id := strings.Split(urlValues, "/")[4]
		session, err := mongo.OpenSession(cfg.MongoDB)

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		//We remove the record bassed on its unique id
		err = mongo.IdRemove(session, cfg.MongoDB.Db, "aps", id)
		mongo.CloseSession(session)

		if err != nil {

			message = "No profile matching the requested id" //If not found we inform the user
			output, err := messageXML(message)

			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}

			code = http.StatusBadRequest
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		} else {

			// Everything went fine and profile was deleted
			message = "Availability Profile was successfully deleted"
			output, err := messageXML(message) //Render the response into XML

			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}

			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}
	} else {

		code = http.StatusUnauthorized
		output = []byte(http.StatusText(http.StatusUnauthorized)) //If wrong api key is passed we return UNAUTHORIZED http status
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

}
