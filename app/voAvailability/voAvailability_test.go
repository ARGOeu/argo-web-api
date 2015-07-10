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

package voAvailability

import (
	"code.google.com/p/gcfg"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strings"
	"testing"
)

// voAvailabilityTestSuite is a utility suite struct used in tests
type voAvailabilityTestSuite struct {
	suite.Suite
	cfg                   config.Config
	tenantcfg             config.MongoConfig
	resp_nokeyprovided    string
	resp_unauthorized     string
	resp_fedcloud_monthly string
	resp_fedcloud_daily   string
	resp_critical_monthly string
	resp_critical_daily   string
}

// voAvailability will bootstrap and provide the testing environment
func (suite *voAvailabilityTestSuite) SetupTest() {

	const coreConfig = `
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
    db = "ar_test_vo_avail"
`
	_ = gcfg.ReadStringInto(&suite.cfg, coreConfig)
	suite.resp_nokeyprovided = "404 page not found"
	suite.resp_unauthorized = "Unauthorized"

	// Connect to mongo coredb
	session, err := mongo.OpenSession(suite.cfg.MongoDB)
	defer mongo.CloseSession(session)
	if err != nil {
		panic(err)
	}

	// Add authentication token to mongo coredb
	seed_auth := bson.M{"name" : "EGI",
		"db_conf" : []bson.M{ bson.M{ "store": "ar", "server" : "127.0.0.1", "port" : 27017, "database" : "argo_EGI_test_vo"} } ,
		"users" : []bson.M{ bson.M{"name" : "Jack Doe", "email" : "jack.doe@example.com", "api_key" : "elmL5K"} }}
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "tenants", seed_auth)

	// TODO: I don't like it here that I rewrite the test data. 
	// However, this is a test for voreports, not for AuthenticateTenant function. 
	suite.tenantcfg.Host = "127.0.0.1"
	suite.tenantcfg.Port = 27017
	suite.tenantcfg.Db = "argo_EGI_test_vo"

	// Add a few data in collection voreports
	c := session.DB(suite.tenantcfg.Db).C("voreports")
	// TODO: Modify test *key* entries to match new Data Architecture Specification
	c.Insert(bson.M{"dt" : 20150501, "v" : "ops", "p" : "ch.cern.sam.ROC_CRITICAL", "ap" : "test-ap1",     "a" : 100,    "r" : 100,    "up" : 0.99306, "u" : 0.00694, "d" : 0 })
	c.Insert(bson.M{"dt" : 20150501, "v" : "ops", "p" : "ch.cern.sam.ROC_CRITICAL", "ap" : "egi-fedcloud", "a" : 100,    "r" : 100,    "up" : 1,       "u" : 0,       "d" : 0 })
	c.Insert(bson.M{"dt" : 20150502, "v" : "ops", "p" : "ch.cern.sam.ROC_CRITICAL", "ap" : "test-ap1",     "a" : 100,    "r" : 100,    "up" : 0.99306, "u" : 0.00694, "d" : 0 })
	c.Insert(bson.M{"dt" : 20150502, "v" : "ops", "p" : "ch.cern.sam.ROC_CRITICAL", "ap" : "egi-fedcloud", "a" : 99.653, "r" : 99.653, "up" : 0.99653, "u" : 0,       "d" : 0 })

}

