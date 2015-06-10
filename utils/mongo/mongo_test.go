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

package mongo

import (
	"testing"

	"code.google.com/p/gcfg"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// This is a utility suite struct used in tests (see pkg "testify")
type mongoTestSuite struct {
	suite.Suite
	cfg    config.Config
	result bson.M
}

// Setup the Test Environment
// This function runs before any test and setups the environment
func (suite *mongoTestSuite) SetupTest() {

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
    db = "argo_core_test_mongo"
    `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	//	suite.result = bson.M{ "name" : "Westeros", "db_conf" : bson.M{ "server" : "localhost", "port" : 27017, "database" : "argo_EGI" } }

	// seed mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Seed database with tenants
	c := session.DB(suite.cfg.MongoDB.Db).C("tenants")
	c.Insert(
		bson.M{"name": "Westeros",
			"db_conf": bson.M{
				"server":   "localhost",
				"port":     27017,
				"database": "argo_GOT",
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
					"api_key": "sansa3",
				},
			}})

	c.Insert(
		bson.M{"name": "EGI",
			"db_conf": bson.M{
				"server":   "localhost",
				"port":     27017,
				"database": "argo_EGI",
			},
			"users": []bson.M{
				bson.M{
					"name":    "Joe Complex",
					"email":   "C.Joe@egi.eu",
					"api_key": "mysecretcombination",
				},
				bson.M{
					"name":    "Josh Plain",
					"email":   "P.Josh@egi.eu",
					"api_key": "itsamysterytoyou",
				},
			}})
}

// Testing query and projection using FindAndProject method.
// During Setup of the test environment the testdb is seeded with
// two tenants ("Westeros","EGI"). We query for a tenant specific
// x-api-key and expect to get back the tenant "name" and "db_conf"
func (suite *mongoTestSuite) TestFindAndProject() {
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	query := bson.M{"users.api_key": "sansa3"}
	projection := bson.M{"_id": 0, "name": 1, "db_conf": 1}
	expected_result := []bson.M{bson.M{"name": "Westeros", "db_conf": bson.M{"server": "localhost", "port": 27017, "database": "argo_GOT"}}}
	result := []bson.M{}
	err = FindAndProject(session, suite.cfg.MongoDB.Db, "tenants", query, projection, "users.api_key", &result)
	suite.Equal(expected_result, result, "Unexpected result")
}

func (suite *mongoTestSuite) TearDownSuite() {
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

func TestMongoTestSuite(t *testing.T) {
	suite.Run(t, new(mongoTestSuite))
}
