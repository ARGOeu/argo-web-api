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

type serviceFlavorAvailabilityTestSuite struct {
	suite.Suite
	cfg          config.Config
	router       *mux.Router
	confHandler  respond.ConfHandler
	tenantDbConf config.MongoConfig
	clientkey    string
}

// Setup the Test Environment
func (suite *serviceFlavorAvailabilityTestSuite) SetupSuite() {

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
	db = "ARGO_test_service_types_availability"
	`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.tenantDbConf.Db = "ARGO_test_service_types_availability_tenant"
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
func (suite *serviceFlavorAvailabilityTestSuite) SetupTest() {

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
					"api_key": "bob_key",
					"roles":   []string{"viewer"},
				},
				{
					"name":    "alice",
					"email":   "alice@tenant-aa.foo",
					"api_key": "alice_key",
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
					"database": "argo_wrong_db_serviceflavoravailability",
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
					"email":   "alice@tenant-bb.foo",
					"api_key": "alice_key",
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
	// Seed database with recomputations
	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("service_ar")

	// Insert seed data
	c.InsertOne(context.TODO(),
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "SF01",
			"supergroup":   "ST01",
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
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "SF02",
			"supergroup":   "ST01",
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
			"name":         "SF03",
			"supergroup":   "ST02",
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
			"name":         "SF01",
			"supergroup":   "ST01",
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
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150623,
			"name":         "SF02",
			"supergroup":   "ST01",
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

// TestListServiceFlavorAvailabilityMonthly tests if daily results are returned correctly
func (suite *serviceFlavorAvailabilityTestSuite) TestListServiceTypeAvailabilityMonthly() {

	request, _ := http.NewRequest("GET", "/api/v5/results/Report_A/groups/ST01/service-types?start-time=2015-06-22T00:00:00Z&end-time=2015-06-23T23:59:59Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	responseBody := response.Body.String()

	serviceFlavorAvailabilityJSON := `{
   "results": [
     {
       "name": "ST01",
       "type": "SITE",
       "service-types": [
         {
           "name": "SF01",
           "type": "service",
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
         },
         {
           "name": "SF02",
           "type": "service",
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
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(serviceFlavorAvailabilityJSON, responseBody, "Response body mismatch")

}

func (suite *serviceFlavorAvailabilityTestSuite) TestListServiceTypeAvailabilityCustom() {

	request, _ := http.NewRequest("GET", "/api/v5/results/Report_A/groups/ST01/service-types?start-time=2015-06-22T00:00:00Z&end-time=2015-06-23T23:59:59Z&granularity=custom", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	responseBody := response.Body.String()
	serviceFlavorAvailabilityJSON := `{
   "results": [
     {
       "name": "ST01",
       "type": "SITE",
       "service-types": [
         {
           "name": "SF01",
           "type": "service",
           "results": [
             {
               "availability": 76.27,
               "reliability": 91.61,
               "unknown": 0.01,
               "uptime": 0.76,
               "downtime": 0.17
             }
           ]
         },
         {
           "name": "SF02",
           "type": "service",
           "results": [
             {
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
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(serviceFlavorAvailabilityJSON, responseBody, "Response body mismatch")

}

// TestListServiceFlavorAvailabilityDaily tests if daily results are returned correctly
func (suite *serviceFlavorAvailabilityTestSuite) TestListServiceTypeAvailabilityDaily() {

	request, _ := http.NewRequest("GET", "/api/v5/results/Report_A/groups/ST01/service-types?start-time=2015-06-22T00:00:00Z&end-time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	serviceFlavorAvailabilityJSON := `{
   "results": [
     {
       "name": "ST01",
       "type": "SITE",
       "service-types": [
         {
           "name": "SF01",
           "type": "service",
           "results": [
             {
               "timestamp": "2015-06-22",
               "availability": 98.26,
               "reliability": 98.26,
               "unknown": 0,
               "uptime": 0.98,
               "downtime": 0
             },
             {
               "timestamp": "2015-06-23",
               "availability": 54.04,
               "reliability": 81.48,
               "unknown": 0.01,
               "uptime": 0.53,
               "downtime": 0.33
             }
           ]
         },
         {
           "name": "SF02",
           "type": "service",
           "results": [
             {
               "timestamp": "2015-06-22",
               "availability": 96.88,
               "reliability": 96.88,
               "unknown": 0,
               "uptime": 0.97,
               "downtime": 0
             },
             {
               "timestamp": "2015-06-23",
               "availability": 100,
               "reliability": 100,
               "unknown": 0,
               "uptime": 1,
               "downtime": 0
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

	request, _ = http.NewRequest("GET", "/api/v5/results/Report_A/groups/ST01/service-types?start-time=2015-06-22T00:00:00Z&end-time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", "AWRONGKEY")
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 401 Unauthorized code
	suite.Equal(401, response.Code, "Incorrect HTTP response code")

}

// TestOptionsServiceFlavor tests responses in case the OPTIONS http verb is used
func (suite *serviceFlavorAvailabilityTestSuite) TestOptionsServiceFlavor() {

	request, _ := http.NewRequest("OPTIONS", "/api/v5/results/Report_A/groups/ST01/service-types", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v5/results/Report_A/groups/ST01/service-types/service_a", strings.NewReader(""))

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

// TestStrictSlashServiceFlavorResults test if not found responses are returned correctly
func (suite *serviceFlavorAvailabilityTestSuite) TestStrictSlashServiceFlavorResults() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/groups/ST01/service-types/?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")
	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	suite.Equal(404, response.Code, "Incorrect HTTP response code")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/groups/ST01/service-types/SF01/?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")
	response = httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	suite.Equal(404, response.Code, "Incorrect HTTP response code")

}

// TearDownTest to tear down every test
func (suite *serviceFlavorAvailabilityTestSuite) TearDownTest() {

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
func (suite *serviceFlavorAvailabilityTestSuite) TearDownSuite() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
}

// TestSuiteResultService responsible for calling the tests
func TestSuiteResultsService(t *testing.T) {
	suite.Run(t, new(serviceFlavorAvailabilityTestSuite))
}