// TestListVoAvailability will run unit tests against the List function
func (suite *voAvailabilityTestSuite) TestListVoAvailability() {

	suite.resp_fedcloud_monthly = ` <root>
   <Profile name="egi-fedcloud">
     <Vo VO="ops">
       <Availability timestamp="2015-05" availability="99.82649900173502" reliability="99.82649900173502"></Availability>
     </Vo>
   </Profile>
 </root>`
	suite.resp_fedcloud_daily   = ` <root>
   <Profile name="egi-fedcloud">
     <Vo VO="ops">
       <Availability timestamp="2015-05-01" availability="100" reliability="100"></Availability>
       <Availability timestamp="2015-05-02" availability="99.653" reliability="99.653"></Availability>
     </Vo>
   </Profile>
 </root>`
	suite.resp_critical_monthly = ` <root>
   <Profile name="test-ap1">
     <Vo VO="ops">
       <Availability timestamp="2015-05" availability="99.99999899301152" reliability="99.99999899301152"></Availability>
     </Vo>
   </Profile>
 </root>`
	suite.resp_critical_daily   = ` <root>
   <Profile name="test-ap1">
     <Vo VO="ops">
       <Availability timestamp="2015-05-01" availability="100" reliability="100"></Availability>
       <Availability timestamp="2015-05-02" availability="100" reliability="100"></Availability>
     </Vo>
   </Profile>
 </root>`

	// Prepare the request object
	request, _ := http.NewRequest("GET", "/api/v1/group_availability?start_time=2015-05-01T00:00:00Z&end_time=2015-05-03T00:00:00Z&availability_profile=egi-fedcloud&group_type=vo&granularity=monthly", nil)
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "elmL5K")
	// Execute the request in the controller
	code, _, output, _ := List(request, suite.cfg)
	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.resp_fedcloud_monthly, string(output), "Response body mismatch")

	// Prepare the request object
	request, _ = http.NewRequest("GET", "/api/v1/group_availability?start_time=2015-05-01T00:00:00Z&end_time=2015-05-03T00:00:00Z&availability_profile=egi-fedcloud&group_type=vo&granularity=daily", nil)
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "elmL5K")
	// Execute the request in the controller
	code, _, output, _ = List(request, suite.cfg)
	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.resp_fedcloud_daily, string(output), "Response body mismatch")

	// Prepare the request object
	request, _ = http.NewRequest("GET", "/api/v1/group_availability?start_time=2015-05-01T00:00:00Z&end_time=2015-05-03T00:00:00Z&availability_profile=test-ap1&group_type=vo&granularity=monthly", nil)
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "elmL5K")
	// Execute the request in the controller
	code, _, output, _ = List(request, suite.cfg)
	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.resp_critical_monthly, string(output), "Response body mismatch")

	// Prepare the request object
	request, _ = http.NewRequest("GET", "/api/v1/group_availability?start_time=2015-05-01T00:00:00Z&end_time=2015-05-03T00:00:00Z&availability_profile=test-ap1&group_type=vo&granularity=daily", nil)
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "elmL5K")
	// Execute the request in the controller
	code, _, output, _ = List(request, suite.cfg)
	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.resp_critical_daily, string(output), "Response body mismatch")

	// Prepare new request object
	request, _ = http.NewRequest("GET", "/api/v1/group_availability", strings.NewReader(""))
	// add the authentication token which is not seeded in testdb
	request.Header.Set("x-api-key", "wrongkey")
	// Execute the request in the controller
	code, _, output, _ = List(request, suite.cfg)
	suite.Equal(401, code, "Should have gotten return code 401 (Unauthorized)")
	suite.Equal(suite.resp_unauthorized, string(output), "Should have gotten reply Unauthorized")

	// Remove the test data from core db not to contaminate other tests
	// Open session to core mongo
	session, err := mongo.OpenSession(suite.cfg.MongoDB)
	if err != nil {
		panic(err)
	}
	defer mongo.CloseSession(session)
	// Open collection authentication
	c := session.DB(suite.cfg.MongoDB.Db).C("authentication")
	// Remove the specific entries inserted during this test
	c.Remove(bson.M{"name": "John Doe"})

	// Remove the test data from tenant db not to contaminate other tests
	// Open session to tenant mongo
	session, err = mgo.Dial(suite.tenantcfg.Host)
	if err != nil {
		panic(err)
	}
	defer mongo.CloseSession(session)
	// Open collection authentication
	c = session.DB(suite.tenantcfg.Db).C("voreports")

	// TODO: change key also here
	// Remove the specific entries inserted during this test
	c.Remove(bson.M{"v" : "ops"})

}

//TearDownTest to tear down every test
func (suite *voAvailabilityTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()

	session, err = mgo.Dial(suite.tenantcfg.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantcfg.Db).DropDatabase()

}

func TestVoAvailabilityTestSuite(t *testing.T) {
	suite.Run(t, new(voAvailabilityTestSuite))
}
