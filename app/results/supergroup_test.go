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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gcfg.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SuperGroupAvailabilityTestSuite struct {
	suite.Suite
	cfg                       config.Config
	router                    *mux.Router
	confHandler               respond.ConfHandler
	tenantDbConf              config.MongoConfig
	tenantpassword            string
	tenantusername            string
	tenantstorename           string
	clientkey                 string
	respRecomputationsCreated string
}

// Setup the Test Environment
func (suite *SuperGroupAvailabilityTestSuite) SetupSuite() {

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
	    db = "ARGO_test_SuperGroup_availability"
	    `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.tenantDbConf.Db = "ARGO_test_SuperGroup_availability_tenant"
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
func (suite *SuperGroupAvailabilityTestSuite) SetupTest() {

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
					"roles":   []string{"viewer"},
				},
				bson.M{
					"name":    "King Joffrey",
					"email":   "g0dk1ng@kingslanding.gov",
					"api_key": "sansa <3",
					"roles":   []string{"viewer"},
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
					"roles":   []string{"viewer"},
				},
				bson.M{
					"name":    "Josh Plain",
					"email":   "P.Josh@egi.eu",
					"api_key": "itsamysterytoyou",
					"roles":   []string{"viewer"},
				},
			}})

	c = session.DB(suite.cfg.MongoDB.Db).C("roles")
	c.Insert(
		bson.M{
			"resource": "results.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "results.get",
			"roles":    []string{"editor", "viewer"},
		})
	// Seed database with recomputations
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
			"name":         "ST04",
			"supergroup":   "GROUP_B",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 30,
			"reliability":  100,
			"weight":       5344,
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
			"name":         "ST05",
			"supergroup":   "GROUP_B",
			"up":           0,
			"down":         0,
			"unknown":      1,
			"availability": 90,
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
			"date":         20150624,
			"name":         "ST05",
			"supergroup":   "GROUP_B",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 40,
			"reliability":  70,
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
			"date":         20150625,
			"name":         "ST05",
			"supergroup":   "GROUP_B",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 40,
			"reliability":  70,
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
					"type": "SITE",
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

// TestListSuperGroupAvailability test if daily results are returned correctly
func (suite *SuperGroupAvailabilityTestSuite) TestListSuperGroupAvailability() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/GROUP/GROUP_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	SuperGrouAvailabilityXML := ` <root>
   <group name="GROUP_A" type="GROUP">
     <results timestamp="2015-06-22" availability="68.13896116893515" reliability="50.413931144915935"></results>
     <results timestamp="2015-06-23" availability="75.36324059247399" reliability="80.8138510808647"></results>
   </group>
 </root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(SuperGrouAvailabilityXML, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/GROUP/GROUP_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	SuperGrouAvailabilityJSON := `{
   "results": [
     {
       "name": "GROUP_A",
       "type": "GROUP",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "68.13896116893515",
           "reliability": "50.413931144915935"
         },
         {
           "timestamp": "2015-06-23",
           "availability": "75.36324059247399",
           "reliability": "80.8138510808647"
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(SuperGrouAvailabilityJSON, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/GROUP/GROUP_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", "AWRONGKEY")
	request.Header.Set("Accept", "application/xml")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 401 Unauthorized code
	suite.Equal(401, response.Code, "Incorrect HTTP response code")

}

// TestListAllSuperGroupAvailability test if daily results are returned correctly
func (suite *SuperGroupAvailabilityTestSuite) TestListAllSuperGroupAvailability() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/GROUP?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	SuperGroupAvailabilityXML := ` <root>
   <group name="GROUP_A" type="GROUP">
     <results timestamp="2015-06-22" availability="68.13896116893515" reliability="50.413931144915935"></results>
     <results timestamp="2015-06-23" availability="75.36324059247399" reliability="80.8138510808647"></results>
   </group>
   <group name="GROUP_B" type="GROUP">
     <results timestamp="2015-06-23" availability="60.79234972677595" reliability="100"></results>
     <results timestamp="2015-06-24" availability="40" reliability="70"></results>
     <results timestamp="2015-06-25" availability="40" reliability="70"></results>
   </group>
 </root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(SuperGroupAvailabilityXML, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/GROUP?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	output := response.Body.String()

	SuperGroupAvailabilityJSON := `{
   "results": [
     {
       "name": "GROUP_A",
       "type": "GROUP",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "68.13896116893515",
           "reliability": "50.413931144915935"
         },
         {
           "timestamp": "2015-06-23",
           "availability": "75.36324059247399",
           "reliability": "80.8138510808647"
         }
       ]
     },
     {
       "name": "GROUP_B",
       "type": "GROUP",
       "results": [
         {
           "timestamp": "2015-06-23",
           "availability": "60.79234972677595",
           "reliability": "100"
         },
         {
           "timestamp": "2015-06-24",
           "availability": "40",
           "reliability": "70"
         },
         {
           "timestamp": "2015-06-25",
           "availability": "40",
           "reliability": "70"
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(SuperGroupAvailabilityJSON, output, "Response body mismatch")

}

// TestListAllSuperGroupAvailabilityCustom test if results are returned correctly for a custom period
func (suite *SuperGroupAvailabilityTestSuite) TestListAllSuperGroupAvailabilityCustom() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/GROUP?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=custom", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	SuperGroupAvailabilityJSON := `{
   "results": [
     {
       "name": "GROUP_A",
       "type": "GROUP",
       "results": [
         {
           "availability": "71.75110088070457",
           "reliability": "65.61389111289031"
         }
       ]
     },
     {
       "name": "GROUP_B",
       "type": "GROUP",
       "results": [
         {
           "availability": "46.930783242258656",
           "reliability": "80"
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(SuperGroupAvailabilityJSON, response.Body.String(), "Response body mismatch")
	fmt.Println(response.Body.String())

}

// TestListSuperGroupAvailabilityErrors tests if errors/exceptions are returned correctly
func (suite *SuperGroupAvailabilityTestSuite) TestListSuperGroupAvailabilityErrors() {

	reportErrorXML := ` <root>
   <message>The report with the name Report_B does not exist</message>
   <code>404</code>
 </root>`

	typeErrorXML := ` <root>
   <message>The report Report_A does not define any group type: supergroup</message>
   <code>404</code>
 </root>`

	typeError1XML := ` <root>
   <message>No results found for given query</message>
   <code>404</code>
 </root>`

	typeError1JSON := `{
   "message": "No results found for given query",
   "code": 404
 }`

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_B/supergroup?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 404 bad request code
	suite.Equal(404, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(reportErrorXML, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/supergroup?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")
	//request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	output := response.Body.String()

	// Check that we must have a 404 bad request code
	suite.Equal(404, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(typeErrorXML, output, "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/GROUP?start_time=2025-06-22T00:00:00Z&end_time=2025-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")
	//request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	output = response.Body.String()

	// Check that we must have a 404 not found code
	suite.Equal(404, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(typeError1XML, output, "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/GROUP?start_time=2025-06-22T00:00:00Z&end_time=2025-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	//request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	output = response.Body.String()

	// Check that we must have a 404 not found code
	suite.Equal(404, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(typeError1JSON, output, "Response body mismatch")

}

func (suite *SuperGroupAvailabilityTestSuite) TestOptionsSuperGroup() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/results/Report_A/GROUP", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v2/results/Report_A/GROUP/GROUP_A", strings.NewReader(""))

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

// TestStrictSlashSuperGroupResults test if not found responses are returned correctly
func (suite *SuperGroupAvailabilityTestSuite) TestStrictSlashSuperGroupResults() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/GROUP/?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")
	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	suite.Equal(404, response.Code, "Incorrect HTTP response code")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/GROUP/GROUP_A/?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")
	response = httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	suite.Equal(404, response.Code, "Incorrect HTTP response code")

}

//TearDownTest to tear down every test
func (suite *SuperGroupAvailabilityTestSuite) TearDownTest() {

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
func (suite *SuperGroupAvailabilityTestSuite) TearDownSuite() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

// TestRecompuptationsTestSuite is responsible for calling the tests
func TestSuperGroupResultsTestSuite(t *testing.T) {
	suite.Run(t, new(SuperGroupAvailabilityTestSuite))
}
