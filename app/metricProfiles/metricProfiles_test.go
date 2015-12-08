/*
 * Copyright (c) 2015 GRNET S.A.
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
 * or implied, of GRNET S.A.
 *
 */

package metricProfiles

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ARGOeu/argo-web-api/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/ARGOeu/argo-web-api/Godeps/_workspace/src/github.com/stretchr/testify/suite"
	"github.com/ARGOeu/argo-web-api/Godeps/_workspace/src/gopkg.in/gcfg.v1"
	"github.com/ARGOeu/argo-web-api/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/ARGOeu/argo-web-api/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
)

// This is a util. suite struct used in tests (see pkg "testify")
type MetricProfilesTestSuite struct {
	suite.Suite
	cfg                       config.Config
	router                    *mux.Router
	confHandler               respond.ConfHandler
	tenantDbConf              config.MongoConfig
	clientkey                 string
	respRecomputationsCreated string
	respUnauthorized          string
}

// Setup the Test Environment
// This function runs before any test and setups the environment
func (suite *MetricProfilesTestSuite) SetupTest() {

	const testConfig = `
    [server]
    bindip = ""
    port = 8080
    maxprocs = 4
    cache = false
    lrucache = 700000000
    gzip = true
	reqsizelimit = 1073741824

    [mongodb]
    host = "127.0.0.1"
    port = 27017
    db = "AR_test_mprof"
    `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.respUnauthorized = "Unauthorized"
	suite.tenantDbConf = config.MongoConfig{
		Host:     "localhost",
		Port:     27017,
		Db:       "AR_test_metric_profiles_tenant",
		Password: "pass",
		Username: "dbuser",
		Store:    "ar",
	}
	suite.clientkey = "123456"

	suite.confHandler = respond.ConfHandler{suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v2").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)

	// seed mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Seed database with tenants
	//TODO: move tests to
	c := session.DB(suite.cfg.MongoDB.Db).C("tenants")
	c.Insert(
		bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"info": bson.M{
				"name":    "GUARDIANS",
				"email":   "email@something2",
				"website": "www.gotg.com",
				"created": "2015-10-20 02:08:04",
				"updated": "2015-10-20 02:08:04"},
			"db_conf": []bson.M{

				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_FOO",
				},
				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_FOO",
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "user1",
					"email":   "user1@email.com",
					"api_key": "USER1KEY",
				},
				bson.M{
					"name":    "user2",
					"email":   "user2@email.com",
					"api_key": "USER2KEY",
				},
			}})
	c.Insert(
		bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50d",
			"info": bson.M{
				"name":    "AVENGERS",
				"email":   "email@something2",
				"website": "www.gotg.com",
				"created": "2015-10-20 02:08:04",
				"updated": "2015-10-20 02:08:04"},
			"db_conf": []bson.M{

				bson.M{
					// "store":    "ar",
					"server":   suite.tenantDbConf.Host,
					"port":     suite.tenantDbConf.Port,
					"database": suite.tenantDbConf.Db,
					"username": suite.tenantDbConf.Username,
					"password": suite.tenantDbConf.Password,
				},
				bson.M{
					"server":   suite.tenantDbConf.Host,
					"port":     suite.tenantDbConf.Port,
					"database": suite.tenantDbConf.Db,
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "user3",
					"email":   "user3@email.com",
					"api_key": suite.clientkey,
				},
				bson.M{
					"name":    "user4",
					"email":   "user4@email.com",
					"api_key": "USER4KEY",
				},
			}})
	// Seed database with metric profiles
	c = session.DB(suite.tenantDbConf.Db).C("metric_profiles")
	c.Insert(
		bson.M{
			"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			"name": "ch.cern.SAM.ROC_CRITICAL",
			"services": []bson.M{
				bson.M{"service": "CREAM-CE",
					"metrics": []string{
						"emi.cream.CREAMCE-JobSubmit",
						"emi.wn.WN-Bi",
						"emi.wn.WN-Csh",
						"emi.wn.WN-SoftVer"},
				},
				bson.M{"service": "SRMv2",
					"metrics": []string{"hr.srce.SRM2-CertLifetime",
						"org.sam.SRM-Del",
						"org.sam.SRM-Get",
						"org.sam.SRM-GetSURLs",
						"org.sam.SRM-GetTURLs",
						"org.sam.SRM-Ls",
						"org.sam.SRM-LsDir",
						"org.sam.SRM-Put"},
				},
			},
		})
	c.Insert(
		bson.M{
			"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"name": "ch.cern.SAM.ROC",
			"services": []bson.M{
				bson.M{"service": "CREAM-CE",
					"metrics": []string{
						"emi.cream.CREAMCE-JobSubmit",
						"emi.wn.WN-Bi",
						"emi.wn.WN-Csh",
						"hr.srce.CADist-Check",
						"hr.srce.CREAMCE-CertLifetime",
						"emi.wn.WN-SoftVer"},
				},
				bson.M{"service": "SRMv2",
					"metrics": []string{"hr.srce.SRM2-CertLifetime",
						"org.sam.SRM-Del",
						"org.sam.SRM-Get",
						"org.sam.SRM-GetSURLs",
						"org.sam.SRM-GetTURLs",
						"org.sam.SRM-Ls",
						"org.sam.SRM-LsDir",
						"org.sam.SRM-Put"},
				},
			},
		})

}

