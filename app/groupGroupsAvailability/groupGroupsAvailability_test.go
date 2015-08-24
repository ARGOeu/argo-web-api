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

package groupGroupsAvailability

import (
	"code.google.com/p/gcfg"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"testing"
)

type EgAvailabilityTestSuite struct {
	suite.Suite
	cfg              config.Config
	tenantDbConf     config.MongoConfig
	tenantPassword   string
	tenantUsername   string
	tenantStorename  string
	clientKey        string
	respUnauthorized string
	responseDaily    string
	responseMonthly  string
	dailyApiCall     string
	monthlyApiCall   string
}

// Setup the Test Environment
// This function runs before any test and setups the environment
// A test configuration object is instantiated using a reference
// to testdb: ARGO_core_test_* .

func (suite *EgAvailabilityTestSuite) SetupTest() {

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
    db = "ARGO_core_test_egavailability"
`

	suite.respUnauthorized = "Unauthorized"
	suite.tenantDbConf.Db = "ARGO_northern_test_egavailability"
	suite.tenantDbConf.Password = "h4shp4ss"
	suite.tenantDbConf.Username = "dbuser"
	suite.tenantDbConf.Store = "ar"
	suite.clientKey = "secretkey"
	suite.responseDaily = ` <root>
   <Job name="Job_A">
     <SuperGroup name="GROUP_A">
       <Availability timestamp="2015-06-22" availability="68.13896116893515" reliability="50.413931144915935"></Availability>
       <Availability timestamp="2015-06-23" availability="75.36324059247399" reliability="80.8138510808647"></Availability>
     </SuperGroup>
   </Job>
 </root>`

	suite.responseMonthly = ` <root>
   <Job name="Job_A">
     <SuperGroup name="GROUP_A">
       <Availability timestamp="2015-06" availability="71.75110088070457" reliability="65.61389111289031"></Availability>
     </SuperGroup>
   </Job>
 </root>`

	suite.dailyApiCall = "/api/v1/group_groups_availability?start_time=2015-06-21T00:00:00Z&end_time=2015-06-24T00:00:00Z&job=Job_A&granularity=daily"
	suite.monthlyApiCall = "/api/v1/group_groups_availability?start_time=2015-06-21T00:00:00Z&end_time=2015-06-24T00:00:00Z&job=Job_A&granularity=monthly"

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	// Connect to mongo testdb
	session, err := mongo.OpenSession(suite.cfg.MongoDB)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Seed database with tenants and test credentials
	c := session.DB(suite.cfg.MongoDB.Db).C("tenants")
	c.Insert(
		bson.M{"name": "Western",
			"db_conf": []bson.M{

				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Western1",
				},
				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Western2",
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "John Snow",
					"email":   "J.Snow@foo.bar",
					"api_key": "wh1t3_w@lk3rs",
				},
				bson.M{
					"name":    "King Joffrey",
					"email":   "g0dk1ng@foo.bar",
					"api_key": "sansa <3",
				},
			}})
	c.Insert(
		bson.M{"name": "Northern",
			"db_conf": []bson.M{

				bson.M{
					// "store":    "ar",
					"server":   "localhost",
					"port":     27017,
					"database": suite.tenantDbConf.Db,
					"username": suite.tenantDbConf.Username,
					"password": suite.tenantDbConf.Password,
				},
				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_northern_test_db",
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "Joe Complex",
					"email":   "C.Joe@foo.bar",
					"api_key": suite.clientKey,
				},
				bson.M{
					"name":    "Josh Plain",
					"email":   "P.Josh@foo.bar",
					"api_key": "testsecretapikey",
				},
			}})

	// Open DB session
	c = session.DB(suite.tenantDbConf.Db).C("endpoint_group_ar")

	// Insert seed data
	c.Insert(
		bson.M{"job": "Job_A",
			"date":         20150622,
			"name":         "ST01",
			"type":         "SITE",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
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
		})
	c.Insert(
		bson.M{"job": "Job_A",
			"date":         20150622,
			"name":         "ST02",
			"type":         "SITE",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
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
		})
	c.Insert(
		bson.M{"job": "Job_A",
			"date":         20150623,
			"name":         "ST01",
			"type":         "SITE",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
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
		})
	c.Insert(
		bson.M{"job": "Job_A",
			"date":         20150623,
			"name":         "ST02",
			"type":         "SITE",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
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
}

// Check if the Daily calculation returns correct data.
func (suite *EgAvailabilityTestSuite) TestReadGroupArDaily() {

	request, _ := http.NewRequest("GET", suite.dailyApiCall, strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientKey)
	code, _, output, _ := List(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the actual with the expected results
	suite.Equal(suite.responseDaily, string(output), "Response body mismatch")
}

// Check if the Monthly calculation returns correct data.
func (suite *EgAvailabilityTestSuite) TestReadGroupArMonthly() {

	request, _ := http.NewRequest("GET", suite.monthlyApiCall, strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientKey)
	code, _, output, _ := List(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the actual with the expected results
	suite.Equal(suite.responseMonthly, string(output), "Response body mismatch")
}

// Use the API call for Daily results to check if the authentication mechanism works as expected.
func (suite *EgAvailabilityTestSuite) TestReadDailyUnauthorized() {

	request, _ := http.NewRequest("GET", suite.dailyApiCall, strings.NewReader(""))
	request.Header.Set("x-api-key", "clientKey")
	code, _, _, err := List(request, suite.cfg)
	// Check that we must have a 200 ok code.
	suite.Equal(401, code, "Internal Server Error")
	// Check if the error message is the expected one.
	suite.Equal(suite.respUnauthorized, err.Error(), "Response body mismatch")
}

// Use the API call for Monthly results to check if the authentication mechanism works as expected.
func (suite *EgAvailabilityTestSuite) TestReadMonthlyUnauthorized() {

	request, _ := http.NewRequest("GET", suite.monthlyApiCall, strings.NewReader(""))
	request.Header.Set("x-api-key", "clientKey")
	code, _, _, err := List(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(401, code, "Internal Server Error")
	// Check if the error message is the expected one.
	suite.Equal(suite.respUnauthorized, err.Error(), "Response body mismatch")
}

// This function is actually called in the end of all tests
// and clears the test environment.
// Mainly it's purpose is to drop the testdb
func (suite *EgAvailabilityTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

// This is the first function called when go test is issued
func TestEgAvailabilityTestSuite(t *testing.T) {
	suite.Run(t, new(EgAvailabilityTestSuite))
}
