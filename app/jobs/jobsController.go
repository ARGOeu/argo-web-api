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

package jobs

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

// Create function is used to implement the create job request.
// The request is an http POST request with the job description
// provided as json structure in the request body
func Create(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Authenticate user's api key and find corresponding tenant
	tenantDbConf, err := authentication.AuthenticateTenant(r.Header, cfg)

	// if authentication procedure fails then
	// return unauthorized http status
	if err != nil {

		output = []byte(http.StatusText(http.StatusUnauthorized))
		//If wrong api key is passed we return UNAUTHORIZED http status
		code = http.StatusUnauthorized
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	//Reading the json input from the request body
	reqBody, err := ioutil.ReadAll(r.Body)
	input := Job{}
	//Unmarshalling the json input into byte form
	err = json.Unmarshal(reqBody, &input)

	// Check if json body is malformed
	if err != nil {
		if err != nil {
			// Msg in xml style, to notify for malformed json
			output, err := messageXML("Malformated json input data")

			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}

			code = http.StatusBadRequest
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}
	}

	// Try to open the mongo session
	session, err := mongo.OpenSession(tenantDbConf)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Prepare structure for storing query results
	results := []Job{}

	// Check if job with the same name exists in datastore
	query := searchName(input.Name)
	err = mongo.Find(session, tenantDbConf.Db, "jobs", query, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// If results are returned for the specific name
	// then we already have an existing job and we must
	// abort creation notifing the user
	if len(results) > 0 {
		// Name was found so print the error message in xml
		output, err = messageXML("Job with the same name already exists")

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		code = http.StatusBadRequest
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err

	}

	// If no job exists with this name create a new one
	query = createJob(input)
	err = mongo.Insert(session, tenantDbConf.Db, "jobs", query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Notify user that the job has been created. In xml style
	output, err = messageXML("Job was successfully created")

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err

}

// List function that implements the http GET request that retrieves
// all avaiable job information
func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Authenticate user's api key and find corresponding tenant
	tenantDbConf, err := authentication.AuthenticateTenant(r.Header, cfg)

	// if authentication procedure fails then
	// return unauthorized http status
	if err != nil {

		output = []byte(http.StatusText(http.StatusUnauthorized))
		//If wrong api key is passed we return UNAUTHORIZED http status
		code = http.StatusUnauthorized
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	// Try to open the mongo session
	session, err := mongo.OpenSession(cfg.MongoDB)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create structure for storing query results
	results := []Job{}
	// Query tenant collection for all available documents.
	// nil query param == match everything
	err = mongo.Find(session, tenantDbConf.Db, "jobs", nil, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// After successfully retrieving the db results
	// call the createView function to render them into idented xml
	output, err = createView(results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// ListOne function that implements the http GET request that retrieves
// all avaiable job information
func ListOne(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Authenticate user's api key and find corresponding tenant
	tenantDbConf, err := authentication.AuthenticateTenant(r.Header, cfg)

	// if authentication procedure fails then
	// return unauthorized http status
	if err != nil {

		output = []byte(http.StatusText(http.StatusUnauthorized))
		//If wrong api key is passed we return UNAUTHORIZED http status
		code = http.StatusUnauthorized
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	//Extracting urlvar "name" from url path
	urlValues := r.URL.Path
	nameFromURL := strings.Split(urlValues, "/")[4]

	// Try to open the mongo session
	session, err := mongo.OpenSession(tenantDbConf)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create structure for storing query results
	results := []Job{}
	// Create a simple query object to query by name
	query := searchName(nameFromURL)
	// Query collection tenants for the specific tenant name
	err = mongo.Find(session, tenantDbConf.Db, "jobs", query, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// If query returned zero result then no tenant matched this name,
	// abort and notify user accordingly
	if len(results) == 0 {

		output, err := messageXML("Job not found")

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		code = http.StatusBadRequest
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	// After successfully retrieving the db results
	// call the createView function to render them into idented xml
	output, err = createView(results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// Update function used to implement update job request.
// This is an http PUT request that gets a specific job's name
// as a urlvar parameter input and a json structure in the request
// body in order to update the datastore document for the specific
// job
func Update(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Authenticate user's api key and find corresponding tenant
	tenantDbConf, err := authentication.AuthenticateTenant(r.Header, cfg)

	// if authentication procedure fails then
	// return unauthorized http status
	if err != nil {

		output = []byte(http.StatusText(http.StatusUnauthorized))
		//If wrong api key is passed we return UNAUTHORIZED http status
		code = http.StatusUnauthorized
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	//Extracting job name from url
	urlValues := r.URL.Path
	nameFromURL := strings.Split(urlValues, "/")[4]

	//Reading the json input
	reqBody, err := ioutil.ReadAll(r.Body)

	input := Job{}
	//Unmarshalling the json input into byte form
	err = json.Unmarshal(reqBody, &input)

	if err != nil {
		if err != nil {
			// User provided malformed json input data
			output, err := messageXML("Malformated json input data")

			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}

			code = http.StatusBadRequest
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}
	}

	// Try to open the mongo session
	session, err := mongo.OpenSession(tenantDbConf)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// We search by name and update
	query := searchName(nameFromURL)
	err = mongo.Update(session, tenantDbConf.Db, "jobs", query, input)

	if err != nil {

		if err.Error() != "not found" {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
		//Render the response into XML
		output, err = messageXML("Job not found")

	} else {
		//Render the response into XML
		output, err = messageXML("Job was successfully updated")
	}

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err

}

// Delete function used to implement remove job request
func Delete(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Authenticate user's api key and find corresponding tenant
	tenantDbConf, err := authentication.AuthenticateTenant(r.Header, cfg)

	// if authentication procedure fails then
	// return unauthorized http status
	if err != nil {

		output = []byte(http.StatusText(http.StatusUnauthorized))
		//If wrong api key is passed we return UNAUTHORIZED http status
		code = http.StatusUnauthorized
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	//Extracting record id from url
	urlValues := r.URL.Path
	nameFromURL := strings.Split(urlValues, "/")[4]

	// Try to open the mongo session
	session, err := mongo.OpenSession(tenantDbConf)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// We search by name and delete the document in db
	query := searchName(nameFromURL)
	info, err := mongo.Remove(session, tenantDbConf.Db, "jobs", query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// info.Removed > 0 means that many documents have been removed
	// If deletion took place we notify user accordingly.
	// Else we notify that no tenant matched the specific name
	if info.Removed > 0 {
		output, err = messageXML("Job was successfully deleted")
	} else {
		output, err = messageXML("Job not found")
	}
	//Render the response into XML
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err

}
