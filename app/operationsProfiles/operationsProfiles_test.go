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
			"uuid":              "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			"name":              "critical",
			"namespace":         "test",
			"endpoint_group":    "sites",
			"metric_operation":  "AND",
			"profile_operation": "AND",
			"metric_profile": bson.M{
				"name": "roc.critical",
				"uuid": "5637d684-1f8e-4a02-a502-720e8f11e432",
			},
			"groups": []bson.M{
				bson.M{"name": "compute",
					"operation": "OR",
					"services": []bson.M{
						bson.M{
							"name":      "CREAM-CE",
							"operation": "AND",
						},
						bson.M{
							"name":      "ARC-CE",
							"operation": "AND",
						},
					}},
				bson.M{"name": "storage",
					"operation": "OR",
					"services": []bson.M{
						bson.M{
							"name":      "SRMv2",
							"operation": "AND",
						},
						bson.M{
							"name":      "SRM",
							"operation": "AND",
						},
					}},
			}})
	c.Insert(
		bson.M{
			"uuid":              "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"name":              "cloud",
			"namespace":         "test",
			"endpoint_group":    "sites",
			"metric_operation":  "AND",
			"profile_operation": "AND",
			"metric_profile": bson.M{
				"name": "roc.critical",
				"uuid": "5637d684-1f8e-4a02-a502-720e8f11e432",
			},
			"groups": []bson.M{
				bson.M{"name": "compute",
					"operation": "OR",
					"services": []bson.M{
						bson.M{
							"name":      "SERVICEA",
							"operation": "AND",
						},
						bson.M{
							"name":      "SERVICEB",
							"operation": "AND",
						},
					}},
				bson.M{"name": "images",
					"operation": "OR",
					"services": []bson.M{
						bson.M{
							"name":      "SERVICEC",
							"operation": "AND",
						},
						bson.M{
							"name":      "SERVICED",
							"operation": "AND",
						},
					}},
			}})

}