func (suite *MetricProfilesTestSuite) TestList() {

	request, _ := http.NewRequest("GET", "/api/v2/metric_profiles", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "name": "ch.cern.SAM.ROC",
   "services": [
    {
     "service": "CREAM-CE",
     "metrics": [
      "emi.cream.CREAMCE-JobSubmit",
      "emi.wn.WN-Bi",
      "emi.wn.WN-Csh",
      "hr.srce.CADist-Check",
      "hr.srce.CREAMCE-CertLifetime",
      "emi.wn.WN-SoftVer"
     ]
    },
    {
     "service": "SRMv2",
     "metrics": [
      "hr.srce.SRM2-CertLifetime",
      "org.sam.SRM-Del",
      "org.sam.SRM-Get",
      "org.sam.SRM-GetSURLs",
      "org.sam.SRM-GetTURLs",
      "org.sam.SRM-Ls",
      "org.sam.SRM-LsDir",
      "org.sam.SRM-Put"
     ]
    }
   ]
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "name": "ch.cern.SAM.ROC_CRITICAL",
   "services": [
    {
     "service": "CREAM-CE",
     "metrics": [
      "emi.cream.CREAMCE-JobSubmit",
      "emi.wn.WN-Bi",
      "emi.wn.WN-Csh",
      "emi.wn.WN-SoftVer"
     ]
    },
    {
     "service": "SRMv2",
     "metrics": [
      "hr.srce.SRM2-CertLifetime",
      "org.sam.SRM-Del",
      "org.sam.SRM-Get",
      "org.sam.SRM-GetSURLs",
      "org.sam.SRM-GetTURLs",
      "org.sam.SRM-Ls",
      "org.sam.SRM-LsDir",
      "org.sam.SRM-Put"
     ]
    }
   ]
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricProfileJSON, output, "Response body mismatch")

}

func (suite *MetricProfilesTestSuite) TestListQueryName() {

	request, _ := http.NewRequest("GET", "/api/v2/metric_profiles?name=ch.cern.SAM.ROC", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "name": "ch.cern.SAM.ROC",
   "services": [
    {
     "service": "CREAM-CE",
     "metrics": [
      "emi.cream.CREAMCE-JobSubmit",
      "emi.wn.WN-Bi",
      "emi.wn.WN-Csh",
      "hr.srce.CADist-Check",
      "hr.srce.CREAMCE-CertLifetime",
      "emi.wn.WN-SoftVer"
     ]
    },
    {
     "service": "SRMv2",
     "metrics": [
      "hr.srce.SRM2-CertLifetime",
      "org.sam.SRM-Del",
      "org.sam.SRM-Get",
      "org.sam.SRM-GetSURLs",
      "org.sam.SRM-GetTURLs",
      "org.sam.SRM-Ls",
      "org.sam.SRM-LsDir",
      "org.sam.SRM-Put"
     ]
    }
   ]
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricProfileJSON, output, "Response body mismatch")

}

