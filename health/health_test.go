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

package health

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
			"resource": "health",
			"roles":    []string{"editor", "viewer"},
		})
	c = suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("consistency")
	c.InsertOne(context.TODO(),
		bson.M{
			"_id":                  conID,
			"auto_check_status":    "OK",
			"auto_check_timestamp": "2025-07-01T00:00:00Z",
			"auto_check_message":   "No flapping items",
		})

}

func (suite *consistencyTestSuite) TestCheckConsistency() {

	expJSON := `{
 "status": "OK",
 "timestamp": ".*",
 "message": "No flapping items"
}`

	request, _ := http.NewRequest("GET", "/api/v2/health", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "all ok")
	// Compare the expected and actual json response
	suite.Regexp(expJSON, output, "Health event")

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