func (suite *OperationsProfilesTestSuite) TestList() {

	request, _ := http.NewRequest("GET", "/api/v2/aggregation_profiles", strings.NewReader(""))
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
   "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "name": "cloud",
   "namespace": "test",
   "endpoint_group": "sites",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "roc.critical",
    "uuid": "5637d684-1f8e-4a02-a502-720e8f11e432"
   },
   "groups": [
    {
     "name": "compute",
     "operation": "OR",
     "services": [
      {
       "name": "SERVICEA",
       "operation": "AND"
      },
      {
       "name": "SERVICEB",
       "operation": "AND"
      }
     ]
    },
    {
     "name": "images",
     "operation": "OR",
     "services": [
      {
       "name": "SERVICEC",
       "operation": "AND"
      },
      {
       "name": "SERVICED",
       "operation": "AND"
      }
     ]
    }
   ]
  },
  {
   "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "name": "critical",
   "namespace": "test",
   "endpoint_group": "sites",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "roc.critical",
    "uuid": "5637d684-1f8e-4a02-a502-720e8f11e432"
   },
   "groups": [
    {
     "name": "compute",
     "operation": "OR",
     "services": [
      {
       "name": "CREAM-CE",
       "operation": "AND"
      },
      {
       "name": "ARC-CE",
       "operation": "AND"
      }
     ]
    },
    {
     "name": "storage",
     "operation": "OR",
     "services": [
      {
       "name": "SRMv2",
       "operation": "AND"
      },
      {
       "name": "SRM",
       "operation": "AND"
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

func (suite *OperationsProfilesTestSuite) NotTestListQueryName() {

	request, _ := http.NewRequest("GET", "/api/v2/aggregation_profiles?name=cloud", strings.NewReader(""))
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
   "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "name": "cloud",
   "namespace": "test",
   "endpoint_group": "sites",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "roc.critical",
    "uuid": "5637d684-1f8e-4a02-a502-720e8f11e432"
   },
   "groups": [
    {
     "name": "compute",
     "operation": "OR",
     "services": [
      {
       "name": "SERVICEA",
       "operation": "AND"
      },
      {
       "name": "SERVICEB",
       "operation": "AND"
      }
     ]
    },
    {
     "name": "images",
     "operation": "OR",
     "services": [
      {
       "name": "SERVICEC",
       "operation": "AND"
      },
      {
       "name": "SERVICED",
       "operation": "AND"
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

func (suite *OperationsProfilesTestSuite) NotTestListOneNotFound() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Not Found",
  "code": "404",
  "details": "item with the specific UUID was not found on the server"
 }
}`

	request, _ := http.NewRequest("GET", "/api/v2/aggregation_profiles/wrong-uuid", strings.NewReader(jsonInput))
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

func (suite *OperationsProfilesTestSuite) NotTestListOne() {

	request, _ := http.NewRequest("GET", "/api/v2/aggregation_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
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
   "name": "critical",
   "namespace": "test",
   "endpoint_group": "sites",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "roc.critical",
    "uuid": "5637d684-1f8e-4a02-a502-720e8f11e432"
   },
   "groups": [
    {
     "name": "compute",
     "operation": "OR",
     "services": [
      {
       "name": "CREAM-CE",
       "operation": "AND"
      },
      {
       "name": "ARC-CE",
       "operation": "AND"
      }
     ]
    },
    {
     "name": "storage",
     "operation": "OR",
     "services": [
      {
       "name": "SRMv2",
       "operation": "AND"
      },
      {
       "name": "SRM",
       "operation": "AND"
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

func (suite *OperationsProfilesTestSuite) NotTestCreateBadJson() {

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

	request, _ := http.NewRequest("POST", "/api/v2/aggregation_profiles", strings.NewReader(jsonInput))
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

func (suite *OperationsProfilesTestSuite) NotTestInvalidCreate() {

	jsonInput := `{
   "name": "yolo",
   "namespace": "testing-namespace",
   "endpoint_group": "test",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "uuid": "6ac7d684-1f8e-4a02-a502-720e8f110007"
   },
   "groups": [
    {
     "name": "tttcompute",
     "operation": "OR",
     "services": [
      {
       "name": "tttCREAM-CE",
       "operation": "AND"
      },
      {
       "name": "tttARC-CE",
       "operation": "AND"
      }
     ]
    },
    {
     "name": "tttstorage",
     "operation": "OR",
     "services": [
      {
       "name": "tttSRMv2",
       "operation": "AND"
      },
      {
       "name": "tttSRM",
       "operation": "AND"
      }
     ]
    }
   ]
  }`

	jsonOutput := `{
 "status": {
  "message": "Referenced metric profile UUID is not found",
  "code": "422"
 }
}`

	request, _ := http.NewRequest("POST", "/api/v2/aggregation_profiles", strings.NewReader(jsonInput))
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

func (suite *OperationsProfilesTestSuite) NotTestCreate() {

	jsonInput := `{
   "name": "yolo",
   "namespace": "testing-namespace",
   "endpoint_group": "test",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"
   },
   "groups": [
    {
     "name": "tttcompute",
     "operation": "OR",
     "services": [
      {
       "name": "tttCREAM-CE",
       "operation": "AND"
      },
      {
       "name": "tttARC-CE",
       "operation": "AND"
      }
     ]
    },
    {
     "name": "tttstorage",
     "operation": "OR",
     "services": [
      {
       "name": "tttSRMv2",
       "operation": "AND"
      },
      {
       "name": "tttSRM",
       "operation": "AND"
      }
     ]
    }
   ]
  }`

	jsonOutput := `{
 "status": {
  "message": "Aggregation Profile successfully created",
  "code": "201"
 },
 "data": {
  "uuid": "{{UUID}}",
  "links": {
   "self": "https:///api/v2/aggregation_profiles/{{UUID}}"
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
   "name": "yolo",
   "namespace": "testing-namespace",
   "endpoint_group": "test",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "critical",
    "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"
   },
   "groups": [
    {
     "name": "tttcompute",
     "operation": "OR",
     "services": [
      {
       "name": "tttCREAM-CE",
       "operation": "AND"
      },
      {
       "name": "tttARC-CE",
       "operation": "AND"
      }
     ]
    },
    {
     "name": "tttstorage",
     "operation": "OR",
     "services": [
      {
       "name": "tttSRMv2",
       "operation": "AND"
      },
      {
       "name": "tttSRM",
       "operation": "AND"
      }
     ]
    }
   ]
  }
 ]
}`

	request, _ := http.NewRequest("POST", "/api/v2/aggregation_profiles", strings.NewReader(jsonInput))
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

	c.Find(bson.M{"name": "yolo"}).One(&result)
	uuid := result["uuid"].(string)

	// Apply uuid to output template and check
	suite.Equal(strings.Replace(jsonOutput, "{{UUID}}", uuid, 2), output, "Response body mismatch")

	// Check that actually the item has been created
	// Call List one with the specific UUID
	request2, _ := http.NewRequest("GET", "/api/v2/aggregation_profiles/"+uuid, strings.NewReader(jsonInput))
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

func (suite *OperationsProfilesTestSuite) NotTestUpdateBadJson() {

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

	request, _ := http.NewRequest("PUT", "/api/v2/aggregation_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
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

func (suite *OperationsProfilesTestSuite) NotTestUpdateNotFound() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Not Found",
  "code": "404",
  "details": "item with the specific UUID was not found on the server"
 }
}`

	request, _ := http.NewRequest("PUT", "/api/v2/aggregation_profiles/wrong-uuid", strings.NewReader(jsonInput))
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

func (suite *OperationsProfilesTestSuite) NotTestInvalidUpdate() {

	jsonInput := `{
   "name": "yolo",
   "namespace": "testing-namespace",
   "endpoint_group": "test",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "testing",
    "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e007"
   },
   "groups": [
    {
     "name": "tttcompute",
     "operation": "OR",
     "services": [
      {
       "name": "tttCREAM-CE",
       "operation": "AND"
      },
      {
       "name": "tttARC-CE",
       "operation": "AND"
      }
     ]
    },
    {
     "name": "tttstorage",
     "operation": "OR",
     "services": [
      {
       "name": "tttSRMv2",
       "operation": "AND"
      },
      {
       "name": "tttSRM",
       "operation": "AND"
      }
     ]
    }
   ]
  }`

	jsonOutput := `{
 "status": {
  "message": "Referenced metric profile UUID is not found",
  "code": "422"
 }
}`

	request, _ := http.NewRequest("PUT", "/api/v2/aggregation_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
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

func (suite *OperationsProfilesTestSuite) NotTestUpdate() {

	jsonInput := `{
   "name": "yolo",
   "namespace": "testing-namespace",
   "endpoint_group": "test",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "testing",
    "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"
   },
   "groups": [
    {
     "name": "tttcompute",
     "operation": "OR",
     "services": [
      {
       "name": "tttCREAM-CE",
       "operation": "AND"
      },
      {
       "name": "tttARC-CE",
       "operation": "AND"
      }
     ]
    },
    {
     "name": "tttstorage",
     "operation": "OR",
     "services": [
      {
       "name": "tttSRMv2",
       "operation": "AND"
      },
      {
       "name": "tttSRM",
       "operation": "AND"
      }
     ]
    }
   ]
  }`

	jsonOutput := `{
 "status": {
  "message": "Aggregation Profile successfully updated",
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
   "name": "yolo",
   "namespace": "testing-namespace",
   "endpoint_group": "test",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "critical",
    "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"
   },
   "groups": [
    {
     "name": "tttcompute",
     "operation": "OR",
     "services": [
      {
       "name": "tttCREAM-CE",
       "operation": "AND"
      },
      {
       "name": "tttARC-CE",
       "operation": "AND"
      }
     ]
    },
    {
     "name": "tttstorage",
     "operation": "OR",
     "services": [
      {
       "name": "tttSRMv2",
       "operation": "AND"
      },
      {
       "name": "tttSRM",
       "operation": "AND"
      }
     ]
    }
   ]
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/aggregation_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
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
	request2, _ := http.NewRequest("GET", "/api/v2/aggregation_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
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

func (suite *OperationsProfilesTestSuite) NotTestDeleteNotFound() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Not Found",
  "code": "404",
  "details": "item with the specific UUID was not found on the server"
 }
}`

	request, _ := http.NewRequest("DELETE", "/api/v2/aggregation_profiles/wrong-uuid", strings.NewReader(jsonInput))
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
