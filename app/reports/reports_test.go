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

package reports

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
type ReportTestSuite struct {
	suite.Suite
	cfg                config.Config
	tenantDbConf       config.MongoConfig
	respReportCreated  string
	respReportUpdated  string
	respReportDeleted  string
	respReportNotFound string
	respUnauthorized   string
	respBadJSON        string
}

// Setup the Test Environment
// This function runs before any test and setups the environment
// A test configuration object is instantiated using a reference
// to testdb: argo_test_reports. Also here is are instantiated some expected
// xml response validation messages (authorization,crud responses).
// Also the testdb is seeded with two reports
func (suite *ReportTestSuite) SetupTest() {

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
    db = "argo_test_reports"
`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.respReportCreated = " <root>\n" +
		"   <Message>Report was successfully created</Message>\n </root>"

	suite.respReportUpdated = " <root>\n" +
		"   <Message>Report was successfully updated</Message>\n </root>"

	suite.respReportDeleted = " <root>\n" +
		"   <Message>Report was successfully deleted</Message>\n </root>"

	suite.respReportNotFound = " <root>\n" +
		"   <Message>Report not found</Message>\n </root>"

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
				"database": "argo_test_reports_db1",
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

	// Now seed the report DEFINITIONS
	c = session.DB(suite.tenantDbConf.Db).C("reports")
	c.Insert(bson.M{
		"name":            "Report_A",
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
		"name":            "Report_B",
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
}

// TestCreateReport function implements testing the http POST create report request.
// Request requires admin authentication and gets as input a json body containing
// all the available information to be added to the datastore
// After the operation succeeds is double-checked
// that the newly created report is correctly retrieved
func (suite *ReportTestSuite) TestCreateReport() {

	// create json input data for the request
	postData := `
  {
    "name":"Foo_Report",
    "tenant":"AVENGERS",
    "profiles":[
      { "name":"metric","value":"profA"},
      { "name":"ap","value":"profB"}
     ],
    "endpoint_group":"SITES",
    "group_of_groups":"NGI",
    "filter_tags":[
      { "name":"production","value":"Y"},
      { "name":"monitored","value":"Y"}
    ]
  }
    `
	// Prepare the request object
	request, _ := http.NewRequest("POST", "", strings.NewReader(postData))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	// Execute the request in the controller
	code, _, output, _ := Create(request, suite.cfg)

	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.respReportCreated, string(output), "Response body mismatch")

	// Double check that you read the newly inserted profile
	// Create a string literal of the expected xml Response
	respXML := `<root>
 <report name="Foo_Report" tenant="AVENGERS" endpoint_group="SITES" group_of_groups="NGI">
  <profiles>
   <profile name="metric" value="profA"></profile>
   <profile name="ap" value="profB"></profile>
  </profiles>
  <filter_tags>
   <tag name="production" value="Y"></tag>
   <tag name="monitored" value="Y"></tag>
  </filter_tags>
 </report>
</root>`

	// Prepare the request object using report name as urlvar in url path
	request, _ = http.NewRequest("GET", "/api/v1/reports/Foo_Report", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")
	// Pass request to controller calling List() handler method
	code, _, output, _ = ListOne(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respXML, string(output), "Response body mismatch")
}

// TestUpdateReport function implements testing the http PUT update report request.
// Request requires admin authentication and gets as input the name of the
// report to be updated and a json body with the update.
// After the operation succeeds is double-checked
// that the specific report has been updated
func (suite *ReportTestSuite) TestUpdateReport() {

	// create json input data for the request
	postData := `
  {
    "name":"Report_A_modified",
    "tenant":"AVENGERS",
    "profiles":[
      { "name":"metric","value":"profA_mod"},
      { "name":"ap","value":"profB_mod"}
     ],
    "endpoint_group":"SITES",
    "group_of_groups":"NGI",
    "filter_tags":[
      { "name":"production","value":"Y"},
      { "name":"monitored","value":"Y"}
    ]
  }
    `
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "/api/v1/reports/Report_A", strings.NewReader(postData))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	// Execute the request in the controller
	code, _, output, _ := Update(request, suite.cfg)

	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.respReportUpdated, string(output), "Response body mismatch")

	// Double check that you read the newly inserted profile
	// Create a string literal of the expected xml Response
	respXML := `<root>
 <report name="Report_A_modified" tenant="AVENGERS" endpoint_group="SITES" group_of_groups="NGI">
  <profiles>
   <profile name="metric" value="profA_mod"></profile>
   <profile name="ap" value="profB_mod"></profile>
  </profiles>
  <filter_tags>
   <tag name="production" value="Y"></tag>
   <tag name="monitored" value="Y"></tag>
  </filter_tags>
 </report>
</root>`

	// Prepare the request object using report name as urlvar in url path
	request, _ = http.NewRequest("GET", "/api/v1/reports/Report_A_modified", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")
	// Pass request to controller calling List() handler method
	code, _, output, _ = ListOne(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respXML, string(output), "Response body mismatch")
}

// TestDeleteReport function implements testing the http DELETE report request.
// Request requires admin authentication and gets as input the name of the
// report to be deleted. After the operation succeeds is double-checked
// that the deleted report is actually missing from the datastore
func (suite *ReportTestSuite) TestDeleteReport() {

	// Prepare the request object
	request, _ := http.NewRequest("DELETE", "/api/v1/reports/Report_B", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	// Execute the request in the controller
	code, _, output, _ := Delete(request, suite.cfg)

	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.respReportDeleted, string(output), "Response body mismatch")

	// Double check that the report is actually removed when you try
	// to retrieve it's information by name
	// Prepare the request object using report name as urlvar in url path
	request, _ = http.NewRequest("GET", "/api/v1/reports/Report_B", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")
	// Pass request to controller calling List() handler method
	code, _, output, _ = ListOne(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(400, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(suite.respReportNotFound, string(output), "Response body mismatch")
}

// TestReadOneReport function implements the testing
// of the get request which retrieves information
// about a specific report (using it's name as input)
func (suite *ReportTestSuite) TestReadOneReport() {

	// Create a string literal of the expected xml Response
	respXML := `<root>
 <report name="Report_A" tenant="AVENGERS" endpoint_group="SITES" group_of_groups="NGI">
  <profiles>
   <profile name="metric" value="profile1"></profile>
   <profile name="ops" value="profile2"></profile>
  </profiles>
  <filter_tags>
   <tag name="name1" value="value1"></tag>
   <tag name="name2" value="value2"></tag>
  </filter_tags>
 </report>
</root>`

	// Prepare the request object using report name as urlvar in url path
	request, _ := http.NewRequest("GET", "/api/v1/reports/Report_A", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")
	// Pass request to controller calling List() handler method
	code, _, output, _ := ListOne(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")

	// Compare the expected and actual xml response
	suite.Equal(respXML, string(output), "Response body mismatch")
}

// TestReadReport function implements the testing
// of the get request which retrieves information
// about all available reports
func (suite *ReportTestSuite) TestReadReports() {

	// Create a string literal of the expected xml Response
	respXML := `<root>
 <report name="Report_A" tenant="AVENGERS" endpoint_group="SITES" group_of_groups="NGI">
  <profiles>
   <profile name="metric" value="profile1"></profile>
   <profile name="ops" value="profile2"></profile>
  </profiles>
  <filter_tags>
   <tag name="name1" value="value1"></tag>
   <tag name="name2" value="value2"></tag>
  </filter_tags>
 </report>
 <report name="Report_B" tenant="AVENGERS" endpoint_group="SITES" group_of_groups="NGI">
  <profiles>
   <profile name="metric" value="profile1"></profile>
   <profile name="ops" value="profile2"></profile>
  </profiles>
  <filter_tags>
   <tag name="name1" value="value1"></tag>
   <tag name="name2" value="value2"></tag>
  </filter_tags>
 </report>
</root>`

	// Prepare the request object
	request, _ := http.NewRequest("GET", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")
	// Pass request to controller calling List() handler method
	code, _, output, _ := List(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respXML, string(output), "Response body mismatch")
}

// TestCreateUnauthorized function tests calling the create report request (POST) and
// providing a wrong api-key. The response should be unauthorized
func (suite *ReportTestSuite) TestCreateUnauthorized() {
	// Prepare the request object (use id2 for path)
	request, _ := http.NewRequest("POST", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "F00T0K3N")

	// Execute the request in the controller
	code, _, output, _ := Create(request, suite.cfg)

	suite.Equal(401, code, "Internal Server Error")
	suite.Equal(suite.respUnauthorized, string(output), "Response body mismatch")
}

// TestUpdateUnauthorized function tests calling the update report request (PUT)
// and providing  a wrong api-key. The response should be unauthorized
func (suite *ReportTestSuite) TestUpdateUnauthorized() {
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "", strings.NewReader("{}"))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "F00T0K3N")

	// Execute the request in the controller
	code, _, output, _ := Update(request, suite.cfg)

	suite.Equal(401, code, "Internal Server Error")
	suite.Equal(suite.respUnauthorized, string(output), "Response body mismatch")
}

// TestDeleteUnauthorized function tests calling the remove report request (DELETE)
// and providing a wrong api-key. The response should be unauthorized
func (suite *ReportTestSuite) TestDeleteUnauthorized() {
	// Prepare the request object
	request, _ := http.NewRequest("DELETE", "", strings.NewReader("{}"))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "F00T0K3N")

	// Execute the request in the controller
	code, _, output, _ := Delete(request, suite.cfg)

	suite.Equal(401, code, "Internal Server Error")
	suite.Equal(suite.respUnauthorized, string(output), "Response body mismatch")
}

// TestCreateBadJson tests calling the create report request (POST) and providing
// bad json input. The response should be malformed json
func (suite *ReportTestSuite) TestCreateBadJson() {
	// Prepare the request object
	request, _ := http.NewRequest("POST", "/api/v1/reports/Report_A", strings.NewReader("{ bad json"))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	// Execute the request in the controller
	code, _, output, _ := Create(request, suite.cfg)

	suite.Equal(400, code, "Internal Server Error")
	suite.Equal(suite.respBadJSON, string(output), "Response body mismatch")
}

// TestUpdateBadJson tests calling the update report request (PUT) and providing
// bad json input. The response should be malformed json
func (suite *ReportTestSuite) TestUpdateBadJson() {
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "/api/v1/reports/Re[prt_A", strings.NewReader("{ bad json"))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	// Execute the request in the controller
	code, _, output, _ := Update(request, suite.cfg)

	suite.Equal(400, code, "Internal Server Error")
	suite.Equal(suite.respBadJSON, string(output), "Response body mismatch")
}

// TestListOneNotFound tests calling the http (GET) report info request
// and provide a non-existing report name. The response should be report not found
func (suite *ReportTestSuite) TestListOneNotFound() {
	// Prepare the request object
	request, _ := http.NewRequest("GET", "/api/v1/reports/BADNAME", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	// Execute the request in the controller
	code, _, output, _ := ListOne(request, suite.cfg)

	suite.Equal(400, code, "Internal Server Error")
	suite.Equal(suite.respReportNotFound, string(output), "Response body mismatch")
}

// TestUpdateNotFound tests calling the http (PUT) update report equest
// and provide a non-existing report name. The response should be report not found
func (suite *ReportTestSuite) TestUpdateNotFound() {
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "/api/v1/reports/BADNAME", strings.NewReader("{}"))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	// Execute the request in the controller
	code, _, output, _ := Update(request, suite.cfg)

	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.respReportNotFound, string(output), "Response body mismatch")
}

// TestDeleteNotFound tests calling the http (PUT) update report request
// and provide a non-existing report name. The response should be report not found
func (suite *ReportTestSuite) TestDeleteNotFound() {
	// Prepare the request object
	request, _ := http.NewRequest("DELETE", "/api/v1/reports/BADNAME", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	// Execute the request in the controller
	code, _, output, _ := Delete(request, suite.cfg)

	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.respReportNotFound, string(output), "Response body mismatch")
}

// This function is actually called in the end of all tests
// and clears the test environment.
// Mainly it's purpose is to drop the testdb
func (suite *ReportTestSuite) TearDownTest() {

	session, _ := mongo.OpenSession(suite.cfg.MongoDB)

	session.DB("argo_test_reports").DropDatabase()
	session.DB("argo_test_reports_db1").DropDatabase()
}

// This is the first function called when go test is issued
func TestReportSuite(t *testing.T) {
	suite.Run(t, new(ReportTestSuite))
}
