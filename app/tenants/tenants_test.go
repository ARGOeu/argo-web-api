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

package tenants

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

// This is a util. suite struct used in tests (see pkg "testify")
type TenantsTestSuite struct {
	suite.Suite
	cfg                 config.Config
	resp_unauthorized   string
	resp_tenantsList    string
}

func (suite *TenantsTestSuite) SetupTest() {

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
    db = "argo_core_test"
`
	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)
	suite.resp_unauthorized = "Unauthorized"
	suite.resp_tenantsList = `<root>
 <Tenant name="PREIS"></Tenant>
 <Tenant name="USDAT"></Tenant>
</root>`

	// Connect to mongo testdb
	session, _ := mongo.OpenSession(suite.cfg.MongoDB)

	// Add authentication token to mongo testdb
	seed_auth := bson.M{"name" : "John Doe", "email" : "john.doe@example.com", "api_key" : "mysecretcombination"}
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "authentication", seed_auth)

	// seed mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Insert first seed profile
	c := session.DB(suite.cfg.MongoDB.Db).C("tenants")
	c.Insert(bson.M{"name" : "USDAT",
		"db_conf" : []bson.M{ bson.M{"server" : "127.0.0.1", "port" : 27017, "database" : "USDAT_test"} } ,
		"users" : []bson.M{ bson.M{"name" : "Jack Doe", "email" : "jack.doe@example.com", "api_key" : "anothersecret"} }})
	c.Insert(bson.M{"name" : "PREIS",
		"db_conf" : []bson.M{ bson.M{"server" : "127.0.0.1", "port" : 27017, "database" : "PREIS_test"} } ,
		"users" : []bson.M{ bson.M{"name" : "Jill Doe", "email" : "jill.doe@example.com", "api_key" : "onemoresecret"} }})

}

func (suite *TenantsTestSuite) TestListTenants() {

	// Prepare the request object
	request, _ := http.NewRequest("GET", "/api/v1/tenants", strings.NewReader(""))
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "mysecretcombination")
	// Execute the request in the controller
	code, _, output, _ := List(request, suite.cfg)
	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.resp_tenantsList, string(output), "Response body mismatch")

	// Prepare new request object
	request, _ = http.NewRequest("GET", "/api/v1/tenants", strings.NewReader(""))
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "wrongkey")
	// Execute the request in the controller
	code, _, output, _ = List(request, suite.cfg)
	suite.Equal(401, code, "Should have gotten return code 401 (Unauthorized)")
	suite.Equal(suite.resp_unauthorized, string(output), "Should have gotten reply Unauthorized")

	// Remove the profile not to contaminate other tests
	// Open session to mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Open collection authentication
	c := session.DB(suite.cfg.MongoDB.Db).C("authentication")
	// Remove the specific entries inserted during this test
	c.Remove(bson.M{"name": "John Doe"})

	// Open collection tenants
	c = session.DB(suite.cfg.MongoDB.Db).C("tenants")
	// Remove the specific entries inserted during this test
	c.Remove(bson.M{"name": "USDAT"})
	c.Remove(bson.M{"name": "PREIS"})

}

//TearDownTest to tear down every test
func (suite *TenantsTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()

}

func TestTenantsTestSuite(t *testing.T) {
	suite.Run(t, new(TenantsTestSuite))
}
