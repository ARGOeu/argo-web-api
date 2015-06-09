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

package availabilityProfiles

import (
	"code.google.com/p/gcfg"
	"encoding/xml"
	"fmt"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// This is a util. suite struct used in tests (see pkg "testify")
type AvProfileTestSuite struct {
	suite.Suite
	cfg                config.Config
	respProfileCreated string
	respProfileUpdated string
	respProfileDeleted string
	respUnauthorized   string
	respNoID           string
	respBadJSON        string
}

// Test structs to represent the expected XML schema
type Grouptest struct {
	XMLName xml.Name
}
type Ortest struct {
	XMLName xml.Name    `xml:"OR"`
	Groups  []Grouptest `xml:"Group"`
}
type Andtest struct {
	XMLName xml.Name `xml:"AND"`
	Ors     []Ortest `xml:"OR"`
}
type Profiletest struct {
	XMLName          xml.Name `xml:"profile"`
	ID               string   `xml:"id,attr"`
	Name             string   `xml:"name,attr"`
	Namespace        string   `xml:"namespace,attr"`
	MetricProfile    string   `xml:"metricprofiles,attr"`
	EndpointGroup    string   `xml:"endpointgroup,attr"`
	MetricOperation  string   `xml:"metricoperation,attr"`
	ProfileOperation string   `xml:"profileoperation,attr"`
	Ands             Andtest
}

type Resulttest struct {
	XMLName      xml.Name      `xml:"root"`
	Profiletests []Profiletest `xml:"profile"`
}

//  ServiceIn struct to represent maps of sercices stored in MongoDB
type ServiceIn struct {
	ServiceSetIn map[string]string `bson:"services"`
	Operator     string            `bson:"operation"`
}

// Prepare maps for services and groups for ap1
var apGroup1 = ServiceIn{
	ServiceSetIn: map[string]string{
		"ap1-service1": "OR",
		"ap1-service2": "AND",
		"ap1-service3": "AND"},
	Operator: "OR"}

var apGroup2 = ServiceIn{
	ServiceSetIn: map[string]string{
		"ap1-service4": "AND",
		"ap1-service5": "OR",
		"ap1-service6": "AND"},
	Operator: "AND"}

var profileGroup1 = map[string]ServiceIn{
	"compute": apGroup1,
	"storage": apGroup2}

// Prepare maps for services and groups for ap2
var apGroup3 = ServiceIn{
	ServiceSetIn: map[string]string{
		"ap2-service1": "OR",
		"ap2-service2": "AND",
		"ap2-service3": "AND"},
	Operator: "OR"}

var apGroup4 = ServiceIn{
	ServiceSetIn: map[string]string{
		"ap2-service4": "AND",
		"ap2-service5": "AND",
		"ap2-service6": "OR"},
	Operator: "OR"}

var profileGroup2 = map[string]ServiceIn{
	"compute": apGroup3,
	"storage": apGroup4}

// Setup the Test Environment
// This function runs before any test and setups the environment
// A test configuration object is instantiated using a reference
// to testdb: AR_test. Also here is are instantiated some expected
// xml response validation messages (authorization,crud responses).
// Also the testdb is seeded with two av.profiles ("ap1","ap2")
// and with an authorization token:"S3CR3T"
func (suite *AvProfileTestSuite) SetupTest() {

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
    db = "AR_test"
`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.respProfileCreated = " <root>\n" +
		"   <Message>Availability Profile record successfully created</Message>\n </root>"

	suite.respProfileUpdated = " <root>\n" +
		"   <Message>Availability Profile was successfully updated</Message>\n </root>"

	suite.respProfileDeleted = " <root>\n" +
		"   <Message>Availability Profile was successfully deleted</Message>\n </root>"

	suite.respUnauthorized = "Unauthorized"

	suite.respNoID = " <root>\n" +
		"   <Message>No profile matching the requested id</Message>\n </root>"

	suite.respBadJSON = " <root>\n" +
		"   <Message>Malformated json input data</Message>\n </root>"

	// Connect to mongo testdb
	session, _ := mongo.OpenSession(suite.cfg.MongoDB)

	// Add authentication token to mongo testdb
	seedAuth := bson.M{"apiKey": "S3CR3T"}
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "authentication", seedAuth)

	// seed mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Open DB session
	c := session.DB(suite.cfg.MongoDB.Db).C("aps")

	// Insert first seed profile
	c.Insert(
		bson.M{"name": "ap1",
			"namespace":        "namespace1",
			"metricprofiles":   []string{"metricprofile01"},
			"endpointgroup":    "sites",
			"metricoperation":  "AND",
			"profileoperation": "AND",
			"groups":           profileGroup1})

	// Insert second seed profile
	c.Insert(
		bson.M{"name": "ap2",
			"namespace":        "namespace2",
			"metricprofiles":   []string{"metricprofile02"},
			"endpointgroup":    "sites",
			"metricoperation":  "AND",
			"profileoperation": "AND",
			"groups":           profileGroup2})

}

// Testing creation of a profile  using POST request.
// During Setup of the test environment the testdb is seeded with
// two availability profiles ("ap1","ap2"). Mongo assigns
// two object _ids on these profiles which cannot predict so,
// we have to read them from the database and insert them in the
// expected xml response string.
func (suite *AvProfileTestSuite) TestCreateProfile() {

	// create json input data for the request
	postData := `
      {
          "name": "fresh_test_profile",
          "namespace": "test_namespace",
          "metricprofiles": ["test_metricprofile"],
          "endpointgroup": "sites",
          "metricoperation": "AND",
          "profileoperation": "AND",
          "groups" : {
            "compute" : {
              "services" : {
                "service1" : "OR",
                "service2" : "AND",
                "service3" : "OR"
                },
                "operation" : "OR"
                },
                "storage" : {
                  "services" : {
                    "service4" : "AND",
                    "service5" : "OR",
                    "service6" : "AND"
                    },
                    "operation" : "OR"
                  }
                }

      }
      `
	// Prepare the request object
	request, _ := http.NewRequest("POST", "", strings.NewReader(postData))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Create(request, suite.cfg)

	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.respProfileCreated, string(output), "Response body mismatch")

	// Remove the profile not to contaminate other tests
	// Open session to mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Open collection aps
	c := session.DB(suite.cfg.MongoDB.Db).C("aps")
	// Remove the specific profile inserted during this test
	c.Remove(bson.M{"name": "fresh_test_profile"})

}

// TODO: The TestReadProfile function will be refactored in future releases.
// The GO runtime randmonizes the map iteration order.
// So when we try to compare the produced xml view
// with a static one declared in unit test it always fails
// because the xml attributes' order is random.
// For now we use two tests. One to check if the XML structure is the desired one
// and the XML can be unmarshaled using the test structs,
// and a comparizon between the XMLs as strings using regex to match the actual attributes values.

// Testing Reading of profile list using GET request.
// During Setup of the test environment the testdb is seeded with
// two availability profiles ("ap1","ap2"). Mongo assigns
// two object _ids on these profiles which cannot predict so,
// we have to read them from the database and insert them in the
// expected xml response string.
func (suite *AvProfileTestSuite) TestReadProfile() {

	// Open a session to mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Open availability profile collection: aps
	c := session.DB(suite.cfg.MongoDB.Db).C("aps")
	// Instantiate a AvProfile struct to hold bson results
	results := AvailabilityProfileOutput{}
	// Query first seed profile - name:ap1
	c.Find(bson.M{"name": "ap1"}).One(&results)
	// Grab from results ObjectId and convert it to string: Hex() method
	id1 := (results.ID.Hex())
	// Query second seed profile - name:ap2
	c.Find(bson.M{"name": "ap2"}).One(&results)
	id2 := (results.ID.Hex())
	// Hold a string multiline literal including the two profile ids retrieved.
	// This would be the representation for the XML expected schema.
	schema := ` <root>
   <profile id="` + id1 + `" name="ap1" namespace="namespace1" metricprofiles="metricprofile01" endpointgroup="sites" metricoperation="AND" profileoperation="AND">
     <AND>
       <OR>
         <Group service_flavor="ap1-service([1-6]{1})" operation="(AND|OR)"></Group>
         <Group service_flavor="ap1-service([1-6]{1})" operation="(AND|OR)"></Group>
         <Group service_flavor="ap1-service([1-6]{1})" operation="(AND|OR)"></Group>
       </OR>
       <OR>
         <Group service_flavor="ap1-service([1-6]{1})" operation="(AND|OR)"></Group>
         <Group service_flavor="ap1-service([1-6]{1})" operation="(AND|OR)"></Group>
         <Group service_flavor="ap1-service([1-6]{1})" operation="(AND|OR)"></Group>
       </OR>
     </AND>
   </profile>
   <profile id="` + id2 + `" name="ap2" namespace="namespace2" metricprofiles="metricprofile02" endpointgroup="sites" metricoperation="AND" profileoperation="AND">
     <AND>
       <OR>
         <Group service_flavor="ap2-service([1-6]{1})" operation="(AND|OR)"></Group>
         <Group service_flavor="ap2-service([1-6]{1})" operation="(AND|OR)"></Group>
         <Group service_flavor="ap2-service([1-6]{1})" operation="(AND|OR)"></Group>
       </OR>
       <OR>
         <Group service_flavor="ap2-service([1-6]{1})" operation="(AND|OR)"></Group>
         <Group service_flavor="ap2-service([1-6]{1})" operation="(AND|OR)"></Group>
         <Group service_flavor="ap2-service([1-6]{1})" operation="(AND|OR)"></Group>
       </OR>
     </AND>
   </profile>
 </root>`

	// Prepare the request object
	request, _ := http.NewRequest("GET", "", strings.NewReader(""))
	// Pass request to controller calling List() handler method
	code, _, output, _ := List(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Unmarshal the produced XML using the test structs
	v := &Resulttest{}

	xmlErr := xml.Unmarshal(output, v)

	// Unmarshal the test schema that will be used in the comparison
	d := &Resulttest{}
	_ = xml.Unmarshal([]byte(schema), d)

	if xmlErr != nil {
		fmt.Printf("error: %v", xmlErr)
		suite.Fail("Unmarshal error: ", xmlErr.Error())
	}

	// Compare the expected and actual XML result
	suite.Regexp(schema, string(output), "Response body mismatch")

	// Compare only the structure of the XML
	cmp := reflect.DeepEqual(v, d)
	if cmp != true {
		suite.Fail("XML schema mismatch")
	}
}

// Testing update of a profile  using POST request.
// During Setup of the test environment the testdb is seeded with
// two availability profiles ("ap1","ap2"). Mongo assigns
// two object _ids on these profiles which cannot predict so,
// we have to read them from the database and insert them in the
// expected xml response string.
func (suite *AvProfileTestSuite) TestUpdateProfile() {

	// We will make update to ap2 profile
	putData := `
      {
          "name": "updated-ap2",
          "namespace": "updated-ap2-namespace",
          "metricprofiles": ["updated-ap2-metricprofile"],
          "endpointgroup": "sites",
          "metricoperation": "AND",
          "profileoperation": "AND",
          "groups" : {
            "compute" : {
              "services" : {
                "updated-service1" : "OR",
                "updated-service2" : "AND"
                },
                "operation" : "OR"
                },
                "storage" : {
                  "services" : {
                    "updated-service3" : "AND",
                    "updated-service4" : "OR",
                    "updated-service5" : "AND"
                    },
                    "operation" : "OR"
                  }
                }
      }
      `
	// Read the id
	// Open a session to mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Open availability profile collection: aps
	c := session.DB(suite.cfg.MongoDB.Db).C("aps")
	// Instantiate a AvProfile struct to hold bson results
	results := AvailabilityProfileOutput{}
	// Query first seed profile - name:ap2
	c.Find(bson.M{"name": "ap2"}).One(&results)
	// Grab from results ObjectId and convert it to string: Hex() method
	id2 := (results.ID.Hex())
	// Prepare the request object (use id2 for path)
	request, _ := http.NewRequest("PUT", "/api/v1/AP/"+id2, strings.NewReader(putData))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")
	// Execute the request in the controller
	code, _, output, _ := Update(request, suite.cfg)
	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.respProfileUpdated, string(output), "Response body mismatch")

	// Reestablish ap2 profile (remove and reinsert)
	c.Remove(bson.M{"name": "ap2"})
	c.Insert(
		bson.M{"name": "ap2",
			"namespace":        "namespace2",
			"metricprofiles":   []string{"metricprofile02"},
			"endpointgroup":    "sites",
			"metricoperation":  "AND",
			"profileoperation": "AND",
			"groups":           profileGroup2})

}

// TODO: The TestDeleteProfile function will be refactored in future releases.
// The GO runtime randmonizes the map iteration order.
// So when we try to compare the produced xml view
// with a static one declared in unit test it always fails
// because the xml attributes' order is random.
// For now we use two tests. One to check if the XML structure is the desired one
// and the XML can be unmarshaled using the test structs,
// and a comparizon between the XMLs as strings using regex to match the actual attributes values.

// Testing deletion of a profile  using DELETE request.
// During Setup of the test environment the testdb is seeded with
// two availability profiles ("ap1","ap2"). Mongo assigns
// two object _ids on these profiles which cannot predict so,
// we have to read them from the database and insert them in the
// expected xml response string. In this test ap2 is deleted and then
// profile list is read again (with get request) to validate that
// only ap1 is left present.
func (suite *AvProfileTestSuite) TestDeleteProfile() {

	// We will delete ap2 profile so we expect to find only ap1 in list

	// Grab the two _ids of profiles: ap1, ap2
	// Open a session to mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Open availability profile collection: aps
	c := session.DB(suite.cfg.MongoDB.Db).C("aps")
	// Instantiate a AvProfile struct to hold bson results
	results := AvailabilityProfileOutput{}
	// Query first seed profile - name:ap1
	c.Find(bson.M{"name": "ap1"}).One(&results)
	// Grab from results ObjectId and convert it to string: Hex() method
	id1 := (results.ID.Hex())
	// Query second seed profile - name:ap2
	c.Find(bson.M{"name": "ap2"}).One(&results)
	id2 := (results.ID.Hex())

	// Prepare the expected xml response after deleting ap2
	schema := ` <root>
   <profile id="` + id1 + `" name="ap1" namespace="namespace1" metricprofiles="metricprofile01" endpointgroup="sites" metricoperation="AND" profileoperation="AND">
     <AND>
       <OR>
         <Group service_flavor="ap1-service([1-6]{1})" operation="(AND|OR)"></Group>
         <Group service_flavor="ap1-service([1-6]{1})" operation="(AND|OR)"></Group>
         <Group service_flavor="ap1-service([1-6]{1})" operation="(AND|OR)"></Group>
       </OR>
       <OR>
         <Group service_flavor="ap1-service([1-6]{1})" operation="(AND|OR)"></Group>
         <Group service_flavor="ap1-service([1-6]{1})" operation="(AND|OR)"></Group>
         <Group service_flavor="ap1-service([1-6]{1})" operation="(AND|OR)"></Group>
       </OR>
     </AND>
   </profile>
 </root>`

	// Prepare the request object (use id2 for path)
	request, _ := http.NewRequest("DELETE", "/api/v1/AP/"+id2, strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Delete(request, suite.cfg)
	// Check proper response that the profile successfully deleted
	suite.Equal(200, code, "Internal Server Error")

	// Compare the expected and actual response
	suite.Equal(suite.respProfileDeleted, string(output), "Response body mismatch")

	// Double-check that the profile is actually missing from the profile list
	request, _ = http.NewRequest("GET", "", strings.NewReader(""))

	// Pass request to controller calling List() handler method
	code, _, output, _ = List(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")

	// Compare the expected and actual XML
	suite.Regexp(schema, string(output), "Response body mismatch")

	// Reestablish ap2 profile (reinsert)
	c.Insert(
		bson.M{"name": "ap2",
			"namespace":        "namespace2",
			"metricprofiles":   []string{"metricprofile02"},
			"endpointgroup":    "sites",
			"metricoperation":  "AND",
			"profileoperation": "AND",
			"groups":           profileGroup2})

}

// This function tests calling the create request (POST) and providing
// a wrong api-key. The response should be unauthorized
func (suite *AvProfileTestSuite) TestCreateUnauthorized() {
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

// This function tests calling the update profile request (PUT) and providing
// a wrong api-key. The response should be unauthorized
func (suite *AvProfileTestSuite) TestUpdateUnauthorized() {
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

// This function tests calling the remove av.profile request (DELETE) and providing
// a wrong api-key. The response should be unauthorized
func (suite *AvProfileTestSuite) TestDeleteUnauthorized() {
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

// This function tests calling the update av.profile request (PUT) and providing
// a wrong profile id. The response should be profile with id doesn't exist
func (suite *AvProfileTestSuite) TestUpdateBadId() {
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "/api/v1/AP/wrongid", strings.NewReader("{}"))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Update(request, suite.cfg)

	suite.Equal(400, code, "Internal Server Error")
	suite.Equal(suite.respNoID, string(output), "Response body mismatch")
}

// This function tests calling the update av.profile request (DELETE) and providing
// a wrong profile id. The response should be profile with id doesn't exist
func (suite *AvProfileTestSuite) TestDeleteBadId() {
	// Prepare the request object
	request, _ := http.NewRequest("DELETE", "/api/v1/AP/wrongid", strings.NewReader("{}"))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Delete(request, suite.cfg)

	suite.Equal(400, code, "Internal Server Error")
	suite.Equal(suite.respNoID, string(output), "Response body mismatch")
}

// This function tests calling the create av.profile request (POST) and providing
// bad json input. The response should be malformed json
func (suite *AvProfileTestSuite) TestCreateBadJson() {
	// Find an existing id
	// Open a session to mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Open availability profile collection: aps
	c := session.DB(suite.cfg.MongoDB.Db).C("aps")
	// Instantiate a AvProfile struct to hold bson results
	results := AvailabilityProfileOutput{}
	// Query first seed profile - name:ap1
	c.Find(bson.M{"name": "ap1"}).One(&results)
	// Grab from results ObjectId and convert it to string: Hex() method
	id1 := (results.ID.Hex())

	// Prepare the request object
	request, _ := http.NewRequest("POST", "/api/v1/AP/"+id1, strings.NewReader("{ bad json"))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Create(request, suite.cfg)

	suite.Equal(400, code, "Internal Server Error")
	suite.Equal(suite.respBadJSON, string(output), "Response body mismatch")
}

// This function tests calling the update av.profile request (PUT) and providing
// bad json input. The response should be malformed json
func (suite *AvProfileTestSuite) TestUpdateBadJson() {
	// Find an existing id
	// Open a session to mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Open availability profile collection: aps
	c := session.DB(suite.cfg.MongoDB.Db).C("aps")
	// Instantiate a AvProfile struct to hold bson results
	results := AvailabilityProfileOutput{}
	// Query first seed profile - name:ap1
	c.Find(bson.M{"name": "ap1"}).One(&results)
	// Grab from results ObjectId and convert it to string: Hex() method
	id1 := (results.ID.Hex())
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "/api/v1/AP/"+id1, strings.NewReader("{ bad json"))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Update(request, suite.cfg)

	suite.Equal(400, code, "Internal Server Error")
	suite.Equal(suite.respBadJSON, string(output), "Response body mismatch")
}

// This function is actually called in the end of all tests
// and clears the test environment.
// Mainly it's purpose is to drop the testdb
func (suite *AvProfileTestSuite) TearDownTest() {

	session, _ := mongo.OpenSession(suite.cfg.MongoDB)
	session.DB("AR_test").DropDatabase()

}

// This is the first function called when go test is issued
func TestAvProfileTestSuite(t *testing.T) {
	suite.Run(t, new(AvProfileTestSuite))
}
