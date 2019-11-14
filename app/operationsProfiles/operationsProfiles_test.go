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
	"io/ioutil"
	"log"
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
func (suite *OperationsProfilesTestSuite) SetupSuite() {

	log.SetOutput(ioutil.Discard)

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
	    db = "AR_test_op_profiles"
	    `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.respUnauthorized = "Unauthorized"
	suite.tenantDbConf = config.MongoConfig{
		Host:     "localhost",
		Port:     27017,
		Db:       "AR_test_operations_profiles_tenant",
		Password: "pass",
		Username: "dbuser",
		Store:    "ar",
	}
	suite.clientkey = "123456"

	suite.confHandler = respond.ConfHandler{suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v2").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)
}

// This function runs before any test and setups the environment
func (suite *OperationsProfilesTestSuite) SetupTest() {

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
		bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"info": bson.M{
				"name":    "GUARDIANS",
				"email":   "email@something2",
				"website": "www.gotg.com",
				"created": "2015-10-20 02:08:04",
				"updated": "2015-10-20 02:08:04"},
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
					"roles":   []string{"editor"},
				},
				bson.M{
					"name":    "user2",
					"email":   "user2@email.com",
					"api_key": "USER2KEY",
					"roles":   []string{"editor"},
				},
			}})
	c.Insert(
		bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50d",
			"info": bson.M{
				"name":    "AVENGERS",
				"email":   "email@something2",
				"website": "www.gotg.com",
				"created": "2015-10-20 02:08:04",
				"updated": "2015-10-20 02:08:04"},
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
					"roles":   []string{"editor"},
				},
				bson.M{
					"name":    "user4",
					"email":   "user4@email.com",
					"api_key": "VIEWERKEY",
					"roles":   []string{"viewer"},
				},
			}})

	c = session.DB(suite.cfg.MongoDB.Db).C("roles")
	c.Insert(
		bson.M{
			"resource": "operationsProfiles.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "operationsProfiles.get",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "operationsProfiles.create",
			"roles":    []string{"editor"},
		})
	c.Insert(
		bson.M{
			"resource": "operationsProfiles.delete",
			"roles":    []string{"editor"},
		})
	c.Insert(
		bson.M{
			"resource": "operationsProfiles.update",
			"roles":    []string{"editor"},
		})

	// Seed database with operations profiles
	c = session.DB(suite.tenantDbConf.Db).C("operations_profiles")
	c.Insert(
		bson.M{
			"id":               "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
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
			"id":               "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
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

	request, _ := http.NewRequest("GET", "/api/v2/operations_profiles/wrong-id", strings.NewReader(jsonInput))
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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
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

	// Apply id to output template and check
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
  "id": "{{ID}}",
  "links": {
   "self": "https:///api/v2/operations_profiles/{{ID}}"
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

	// Grab ID from mongodb
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	defer session.Close()
	if err != nil {
		panic(err)
	}

	// Retrieve id from database
	var result map[string]interface{}
	c := session.DB(suite.tenantDbConf.Db).C("operations_profiles")

	c.Find(bson.M{"name": "tops1"}).One(&result)
	id := result["id"].(string)

	// Apply id to output template and check
	suite.Equal(strings.Replace(jsonOutput, "{{ID}}", id, 2), output, "Response body mismatch")

	// Check that actually the item has been created
	// Call List one with the specific ID
	request2, _ := http.NewRequest("GET", "/api/v2/operations_profiles/"+id, strings.NewReader(jsonInput))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)

	code2 := response2.Code
	output2 := response2.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(200, code2, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(strings.Replace(jsonCreated, "{{ID}}", id, 1), output2, "Response body mismatch")
}

func (suite *OperationsProfilesTestSuite) TestCreateNameAlreadyExists() {

	jsonInput := `{
   "name": "ops1",
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
  "message": "Conflict",
  "code": "409"
 },
 "errors": [
  {
   "message": "Conflict",
   "code": "409",
   "details": "Operations profile with the same name already exists"
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

	suite.Equal(409, code)
	suite.Equal(jsonOutput, output)
}

func (suite *OperationsProfilesTestSuite) TestUpdateNameAlreadyExists() {

	jsonInput := `{
   "name": "ops1",
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
  "message": "Conflict",
  "code": "409"
 },
 "errors": [
  {
   "message": "Conflict",
   "code": "409",
   "details": "Operations profile with the same name already exists"
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

	suite.Equal(409, code)
	suite.Equal(jsonOutput, output)
}

func (suite *OperationsProfilesTestSuite) TestUpdateBadJson() {

	jsonInput := `{
   "name": "yolo",
   "namespace": "testin
    `
	jsonOutput := `{
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

	request, _ := http.NewRequest("PUT", "/api/v2/operations_profiles/wrong-id", strings.NewReader(jsonInput))
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

	// Apply id to output template and check
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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
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

	// Apply id to output template and check
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

	request, _ := http.NewRequest("DELETE", "/api/v2/operations_profiles/wrong-id", strings.NewReader(jsonInput))
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

func (suite *OperationsProfilesTestSuite) TestDelete() {

	request, _ := http.NewRequest("DELETE", "/api/v2/operations_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Operations Profile Successfully Deleted",
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
	err = c.Find(bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"}).One(&result)

	suite.NotEqual(err, nil, "No not found error")
	suite.Equal(err.Error(), "not found", "No not found error")
}

func (suite *OperationsProfilesTestSuite) TestOptionsOperationsProfiles() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/operations_profiles", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, POST, DELETE, PUT, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v2/operations_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	headers = response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, POST, DELETE, PUT, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

}

func (suite *OperationsProfilesTestSuite) TestCreateForbidViewer() {

	jsonInput := `{
  "name": "test_profile",
  "namespace [
    `

	jsonOutput := `{
 "status": {
  "message": "Access to the resource is Forbidden",
  "code": "403"
 }
}`

	request, _ := http.NewRequest("POST", "/api/v2/operations_profiles", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", "VIEWERKEY")
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(403, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")
}

func (suite *OperationsProfilesTestSuite) TestUpdateForbidViewer() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Access to the resource is Forbidden",
  "code": "403"
 }
}`

	request, _ := http.NewRequest("PUT", "/api/v2/operations_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e007", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", "VIEWERKEY")
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(403, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *OperationsProfilesTestSuite) TestDeleteForbidViewer() {

	request, _ := http.NewRequest("DELETE", "/api/v2/operations_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
	request.Header.Set("x-api-key", "VIEWERKEY")
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Access to the resource is Forbidden",
  "code": "403"
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(403, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricProfileJSON, output, "Response body mismatch")
}

//TearDownTest to tear down every test
func (suite *OperationsProfilesTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}

	tenantDB := session.DB(suite.tenantDbConf.Db)
	mainDB := session.DB(suite.cfg.MongoDB.Db)

	cols, err := tenantDB.CollectionNames()
	for _, col := range cols {
		tenantDB.C(col).RemoveAll(nil)
	}

	cols, err = mainDB.CollectionNames()
	for _, col := range cols {
		mainDB.C(col).RemoveAll(nil)
	}

}

//TearDownTest to tear down every test
func (suite *OperationsProfilesTestSuite) TearDownSuite() {

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
