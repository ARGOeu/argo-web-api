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

package statusEndpoints

import (
	"net/http"
	"strings"
	"testing"

	"code.google.com/p/gcfg"
	"github.com/ARGOeu/argo-web-api/utils/authentication"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// This is a util. suite struct used in tests (see pkg "testify")
type StatusEndpointsTestSuite struct {
	suite.Suite
	cfg              config.Config
	tenantDbConf     config.MongoConfig
}

// Setup the Test Environment
// This function runs before any test and setups the environment
// A test configuration object is instantiated using a reference
// to testdb: argo_test_jobs. Also here is are instantiated some expected
// xml response validation messages (authorization,crud responses).
// Also the testdb is seeded with two jobs
func (suite *StatusEndpointsTestSuite) SetupTest() {

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
    db = "argo_core_test_endpoints"
`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	// Connect to mongo testdb
	session, _ := mongo.OpenSession(suite.cfg.MongoDB)

	// Add authentication token to mongo testdb
	seedAuth := bson.M{"api_key": "S3CR3T"}
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "authentication", seedAuth)

	// seed mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// seed a tenant to use
	c := session.DB(suite.cfg.MongoDB.Db).C("tenants")
	c.Insert(bson.M{
		"name": "AVENGERS",
		"db_conf": []bson.M{
			bson.M{
				"store":    "ar",
				"server":   "localhost",
				"port":     27017,
				"database": "argo_tenant1_endpoints_db1",
				"username": "admin",
				"password": "3NCRYPT3D"},
			bson.M{
				"store":    "status",
				"server":   "b.mongodb.org",
				"port":     27017,
				"database": "status_db",
				"username": "admin",
				"password": "3NCRYPT3D"},
		},
		"users": []bson.M{
			bson.M{
				"name":    "cap",
				"email":   "cap@email.com",
				"api_key": "C4PK3Y"},
			bson.M{
				"name":    "thor",
				"email":   "thor@email.com",
				"api_key": "TH0RK3Y"},
		}})
		c.Insert(bson.M{
			"name": "THEOTHERS",
			"db_conf": []bson.M{
				bson.M{
					"store":    "ar",
					"server":   "localhost",
					"port":     27017,
					"database": "argo_tenant2_endpoints_db1",
					"username": "admin",
					"password": "UN3NCRYPT3D"},
				bson.M{
					"store":    "status",
					"server":   "c.mongodb.org",
					"port":     27017,
					"database": "status_db",
					"username": "admin",
					"password": "UN3NCRYPT3D"},
			},
			"users": []bson.M{
				bson.M{
					"name":    "chapie",
					"email":   "chapie@email.com",
					"api_key": "C3POK3Y"},
				bson.M{
					"name":    "saul",
					"email":   "saul@email.com",
					"api_key": "SAULK3Y"},
			}})

	// get dbconfiguration based on the tenant
	// Prepare the request object
	request, _ := http.NewRequest("GET", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")
	// authenticate user's api key and find corresponding tenant
	suite.tenantDbConf, err = authentication.AuthenticateTenant(request.Header, suite.cfg)

	// Now seed metric data
	c = session.DB(suite.tenantDbConf.Db).C("status_endpoints")
	c.Insert(bson.M{
		"job":                 "JOB_A",
		"date_int":            20150501,
		"timestamp":           "2015-05-01T00:00:00Z",
		"supergroup":          "NGI_GRNET",
		"endpoint_group":      "HG-03-AUTH",
		"group_type":          "NGI",
		"endpoint_group_type": "SITES",
		"service":             "CREAM-CE",
		"hostname":            "cream.afroditi.hellasgrid.gr",
		"status":              "OK",
		"time_int":            0,
	})
	c.Insert(bson.M{
		"job":                 "JOB_A",
		"date_int":            20150501,
		"timestamp":           "2015-05-01T01:00:00Z",
		"supergroup":          "NGI_GRNET",
		"endpoint_group":      "HG-03-AUTH",
		"group_type":          "NGI",
		"endpoint_group_type": "SITES",
		"service":             "CREAM-CE",
		"hostname":            "cream.afroditi.hellasgrid.gr",
		"status":              "CRITICAL",
		"time_int":            10000,
	})
	c.Insert(bson.M{
		"job":                 "JOB_A",
		"date_int":            20150501,
		"timestamp":           "2015-05-01T05:00:00Z",
		"supergroup":          "NGI_GRNET",
		"endpoint_group":      "HG-03-AUTH",
		"group_type":          "NGI",
		"endpoint_group_type": "SITES",
		"service":             "CREAM-CE",
		"hostname":            "cream.afroditi.hellasgrid.gr",
		"status":              "OK",
		"time_int":            50000,
	})
	c.Insert(bson.M{
		"job":                 "JOB_A",
		"date_int":            20150501,
		"timestamp":           "2015-05-01T06:00:00Z",
		"supergroup":          "NGI_GRNET",
		"endpoint_group":      "HG-03-AUTH",
		"group_type":          "NGI",
		"endpoint_group_type": "SITES",
		"service":             "Site-BDII",
		"hostname":            "cream.afroditi.hellasgrid.gr",
		"status":              "OK",
		"time_int":            60000,
	})
	
	// add now anothee authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C3POK3Y")
	// authenticate user's api key and find corresponding tenant
	suite.tenantDbConf, err = authentication.AuthenticateTenant(request.Header, suite.cfg)
	c = session.DB(suite.tenantDbConf.Db).C("status_endpoints")
	c.Insert(bson.M{
		"job":                 "JOB_A",
		"date_int":            20150501,
		"timestamp":           "2015-05-01T00:00:10Z",
		"supergroup":          "NORTH",
		"endpoint_group":      "JUELICH",
		"group_type":          "REGION",
		"endpoint_group_type": "SITES",
		"service":             "iRods",
		"hostname":            "irods01.juelich.de",
		"status":              "OK",
		"time_int":            10,
	})
	c.Insert(bson.M{
		"job":                 "JOB_A",
		"date_int":            20150501,
		"timestamp":           "2015-05-01T01:01:00Z",
		"supergroup":          "NORTH",
		"endpoint_group":      "JUELICH",
		"group_type":          "REGION",
		"endpoint_group_type": "SITES",
		"service":             "iRods",
		"hostname":            "irods01.juelich.de",
		"status":              "CRITICAL",
		"time_int":            10100,
	})
	c.Insert(bson.M{
		"job":                 "JOB_A",
		"date_int":            20150511,
		"timestamp":           "2015-05-11T01:01:00Z",
		"supergroup":          "NORTH",
		"endpoint_group":      "JUELICH",
		"group_type":          "REGION",
		"endpoint_group_type": "SITES",
		"service":             "iRods",
		"hostname":            "irods01.juelich.de",
		"status":              "CRITICAL",
		"time_int":            10100,
	})
}

func (suite *StatusEndpointsTestSuite) TestListStatusEndpoints() {

	respXML1 := ` <root>
   <job name="JOB_A">
     <endpoint hostname="cream.afroditi.hellasgrid.gr" service="CREAM-CE">
       <status timestamp="2015-05-01T00:00:00Z" status="OK"></status>
       <status timestamp="2015-05-01T01:00:00Z" status="CRITICAL"></status>
       <status timestamp="2015-05-01T05:00:00Z" status="OK"></status>
     </endpoint>
   </job>
 </root>`

	respXML2 := ` <root>
   <job name="JOB_A">
     <endpoint hostname="irods01.juelich.de" service="iRods">
       <status timestamp="2015-05-01T00:00:10Z" status="OK"></status>
       <status timestamp="2015-05-01T01:01:00Z" status="CRITICAL"></status>
     </endpoint>
   </job>
 </root>`

	fullurl1 := "/api/v1/status/endpoints/timeline/cream.afroditi.hellasgrid.gr/CREAM-CE?" +
		"job=JOB_A&start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:59:59Z"

	fullurl2 := "/api/v1/status/endpoints/timeline/irods01.juelich.de/iRods?" +
		"job=JOB_A&start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:59:59Z"

	// Prepare the request object for fist tenant
	request, _ := http.NewRequest("GET", fullurl1, strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")
	// Pass request to controller calling List() handler method
	code, _, output, _ := List(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respXML1, string(output), "Response body mismatch")

	// Prepare the request object for second tenant
	request, _ = http.NewRequest("GET", fullurl2, strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C3POK3Y")
	// Pass request to controller calling List() handler method
	code, _, output, _ = List(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respXML2, string(output), "Response body mismatch")

	// Prepare the request object for accessing a tenant without proper authorization
	request, _ = http.NewRequest("GET", fullurl1, strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "WRONGKEY")
	// Pass request to controller calling List() handler method
	code, _, _ , _ = List(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(401, code, "Should have gotten return code 401 (Unauthorized)")
}


// This function is actually called in the end of all tests
// and clears the test environment.
// Mainly it's purpose is to drop the testdb
func (suite *StatusEndpointsTestSuite) TearDownTest() {

	session, _ := mongo.OpenSession(suite.cfg.MongoDB)

	session.DB("argo_core_test_endpoints").DropDatabase()
	session.DB("argo_tenant1_endpoints_db1").DropDatabase()
	session.DB("argo_tenant2_endpoints_db1").DropDatabase()
}

// This is the first function called when go test is issued
func TestJobsSuite(t *testing.T) {
	suite.Run(t, new(StatusEndpointsTestSuite))
}