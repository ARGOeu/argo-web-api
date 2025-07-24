/*
 * Copyright (c) 2018 GRNET S.A.
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
 * or implied, of GRNET S.A.
 *
 */

package latest

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/authentication"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/store"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/gcfg.v1"
)

// This is a util. suite struct used in tests (see pkg "testify")
type LatestTestSuite struct {
	suite.Suite
	cfg         config.Config
	router      *mux.Router
	confHandler respond.ConfHandler
}

// Setup the Test Environment
// This function runs before any test and setups the environment
// A test configuration object is instantiated using a reference
// to testdb: argo_test_details. Also here is are instantiated some expected
// xml response validation messages (authorization,crud responses).
// Also the testdb is seeded with tenants,reports,metric_profiles and status_metrics
func (suite *LatestTestSuite) SetupTest() {

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
    db = "argotest_latest"
`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(true).PathPrefix("/api/v2/latest").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)

	// Add authentication token to mongo testdb
	seedAuth := bson.M{"api_key": "S3CR3T"}
	coreCol := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("authentication")

	coreCol.InsertOne(context.TODO(), seedAuth)

	// seed a tenant to use
	c := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("tenants")
	c.InsertOne(context.TODO(), bson.M{
		"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
		"info": bson.M{
			"name":    "GUARDIANS",
			"email":   "email@something2",
			"website": "www.gotg.com",
			"created": "2015-10-20 02:08:04",
			"updated": "2015-10-20 02:08:04"},
		"db_conf": []bson.M{
			{
				"store":    "main",
				"server":   "localhost",
				"port":     27017,
				"database": "argotest_metrics_egi",
				"username": "",
				"password": ""},
		},
		"users": []bson.M{
			{
				"name":    "egi_user",
				"email":   "egi_user@email.com",
				"roles":   []string{"viewer"},
				"api_key": "KEY1"},
		}})

	c.InsertOne(context.TODO(), bson.M{
		"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50d",
		"info": bson.M{
			"name":    "AVENGERS",
			"email":   "email@something2",
			"website": "www.gotg.com",
			"created": "2015-10-20 02:08:04",
			"updated": "2015-10-20 02:08:04"},
		"db_conf": []bson.M{
			{
				"store":    "main",
				"server":   "localhost",
				"port":     27017,
				"database": "argotest_metrics_tenant2",
				"username": "",
				"password": ""},
		},
		"users": []bson.M{
			{
				"name":    "tenant2_user",
				"email":   "tenant2_user@email.com",
				"roles":   []string{"viewer"},
				"api_key": "KEY2"},
		}})
	c = suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("roles")
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "latest.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "latest.get",
			"roles":    []string{"editor", "viewer"},
		})
	// get dbconfiguration based on the tenant
	// Prepare the request object
	request, _ := http.NewRequest("GET", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY1")
	// authenticate user's api key and find corresponding tenant
	tenantDbConf1, _, _ := authentication.AuthenticateTenant(request.Header, suite.cfg)

	// Now seed the report DEFINITIONS
	c = suite.cfg.MongoClient.Database(tenantDbConf1.Db).Collection("reports")
	c.InsertOne(context.TODO(), bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"info": bson.M{
			"name":        "Report_A",
			"description": "report aaaaa",
			"created":     "2015-9-10 13:43:00",
			"updated":     "2015-10-11 13:43:00",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "NGI",
				"group": bson.M{
					"type": "SITES",
				},
			},
		},
		"profiles": []bson.M{
			{
				"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
				"type": "metric",
				"name": "profile1"},
			{
				"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e523",
				"type": "operations",
				"name": "profile2"},
			{
				"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e50q",
				"type": "aggregation",
				"name": "profile3"},
		},
		"filter_tags": []bson.M{
			{
				"name":  "name1",
				"value": "value1"},
			{
				"name":  "name2",
				"value": "value2"},
		}})

	// seed the status detailed metric data
	c = suite.cfg.MongoClient.Database(tenantDbConf1.Db).Collection("status_metrics")
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"monitoring_box":     "nagios3.hellasgrid.gr",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T00:00:00Z",
		"service":            "CREAM-CE",
		"host":               "cream01.afroditi.gr",
		"endpoint_group":     "HG-03-AUTH",
		"metric":             "emi.cream.CREAMCE-JobSubmit",
		"status":             "OK",
		"time_integer":       0,
		"previous_state":     "OK",
		"previous_timestamp": "2015-04-30T23:59:00Z",
		"summary":            "Cream status is ok",
		"message":            "Cream job submission test return value of ok",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"monitoring_box":     "nagios3.hellasgrid.gr",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T01:00:00Z",
		"service":            "CREAM-CE",
		"host":               "cream01.afroditi.gr",
		"endpoint_group":     "HG-03-AUTH",
		"metric":             "emi.cream.CREAMCE-JobSubmit",
		"status":             "CRITICAL",
		"time_integer":       10000,
		"previous_state":     "OK",
		"previous_timestamp": "2015-05-01T00:00:00Z",
		"summary":            "Cream status is CRITICAL",
		"message":            "Cream job submission test failed",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"monitoring_box":     "nagios3.hellasgrid.gr",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T05:00:00Z",
		"service":            "CREAM-CE",
		"host":               "cream01.afroditi.gr",
		"endpoint_group":     "HG-03-AUTH",
		"metric":             "emi.cream.CREAMCE-JobSubmit",
		"status":             "OK",
		"time_integer":       50000,
		"previous_state":     "CRITICAL",
		"previous_timestamp": "2015-05-01T01:00:00Z",
		"summary":            "Cream status is ok",
		"message":            "Cream job submission test return value of ok",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"monitoring_box":     "nagios3.tenant2.eu",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T23:20:00Z",
		"service":            "someService-A",
		"host":               "someservice.example.gr",
		"endpoint_group":     "EL-01-AUTH",
		"metric":             "someService-FileTransfer",
		"status":             "OK",
		"time_integer":       232000,
		"previous_state":     "OK",
		"previous_timestamp": "2015-05-01T22:20:00Z",
		"summary":            "someService status is ok",
		"message":            "someService data upload test return value of ok",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"monitoring_box":     "nagios3.tenant2.eu",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T00:00:00Z",
		"service":            "someService-A",
		"host":               "someservice.example.gr",
		"endpoint_group":     "EL-01-AUTH",
		"metric":             "someService-FileTransfer",
		"status":             "OK",
		"time_integer":       0,
		"previous_state":     "OK",
		"previous_timestamp": "2015-05-01T23:20:00Z",
		"summary":            "someService status is ok",
		"message":            "someService data upload test return value of ok",
	})

	// get dbconfiguration based on the tenant
	// Prepare the request object
	request, _ = http.NewRequest("GET", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY2")
	// authenticate user's api key and find corresponding tenant
	tenantDbConf2, _, _ := authentication.AuthenticateTenant(request.Header, suite.cfg)

	// Now seed the reports DEFINITIONS
	c = suite.cfg.MongoClient.Database(tenantDbConf2.Db).Collection("reports")
	c.InsertOne(context.TODO(), bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"info": bson.M{
			"name":        "Report_B",
			"description": "report aaaaa",
			"created":     "2015-9-10 13:43:00",
			"updated":     "2015-10-11 13:43:00",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "EUDAT_GROUPS",
				"group": bson.M{
					"type": "EUDAT_SITES",
				},
			},
		},
		"profiles": []bson.M{
			{
				"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
				"type": "metric",
				"name": "eudat.CRITICAL"},
			{
				"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e523",
				"type": "operations",
				"name": "profile2"},
			{
				"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e50q",
				"type": "aggregation",
				"name": "profile3"},
		},
		"filter_tags": []bson.M{
			{
				"name":  "name1",
				"value": "value1"},
			{
				"name":  "name2",
				"value": "value2"},
		}})

	// seed the status detailed metric data
	c = suite.cfg.MongoClient.Database(tenantDbConf2.Db).Collection("status_metrics")
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"monitoring_box":     "nagios3.tenant2.eu",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T00:00:00Z",
		"service":            "someService",
		"host":               "someservice.example.gr",
		"endpoint_group":     "EL-01-AUTH",
		"metric":             "someService-FileTransfer",
		"status":             "OK",
		"time_integer":       0,
		"previous_state":     "OK",
		"previous_timestamp": "2015-04-30T23:59:00Z",
		"summary":            "someService status is ok",
		"message":            "someService data upload test return value of ok",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"monitoring_box":     "nagios3.tenant2.eu",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T01:00:00Z",
		"service":            "someService",
		"host":               "someservice.example.gr",
		"endpoint_group":     "EL-01-AUTH",
		"metric":             "someService-FileTransfer",
		"status":             "UNKNOWN",
		"time_integer":       10000,
		"previous_state":     "OK",
		"previous_timestamp": "2015-05-01T00:00:00Z",
		"summary":            "someService status is CRITICAL",
		"message":            "someService data upload test failed",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"monitoring_box":     "nagios3.tenant2.eu",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T05:00:00Z",
		"service":            "someService",
		"host":               "someservice.example.gr",
		"endpoint_group":     "EL-01-AUTH",
		"metric":             "someService-FileTransfer",
		"status":             "CRITICAL",
		"time_integer":       50000,
		"previous_state":     "CRITICAL",
		"previous_timestamp": "2015-05-01T01:00:00Z",
		"summary":            "someService status is ok",
		"message":            "someService data upload test return value of ok",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"monitoring_box":     "nagios3.tenant2.eu",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T05:00:00Z",
		"service":            "someService",
		"host":               "someservice.example.gr",
		"endpoint_group":     "EL-01-AUTH",
		"metric":             "someService-FileTransfer",
		"status":             "MISSING",
		"time_integer":       80000,
		"previous_state":     "MISSING",
		"previous_timestamp": "2015-05-01T02:00:00Z",
		"summary":            "someService status is ok",
		"message":            "someService data upload test return value of ok",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"monitoring_box":     "nagios3.tenant2.eu",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T05:00:00Z",
		"service":            "someService",
		"host":               "someservice.example.gr",
		"endpoint_group":     "EL-01-AUTH",
		"metric":             "someService-FileTransfer",
		"status":             "WARNING",
		"time_integer":       90000,
		"previous_state":     "WARNING",
		"previous_timestamp": "2015-05-01T03:00:00Z",
		"summary":            "someService status is ok",
		"message":            "someService data upload test return value of ok",
	})

}

func (suite *LatestTestSuite) TestListLatest() {

	respJSON1 := `{
 "status": {
  "message": "application/json",
  "code": "200"
 },
 "data": {
  "metric_data": [
   {
    "endpoint_group": "HG-03-AUTH",
    "service": "CREAM-CE",
    "endpoint": "cream01.afroditi.gr",
    "metric": "emi.cream.CREAMCE-JobSubmit",
    "timestamp": "2015-05-01T05:00:00Z",
    "status": "OK",
    "summary": "Cream status is ok",
    "message": "Cream job submission test return value of ok"
   },
   {
    "endpoint_group": "HG-03-AUTH",
    "service": "CREAM-CE",
    "endpoint": "cream01.afroditi.gr",
    "metric": "emi.cream.CREAMCE-JobSubmit",
    "timestamp": "2015-05-01T01:00:00Z",
    "status": "CRITICAL",
    "summary": "Cream status is CRITICAL",
    "message": "Cream job submission test failed"
   },
   {
    "endpoint_group": "HG-03-AUTH",
    "service": "CREAM-CE",
    "endpoint": "cream01.afroditi.gr",
    "metric": "emi.cream.CREAMCE-JobSubmit",
    "timestamp": "2015-05-01T00:00:00Z",
    "status": "OK",
    "summary": "Cream status is ok",
    "message": "Cream job submission test return value of ok"
   }
  ]
 }
}`
	respJSON2 := `{
 "status": {
  "message": "application/json",
  "code": "200"
 },
 "data": {
  "metric_data": [
   {
    "endpoint_group": "EL-01-AUTH",
    "service": "someService",
    "endpoint": "someservice.example.gr",
    "metric": "someService-FileTransfer",
    "timestamp": "2015-05-01T05:00:00Z",
    "status": "WARNING",
    "summary": "someService status is ok",
    "message": "someService data upload test return value of ok"
   },
   {
    "endpoint_group": "EL-01-AUTH",
    "service": "someService",
    "endpoint": "someservice.example.gr",
    "metric": "someService-FileTransfer",
    "timestamp": "2015-05-01T05:00:00Z",
    "status": "MISSING",
    "summary": "someService status is ok",
    "message": "someService data upload test return value of ok"
   },
   {
    "endpoint_group": "EL-01-AUTH",
    "service": "someService",
    "endpoint": "someservice.example.gr",
    "metric": "someService-FileTransfer",
    "timestamp": "2015-05-01T05:00:00Z",
    "status": "CRITICAL",
    "summary": "someService status is ok",
    "message": "someService data upload test return value of ok"
   },
   {
    "endpoint_group": "EL-01-AUTH",
    "service": "someService",
    "endpoint": "someservice.example.gr",
    "metric": "someService-FileTransfer",
    "timestamp": "2015-05-01T01:00:00Z",
    "status": "UNKNOWN",
    "summary": "someService status is CRITICAL",
    "message": "someService data upload test failed"
   },
   {
    "endpoint_group": "EL-01-AUTH",
    "service": "someService",
    "endpoint": "someservice.example.gr",
    "metric": "someService-FileTransfer",
    "timestamp": "2015-05-01T00:00:00Z",
    "status": "OK",
    "summary": "someService status is ok",
    "message": "someService data upload test return value of ok"
   }
  ]
 }
}`

	respJSON3 := `{
 "status": {
  "message": "application/json",
  "code": "200"
 },
 "data": {
  "metric_data": [
   {
    "endpoint_group": "EL-01-AUTH",
    "service": "someService",
    "endpoint": "someservice.example.gr",
    "metric": "someService-FileTransfer",
    "timestamp": "2015-05-01T05:00:00Z",
    "status": "WARNING",
    "summary": "someService status is ok",
    "message": "someService data upload test return value of ok"
   }
  ]
 }
}`

	respJSON4 := `{
 "status": {
  "message": "application/json",
  "code": "200"
 },
 "data": {
  "metric_data": [
   {
    "endpoint_group": "EL-01-AUTH",
    "service": "someService",
    "endpoint": "someservice.example.gr",
    "metric": "someService-FileTransfer",
    "timestamp": "2015-05-01T05:00:00Z",
    "status": "WARNING",
    "summary": "someService status is ok",
    "message": "someService data upload test return value of ok"
   },
   {
    "endpoint_group": "EL-01-AUTH",
    "service": "someService",
    "endpoint": "someservice.example.gr",
    "metric": "someService-FileTransfer",
    "timestamp": "2015-05-01T05:00:00Z",
    "status": "MISSING",
    "summary": "someService status is ok",
    "message": "someService data upload test return value of ok"
   },
   {
    "endpoint_group": "EL-01-AUTH",
    "service": "someService",
    "endpoint": "someservice.example.gr",
    "metric": "someService-FileTransfer",
    "timestamp": "2015-05-01T05:00:00Z",
    "status": "CRITICAL",
    "summary": "someService status is ok",
    "message": "someService data upload test return value of ok"
   },
   {
    "endpoint_group": "EL-01-AUTH",
    "service": "someService",
    "endpoint": "someservice.example.gr",
    "metric": "someService-FileTransfer",
    "timestamp": "2015-05-01T01:00:00Z",
    "status": "UNKNOWN",
    "summary": "someService status is CRITICAL",
    "message": "someService data upload test failed"
   }
  ]
 }
}`

	respJSON5 := `{
 "status": {
  "message": "application/json",
  "code": "200"
 },
 "data": {
  "metric_data": [
   {
    "endpoint_group": "EL-01-AUTH",
    "service": "someService",
    "endpoint": "someservice.example.gr",
    "metric": "someService-FileTransfer",
    "timestamp": "2015-05-01T05:00:00Z",
    "status": "CRITICAL",
    "summary": "someService status is ok",
    "message": "someService data upload test return value of ok"
   }
  ]
 }
}`

	respJSON6 := `{
 "status": {
  "message": "application/json",
  "code": "200"
 },
 "data": {
  "metric_data": [
   {
    "endpoint_group": "EL-01-AUTH",
    "service": "someService",
    "endpoint": "someservice.example.gr",
    "metric": "someService-FileTransfer",
    "timestamp": "2015-05-01T00:00:00Z",
    "status": "OK",
    "summary": "someService status is ok",
    "message": "someService data upload test return value of ok"
   }
  ]
 }
}`

	respJSON7 := `{
 "status": {
  "message": "application/json",
  "code": "200"
 },
 "data": {
  "metric_data": [
   {
    "endpoint_group": "HG-03-AUTH",
    "service": "CREAM-CE",
    "endpoint": "cream01.afroditi.gr",
    "metric": "emi.cream.CREAMCE-JobSubmit",
    "timestamp": "2015-05-01T05:00:00Z",
    "status": "OK",
    "summary": "Cream status is ok",
    "message": "Cream job submission test return value of ok"
   }
  ]
 }
}`

	respJSON8 := `{
 "status": {
  "message": "application/json",
  "code": "200"
 },
 "data": {
  "metric_data": [
   {
    "endpoint_group": "EL-01-AUTH",
    "service": "someService-A",
    "endpoint": "someservice.example.gr",
    "metric": "someService-FileTransfer",
    "timestamp": "2015-05-01T23:20:00Z",
    "status": "OK",
    "summary": "someService status is ok",
    "message": "someService data upload test return value of ok"
   },
   {
    "endpoint_group": "HG-03-AUTH",
    "service": "CREAM-CE",
    "endpoint": "cream01.afroditi.gr",
    "metric": "emi.cream.CREAMCE-JobSubmit",
    "timestamp": "2015-05-01T05:00:00Z",
    "status": "OK",
    "summary": "Cream status is ok",
    "message": "Cream job submission test return value of ok"
   }
  ]
 }
}`

	respJSON9 := `{
 "status": {
  "message": "application/json",
  "code": "200"
 },
 "data": {
  "metric_data": [
   {
    "endpoint_group": "EL-01-AUTH",
    "service": "someService-A",
    "endpoint": "someservice.example.gr",
    "metric": "someService-FileTransfer",
    "timestamp": "2015-05-01T23:20:00Z",
    "status": "OK",
    "summary": "someService status is ok",
    "message": "someService data upload test return value of ok"
   },
   {
    "endpoint_group": "HG-03-AUTH",
    "service": "CREAM-CE",
    "endpoint": "cream01.afroditi.gr",
    "metric": "emi.cream.CREAMCE-JobSubmit",
    "timestamp": "2015-05-01T05:00:00Z",
    "status": "OK",
    "summary": "Cream status is ok",
    "message": "Cream job submission test return value of ok"
   }
  ]
 }
}`

	respJSON10 := `{
 "status": {
  "message": "application/json",
  "code": "200"
 },
 "data": {
  "metric_data": [
   {
    "endpoint_group": "EL-01-AUTH",
    "service": "someService",
    "endpoint": "someservice.example.gr",
    "metric": "someService-FileTransfer",
    "timestamp": "2015-05-01T05:00:00Z",
    "status": "WARNING",
    "summary": "someService status is ok",
    "message": "someService data upload test return value of ok"
   }
  ]
 }
}`

	respUnauthorized := `{
 "status": {
  "message": "Unauthorized",
  "code": "401",
  "details": "You need to provide a correct authentication token using the header 'x-api-key'"
 }
}`

	fullurl1 := "/api/v2/latest/Report_A/SITES/HG-03-AUTH" +
		"?date=2015-05-01T00:00:00Z&strict=false"

	fullurl2 := "/api/v2/latest/Report_B/EUDAT_SITES/EL-01-AUTH" +
		"?date=2015-05-01T00:00:00Z&strict=false"

	fullurl3 := "/api/v2/latest/Report_B/EUDAT_SITES/EL-01-AUTH" +
		"?date=2015-05-01T00:00:00Z&limit=1&strict=false"

	fullurl4 := "/api/v2/latest/Report_B/EUDAT_SITES/EL-01-AUTH" +
		"?date=2015-05-01T00:00:00Z&filter=non-ok"

	fullurl5 := "/api/v2/latest/Report_B/EUDAT_SITES/EL-01-AUTH" +
		"?date=2015-05-01T00:00:00Z&filter=critical"

	fullurl6 := "/api/v2/latest/Report_B/EUDAT_SITES/EL-01-AUTH" +
		"?date=2015-05-01T00:00:00Z&filter=ok"

	fullurl7 := "/api/v2/latest/Report_A/SITES/HG-03-AUTH" +
		"?date=2015-05-01T00:00:00Z&strict=true"

	fullurl8 := "/api/v2/latest/Report_A/SITES" +
		"?date=2015-05-01T00:00:00Z&strict=true"

	fullurl9 := "/api/v2/latest/Report_A/SITES" +
		"?date=2015-05-01T00:00:00Z&strict=true&limit=2"

	fullurl10 := "/api/v2/latest/Report_B/EUDAT_SITES/EL-01-AUTH" +
		"?date=2015-05-01T00:00:00Z&filter=non-ok&strict=true"

	// 1. EGI JSON REQUEST
	// init the response placeholder
	response := httptest.NewRecorder()
	// Prepare the request object for second tenant
	request, _ := http.NewRequest("GET", fullurl1, strings.NewReader(""))
	// add json accept header
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY1")
	// Serve the http request
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respJSON1, response.Body.String(), "Response body mismatch")

	// 2. EUDAT JSON REQUEST
	// init the response placeholder
	response = httptest.NewRecorder()
	// Prepare the request object for second tenant
	request, _ = http.NewRequest("GET", fullurl2, strings.NewReader(""))
	// add json accept header
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY2")
	// Serve the http request
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respJSON2, response.Body.String(), "Response body mismatch")

	// 3. EUDAT limit = 1
	// init the response placeholder
	response = httptest.NewRecorder()
	// Prepare the request object for second tenant
	request, _ = http.NewRequest("GET", fullurl3, strings.NewReader(""))
	// add json accept header
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY2")
	// Serve the http request
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respJSON3, response.Body.String(), "Response body mismatch")

	// 4. EUDAT non-ok
	// init the response placeholder
	response = httptest.NewRecorder()
	// Prepare the request object for second tenant
	request, _ = http.NewRequest("GET", fullurl4, strings.NewReader(""))
	// add json accept header
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY2")
	// Serve the http request
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respJSON4, response.Body.String(), "Response body mismatch")

	// 5. EUDAT non-ok
	// init the response placeholder
	response = httptest.NewRecorder()
	// Prepare the request object for second tenant
	request, _ = http.NewRequest("GET", fullurl5, strings.NewReader(""))
	// add json accept header
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY2")
	// Serve the http request
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respJSON5, response.Body.String(), "Response body mismatch")

	// 6. EUDAT non-ok
	// init the response placeholder
	response = httptest.NewRecorder()
	// Prepare the request object for second tenant
	request, _ = http.NewRequest("GET", fullurl6, strings.NewReader(""))
	// add json accept header
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY2")
	// Serve the http request
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respJSON6, response.Body.String(), "Response body mismatch")

	// 6b. WRONG KEY REQUEST
	// init the response placeholder
	response = httptest.NewRecorder()
	// Prepare the request object for second tenant
	request, _ = http.NewRequest("GET", fullurl2, strings.NewReader(""))
	// add json accept header
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEYISWRONG")
	// Serve the http request
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(401, response.Code, "Response code mismatch")
	// Compare the expected and actual xml response
	suite.Equal(respUnauthorized, response.Body.String(), "Response body mismatch")

	// 7. EGI JSON REQUEST - strict mode
	// init the response placeholder
	response = httptest.NewRecorder()
	// Prepare the request object for second tenant
	request, _ = http.NewRequest("GET", fullurl7, strings.NewReader(""))
	// add json accept header
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY1")
	// Serve the http request
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respJSON7, response.Body.String(), "Response body mismatch")

	// 8. EGI JSON REQUEST - strict mode - all sites
	// init the response placeholder
	response = httptest.NewRecorder()
	// Prepare the request object for second tenant
	request, _ = http.NewRequest("GET", fullurl8, strings.NewReader(""))
	// add json accept header
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY1")
	// Serve the http request
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respJSON8, response.Body.String(), "Response body mismatch")

	// 9. EGI JSON REQUEST - strict mode - all sites but with limit (strict honors limits)
	// init the response placeholder
	response = httptest.NewRecorder()
	// Prepare the request object for second tenant
	request, _ = http.NewRequest("GET", fullurl9, strings.NewReader(""))
	// add json accept header
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY1")
	// Serve the http request
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respJSON9, response.Body.String(), "Response body mismatch")

	// 10. EGI JSON REQUEST - strict mode - honor non-ok values only
	// init the response placeholder
	response = httptest.NewRecorder()
	// Prepare the request object for second tenant
	request, _ = http.NewRequest("GET", fullurl10, strings.NewReader(""))
	// add json accept header
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY2")
	// Serve the http request
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respJSON10, response.Body.String(), "Response body mismatch")

}

// This function is actually called in the end of all tests
// and clears the test environment.
// Mainly it's purpose is to drop the testdb
func (suite *LatestTestSuite) TearDownTest() {

	suite.cfg.MongoClient.Database("argotest_metrics").Drop(context.TODO())
	suite.cfg.MongoClient.Database("argotest_metrics_tenant2").Drop(context.TODO())
	suite.cfg.MongoClient.Database("argotest_metrics_egi").Drop(context.TODO())
}

// This is the first function called when go test is issued
func TestSuiteLatest(t *testing.T) {
	suite.Run(t, new(LatestTestSuite))
}
