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
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/store"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/gcfg.v1"
)

// AuthenticationProfileTestSuite is a utility suite struct used in tests
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

// SetupTest will bootstrap and provide the testing environment
func (suite *AuthenticationProfileTestSuite) SetupTest() {

	log.SetOutput(io.Discard)

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
    db = "argo_core_test_authenticate"
    `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.respUnauthorized = "Unauthorized"
	suite.tenantdb = "argo_egi_AR_data"
	suite.clientkey = "mysecretcombination"
	suite.tenantpassword = "h4shp4ss"
	suite.tenantusername = "johndoe"
	suite.tenantstorename = "ar"

	// seed mongo
	client := store.GetMongoClient(suite.cfg.MongoDB)

	suite.cfg.MongoClient = client
	// Seed database with tenants
	//TODO: move tests to
	c := client.Database(suite.cfg.MongoDB.Db).Collection("tenants")
	c.InsertOne(context.TODO(),
		bson.M{"name": "Westeros",
			"db_conf": []bson.M{
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Westeros1",
				},
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Westeros2",
				},
			},
			"users": []bson.M{
				{
					"name":    "John Snow",
					"email":   "J.Snow@brothers.wall",
					"api_key": "wh1t3_w@lk3rs",
				},
				{
					"name":    "King Joffrey",
					"email":   "g0dk1ng@kingslanding.gov",
					"api_key": "sansa <3",
				},
			}})
	c.InsertOne(context.TODO(),
		bson.M{"name": "EGI",
			"db_conf": []bson.M{
				{
					"store":    suite.tenantstorename,
					"server":   "localhost",
					"port":     27017,
					"database": suite.tenantdb,
					"username": suite.tenantusername,
					"password": suite.tenantpassword,
				},
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_egi_metric_data",
				},
			},
			"users": []bson.M{
				{
					"name":    "Joe Complex",
					"email":   "C.Joe@egi.eu",
					"api_key": suite.clientkey,
				},
				{
					"name":    "Josh Plain",
					"email":   "P.Josh@egi.eu",
					"api_key": "itsamysterytoyou",
				},
			}})
	c = client.Database(suite.cfg.MongoDB.Db).Collection("authentication")
	c.InsertOne(context.TODO(),
		bson.M{
			"name":    "Igano Kabamaru",
			"email":   "igano@kabamaru.io",
			"api_key": "makaronada",
		},
	)
	c.InsertOne(context.TODO(),
		bson.M{
			"name":    "Optimus Prime",
			"email":   "prime@autobots.com",
			"api_key": "megatron_sucks",
		},
	)
}

// TestAdminAuthentication performs unit tests against the AuthenticateAdmin function
func (suite *AuthenticationProfileTestSuite) TestAdminAuthentication() {

	request, _ := http.NewRequest("GET", "", strings.NewReader(""))
	request.Header.Set("x-api-key", "megatron_is_a_fool")
	suite.Equal(AuthenticateAdmin(request.Header, suite.cfg), false, "authetication problem")

	request.Header.Set("x-api-key", "makaronada")
	suite.Equal(AuthenticateAdmin(request.Header, suite.cfg), true, "authetication problem")

}

// TestTenantAuthentication performs unit tests against the AuthenticateTenant function
func (suite *AuthenticationProfileTestSuite) TestTenantAuthentication() {

	request, _ := http.NewRequest("GET", "", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)

	tenantdbconfig, _, err := AuthenticateTenant(request.Header, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(nil, err, "authetication problem")

	//Check the data is conveyed correctly
	suite.Regexp(tenantdbconfig.Db, suite.tenantdb, "Database mismatch")
	suite.Regexp(tenantdbconfig.Username, suite.tenantusername, "Username mismatch")
	suite.Regexp(tenantdbconfig.Password, suite.tenantpassword, "Password mismatch")
	suite.Regexp(tenantdbconfig.Store, suite.tenantstorename, "Store db mismatch")
}

// TearDownTest to tear down every test
func (suite *AuthenticationProfileTestSuite) TearDownTest() {

	client := store.GetMongoClient(suite.cfg.MongoDB)
	defer client.Disconnect(context.TODO())

	client.Database(suite.tenantdb).Drop(context.TODO())
	client.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())

}

func TestSuiteAuthentication(t *testing.T) {
	suite.Run(t, new(AuthenticationProfileTestSuite))
}
