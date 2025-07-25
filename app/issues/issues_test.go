/*
 * Copyright (c) 2015 GRNET S.A.
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

package issues

import (
	"context"
	"fmt"
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
type IssuesTestSuite struct {
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
func (suite *IssuesTestSuite) SetupTest() {

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
    db = "argotest_issues"
`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(true).PathPrefix("/api/v2/issues").Subrouter()
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
				"database": "argotest_flatendpoints_egi",
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
				"database": "argotest_flatendpoints_eudat",
				"username": "",
				"password": ""},
		},
		"users": []bson.M{
			{
				"name":    "eudat_user",
				"email":   "eudat_user@email.com",
				"roles":   []string{"viewer"},
				"api_key": "KEY2"},
		}})

	c = suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("roles")
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "issues.list_endpoints",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "issues.list_group_metrics",
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
	c = suite.cfg.MongoClient.Database(tenantDbConf1.Db).Collection("status_endpoints")
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"info":           bson.M{"Url": "http://example.foo/path/to/service"},
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream01.afroditi.gr",

		"status": "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T01:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream01.afroditi.gr",
		"info":           bson.M{"Url": "http://example.foo/path/to/service"},

		"status": "CRITICAL",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T05:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream01.afroditi.gr",
		"info":           bson.M{"Url": "http://example.foo/path/to/service"},

		"status": "WARNING",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream02.afroditi.gr",

		"status": "WARNING",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T08:47:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream02.afroditi.gr",

		"status": "WARNING",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T12:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream02.afroditi.gr",

		"status": "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream03.afroditi.gr",

		"status": "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T04:40:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream03.afroditi.gr",

		"status": "UNKNOWN",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T06:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream03.afroditi.gr",

		"status": "CRITICAL",
	})

	// seed the status detailed metric data
	c = suite.cfg.MongoClient.Database(tenantDbConf1.Db).Collection("status_metrics")
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"info":           bson.M{"Url": "http://example.foo/path/to/service"},
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream01.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T01:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream01.afroditi.gr",
		"info":           bson.M{"Url": "http://example.foo/path/to/service"},
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "CRITICAL",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T05:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream01.afroditi.gr",
		"info":           bson.M{"Url": "http://example.foo/path/to/service"},
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "WARNING",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream02.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "WARNING",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T08:47:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream02.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "WARNING",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T12:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream02.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream03.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T04:40:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream03.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "UNKNOWN",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T06:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream03.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "CRITICAL",
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
	c = suite.cfg.MongoClient.Database(tenantDbConf2.Db).Collection("status_endpoints")
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "EL-01-AUTH",
		"service":        "srv.typeA",
		"host":           "host01.eudat.gr",
		"metric":         "typeA.metric.Memory",
		"status":         "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T01:00:00Z",
		"endpoint_group": "EL-01-AUTH",
		"service":        "srv.typeA",
		"host":           "host01.eudat.gr",
		"metric":         "typeA.metric.Memory",
		"status":         "WARNING",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T05:00:00Z",
		"endpoint_group": "EL-01-AUTH",
		"service":        "srv.typeA",
		"host":           "host01.eudat.gr",
		"metric":         "typeA.metric.Memory",
		"status":         "CRITICAL",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "EL-01-AUTH",
		"service":        "srv.typeA",
		"host":           "host02.eudat.gr",
		"metric":         "typeA.metric.Memory",
		"status":         "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T01:00:00Z",
		"endpoint_group": "EL-01-AUTH",
		"service":        "srv.typeA",
		"host":           "host02.eudat.gr",
		"metric":         "typeA.metric.Memory",
		"status":         "CRITICAL",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T05:00:00Z",
		"endpoint_group": "EL-01-AUTH",
		"service":        "srv.typeA",
		"host":           "host02.eudat.gr",
		"metric":         "typeA.metric.Memory",
		"status":         "WARNING",
	})

}

func (suite *IssuesTestSuite) TestIssuesEndpoints() {

	type expReq struct {
		method string
		url    string
		code   int
		result string
		key    string
	}

	expReqs := []expReq{

		{
			method: "GET",
			url:    "/api/v2/issues/Report_A/endpoints?date=2015-05-01",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "timestamp": "2015-05-01T05:00:00Z",
   "endpoint_group": "HG-03-AUTH",
   "service": "CREAM-CE",
   "endpoint": "cream01.afroditi.gr",
   "status": "WARNING",
   "info": {
    "Url": "http://example.foo/path/to/service"
   }
  },
  {
   "timestamp": "2015-05-01T06:00:00Z",
   "endpoint_group": "HG-03-AUTH",
   "service": "CREAM-CE",
   "endpoint": "cream03.afroditi.gr",
   "status": "CRITICAL"
  }
 ]
}`,
		},
		{
			method: "GET",
			url:    "/api/v2/issues/Report_A/endpoints?date=2015-05-01&filter=CRITICAL",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "timestamp": "2015-05-01T06:00:00Z",
   "endpoint_group": "HG-03-AUTH",
   "service": "CREAM-CE",
   "endpoint": "cream03.afroditi.gr",
   "status": "CRITICAL"
  }
 ]
}`,
		},
		{
			method: "GET",
			url:    "/api/v2/issues/Report_A/endpoints?date=2015-05-01&filter=WARNING",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "timestamp": "2015-05-01T05:00:00Z",
   "endpoint_group": "HG-03-AUTH",
   "service": "CREAM-CE",
   "endpoint": "cream01.afroditi.gr",
   "status": "WARNING",
   "info": {
    "Url": "http://example.foo/path/to/service"
   }
  }
 ]
}`,
		},
		{
			method: "GET",
			url:    "/api/v2/issues/Report_A/groups/HG-03-AUTH/metrics?date=2015-05-01",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "service": "CREAM-CE",
   "hostname": "cream01.afroditi.gr",
   "metric": "emi.cream.CREAMCE-JobSubmit",
   "status": "CRITICAL",
   "info": {
    "Url": "http://example.foo/path/to/service"
   }
  },
  {
   "service": "CREAM-CE",
   "hostname": "cream01.afroditi.gr",
   "metric": "emi.cream.CREAMCE-JobSubmit",
   "status": "WARNING",
   "info": {
    "Url": "http://example.foo/path/to/service"
   }
  },
  {
   "service": "CREAM-CE",
   "hostname": "cream02.afroditi.gr",
   "metric": "emi.cream.CREAMCE-JobSubmit",
   "status": "WARNING"
  },
  {
   "service": "CREAM-CE",
   "hostname": "cream03.afroditi.gr",
   "metric": "emi.cream.CREAMCE-JobSubmit",
   "status": "CRITICAL"
  },
  {
   "service": "CREAM-CE",
   "hostname": "cream03.afroditi.gr",
   "metric": "emi.cream.CREAMCE-JobSubmit",
   "status": "UNKNOWN"
  }
 ]
}`,
		},
		{
			method: "GET",
			url:    "/api/v2/issues/Report_A/groups/HG-03-AUTH/metrics?date=2015-05-01&filter=CRITICAL",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "service": "CREAM-CE",
   "hostname": "cream01.afroditi.gr",
   "metric": "emi.cream.CREAMCE-JobSubmit",
   "status": "CRITICAL",
   "info": {
    "Url": "http://example.foo/path/to/service"
   }
  },
  {
   "service": "CREAM-CE",
   "hostname": "cream03.afroditi.gr",
   "metric": "emi.cream.CREAMCE-JobSubmit",
   "status": "CRITICAL"
  }
 ]
}`,
		},
		{
			method: "GET",
			url:    "/api/v2/issues/Report_A/endpoints?date=2015-05-01&filter=MISSING",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": []
}`,
		},
	}

	for _, expReq := range expReqs {
		request, _ := http.NewRequest(expReq.method, expReq.url, strings.NewReader(""))
		request.Header.Set("x-api-key", expReq.key)
		request.Header.Set("Accept", "application/json")

		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		suite.Equal(expReq.code, response.Code, "Incorrect HTTP response code")
		// Compare the expected and actual xml response
		if !(suite.Equal(expReq.result, response.Body.String(), "Response body mismatch")) {
			fmt.Println(response.Body.String())
		}

	}

}

func (suite *IssuesTestSuite) TestOptionsIssuesEndpoints() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/issues/Report_A/endpoints", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

}

func (suite *IssuesTestSuite) TestOptionsIssuesGroups() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/issues/Report_A/groups/test_group/metrics", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

}

// This function is actually called in of each test
// and clears the test environment.
// Mainly it's purpose is to drop the testdb
func (suite *IssuesTestSuite) TearDownTest() {

	suite.cfg.MongoClient.Database("argotest_flatendpoints").Drop(context.TODO())
	suite.cfg.MongoClient.Database("argotest_flatendpoints_eudat").Drop(context.TODO())
	suite.cfg.MongoClient.Database("argotest_flatendpoints_egi").Drop(context.TODO())
}

// This is the first function called when go test is issued
func TestSuiteIssues(t *testing.T) {
	suite.Run(t, new(IssuesTestSuite))
}
