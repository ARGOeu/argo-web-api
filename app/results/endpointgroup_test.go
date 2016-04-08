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

package results

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

type endpointGroupAvailabilityTestSuite struct {
	suite.Suite
	cfg             config.Config
	router          *mux.Router
	confHandler     respond.ConfHandler
	tenantDbConf    config.MongoConfig
	tenantpassword  string
	tenantusername  string
	tenantstorename string
	clientkey       string
}

// Setup the Test Environment
func (suite *endpointGroupAvailabilityTestSuite) SetupSuite() {

	const testConfig = `
	    [server]
	    bindip = ""
	    port = 8080
	    maxprocs = 4
	    cache = false
	    lrucache = 700000000
	    gzip = true
	    [mongodb]
	    host = "127.0.0.1"
	    port = 27017
	    db = "ARGO_test_endpointGroup_availability"
	    `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.tenantDbConf.Db = "ARGO_test_endpointGroup_availability_tenant"
	suite.tenantDbConf.Password = "h4shp4ss"
	suite.tenantDbConf.Username = "dbuser"
	suite.tenantDbConf.Store = "ar"
	suite.clientkey = "secretkey"

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v2/results").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)
}

// This function runs before any test and setups the environment
func (suite *endpointGroupAvailabilityTestSuite) SetupTest() {

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
		bson.M{"name": "Westeros",
			"db_conf": []bson.M{

				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Westeros1",
				},
				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Westeros2",
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "John Snow",
					"email":   "J.Snow@brothers.wall",
					"api_key": "wh1t3_w@lk3rs",
				},
				bson.M{
					"name":    "King Joffrey",
					"email":   "g0dk1ng@kingslanding.gov",
					"api_key": "sansa <3",
				},
			}})
	c.Insert(
		bson.M{"name": "EGI",
			"db_conf": []bson.M{

				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": suite.tenantDbConf.Db,
					"username": suite.tenantDbConf.Username,
					"password": suite.tenantDbConf.Password,
				},
				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_wrong_db_endpointgrouavailability",
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "Joe Complex",
					"email":   "C.Joe@egi.eu",
					"api_key": suite.clientkey,
				},
				bson.M{
					"name":    "Josh Plain",
					"email":   "P.Josh@egi.eu",
					"api_key": "itsamysterytoyou",
				},
			}})

	// Seed tenant database with data
	c = session.DB(suite.tenantDbConf.Db).C("endpoint_group_ar")

	// Insert seed data
	c.Insert(
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "ST01",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 66.7,
			"reliability":  54.6,
			"weight":       5634,
			"tags": []bson.M{
				bson.M{
					"name":  "",
					"value": "",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "ST02",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 70,
			"reliability":  45,
			"weight":       4356,
			"tags": []bson.M{
				bson.M{
					"name":  "",
					"value": "",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150623,
			"name":         "ST01",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 100,
			"reliability":  100,
			"weight":       5634,
			"tags": []bson.M{
				bson.M{
					"name":  "",
					"value": "",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150623,
			"name":         "ST02",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 43.5,
			"reliability":  56,
			"weight":       4356,
			"tags": []bson.M{
				bson.M{
					"name":  "",
					"value": "",
				},
			},
		})

	c = session.DB(suite.tenantDbConf.Db).C("reports")

	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a49436",
		"info": bson.M{
			"name":        "Report_A",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "GROUP",
				"group": bson.M{
					"type": "SITES",
				},
			},
		},
		"profiles": []bson.M{
			bson.M{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			bson.M{
				"name":  "name1",
				"value": "value1"},
			bson.M{
				"name":  "name2",
				"value": "value2"},
		}})

}

// TestListEndpointGroupAvailability test if daily results are returned correctly
func (suite *endpointGroupAvailabilityTestSuite) TestListEndpointGroupAvailability() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/SITES/ST01?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	endpointGroupAvailabilityXML := ` <root>
   <group name="GROUP_A" type="GROUP">
     <group name="ST01" type="SITES">
       <results timestamp="2015-06-22" availability="66.7" reliability="54.6" unknown="0" uptime="1" downtime="0"></results>
       <results timestamp="2015-06-23" availability="100" reliability="100" unknown="0" uptime="1" downtime="0"></results>
     </group>
   </group>
 </root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointGroupAvailabilityXML, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/SITES/ST01?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	endpointGroupAvailabilityJSON := `{
   "root": [
     {
       "name": "GROUP_A",
       "type": "GROUP",
       "endpoints": [
         {
           "name": "ST01",
           "type": "SITES",
           "results": [
             {
               "timestamp": "2015-06-22",
               "availability": "66.7",
               "reliability": "54.6",
               "unknown": "0",
               "uptime": "1",
               "downtime": "0"
             },
             {
               "timestamp": "2015-06-23",
               "availability": "100",
               "reliability": "100",
               "unknown": "0",
               "uptime": "1",
               "downtime": "0"
             }
           ]
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointGroupAvailabilityJSON, response.Body.String(), "Response body mismatch")
	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/SITES/ST01?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", "AWRONGKEY")
	request.Header.Set("Accept", "application/xml")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	unauthorizedresponseXML := ` <root>
   <message>Unauthorized</message>
   <code>401</code>
 </root>`

	// Check that we must have a 401 Unauthorized code
	suite.Equal(401, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(unauthorizedresponseXML, response.Body.String(), "Response body mismatch")

}

// TestListAllEndpointGroupAvailability test if daily results are returned correctly
func (suite *endpointGroupAvailabilityTestSuite) TestListAllEndpointGroupAvailability() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/SITES?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	endpointGroupAvailabilityXML := ` <root>
   <group name="GROUP_A" type="GROUP">
     <group name="ST01" type="SITES">
       <results timestamp="2015-06-22" availability="66.7" reliability="54.6" unknown="0" uptime="1" downtime="0"></results>
       <results timestamp="2015-06-23" availability="100" reliability="100" unknown="0" uptime="1" downtime="0"></results>
     </group>
     <group name="ST02" type="SITES">
       <results timestamp="2015-06-22" availability="70" reliability="45" unknown="0" uptime="1" downtime="0"></results>
       <results timestamp="2015-06-23" availability="43.5" reliability="56" unknown="0" uptime="1" downtime="0"></results>
     </group>
   </group>
 </root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointGroupAvailabilityXML, response.Body.String(), "Response body mismatch")

}

// TestListEndpointGroupAvailabilityErrors tests if errors/exceptions are returned correctly
func (suite *endpointGroupAvailabilityTestSuite) TestListEndpointGroupAvailabilityErrors() {

	reportErrorXML := ` <root>
   <message>The report with the name Report_B does not exist</message>
   <code>404</code>
 </root>`

	typeError1XML := `<root>
 <status>
  <message>Bad Request</message>
  <code>400</code>
 </status>
 <errors>
  <error>
   <message>Endpoint Group type not in report</message>
   <code>400</code>
   <details>Endpoint Group type CLUSTERS not present in report Report_A. Try using SITES instead</details>
  </error>
 </errors>
</root>`

	typeError2XML := ` <root>
   <message>The report Report_A does not define any group type: CLUSTERS</message>
   <code>404</code>
 </root>`

	typeError3XML := ` <root>
   <message>No results found for given query</message>
   <code>404</code>
 </root>`

	typeError3JSON := `{
   "message": "No results found for given query",
   "code": 404
 }`

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_B/SITES?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 404 bad request code
	suite.Equal(404, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(reportErrorXML, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_B/GROUP/GROUP_A/SITES?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 404 bad request code
	suite.Equal(404, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(reportErrorXML, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/CLUSTERS?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 404 not found code
	suite.Equal(404, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(typeError2XML, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/GROUP/GROUP_A/CLUSTERS?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	output := response.Body.String()

	// Check that we must have a 400 bad request code
	suite.Equal(400, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(typeError1XML, output, "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/GROUP/GROUP_A/SITES?start_time=2025-06-22T00:00:00Z&end_time=2025-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	output = response.Body.String()

	// Check that we must have a 404 not found code
	suite.Equal(404, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(typeError3XML, output, "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/GROUP/GROUP_A/SITES?start_time=2025-06-22T00:00:00Z&end_time=2025-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	output = response.Body.String()

	// Check that we must have a 404 not found code
	suite.Equal(404, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(typeError3JSON, output, "Response body mismatch")

}

// TestOptionsEndpointGroup tests responses in case the OPTIONS http verb is used
func (suite *endpointGroupAvailabilityTestSuite) TestOptionsEndpointGroup() {

	request, _ := http.NewRequest("OPTIONS", "/api/v2/results/Report_A/SITES", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v2/results/Report_A/SITES/ST01", strings.NewReader(""))

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	headers = response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v2/results/Report_A/GROUP/GROUP_A/SITES", strings.NewReader(""))

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	headers = response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v2/results/Report_A/GROUP/GROUP_A/SITES/ST01", strings.NewReader(""))

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	headers = response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

}

// TestStrictSlashEndpointGroupResults test if not found responses are returned correctly
func (suite *endpointGroupAvailabilityTestSuite) TestStrictSlashEndpointGroupResults() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/SITES/?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")
	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	suite.Equal(404, response.Code, "Incorrect HTTP response code")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/SITES/ST01/?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")
	response = httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	suite.Equal(404, response.Code, "Incorrect HTTP response code")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/GROUP/GROUP_A/SITES/?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")
	response = httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	suite.Equal(404, response.Code, "Incorrect HTTP response code")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/GROUP/GROUP_A/SITES/SITES/ST01/?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")
	response = httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	suite.Equal(404, response.Code, "Incorrect HTTP response code")

}

//TearDownTest to tear down every test
func (suite *endpointGroupAvailabilityTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}

	tenantDB := session.DB(suite.tenantDbConf.Db)
	mainDB := session.DB(suite.cfg.MongoDB.Db)

	cols, err := tenantDB.CollectionNames()
	for _, col := range cols {
		tenantDB.C(col).RemoveAll(nil)
	}

	cols, err = mainDB.CollectionNames()
	for _, col := range cols {
		mainDB.C(col).RemoveAll(nil)
	}

}

//TearDownTest to tear down every test
func (suite *endpointGroupAvailabilityTestSuite) TearDownSuite() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

// TestEndpointGroupsTestSuite is responsible for calling the tests
func TestEndpointGroupsTestSuite(t *testing.T) {
	suite.Run(t, new(endpointGroupAvailabilityTestSuite))
}
