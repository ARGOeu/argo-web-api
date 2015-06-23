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

package ngiAvailability

import (
	"code.google.com/p/gcfg"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"
	"github.com/stretchr/testify/suite"
	//"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"testing"
)

type EgAvailabilityTestSuite struct {
	suite.Suite
	cfg              config.Config
	tenantDbConf     config.MongoConfig
	tenantPassword   string
	tenantUsername   string
	tenantStorename  string
	clientKey        string
	respUnauthorized string
}

// Setup the Test Environment
// This function runs before any test and setups the environment
// A test configuration object is instantiated using a reference
// to testdb: ARGO_core_test_* .

func (suite *EgAvailabilityTestSuite) SetupTest() {

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
    db = "ARGO_core_test_egavailability"
`

	suite.respUnauthorized = "Unauthorized"
	suite.tenantDbConf.Db = "ARGO_northern_test_egavailability"
	suite.tenantDbConf.Password = "h4shp4ss"
	suite.tenantDbConf.Username = "dbuser"
	suite.tenantDbConf.Store = "ar"
	suite.clientKey = "secretkey"

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	// Connect to mongo testdb
	session, err := mongo.OpenSession(suite.cfg.MongoDB)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Seed database with tenants and test credentials
	c := session.DB(suite.cfg.MongoDB.Db).C("tenants")
	c.Insert(
		bson.M{"name": "Western",
			"db_conf": []bson.M{

				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Western1",
				},
				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Western2",
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "John Snow",
					"email":   "J.Snow@foo.bar",
					"api_key": "wh1t3_w@lk3rs",
				},
				bson.M{
					"name":    "King Joffrey",
					"email":   "g0dk1ng@foo.bar",
					"api_key": "sansa <3",
				},
			}})
	c.Insert(
		bson.M{"name": "Northern",
			"db_conf": []bson.M{

				bson.M{
					// "store":    "ar",
					"server":   "localhost",
					"port":     27017,
					"database": suite.tenantDbConf.Db,
					"username": suite.tenantDbConf.Username,
					"password": suite.tenantDbConf.Password,
				},
				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_northern_test_db",
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "Joe Complex",
					"email":   "C.Joe@foo.bar",
					"api_key": suite.clientKey,
				},
				bson.M{
					"name":    "Josh Plain",
					"email":   "P.Josh@foo.bar",
					"api_key": "testsecretapikey",
				},
			}})

	// Open DB session
	c = session.DB(suite.tenantDbConf.Db).C("endpoint_group_ar")

	// Insert seed data
	c.Insert(
		bson.M{"job": "Northern_job1",
			"date":         "20150623",
			"name":         "ST01",
			"type":         "SITE",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 66.7,
			"reliability":  54.6,
			"weights": []bson.M{

				bson.M{
					"name":  "bench",
					"value": 56644,
				},
			},
			"tags": []bson.M{
				bson.M{
					"name":  "",
					"value": "",
				},
			},
		})
	c.Insert(
		bson.M{"job": "Northern_job2",
			"date":         "20150623",
			"name":         "ST02",
			"type":         "SITE",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 100,
			"reliability":  100,
			"weights": []bson.M{

				bson.M{
					"name":  "bench",
					"value": 56644,
				},
			},
			"tags": []bson.M{

				bson.M{
					"name":  "foo2",
					"value": "Y",
				},
				bson.M{
					"name":  "bar3",
					"value": "N",
				},
			},
		})
}

func (suite *EgAvailabilityTestSuite) TestReadGroupAr() {
	// Open a session to mongo
	session, err := mongo.OpenSession(suite.cfg.MongoDB)
	if err != nil {
		panic(err)
	}
	defer mongo.CloseSession(session)

}

func (suite *EgAvailabilityTestSuite) TestReadGroupArDaily() {
	// Open a session to mongo
	session, err := mongo.OpenSession(suite.cfg.MongoDB)
	if err != nil {
		panic(err)
	}
	defer mongo.CloseSession(session)

	filter := prepareFilter(input)

	query := []bson.M{
		{"$match": filter},
		{"$project": bson.M{"dt": 1, "a": 1, "r": 1, "ap": 1, "n": 1, "hs": bson.M{"$add": list{"$hs", 1}}}},
		{"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "n": "$n", "ap": "$ap"},
			"a": bson.M{"$sum": bson.M{"$multiply": list{"$a", "$hs"}}}, "r": bson.M{"$sum": bson.M{"$multiply": list{"$r", "$hs"}}}, "hs": bson.M{"$sum": "$hs"}}},
		{"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ap": "$_id.ap", "a": bson.M{"$divide": list{"$a", "$hs"}},
			"r": bson.M{"$divide": list{"$r", "$hs"}}}},
		{"$sort": bson.D{{"ap", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}

	//query := []bson.M{{"$match": q}, {"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "n": "$n", "ns": "$ns", "p": "$p"}, "a": bson.M{"$sum": bson.M{"$multiply": list{"$a", "$hs"}}}, 		"r": bson.M{"$sum": bson.M{"$multiply": list{"$r", "$hs"}}}, "hs": bson.M{"$sum": "$hs"}}}, {"$match": bson.M{"hs": bson.M{"$gt": 0}}}, {"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ns": "$_id.ns", "p": 		"$_id.p", "a": bson.M{"$divide": list{"$a", "$hs"}}, "r": bson.M{"$divide": list{"$r", "$hs"}}}}, {"$sort": bson.D{{"p", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}

	return query

}

// This function is actually called in the end of all tests
// and clears the test environment.
// Mainly it's purpose is to drop the testdb
// func (suite *EgAvailabilityTestSuite) TearDownTest() {
//
// 	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
// 	if err != nil {
// 		panic(err)
// 	}
// 	session.DB(suite.tenantDbConf.Db).DropDatabase()
// 	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
// }

// This is the first function called when go test is issued
func TestEgAvailabilityTestSuite(t *testing.T) {
	suite.Run(t, new(EgAvailabilityTestSuite))
}
