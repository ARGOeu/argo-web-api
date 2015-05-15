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
	"github.com/argoeu/ar-web-api/utils/config"
	"github.com/argoeu/ar-web-api/utils/mongo"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strings"
	"testing"
)

// This is a util. suite struct used in tests (see pkg "testify")
type AvProfileTestSuite struct {
	suite.Suite
	cfg                 config.Config
	resp_profileCreated string
	resp_profileUpdated string
	resp_profileDeleted string
	resp_unauthorized   string
	resp_no_id          string
	resp_bad_json       string
}

// Setup the Test Enviroment
// This function runs before any test and setups the enviroment
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

	suite.resp_profileCreated = " <root>\n" +
		"   <Message>Availability Profile record successfully created</Message>\n </root>"

	suite.resp_profileUpdated = " <root>\n" +
		"   <Message>Availability Profile was successfully updated</Message>\n </root>"

	suite.resp_profileDeleted = " <root>\n" +
		"   <Message>Availability Profile was successfully deleted</Message>\n </root>"

	suite.resp_unauthorized = "Unauthorized"

	suite.resp_no_id = " <root>\n" +
		"   <Message>No profile matching the requested id</Message>\n </root>"

	suite.resp_bad_json = " <root>\n" +
		"   <Message>Malformated json input data</Message>\n </root>"

	// Connect to mongo testdb
	session, _ := mongo.OpenSession(suite.cfg)

	// Add authentication token to mongo testdb
	seed_auth := bson.M{"apiKey": "S3CR3T"}
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "authentication", seed_auth)

	// seed mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Insert first seed profile
	c := session.DB(suite.cfg.MongoDB.Db).C("aps")
	c.Insert(bson.M{"name": "ap1", "namespace": "namespace1", "poems": []string{"poem01"},
		"groups": [][]string{
			[]string{"ap1-service1", "ap1-service2", "ap1-service3"},
			[]string{"ap1-service4", "ap1-service5", "ap1-service6"}}})
	// Insert first seed profile
	c.Insert(bson.M{"name": "ap2", "namespace": "namespace2", "poems": []string{"poem02"},
		"groups": [][]string{
			[]string{"ap2-service1", "ap2-service2", "ap2-service3"},
			[]string{"ap2-service4", "ap2-service5", "ap2-service6"}}})

}

// Testing creation of a profile  using POST request.
// During Setup of the test enviroment the testdb is seeded with
// two availability profiles ("ap1","ap2"). Mongo assigns
// two object _ids on these profiles which cannot predict so,
// we have to read them from the database and insert them in the
// expected xml response string.
func (suite *AvProfileTestSuite) TestCreateProfile() {

	// create json input data for the request
	post_data := `
    {
        "name": "fresh_test_profile",
        "namespace": "test_namespace",
        "poems": ["test_poem"],
        "groups": [
            [
               "service1",
               "service2",
               "service3"
            ],
            [
                "service4",
                "service5",
                "service6"
            ]
        ]
    }
    `
	// Prepare the request object
	request, _ := http.NewRequest("POST", "", strings.NewReader(post_data))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Create(request, suite.cfg)

	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.resp_profileCreated, string(output), "Response body mismatch")

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

