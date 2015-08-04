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

	"code.google.com/p/gcfg"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type endpointGroupAvailabilityTestSuite struct {
	suite.Suite
	cfg                       config.Config
	router                    *mux.Router
	confHandler               respond.ConfHandler
	tenantDbConf              config.MongoConfig
	tenantpassword            string
	tenantusername            string
	tenantstorename           string
	clientkey                 string
	respRecomputationsCreated string
}

// Setup the Test Environment
// This function runs before any test and setups the environment
func (suite *endpointGroupAvailabilityTestSuite) SetupTest() {

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
    db = "ARGO_test_endpointGroup_availability"
    `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.tenantDbConf.Db = "ARGO_test_endpointGroup_availability_tenant"
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
					"database": "argo_wrong_db_endpointgrouavailability",
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
	c = session.DB(suite.tenantDbConf.Db).C("endpoint_group_ar")

	// Insert seed data
	c.Insert(
		bson.M{
			"report":       "Report_A",
			"date":         20150622,
			"name":         "ST01",
			"type":         "SITE",
			"supergroup":   "GROUP_A",
			"uptime":       1,
			"downtime":     0,
			"unknown":      0,
			"availability": 66.7,
			"reliability":  54.6,
			"weights":      5634,
			"tags": []bson.M{
				bson.M{
					"name":  "",
					"value": "",
				},
			},
		},
		bson.M{
			"report":       "Report_A",
			"date":         20150622,
			"name":         "ST02",
			"type":         "SITE",
			"supergroup":   "GROUP_A",
			"uptime":       1,
			"downtime":     0,
			"unknown":      0,
			"availability": 70,
			"reliability":  45,
			"weights":      4356,
			"tags": []bson.M{
				bson.M{
					"name":  "",
					"value": "",
				},
			},
		},
		bson.M{
			"report":       "Report_A",
			"date":         20150623,
			"name":         "ST01",
			"type":         "SITE",
			"supergroup":   "GROUP_A",
			"uptime":       1,
			"downtime":     0,
			"unknown":      0,
			"availability": 100,
			"reliability":  100,
			"weights":      5634,
			"tags": []bson.M{
				bson.M{
					"name":  "",
					"value": "",
				},
			},
		},
		bson.M{
			"report":       "Report_A",
			"date":         20150623,
			"name":         "ST02",
			"type":         "SITE",
			"supergroup":   "GROUP_A",
			"uptime":       1,
			"downtime":     0,
			"unknown":      0,
			"availability": 43.5,
			"reliability":  56,
			"weights":      4356,
			"tags": []bson.M{
				bson.M{
					"name":  "",
					"value": "",
				},
			},
		})

	c = session.DB(suite.tenantDbConf.Db).C("reports")

	c.Insert(
		bson.M{
			"name":   "Report_A",
			"tenant": "EGI",
			"profiles": bson.M{
				"availability": "ap1",
				"metrics":      "ch.cern.sam.ROC_CRITICAL",
				"operations":   "ops1",
			},
			"endpoints_group": "SITE",
			"group_of_groups": "GROUP",
			"filter_tags": []bson.M{
				{"name": "production", "value": "Y"},
				{"name": "monitored", "value": "Y"},
			},
		})

}

// TestListEndpointGroupAvailability test if daily results are returned correctly
func (suite *endpointGroupAvailabilityTestSuite) TestListEndpointGroupAvailability() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	endpointGrouAvailabitiyXML := ` <root>
   <group name="GROUP_A" type="GROUP">
     <group name="ST01" type="SITE">
       <results timestamp="2015-06-22" availability="66.7" reliability="54.6" unknown="0" uptime="1" downtime="0"></results>
       <results timestamp="2015-06-23" availability="100" reliability="100" unknown="0" uptime="1" downtime="0"></results>
     </group>
   </group>
 </root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointGrouAvailabitiyXML, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/SITE/ST01?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	endpointGrouAvailabitiyJSON := `{
   "root": [
     {
       "name": "GROUP_A",
       "type": "GROUP",
       "endpoints": [
         {
           "name": "ST01",
           "type": "SITE",
           "results": [
             {
               "timestamp": "2015-06-22",
               "availability": "66.7",
               "reliability": "54.6",
               "unknown": "0",
               "uptime": "1",
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
	suite.Equal(endpointGrouAvailabitiyJSON, response.Body.String(), "Response body mismatch")

}

// TestListAllEndpointGroupAvailability test if daily results are returned correctly
func (suite *endpointGroupAvailabilityTestSuite) TestListAllEndpointGroupAvailability() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/SITE?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	endpointGrouAvailabitiyXML := ` <root>
   <group name="GROUP_A" type="GROUP">
     <group name="ST01" type="SITE">
       <results timestamp="2015-06-22" availability="66.7" reliability="54.6" unknown="0" uptime="1" downtime="0"></results>
       <results timestamp="2015-06-23" availability="100" reliability="100" unknown="0" uptime="1" downtime="0"></results>
     </group>
     <group name="ST02" type="SITE">
       <results timestamp="2015-06-22" availability="70" reliability="45" unknown="0" uptime="1" downtime="0"></results>
       <results timestamp="2015-06-23" availability="43.5" reliability="56" unknown="0" uptime="1" downtime="0"></results>
     </group>
   </group>
 </root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointGrouAvailabitiyXML, response.Body.String(), "Response body mismatch")

}

//TearDownTest to tear down every test
func (suite *endpointGroupAvailabilityTestSuite) TearDownTest() {

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
func (suite *endpointGroupAvailabilityTestSuite) TearDownSuite() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	defer session.Close()

	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()

}

// TestRecompuptationsTestSuite is responsible for calling the tests
func TestRecompuptationsTestSuite(t *testing.T) {
	suite.Run(t, new(endpointGroupAvailabilityTestSuite))
}
