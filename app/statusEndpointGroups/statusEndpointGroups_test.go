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

package statusEndpointGroups

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
type StatusEndpointGroupsTestSuite struct {
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
func (suite *StatusEndpointGroupsTestSuite) SetupTest() {

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
    db = "argotest_egroups"
`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(true).PathPrefix("/api/v2/status").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)

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
				"database": "argotest_egroups_egi",
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
				"database": "argotest_egroups_eudat",
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
			"resource": "status.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "status.get",
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

	// seed the status detailed metric data
	c = suite.cfg.MongoClient.Database(t1conf.Db).Collection("status_endpoint_groups")
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"status":         "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T01:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"status":         "CRITICAL",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T05:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"status":         "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "HG-02-AUTH",
		"status":         "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T03:00:00Z",
		"endpoint_group": "HG-02-AUTH",
		"status":         "WARNING",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T17:53:00Z",
		"endpoint_group":     "HG-02-AUTH",
		"status":             "CRITICAL",
		"has_threshold_rule": true,
	})

	// get dbconfiguration based on the tenant
	// Prepare the request object
	request, _ = http.NewRequest("GET", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY2")
	// authenticate user's api key and find corresponding tenant
	t2conf, _, _ := authentication.AuthenticateTenant(request.Header, suite.cfg)

	// Now seed the reports DEFINITIONS
	c = suite.cfg.MongoClient.Database(t2conf.Db).Collection("reports")
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
	c = suite.cfg.MongoClient.Database(t2conf.Db).Collection("status_endpoint_groups")
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "EL-01-AUTH",
		"status":         "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T01:00:00Z",
		"endpoint_group": "EL-01-AUTH",
		"status":         "CRITICAL",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T05:00:00Z",
		"endpoint_group": "EL-01-AUTH",
		"status":         "OK",
	})

}

func (suite *StatusEndpointGroupsTestSuite) TestListStatusEndpointGroups() {
	respXML1 := `<root>
 <group name="HG-03-AUTH" type="SITES">
  <status timestamp="2015-05-01T00:00:00Z" value="OK"></status>
  <status timestamp="2015-05-01T01:00:00Z" value="CRITICAL"></status>
  <status timestamp="2015-05-01T05:00:00Z" value="OK"></status>
  <status timestamp="2015-05-01T23:59:59Z" value="OK"></status>
 </group>
</root>`

	respXML2 := `<root>
 <group name="EL-01-AUTH" type="EUDAT_SITES">
  <status timestamp="2015-05-01T00:00:00Z" value="OK"></status>
  <status timestamp="2015-05-01T01:00:00Z" value="CRITICAL"></status>
  <status timestamp="2015-05-01T05:00:00Z" value="OK"></status>
  <status timestamp="2015-05-01T23:59:59Z" value="OK"></status>
 </group>
</root>`

	respJSON1 := `{
 "groups": [
  {
   "name": "HG-03-AUTH",
   "type": "SITES",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T01:00:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T05:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  }
 ]
}`
	respJSON2 := `{
 "groups": [
  {
   "name": "EL-01-AUTH",
   "type": "EUDAT_SITES",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T01:00:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T05:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  }
 ]
}`

	respUnauthorized := `{
 "status": {
  "message": "Unauthorized",
  "code": "401",
  "details": "You need to provide a correct authentication token using the header 'x-api-key'"
 }
}`

	fullurl1 := "/api/v2/status/Report_A/SITES/HG-03-AUTH" +
		"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z"

	fullurl2 := "/api/v2/status/Report_B/EUDAT_SITES/EL-01-AUTH" +
		"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z"

	fullurl3 := "/api/v2/status/Report_B/EUDAT_SITES" +
		"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z"

	// 1. EGI XML REQUEST
	// init the response placeholder
	response := httptest.NewRecorder()
	// Prepare the request object for fist tenant
	request, _ := http.NewRequest("GET", fullurl1, strings.NewReader(""))
	// add accept xml header
	request.Header.Set("Accept", "application/xml")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY1")
	// Serve the http request
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respXML1, response.Body.String(), "Response body mismatch")

	// 2. EUDAT XML REQUEST
	// init the response placeholder
	response = httptest.NewRecorder()
	// Prepare the request object for second tenant
	request, _ = http.NewRequest("GET", fullurl2, strings.NewReader(""))
	// add accept xml header
	request.Header.Set("Accept", "application/xml")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY2")
	// Serve the http request
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respXML2, response.Body.String(), "Response body mismatch")

	// 3. EGI JSON REQUEST
	// init the response placeholder
	response = httptest.NewRecorder()
	// Prepare the request object for second tenant
	request, _ = http.NewRequest("GET", fullurl1, strings.NewReader(""))
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

	// 4. EUDAT JSON REQUEST
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

	// 5. EUDAT ALL JSON REQUEST
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
	suite.Equal(respJSON2, response.Body.String(), "Response body mismatch")

	// 6. WRONG KEY REQUEST
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

}

func (suite *StatusEndpointGroupsTestSuite) TestLatestResults() {

	fullurl1 := "/api/v2/status/Report_A/SITES/HG-03-AUTH" +
		"?start_time=2015-05-03T00:00:00Z&end_time=2015-05-03T23:00:00Z"

	fullurl2 := "/api/v2/status/Report_A/SITES/HG-03-AUTH" +
		"?start_time=2015-05-02T00:00:00Z&end_time=2015-05-02T23:00:00Z"

	respJSON1 := `{
 "groups": [
  {
   "name": "HG-03-AUTH",
   "type": "SITES",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T01:00:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T05:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  }
 ]
}`

	respEmptyJSON := `{
   "groups": null
 }`

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
	suite.Equal(respEmptyJSON, response.Body.String(), "Response body mismatch")

	// init the response placeholder
	response2 := httptest.NewRecorder()
	// Prepare the request object for second tenant
	request2, _ := http.NewRequest("GET", fullurl2, strings.NewReader(""))
	// add json accept header
	request2.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request2.Header.Set("x-api-key", "KEY1")
	// Serve the http request
	suite.router.ServeHTTP(response2, request2)
	// Check that we must have a 200 ok code
	suite.Equal(200, response2.Code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respJSON1, response2.Body.String(), "Response body mismatch")

}

func (suite *StatusEndpointGroupsTestSuite) TestMultipleItems() {

	fullurl1 := "/api/v2/status/Report_A/SITES" +
		"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z"

	respJSON1 := `{
 "groups": [
  {
   "name": "HG-03-AUTH",
   "type": "SITES",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T01:00:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T05:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  },
  {
   "name": "HG-02-AUTH",
   "type": "SITES",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T03:00:00Z",
     "value": "WARNING"
    },
    {
     "timestamp": "2015-05-01T17:53:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "CRITICAL"
    }
   ]
  }
 ]
}`

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

}

func (suite *StatusEndpointGroupsTestSuite) TestThresholds() {

	fullurl1 := "/api/v2/status/Report_A/SITES" +
		"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z&view=details"

	respJSON1 := `{
 "groups": [
  {
   "name": "HG-03-AUTH",
   "type": "SITES",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T01:00:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T05:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  },
  {
   "name": "HG-02-AUTH",
   "type": "SITES",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T03:00:00Z",
     "value": "WARNING"
    },
    {
     "timestamp": "2015-05-01T17:53:00Z",
     "value": "CRITICAL",
     "affected_by_threshold_rule": true
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "CRITICAL"
    }
   ]
  }
 ]
}`

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

}

func (suite *StatusEndpointGroupsTestSuite) TestOptionsStatusEndpointGroups() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/status/Report_A/GROUP", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v2/status/Report_A/GROUP/GROUP_A", strings.NewReader(""))

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	headers = response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

}

// This function is actually called in the end of all tests
// and clears the test environment.
// Mainly it's purpose is to drop the testdb
func (suite *StatusEndpointGroupsTestSuite) TearDownTest() {

	suite.cfg.MongoClient.Database("argotest_egroups").Drop(context.TODO())
	suite.cfg.MongoClient.Database("argotest_egroups_eudat").Drop(context.TODO())
	suite.cfg.MongoClient.Database("argotest_egroups_egi").Drop(context.TODO())
}

// This is the first function called when go test is issued
func TestSuiteStatusEndpointGroups(t *testing.T) {
	suite.Run(t, new(StatusEndpointGroupsTestSuite))
}
