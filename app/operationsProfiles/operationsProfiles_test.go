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

package operationsProfiles

import (
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

// This is a util. suite struct used in tests (see pkg "testify")
type OperationsProfilesTestSuite struct {
	suite.Suite
	cfg                       config.Config
	router                    *mux.Router
	confHandler               respond.ConfHandler
	tenantDbConf              config.MongoConfig
	clientkey                 string
	respRecomputationsCreated string
	respUnauthorized          string
}

// Setup the Test Environment
// This function runs before any test and setups the environment
func (suite *OperationsProfilesTestSuite) SetupTest() {

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
    db = "AR_test_recomputations"
    `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.respUnauthorized = "Unauthorized"
	suite.tenantDbConf = config.MongoConfig{
		Host:     "localhost",
		Port:     27017,
		Db:       "AR_test_aggregation_profiles_tenant",
		Password: "pass",
		Username: "dbuser",
		Store:    "ar",
	}
	suite.clientkey = "123456"

	suite.confHandler = respond.ConfHandler{suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v2").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)

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
		bson.M{"name": "FOO",
			"db_conf": []bson.M{

				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_FOO",
				},
				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_FOO",
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "user1",
					"email":   "user1@email.com",
					"api_key": "USER1KEY",
				},
				bson.M{
					"name":    "user2",
					"email":   "user2@email.com",
					"api_key": "USER2KEY",
				},
			}})
	c.Insert(
		bson.M{"name": "BAR",
			"db_conf": []bson.M{

				bson.M{
					// "store":    "ar",
					"server":   suite.tenantDbConf.Host,
					"port":     suite.tenantDbConf.Port,
					"database": suite.tenantDbConf.Db,
					"username": suite.tenantDbConf.Username,
					"password": suite.tenantDbConf.Password,
				},
				bson.M{
					"server":   suite.tenantDbConf.Host,
					"port":     suite.tenantDbConf.Port,
					"database": suite.tenantDbConf.Db,
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "user3",
					"email":   "user3@email.com",
					"api_key": suite.clientkey,
				},
				bson.M{
					"name":    "user4",
					"email":   "user4@email.com",
					"api_key": "USER4KEY",
				},
			}})
	// Seed database with operations profiles
	c = session.DB(suite.tenantDbConf.Db).C("operations_profiles")
	c.Insert(
		bson.M{
			"uuid":             "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			"name":             "ops1",
			"available_states": []string{"A,B,C"},
			"defaults": bson.M{
				"missing": "A",
				"down":    "B",
				"unknown": "C"},
			"operations": []bson.M{
				bson.M{
					"name": "AND",
					"truth_table": []bson.M{
						bson.M{
							"a": "A",
							"b": "B",
							"x": "B",
						},
						bson.M{
							"a": "A",
							"b": "C",
							"x": "C",
						},
						bson.M{
							"a": "B",
							"b": "C",
							"x": "C",
						}}},
				bson.M{
					"name": "OR",
					"truth_table": []bson.M{
						bson.M{
							"a": "A",
							"b": "B",
							"x": "A",
						},
						bson.M{
							"a": "A",
							"b": "C",
							"x": "A",
						},
						bson.M{
							"a": "B",
							"b": "C",
							"x": "B",
						}}},
			}})

	// Seed database with operations profiles
	c = session.DB(suite.tenantDbConf.Db).C("operations_profiles")
	c.Insert(
		bson.M{
			"uuid":             "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"name":             "ops2",
			"available_states": []string{"X,Y,Z"},
			"defaults": bson.M{
				"missing": "X",
				"down":    "Y",
				"unknown": "Z"},
			"operations": []bson.M{
				bson.M{
					"name": "AND",
					"truth_table": []bson.M{
						bson.M{
							"a": "X",
							"b": "Y",
							"x": "Y",
						},
						bson.M{
							"a": "X",
							"b": "Z",
							"x": "Z",
						},
						bson.M{
							"a": "Y",
							"b": "Z",
							"x": "Z",
						}}},
				bson.M{
					"name": "OR",
					"truth_table": []bson.M{
						bson.M{
							"a": "X",
							"b": "Y",
							"x": "X",
						},
						bson.M{
							"a": "X",
							"b": "Z",
							"x": "X",
						},
						bson.M{
							"a": "Y",
							"b": "Z",
							"x": "Y",
						}}},
			}})
}

func (suite *OperationsProfilesTestSuite) TestList() {

	request, _ := http.NewRequest("GET", "/api/v2/operations_profiles", strings.NewReader(""))
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
   "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "name": "ops1",
   "available_states": [
    "A,B,C"
   ],
   "defaults": {
    "down": "B",
    "missing": "A",
    "unknown": "C"
   },
   "operations": [
    {
     "name": "AND",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "B"
      },
      {
       "a": "A",
       "b": "C",
       "x": "C"
      },
      {
       "a": "B",
       "b": "C",
       "x": "C"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "A"
      },
      {
       "a": "A",
       "b": "C",
       "x": "A"
      },
      {
       "a": "B",
       "b": "C",
       "x": "B"
      }
     ]
    }
   ]
  },
  {
   "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "name": "ops2",
   "available_states": [
    "X,Y,Z"
   ],
   "defaults": {
    "down": "Y",
    "missing": "X",
    "unknown": "Z"
   },
   "operations": [
    {
     "name": "AND",
     "truth_table": [
      {
       "a": "X",
       "b": "Y",
       "x": "Y"
      },
      {
       "a": "X",
       "b": "Z",
       "x": "Z"
      },
      {
       "a": "Y",
       "b": "Z",
       "x": "Z"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "X",
       "b": "Y",
       "x": "X"
      },
      {
       "a": "X",
       "b": "Z",
       "x": "X"
      },
      {
       "a": "Y",
       "b": "Z",
       "x": "Y"
      }
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

func (suite *OperationsProfilesTestSuite) TestListQueryName() {

	request, _ := http.NewRequest("GET", "/api/v2/operations_profiles?name=ops1", strings.NewReader(""))
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
   "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "name": "ops1",
   "available_states": [
    "A,B,C"
   ],
   "defaults": {
    "down": "B",
    "missing": "A",
    "unknown": "C"
   },
   "operations": [
    {
     "name": "AND",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "B"
      },
      {
       "a": "A",
       "b": "C",
       "x": "C"
      },
      {
       "a": "B",
       "b": "C",
       "x": "C"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "A"
      },
      {
       "a": "A",
       "b": "C",
       "x": "A"
      },
      {
       "a": "B",
       "b": "C",
       "x": "B"
      }
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

func (suite *OperationsProfilesTestSuite) TestListOneNotFound() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Not Found",
  "code": "404",
  "details": "item with the specific UUID was not found on the server"
 }
}`

	request, _ := http.NewRequest("GET", "/api/v2/operations_profiles/wrong-uuid", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(404, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *OperationsProfilesTestSuite) TestListOne() {

	request, _ := http.NewRequest("GET", "/api/v2/operations_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
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
   "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "name": "ops1",
   "available_states": [
    "A,B,C"
   ],
   "defaults": {
    "down": "B",
    "missing": "A",
    "unknown": "C"
   },
   "operations": [
    {
     "name": "AND",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "B"
      },
      {
       "a": "A",
       "b": "C",
       "x": "C"
      },
      {
       "a": "B",
       "b": "C",
       "x": "C"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "A"
      },
      {
       "a": "A",
       "b": "C",
       "x": "A"
      },
      {
       "a": "B",
       "b": "C",
       "x": "B"
      }
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

func (suite *OperationsProfilesTestSuite) TestCreateBadJson() {

	jsonInput := `{
  "name": "test_profile",
  "namespace [
    `

	jsonOutput := `{
 "status": {
  "message": "Bad Request",
  "code": "400",
  "details": "Request Body contains malformed JSON, thus rendering the Request Bad"
 }
}`

	request, _ := http.NewRequest("POST", "/api/v2/operations_profiles", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(400, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *OperationsProfilesTestSuite) TestInvalidCreate() {

	jsonInput := `{
   "name": "ops1",
   "available_states": [
    "A","B","C","C"
   ],
   "defaults": {
    "down": "B",
    "missing": "FOO",
    "unknown": "C"
   },
   "operations": [
    {
     "name": "AND",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "B"
      },
      {
       "a": "A",
       "b": "C",
       "x": "C"
      },
      {
       "a": "B",
       "b": "BAR",
       "x": "C"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "A"
      },
      {
       "a": "A",
       "b": "C",
       "x": "A"
      },
      {
       "a": "B",
       "b": "CAR",
       "x": "B"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "A"
      },
      {
       "a": "A",
       "b": "C",
       "x": "A"
      },
      {
       "a": "B",
       "b": "C",
       "x": "B"
      }
     ]
    }
   ]
  }`

	jsonOutput := `{
 "status": {
  "message": "Validation Error",
  "code": "422"
 },
 "errors": [
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "State:C is duplicated"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Operation:OR is duplicated"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Default Missing State: FOO not in available States"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "In Operation: AND, statement member b: BAR contains undeclared state"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "In Operation: OR, statement member b: CAR contains undeclared state"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:A in operation: AND"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:B in operation: AND"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:C in operation: AND"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:A in operation: OR"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:B in operation: OR"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:C in operation: OR"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:A in operation: OR"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:B in operation: OR"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:C in operation: OR"
  }
 ]
}`

	request, _ := http.NewRequest("POST", "/api/v2/operations_profiles", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(422, code, "Internal Server Error")

	// Apply uuid to output template and check
	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *OperationsProfilesTestSuite) TestCreate() {

	jsonInput := `{
   "name": "tops1",
   "available_states": [
    "A","B","C"
   ],
   "defaults": {
    "down": "B",
    "missing": "A",
    "unknown": "C"
   },
   "operations": [
    {
     "name": "AND",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "B"
      },
      {
       "a": "A",
       "b": "C",
       "x": "C"
      },
      {
       "a": "B",
       "b": "C",
       "x": "C"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "A"
      },
      {
       "a": "A",
       "b": "C",
       "x": "A"
      },
      {
       "a": "B",
       "b": "C",
       "x": "B"
      }
     ]
    }
   ]
  }`

	jsonOutput := `{
 "status": {
  "message": "Operations Profile successfully created",
  "code": "201"
 },
 "data": {
  "uuid": "{{UUID}}",
  "links": {
   "self": "https:///api/v2/operations_profiles/{{UUID}}"
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
   "uuid": "{{UUID}}",
   "name": "tops1",
   "available_states": [
    "A",
    "B",
    "C"
   ],
   "defaults": {
    "down": "B",
    "missing": "A",
    "unknown": "C"
   },
   "operations": [
    {
     "name": "AND",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "B"
      },
      {
       "a": "A",
       "b": "C",
       "x": "C"
      },
      {
       "a": "B",
       "b": "C",
       "x": "C"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "A"
      },
      {
       "a": "A",
       "b": "C",
       "x": "A"
      },
      {
       "a": "B",
       "b": "C",
       "x": "B"
      }
     ]
    }
   ]
  }
 ]
}`

	request, _ := http.NewRequest("POST", "/api/v2/operations_profiles", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(201, code, "Internal Server Error")
	// Compare the expected and actual json response

	// Grab UUID from mongodb
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	defer session.Close()
	if err != nil {
		panic(err)
	}

	// Retrieve uuid from database
	var result map[string]interface{}
	c := session.DB(suite.tenantDbConf.Db).C("operations_profiles")

	c.Find(bson.M{"name": "tops1"}).One(&result)
	uuid := result["uuid"].(string)

	// Apply uuid to output template and check
	suite.Equal(strings.Replace(jsonOutput, "{{UUID}}", uuid, 2), output, "Response body mismatch")

	// Check that actually the item has been created
	// Call List one with the specific UUID
	request2, _ := http.NewRequest("GET", "/api/v2/operations_profiles/"+uuid, strings.NewReader(jsonInput))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)

	code2 := response2.Code
	output2 := response2.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(200, code2, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(strings.Replace(jsonCreated, "{{UUID}}", uuid, 1), output2, "Response body mismatch")
}

func (suite *OperationsProfilesTestSuite) TestUpdateBadJson() {

	jsonInput := `{
   "name": "yolo",
   "namespace": "testin
    `

	jsonOutput := `{
 "status": {
  "message": "Bad Request",
  "code": "400",
  "details": "Request Body contains malformed JSON, thus rendering the Request Bad"
 }
}`

	request, _ := http.NewRequest("PUT", "/api/v2/operations_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(400, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *OperationsProfilesTestSuite) TestUpdateNotFound() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Not Found",
  "code": "404",
  "details": "item with the specific UUID was not found on the server"
 }
}`

	request, _ := http.NewRequest("PUT", "/api/v2/operations_profiles/wrong-uuid", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(404, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *OperationsProfilesTestSuite) TestInvalidUpdate() {

	jsonInput := `{
   "name": "tops1",
   "available_states": [
    "A","B","C"
   ],
   "defaults": {
    "down": "D",
    "missing": "A",
    "unknown": "C"
   },
   "operations": [
    {
     "name": "AND",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "B"
      },
      {
       "a": "A",
       "b": "X",
       "x": "C"
      },
      {
       "a": "B",
       "b": "C",
       "x": "C"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "A"
      },
      {
       "a": "A",
       "b": "Z",
       "x": "A"
      },
      {
       "a": "B",
       "b": "C",
       "x": "B"
      }
     ]
    }
   ]
  }`

	jsonOutput := `{
 "status": {
  "message": "Validation Error",
  "code": "422"
 },
 "errors": [
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Default Down State: D not in available States"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "In Operation: AND, statement member b: X contains undeclared state"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "In Operation: OR, statement member b: Z contains undeclared state"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:C in operation: AND"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:C in operation: OR"
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/operations_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(422, code, "Internal Server Error")
	// Compare the expected and actual json response

	// Apply uuid to output template and check
	suite.Equal(jsonOutput, output, "Response body mismatch")
}

func (suite *OperationsProfilesTestSuite) TestUpdate() {

	jsonInput := `{
	 "name": "tops1",
	 "available_states": [
		"A","B","C"
	 ],
	 "defaults": {
		"down": "B",
		"missing": "A",
		"unknown": "C"
	 },
	 "operations": [
		{
		 "name": "AND",
		 "truth_table": [
			{
			 "a": "A",
			 "b": "B",
			 "x": "B"
			},
			{
			 "a": "A",
			 "b": "C",
			 "x": "C"
			},
			{
			 "a": "B",
			 "b": "C",
			 "x": "C"
			}
		 ]
		},
		{
		 "name": "OR",
		 "truth_table": [
			{
			 "a": "A",
			 "b": "B",
			 "x": "A"
			},
			{
			 "a": "A",
			 "b": "C",
			 "x": "A"
			},
			{
			 "a": "B",
			 "b": "C",
			 "x": "B"
			}
		 ]
		}
	 ]
	}`

	jsonOutput := `{
 "status": {
  "message": "Operations Profile successfully updated",
  "code": "200"
 }
}`

	jsonUpdated := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "name": "tops1",
   "available_states": [
    "A",
    "B",
    "C"
   ],
   "defaults": {
    "down": "B",
    "missing": "A",
    "unknown": "C"
   },
   "operations": [
    {
     "name": "AND",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "B"
      },
      {
       "a": "A",
       "b": "C",
       "x": "C"
      },
      {
       "a": "B",
       "b": "C",
       "x": "C"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "A"
      },
      {
       "a": "A",
       "b": "C",
       "x": "A"
      },
      {
       "a": "B",
       "b": "C",
       "x": "B"
      }
     ]
    }
   ]
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/operations_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response

	// Apply uuid to output template and check
	suite.Equal(jsonOutput, output, "Response body mismatch")

	// Check that the item has actually updated
	// run a list specific
	request2, _ := http.NewRequest("GET", "/api/v2/operations_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)

	code2 := response2.Code
	output2 := response2.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(200, code2, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(jsonUpdated, output2, "Response body mismatch")
}

func (suite *OperationsProfilesTestSuite) TestDeleteNotFound() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Not Found",
  "code": "404",
  "details": "item with the specific UUID was not found on the server"
 }
}`

	request, _ := http.NewRequest("DELETE", "/api/v2/operations_profiles/wrong-uuid", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(404, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *OperationsProfilesTestSuite) NotTestDelete() {

	request, _ := http.NewRequest("DELETE", "/api/v2/aggregation_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Aggregation Profile Successfully Deleted",
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
	c := session.DB(suite.tenantDbConf.Db).C("operations_profiles")
	err = c.Find(bson.M{"uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"}).One(&result)

	suite.NotEqual(err, nil, "No not found error")
	suite.Equal(err.Error(), "not found", "No not found error")
}

//TearDownTest to tear down every test
func (suite *OperationsProfilesTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

func TestOperationsProfilesTestSuite(t *testing.T) {
	suite.Run(t, new(OperationsProfilesTestSuite))
}
