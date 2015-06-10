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

package tenants

import (
	"net/http"
	"strings"
	"testing"

	"code.google.com/p/gcfg"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// This is a util. suite struct used in tests (see pkg "testify")
type TenantTestSuite struct {
	suite.Suite
	cfg                config.Config
	respTenantCreated  string
	respTenantUpdated  string
	respTenantDeleted  string
	respTenantNotFound string
	respUnauthorized   string
	respBadJSON        string
}

// Setup the Test Environment
// This function runs before any test and setups the environment
// A test configuration object is instantiated using a reference
// to testdb: AR_test_tenants. Also here is are instantiated some expected
// xml response validation messages (authorization,crud responses).
// Also the testdb is seeded with two tenants
// and with an authorization token:"S3CR3T"
func (suite *TenantTestSuite) SetupTest() {

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
    db = "argo_test_tenants"
`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.respTenantCreated = " <root>\n" +
		"   <Message>Tenant was successfully created</Message>\n </root>"

	suite.respTenantUpdated = " <root>\n" +
		"   <Message>Tenant was successfully updated</Message>\n </root>"

	suite.respTenantDeleted = " <root>\n" +
		"   <Message>Tenant was successfully deleted</Message>\n </root>"

	suite.respTenantNotFound = " <root>\n" +
		"   <Message>Tenant not found</Message>\n </root>"

	suite.respBadJSON = " <root>\n" +
		"   <Message>Malformated json input data</Message>\n </root>"

	suite.respUnauthorized = "Unauthorized"

	// Connect to mongo testdb
	session, _ := mongo.OpenSession(suite.cfg.MongoDB)

	// Add authentication token to mongo testdb
	seedAuth := bson.M{"api_key": "S3CR3T"}
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "authentication", seedAuth)

	// seed mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// seed first tenant
	c := session.DB(suite.cfg.MongoDB.Db).C("tenants")
	c.Insert(bson.M{
		"name": "AVENGERS",
		"db_conf": []bson.M{
			bson.M{
				"store":    "ar",
				"server":   "a.mongodb.org",
				"port":     27017,
				"database": "ar_db",
				"username": "admin",
				"password": "3NCRYPT3D"},
			bson.M{
				"store":    "status",
				"server":   "b.mongodb.org",
				"port":     27017,
				"database": "status_db",
				"username": "admin",
				"password": "3NCRYPT3D"},
		},
		"users": []bson.M{
			bson.M{
				"name":    "cap",
				"email":   "cap@email.com",
				"api_key": "C4PK3Y"},
			bson.M{
				"name":    "thor",
				"email":   "thor@email.com",
				"api_key": "TH0RK3Y"},
		}})

	// seed second tenant
	c.Insert(bson.M{
		"name": "GUARDIANS",
		"db_conf": []bson.M{
			bson.M{
				"store":    "ar",
				"server":   "a.mongodb.org",
				"port":     27017,
				"database": "ar_db",
				"username": "admin",
				"password": "3NCRYPT3D"},
			bson.M{
				"store":    "status",
				"server":   "b.mongodb.org",
				"port":     27017,
				"database": "status_db",
				"username": "admin",
				"password": "3NCRYPT3D"},
		},
		"users": []bson.M{
			bson.M{
				"name":    "groot",
				"email":   "groot@email.com",
				"api_key": "GR00TK3Y"},
			bson.M{
				"name":    "starlord",
				"email":   "starlord@email.com",
				"api_key": "ST4RL0RDK3Y"},
		}})
}

// TestCreateTenant function implements testing the http POST create tenant request.
// Request requires admin authentication and gets as input a json body containing
// all the available information to be added to the datastore
// After the operation succeeds is double-checked
// that the newly created tenant is correctly retrieved
func (suite *TenantTestSuite) TestCreateTenant() {

	// create json input data for the request
	postData := `
  {
      "name": "MUTANTS",
      "db_conf": [
        {
          "store":"ar",
          "server":"localhost",
          "port":27017,
          "database":"ar_db",
          "username":"admin",
          "password":"3NCRYPT3D"
        },
        {
          "store":"status",
          "server":"localhost",
          "port":27017,
          "database":"status_db",
          "username":"admin",
          "password":"3NCRYPT3D"
        }],
      "users": [
          {
            "name":"xavier",
            "email":"xavier@email.com",
            "api_key":"X4V13R"
          },
          {
            "name":"magneto",
            "email":"magneto@email.com",
            "api_key":"M4GN3T0"
          }]
  }
    `
	// Prepare the request object
	request, _ := http.NewRequest("POST", "", strings.NewReader(postData))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Create(request, suite.cfg)

	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.respTenantCreated, string(output), "Response body mismatch")

	// Double check that you read the newly inserted profile
	// Create a string literal of the expected xml Response
	respXML := `<root>
 <tenant name="MUTANTS">
  <db_confs>
   <db_conf store="ar" server="localhost" port="27017" database="ar_db" username="admin" password="3NCRYPT3D"></db_conf>
   <db_conf store="status" server="localhost" port="27017" database="status_db" username="admin" password="3NCRYPT3D"></db_conf>
  </db_confs>
  <users>
   <user name="xavier" email="xavier@email.com" api_key="X4V13R"></user>
   <user name="magneto" email="magneto@email.com" api_key="M4GN3T0"></user>
  </users>
 </tenant>
</root>`

	// Prepare the request object using tenant name as urlvar in url path
	request, _ = http.NewRequest("GET", "/api/v1/tenants/MUTANTS", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")
	// Pass request to controller calling List() handler method
	code, _, output, _ = ListOne(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respXML, string(output), "Response body mismatch")
}

// TestUpdateTenant function implements testing the http PUT update tenant request.
// Request requires admin authentication and gets as input the name of the
// tenant to be updated and a json body with the update.
// After the operation succeeds is double-checked
// that the specific tenant has been updated
func (suite *TenantTestSuite) TestUpdateTenant() {

	// create json input data for the request
	postData := `
  {
      "name": "AVENGERS_modified",
      "db_conf": [
        {
          "store":"ar_mod",
          "server":"localhost",
          "port":27017,
          "database":"ar_db",
          "username":"admin",
          "password":"3NCRYPT3D"
        },
        {
          "store":"status_mod",
          "server":"localhost",
          "port":27017,
          "database":"status_db",
          "username":"admin",
          "password":"3NCRYPT3D"
        }],
      "users": [
          {
            "name":"thor",
            "email":"thor@email.com",
            "api_key":"TH0RK3Y"
          },
          {
            "name":"cap",
            "email":"cap@email.com",
            "api_key":"C4PK3Y"
          },
          {
            "name":"hulk",
            "email":"hulk@email.com",
            "api_key":"HULKK3Y"
          }]
  }`
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "/api/v1/tenants/AVENGERS", strings.NewReader(postData))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Update(request, suite.cfg)

	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.respTenantUpdated, string(output), "Response body mismatch")

	// Double check that you read the newly updated profile
	// Create a string literal of the expected xml Response
	respXML := `<root>
 <tenant name="AVENGERS_modified">
  <db_confs>
   <db_conf store="ar_mod" server="localhost" port="27017" database="ar_db" username="admin" password="3NCRYPT3D"></db_conf>
   <db_conf store="status_mod" server="localhost" port="27017" database="status_db" username="admin" password="3NCRYPT3D"></db_conf>
  </db_confs>
  <users>
   <user name="thor" email="thor@email.com" api_key="TH0RK3Y"></user>
   <user name="cap" email="cap@email.com" api_key="C4PK3Y"></user>
   <user name="hulk" email="hulk@email.com" api_key="HULKK3Y"></user>
  </users>
 </tenant>
</root>`

	// Prepare the request object using tenant name as urlvar in url path
	request, _ = http.NewRequest("GET", "/api/v1/tenants/AVENGERS_modified", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")
	// Pass request to controller calling List() handler method
	code, _, output, _ = ListOne(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respXML, string(output), "Response body mismatch")
}

// TestDeleteTenant function implements testing the http DELETE tenant request.
// Request requires admin authentication and gets as input the name of the
// tenant to be deleted. After the operation succeeds is double-checked
// that the deleted tenant is actually missing from the datastore
func (suite *TenantTestSuite) TestDeleteTenant() {

	// Prepare the request object
	request, _ := http.NewRequest("DELETE", "/api/v1/tenants/AVENGERS", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Delete(request, suite.cfg)

	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.respTenantDeleted, string(output), "Response body mismatch")

	// Double check that the tenant is actually removed when you try
	// to retrieve it's information by name
	// Prepare the request object using tenant name as urlvar in url path
	request, _ = http.NewRequest("GET", "/api/v1/tenants/AVENGERS", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")
	// Pass request to controller calling List() handler method
	code, _, output, _ = ListOne(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(400, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(suite.respTenantNotFound, string(output), "Response body mismatch")
}

// TestReadOneTeanant function implements the testing
// of the get request which retrieves information
// about a specific tenant (using it's name as input)
func (suite *TenantTestSuite) TestReadOneTenant() {

	// Open a session to mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Create a string literal of the expected xml Response
	respXML := `<root>
 <tenant name="GUARDIANS">
  <db_confs>
   <db_conf store="ar" server="a.mongodb.org" port="27017" database="ar_db" username="admin" password="3NCRYPT3D"></db_conf>
   <db_conf store="status" server="b.mongodb.org" port="27017" database="status_db" username="admin" password="3NCRYPT3D"></db_conf>
  </db_confs>
  <users>
   <user name="groot" email="groot@email.com" api_key="GR00TK3Y"></user>
   <user name="starlord" email="starlord@email.com" api_key="ST4RL0RDK3Y"></user>
  </users>
 </tenant>
</root>`

	// Prepare the request object using tenant name as urlvar in url path
	request, _ := http.NewRequest("GET", "/api/v1/tenants/GUARDIANS", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")
	// Pass request to controller calling List() handler method
	code, _, output, _ := ListOne(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respXML, string(output), "Response body mismatch")
}

// TestReadTeanants function implements the testing
// of the get request which retrieves all tenant information
func (suite *TenantTestSuite) TestReadTenants() {

	// Open a session to mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Create a string literal of the expected xml Response
	respXML := `<root>
 <tenant name="AVENGERS">
  <db_confs>
   <db_conf store="ar" server="a.mongodb.org" port="27017" database="ar_db" username="admin" password="3NCRYPT3D"></db_conf>
   <db_conf store="status" server="b.mongodb.org" port="27017" database="status_db" username="admin" password="3NCRYPT3D"></db_conf>
  </db_confs>
  <users>
   <user name="cap" email="cap@email.com" api_key="C4PK3Y"></user>
   <user name="thor" email="thor@email.com" api_key="TH0RK3Y"></user>
  </users>
 </tenant>
 <tenant name="GUARDIANS">
  <db_confs>
   <db_conf store="ar" server="a.mongodb.org" port="27017" database="ar_db" username="admin" password="3NCRYPT3D"></db_conf>
   <db_conf store="status" server="b.mongodb.org" port="27017" database="status_db" username="admin" password="3NCRYPT3D"></db_conf>
  </db_confs>
  <users>
   <user name="groot" email="groot@email.com" api_key="GR00TK3Y"></user>
   <user name="starlord" email="starlord@email.com" api_key="ST4RL0RDK3Y"></user>
  </users>
 </tenant>
</root>`

	// Prepare the request object
	request, _ := http.NewRequest("GET", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")
	// Pass request to controller calling List() handler method
	code, _, output, _ := List(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respXML, string(output), "Response body mismatch")
}

// TestCreateUnauthorized function tests calling the create tenant request (POST) and
// providing a wrong api-key. The response should be unauthorized
func (suite *TenantTestSuite) TestCreateUnauthorized() {
	// Prepare the request object (use id2 for path)
	request, _ := http.NewRequest("POST", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "F00T0K3N")

	// Execute the request in the controller
	code, _, output, _ := Create(request, suite.cfg)

	suite.Equal(401, code, "Internal Server Error")
	suite.Equal(suite.respUnauthorized, string(output), "Response body mismatch")
}

// TestUpdateUnauthorized function tests calling the update tenant request (PUT)
// and providing  a wrong api-key. The response should be unauthorized
func (suite *TenantTestSuite) TestUpdateUnauthorized() {
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "", strings.NewReader("{}"))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "F00T0K3N")

	// Execute the request in the controller
	code, _, output, _ := Update(request, suite.cfg)

	suite.Equal(401, code, "Internal Server Error")
	suite.Equal(suite.respUnauthorized, string(output), "Response body mismatch")
}

// TestDeleteUnauthorized function tests calling the remove tenant request (DELETE)
// and providing a wrong api-key. The response should be unauthorized
func (suite *TenantTestSuite) TestDeleteUnauthorized() {
	// Prepare the request object
	request, _ := http.NewRequest("DELETE", "", strings.NewReader("{}"))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "F00T0K3N")

	// Execute the request in the controller
	code, _, output, _ := Delete(request, suite.cfg)

	suite.Equal(401, code, "Internal Server Error")
	suite.Equal(suite.respUnauthorized, string(output), "Response body mismatch")
}

// TestCreateBadJson tests calling the create tenant request (POST) and providing
// bad json input. The response should be malformed json
func (suite *TenantTestSuite) TestCreateBadJson() {
	// Prepare the request object
	request, _ := http.NewRequest("POST", "/api/v1/tenants/AVENGERS", strings.NewReader("{ bad json"))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Create(request, suite.cfg)

	suite.Equal(400, code, "Internal Server Error")
	suite.Equal(suite.respBadJSON, string(output), "Response body mismatch")
}

// TestUpdateBadJson tests calling the update tenant request (PUT) and providing
// bad json input. The response should be malformed json
func (suite *TenantTestSuite) TestUpdateBadJson() {
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "/api/v1/tenants/AVENGERS", strings.NewReader("{ bad json"))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Update(request, suite.cfg)

	suite.Equal(400, code, "Internal Server Error")
	suite.Equal(suite.respBadJSON, string(output), "Response body mismatch")
}

// TestListOneNotFound tests calling the http (GET) tenant info request
// and provide a non-existing tenant name. The response should be tenant not found
func (suite *TenantTestSuite) TestListOneNotFound() {
	// Prepare the request object
	request, _ := http.NewRequest("GET", "/api/v1/tenants/BADNAME", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := ListOne(request, suite.cfg)

	suite.Equal(400, code, "Internal Server Error")
	suite.Equal(suite.respTenantNotFound, string(output), "Response body mismatch")
}

// TestUpdateNotFound tests calling the http (PUT) update tenant request
// and provide a non-existing tenant name. The response should be tenant not found
func (suite *TenantTestSuite) TestUpdateNotFound() {
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "/api/v1/tenants/BADNAME", strings.NewReader("{}"))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Update(request, suite.cfg)

	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.respTenantNotFound, string(output), "Response body mismatch")
}

// TestDeleteNotFound tests calling the http (PUT) update tenant request
// and provide a non-existing tenant name. The response should be tenant not found
func (suite *TenantTestSuite) TestDeleteNotFound() {
	// Prepare the request object
	request, _ := http.NewRequest("DELETE", "/api/v1/tenants/BADNAME", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Delete(request, suite.cfg)

	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.respTenantNotFound, string(output), "Response body mismatch")
}

// This function is actually called in the end of all tests
// and clears the test environment.
// Mainly it's purpose is to drop the testdb
func (suite *TenantTestSuite) TearDownTest() {

	session, _ := mongo.OpenSession(suite.cfg.MongoDB)

	session.DB("argo_test_tenants").DropDatabase()

}

// This is the first function called when go test is issued
func TestTenantsSuite(t *testing.T) {
	suite.Run(t, new(TenantTestSuite))
}
