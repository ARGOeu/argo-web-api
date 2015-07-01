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

package endpointGroupAvailability

import (
	"net/http"
	"strings"
	"testing"

	"code.google.com/p/gcfg"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type endpointGroupAvailabilityTestSuite struct {
	suite.Suite
	cfg                       config.Config
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
			"job":          "Job_A",
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
		})
	c.Insert(
		bson.M{
			"job":          "Job_A",
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
		})
	c.Insert(
		bson.M{
			"job":          "Job_A",
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
		})
	c.Insert(
		bson.M{
			"job":          "Job_A",
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
	c.Insert(
		bson.M{
			"job":          "Job_A",
			"date":         20150623,
			"name":         "VO01",
			"type":         "VO",
			"supergroup":   "GROUP_C",
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
	c.Insert(
		bson.M{
			"job":          "Job_A",
			"date":         20150623,
			"name":         "VO02",
			"type":         "VO",
			"supergroup":   "GROUP_C",
			"uptime":       1,
			"downtime":     0,
			"unknown":      0,
			"availability": 43.5,
			"reliability":  56,
			"weights":      4987,
			"tags": []bson.M{
				bson.M{
					"name":  "",
					"value": "",
				},
			},
		})
	c.Insert(
		bson.M{
			"job":          "Job_A",
			"date":         20150623,
			"name":         "VO03",
			"type":         "VO",
			"supergroup":   "GROUP_C",
			"uptime":       0,
			"downtime":     1,
			"unknown":      0,
			"availability": 73.5,
			"reliability":  59,
			"weights":      2345,
			"tags": []bson.M{
				bson.M{
					"name":  "",
					"value": "",
				},
			},
		})

}

// TestListEndpointGroupAvailability test if daily results are returned correctly
func (suite *endpointGroupAvailabilityTestSuite) TestListEndpointGroupTypeSiteAvailability() {

	request, _ := http.NewRequest("GET", "/api/v1/endpoint_group_availability?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&job=Job_A&granularity=daily&type=SITE", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	code, _, output, _ := List(request, suite.cfg)

	endpointGrouAvailabitiyXML := ` <root>
   <Job name="Job_A">
     <EndpointGroup name="ST01" SuperGroup="GROUP_A" type="SITE">
       <Availability timestamp="2015-06-22" availability="66.7" reliability="54.6"></Availability>
       <Availability timestamp="2015-06-23" availability="100" reliability="100"></Availability>
     </EndpointGroup>
     <EndpointGroup name="ST02" SuperGroup="GROUP_A" type="SITE">
       <Availability timestamp="2015-06-22" availability="70" reliability="45"></Availability>
       <Availability timestamp="2015-06-23" availability="43.5" reliability="56"></Availability>
     </EndpointGroup>
   </Job>
 </root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(endpointGrouAvailabitiyXML, string(output), "Response body mismatch")

}

func (suite *endpointGroupAvailabilityTestSuite) TestListEndpointGroupTypeVOAvailability() {

	request, _ := http.NewRequest("GET", "/api/v1/endpoint_group_availability?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&job=Job_A&granularity=daily&type=VO", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	code, _, output, _ := List(request, suite.cfg)

	endpointGrouAvailabitiyXML := ` <root>
   <Job name="Job_A">
     <EndpointGroup name="VO01" SuperGroup="GROUP_C" type="VO">
       <Availability timestamp="2015-06-23" availability="43.5" reliability="56"></Availability>
     </EndpointGroup>
     <EndpointGroup name="VO02" SuperGroup="GROUP_C" type="VO">
       <Availability timestamp="2015-06-23" availability="43.5" reliability="56"></Availability>
     </EndpointGroup>
     <EndpointGroup name="VO03" SuperGroup="GROUP_C" type="VO">
       <Availability timestamp="2015-06-23" availability="73.5" reliability="59"></Availability>
     </EndpointGroup>
   </Job>
 </root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(endpointGrouAvailabitiyXML, string(output), "Response body mismatch")

}

func (suite *endpointGroupAvailabilityTestSuite) TestListEndpointGroupTypeVOMonthlyAvailability() {

	request, _ := http.NewRequest("GET", "/api/v1/endpoint_group_availability?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&job=Job_A&granularity=monthly&type=VO", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	code, _, output, _ := List(request, suite.cfg)

	endpointGrouAvailabitiyXML := ` <root>
   <Job name="Job_A">
     <EndpointGroup name="VO01" SuperGroup="GROUP_C" type="VO">
       <Availability timestamp="2015-06" availability="99.99999900000002" reliability="99.99999900000002"></Availability>
     </EndpointGroup>
     <EndpointGroup name="VO02" SuperGroup="GROUP_C" type="VO">
       <Availability timestamp="2015-06" availability="99.99999900000002" reliability="99.99999900000002"></Availability>
     </EndpointGroup>
     <EndpointGroup name="VO03" SuperGroup="GROUP_C" type="VO">
       <Availability timestamp="2015-06" availability="0" reliability="0"></Availability>
     </EndpointGroup>
   </Job>
 </root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(endpointGrouAvailabitiyXML, string(output), "Response body mismatch")

}

//TearDownTest to tear down every test
func (suite *endpointGroupAvailabilityTestSuite) TearDownTest() {

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