func (suite *MetricProfilesTestSuite) TestListOneNotFound() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Not Found",
  "code": "404",
  "details": "item with the specific ID was not found on the server"
 }
}`

	request, _ := http.NewRequest("GET", "/api/v2/metric_profiles/wrong-id", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(404, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *MetricProfilesTestSuite) TestListOne() {

	request, _ := http.NewRequest("GET", "/api/v2/metric_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "name": "ch.cern.SAM.ROC_CRITICAL",
   "services": [
    {
     "service": "CREAM-CE",
     "metrics": [
      "emi.cream.CREAMCE-JobSubmit",
      "emi.wn.WN-Bi",
      "emi.wn.WN-Csh",
      "emi.wn.WN-SoftVer"
     ]
    },
    {
     "service": "SRMv2",
     "metrics": [
      "hr.srce.SRM2-CertLifetime",
      "org.sam.SRM-Del",
      "org.sam.SRM-Get",
      "org.sam.SRM-GetSURLs",
      "org.sam.SRM-GetTURLs",
      "org.sam.SRM-Ls",
      "org.sam.SRM-LsDir",
      "org.sam.SRM-Put"
     ]
    }
   ]
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricProfileJSON, output, "Response body mismatch")

}

func (suite *MetricProfilesTestSuite) TestCreateBadJson() {

	jsonInput := `{
  "name": "test_profile",
  "services": [
    {
      "service": "Service-A",
      "metrics": [
        "metric.A.1",
        "metric.A.2",
        "metric.A.3",
        "metric.A.4"
      ]
    },
    {
      "service": "Service-B",
      "metrics": [
        "metric.B.1",
        "metric.B.2"
    `

	jsonOutput := `{
 "status": {
  "message": "Bad Request",
  "code": "400",
  "details": "Request Body contains malformed JSON, thus rendering the Request Bad"
 }
}`

	request, _ := http.NewRequest("POST", "/api/v2/metric_profiles", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(400, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *MetricProfilesTestSuite) TestCreate() {

	jsonInput := `{
  "name": "test_profile",
  "services": [
    {
      "service": "Service-A",
      "metrics": [
        "metric.A.1",
        "metric.A.2",
        "metric.A.3",
        "metric.A.4"
      ]
    },
    {
      "service": "Service-B",
      "metrics": [
        "metric.B.1",
        "metric.B.2"
      ]
    }
  ]
}`

	jsonOutput := `{
 "status": {
  "message": "Metric Profile successfully created",
  "code": "201"
 },
 "data": {
  "id": "{{id}}",
  "links": {
   "self": "https:///api/v2/metric_profiles/{{id}}"
  }
 }
}`

	jsonCreated := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "{{id}}",
   "name": "test_profile",
   "services": [
    {
     "service": "Service-A",
     "metrics": [
      "metric.A.1",
      "metric.A.2",
      "metric.A.3",
      "metric.A.4"
     ]
    },
    {
     "service": "Service-B",
     "metrics": [
      "metric.B.1",
      "metric.B.2"
     ]
    }
   ]
  }
 ]
}`

	request, _ := http.NewRequest("POST", "/api/v2/metric_profiles", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(201, code, "Internal Server Error")
	// Compare the expected and actual json response

	// Grab id from mongodb
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	defer session.Close()
	if err != nil {
		panic(err)
	}
	// Retrieve id from database
	var result map[string]interface{}
	c := session.DB(suite.tenantDbConf.Db).C("metric_profiles")
	c.Find(bson.M{"name": "test_profile"}).One(&result)
	id := result["id"].(string)

	// Apply id to output template and check
	suite.Equal(strings.Replace(jsonOutput, "{{id}}", id, 2), output, "Response body mismatch")

	// Check that actually the item has been created
	// Call List one with the specific id
	request2, _ := http.NewRequest("GET", "/api/v2/metric_profiles/"+id, strings.NewReader(jsonInput))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)

	code2 := response2.Code
	output2 := response2.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(200, code2, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(strings.Replace(jsonCreated, "{{id}}", id, 2), output2, "Response body mismatch")

}

func (suite *MetricProfilesTestSuite) TestUpdateBadJson() {

	jsonInput := `{
  "name": "test_profile",
  "services": [
    {
      "service": "Service-A",
      "metrics": [
        "metric.A.1",
        "metric.A.2",
        "metric.A.3",
        "metric.A.4"
      ]
    },
    {
      "service": "Service-B",
      "metrics": [
        "metric.B.1",
        "metric.B.2"
    `

	jsonOutput := `{
 "status": {
  "message": "Bad Request",
  "code": "400",
  "details": "Request Body contains malformed JSON, thus rendering the Request Bad"
 }
}`

	request, _ := http.NewRequest("PUT", "/api/v2/metric_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(400, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *MetricProfilesTestSuite) TestUpdateNotFound() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Not Found",
  "code": "404",
  "details": "item with the specific ID was not found on the server"
 }
}`

	request, _ := http.NewRequest("PUT", "/api/v2/metric_profiles/wrong-id", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(404, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *MetricProfilesTestSuite) TestUpdate() {

	jsonInput := `{
  "name": "test_profile",
  "services": [
    {
      "service": "Service-A",
      "metrics": [
        "metric.A.1",
        "metric.A.2",
        "metric.A.3",
        "metric.A.4"
      ]
    },
    {
      "service": "Service-B",
      "metrics": [
        "metric.B.1",
        "metric.B.2"
      ]
    }
  ]
}`

	jsonOutput := `{
 "status": {
  "message": "Metric Profile successfully updated",
  "code": "200"
 }
}`

	jsonUpdated := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "name": "test_profile",
   "services": [
    {
     "service": "Service-A",
     "metrics": [
      "metric.A.1",
      "metric.A.2",
      "metric.A.3",
      "metric.A.4"
     ]
    },
    {
     "service": "Service-B",
     "metrics": [
      "metric.B.1",
      "metric.B.2"
     ]
    }
   ]
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/metric_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response

	// Apply id to output template and check
	suite.Equal(jsonOutput, output, "Response body mismatch")

	// Check that the item has actually updated
	// run a list specific
	request2, _ := http.NewRequest("GET", "/api/v2/metric_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)

	code2 := response2.Code
	output2 := response2.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(200, code2, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonUpdated, output2, "Response body mismatch")

}

func (suite *MetricProfilesTestSuite) TestDeleteNotFound() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Not Found",
  "code": "404",
  "details": "item with the specific ID was not found on the server"
 }
}`

	request, _ := http.NewRequest("DELETE", "/api/v2/metric_profiles/wrong-id", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(404, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *MetricProfilesTestSuite) TestDelete() {

	request, _ := http.NewRequest("DELETE", "/api/v2/metric_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Metric Profile Successfully Deleted",
  "code": "200"
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricProfileJSON, output, "Response body mismatch")

	// check that the element has actually been Deleted
	// connect to mongodb
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	defer session.Close()
	if err != nil {
		panic(err)
	}
	// try to retrieve item
	var result map[string]interface{}
	c := session.DB(suite.tenantDbConf.Db).C("metric_profiles")
	err = c.Find(bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"}).One(&result)

	suite.NotEqual(err, nil, "No not found error")
	suite.Equal(err.Error(), "not found", "No not found error")
}

//TearDownTest to tear down every test
func (suite *MetricProfilesTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

func TestMetricProfilesTestSuite(t *testing.T) {
	suite.Run(t, new(MetricProfilesTestSuite))
}
