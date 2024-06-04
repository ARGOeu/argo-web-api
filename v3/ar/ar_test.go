/*
 * Copyright (c) 2022 GRNET S.A.
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

package ar

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/store"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/gcfg.v1"
)

type AvailabilityTestSuite struct {
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
func (suite *AvailabilityTestSuite) SetupSuite() {

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
		 db = "ARGO_test_arV3"
		 `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.tenantDbConf.Db = "ARGO_test_arV3_tenant"
	suite.tenantDbConf.Password = "h4shp4ss"
	suite.tenantDbConf.Username = "dbuser"
	suite.tenantDbConf.Store = "ar"
	suite.clientkey = "secretkey"

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v3/results").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)
}

// This function runs before any test and setups the environment
func (suite *AvailabilityTestSuite) SetupTest() {

	log.SetOutput(io.Discard)

	// Seed database with tenants
	//TODO: move tests to
	c := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("tenants")
	c.InsertOne(context.TODO(),
		bson.M{"name": "Westeros",
			"db_conf": []bson.M{
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Westeros1",
				},
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Westeros2",
				},
			},
			"users": []bson.M{
				{
					"name":    "John Snow",
					"email":   "J.Snow@brothers.wall",
					"api_key": "wh1t3_w@lk3rs",
					"roles":   []string{"viewer"},
				},
				{
					"name":    "King Joffrey",
					"email":   "g0dk1ng@kingslanding.gov",
					"api_key": "sansa <3",
					"roles":   []string{"viewer"},
				},
			}})
	c.InsertOne(context.TODO(),
		bson.M{"name": "EGI",
			"db_conf": []bson.M{
				{
					"server":   "localhost",
					"port":     27017,
					"database": suite.tenantDbConf.Db,
					"username": suite.tenantDbConf.Username,
					"password": suite.tenantDbConf.Password,
				},
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_wrong_db_endpointgrouavailability",
				},
			},
			"users": []bson.M{
				{
					"name":    "Joe Complex",
					"email":   "C.Joe@egi.eu",
					"api_key": suite.clientkey,
					"roles":   []string{"viewer"},
				},
				{
					"name":    "Josh Plain",
					"email":   "P.Josh@egi.eu",
					"api_key": "itsamysterytoyou",
					"roles":   []string{"viewer"},
				},
			}})

	c = suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("roles")
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "v3.ar.list-by-id",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "v3.ar.list",
			"roles":    []string{"editor", "viewer"},
		})

	// Seed tenant database with data
	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("endpoint_group_ar")

	// Insert seed data
	c.InsertOne(context.TODO(),
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
				{
					"name":  "",
					"value": "",
				},
			},
		})
	c.InsertOne(context.TODO(),
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
		})
	c.InsertOne(context.TODO(),
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
				{
					"name":  "",
					"value": "",
				},
			},
		})
	c.InsertOne(context.TODO(),
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
				{
					"name":  "",
					"value": "",
				},
			},
		})

	// Seed endpoint data
	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("endpoint_ar")

	// Insert seed data
	c.InsertOne(context.TODO(),
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "host01",
			"service":      "service01",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 66.7,
			"reliability":  54.6,
			"weight":       5634,
			"tags": []bson.M{
				{
					"name":  "",
					"value": "",
				},
			},
			"info": bson.M{"ID": "special-queue"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150623,
			"name":         "host01",
			"service":      "service01",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 100,
			"reliability":  100,
			"weight":       5634,
			"tags": []bson.M{
				{
					"name":  "",
					"value": "",
				},
			},
			"info": bson.M{"ID": "special-queue"},
		})

	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("reports")

	c.InsertOne(context.TODO(), bson.M{
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
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"name":  "name1",
				"value": "value1"},
			{
				"name":  "name2",
				"value": "value2"},
		}})

}

// TestListEndpointGroupAvailability test if daily results are returned correctly
func (suite *AvailabilityTestSuite) TestListEndpointGroupAvailability() {

	request, _ := http.NewRequest("GET", "/api/v3/results/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	endpointGroupAvailabilityA := `{
   "results": [
     {
       "name": "GROUP_A",
       "type": "GROUP",
       "results": [
         {
           "date": "2015-06-22",
           "availability": "68.13896116893515",
           "reliability": "50.413931144915935"
         },
         {
           "date": "2015-06-23",
           "availability": "75.36324059247399",
           "reliability": "80.8138510808647"
         }
       ],
       "groups": [
         {
           "name": "ST01",
           "type": "SITES",
           "results": [
             {
               "date": "2015-06-22",
               "availability": "66.7",
               "reliability": "54.6",
               "unknown": "0",
               "uptime": "1",
               "downtime": "0"
             },
             {
               "date": "2015-06-23",
               "availability": "100",
               "reliability": "100",
               "unknown": "0",
               "uptime": "1",
               "downtime": "0"
             }
           ]
         },
         {
           "name": "ST02",
           "type": "SITES",
           "results": [
             {
               "date": "2015-06-22",
               "availability": "70",
               "reliability": "45",
               "unknown": "0",
               "uptime": "1",
               "downtime": "0"
             },
             {
               "date": "2015-06-23",
               "availability": "43.5",
               "reliability": "56",
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
	suite.Equal(endpointGroupAvailabilityA, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v3/results/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	endpointGroupAvailabilityJSON := `{
   "results": [
     {
       "name": "GROUP_A",
       "type": "GROUP",
       "results": [
         {
           "date": "2015-06",
           "availability": "71.75110088070457",
           "reliability": "65.61389111289031"
         }
       ],
       "groups": [
         {
           "name": "ST01",
           "type": "SITES",
           "results": [
             {
               "date": "2015-06",
               "availability": "99.99999900000002",
               "reliability": "99.99999900000002",
               "unknown": "0",
               "uptime": "1",
               "downtime": "0"
             }
           ]
         },
         {
           "name": "ST02",
           "type": "SITES",
           "results": [
             {
               "date": "2015-06",
               "availability": "99.99999900000002",
               "reliability": "99.99999900000002",
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

	request, _ = http.NewRequest("GET", "/api/v3/results/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", "AWRONGKEY")
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	unauthorizedresponse := `{
 "status": {
  "message": "Unauthorized",
  "code": "401",
  "details": "You need to provide a correct authentication token using the header 'x-api-key'"
 }
}`

	// Check that we must have a 401 Unauthorized code
	suite.Equal(401, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(unauthorizedresponse, response.Body.String(), "Response body mismatch")

}

// TestListEndpointAvailabilityCustom test if daily results are returned correctly for a specific id
func (suite *AvailabilityTestSuite) TestListGroupAvailCustom() {

	request, _ := http.NewRequest("GET", "/api/v3/results/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	expected := `{
   "results": [
     {
       "name": "GROUP_A",
       "type": "GROUP",
       "results": [
         {
           "date": "2015-06",
           "availability": "71.75110088070457",
           "reliability": "65.61389111289031"
         }
       ],
       "groups": [
         {
           "name": "ST01",
           "type": "SITES",
           "results": [
             {
               "date": "2015-06",
               "availability": "99.99999900000002",
               "reliability": "99.99999900000002",
               "unknown": "0",
               "uptime": "1",
               "downtime": "0"
             }
           ]
         },
         {
           "name": "ST02",
           "type": "SITES",
           "results": [
             {
               "date": "2015-06",
               "availability": "99.99999900000002",
               "reliability": "99.99999900000002",
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
	suite.Equal(expected, response.Body.String(), "Response body mismatch")

}

// TestListEndpointAvailability test if daily results are returned correctly for a specific id
func (suite *AvailabilityTestSuite) TestListEndpointAvailabilityDaily() {

	request, _ := http.NewRequest("GET", "/api/v3/results/Report_A/id/special-queue?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	expected := `{
   "id": "special-queue",
   "endpoints": [
     {
       "name": "host01",
       "service": "service01",
       "supergroup": "GROUP_A",
       "info": {
         "ID": "special-queue"
       },
       "results": [
         {
           "date": "2015-06-22",
           "availability": "66.7",
           "reliability": "54.6",
           "unknown": "0",
           "uptime": "1",
           "downtime": "0"
         },
         {
           "date": "2015-06-23",
           "availability": "100",
           "reliability": "100",
           "unknown": "0",
           "uptime": "1",
           "downtime": "0"
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(expected, response.Body.String(), "Response body mismatch")

}

// TestListEndpointAvailabilityMonthly test if monthly results are returned correctly for a specific id
func (suite *AvailabilityTestSuite) TestListEndpointAvailabilityMonthly() {

	request, _ := http.NewRequest("GET", "/api/v3/results/Report_A/id/special-queue?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	expected := `{
   "id": "special-queue",
   "endpoints": [
     {
       "name": "host01",
       "service": "service01",
       "supergroup": "GROUP_A",
       "info": {
         "ID": "special-queue"
       },
       "results": [
         {
           "date": "2015-06",
           "availability": "99.99999900000002",
           "reliability": "99.99999900000002",
           "unknown": "0",
           "uptime": "1",
           "downtime": "0"
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(expected, response.Body.String(), "Response body mismatch")

}

func (suite *AvailabilityTestSuite) TestListEndpointAvailabilityCustom() {

	request, _ := http.NewRequest("GET", "/api/v3/results/Report_A/id/special-queue?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=custom", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	expected := `{
   "id": "special-queue",
   "endpoints": [
     {
       "name": "host01",
       "service": "service01",
       "supergroup": "GROUP_A",
       "info": {
         "ID": "special-queue"
       },
       "results": [
         {
           "availability": "99.99999900000002",
           "reliability": "99.99999900000002",
           "unknown": "0",
           "uptime": "1",
           "downtime": "0"
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(expected, response.Body.String(), "Response body mismatch")

}

// TestListEndpointIdNotFound tests the response if there are no endpoints with the resource-id requested
func (suite *AvailabilityTestSuite) TestListEndpointIdNotFound() {

	request, _ := http.NewRequest("GET", "/api/v3/results/Report_A/id/another-queue?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	expected := `{
   "message": "No endpoints found with resource-id: another-queue",
   "code": 404
 }`

	// Check that we must have a 200 ok code
	suite.Equal(404, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(expected, response.Body.String(), "Response body mismatch")

}

// TestOmitingStartEnd tests if the validator catches omissions of start or end date
func (suite *AvailabilityTestSuite) TestOmittingStartEnd() {

	urls := []string{
		"/api/v3/results/Report_A/id/another-queue?end_time=2015-06-23T23:00:00Z",
		"/api/v3/results/Report_A/id/another-queue?start_time=2015-06-23T23:00:00Z",
		"/api/v3/results/Report_A",
		"/api/v3/results/Report_A?start_time=2015-06-23T23:00:00Z",
		"/api/v3/results/Report_A?end_time=2015-06-23T23:00:00Z",
	}

	expected := `{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "No time span set",
   "code": "400",
   "details": "Please use start_time and end_time url parameters to set the prefered time span"
  }
 ]
}`

	for _, url := range urls {
		request, _ := http.NewRequest("GET", url, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")

		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		// Check that we must have a 200 ok code
		suite.Equal(400, response.Code, "Incorrect HTTP response code")
		// Compare the expected and actual xml response
		suite.Equal(expected, response.Body.String(), "Response body mismatch")

	}

}

// TestOptions tests responses in case the OPTIONS http verb is used
func (suite *AvailabilityTestSuite) TestOptions() {

	request, _ := http.NewRequest("OPTIONS", "/api/v3/results/Report_A", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

}

// TearDownTest to tear down every test
func (suite *AvailabilityTestSuite) TearDownTest() {

	mainDB := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db)
	cols, err := mainDB.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		panic(err)
	}

	for _, col := range cols {
		mainDB.Collection(col).Drop(context.TODO())
	}

	tenantDB := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db)
	cols, err = tenantDB.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		panic(err)
	}

	for _, col := range cols {
		tenantDB.Collection(col).Drop(context.TODO())
	}

}

// TearDownTest to tear down every test
func (suite *AvailabilityTestSuite) TearDownSuite() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
}

// TestEndpointGroupsTestSuite is responsible for calling the tests
func TestSuiteAR(t *testing.T) {
	suite.Run(t, new(AvailabilityTestSuite))
}
