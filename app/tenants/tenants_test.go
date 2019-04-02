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
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gcfg.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// This is a util. suite struct used in tests (see pkg "testify")
type TenantTestSuite struct {
	suite.Suite
	cfg                         config.Config
	respTenantCreated           string
	respTenantUpdated           string
	respTenantDeleted           string
	respTenantNotFound          string
	respUnauthorized            string
	respBadJSON                 string
	respTenantNameConflict      string
	respTenantUsersKeysConflict func(string) string
	clientkey                   string
	router                      *mux.Router
	confHandler                 respond.ConfHandler
}

// Setup the Test Environment
func (suite *TenantTestSuite) SetupSuite() {
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
	db = "argo_test_tenants"
	`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.confHandler = respond.ConfHandler{suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v2/admin").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)

	suite.clientkey = "S3CR3T"

	suite.respUnauthorized = `{
 "status": {
  "message": "Unauthorized",
  "code": "401",
  "details": "You need to provide a correct authentication token using the header 'x-api-key'"
 }
}`

	suite.respBadJSON = `{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "Bad Request",
   "code": "400",
   "details": "Request Body contains malformed JSON, thus rendering the Request Bad"
  }
 ]
}`

	suite.respTenantNotFound = `{
 "status": {
  "message": "Not Found",
  "code": "404"
 },
 "errors": [
  {
   "message": "Not Found",
   "code": "404",
   "details": "item with the specific ID was not found on the server"
  }
 ]
}`

	suite.respTenantNameConflict = `{
 "status": {
  "message": "Conflict",
  "code": "409"
 },
 "errors": [
  {
   "message": "Conflict",
   "code": "409",
   "details": "Tenant with same name already exists"
  }
 ]
}`

	suite.respTenantUsersKeysConflict = func(s string) string {

		resp := `{
 "status": {
  "message": "Conflict",
  "code": "409"
 },
 "errors": [
  {
   "message": "Conflict",
   "code": "409",
   "details": "More than one users found using the key: {{key}}"
  }
 ]
}`
		return strings.Replace(resp, "{{key}}", s, 1)
	}
}

// This function runs before any test and setups the environment
// A test configuration object is instantiated using a reference
// to testdb: AR_test_tenants. Also here is are instantiated some expected
// xml response validation messages (authorization,crud responses).
// Also the testdb is seeded with two tenants
// and with an authorization token:"S3CR3T"
func (suite *TenantTestSuite) SetupTest() {

	log.SetOutput(ioutil.Discard)

	// Connect to mongo testdb
	session, _ := mongo.OpenSession(suite.cfg.MongoDB)

	// Add authentication token to mongo testdb
	seedAuth := bson.M{"api_key": "S3CR3T"}
	seedResAuth := bson.M{"api_key": "R3STRICT3D", "restricted": true}
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "authentication", seedAuth)
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "authentication", seedResAuth)

	// seed mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB(suite.cfg.MongoDB.Db).C("roles")
	c.Insert(
		bson.M{
			"resource": "tenants.list",
			"roles":    []string{"super_admin", "super_admin_restricted"},
		})
	c.Insert(
		bson.M{
			"resource": "tenants.get",
			"roles":    []string{"super_admin"},
		})
	c.Insert(
		bson.M{
			"resource": "tenants.create",
			"roles":    []string{"super_admin"},
		})
	c.Insert(
		bson.M{
			"resource": "tenants.delete",
			"roles":    []string{"super_admin"},
		})
	c.Insert(
		bson.M{
			"resource": "tenants.update",
			"roles":    []string{"super_admin"},
		})
	c.Insert(
		bson.M{
			"resource": "tenants.update_status",
			"roles":    []string{"super_admin"},
		})
	c.Insert(
		bson.M{
			"resource": "tenants.get_status",
			"roles":    []string{"super_admin"},
		})
	// seed first tenant
	c = session.DB(suite.cfg.MongoDB.Db).C("tenants")
	c.Insert(bson.M{
		"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
		"info": bson.M{
			"name":    "AVENGERS",
			"email":   "email@something",
			"website": "www.avengers.com",
			"created": "2015-10-20 02:08:04",
			"updated": "2015-10-20 02:08:04"},
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
				"id":      "acb74194-553a-11e9-8647-d663bd873d93",
				"name":    "cap",
				"email":   "cap@email.com",
				"api_key": "C4PK3Y",
				"roles":   []string{"admin"}},
			bson.M{
				"id":      "acb74432-553a-11e9-8647-d663bd873d93",
				"name":    "thor",
				"email":   "thor@email.com",
				"api_key": "TH0RK3Y",
				"roles":   []string{"admin"}},
		}})

	// seed second tenant
	c.Insert(bson.M{
		"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
		"info": bson.M{
			"name":    "GUARDIANS",
			"email":   "email@something2",
			"website": "www.gotg.com",
			"created": "2015-10-20 02:08:04",
			"updated": "2015-10-20 02:08:04"},
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
				"id":      "acb7459a-553a-11e9-8647-d663bd873d93",
				"name":    "groot",
				"email":   "groot@email.com",
				"api_key": "GR00TK3Y",
				"roles":   []string{"admin"}},
			bson.M{
				"id":      "acb74702-553a-11e9-8647-d663bd873d93",
				"name":    "starlord",
				"email":   "starlord@email.com",
				"api_key": "ST4RL0RDK3Y",
				"roles":   []string{"admin"}},
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
      "info":{
				"name":"mutants",
				"email":"yo@yo",
				"website":"website"
			},
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
            "api_key":"X4V13R",
            "roles": [
                "admin"
            ]
          },
          {
            "name":"magneto",
            "email":"magneto@email.com",
            "api_key":"M4GN3T0",
            "roles": [
                "admin"
            ]
          }]
  }`

	jsonOutput := `{
 "status": {
  "message": "Tenant was successfully created",
  "code": "201"
 },
 "data": {
  "id": "{{ID}}",
  "links": {
   "self": "https:///api/v2/admin/tenants/{{ID}}"
  }
 }
}`

	jsonCreated := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "{{ID}}",
   "info": {
    "name": "mutants",
    "email": "yo@yo",
    "website": "website",
    "created": "{{TIMESTAMP}}",
    "updated": "{{TIMESTAMP}}"
   },
   "db_conf": [
    {
     "store": "ar",
     "server": "localhost",
     "port": 27017,
     "database": "ar_db",
     "username": "admin",
     "password": "3NCRYPT3D"
    },
    {
     "store": "status",
     "server": "localhost",
     "port": 27017,
     "database": "status_db",
     "username": "admin",
     "password": "3NCRYPT3D"
    }
   ],
   "users": [
    {
     "id": "{{UUID-1}}",
     "name": "xavier",
     "email": "xavier@email.com",
     "api_key": "X4V13R",
     "roles": [
      "admin"
     ]
    },
    {
     "id": "{{UUID-2}}",
     "name": "magneto",
     "email": "magneto@email.com",
     "api_key": "M4GN3T0",
     "roles": [
      "admin"
     ]
    }
   ]
  }
 ]
}`

	request, _ := http.NewRequest("POST", "/api/v2/admin/tenants", strings.NewReader(postData))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Grab ID from mongodb
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	defer session.Close()
	if err != nil {
		panic(err)
	}

	// Retrieve id from database
	var result = Tenant{}
	c := session.DB(suite.cfg.MongoDB.Db).C("tenants")

	c.Find(bson.M{"info.name": "mutants"}).One(&result)
	//id := result["id"].(string)
	//info := result["info"].(map[string]interface{})
	//timestamp := info["created"].(string)

	code := response.Code
	output := response.Body.String()

	suite.Equal(201, code, "Internal Server Error")
	// Apply id to output template and check
	suite.Equal(strings.Replace(jsonOutput, "{{ID}}", result.ID, 2), output, "Response body mismatch")

	// Check that actually the item has been created
	// Call List one with the specific ID
	request2, _ := http.NewRequest("GET", "/api/v2/admin/tenants/"+result.ID, strings.NewReader(""))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)

	code2 := response2.Code
	output2 := response2.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(200, code2, "Internal Server Error")

	jsonCreated = strings.Replace(jsonCreated, "{{ID}}", result.ID, 1)
	jsonCreated = strings.Replace(jsonCreated, "{{TIMESTAMP}}", result.Info.Created, 2)
	jsonCreated = strings.Replace(jsonCreated, "{{UUID-1}}", result.Users[0].ID, 1)
	jsonCreated = strings.Replace(jsonCreated, "{{UUID-2}}", result.Users[1].ID, 1)

	// Compare the expected and actual json response
	suite.Equal(jsonCreated, output2, "Response body mismatch")

}

func (suite *TenantTestSuite) TestCreateTenantDuplicateUsersKeys() {

	// create json input data for the request
	postData := `
  {
      "info":{
				"name":"mutants",
				"email":"yo@yo",
				"website":"website"
			},
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
            "api_key":"X4V13R",
            "roles": [
                "admin"
            ]
          },
          {
            "name":"magneto",
            "email":"magneto@email.com",
            "api_key":"X4V13R",
            "roles": [
                "admin"
            ]
          }]
  }`

	request, _ := http.NewRequest("POST", "/api/v2/admin/tenants", strings.NewReader(postData))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(409, code)
	suite.Equal(suite.respTenantUsersKeysConflict("X4V13R"), output)
}

// TestCreateTenantDuplicateUsersKeysWithDB tests the case where a duplicate key was found in the store
func (suite *TenantTestSuite) TestCreateTenantDuplicateUsersKeysWithDB() {

	// create json input data for the request
	postData := `
  {
      "info":{
				"name":"mutants",
				"email":"yo@yo",
				"website":"website"
			},
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
            "api_key":"X4V13R",
            "roles": [
                "admin"
            ]
          },
          {
            "name":"magneto",
            "email":"magneto@email.com",
            "api_key":"GR00TK3Y",
            "roles": [
                "admin"
            ]
          }]
  }`

	request, _ := http.NewRequest("POST", "/api/v2/admin/tenants", strings.NewReader(postData))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(409, code)
	suite.Equal(suite.respTenantUsersKeysConflict("GR00TK3Y"), output)
}

func (suite *TenantTestSuite) TestCreateTenantAlreadyExistingName() {

	// create json input data for the request
	postData := `
  {
      "info":{
				"name":"GUARDIANS",
				"email":"yo@yo",
				"website":"website"
			},
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
  }`

	request, _ := http.NewRequest("POST", "/api/v2/admin/tenants", strings.NewReader(postData))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(409, code)
	suite.Equal(suite.respTenantNameConflict, output)

}

func (suite *TenantTestSuite) TestUpdateTenantStatus() {

	// create json input data for the request
	putData := `
	{
      "ams": {
          "metric_data": {
              "ingestion": true,
              "publishing": true,
              "status_streaming": true,
              "messages_arrived": 100
          },
          "sync_data": {
              "ingestion": true,
              "publishing": true,
              "status_streaming": true,
              "messages_arrived": 200
          }
      },
      "hdfs": {
          "metric_data": true,
          "sync_data": {
          	"Critical": {
          			"downtimes": true,
          			"group_endpoints": true,
          			"blank_recomputation": true,
								"configuration_profile": true,
          			"group_groups": true,
          			"weights": true,
          			"operations_profile": true,
          			"metric_profile": true,
          			"aggregation_profile": true

          	}
          	}

      },
      "engine_config": true,
      "last_check": "2018-08-10T12:32:45Z"

}`

	jsonOutput := `{
 "status": {
  "message": "Tenant successfully updated",
  "code": "200"
 }
}`

	jsonUpdatedTenant := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "info": {
    "name": "GUARDIANS",
    "email": "email@something2",
    "website": "www.gotg.com",
    "created": "2015-10-20 02:08:04",
    "updated": "2015-10-20 02:08:04"
   },
   "status": {
    "total_status": true,
    "ams": {
     "metric_data": {
      "ingestion": true,
      "publishing": true,
      "status_streaming": true,
      "messages_arrived": 100
     },
     "sync_data": {
      "ingestion": true,
      "publishing": true,
      "status_streaming": true,
      "messages_arrived": 200
     }
    },
    "hdfs": {
     "metric_data": true,
     "sync_data": {
      "Critical": {
       "aggregation_profile": true,
       "blank_recomputation": true,
       "configuration_profile": true,
       "downtimes": true,
       "group_endpoints": true,
       "group_groups": true,
       "metric_profile": true,
       "operations_profile": true,
       "weights": true
      }
     }
    },
    "engine_config": true,
    "last_check": "2018-08-10T12:32:45Z"
   }
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/admin/tenants/6ac7d684-1f8e-4a02-a502-720e8f11e50c/status", strings.NewReader(putData))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(200, code, "Wrong code in response")
	suite.Equal(jsonOutput, output, "Response body mismatch")

	// Retrieve updated information
	request2, _ := http.NewRequest("GET", "/api/v2/admin/tenants/6ac7d684-1f8e-4a02-a502-720e8f11e50c/status", strings.NewReader(""))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)

	code2 := response2.Code
	output2 := response2.Body.String()

	suite.Equal(200, code2, "Wrong code in response")
	suite.Equal(jsonUpdatedTenant, output2, "Response body mismatch")
}

// TestUpdateTenant function implements testing the http PUT update tenant request.
// Request requires admin authentication and gets as input the name of the
// tenant to be updated and a json body with the update.
// After the operation succeeds is double-checked
// that the specific tenant has been updated
func (suite *TenantTestSuite) TestUpdateTenant() {

	// create json input data for the request
	putData := `
  {
      "info":{
				"name":"new_mutants",
				"email":"yo@yo",
				"website":"website"
			},
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
            "api_key":"X4V13R",
            "roles": [
                "admin"
            ]
          },
          {
            "name":"magneto",
            "email":"magneto@email.com",
            "api_key":"M4GN3T0",
            "roles": [
                "admin"
            ]
          }]
  }`

	jsonOutput := `{
 "status": {
  "message": "Tenant successfully updated",
  "code": "200"
 }
}`

	request, _ := http.NewRequest("PUT", "/api/v2/admin/tenants/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(putData))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *TenantTestSuite) TestUpdateTenantDuplicateUsersKeys() {

	// create json input data for the request
	putData := `
  {
      "info":{
				"name":"new_mutants",
				"email":"yo@yo",
				"website":"website"
			},
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
            "api_key":"X4V13R",
            "roles": [
                "admin"
            ]
          },
          {
            "name":"magneto",
            "email":"magneto@email.com",
            "api_key":"X4V13R",
            "roles": [
                "admin"
            ]
          }]
  }`

	request, _ := http.NewRequest("PUT", "/api/v2/admin/tenants/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(putData))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(409, code)
	// Compare the expected and actual xml response
	suite.Equal(suite.respTenantUsersKeysConflict("X4V13R"), output)

}

// TestUpdateTenantDuplicateUsersKeysWithDB tests the case where a duplicate key was found in the store
func (suite *TenantTestSuite) TestUpdateTenantDuplicateUsersKeysWithDB() {

	// create json input data for the request
	putData := `
  {
      "info":{
				"name":"new_mutants",
				"email":"yo@yo",
				"website":"website"
			},
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
            "api_key":"X4V13R",
            "roles": [
                "admin"
            ]
          },
          {
            "name":"magneto",
            "email":"magneto@email.com",
            "api_key":"TH0RK3Y",
            "roles": [
                "admin"
            ]
          }]
  }`

	request, _ := http.NewRequest("PUT", "/api/v2/admin/tenants/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(putData))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(409, code)
	// Compare the expected and actual xml response
	suite.Equal(suite.respTenantUsersKeysConflict("TH0RK3Y"), output)

}

func (suite *TenantTestSuite) TestUpdateTenantAlreadyExistingName() {

	// create json input data for the request
	postData := `
  {
      "info":{
				"name":"AVENGERS",
				"email":"yo@yo",
				"website":"website"
			},
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
  }`

	request, _ := http.NewRequest("PUT", "/api/v2/admin/tenants/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(postData))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(409, code)
	suite.Equal(suite.respTenantNameConflict, output)

}

// TestDeleteTenant function implements testing the http DELETE tenant request.
// Request requires admin authentication and gets as input the name of the
// tenant to be deleted. After the operation succeeds is double-checked
// that the deleted tenant is actually missing from the datastore
func (suite *TenantTestSuite) TestDeleteTenant() {

	request, _ := http.NewRequest("DELETE", "/api/v2/admin/tenants/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Tenant Successfully Deleted",
  "code": "200"
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricProfileJSON, output, "Response body mismatch")

	// check that the element has actually been Deleted
	// connect to mongodb
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	defer session.Close()
	if err != nil {
		panic(err)
	}
	// try to retrieve item
	var result map[string]interface{}
	c := session.DB(suite.cfg.MongoDB.Db).C("tenants")
	err = c.Find(bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"}).One(&result)

	suite.NotEqual(err, nil, "No not found error")
	suite.Equal(err.Error(), "not found", "No not found error")
}

// TestReadTeanants function implements the testing
// of the get request which retrieves all tenant information
func (suite *TenantTestSuite) TestListRestrictedTenants() {

	request, _ := http.NewRequest("GET", "/api/v2/admin/tenants", strings.NewReader(""))
	// emulate a restricted super admin user
	request.Header.Set("x-api-key", "R3STRICT3D")
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	profileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "info": {
    "name": "AVENGERS",
    "email": "email@something",
    "website": "www.avengers.com",
    "created": "2015-10-20 02:08:04",
    "updated": "2015-10-20 02:08:04"
   }
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "info": {
    "name": "GUARDIANS",
    "email": "email@something2",
    "website": "www.gotg.com",
    "created": "2015-10-20 02:08:04",
    "updated": "2015-10-20 02:08:04"
   }
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(profileJSON, output, "Response body mismatch")

}

// TestReadTeanants function implements the testing
// of the get request which retrieves all tenant information
func (suite *TenantTestSuite) TestListTenants() {

	request, _ := http.NewRequest("GET", "/api/v2/admin/tenants", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	profileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "info": {
    "name": "AVENGERS",
    "email": "email@something",
    "website": "www.avengers.com",
    "created": "2015-10-20 02:08:04",
    "updated": "2015-10-20 02:08:04"
   },
   "db_conf": [
    {
     "store": "ar",
     "server": "a.mongodb.org",
     "port": 27017,
     "database": "ar_db",
     "username": "admin",
     "password": "3NCRYPT3D"
    },
    {
     "store": "status",
     "server": "b.mongodb.org",
     "port": 27017,
     "database": "status_db",
     "username": "admin",
     "password": "3NCRYPT3D"
    }
   ],
   "users": [
    {
     "id": "acb74194-553a-11e9-8647-d663bd873d93",
     "name": "cap",
     "email": "cap@email.com",
     "api_key": "C4PK3Y",
     "roles": [
      "admin"
     ]
    },
    {
     "id": "acb74432-553a-11e9-8647-d663bd873d93",
     "name": "thor",
     "email": "thor@email.com",
     "api_key": "TH0RK3Y",
     "roles": [
      "admin"
     ]
    }
   ]
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "info": {
    "name": "GUARDIANS",
    "email": "email@something2",
    "website": "www.gotg.com",
    "created": "2015-10-20 02:08:04",
    "updated": "2015-10-20 02:08:04"
   },
   "db_conf": [
    {
     "store": "ar",
     "server": "a.mongodb.org",
     "port": 27017,
     "database": "ar_db",
     "username": "admin",
     "password": "3NCRYPT3D"
    },
    {
     "store": "status",
     "server": "b.mongodb.org",
     "port": 27017,
     "database": "status_db",
     "username": "admin",
     "password": "3NCRYPT3D"
    }
   ],
   "users": [
    {
     "id": "acb7459a-553a-11e9-8647-d663bd873d93",
     "name": "groot",
     "email": "groot@email.com",
     "api_key": "GR00TK3Y",
     "roles": [
      "admin"
     ]
    },
    {
     "id": "acb74702-553a-11e9-8647-d663bd873d93",
     "name": "starlord",
     "email": "starlord@email.com",
     "api_key": "ST4RL0RDK3Y",
     "roles": [
      "admin"
     ]
    }
   ]
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(profileJSON, output, "Response body mismatch")

}

func (suite *TenantTestSuite) TestListTenantStatus() {

	request, _ := http.NewRequest("GET", "/api/v2/admin/tenants/6ac7d684-1f8e-4a02-a502-720e8f11e50b/status", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	profileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "info": {
    "name": "AVENGERS",
    "email": "email@something",
    "website": "www.avengers.com",
    "created": "2015-10-20 02:08:04",
    "updated": "2015-10-20 02:08:04"
   },
   "status": {
    "total_status": false,
    "ams": {
     "metric_data": {
      "ingestion": false,
      "publishing": false,
      "status_streaming": false,
      "messages_arrived": 0
     },
     "sync_data": {
      "ingestion": false,
      "publishing": false,
      "status_streaming": false,
      "messages_arrived": 0
     }
    },
    "hdfs": {
     "metric_data": false
    },
    "engine_config": false,
    "last_check": ""
   }
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(profileJSON, output, "Response body mismatch")

}

// TestCreateUnauthorized function tests calling the create tenant request (POST) and
// providing a wrong api-key. The response should be unauthorized
func (suite *TenantTestSuite) TestCreateUnauthorized() {
	request, _ := http.NewRequest("POST", "/api/v2/admin/tenants", strings.NewReader(""))
	request.Header.Set("x-api-key", "FOO")
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	suite.Equal(401, code, "Internal Server Error")

	suite.Equal(suite.respUnauthorized, output, "Response body mismatch")
}

// TestUpdateUnauthorized function tests calling the update tenant request (PUT)
// and providing  a wrong api-key. The response should be unauthorized
func (suite *TenantTestSuite) TestUpdateUnauthorized() {
	request, _ := http.NewRequest("PUT", "/api/v2/admin/tenants/id", strings.NewReader(""))
	request.Header.Set("x-api-key", "FOO")
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(401, code, "Internal Server Error")
	suite.Equal(suite.respUnauthorized, output, "Response body mismatch")

}

// TestDeleteUnauthorized function tests calling the remove tenant request (DELETE)
// and providing a wrong api-key. The response should be unauthorized
func (suite *TenantTestSuite) TestDeleteUnauthorized() {
	request, _ := http.NewRequest("DELETE", "/api/v2/admin/tenants/id", strings.NewReader(""))
	request.Header.Set("x-api-key", "FOO")
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(401, code, "Internal Server Error")
	suite.Equal(suite.respUnauthorized, output, "Response body mismatch")
}

// TestCreateBadJson tests calling the create tenant request (POST) and providing
// bad json input. The response should be malformed json
func (suite *TenantTestSuite) TestCreateBadJson() {
	jsonInput := "{bad json:{}"
	request, _ := http.NewRequest("POST", "/api/v2/admin/tenants", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(400, code, "Internal Server Error")
	suite.Equal(suite.respBadJSON, output, "Response body mismatch")

}

// TestUpdateBadJson tests calling the update tenant request (PUT) and providing
// bad json input. The response should be malformed json
func (suite *TenantTestSuite) TestUpdateBadJson() {
	jsonInput := "{bad json:{}"
	request, _ := http.NewRequest("PUT", "/api/v2/admin/tenants/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(400, code, "Internal Server Error")
	suite.Equal(suite.respBadJSON, output, "Response body mismatch")
}

// TestListOneNotFound tests calling the http (GET) tenant info request
// and provide a non-existing tenant name. The response should be tenant not found
func (suite *TenantTestSuite) TestListOneNotFound() {
	request, _ := http.NewRequest("DELETE", "/api/v2/admin/tenants/BADID", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(404, code, "Internal Server Error")
	suite.Equal(suite.respTenantNotFound, output, "Response body mismatch")
}

// TestUpdateNotFound tests calling the http (PUT) update tenant request
// and provide a non-existing tenant name. The response should be tenant not found
func (suite *TenantTestSuite) TestUpdateNotFound() {
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "/api/v2/admin/tenants/BADID", strings.NewReader("{}"))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", suite.clientkey)

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(404, code, "Internal Server Error")
	suite.Equal(suite.respTenantNotFound, output, "Response body mismatch")
}

// TestDeleteNotFound tests calling the http (PUT) update tenant request
// and provide a non-existing tenant name. The response should be tenant not found
func (suite *TenantTestSuite) TestDeleteNotFound() {
	request, _ := http.NewRequest("DELETE", "/api/v2/admin/tenants/id", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(404, code, "Internal Server Error")
	suite.Equal(suite.respTenantNotFound, output, "Response body mismatch")
}

func (suite *TenantTestSuite) TestOptionsTenant() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/admin/tenants", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET,POST,PUT,DELETE,OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")
}

//TearDownTest to tear down every test
func (suite *TenantTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}

	mainDB := session.DB(suite.cfg.MongoDB.Db)

	cols, err := mainDB.CollectionNames()
	for _, col := range cols {
		mainDB.C(col).RemoveAll(nil)
	}

}

//TearDownTest to tear down every test
func (suite *TenantTestSuite) TearDownSuite() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}

	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

// This is the first function called when go test is issued
func TestTenantsSuite(t *testing.T) {
	suite.Run(t, new(TenantTestSuite))
}
