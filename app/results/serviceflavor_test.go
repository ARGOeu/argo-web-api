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

package results

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

type serviceFlavorAvailabilityTestSuite struct {
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
// This function runs before any test and setups the environment
func (suite *serviceFlavorAvailabilityTestSuite) SetupTest() {

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
    db = "ARGO_test_serviceFlavor_availability"
    `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.tenantDbConf.Db = "ARGO_test_serviceFlavor_availability_tenant"
	suite.tenantDbConf.Password = "h4shp4ss"
	suite.tenantDbConf.Username = "dbuser"
	suite.tenantDbConf.Store = "ar"
	suite.clientkey = "secretkey"

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(true).PathPrefix("/api/v2/results").Subrouter()
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
				},
				bson.M{
					"name":    "King Joffrey",
					"email":   "g0dk1ng@kingslanding.gov",
					"api_key": "sansa <3",
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
					"database": "argo_wrong_db_serviceflavoravailability",
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "Joe Complex",
					"email":   "C.Joe@egi.eu",
					"api_key": suite.clientkey,
				},
				bson.M{
					"name":    "Josh Plain",
					"email":   "P.Josh@egi.eu",
					"api_key": "itsamysterytoyou",
				},
			}})
	// Seed database with recomputations
	c = session.DB(suite.tenantDbConf.Db).C("service_ar")

	// Insert seed data
	c.Insert(
		bson.M{
			"report":       "Report_A",
			"date":         20150622,
			"name":         "SF01",
			"supergroup":   "ST01",
			"up":           0.98264,
			"down":         0,
			"unknown":      0,
			"availability": 98.26389,
			"reliability":  98.26389,
			"tags": []bson.M{
				bson.M{
					"name":  "production",
					"value": "Y",
				},
			},
		},
		bson.M{
			"report":       "Report_A",
			"date":         20150622,
			"name":         "SF02",
			"supergroup":   "ST01",
			"up":           0.96875,
			"down":         0,
			"unknown":      0,
			"availability": 96.875,
			"reliability":  96.875,
			"tags": []bson.M{
				bson.M{
					"name":  "production",
					"value": "Y",
				},
			},
		},
		bson.M{
			"report":       "Report_A",
			"date":         20150622,
			"name":         "SF03",
			"supergroup":   "ST02",
			"up":           0.96875,
			"down":         0,
			"unknown":      0,
			"availability": 96.875,
			"reliability":  96.875,
			"tags": []bson.M{
				bson.M{
					"name":  "production",
					"value": "Y",
				},
			},
		},
		bson.M{
			"report":       "Report_A",
			"date":         20150623,
			"name":         "SF01",
			"supergroup":   "ST01",
			"up":           0.53472,
			"down":         0.33333,
			"unknown":      0.01042,
			"availability": 54.03509,
			"reliability":  81.48148,
			"tags": []bson.M{
				bson.M{
					"name":  "production",
					"value": "Y",
				},
			},
		},
		bson.M{
			"report":       "Report_A",
			"date":         20150623,
			"name":         "SF02",
			"supergroup":   "ST01",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 100,
			"reliability":  100,
			"tags": []bson.M{
				bson.M{
					"name":  "production",
					"value": "Y",
				},
			},
		})

	c = session.DB(suite.tenantDbConf.Db).C("reports")

	c.Insert(bson.M{
		"info": bson.M{
			"name":        "Report_A",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "GROUP",
				"group": bson.M{
					"type": "SITE",
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
}

// TestListServiceFlavorAvailabilityMonthly tests if daily results are returned correctly
func (suite *serviceFlavorAvailabilityTestSuite) TestListServiceFlavorAvailabilityMonthly() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/services?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	responseBody := response.Body.String()
	serviceFlavorAvailabilityXML := ` <root>
   <group name="ST01" type="SITE">
     <group name="SF01" type="service">
       <results timestamp="2015-06" availability="76.26534166743393" reliability="91.61418757296076" unknown="0.00521" uptime="0.75868" downtime="0.166665"></results>
     </group>
     <group name="SF02" type="service">
       <results timestamp="2015-06" availability="98.43749901562502" reliability="98.43749901562502" unknown="0" uptime="0.984375" downtime="0"></results>
     </group>
   </group>
 </root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(serviceFlavorAvailabilityXML, responseBody, "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/services?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	responseBody = response.Body.String()

	serviceFlavorAvailabilityJSON := `{
   "root": [
     {
       "name": "ST01",
       "type": "SITE",
       "serviceflavors": [
         {
           "name": "SF01",
           "type": "service",
           "results": [
             {
               "timestamp": "2015-06",
               "availability": "76.26534166743393",
               "reliability": "91.61418757296076",
               "unknown": "0.00521",
               "uptime": "0.75868",
               "downtime": "0.166665"
             }
           ]
         },
         {
           "name": "SF02",
           "type": "service",
           "results": [
             {
               "timestamp": "2015-06",
               "availability": "98.43749901562502",
               "reliability": "98.43749901562502",
               "unknown": "0",
               "uptime": "0.984375",
               "downtime": "0"
             }
           ]
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(serviceFlavorAvailabilityJSON, responseBody, "Response body mismatch")

}

// TestListServiceFlavorAvailabilityDaily tests if daily results are returned correctly
func (suite *serviceFlavorAvailabilityTestSuite) TestListServiceFlavorAvailabilityDaily() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/services?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	serviceFlavorAvailabilityXML := ` <root>
   <group name="ST01" type="SITE">
     <group name="SF01" type="service">
       <results timestamp="2015-06-22" availability="98.26389" reliability="98.26389" unknown="0" uptime="0.98264" downtime="0"></results>
       <results timestamp="2015-06-23" availability="54.03509" reliability="81.48148" unknown="0.01042" uptime="0.53472" downtime="0.33333"></results>
     </group>
     <group name="SF02" type="service">
       <results timestamp="2015-06-22" availability="96.875" reliability="96.875" unknown="0" uptime="0.96875" downtime="0"></results>
       <results timestamp="2015-06-23" availability="100" reliability="100" unknown="0" uptime="1" downtime="0"></results>
     </group>
   </group>
 </root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(serviceFlavorAvailabilityXML, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/services?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	serviceFlavorAvailabilityJSON := `{
   "root": [
     {
       "name": "ST01",
       "type": "SITE",
       "serviceflavors": [
         {
           "name": "SF01",
           "type": "service",
           "results": [
             {
               "timestamp": "2015-06-22",
               "availability": "98.26389",
               "reliability": "98.26389",
               "unknown": "0",
               "uptime": "0.98264",
               "downtime": "0"
             },
             {
               "timestamp": "2015-06-23",
               "availability": "54.03509",
               "reliability": "81.48148",
               "unknown": "0.01042",
               "uptime": "0.53472",
               "downtime": "0.33333"
             }
           ]
         },
         {
           "name": "SF02",
           "type": "service",
           "results": [
             {
               "timestamp": "2015-06-22",
               "availability": "96.875",
               "reliability": "96.875",
               "unknown": "0",
               "uptime": "0.96875",
               "downtime": "0"
             },
             {
               "timestamp": "2015-06-23",
               "availability": "100",
               "reliability": "100",
               "unknown": "0",
               "uptime": "1",
               "downtime": "0"
             }
           ]
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(serviceFlavorAvailabilityJSON, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01/services?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", "AWRONGKEY")
	request.Header.Set("Accept", "application/xml")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 401 Unauthorized code
	suite.Equal(401, response.Code, "Incorrect HTTP response code")

}

// TestListServiceFlavorAvailabilityErrors tests if errors/exceptions are returned correctly
func (suite *serviceFlavorAvailabilityTestSuite) TestListServiceFlavorAvailabilityErrors() {

	reportErrorXML := ` <root>
   <message>The report with the name Report_B does not exist</message>
 </root>`

	typeErrorXML := ` <root>
   <message>The report Report_A does not define endpoint group type: Site. Try using SITE instead.</message>
 </root>`

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_B/SITE/ST01/services?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 400 bad request code
	suite.Equal(400, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(reportErrorXML, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/Site/ST01/services?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	//request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 400 bad request code
	suite.Equal(400, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(typeErrorXML, response.Body.String(), "Response body mismatch")

}

//TearDownTest to tear down every test
func (suite *serviceFlavorAvailabilityTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	defer session.Close()

	if err != nil {
		panic(err)
	}

	cols, err := session.DB(suite.tenantDbConf.Db).CollectionNames()
	for _, col := range cols {
		session.DB(suite.tenantDbConf.Db).C(col).RemoveAll(bson.M{})
	}
	cols, err = session.DB(suite.cfg.MongoDB.Db).CollectionNames()
	for _, col := range cols {
		session.DB(suite.cfg.MongoDB.Db).C(col).RemoveAll(bson.M{})
	}

}

//TearDownTest to tear down every test
func (suite *serviceFlavorAvailabilityTestSuite) TearDownSuite() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	defer session.Close()

	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()

}

// TestServiceFlavorAvailabilityTestSuite is responsible for calling the tests
func TestServiceFlavorAvailabilityTestSuite(t *testing.T) {
	suite.Run(t, new(serviceFlavorAvailabilityTestSuite))
}
