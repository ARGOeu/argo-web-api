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

package resultsV5

import (
	"context"
	"fmt"
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

type endpointAvailabilityTestSuite struct {
	suite.Suite
	cfg          config.Config
	router       *mux.Router
	confHandler  respond.ConfHandler
	tenantDbConf config.MongoConfig
	clientkey    string
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

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.tenantDbConf.Db = "ARGO_test_new_endpoint_availability_tenant"
	suite.tenantDbConf.Password = "h4shp4ss"
	suite.tenantDbConf.Username = "dbuser"
	suite.tenantDbConf.Store = "ar"
	suite.clientkey = "secretkey"

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v5/results").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)
}

// This function runs before any test and setups the environment
func (suite *endpointAvailabilityTestSuite) SetupTest() {

	log.SetOutput(io.Discard)

	// Seed database with tenants
	//TODO: move tests to
	c := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("tenants")
	c.InsertOne(context.TODO(),
		bson.M{"name": "TENANT-AA",
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
					"name":    "bob",
					"email":   "bob@tenant-aa.foo",
					"api_key": "bob_secret",
					"roles":   []string{"viewer"},
				},
				{
					"name":    "alice",
					"email":   "alice@tenant-aa.foo",
					"api_key": "alice_secret",
					"roles":   []string{"viewer"},
				},
			}})
	c.InsertOne(context.TODO(),
		bson.M{"name": "TENANT-BB",
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
					"database": "argo_wrong_db_endpoint_availability",
				},
			},
			"users": []bson.M{

				{
					"name":    "bob",
					"email":   "bob@tenant-bb.foo",
					"api_key": suite.clientkey,
					"roles":   []string{"viewer"},
				},
				{
					"name":    "alice",
					"email":   "alice@tenant-cc.foo",
					"api_key": "alice_secret2",
					"roles":   []string{"viewer"},
				},
			}})

	c = suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("roles")
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "results.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "results.get",
			"roles":    []string{"editor", "viewer"},
		})

	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("endpoint_ar")

	// Insert seed data
	c.InsertOne(context.TODO(),
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
				{
					"name":  "production",
					"value": "Y",
				},
			},
			"info": bson.M{
				"Url": "https://foo.example.url",
			},
		})
	c.InsertOne(context.TODO(),
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
				{
					"name":  "production",
					"value": "Y",
				},
			},
		})
	c.InsertOne(context.TODO(),
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
				{
					"name":  "production",
					"value": "Y",
				},
			},
		})
	c.InsertOne(context.TODO(),
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
				{
					"name":  "production",
					"value": "Y",
				},
			},
		})
	c.InsertOne(context.TODO(),
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
				{
					"name":  "production",
					"value": "Y",
				},
			},
		})
	c.InsertOne(context.TODO(),
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
				{
					"name":  "production",
					"value": "Y",
				},
			},
			"info": bson.M{
				"Url": "https://foo.example.url",
			},
		})
	c.InsertOne(context.TODO(),
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
				{
					"name":  "production",
					"value": "Y",
				},
			},
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
					"type": "SITE",
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

// TestListEndpointAvailabilityMonthly tests if monthly results are returned correctly
func (suite *endpointAvailabilityTestSuite) TestListEndpointAvailabilityMonthly() {

	request, _ := http.NewRequest("GET", "/api/v5/results/Report_A/endpoints/e01?start-time=2015-06-22T00:00:00Z&end-time=2015-06-23T23:59:59Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	responseBody := response.Body.String()
	endpointAvailabilityJSON1 := `{
   "results": [
     {
       "name": "ST01",
       "type": "SITE",
       "service-types": [
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
                   "availability": 76.27,
                   "reliability": 91.61,
                   "unknown": 0.01,
                   "uptime": 0.76,
                   "downtime": 0.17
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
	// Compare the expected and actual json response
	suite.Equal(endpointAvailabilityJSON1, responseBody, "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v5/results/Report_A/endpoints/e02?start-time=2015-06-22T00:00:00Z&end-time=2015-06-23T23:59:59Z&granularity=monthly", strings.NewReader(""))
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
       "service-types": [
         {
           "name": "service_a",
           "type": "service",
           "endpoints": [
             {
               "name": "e02",
               "type": "endpoint",
               "results": [
                 {
                   "timestamp": "2015-06",
                   "availability": 98.44,
                   "reliability": 98.44,
                   "unknown": 0,
                   "uptime": 0.98,
                   "downtime": 0
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

	expectedMonthly := `{
   "results": [
     {
       "name": "ST01",
       "type": "SITE",
       "service-types": [
         {
           "name": "service_a",
           "type": "service",
           "endpoints": [
             {
               "name": "e02",
               "type": "endpoint",
               "results": [
                 {
                   "timestamp": "2015-06",
                   "availability": 98.44,
                   "reliability": 98.44,
                   "unknown": 0,
                   "uptime": 0.98,
                   "downtime": 0
                 }
               ]
             }
           ]
         }
       ]
     }
   ]
 }`

	// Test monthly a/r
	request, _ = http.NewRequest("GET", "/api/v5/results/Report_A/endpoints/e02?start-time=2015-06-22T00:00:00Z&end-time=2015-06-23T23:59:59Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	responseBody = response.Body.String()

	fmt.Println(responseBody)

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(expectedMonthly, responseBody, "Response body mismatch")

	// Test endpoint ar through specific group

	expectedSpecGroup := `{
   "results": [
     {
       "name": "STX",
       "type": "SITE",
       "service-types": [
         {
           "name": "service_b",
           "type": "service",
           "endpoints": [
             {
               "name": "e03",
               "type": "endpoint",
               "results": [
                 {
                   "timestamp": "2015-06",
                   "availability": 96.87,
                   "reliability": 96.87,
                   "unknown": 0,
                   "uptime": 0.97,
                   "downtime": 0
                 }
               ]
             }
           ]
         }
       ]
     }
   ]
 }`
	request, _ = http.NewRequest("GET", "/api/v5/results/Report_A/groups/STX/endpoints?start-time=2015-06-22T00:00:00Z&end-time=2015-06-23T23:59:59Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	responseBody = response.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(expectedSpecGroup, responseBody, "Response body mismatch")

}

func (suite *SuperGroupAvailabilityTestSuite) TestOptionsEndpoints() {
	request, _ := http.NewRequest("OPTIONS", "/api/v5/results/Report_A/endpoints", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v5/results/Report_A/endpoints/e01", strings.NewReader(""))

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	headers = response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

}

// TearDownTest to tear down every test
func (suite *endpointAvailabilityTestSuite) TearDownTest() {

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
func (suite *endpointAvailabilityTestSuite) TearDownSuite() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
}

// TestSuiteResultsEndpoint is responsible for calling the tests
func TestSuiteResultsEndpoint(t *testing.T) {
	suite.Run(t, new(endpointAvailabilityTestSuite))
}
