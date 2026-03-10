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

package nodes

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

type StatusTestSuite struct {
	suite.Suite
	cfg          config.Config
	router       *mux.Router
	confHandler  respond.ConfHandler
	tenantDbConf config.MongoConfig

	clientkey string
}

// Setup the Test Environment
func (suite *StatusTestSuite) SetupSuite() {

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
		   db = "ARGO_test_node_status"
		   `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.tenantDbConf.Db = "ARGO_test_statusV4_tenant"
	suite.tenantDbConf.Password = "h4shp4ss"
	suite.tenantDbConf.Username = "dbuser"
	suite.tenantDbConf.Store = "ar"
	suite.clientkey = "secretkey"

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v4").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)
}

// This function runs before any test and setups the environment
func (suite *StatusTestSuite) SetupTest() {

	log.SetOutput(io.Discard)

	// Seed database with tenants
	//TODO: move tests to
	c := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("tenants")
	c.InsertOne(context.TODO(),
		bson.M{"name": "TENANTA",
			"db_conf": []bson.M{
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Tenant1",
				},
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Tenant2",
				},
			},
			"node": bson.M{
				"id":   "node-A",
				"name": "NODEA",
			},
			"users": []bson.M{
				{
					"name":    "bob",
					"email":   "bob@foo",
					"api_key": "bob-token",
					"roles":   []string{"viewer"},
				},
				{
					"name":    "alice",
					"email":   "alice@foo",
					"api_key": "alice-token",
					"roles":   []string{"viewer"},
				},
			}})
	c.InsertOne(context.TODO(),
		bson.M{"name": "TENANTB",
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
			"node": bson.M{
				"id":   "node-B",
				"name": "NODEB",
			},
			"users": []bson.M{
				{
					"name":    "john",
					"email":   "john@foo",
					"api_key": suite.clientkey,
					"roles":   []string{"viewer"},
				},
				{
					"name":    "george",
					"email":   "george@foo",
					"api_key": "itsamysterytoyou",
					"roles":   []string{"viewer"},
				},
			}})

	c = suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("roles")
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "v4.nodes.availability",
			"roles":    []string{"super_admin", "admin", "editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "v4.nodes.uptime",
			"roles":    []string{"super_admin", "admin", "editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "v4.nodes.status",
			"roles":    []string{"super_admin", "admin", "editor", "viewer"},
		})

	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("reports")

	c.InsertOne(context.TODO(), bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a494364",
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

	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection(nodeReportColName)
	c.InsertOne(context.TODO(), bson.M{
		"_id":       "node-report",
		"report_id": "eba61a9e-22e9-4521-9e47-ecaa4a494364",
	})
	// Seed tenant database with data
	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection(statusGroupColName)

	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "SITEA",
		"status":         "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T01:00:00Z",
		"endpoint_group": "SITEA",
		"status":         "CRITICAL",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T05:00:00Z",
		"endpoint_group": "SITEA",
		"status":         "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "SITEB",
		"status":         "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T03:00:00Z",
		"endpoint_group": "SITEB",
		"status":         "WARNING",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T17:53:00Z",
		"endpoint_group":     "SITEB",
		"status":             "CRITICAL",
		"has_threshold_rule": true,
	})

}

// TestListStatus test the status results
func (suite *StatusTestSuite) TestListStatus() {

	request, _ := http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/status?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:59:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	expResponse := `{
 "data": [
  {
   "name": "SITEA",
   "results": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T01:00:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T05:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  },
  {
   "name": "SITEB",
   "results": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T03:00:00Z",
     "value": "WARNING"
    },
    {
     "timestamp": "2015-05-01T17:53:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "CRITICAL"
    }
   ]
  }
 ]
}`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual json response
	suite.Equal(expResponse, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/status?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z", strings.NewReader(""))
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

	// Case of bad start_time input
	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/status?start_time=2015-06-20AT12:00:00Z&end_time=2015-06-23T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	badRequestResponse := `{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "Bad Request",
   "code": "400",
   "details": "Error parsing start_time=2015-06-20AT12:00:00Z - please use zulu format like 2006-01-02T15:04:05Z"
  }
 ]
}`
	// Check that we must have a 401 Unauthorized code
	suite.Equal(400, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(badRequestResponse, response.Body.String(), "Response body mismatch")

	// Case of bad end_time input
	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/status?start_time=2015-06-20T12:00:00Z&end_time=2015-06T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	badRequestResponse = `{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "Bad Request",
   "code": "400",
   "details": "Error parsing end_time=2015-06T23:00:00Z - please use zulu format like 2006-01-02T15:04:05Z"
  }
 ]
}`
	// Check that we must have a 401 Unauthorized code
	suite.Equal(400, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(badRequestResponse, response.Body.String(), "Response body mismatch")

}

// TestOptions tests responses in case the OPTIONS http verb is used
func (suite *StatusTestSuite) TestOptions() {

	request, _ := http.NewRequest("OPTIONS", "/api/v4/nodes/NODEB/capabilities/status", strings.NewReader(""))

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
func (suite *StatusTestSuite) TearDownTest() {

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
func (suite *StatusTestSuite) TearDownSuite() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
}

// TestEndpointGroupsTestSuite is responsible for calling the tests
func TestSuiteStatus(t *testing.T) {
	suite.Run(t, new(StatusTestSuite))
}
