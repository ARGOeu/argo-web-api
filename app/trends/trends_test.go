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

package trends

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
	"gopkg.in/gcfg.v1"
	"gopkg.in/mgo.v2/bson"
)

// This is a util. suite struct used in tests (see pkg "testify")
type TrendsTestSuite struct {
	suite.Suite
	cfg          config.Config
	router       *mux.Router
	confHandler  respond.ConfHandler
	tenantDbConf config.MongoConfig
}

// Setup the Test Environment
// This function runs before any test and setups the environment
// A test configuration object is instantiated using a reference
// to testdb: argo_test_details. Also here is are instantiated some expected
// xml response validation messages (authorization,crud responses).
// Also the testdb is seeded with tenants,reports,metric_profiles and status_metrics
func (suite *TrendsTestSuite) SetupTest() {

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
    db = "argotest_trends"
`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(true).PathPrefix("/api/v2/trends").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)

	// Connect to mongo testdb
	authCol := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("authentication")

	// Add authentication token to mongo testdb
	seedAuth := bson.M{"api_key": "S3CR3T"}
	authCol.InsertOne(context.TODO(), seedAuth)

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
				"database": "argotest_trends_tenant",
				"username": "",
				"password": ""},
		},
		"users": []bson.M{
			{
				"name":    "tenant_user",
				"email":   "tenant_user@email.com",
				"roles":   []string{"viewer"},
				"api_key": "KEY1"},
		}})

	c = suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("roles")
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "trends.flapping_metrics",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "trends.flapping_metrics_tags",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "trends.flapping_endpoints",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "trends.flapping_services",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "trends.flapping_endpoint_groups",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "trends.status_metrics",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "trends.status_endpoints",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "trends.status_services",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "trends.status_metrics_tags",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "trends.status_endpoint_groups",
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
	t1conf, _, _ := authentication.AuthenticateTenant(request.Header, suite.cfg)

	// Now seed the report DEFINITIONS
	c = suite.cfg.MongoClient.Database(t1conf.Db).Collection("reports")
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

	// seed the status detailed trends for endpoint group data
	c = suite.cfg.MongoClient.Database(t1conf.Db).Collection("status_trends_groups")
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"status":   "CRITICAL",
		"duration": 55,
	})

	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-B",
		"status":   "UNKNOWN",
		"duration": 12,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-XB",
		"status":   "CRITICAL",
		"duration": 25,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"status":   "WARNING",
		"duration": 55,
	})

	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-C",
		"status":   "WARNING",
		"duration": 12,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-B",
		"status":   "UNKNOWN",
		"duration": 5,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"status":   "UNKNOWN",
		"duration": 45,
	})

	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-C",
		"status":   "WARNING",
		"duration": 8,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-B",
		"status":   "CRITICAL",
		"duration": 7,
	})

	// seed the status detailed trends service data
	c = suite.cfg.MongoClient.Database(t1conf.Db).Collection("status_trends_services")
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"service":  "service-A",
		"status":   "CRITICAL",
		"duration": 55,
	})

	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"service":  "service-B",
		"status":   "UNKNOWN",
		"duration": 12,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-XB",
		"service":  "service-XA",
		"status":   "CRITICAL",
		"duration": 25,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"service":  "service-A",
		"status":   "WARNING",
		"duration": 55,
	})

	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"service":  "service-B",
		"status":   "WARNING",
		"duration": 12,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-B",
		"service":  "service-A",
		"status":   "UNKNOWN",
		"duration": 5,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"service":  "service-A",
		"status":   "UNKNOWN",
		"duration": 45,
	})

	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"service":  "service-B",
		"status":   "WARNING",
		"duration": 8,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-B",
		"service":  "service-A",
		"status":   "CRITICAL",
		"duration": 7,
	})

	// seed the status detailed trends endpoint data
	c = suite.cfg.MongoClient.Database(t1conf.Db).Collection("status_trends_endpoints")
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"status":   "CRITICAL",
		"duration": 55,
	})

	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"service":  "service-B",
		"endpoint": "hostb.example.foo",

		"status":   "UNKNOWN",
		"duration": 12,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-XB",
		"service":  "service-XA",
		"endpoint": "hosta.examplex2.foo",

		"status":   "CRITICAL",
		"duration": 25,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",

		"status":   "WARNING",
		"duration": 55,
	})

	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"service":  "service-B",
		"endpoint": "hostb.example.foo",
		"status":   "WARNING",
		"duration": 12,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-B",
		"service":  "service-A",
		"endpoint": "hosta.example2.foo",
		"status":   "UNKNOWN",
		"duration": 5,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"status":   "UNKNOWN",
		"duration": 45,
	})

	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"service":  "service-B",
		"endpoint": "hostb.example.foo",
		"status":   "WARNING",
		"duration": 8,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-B",
		"service":  "service-A",
		"endpoint": "hosta.example2.foo",
		"status":   "CRITICAL",
		"duration": 7,
	})

	// seed the status detailed trends metric data
	c = suite.cfg.MongoClient.Database(t1conf.Db).Collection("status_trends_metrics")
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"metric":   "check-1",
		"status":   "CRITICAL",
		"trends":   55,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"metric":   "check-2",
		"status":   "WARNING",
		"trends":   40,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"service":  "service-B",
		"endpoint": "hostb.example.foo",
		"metric":   "web-check",
		"status":   "UNKNOWN",
		"trends":   12,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-XB",
		"service":  "service-XA",
		"endpoint": "hosta.examplex2.foo",
		"metric":   "web-check",
		"status":   "CRITICAL",
		"trends":   25,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"metric":   "check-1",
		"status":   "WARNING",
		"trends":   55,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"metric":   "check-2",
		"status":   "CRITICAL",
		"trends":   40,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"service":  "service-B",
		"endpoint": "hostb.example.foo",
		"metric":   "web-check",
		"status":   "WARNING",
		"trends":   12,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-B",
		"service":  "service-A",
		"endpoint": "hosta.example2.foo",
		"metric":   "web-check",
		"status":   "UNKNOWN",
		"trends":   5,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"metric":   "check-1",
		"status":   "UNKNOWN",
		"trends":   45,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"metric":   "check-2",
		"status":   "UNKNOWN",
		"trends":   32,
		"tags":     []string{"STORAGE"},
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"service":  "service-B",
		"endpoint": "hostb.example.foo",
		"metric":   "web-check",
		"status":   "WARNING",
		"trends":   8,
		"tags":     []string{"NETWORK", "HTTP"},
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-B",
		"service":  "service-A",
		"endpoint": "hosta.example2.foo",
		"metric":   "web-check",
		"status":   "CRITICAL",
		"trends":   7,
	})

	// seed the status detailed trends metric data
	c = suite.cfg.MongoClient.Database(t1conf.Db).Collection("flipflop_trends_metrics")
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"metric":   "check-1",
		"flipflop": 55,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"metric":   "check-2",
		"flipflop": 40,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"service":  "service-B",
		"endpoint": "hostb.example.foo",
		"metric":   "web-check",
		"flipflop": 12,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-XB",
		"service":  "service-XA",
		"endpoint": "hosta.examplex2.foo",
		"metric":   "web-check",
		"flipflop": 25,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"metric":   "check-1",
		"flipflop": 55,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"metric":   "check-2",
		"flipflop": 40,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"service":  "service-B",
		"endpoint": "hostb.example.foo",
		"metric":   "web-check",
		"flipflop": 12,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-B",
		"service":  "service-A",
		"endpoint": "hosta.example2.foo",
		"metric":   "web-check",
		"flipflop": 5,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"metric":   "check-1",
		"flipflop": 45,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"metric":   "check-2",
		"tags":     []string{"MEMORY"},
		"flipflop": 32,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"service":  "service-B",
		"endpoint": "hostb.example.foo",
		"metric":   "web-check",
		"tags":     []string{"NETWORK", "HTTP"},
		"flipflop": 8,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-B",
		"service":  "service-A",
		"endpoint": "hosta.example2.foo",
		"metric":   "web-check",
		"tags":     []string{"NETWORK", "HTTP"},
		"flipflop": 7,
	})

	// seed the status detailed trends endpoint data
	c = suite.cfg.MongoClient.Database(t1conf.Db).Collection("flipflop_trends_endpoints")
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"flipflop": 25,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"service":  "service-B",
		"endpoint": "hostb.example.foo",
		"flipflop": 2,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-XB",
		"service":  "service-XA",
		"endpoint": "hosta.exampleX2.foo",
		"flipflop": 35,
	})

	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"flipflop": 55,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"service":  "service-B",
		"endpoint": "hostb.example.foo",
		"flipflop": 12,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-B",
		"service":  "service-A",
		"endpoint": "hosta.example2.foo",
		"flipflop": 5,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"service":  "service-A",
		"endpoint": "hosta.example.foo",
		"flipflop": 48,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"service":  "service-B",
		"endpoint": "hostb.example.foo",
		"flipflop": 7,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-B",
		"service":  "service-A",
		"endpoint": "hosta.example2.foo",
		"flipflop": 3,
	})

	// seed the status detailed trends service data
	c = suite.cfg.MongoClient.Database(t1conf.Db).Collection("flipflop_trends_services")
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"service":  "service-A",
		"flipflop": 25,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"service":  "service-B",
		"flipflop": 16,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-XB",
		"service":  "service-XA",
		"flipflop": 3,
	})

	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"service":  "service-A",
		"flipflop": 55,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"service":  "service-B",
		"flipflop": 12,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-B",
		"service":  "service-A",
		"flipflop": 5,
	})

	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"service":  "service-A",
		"flipflop": 43,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"service":  "service-B",
		"flipflop": 11,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-B",
		"service":  "service-A",
		"flipflop": 4,
	})

	// seed the status detailed trends group data
	c = suite.cfg.MongoClient.Database(t1conf.Db).Collection("flipflop_trends_endpoint_groups")
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-A",
		"flipflop": 35,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150401,
		"group":    "SITE-XB",
		"flipflop": 3,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-A",
		"flipflop": 55,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150501,
		"group":    "SITE-B",
		"flipflop": 5,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-A",
		"flipflop": 11,
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":   "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date":     20150502,
		"group":    "SITE-B",
		"flipflop": 4,
	})

}

func (suite *TrendsTestSuite) TestTrends() {

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
			url:    "/api/v2/trends/Report_A/flapping/metrics?date=2015-05-01",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "endpoint_group": "SITE-A",
   "service": "service-A",
   "endpoint": "hosta.example.foo",
   "metric": "check-1",
   "flapping": 55
  },
  {
   "endpoint_group": "SITE-A",
   "service": "service-A",
   "endpoint": "hosta.example.foo",
   "metric": "check-2",
   "flapping": 40
  },
  {
   "endpoint_group": "SITE-A",
   "service": "service-B",
   "endpoint": "hostb.example.foo",
   "metric": "web-check",
   "flapping": 12
  },
  {
   "endpoint_group": "SITE-B",
   "service": "service-A",
   "endpoint": "hosta.example2.foo",
   "metric": "web-check",
   "flapping": 5
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/endpoints?date=2015-05-01",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "endpoint_group": "SITE-A",
   "service": "service-A",
   "endpoint": "hosta.example.foo",
   "flapping": 55
  },
  {
   "endpoint_group": "SITE-A",
   "service": "service-B",
   "endpoint": "hostb.example.foo",
   "flapping": 12
  },
  {
   "endpoint_group": "SITE-B",
   "service": "service-A",
   "endpoint": "hosta.example2.foo",
   "flapping": 5
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/services?date=2015-05-01",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "endpoint_group": "SITE-A",
   "service": "service-A",
   "flapping": 55
  },
  {
   "endpoint_group": "SITE-A",
   "service": "service-B",
   "flapping": 12
  },
  {
   "endpoint_group": "SITE-B",
   "service": "service-A",
   "flapping": 5
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/groups?date=2015-05-01",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "endpoint_group": "SITE-A",
   "flapping": 55
  },
  {
   "endpoint_group": "SITE-B",
   "flapping": 5
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/metrics?start_date=2015-04-01&end_date=2015-05-02&granularity=monthly",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-1",
     "flapping": 55
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "flapping": 40
    },
    {
     "endpoint_group": "SITE-XB",
     "service": "service-XA",
     "endpoint": "hosta.examplex2.foo",
     "metric": "web-check",
     "flapping": 25
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "metric": "web-check",
     "flapping": 12
    }
   ]
  },
  {
   "date": "2015-05",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-1",
     "flapping": 100
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "flapping": 72
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "metric": "web-check",
     "flapping": 20
    },
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "metric": "web-check",
     "flapping": 12
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/metrics/tags?start_date=2015-04-01&end_date=2015-05-02&granularity=monthly",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-05",
   "tag": "HTTP",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "metric": "web-check",
     "flapping": 8
    },
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "metric": "web-check",
     "flapping": 7
    }
   ]
  },
  {
   "date": "2015-05",
   "tag": "MEMORY",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "flapping": 32
    }
   ]
  },
  {
   "date": "2015-05",
   "tag": "NETWORK",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "metric": "web-check",
     "flapping": 8
    },
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "metric": "web-check",
     "flapping": 7
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/metrics?start_date=2015-04-01&end_date=2015-05-02&top=3&granularity=monthly",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-1",
     "flapping": 55
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "flapping": 40
    },
    {
     "endpoint_group": "SITE-XB",
     "service": "service-XA",
     "endpoint": "hosta.examplex2.foo",
     "metric": "web-check",
     "flapping": 25
    }
   ]
  },
  {
   "date": "2015-05",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-1",
     "flapping": 100
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "flapping": 72
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "metric": "web-check",
     "flapping": 20
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/metrics?start_date=2015-05-01&end_date=2015-05-02&top=3",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "endpoint_group": "SITE-A",
   "service": "service-A",
   "endpoint": "hosta.example.foo",
   "metric": "check-1",
   "flapping": 100
  },
  {
   "endpoint_group": "SITE-A",
   "service": "service-A",
   "endpoint": "hosta.example.foo",
   "metric": "check-2",
   "flapping": 72
  },
  {
   "endpoint_group": "SITE-A",
   "service": "service-B",
   "endpoint": "hostb.example.foo",
   "metric": "web-check",
   "flapping": 20
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/metrics/tags?start_date=2015-05-01&end_date=2015-05-02&top=3",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "tag": "HTTP",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "metric": "web-check",
     "flapping": 8
    },
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "metric": "web-check",
     "flapping": 7
    }
   ]
  },
  {
   "tag": "MEMORY",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "flapping": 32
    }
   ]
  },
  {
   "tag": "NETWORK",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "metric": "web-check",
     "flapping": 8
    },
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "metric": "web-check",
     "flapping": 7
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/metrics?start_date=2015-05-01&end_date=2015-05-02",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "endpoint_group": "SITE-A",
   "service": "service-A",
   "endpoint": "hosta.example.foo",
   "metric": "check-1",
   "flapping": 100
  },
  {
   "endpoint_group": "SITE-A",
   "service": "service-A",
   "endpoint": "hosta.example.foo",
   "metric": "check-2",
   "flapping": 72
  },
  {
   "endpoint_group": "SITE-A",
   "service": "service-B",
   "endpoint": "hostb.example.foo",
   "metric": "web-check",
   "flapping": 20
  },
  {
   "endpoint_group": "SITE-B",
   "service": "service-A",
   "endpoint": "hosta.example2.foo",
   "metric": "web-check",
   "flapping": 12
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/endpoints?start_date=2015-05-01&end_date=2015-05-02",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "endpoint_group": "SITE-A",
   "service": "service-A",
   "endpoint": "hosta.example.foo",
   "flapping": 103
  },
  {
   "endpoint_group": "SITE-A",
   "service": "service-B",
   "endpoint": "hostb.example.foo",
   "flapping": 19
  },
  {
   "endpoint_group": "SITE-B",
   "service": "service-A",
   "endpoint": "hosta.example2.foo",
   "flapping": 8
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/endpoints?start_date=2015-04-01&end_date=2015-05-02&granularity=monthly",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "top": [
    {
     "endpoint_group": "SITE-XB",
     "service": "service-XA",
     "endpoint": "hosta.exampleX2.foo",
     "flapping": 35
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "flapping": 25
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "flapping": 2
    }
   ]
  },
  {
   "date": "2015-05",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "flapping": 103
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "flapping": 19
    },
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "flapping": 8
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/endpoints?start_date=2015-04-01&end_date=2015-05-02&top=2&granularity=monthly",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "top": [
    {
     "endpoint_group": "SITE-XB",
     "service": "service-XA",
     "endpoint": "hosta.exampleX2.foo",
     "flapping": 35
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "flapping": 25
    }
   ]
  },
  {
   "date": "2015-05",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "flapping": 103
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "flapping": 19
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/endpoints?start_date=2015-05-01&end_date=2015-05-02&top=2",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "endpoint_group": "SITE-A",
   "service": "service-A",
   "endpoint": "hosta.example.foo",
   "flapping": 103
  },
  {
   "endpoint_group": "SITE-A",
   "service": "service-B",
   "endpoint": "hostb.example.foo",
   "flapping": 19
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/services?start_date=2015-05-01&end_date=2015-05-02",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "endpoint_group": "SITE-A",
   "service": "service-A",
   "flapping": 98
  },
  {
   "endpoint_group": "SITE-A",
   "service": "service-B",
   "flapping": 23
  },
  {
   "endpoint_group": "SITE-B",
   "service": "service-A",
   "flapping": 9
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/services?start_date=2015-04-01&end_date=2015-05-02&granularity=monthly",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "flapping": 25
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "flapping": 16
    },
    {
     "endpoint_group": "SITE-XB",
     "service": "service-XA",
     "flapping": 3
    }
   ]
  },
  {
   "date": "2015-05",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "flapping": 98
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "flapping": 23
    },
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "flapping": 9
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/services?start_date=2015-04-01&end_date=2015-05-02&top=2&granularity=monthly",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "flapping": 25
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "flapping": 16
    }
   ]
  },
  {
   "date": "2015-05",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "flapping": 98
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "flapping": 23
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/services?start_date=2015-05-01&end_date=2015-05-02&top=2",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "endpoint_group": "SITE-A",
   "service": "service-A",
   "flapping": 98
  },
  {
   "endpoint_group": "SITE-A",
   "service": "service-B",
   "flapping": 23
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/groups?start_date=2015-05-01",
			code:   400,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "Bad Request",
   "code": "400",
   "details": "please use either a date url parameter or a combination of start_date and end_date parameters to declare range"
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/groups?end_date=2015-05-01",
			code:   400,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "Bad Request",
   "code": "400",
   "details": "please use either a date url parameter or a combination of start_date and end_date parameters to declare range"
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/groups?start_date=2015-05-01&end_date=2015-05-02",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "endpoint_group": "SITE-A",
   "flapping": 66
  },
  {
   "endpoint_group": "SITE-B",
   "flapping": 9
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/groups?start_date=2015-04-01&end_date=2015-05-02&granularity=monthly",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "flapping": 35
    },
    {
     "endpoint_group": "SITE-XB",
     "flapping": 3
    }
   ]
  },
  {
   "date": "2015-05",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "flapping": 66
    },
    {
     "endpoint_group": "SITE-B",
     "flapping": 9
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/groups?start_date=2015-04-01&end_date=2015-05-02&top=1&granularity=monthly",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "flapping": 35
    }
   ]
  },
  {
   "date": "2015-05",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "flapping": 66
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/flapping/groups?start_date=2015-05-01&end_date=2015-05-02&top=1",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "endpoint_group": "SITE-A",
   "flapping": 66
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/metrics?start_date=2015-05-01&end_date=2015-05-02",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "status": "CRITICAL",
     "events": 40
    },
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "metric": "web-check",
     "status": "CRITICAL",
     "events": 7
    }
   ]
  },
  {
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-1",
     "status": "UNKNOWN",
     "events": 45
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "status": "UNKNOWN",
     "events": 32
    },
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "metric": "web-check",
     "status": "UNKNOWN",
     "events": 5
    }
   ]
  },
  {
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-1",
     "status": "WARNING",
     "events": 55
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "metric": "web-check",
     "status": "WARNING",
     "events": 20
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/metrics?start_date=2015-05-01&end_date=2015-05-02&top=1",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "status": "CRITICAL",
     "events": 40
    }
   ]
  },
  {
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-1",
     "status": "UNKNOWN",
     "events": 45
    }
   ]
  },
  {
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-1",
     "status": "WARNING",
     "events": 55
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/metrics?start_date=2015-04-01&end_date=2015-05-02&granularity=monthly",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-1",
     "status": "CRITICAL",
     "events": 55
    },
    {
     "endpoint_group": "SITE-XB",
     "service": "service-XA",
     "endpoint": "hosta.examplex2.foo",
     "metric": "web-check",
     "status": "CRITICAL",
     "events": 25
    }
   ]
  },
  {
   "date": "2015-04",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "metric": "web-check",
     "status": "UNKNOWN",
     "events": 12
    }
   ]
  },
  {
   "date": "2015-04",
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "status": "WARNING",
     "events": 40
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "status": "CRITICAL",
     "events": 40
    },
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "metric": "web-check",
     "status": "CRITICAL",
     "events": 7
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-1",
     "status": "UNKNOWN",
     "events": 45
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "status": "UNKNOWN",
     "events": 32
    },
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "metric": "web-check",
     "status": "UNKNOWN",
     "events": 5
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-1",
     "status": "WARNING",
     "events": 55
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "metric": "web-check",
     "status": "WARNING",
     "events": 20
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/metrics/tags?start_date=2015-04-01&end_date=2015-05-02&granularity=monthly&top=5",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-05",
   "status": "UNKNOWN",
   "tag": "STORAGE",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "status": "UNKNOWN",
     "events": 32
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "WARNING",
   "tag": "HTTP",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "metric": "web-check",
     "status": "WARNING",
     "events": 8
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "WARNING",
   "tag": "NETWORK",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "metric": "web-check",
     "status": "WARNING",
     "events": 8
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/metrics?start_date=2015-04-01&end_date=2015-05-02&granularity=monthly&top=1",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-1",
     "status": "CRITICAL",
     "events": 55
    }
   ]
  },
  {
   "date": "2015-04",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "metric": "web-check",
     "status": "UNKNOWN",
     "events": 12
    }
   ]
  },
  {
   "date": "2015-04",
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "status": "WARNING",
     "events": 40
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "status": "CRITICAL",
     "events": 40
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-1",
     "status": "UNKNOWN",
     "events": 45
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-1",
     "status": "WARNING",
     "events": 55
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/metrics?date=2015-05-01",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-2",
     "status": "CRITICAL",
     "events": 40
    }
   ]
  },
  {
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "metric": "web-check",
     "status": "UNKNOWN",
     "events": 5
    }
   ]
  },
  {
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "metric": "check-1",
     "status": "WARNING",
     "events": 55
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "metric": "web-check",
     "status": "WARNING",
     "events": 12
    }
   ]
  }
 ]
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
		suite.Equal(expReq.result, response.Body.String(), "Response body mismatch")

	}

}

func (suite *TrendsTestSuite) TestStatusEndpointTrends() {

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
			url:    "/api/v2/trends/Report_A/status/endpoints?start_date=2015-05-01&end_date=2015-05-02",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "status": "CRITICAL",
     "duration_in_minutes": 7
    }
   ]
  },
  {
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "status": "UNKNOWN",
     "duration_in_minutes": 45
    },
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "status": "UNKNOWN",
     "duration_in_minutes": 5
    }
   ]
  },
  {
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "status": "WARNING",
     "duration_in_minutes": 55
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "status": "WARNING",
     "duration_in_minutes": 20
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/endpoints?start_date=2015-05-01&end_date=2015-05-02&top=1",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "status": "CRITICAL",
     "duration_in_minutes": 7
    }
   ]
  },
  {
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "status": "UNKNOWN",
     "duration_in_minutes": 45
    }
   ]
  },
  {
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "status": "WARNING",
     "duration_in_minutes": 55
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/endpoints?start_date=2015-04-01&end_date=2015-05-02&granularity=monthly",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "status": "CRITICAL",
     "duration_in_minutes": 55
    },
    {
     "endpoint_group": "SITE-XB",
     "service": "service-XA",
     "endpoint": "hosta.examplex2.foo",
     "status": "CRITICAL",
     "duration_in_minutes": 25
    }
   ]
  },
  {
   "date": "2015-04",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "status": "UNKNOWN",
     "duration_in_minutes": 12
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "status": "CRITICAL",
     "duration_in_minutes": 7
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "status": "UNKNOWN",
     "duration_in_minutes": 45
    },
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "status": "UNKNOWN",
     "duration_in_minutes": 5
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "status": "WARNING",
     "duration_in_minutes": 55
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "status": "WARNING",
     "duration_in_minutes": 20
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/endpoints?start_date=2015-04-01&end_date=2015-05-02&granularity=monthly&top=1",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "status": "CRITICAL",
     "duration_in_minutes": 55
    }
   ]
  },
  {
   "date": "2015-04",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "status": "UNKNOWN",
     "duration_in_minutes": 12
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "status": "CRITICAL",
     "duration_in_minutes": 7
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "status": "UNKNOWN",
     "duration_in_minutes": 45
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "status": "WARNING",
     "duration_in_minutes": 55
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/endpoints?date=2015-05-01",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "endpoint": "hosta.example2.foo",
     "status": "UNKNOWN",
     "duration_in_minutes": 5
    }
   ]
  },
  {
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "endpoint": "hosta.example.foo",
     "status": "WARNING",
     "duration_in_minutes": 55
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "endpoint": "hostb.example.foo",
     "status": "WARNING",
     "duration_in_minutes": 12
    }
   ]
  }
 ]
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
		// Compare the expected and actual json response
		suite.Equal(expReq.result, response.Body.String(), "Response body mismatch")

	}

}

func (suite *TrendsTestSuite) TestStatusServiceTrends() {

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
			url:    "/api/v2/trends/Report_A/status/services?start_date=2015-05-01&end_date=2015-05-02",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "status": "CRITICAL",
     "duration_in_minutes": 7
    }
   ]
  },
  {
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "status": "UNKNOWN",
     "duration_in_minutes": 45
    },
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "status": "UNKNOWN",
     "duration_in_minutes": 5
    }
   ]
  },
  {
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "status": "WARNING",
     "duration_in_minutes": 55
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "status": "WARNING",
     "duration_in_minutes": 20
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/services?start_date=2015-05-01&end_date=2015-05-02&top=1",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "status": "CRITICAL",
     "duration_in_minutes": 7
    }
   ]
  },
  {
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "status": "UNKNOWN",
     "duration_in_minutes": 45
    }
   ]
  },
  {
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "status": "WARNING",
     "duration_in_minutes": 55
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/services?start_date=2015-04-01&end_date=2015-05-02&granularity=monthly",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "status": "CRITICAL",
     "duration_in_minutes": 55
    },
    {
     "endpoint_group": "SITE-XB",
     "service": "service-XA",
     "status": "CRITICAL",
     "duration_in_minutes": 25
    }
   ]
  },
  {
   "date": "2015-04",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "status": "UNKNOWN",
     "duration_in_minutes": 12
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "status": "CRITICAL",
     "duration_in_minutes": 7
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "status": "UNKNOWN",
     "duration_in_minutes": 45
    },
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "status": "UNKNOWN",
     "duration_in_minutes": 5
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "status": "WARNING",
     "duration_in_minutes": 55
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "status": "WARNING",
     "duration_in_minutes": 20
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/services?start_date=2015-04-01&end_date=2015-05-02&granularity=monthly&top=1",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "status": "CRITICAL",
     "duration_in_minutes": 55
    }
   ]
  },
  {
   "date": "2015-04",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "status": "UNKNOWN",
     "duration_in_minutes": 12
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "status": "CRITICAL",
     "duration_in_minutes": 7
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "status": "UNKNOWN",
     "duration_in_minutes": 45
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "status": "WARNING",
     "duration_in_minutes": 55
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/services?date=2015-05-01",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "service": "service-A",
     "status": "UNKNOWN",
     "duration_in_minutes": 5
    }
   ]
  },
  {
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "status": "WARNING",
     "duration_in_minutes": 55
    },
    {
     "endpoint_group": "SITE-A",
     "service": "service-B",
     "status": "WARNING",
     "duration_in_minutes": 12
    }
   ]
  }
 ]
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
		suite.Equal(expReq.result, response.Body.String(), "Response body mismatch")

	}

}

func (suite *TrendsTestSuite) TestStatusEgroupTrends() {

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
			url:    "/api/v2/trends/Report_A/status/groups?start_date=2015-05-01&end_date=2015-05-02",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "status": "CRITICAL",
     "duration_in_minutes": 7
    }
   ]
  },
  {
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "status": "UNKNOWN",
     "duration_in_minutes": 45
    },
    {
     "endpoint_group": "SITE-B",
     "status": "UNKNOWN",
     "duration_in_minutes": 5
    }
   ]
  },
  {
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "status": "WARNING",
     "duration_in_minutes": 55
    },
    {
     "endpoint_group": "SITE-C",
     "status": "WARNING",
     "duration_in_minutes": 20
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/groups?start_date=2015-05-01&end_date=2015-05-02&top=1",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "status": "CRITICAL",
     "duration_in_minutes": 7
    }
   ]
  },
  {
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "status": "UNKNOWN",
     "duration_in_minutes": 45
    }
   ]
  },
  {
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "status": "WARNING",
     "duration_in_minutes": 55
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/groups?start_date=2015-04-01&end_date=2015-05-02&granularity=monthly",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "status": "CRITICAL",
     "duration_in_minutes": 55
    },
    {
     "endpoint_group": "SITE-XB",
     "status": "CRITICAL",
     "duration_in_minutes": 25
    }
   ]
  },
  {
   "date": "2015-04",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "status": "UNKNOWN",
     "duration_in_minutes": 12
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "status": "CRITICAL",
     "duration_in_minutes": 7
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "status": "UNKNOWN",
     "duration_in_minutes": 45
    },
    {
     "endpoint_group": "SITE-B",
     "status": "UNKNOWN",
     "duration_in_minutes": 5
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "status": "WARNING",
     "duration_in_minutes": 55
    },
    {
     "endpoint_group": "SITE-C",
     "status": "WARNING",
     "duration_in_minutes": 20
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/groups?start_date=2015-04-01&end_date=2015-05-02&granularity=monthly&top=1",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "status": "CRITICAL",
     "duration_in_minutes": 55
    }
   ]
  },
  {
   "date": "2015-04",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "status": "UNKNOWN",
     "duration_in_minutes": 12
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "CRITICAL",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "status": "CRITICAL",
     "duration_in_minutes": 7
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "status": "UNKNOWN",
     "duration_in_minutes": 45
    }
   ]
  },
  {
   "date": "2015-05",
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "status": "WARNING",
     "duration_in_minutes": 55
    }
   ]
  }
 ]
}`,
		},

		{
			method: "GET",
			url:    "/api/v2/trends/Report_A/status/groups?date=2015-05-01",
			code:   200,
			key:    "KEY1",
			result: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "status": "UNKNOWN",
   "top": [
    {
     "endpoint_group": "SITE-B",
     "status": "UNKNOWN",
     "duration_in_minutes": 5
    }
   ]
  },
  {
   "status": "WARNING",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "status": "WARNING",
     "duration_in_minutes": 55
    },
    {
     "endpoint_group": "SITE-C",
     "status": "WARNING",
     "duration_in_minutes": 12
    }
   ]
  }
 ]
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
		suite.Equal(expReq.result, response.Body.String(), "Response body mismatch")

	}

}