// Testing Reading of profile list using GET request.
// During Setup of the test enviroment the testdb is seeded with
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
	// Hold a string multiline literal including the two profile ids retrieved
	profile_list_xml := ` <root>
   <profile id="` + id1 + `" name="ap1" namespace="namespace1" poems="poem01">
     <AND>
       <OR>
         <Group service_flavor="ap1-service1"></Group>
         <Group service_flavor="ap1-service2"></Group>
         <Group service_flavor="ap1-service3"></Group>
       </OR>
       <OR>
         <Group service_flavor="ap1-service4"></Group>
         <Group service_flavor="ap1-service5"></Group>
         <Group service_flavor="ap1-service6"></Group>
       </OR>
     </AND>
   </profile>
   <profile id="` + id2 + `" name="ap2" namespace="namespace2" poems="poem02">
     <AND>
       <OR>
         <Group service_flavor="ap2-service1"></Group>
         <Group service_flavor="ap2-service2"></Group>
         <Group service_flavor="ap2-service3"></Group>
       </OR>
       <OR>
         <Group service_flavor="ap2-service4"></Group>
         <Group service_flavor="ap2-service5"></Group>
         <Group service_flavor="ap2-service6"></Group>
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
	// Compare the expected and actual xml response
	suite.Equal(profile_list_xml, string(output), "Response body mismatch")
}

// Testing update of a profile  using POST request.
// During Setup of the test enviroment the testdb is seeded with
// two availability profiles ("ap1","ap2"). Mongo assigns
// two object _ids on these profiles which cannot predict so,
// we have to read them from the database and insert them in the
// expected xml response string.
func (suite *AvProfileTestSuite) TestUpdateProfile() {

	// We will make update to ap2 profile
	put_data := `
    {
        "name": "updated-ap2",
        "namespace": "updated-ap2-namespace",
        "poems": ["updated-ap2-poem"],
        "groups": [
            [
               "updated-srv1",
               "updated-srv2"
            ],
            [
                "updated-srv3",
                "updated-srv4",
                "updated-srv5",
                "updated-srv6"
            ]
        ]
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
	// Query first seed profile - name:ap1
	c.Find(bson.M{"name": "ap2"}).One(&results)
	// Grab from results ObjectId and convert it to string: Hex() method
	id2 := (results.ID.Hex())
	// Prepare the request object (use id2 for path)
	request, _ := http.NewRequest("PUT", "/api/v1/AP/"+id2, strings.NewReader(put_data))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Update(request, suite.cfg)

	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.resp_profileUpdated, string(output), "Response body mismatch")

	// Reestablish ap2 profile (remove and reinsert)
	c.Remove(bson.M{"name": "ap2"})
	c.Insert(bson.M{"name": "ap2", "namespace": "namespace2", "poems": []string{"poem02"},
		"groups": [][]string{
			[]string{"ap2-service1", "ap2-service2", "ap2-service3"},
			[]string{"ap2-service4", "ap2-service5", "ap2-service6"}}})

}

// Testing deletion of a profile  using DELETE request.
// During Setup of the test enviroment the testdb is seeded with
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
	profile_list_xml := ` <root>
   <profile id="` + id1 + `" name="ap1" namespace="namespace1" poems="poem01">
     <AND>
       <OR>
         <Group service_flavor="ap1-service1"></Group>
         <Group service_flavor="ap1-service2"></Group>
         <Group service_flavor="ap1-service3"></Group>
       </OR>
       <OR>
         <Group service_flavor="ap1-service4"></Group>
         <Group service_flavor="ap1-service5"></Group>
         <Group service_flavor="ap1-service6"></Group>
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
	suite.Equal(suite.resp_profileDeleted, string(output), "Response body mismatch")

	// Double-check that the profile is actually missing from the profile list
	request, _ = http.NewRequest("GET", "", strings.NewReader(""))

	// Pass request to controller calling List() handler method
	code, _, output, _ = List(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(profile_list_xml, string(output), "Response body mismatch")

	// Reestablish ap2 profile (reinsert)
	c.Insert(bson.M{"name": "ap2", "namespace": "namespace2", "poems": []string{"poem02"},
		"groups": [][]string{
			[]string{"ap2-service1", "ap2-service2", "ap2-service3"},
			[]string{"ap2-service4", "ap2-service5", "ap2-service6"}}})

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
	suite.Equal(suite.resp_unauthorized, string(output), "Response body mismatch")
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
	suite.Equal(suite.resp_unauthorized, string(output), "Response body mismatch")
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
	suite.Equal(suite.resp_unauthorized, string(output), "Response body mismatch")
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
	suite.Equal(suite.resp_no_id, string(output), "Response body mismatch")
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
	suite.Equal(suite.resp_no_id, string(output), "Response body mismatch")
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
	suite.Equal(suite.resp_bad_json, string(output), "Response body mismatch")
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
	suite.Equal(suite.resp_bad_json, string(output), "Response body mismatch")
}

// This function is actually called in the end of all tests
// and clears the test enviroment.
// Mainly it's purpose is to drop the testdb
func (suite *AvProfileTestSuite) TearDownTest() {

	session, _ := mongo.OpenSession(suite.cfg)

	session.DB("AR_test").DropDatabase()

}

// This is the first function called when go test is issued
func TestAvProfileTestSuite(t *testing.T) {
	suite.Run(t, new(AvProfileTestSuite))
}
