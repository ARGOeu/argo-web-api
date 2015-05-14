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
	"fmt"
	"github.com/argoeu/ar-web-api/utils/config"
	"github.com/argoeu/ar-web-api/utils/mongo"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strings"
	"testing"
)

// Build a test suite as a struct with specific enviroment parameters
type AvProfileTestSuite struct {
	suite.Suite
	cfg                 config.Config
	resp_profileCreated string
	reps_unauthorized   string
}

// Setup the test suite
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

	suite.resp_profileCreated = " <root>\n"
	+"<Message>Availability Profile record successfully created</Message>\n"
	+"</root>"

	suite.reps_unauthorized = "Unauthorized"

	// Connect to mongo testdb
	session, _ := mongo.OpenSession(suite.cfg)

	// Add authentication token to mongo testdb
	seed_auth := bson.M{"apiKey": "S3CR3T"}
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "authentication", seed_auth)

	// Ready first profile
	prof1 := `
    {
        "name": "ap1",
        "namespace": "namespace1",
        "poems": ["poem01"],
        "groups": [
            [
               "ap1-service1",
               "ap1-service2",
               "ap1-service3"
            ],
            [
                "ap1-service4",
                "ap1-service5",
                "ap1-service6"
            ]
        ]
    }
    `
	// Ready second profile
	prof2 := `
    {
        "name": "ap2",
        "namespace": "namespace2",
        "poems": ["poem02"],
        "groups": [
            [
               "ap2-service1",
               "ap2-service2",
               "ap2-service3"
            ],
            [
                "ap2-service4",
                "ap2-service5",
                "ap2-service6"
            ]
        ]
    }
    `
	//Execture requests
	// scaffold the request object
	request1, _ := http.NewRequest("POST", "", strings.NewReader(prof1))
	// add the content-type header to application/json
	request1.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request1.Header.Set("x-api-key", "S3CR3T")
	// Execute the request in the controller
	Create(request1, suite.cfg)

	// scaffold the request object
	request2, _ := http.NewRequest("POST", "", strings.NewReader(prof2))
	// add the content-type header to application/json
	request2.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request2.Header.Set("x-api-key", "S3CR3T")
	// Execute the request in the controller
	Create(request2, suite.cfg)

}

// Testing the insertion of profiles using POST
func (suite *AvProfileTestSuite) TestCreateProfile() {

	// create request
	post_body := `
    {
        "name": "test_profile",
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
	// scaffold the request object
	request, _ := http.NewRequest("POST", "", strings.NewReader(post_body))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Create(request, suite.cfg)

	suite.Equal(200, code, "Internal Server Error")
	fmt.Println("testing")
	suite.Equal(string(output), suite.resp_profileCreated, "Response body mismatch")
}

func (suite *AvProfileTestSuite) TestInsertProfile() {

	// create request
	post_body := `
    {
        "name": "test_profile",
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
	// scaffold the request object
	request, _ := http.NewRequest("POST", "", strings.NewReader(post_body))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "S3CR3T")

	// Execute the request in the controller
	code, _, output, _ := Create(request, suite.cfg)

	suite.Equal(200, code, "Internal Server Error")
	fmt.Println("testing")
	suite.Equal(string(output), suite.resp_profileCreated, "Response body mismatch")
}

func (suite *AvProfileTestSuite) TearDownTest() {

	session, _ := mongo.OpenSession(suite.cfg)

	session.DB("AR_test").DropDatabase()

}

// Here is the normal golang testing function that runs when
// go test is issued
func TestAvProfileTestSuite(t *testing.T) {
	suite.Run(t, new(AvProfileTestSuite))
}
