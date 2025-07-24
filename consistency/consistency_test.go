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

package consistency

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

type consistencyTestSuite struct {
	suite.Suite
	cfg          config.Config
	router       *mux.Router
	confHandler  respond.ConfHandler
	tenantDbConf config.MongoConfig
	clientkey    string
	tenant1key   string
	tenant2key   string
	tenant1db    string
	tenant2db    string
}

// Setup the Test Environment
func (suite *consistencyTestSuite) SetupSuite() {

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
	db = "ARGO_test_topology_test"
	`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.tenantDbConf.Db = "ARGO_consistency_check"
	suite.tenantDbConf.Password = "h4shp4ss"
	suite.tenantDbConf.Username = "dbuser"
	suite.tenantDbConf.Store = "ar"
	suite.clientkey = "viewerkey"
	suite.tenant1key = "tenant1key"
	suite.tenant2key = "tenant2key"
	suite.tenant1db = "argo_tenant1"
	suite.tenant2db = "argo_tenant2"

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v2").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)
}

// This function runs before any test and setups the environment
func (suite *consistencyTestSuite) SetupTest() {

	log.SetOutput(io.Discard)

	// Seed database with tenants
	//TODO: move tests to
	c := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("tenants")
	c.InsertOne(context.TODO(),
		bson.M{"name": "TestTenant",
			"db_conf": []bson.M{

				{
					"server":   "localhost",
					"port":     27017,
					"database": suite.tenant1db,
				},
			},
			"users": []bson.M{

				{
					"name":    "Viewer",
					"email":   "viewer@example.com",
					"api_key": "viewerkey",
					"roles":   []string{"viewer"},
				},
				{
					"name":    "auto check service",
					"email":   "check@example.com",
					"api_key": "checkkey",
					"roles":   []string{"consistency-check"},
				},
				{
					"name":    "ack admin",
					"email":   "ack@example.com",
					"api_key": "ackkey",
					"roles":   []string{"consistency-ack"},
				},
			}})
	c = suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("roles")
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "consistency.result",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "consistency.auto-check",
			"roles":    []string{"consistency-check"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "consistency.ack",
			"roles":    []string{"consistency-ack"},
		})

}

func (suite *consistencyTestSuite) TestCheckConsistency() {

	expJSON := `{
   "message": "Constistency information is not yet available",
   "code": 404
 }`

	request, _ := http.NewRequest("GET", "/api/v2/consistency", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 404 not found code
	suite.Equal(404, code, "Not Found")
	// Compare the expected and actual json response
	suite.Equal(expJSON, output, "Event Not Found")

	// Try an ack without having an auto-check event

	inputAck := `{
	"status": "OK",
	"message": "This is not a fault of monitoring"
	}`

	expNotAck := `{
   "message": "There is no auto check event yet to acknowledge",
   "code": 404
 }`

	request, _ = http.NewRequest("POST", "/api/v2/consistency/ack", strings.NewReader(inputAck))
	request.Header.Set("x-api-key", "ackkey")
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	// Check that we must have a 404 not foudn code
	suite.Equal(404, code, "not found")
	// Compare the expected and actual json response
	suite.Equal(expNotAck, output, "Ack event should not submit")

	// Post an auto check

	inputJSON := `{
	"status": "OK",
	"message": "Flapping items are lower"
	}`

	expJSON = `{
   "message": "The Auto Check event was posted succesfully",
   "code": 200
 }`

	request, _ = http.NewRequest("POST", "/api/v2/consistency/auto-check", strings.NewReader(inputJSON))
	request.Header.Set("x-api-key", "checkkey")
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	// Check that we must have a 404 not found code
	suite.Equal(200, code, "Not Found")
	// Compare the expected and actual json response
	suite.Equal(expJSON, output, "Event Not Found")

	// Get the check result

	expJSON2 := `{
 "status": "OK",
 "timestamp": ".*",
 "message": "Flapping items are lower"
}`

	request, _ = http.NewRequest("GET", "/api/v2/consistency", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	// Check that we must have a 404 not found code
	suite.Equal(200, code, "ok")
	// Compare the expected and actual json response
	suite.Regexp(expJSON2, output, "Get consistency check result")

	// Adding another auto check event that is critical

	inputCritical := `{
	"status": "CRITICAL",
	"message": "Flapping items are over 50 percent"
	}`

	request, _ = http.NewRequest("POST", "/api/v2/consistency/auto-check", strings.NewReader(inputCritical))
	request.Header.Set("x-api-key", "checkkey")
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	// Check that we must have a 404 not found code
	suite.Equal(200, code, "Not Found")
	// Compare the expected and actual json response
	suite.Equal(expJSON, output, "Event Not Found")

	// get the new consistency check result with verbosity

	expJSON3 := `{
 "status": "CRITICAL",
 "timestamp": ".*",
 "message": "Flapping items are over 50 percent",
 "auto_check_status": "CRITICAL",
 "auto_check_mesage": "Flapping items are over 50 percent",
 "auto_check_timestamp": ".*"
}`

	request, _ = http.NewRequest("GET", "/api/v2/consistency?verbose", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	// Check that we must have a 200 creation code
	suite.Equal(200, code, "ok")
	// Compare the expected and actual json response
	suite.Regexp(expJSON3, output, "Get consistency check result")

	// add an ack result

	// Adding another auto check event that is critical

	inputAck2 := `{
	"status": "OK",
	"message": "This is not a fault of monitoring",
	"timeout_hours": 1
	}`

	expJSON4 := `{
   "message": "The Ack event was posted succesfully",
   "code": 200
 }`

	request, _ = http.NewRequest("POST", "/api/v2/consistency/ack", strings.NewReader(inputAck2))
	request.Header.Set("x-api-key", "ackkey")
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	// Check that we must have a 200 creation code
	suite.Equal(200, code, "ok")
	// Compare the expected and actual json response
	suite.Equal(expJSON4, output, "Ack Event created")

	// get the new consistency check result with verbosity

	expJSON5 := `{
 "status": "OK",
 "timestamp": ".*",
 "message": "This is not a fault of monitoring",
 "auto_check_status": "CRITICAL",
 "auto_check_mesage": "Flapping items are over 50 percent",
 "auto_check_timestamp": ".*",
 "ack_status": "OK",
 "ack_message": "This is not a fault of monitoring",
 "ack_timestamp": ".*",
 "ack_timeout_hours": 1
}`

	request, _ = http.NewRequest("GET", "/api/v2/consistency?verbose", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	// Check that we must have a 200 creation code
	suite.Equal(200, code, "ok")
	// Compare the expected and actual json response
	suite.Regexp(expJSON5, output, "Get consistency check result")

	// add old date to timestamp and ack
	conCol := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection(conColName)
	conCol.UpdateOne(context.TODO(), bson.M{"_id": conID},
		bson.M{"$set": bson.M{"ack_timestamp": "2010-01-01T00:02:00Z"}})

	// get the new consistency check result with verbosity

	expJSON6 := `{
 "status": "CRITICAL",
 "timestamp": ".*",
 "message": "Flapping items are over 50 percent",
 "auto_check_status": "CRITICAL",
 "auto_check_mesage": "Flapping items are over 50 percent",
 "auto_check_timestamp": ".*"
}`

	request, _ = http.NewRequest("GET", "/api/v2/consistency?verbose", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	// Check that we must have a 200 creation code
	suite.Equal(200, code, "ok")
	// Compare the expected and actual json response
	suite.Regexp(expJSON6, output, "Get consistency check result")
}

// TearDownTest to tear down every test
func (suite *consistencyTestSuite) TearDownSuite() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenant1db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenant2db).Drop(context.TODO())

}

// TestTopologyTestSuite is responsible for calling the tests
func TestSuiteTopology(t *testing.T) {
	suite.Run(t, new(consistencyTestSuite))
}
