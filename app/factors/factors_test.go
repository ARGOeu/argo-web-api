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

package factors

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
type FactorsTestSuite struct {
	suite.Suite
	cfg                 config.Config
	tenantcfg           config.MongoConfig
	resp_nokeyprovided  string
	resp_unauthorized   string
	resp_factorsList    string
}

func (suite *FactorsTestSuite) SetupTest() {

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
    db = "argo_core_test_factors"
`
	_ = gcfg.ReadStringInto(&suite.cfg, coreConfig)
	suite.resp_nokeyprovided = "404 page not found"
	suite.resp_unauthorized = "Unauthorized"
	suite.resp_factorsList = `<root>
 <Factor site="CETA-GRID" weight="5406"></Factor>
 <Factor site="CFP-IST" weight="1019"></Factor>
 <Factor site="CIEMAT-LCG2" weight="14595"></Factor>
</root>`

	// Connect to mongo coredb
	session, err := mongo.OpenSession(suite.cfg.MongoDB)
	defer session.Close()

	// Add authentication token to mongo coredb
	seed_auth := bson.M{"name" : "TEST",
		"db_conf" : []bson.M{ bson.M{"server" : "127.0.0.1", "port" : 27017, "database" : "AR_test"} } ,
		"users" : []bson.M{ bson.M{"name" : "Jack Doe", "email" : "jack.doe@example.com", "api_key" : "secret"} }}
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "tenants", seed_auth)

	mongo.CloseSession(session)

	// TODO: I don't like it here that I rewrite the test data. 
	// However, this is a test for factors, not for AuthenticateTenant function. 
	suite.tenantcfg.Host = "127.0.0.1"
	suite.tenantcfg.Port = 27017
	suite.tenantcfg.Db = "AR_test"

	// seed tenantdb mongo
	session, err = mgo.Dial(suite.tenantcfg.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Add a few factors in collection
	c := session.DB(suite.tenantcfg.Db).C("weights")
	c.Insert(bson.M{ "hepspec" : 14595, "name" : "CIEMAT-LCG2" })
	c.Insert(bson.M{ "hepspec" : 1019, "name" : "CFP-IST" })
	c.Insert(bson.M{ "hepspec" : 5406, "name" : "CETA-GRID" })

	mongo.CloseSession(session)

}

func (suite *FactorsTestSuite) TestListFactors() {

	// Prepare the request object
	request, _ := http.NewRequest("GET", "/api/v1/factors", strings.NewReader(""))
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "secret")
	// Execute the request in the controller
	code, _, output, _ := List(request, suite.cfg)
	suite.Equal(200, code, "Internal Server Error")
	suite.Equal(suite.resp_factorsList, string(output), "Response body mismatch")

	// Prepare new request object
	request, _ = http.NewRequest("GET", "/api/v1/tenants", strings.NewReader(""))
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "wrongkey")
	// Execute the request in the controller
	code, _, output, _ = List(request, suite.cfg)
	suite.Equal(401, code, "Should have gotten return code 401 (Unauthorized)")
	suite.Equal(suite.resp_unauthorized, string(output), "Should have gotten reply Unauthorized")

	// Remove the test data from core db not to contaminate other tests
	// Open session to core mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
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
	defer session.Close()
	// Open collection authentication
	c = session.DB(suite.tenantcfg.Db).C("weights")
	// Remove the specific entries inserted during this test
	c.Remove(bson.M{"name": "CIEMAT-LCG2"})
	c.Remove(bson.M{"name": "CFP-IST"})
	c.Remove(bson.M{"name": "CETA-GRID"})
}

//TearDownTest to tear down every test
func (suite *FactorsTestSuite) TearDownTest() {

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

func TestFactorsTestSuite(t *testing.T) {
	suite.Run(t, new(FactorsTestSuite))
}
