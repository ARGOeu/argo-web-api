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

package aggregationProfiles

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
type AggregationProfilesTestSuite struct {
	suite.Suite
	cfg                       config.Config
	router                    *mux.Router
	confHandler               respond.ConfHandler
	tenantDbConf              config.MongoConfig
	clientkey                 string
	respRecomputationsCreated string
	respUnauthorized          string
}

// SetupSuite Setup the Test Environment
func (suite *AggregationProfilesTestSuite) SetupSuite() {
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
    db = "AR_test_aggr_prof"
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

	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v2").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)
}

// This function runs before any test and setups the environment
func (suite *AggregationProfilesTestSuite) SetupTest() {

	log.SetOutput(ioutil.Discard)

	// seed mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Seed database with tenants
	//TODO: move tests to

	//seed roles
	// Seed database with metric profiles
	c := session.DB(suite.cfg.MongoDB.Db).C("roles")
	c.Insert(
		bson.M{
			"resource": "aggregationProfiles.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "aggregationProfiles.get",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "aggregationProfiles.create",
			"roles":    []string{"editor"},
		})
	c.Insert(
		bson.M{
			"resource": "aggregationProfiles.delete",
			"roles":    []string{"editor"},
		})
	c.Insert(
		bson.M{
			"resource": "aggregationProfiles.update",
			"roles":    []string{"editor"},
		})

	c = session.DB(suite.cfg.MongoDB.Db).C("tenants")
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
	// Seed database with metric profiles
	c = session.DB(suite.tenantDbConf.Db).C("aggregation_profiles")
	c.EnsureIndexKey("-date_integer", "id")
	c.Insert(
		bson.M{
			"id":                "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			"date_integer":      20191004,
			"date":              "2019-10-04",
			"name":              "critical",
			"namespace":         "test",
			"endpoint_group":    "sites",
			"metric_operation":  "AND",
			"profile_operation": "AND",
			"metric_profile": bson.M{
				"name": "roc.critical",
				"id":   "5637d684-1f8e-4a02-a502-720e8f11e432",
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
			"id":                "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			"date_integer":      20191104,
			"date":              "2019-11-04",
			"name":              "critical",
			"namespace":         "test",
			"endpoint_group":    "sites",
			"metric_operation":  "AND",
			"profile_operation": "AND",
			"metric_profile": bson.M{
				"name": "roc.critical",
				"id":   "5637d684-1f8e-4a02-a502-720e8f11e432",
			},
			"groups": []bson.M{
				bson.M{"name": "compute2",
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
				bson.M{"name": "storage2",
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
			"id":                "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"date_integer":      20190404,
			"date":              "2019-04-04",
			"name":              "cloud",
			"namespace":         "test",
			"endpoint_group":    "sites",
			"metric_operation":  "AND",
			"profile_operation": "AND",
			"metric_profile": bson.M{
				"name": "roc.critical",
				"id":   "5637d684-1f8e-4a02-a502-720e8f11e432",
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
	c.Insert(
		bson.M{
			"id":                "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"date_integer":      20190504,
			"date":              "2019-05-04",
			"name":              "cloud",
			"namespace":         "test",
			"endpoint_group":    "sites",
			"metric_operation":  "AND",
			"profile_operation": "AND",
			"metric_profile": bson.M{
				"name": "roc.critical",
				"id":   "5637d684-1f8e-4a02-a502-720e8f11e432",
			},
			"groups": []bson.M{
				bson.M{"name": "compute2",
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
				bson.M{"name": "images2",
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

	// Seed database with metric profiles
	c = session.DB(suite.tenantDbConf.Db).C("metric_profiles")
	c.Insert(
		bson.M{
			"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			"name": "ch.cern.SAM.ROC_CRITICAL",
			"services": []bson.M{
				bson.M{"service": "CREAM-CE",
					"metrics": []string{
						"emi.cream.CREAMCE-JobSubmit",
						"emi.wn.WN-Bi",
						"emi.wn.WN-Csh",
						"emi.wn.WN-SoftVer"},
				},
				bson.M{"service": "SRMv2",
					"metrics": []string{"hr.srce.SRM2-CertLifetime",
						"org.sam.SRM-Del",
						"org.sam.SRM-Get",
						"org.sam.SRM-GetSURLs",
						"org.sam.SRM-GetTURLs",
						"org.sam.SRM-Ls",
						"org.sam.SRM-LsDir",
						"org.sam.SRM-Put"},
				},
			},
		})

}

func (suite *AggregationProfilesTestSuite) TestList() {

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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-11-04",
   "name": "critical",
   "namespace": "test",
   "endpoint_group": "sites",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "roc.critical",
    "id": "5637d684-1f8e-4a02-a502-720e8f11e432"
   },
   "groups": [
    {
     "name": "compute2",
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
     "name": "storage2",
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
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "date": "2019-05-04",
   "name": "cloud",
   "namespace": "test",
   "endpoint_group": "sites",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "roc.critical",
    "id": "5637d684-1f8e-4a02-a502-720e8f11e432"
   },
   "groups": [
    {
     "name": "compute2",
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
     "name": "images2",
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

func (suite *AggregationProfilesTestSuite) TestBadDate() {

	badDate := `{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "Bad Request",
   "code": "400",
   "details": "date parameter value: 2020-02 is not in the valid form of YYYY-MM-DD"
  }
 ]
}`

	type reqHeader struct {
		Method string
		Path   string
		Data   string
	}

	requests := []reqHeader{
		reqHeader{Method: "GET", Path: "/api/v2/aggregation_profiles?date=2020-02", Data: ""},
		reqHeader{Method: "GET", Path: "/api/v2/aggregation_profiles/some-uuid?date=2020-02", Data: ""},
		reqHeader{Method: "POST", Path: "/api/v2/aggregation_profiles?date=2020-02", Data: ""},
		reqHeader{Method: "PUT", Path: "/api/v2/aggregation_profiles/some-id?date=2020-02", Data: ""},
	}

	for _, r := range requests {
		request, _ := http.NewRequest(r.Method, r.Path, strings.NewReader(r.Data))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		code := response.Code
		output := response.Body.String()

		// Check that we must have a 200 ok code
		suite.Equal(400, code, "Internal Server Error")
		// Compare the expected and actual json response
		suite.Equal(badDate, output, "Response body mismatch")

	}

}

func (suite *AggregationProfilesTestSuite) TestListEmpty() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	defer session.Close()
	if err != nil {
		panic(err)
	}

	c := session.DB(suite.tenantDbConf.Db).C("aggregation_profiles")
	c.DropCollection()

	request, _ := http.NewRequest("GET", "/api/v2/aggregation_profiles", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	emptyList := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": []
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(emptyList, output, "Response body mismatch")

}

func (suite *AggregationProfilesTestSuite) TestListQueryName() {

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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "date": "2019-05-04",
   "name": "cloud",
   "namespace": "test",
   "endpoint_group": "sites",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "roc.critical",
    "id": "5637d684-1f8e-4a02-a502-720e8f11e432"
   },
   "groups": [
    {
     "name": "compute2",
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
     "name": "images2",
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

func (suite *AggregationProfilesTestSuite) TestListOneNotFound() {

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

	request, _ := http.NewRequest("GET", "/api/v2/aggregation_profiles/wrong-id", strings.NewReader(jsonInput))
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

func (suite *AggregationProfilesTestSuite) TestListOne() {

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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-11-04",
   "name": "critical",
   "namespace": "test",
   "endpoint_group": "sites",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "roc.critical",
    "id": "5637d684-1f8e-4a02-a502-720e8f11e432"
   },
   "groups": [
    {
     "name": "compute2",
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
     "name": "storage2",
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

func (suite *AggregationProfilesTestSuite) TestCreateForbidViewer() {

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

	request, _ := http.NewRequest("POST", "/api/v2/aggregation_profiles", strings.NewReader(jsonInput))
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

func (suite *AggregationProfilesTestSuite) TestUpdateForbidViewer() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Access to the resource is Forbidden",
  "code": "403"
 }
}`

	request, _ := http.NewRequest("PUT", "/api/v2/aggregation_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e007", strings.NewReader(jsonInput))
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

func (suite *AggregationProfilesTestSuite) TestDeleteForbidViewer() {

	request, _ := http.NewRequest("DELETE", "/api/v2/aggregation_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
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

func (suite *AggregationProfilesTestSuite) TestCreateBadJson() {

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

func (suite *AggregationProfilesTestSuite) TestInvalidCreate() {

	jsonInput := `{
   "name": "yolo",
   "namespace": "testing-namespace",
   "endpoint_group": "test",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "id": "6ac7d684-1f8e-4a02-a502-720e8f110007"
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
  "message": "Unprocessable Entity",
  "code": "422"
 },
 "errors": [
  {
   "message": "Unprocessable Entity",
   "code": "422",
   "details": "Referenced metric profile ID is not found"
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
	suite.Equal(422, code, "Internal Server Error")

	// Apply id to output template and check
	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *AggregationProfilesTestSuite) TestCreate() {

	jsonInput := `{
   "name": "yolo",
   "namespace": "testing-namespace",
   "endpoint_group": "test",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"
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
  "id": "{{id}}",
  "links": {
   "self": "https:///api/v2/aggregation_profiles/{{id}}"
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
   "id": "{{id}}",
   "date": "2019-03-03",
   "name": "yolo",
   "namespace": "testing-namespace",
   "endpoint_group": "test",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "ch.cern.SAM.ROC_CRITICAL",
    "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"
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

	request, _ := http.NewRequest("POST", "/api/v2/aggregation_profiles?date=2019-03-03", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(201, code, "Internal Server Error")
	// Compare the expected and actual json response

	// Grab id from mongodb
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	defer session.Close()
	if err != nil {
		panic(err)
	}

	// Retrieve id from database
	var result map[string]interface{}
	c := session.DB(suite.tenantDbConf.Db).C("aggregation_profiles")

	c.Find(bson.M{"name": "yolo"}).One(&result)
	id := result["id"].(string)

	// Apply id to output template and check
	suite.Equal(strings.Replace(jsonOutput, "{{id}}", id, 2), output, "Response body mismatch")

	// Check that actually the item has been created
	// Call List one with the specific id
	request2, _ := http.NewRequest("GET", "/api/v2/aggregation_profiles/"+id, strings.NewReader(jsonInput))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)

	code2 := response2.Code
	output2 := response2.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(200, code2, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(strings.Replace(jsonCreated, "{{id}}", id, 1), output2, "Response body mismatch")
}

func (suite *AggregationProfilesTestSuite) TestCreateNameAlreadyExists() {

	jsonInput := `{
   "name": "critical",
   "namespace": "testing-namespace",
   "endpoint_group": "test",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"
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
  "message": "Conflict",
  "code": "409"
 },
 "errors": [
  {
   "message": "Conflict",
   "code": "409",
   "details": "Aggregation profile with the same name already exists"
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

	suite.Equal(409, code)
	suite.Equal(jsonOutput, output)
}

func (suite *AggregationProfilesTestSuite) TestUpdateNameAlreadyExists() {

	jsonInput := `{
   "name": "critical",
   "namespace": "testing-namespace",
   "endpoint_group": "test",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"
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
  "message": "Conflict",
  "code": "409"
 },
 "errors": [
  {
   "message": "Conflict",
   "code": "409",
   "details": "Aggregation profile with the same name already exists"
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

	suite.Equal(409, code)
	suite.Equal(jsonOutput, output)
}

func (suite *AggregationProfilesTestSuite) TestInvalidUpdate() {

	jsonInput := `{
   "name": "yolo",
   "namespace": "testing-namespace",
   "endpoint_group": "test",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "id": "6ac7d684-1f8e-4a02-a502-720e8f110007"
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
  "message": "Unprocessable Entity",
  "code": "422"
 },
 "errors": [
  {
   "message": "Unprocessable Entity",
   "code": "422",
   "details": "Referenced metric profile ID is not found"
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
	suite.Equal(422, code, "Internal Server Error")

	// Apply id to output template and check
	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *AggregationProfilesTestSuite) TestUpdateBadJson() {

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

func (suite *AggregationProfilesTestSuite) TestUpdateNotFound() {

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

	request, _ := http.NewRequest("PUT", "/api/v2/aggregation_profiles/wrong-id", strings.NewReader(jsonInput))
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

func (suite *AggregationProfilesTestSuite) NotTestInvalidUpdate() {

	jsonInput := `{
   "name": "yolo",
   "namespace": "testing-namespace",
   "endpoint_group": "test",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "testing",
    "id": "6ac7d684-1f8e-4a02-a502-720e8f11e007"
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
  "message": "Unprocessable Entity",
  "code": "422"
 },
 "errors": [
  {
   "message": "Unprocessable Entity",
   "code": "422",
   "details": "Referenced metric profile ID is not found"
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
	suite.Equal(422, code, "Internal Server Error")
	// Compare the expected and actual json response

	// Apply id to output template and check
	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *AggregationProfilesTestSuite) TestUpdate() {

	jsonInput := `{
   "name": "yolo",
   "namespace": "testing-namespace",
   "endpoint_group": "test",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "testing",
    "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"
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
  "message": "Aggregations Profile successfully updated (new history snapshot)",
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
   "date": "2019-12-12",
   "name": "yolo",
   "namespace": "testing-namespace",
   "endpoint_group": "test",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "ch.cern.SAM.ROC_CRITICAL",
    "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"
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

	request, _ := http.NewRequest("PUT", "/api/v2/aggregation_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c?date=2019-12-12", strings.NewReader(jsonInput))
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
	request2, _ := http.NewRequest("GET", "/api/v2/aggregation_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c?date=2019-12-12", strings.NewReader(jsonInput))
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

func (suite *AggregationProfilesTestSuite) TestDeleteNotFound() {

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

	request, _ := http.NewRequest("DELETE", "/api/v2/aggregation_profiles/wrong-id", strings.NewReader(jsonInput))
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

func (suite *AggregationProfilesTestSuite) TestDelete() {

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
	c := session.DB(suite.tenantDbConf.Db).C("aggregation_profiles")
	err = c.Find(bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"}).One(&result)

	suite.NotEqual(err, nil, "No not found error")
	suite.Equal(err.Error(), "not found", "No not found error")
}

func (suite *AggregationProfilesTestSuite) TestOptionsAggregationProfiles() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/aggregation_profiles", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, POST, DELETE, PUT, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v2/aggregation_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))

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

//TearDownTest to tear down every test
func (suite *AggregationProfilesTestSuite) TearDownTest() {

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
func (suite *AggregationProfilesTestSuite) TearDownSuite() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

func TestAggregationProfilesTestSuite(t *testing.T) {
	suite.Run(t, new(AggregationProfilesTestSuite))
}