func (suite *TrendsTestSuite) TestOptionsStatusTrendsEgroups() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/trends/Report_A/status/groups", strings.NewReader(""))

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

func (suite *TrendsTestSuite) TestOptionsStatusTrendsServices() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/trends/Report_A/status/services", strings.NewReader(""))

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

func (suite *TrendsTestSuite) TestOptionsStatusTrendsEndpoints() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/trends/Report_A/status/endpoints", strings.NewReader(""))

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

func (suite *TrendsTestSuite) TestOptionsStatusTrendsMetrics() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/trends/Report_A/status/metrics", strings.NewReader(""))

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

func (suite *TrendsTestSuite) TestOptionsTrendsMetrics() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/trends/Report_A/flapping/metrics", strings.NewReader(""))

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

func (suite *TrendsTestSuite) TestOptionsTrendsEndpoints() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/trends/Report_A/flapping/endpoints", strings.NewReader(""))

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

func (suite *TrendsTestSuite) TestOptionsTrendsServices() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/trends/Report_A/flapping/services", strings.NewReader(""))

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

func (suite *TrendsTestSuite) TestOptionsTrendsEndpointGroups() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/trends/Report_A/flapping/groups", strings.NewReader(""))

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

// This function is actually called in the end of all tests
// and clears the test environment.
// Mainly it's purpose is to drop the testdb
func (suite *TrendsTestSuite) TearDownTest() {

	suite.cfg.MongoClient.Database("argotest_trends").Drop(context.TODO())
	suite.cfg.MongoClient.Database("argotest_trends_tenant").Drop(context.TODO())
}

// This is the first function called when go test is issued
func TestSuiteTrends(t *testing.T) {
	suite.Run(t, new(TrendsTestSuite))
}
