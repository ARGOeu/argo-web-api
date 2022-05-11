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

package status

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

type StatusTestSuite struct {
	suite.Suite
	cfg             config.Config
	router          *mux.Router
	confHandler     respond.ConfHandler
	tenantDbConf    config.MongoConfig
	tenantpassword  string
	tenantusername  string
	tenantstorename string
	clientkey       string
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
		  db = "ARGO_test_statusV3"
		  `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.tenantDbConf.Db = "ARGO_test_statusV3_tenant"
	suite.tenantDbConf.Password = "h4shp4ss"
	suite.tenantDbConf.Username = "dbuser"
	suite.tenantDbConf.Store = "ar"
	suite.clientkey = "secretkey"

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v3/status").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)
}

// This function runs before any test and setups the environment
func (suite *StatusTestSuite) SetupTest() {

	log.SetOutput(ioutil.Discard)

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
		bson.M{"name": "Westeros",
			"db_conf": []bson.M{

				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Westeros1",
				},
				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Westeros2",
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "John Snow",
					"email":   "J.Snow@brothers.wall",
					"api_key": "wh1t3_w@lk3rs",
					"roles":   []string{"viewer"},
				},
				bson.M{
					"name":    "King Joffrey",
					"email":   "g0dk1ng@kingslanding.gov",
					"api_key": "sansa <3",
					"roles":   []string{"viewer"},
				},
			}})
	c.Insert(
		bson.M{"name": "EGI",
			"db_conf": []bson.M{

				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": suite.tenantDbConf.Db,
					"username": suite.tenantDbConf.Username,
					"password": suite.tenantDbConf.Password,
				},
				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_wrong_db_endpointgrouavailability",
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "Joe Complex",
					"email":   "C.Joe@egi.eu",
					"api_key": suite.clientkey,
					"roles":   []string{"viewer"},
				},
				bson.M{
					"name":    "Josh Plain",
					"email":   "P.Josh@egi.eu",
					"api_key": "itsamysterytoyou",
					"roles":   []string{"viewer"},
				},
			}})

	c = session.DB(suite.cfg.MongoDB.Db).C("roles")
	c.Insert(
		bson.M{
			"resource": "v3.status.list",
			"roles":    []string{"editor", "viewer"},
		})

	c = session.DB(suite.tenantDbConf.Db).C("reports")

	c.Insert(bson.M{
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
			bson.M{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			bson.M{
				"name":  "name1",
				"value": "value1"},
			bson.M{
				"name":  "name2",
				"value": "value2"},
		}})

	// Seed tenant database with data
	c = session.DB(suite.tenantDbConf.Db).C(statusGroupColName)

	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "SITEA",
		"status":         "OK",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T01:00:00Z",
		"endpoint_group": "SITEA",
		"status":         "CRITICAL",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T05:00:00Z",
		"endpoint_group": "SITEA",
		"status":         "OK",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "SITEB",
		"status":         "OK",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T03:00:00Z",
		"endpoint_group": "SITEB",
		"status":         "WARNING",
	})
	c.Insert(bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T17:53:00Z",
		"endpoint_group":     "SITEB",
		"status":             "CRITICAL",
		"has_threshold_rule": true,
	})

	// seed the endpoints
	c = session.DB(suite.tenantDbConf.Db).C(statusEndpointColName)
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"info":           bson.M{"Url": "http://example.foo/path/to/service"},
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "SITEA",
		"service":        "CREAM-CE",
		"host":           "cream01.example.foo",

		"status": "OK",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T01:00:00Z",
		"endpoint_group": "SITEA",
		"service":        "CREAM-CE",
		"host":           "cream01.example.foo",
		"info":           bson.M{"Url": "http://example.foo/path/to/service"},

		"status": "CRITICAL",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T05:00:00Z",
		"endpoint_group": "SITEA",
		"service":        "CREAM-CE",
		"host":           "cream01.example.foo",
		"info":           bson.M{"Url": "http://example.foo/path/to/service"},

		"status": "OK",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "SITEA",
		"service":        "CREAM-CE",
		"host":           "cream02.example.foo",

		"status": "OK",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T08:47:00Z",
		"endpoint_group": "SITEA",
		"service":        "CREAM-CE",
		"host":           "cream02.example.foo",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "WARNING",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T12:00:00Z",
		"endpoint_group": "SITEA",
		"service":        "CREAM-CE",
		"host":           "cream02.example.foo",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "OK",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "SITEB",
		"service":        "CREAM-CE",
		"host":           "cream03.example.foo",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "OK",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T04:40:00Z",
		"endpoint_group": "SITEB",
		"service":        "CREAM-CE",
		"host":           "cream03.example.foo",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "UNKNOWN",
	})
	c.Insert(bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T06:00:00Z",
		"endpoint_group":     "SITEB",
		"service":            "CREAM-CE",
		"host":               "cream03.example.foo",
		"metric":             "emi.cream.CREAMCE-JobSubmit",
		"status":             "CRITICAL",
		"has_threshold_rule": true,
	})

}

// TestListStatus test the status results
func (suite *StatusTestSuite) TestListStatus() {

	request, _ := http.NewRequest("GET", "/api/v3/status/Report_A?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:59:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	expResponse := `{
 "groups": [
  {
   "name": "SITEA",
   "type": "SITES",
   "statuses": [
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
   ],
   "endpoints": [
    {
     "hostname": "cream01.example.foo",
     "service": "CREAM-CE",
     "info": {
      "Url": "http://example.foo/path/to/service"
     },
     "statuses": [
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
     "hostname": "cream02.example.foo",
     "service": "CREAM-CE",
     "statuses": [
      {
       "timestamp": "2015-05-01T00:00:00Z",
       "value": "OK"
      },
      {
       "timestamp": "2015-05-01T08:47:00Z",
       "value": "WARNING"
      },
      {
       "timestamp": "2015-05-01T12:00:00Z",
       "value": "OK"
      },
      {
       "timestamp": "2015-05-01T23:59:59Z",
       "value": "OK"
      }
     ]
    }
   ]
  },
  {
   "name": "SITEB",
   "type": "SITES",
   "statuses": [
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
   ],
   "endpoints": [
    {
     "hostname": "cream03.example.foo",
     "service": "CREAM-CE",
     "statuses": [
      {
       "timestamp": "2015-05-01T00:00:00Z",
       "value": "OK"
      },
      {
       "timestamp": "2015-05-01T04:40:00Z",
       "value": "UNKNOWN"
      },
      {
       "timestamp": "2015-05-01T06:00:00Z",
       "value": "CRITICAL"
      },
      {
       "timestamp": "2015-05-01T23:59:59Z",
       "value": "CRITICAL"
      }
     ]
    }
   ]
  }
 ]
}`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual json response
	suite.Equal(expResponse, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v3/status/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z", strings.NewReader(""))
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

}

// TestOptions tests responses in case the OPTIONS http verb is used
func (suite *StatusTestSuite) TestOptions() {

	request, _ := http.NewRequest("OPTIONS", "/api/v3/status/Report_A", strings.NewReader(""))

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

//TearDownTest to tear down every test
func (suite *StatusTestSuite) TearDownTest() {

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
func (suite *StatusTestSuite) TearDownSuite() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

// TestEndpointGroupsTestSuite is responsible for calling the tests
func TestStatusTestSuite(t *testing.T) {
	suite.Run(t, new(StatusTestSuite))
}
