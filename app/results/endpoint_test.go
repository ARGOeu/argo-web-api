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
	"io/ioutil"
	"log"
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

type endpointAvailabilityTestSuite struct {
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
func (suite *endpointAvailabilityTestSuite) SetupSuite() {

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
	 db = "ARGO_test_endpoint_availability"
	 `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.tenantDbConf.Db = "ARGO_test_endpoint_availability_tenant"
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
func (suite *endpointAvailabilityTestSuite) SetupTest() {

	log.SetOutput(ioutil.Discard)

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
					"database": "argo_wrong_db_serviceflavoravailability",
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

	c = session.DB(suite.tenantDbConf.Db).C("endpoint_ar")

	// Insert seed data
	c.Insert(
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "e01",
			"supergroup":   "ST01",
			"service":      "service_a",
			"up":           0.98264,
			"down":         0,
			"unknown":      0,
			"availability": 98.26389,
			"reliability":  98.26389,
			"tags": []bson.M{
				bson.M{
					"name":  "production",
					"value": "Y",
				},
			},
			"info": bson.M{
				"Url": "https://foo.example.url",
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "e02",
			"supergroup":   "ST01",
			"service":      "service_a",
			"up":           0.96875,
			"down":         0,
			"unknown":      0,
			"availability": 96.875,
			"reliability":  96.875,
			"tags": []bson.M{
				bson.M{
					"name":  "production",
					"value": "Y",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "e03",
			"supergroup":   "ST02",
			"service":      "service_b",
			"up":           0.96875,
			"down":         0,
			"unknown":      0,
			"availability": 96.875,
			"reliability":  96.875,
			"tags": []bson.M{
				bson.M{
					"name":  "production",
					"value": "Y",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "e03",
			"supergroup":   "ST02",
			"service":      "service_x",
			"up":           0.96875,
			"down":         0,
			"unknown":      0,
			"availability": 96.875,
			"reliability":  96.875,
			"tags": []bson.M{
				bson.M{
					"name":  "production",
					"value": "Y",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "e03",
			"supergroup":   "STX",
			"service":      "service_b",
			"up":           0.96875,
			"down":         0,
			"unknown":      0,
			"availability": 96.875,
			"reliability":  96.875,
			"tags": []bson.M{
				bson.M{
					"name":  "production",
					"value": "Y",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150623,
			"name":         "e01",
			"supergroup":   "ST01",
			"service":      "service_a",
			"up":           0.53472,
			"down":         0.33333,
			"unknown":      0.01042,
			"availability": 54.03509,
			"reliability":  81.48148,
			"tags": []bson.M{
				bson.M{
					"name":  "production",
					"value": "Y",
				},
			},
			"info": bson.M{
				"Url": "https://foo.example.url",
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150623,
			"name":         "e02",
			"supergroup":   "ST01",
			"service":      "service_a",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 100,
			"reliability":  100,
			"tags": []bson.M{
				bson.M{
					"name":  "production",
					"value": "Y",
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

// TestListEndpointAvailabilityMonthly tests if monthly results are returned correctly
func (suite *endpointAvailabilityTestSuite) TestListEndpointAvailabilityMonthly() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/services/service_a/endpoints?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	responseBody := response.Body.String()
	endpointAvailabilityXML := ` <root>
   <group name="ST01" type="SITE">
     <group name="service_a" type="service">
       <group name="e01" type="endpoint">
         <results timestamp="2015-06" availability="76.26534166743393" reliability="91.61418757296076" unknown="0.00521" uptime="0.75868" downtime="0.166665"></results>
       </group>
       <group name="e02" type="endpoint">
         <results timestamp="2015-06" availability="98.43749901562502" reliability="98.43749901562502" unknown="0" uptime="0.984375" downtime="0"></results>
       </group>
     </group>
   </group>
 </root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointAvailabilityXML, responseBody, "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/services/service_a/endpoints?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	responseBody = response.Body.String()

	endpointAvailabilityJSON := `{
   "results": [
     {
       "name": "ST01",
       "type": "SITE",
       "serviceflavors": [
         {
           "name": "service_a",
           "type": "service",
           "endpoints": [
             {
               "name": "e01",
               "type": "endpoint",
               "info": {
                 "Url": "https://foo.example.url"
               },
               "results": [
                 {
                   "timestamp": "2015-06",
                   "availability": "76.26534166743393",
                   "reliability": "91.61418757296076",
                   "unknown": "0.00521",
                   "uptime": "0.75868",
                   "downtime": "0.166665"
                 }
               ]
             },
             {
               "name": "e02",
               "type": "endpoint",
               "results": [
                 {
                   "timestamp": "2015-06",
                   "availability": "98.43749901562502",
                   "reliability": "98.43749901562502",
                   "unknown": "0",
                   "uptime": "0.984375",
                   "downtime": "0"
                 }
               ]
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
	suite.Equal(endpointAvailabilityJSON, responseBody, "Response body mismatch")

	// Test the quick path to site endpoint a/r
	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/endpoints?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	responseBody = response.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointAvailabilityJSON, responseBody, "Response body mismatch")

}

func (suite *endpointAvailabilityTestSuite) TestListEndpointAvailabilityCustom() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/services/service_a/endpoints?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z&granularity=custom", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	responseBody := response.Body.String()
	endpointAvailabilityJSON := `{
   "results": [
     {
       "name": "ST01",
       "type": "SITE",
       "serviceflavors": [
         {
           "name": "service_a",
           "type": "service",
           "endpoints": [
             {
               "name": "e01",
               "type": "endpoint",
               "info": {
                 "Url": "https://foo.example.url"
               },
               "results": [
                 {
                   "availability": "76.26534166743393",
                   "reliability": "91.61418757296076",
                   "unknown": "0.00521",
                   "uptime": "0.75868",
                   "downtime": "0.166665"
                 }
               ]
             },
             {
               "name": "e02",
               "type": "endpoint",
               "results": [
                 {
                   "availability": "98.43749901562502",
                   "reliability": "98.43749901562502",
                   "unknown": "0",
                   "uptime": "0.984375",
                   "downtime": "0"
                 }
               ]
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
	suite.Equal(endpointAvailabilityJSON, responseBody, "Response body mismatch")

}
func (suite *endpointAvailabilityTestSuite) TestFlatAllEndpoints() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	resultsJSON := `{
   "results": [
     {
       "name": "e01",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "info": {
         "Url": "https://foo.example.url"
       },
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "98.26389",
           "reliability": "98.26389",
           "unknown": "0",
           "uptime": "0.98264",
           "downtime": "0"
         },
         {
           "timestamp": "2015-06-23",
           "availability": "54.03509",
           "reliability": "81.48148",
           "unknown": "0.01042",
           "uptime": "0.53472",
           "downtime": "0.33333"
         }
       ]
     },
     {
       "name": "e02",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
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
     },
     {
       "name": "e03",
       "service": "service_b",
       "supergroup": "ST02",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     },
     {
       "name": "e03",
       "service": "service_b",
       "supergroup": "STX",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     },
     {
       "name": "e03",
       "service": "service_x",
       "supergroup": "ST02",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(resultsJSON, response.Body.String(), "Response body mismatch")

}

func (suite *endpointAvailabilityTestSuite) TestFlatAllEndpointsPaginated() {

	type expReq struct {
		method string
		url    string
		code   int
		result string
	}

	expReqs := []expReq{
		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily&pageSize=1",
			result: `{
   "results": [
     {
       "name": "e01",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "info": {
         "Url": "https://foo.example.url"
       },
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "98.26389",
           "reliability": "98.26389",
           "unknown": "0",
           "uptime": "0.98264",
           "downtime": "0"
         }
       ]
     }
   ],
   "nextPageToken": "MQ==",
   "pageSize": 1
 }`,
			code: 200,
		},
		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily&pageSize=1&nextPageToken=MQ==",
			result: `{
   "results": [
     {
       "name": "e01",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "info": {
         "Url": "https://foo.example.url"
       },
       "results": [
         {
           "timestamp": "2015-06-23",
           "availability": "54.03509",
           "reliability": "81.48148",
           "unknown": "0.01042",
           "uptime": "0.53472",
           "downtime": "0.33333"
         }
       ]
     }
   ],
   "nextPageToken": "Mg==",
   "pageSize": 1
 }`,
			code: 200,
		},

		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily&pageSize=1&nextPageToken=Mg==",
			result: `{
   "results": [
     {
       "name": "e02",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     }
   ],
   "nextPageToken": "Mw==",
   "pageSize": 1
 }`,
			code: 200,
		},

		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily&pageSize=1&nextPageToken=Mw==",
			result: `{
   "results": [
     {
       "name": "e02",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "results": [
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
   ],
   "nextPageToken": "NA==",
   "pageSize": 1
 }`,
			code: 200,
		},

		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily&pageSize=1&nextPageToken=NA==",
			result: `{
   "results": [
     {
       "name": "e03",
       "service": "service_b",
       "supergroup": "ST02",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     }
   ],
   "nextPageToken": "NQ==",
   "pageSize": 1
 }`,
			code: 200,
		},

		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily&pageSize=1&nextPageToken=NQ==",
			result: `{
   "results": [
     {
       "name": "e03",
       "service": "service_b",
       "supergroup": "STX",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     }
   ],
   "nextPageToken": "Ng==",
   "pageSize": 1
 }`,
			code: 200,
		},

		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily&pageSize=1&nextPageToken=Ng==",
			result: `{
   "results": [
     {
       "name": "e03",
       "service": "service_x",
       "supergroup": "ST02",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     }
   ],
   "pageSize": 1
 }`,
			code: 200,
		},
	}

	for _, expReq := range expReqs {
		request, _ := http.NewRequest(expReq.method, expReq.url, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")

		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		suite.Equal(expReq.code, response.Code, "Incorrect HTTP response code")
		// Compare the expected and actual xml response
		suite.Equal(expReq.result, response.Body.String(), "Response body mismatch")

	}

}

func (suite *endpointAvailabilityTestSuite) TestFlatAllEndpointsPaginated3() {

	type expReq struct {
		method string
		url    string
		code   int
		result string
	}

	expReqs := []expReq{
		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily&pageSize=3",
			result: `{
   "results": [
     {
       "name": "e01",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "info": {
         "Url": "https://foo.example.url"
       },
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "98.26389",
           "reliability": "98.26389",
           "unknown": "0",
           "uptime": "0.98264",
           "downtime": "0"
         },
         {
           "timestamp": "2015-06-23",
           "availability": "54.03509",
           "reliability": "81.48148",
           "unknown": "0.01042",
           "uptime": "0.53472",
           "downtime": "0.33333"
         }
       ]
     },
     {
       "name": "e02",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     }
   ],
   "nextPageToken": "Mw==",
   "pageSize": 3
 }`,
			code: 200,
		},
		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily&pageSize=3&nextPageToken=Mw==",
			result: `{
   "results": [
     {
       "name": "e02",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-23",
           "availability": "100",
           "reliability": "100",
           "unknown": "0",
           "uptime": "1",
           "downtime": "0"
         }
       ]
     },
     {
       "name": "e03",
       "service": "service_b",
       "supergroup": "ST02",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     },
     {
       "name": "e03",
       "service": "service_b",
       "supergroup": "STX",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     }
   ],
   "nextPageToken": "Ng==",
   "pageSize": 3
 }`,
			code: 200,
		},
		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily&pageSize=3&nextPageToken=Ng==",
			result: `{
   "results": [
     {
       "name": "e03",
       "service": "service_x",
       "supergroup": "ST02",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     }
   ],
   "pageSize": 3
 }`,
			code: 200,
		},
	}

	for _, expReq := range expReqs {
		request, _ := http.NewRequest(expReq.method, expReq.url, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")

		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		suite.Equal(expReq.code, response.Code, "Incorrect HTTP response code")
		// Compare the expected and actual xml response
		suite.Equal(expReq.result, response.Body.String(), "Response body mismatch")

	}

}

func (suite *endpointAvailabilityTestSuite) TestMonthlyFlatAllEndpointsPaginatedMonthly() {

	type expReq struct {
		method string
		url    string
		code   int
		result string
	}

	expReqs := []expReq{
		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=monthly&pageSize=-1",
			result: `{
   "results": [
     {
       "name": "e01",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "info": {
         "Url": "https://foo.example.url"
       },
       "results": [
         {
           "timestamp": "2015-06",
           "availability": "76.26534166743393",
           "reliability": "91.61418757296076",
           "unknown": "0.00521",
           "uptime": "0.75868",
           "downtime": "0.166665"
         }
       ]
     },
     {
       "name": "e02",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06",
           "availability": "98.43749901562502",
           "reliability": "98.43749901562502",
           "unknown": "0",
           "uptime": "0.984375",
           "downtime": "0"
         }
       ]
     },
     {
       "name": "e03",
       "service": "service_b",
       "supergroup": "ST02",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06",
           "availability": "96.87499903125001",
           "reliability": "96.87499903125001",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     },
     {
       "name": "e03",
       "service": "service_b",
       "supergroup": "STX",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06",
           "availability": "96.87499903125001",
           "reliability": "96.87499903125001",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     },
     {
       "name": "e03",
       "service": "service_x",
       "supergroup": "ST02",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06",
           "availability": "96.87499903125001",
           "reliability": "96.87499903125001",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     }
   ]
 }`,
			code: 200,
		},
		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=monthly&pageSize=2",
			result: `{
   "results": [
     {
       "name": "e01",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "info": {
         "Url": "https://foo.example.url"
       },
       "results": [
         {
           "timestamp": "2015-06",
           "availability": "76.26534166743393",
           "reliability": "91.61418757296076",
           "unknown": "0.00521",
           "uptime": "0.75868",
           "downtime": "0.166665"
         }
       ]
     },
     {
       "name": "e02",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06",
           "availability": "98.43749901562502",
           "reliability": "98.43749901562502",
           "unknown": "0",
           "uptime": "0.984375",
           "downtime": "0"
         }
       ]
     }
   ],
   "nextPageToken": "Mg==",
   "pageSize": 2
 }`,
			code: 200,
		},

		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=monthly&pageSize=2&nextPageToken=Mg==",
			result: `{
   "results": [
     {
       "name": "e03",
       "service": "service_b",
       "supergroup": "ST02",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06",
           "availability": "96.87499903125001",
           "reliability": "96.87499903125001",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     },
     {
       "name": "e03",
       "service": "service_b",
       "supergroup": "STX",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06",
           "availability": "96.87499903125001",
           "reliability": "96.87499903125001",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     }
   ],
   "nextPageToken": "NA==",
   "pageSize": 2
 }`,
			code: 200,
		},

		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=monthly&pageSize=2&nextPageToken=NA==",
			result: `{
   "results": [
     {
       "name": "e03",
       "service": "service_x",
       "supergroup": "ST02",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06",
           "availability": "96.87499903125001",
           "reliability": "96.87499903125001",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     }
   ],
   "pageSize": 2
 }`,
			code: 200,
		},
	}

	for _, expReq := range expReqs {
		request, _ := http.NewRequest(expReq.method, expReq.url, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")

		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		suite.Equal(expReq.code, response.Code, "Incorrect HTTP response code")
		// Compare the expected and actual xml response
		suite.Equal(expReq.result, response.Body.String(), "Response body mismatch")

	}

}

func (suite *endpointAvailabilityTestSuite) TestCustomFlatAllEndpointsPaginated() {

	type expReq struct {
		method string
		url    string
		code   int
		result string
	}

	expReqs := []expReq{
		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=custom&pageSize=-1",
			result: `{
   "results": [
     {
       "name": "e01",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "info": {
         "Url": "https://foo.example.url"
       },
       "results": [
         {
           "availability": "76.26534166743393",
           "reliability": "91.61418757296076",
           "unknown": "0.00521",
           "uptime": "0.75868",
           "downtime": "0.166665"
         }
       ]
     },
     {
       "name": "e02",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "results": [
         {
           "availability": "98.43749901562502",
           "reliability": "98.43749901562502",
           "unknown": "0",
           "uptime": "0.984375",
           "downtime": "0"
         }
       ]
     },
     {
       "name": "e03",
       "service": "service_b",
       "supergroup": "ST02",
       "type": "endpoint",
       "results": [
         {
           "availability": "96.87499903125001",
           "reliability": "96.87499903125001",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     },
     {
       "name": "e03",
       "service": "service_b",
       "supergroup": "STX",
       "type": "endpoint",
       "results": [
         {
           "availability": "96.87499903125001",
           "reliability": "96.87499903125001",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     },
     {
       "name": "e03",
       "service": "service_x",
       "supergroup": "ST02",
       "type": "endpoint",
       "results": [
         {
           "availability": "96.87499903125001",
           "reliability": "96.87499903125001",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     }
   ]
 }`,
			code: 200,
		},
	}

	for _, expReq := range expReqs {
		request, _ := http.NewRequest(expReq.method, expReq.url, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")

		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		suite.Equal(expReq.code, response.Code, "Incorrect HTTP response code")
		// Compare the expected and actual xml response
		suite.Equal(expReq.result, response.Body.String(), "Response body mismatch")

	}

}

func (suite *endpointAvailabilityTestSuite) TestFlatAllEndpointsPaginated4() {

	type expReq struct {
		method string
		url    string
		code   int
		result string
	}

	expReqs := []expReq{
		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily&pageSize=4",
			result: `{
   "results": [
     {
       "name": "e01",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "info": {
         "Url": "https://foo.example.url"
       },
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "98.26389",
           "reliability": "98.26389",
           "unknown": "0",
           "uptime": "0.98264",
           "downtime": "0"
         },
         {
           "timestamp": "2015-06-23",
           "availability": "54.03509",
           "reliability": "81.48148",
           "unknown": "0.01042",
           "uptime": "0.53472",
           "downtime": "0.33333"
         }
       ]
     },
     {
       "name": "e02",
       "service": "service_a",
       "supergroup": "ST01",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
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
   ],
   "nextPageToken": "NA==",
   "pageSize": 4
 }`,
			code: 200,
		},
		expReq{
			method: "GET",
			url:    "/api/v2/results/Report_A/endpoints?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily&pageSize=4&nextPageToken=NA==",
			result: `{
   "results": [
     {
       "name": "e03",
       "service": "service_b",
       "supergroup": "ST02",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     },
     {
       "name": "e03",
       "service": "service_b",
       "supergroup": "STX",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     },
     {
       "name": "e03",
       "service": "service_x",
       "supergroup": "ST02",
       "type": "endpoint",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": "96.875",
           "reliability": "96.875",
           "unknown": "0",
           "uptime": "0.96875",
           "downtime": "0"
         }
       ]
     }
   ],
   "pageSize": 4
 }`,
			code: 200,
		},
	}

	for _, expReq := range expReqs {
		request, _ := http.NewRequest(expReq.method, expReq.url, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")

		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		suite.Equal(expReq.code, response.Code, "Incorrect HTTP response code")
		// Compare the expected and actual xml response
		suite.Equal(expReq.result, response.Body.String(), "Response body mismatch")

	}

}

// TestListServiceFlavorAvailabilityDaily tests if daily results are returned correctly
func (suite *endpointAvailabilityTestSuite) TestListEndpointAvailabilityDaily() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/services/service_a/endpoints?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	endpointAvailabilityXML := ` <root>
   <group name="ST01" type="SITE">
     <group name="service_a" type="service">
       <group name="e01" type="endpoint">
         <results timestamp="2015-06-22" availability="98.26389" reliability="98.26389" unknown="0" uptime="0.98264" downtime="0"></results>
         <results timestamp="2015-06-23" availability="54.03509" reliability="81.48148" unknown="0.01042" uptime="0.53472" downtime="0.33333"></results>
       </group>
       <group name="e02" type="endpoint">
         <results timestamp="2015-06-22" availability="96.875" reliability="96.875" unknown="0" uptime="0.96875" downtime="0"></results>
         <results timestamp="2015-06-23" availability="100" reliability="100" unknown="0" uptime="1" downtime="0"></results>
       </group>
     </group>
   </group>
 </root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointAvailabilityXML, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/services/service_a/endpoints?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	serviceFlavorAvailabilityJSON := `{
   "results": [
     {
       "name": "ST01",
       "type": "SITE",
       "serviceflavors": [
         {
           "name": "service_a",
           "type": "service",
           "endpoints": [
             {
               "name": "e01",
               "type": "endpoint",
               "info": {
                 "Url": "https://foo.example.url"
               },
               "results": [
                 {
                   "timestamp": "2015-06-22",
                   "availability": "98.26389",
                   "reliability": "98.26389",
                   "unknown": "0",
                   "uptime": "0.98264",
                   "downtime": "0"
                 },
                 {
                   "timestamp": "2015-06-23",
                   "availability": "54.03509",
                   "reliability": "81.48148",
                   "unknown": "0.01042",
                   "uptime": "0.53472",
                   "downtime": "0.33333"
                 }
               ]
             },
             {
               "name": "e02",
               "type": "endpoint",
               "results": [
                 {
                   "timestamp": "2015-06-22",
                   "availability": "96.875",
                   "reliability": "96.875",
                   "unknown": "0",
                   "uptime": "0.96875",
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
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(serviceFlavorAvailabilityJSON, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/services?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", "AWRONGKEY")
	request.Header.Set("Accept", "application/xml")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 401 Unauthorized code
	suite.Equal(401, response.Code, "Incorrect HTTP response code")

}

// TestListServiceFlavorAvailabilityErrors tests if errors/exceptions are returned correctly
func (suite *endpointAvailabilityTestSuite) TestListEndpointAvailabilityErrors() {

	reportErrorXML := ` <root>
   <message>The report with the name Report_B does not exist</message>
   <code>404</code>
 </root>`

	typeErrorXML := `<root>
 <status>
  <message>Bad Request</message>
  <code>400</code>
 </status>
 <errors>
  <error>
   <message>Endpoint Group type not in report</message>
   <code>400</code>
   <details>Endpoint Group type Site not present in report Report_A. Try using SITE instead</details>
  </error>
 </errors>
</root>`

	typeError1XML := ` <root>
   <message>No results found for given query</message>
   <code>404</code>
 </root>`

	typeError1JSON := `{
   "message": "No results found for given query",
   "code": 404
 }`

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_B/SITE/ST01/services/service_a/endpoints?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 404 not found code
	suite.Equal(404, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(reportErrorXML, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/Site/ST01/services/service_a/endpoints?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("Accept", "application/xml")
	request.Header.Set("x-api-key", suite.clientkey)

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	output := response.Body.String()

	// Check that we must have a 400 bad request code
	suite.Equal(400, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(typeErrorXML, output, "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/services?start_time=2025-06-22T00:00:00Z&end_time=2025-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("Accept", "application/xml")
	request.Header.Set("x-api-key", suite.clientkey)

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	output = response.Body.String()

	// Check that we must have a 404 not found code
	suite.Equal(404, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(typeError1XML, output, "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/services?start_time=2025-06-22T00:00:00Z&end_time=2025-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("Accept", "application/json")
	request.Header.Set("x-api-key", suite.clientkey)

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	output = response.Body.String()

	// Check that we must have a 404 not found code
	suite.Equal(404, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(typeError1JSON, output, "Response body mismatch")

}

// TestOptionsServiceFlavor tests responses in case the OPTIONS http verb is used
func (suite *endpointAvailabilityTestSuite) TestOptionsServiceFlavor() {

	request, _ := http.NewRequest("OPTIONS", "/api/v2/results/Report_A/SITE/ST01/services/service_a/endpoints", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v2/results/Report_A/SITE/ST01/services/service_a/endpoints/e01", strings.NewReader(""))

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	headers = response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v2/results/Report_A/GROUP/GROUP_A/SITE/ST01/services/service_a/endpoints", strings.NewReader(""))

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	headers = response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v2/results/Report_A/GROUP/GROUP_A/SITE/ST01/services/service_a/endpoints/e01", strings.NewReader(""))

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

// TestStrictSlashServiceFlavorResults test if not found responses are returned correctly
func (suite *endpointAvailabilityTestSuite) TestStrictSlashServiceFlavorResults() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/services/SF01/endpoints/?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")
	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	suite.Equal(404, response.Code, "Incorrect HTTP response code")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/services/SF01/endpoints/e01/?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")
	response = httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	suite.Equal(404, response.Code, "Incorrect HTTP response code")

}

//TearDownTest to tear down every test
func (suite *endpointAvailabilityTestSuite) TearDownTest() {

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
func (suite *endpointAvailabilityTestSuite) TearDownSuite() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

// TestEndpointAvailabilityTestSuite is responsible for calling the tests
func TestEndpointAvailabilityTestSuite(t *testing.T) {
	suite.Run(t, new(endpointAvailabilityTestSuite))
}
