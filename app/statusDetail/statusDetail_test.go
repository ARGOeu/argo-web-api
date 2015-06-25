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

package statusDetail

import (
	"net/http"
	"strings"
	"testing"

	"code.google.com/p/gcfg"
	"github.com/argoeu/argo-web-api/utils/authentication"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// This is a util. suite struct used in tests (see pkg "testify")
type StatusDetailTestSuite struct {
	suite.Suite
	cfg              config.Config
	tenantDbConf     config.MongoConfig
	respUnauthorized string
	respBadJSON      string
}

// Setup the Test Environment
// This function runs before any test and setups the environment
// A test configuration object is instantiated using a reference
// to testdb: argo_test_jobs. Also here is are instantiated some expected
// xml response validation messages (authorization,crud responses).
// Also the testdb is seeded with two jobs
func (suite *StatusDetailTestSuite) SetupTest() {

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
    db = "argotest_detail"
`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.respBadJSON = " <root>\n" +
		"   <Message>Malformated json input data</Message>\n </root>"

	suite.respUnauthorized = "Unauthorized"

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
				"database": "argotest_detail_db1",
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

	// get dbconfiguration based on the tenant
	// Prepare the request object
	request, _ := http.NewRequest("GET", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")
	// authenticate user's api key and find corresponding tenant
	suite.tenantDbConf, err = authentication.AuthenticateTenant(request.Header, suite.cfg)

	// Now seed the job DEFINITIONS
	c = session.DB(suite.tenantDbConf.Db).C("jobs")
	c.Insert(bson.M{
		"name":            "Job_A",
		"tenant":          "AVENGERS",
		"endpoint_group":  "SITES",
		"group_of_groups": "NGI",
		"profiles": []bson.M{
			bson.M{
				"name":  "metric",
				"value": "profile1"},
			bson.M{
				"name":  "ops",
				"value": "profile2"},
		},
		"filter_tags": []bson.M{
			bson.M{
				"name":  "name1",
				"value": "value1"},
			bson.M{
				"name":  "name2",
				"value": "value2"},
		}})

	c.Insert(bson.M{
		"name":            "Job_B",
		"tenant":          "AVENGERS",
		"endpoint_group":  "SITES",
		"group_of_groups": "NGI",
		"profiles": []bson.M{
			bson.M{
				"name":  "metric",
				"value": "profile1"},
			bson.M{
				"name":  "ops",
				"value": "profile2"},
		},
		"filter_tags": []bson.M{
			bson.M{
				"name":  "name1",
				"value": "value1"},
			bson.M{
				"name":  "name2",
				"value": "value2"},
		}})

	// Now seed metric data
	// Now seed the job DEFINITIONS
	c = session.DB(suite.tenantDbConf.Db).C("status_metric")
	c.Insert(bson.M{
		"job":                 "jobA",
		"timestamp":           "2015-05-01T00:00:00Z",
		"group":               "NGI_GRNET",
		"endpoint_group":      "HG-03-AUTH",
		"group_type":          "NGI",
		"endpoint_group_type": "SITE",
		"service":             "CREAM-CE",
		"hostname":            "cream01.afroditi.gr",
		"metric":              "metric.jobsumbit",
		"status":              "OK",
		"time_int":            "0",
		"prev_status":         "OK",
		"prev_timestamp":      "2015-04-30T23:59:00Z",
	})
	c.Insert(bson.M{
		"job":                 "jobA",
		"timestamp":           "2015-05-01T01:00:00Z",
		"group":               "NGI_GRNET",
		"endpoint_group":      "HG-03-AUTH",
		"group_type":          "NGI",
		"endpoint_group_type": "SITE",
		"service":             "CREAM-CE",
		"hostname":            "cream01.afroditi.gr",
		"metric":              "metric.jobsumbit",
		"status":              "CRITICAL",
		"time_int":            "10000",
		"prev_status":         "OK",
		"prev_timestamp":      "2015-05-01T00:00:00Z",
	})
	c.Insert(bson.M{
		"job":                 "jobA",
		"timestamp":           "2015-05-01T05:00:00Z",
		"group":               "NGI_GRNET",
		"endpoint_group":      "HG-03-AUTH",
		"group_type":          "NGI",
		"endpoint_group_type": "SITE",
		"service":             "CREAM-CE",
		"hostname":            "cream01.afroditi.gr",
		"metric":              "metric.jobsumbit",
		"status":              "OK",
		"time_int":            "50000",
		"prev_status":         "CRITICAL",
		"prev_timestamp":      "2015-05-01T01:00:00Z",
	})

}

func (suite *StatusDetailTestSuite) TestListStatusDetail() {

}

// This function is actually called in the end of all tests
// and clears the test environment.
// Mainly it's purpose is to drop the testdb
func (suite *StatusDetailTestSuite) TearDownTest() {

	session, _ := mongo.OpenSession(suite.cfg.MongoDB)

	session.DB("argotest_detail").DropDatabase()
	session.DB("argotest_detail_db1").DropDatabase()
}

// This is the first function called when go test is issued
func TestJobsSuite(t *testing.T) {
	suite.Run(t, new(StatusDetailTestSuite))
}
