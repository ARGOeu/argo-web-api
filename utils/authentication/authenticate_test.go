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

package authentication

import (
	"net/http"
	"strings"
	"testing"

	"code.google.com/p/gcfg"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// This is a util. suite struct used in tests (see pkg "testify")
type AuthenticationProfileTestSuite struct {
	suite.Suite
	cfg              config.Config
	tenantdb         string
	tenantpassword   string
	tenantusername   string
	tenantstorename  string
	clientkey        string
	respUnauthorized string
}

// Setup the Test Environment
// This function runs before any test and setups the environment
func (suite *AuthenticationProfileTestSuite) SetupTest() {

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
    db = "AR_test_Authentication"
    `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.respUnauthorized = "Unauthorized"
	suite.tenantdb = "argo_egi_AR_data"
	suite.clientkey = "mysecretcombination"
	suite.tenantpassword = "h4shp4ss"
	suite.tenantusername = "johndoe"

	// seed mongo
	session, err := mongo.OpenSession(suite.cfg.MongoDB)
	if err != nil {
		panic(err)
	}
	defer mongo.CloseSession(session)

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
					"store":    "ar",
					"server":   "localhost",
					"port":     27017,
					"database": suite.tenantdb,
					"username": suite.tenantusername,
					"password": suite.tenantpassword,
				},
				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_egi_metric_data",
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

}

func (suite *AuthenticationProfileTestSuite) TestAuthentication() {

	request, _ := http.NewRequest("GET", "", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)

	tenantdbconfig, err := AuthenticateTenant(request.Header, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(nil, err, "authetication probelm")

	//Check the data is conveyed correctly
	suite.Regexp(tenantdbconfig.Db, suite.tenantdb, "Database mismatch")
	suite.Regexp(tenantdbconfig.Username, suite.tenantusername, "Username mismatch")
	suite.Regexp(tenantdbconfig.Password, suite.tenantpassword, "Password mismatch")

}

//TearDownTest to tear down every test
func (suite *AuthenticationProfileTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantdb).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()

}

func TestRecompuptationsTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationProfileTestSuite))
}
