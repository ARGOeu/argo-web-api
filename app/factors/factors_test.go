/*
 * Copyright (c) 2014 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
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
 * or implied, of either GRNET S.A., SRCE or IN2P3 CNRS Computing
 * Centre
 *
 * The work represented by this source file is partially funded by
 * the EGI-InSPIRE project through the European Commission's 7th
 * Framework Programme (contract # INFSO-RI-261323)
 */

package factors

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

// FactorsTestSuite is a utility suite struct used in tests
type FactorsTestSuite struct {
	suite.Suite
	cfg                 config.Config
	tenantDbConf        config.MongoConfig
	router              *mux.Router
	confHandler         respond.ConfHandler
	respNokeyprovided   string
	respUnauthorized    string
	respFactorsListXML  string
	respFactorsListJSON string
}

func (suite *FactorsTestSuite) SetupSuite() {

	log.SetOutput(io.Discard)

	const coreConfig = `
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
	    db = "argo_core_test_factors"
	`
	_ = gcfg.ReadStringInto(&suite.cfg, coreConfig)
	suite.respNokeyprovided = "404 page not found"
	suite.respUnauthorized = `{
 "status": {
  "message": "Unauthorized",
  "code": "401",
  "details": "You need to provide a correct authentication token using the header 'x-api-key'"
 }
}`

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v2").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)

	// TODO: I don't like it here that I rewrite the test data.
	// However, this is a test for factors, not for AuthenticateTenant function.
	suite.tenantDbConf.Host = "127.0.0.1"
	suite.tenantDbConf.Port = 27017
	suite.tenantDbConf.Db = "AR_test"
}

// SetupTest will bootstrap and provide the testing environment
func (suite *FactorsTestSuite) SetupTest() {

	c := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("tenants")
	// Add authentication token to mongo coredb
	seedAuth := bson.M{"name": "TEST",
		"db_conf": []bson.M{{"server": "127.0.0.1", "port": 27017, "database": "AR_test"}},
		"users":   []bson.M{{"name": "Jack Doe", "email": "jack.doe@example.com", "api_key": "secret", "roles": []string{"viewer"}}}}
	c.InsertOne(context.TODO(), seedAuth)

	// Add a few factors in collection
	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("weights")
	c.InsertOne(context.TODO(), bson.M{"hepspec": 14595, "name": "CIEMAT-LCG2"})
	c.InsertOne(context.TODO(), bson.M{"hepspec": 1019, "name": "CFP-IST"})
	c.InsertOne(context.TODO(), bson.M{"hepspec": 5406, "name": "CETA-GRID"})

	c = suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("roles")
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "factors.list",
			"roles":    []string{"editor", "viewer"},
		})

}

// TestListFactors will run unit tests against the List function
func (suite *FactorsTestSuite) TestListFactors() {

	suite.respFactorsListXML = `<root>
 <Factor site="CETA-GRID" weight="5406"></Factor>
 <Factor site="CFP-IST" weight="1019"></Factor>
 <Factor site="CIEMAT-LCG2" weight="14595"></Factor>
</root>`

	suite.respFactorsListJSON = `{
 "factors": [
  {
   "site": "CETA-GRID",
   "weight": "5406"
  },
  {
   "site": "CFP-IST",
   "weight": "1019"
  },
  {
   "site": "CIEMAT-LCG2",
   "weight": "14595"
  }
 ]
}`

	// Prepare the request object
	request, _ := http.NewRequest("GET", "/api/v2/factors", strings.NewReader(""))
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "secret")
	request.Header.Set("Accept", "application/xml")
	// Execute the request in the controller
	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	code := response.Code
	output := response.Body.String()
	suite.Equal(200, code, "Return status code mismatch")
	suite.Equal(suite.respFactorsListXML, string(output), "Response body mismatch")

	// Prepare new request object
	request, _ = http.NewRequest("GET", "/api/v2/factors", strings.NewReader(""))
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "secret")
	request.Header.Set("Accept", "application/json")
	// Execute the request in the controller
	response = httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	code = response.Code
	output = response.Body.String()
	suite.Equal(200, code, "Return status code mismatch")
	suite.Equal(suite.respFactorsListJSON, string(output), "Response body mismatch")

	// Prepare new request object
	request, _ = http.NewRequest("GET", "/api/v2/factors", strings.NewReader(""))
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "wrongkey")
	request.Header.Set("Accept", "application/json")
	// Execute the request in the controller
	response = httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	code = response.Code
	output = response.Body.String()
	suite.Equal(401, code, "Return status code mismatch")
	suite.Equal(suite.respUnauthorized, string(output), "Response body mismatch")

}

func (suite *FactorsTestSuite) TestOptionsFactors() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/factors", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

}

func (suite *FactorsTestSuite) TestTrailingSlashFactors() {

	request, _ := http.NewRequest("GET", "/api/v2/factors/", strings.NewReader(""))
	request.Header.Set("x-api-key", "secret")
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code

	suite.Equal(404, code, "Error in response code")

}

// TearDownSuite do things after each test ends
func (suite *FactorsTestSuite) TearDownTest() {

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

// TearDownSuite do things after suite ends
func (suite *FactorsTestSuite) TearDownSuite() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
}

func TestSuiteFactors(t *testing.T) {
	suite.Run(t, new(FactorsTestSuite))
}
