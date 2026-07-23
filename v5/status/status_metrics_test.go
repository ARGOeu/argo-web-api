package statusV5

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
type StatusMetricsTestSuite struct {
	suite.Suite
	cfg         config.Config
	router      *mux.Router
	confHandler respond.ConfHandler
}

func (suite *StatusMetricsTestSuite) SetupTest() {

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
    db = "argotest_status_metrics"
`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(true).PathPrefix("/api/v5/status").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)

	// Add authentication token to mongo testdb
	authCol := suite.cfg.MongoClient.Database((suite.cfg.MongoDB.Db)).Collection("authentication")
	// Add authentication token to mongo testdb
	seedAuth := bson.M{"api_key": "S3CR3T"}
	authCol.InsertOne(context.TODO(), seedAuth)

	// seed a tenant to use
	c := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("tenants")
	c.InsertOne(context.TODO(), bson.M{
		"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
		"info": bson.M{
			"name":    "TENANT-AA",
			"email":   "example@tenant-a.foo",
			"website": "www.tenant-aa.com",
			"created": "2015-10-20 02:08:04",
			"updated": "2015-10-20 02:08:04"},
		"db_conf": []bson.M{
			{
				"store":    "main",
				"server":   "localhost",
				"port":     27017,
				"database": "argotest_status_metrics_ta",
				"username": "",
				"password": ""},
		},
		"users": []bson.M{
			{
				"name":    "bob",
				"email":   "bob@example.foo",
				"roles":   []string{"viewer"},
				"api_key": "KEY1"},
		}})

	c.InsertOne(context.TODO(), bson.M{
		"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50d",
		"info": bson.M{
			"name":    "TENANT-BB",
			"email":   "email@tenant-bb.foo",
			"website": "www.tenant-bb.foo",
			"created": "2015-10-20 02:08:04",
			"updated": "2015-10-20 02:08:04"},
		"db_conf": []bson.M{
			bson.M{
				"store":    "main",
				"server":   "localhost",
				"port":     27017,
				"database": "argotest_status_metrics_tb",
				"username": "",
				"password": ""},
		},
		"users": []bson.M{
			bson.M{
				"name":    "bob",
				"email":   "bob@email.foo",
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
	c = suite.cfg.MongoClient.Database(t1conf.Db).Collection("status_metrics")
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"monitoring_box":     "sensu",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T00:00:00Z",
		"service":            "CREAM-CE",
		"host":               "cream01.example.foo",
		"endpoint_group":     "GROUP-03",
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
		"monitoring_box":     "sensu",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T01:00:00Z",
		"service":            "CREAM-CE",
		"host":               "cream01.example.foo",
		"endpoint_group":     "GROUP-03",
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
		"monitoring_box":     "sensu",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T05:00:00Z",
		"service":            "CREAM-CE",
		"host":               "cream01.example.foo",
		"endpoint_group":     "GROUP-03",
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
		"monitoring_box":     "sensu",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T00:00:00Z",
		"service":            "CREAM-CE",
		"host":               "cream01.example.foo",
		"endpoint_group":     "GROUP-03",
		"metric":             "emi.cream.CREAMCE-JobCancel",
		"status":             "OK",
		"time_integer":       0,
		"previous_state":     "OK",
		"previous_timestamp": "2015-04-30T23:59:00Z",
		"summary":            "Cream status is ok",
		"message":            "Cream job submission test return value of ok",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":                 "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"monitoring_box":         "sensu",
		"date_integer":           20150501,
		"timestamp":              "2015-05-01T01:00:00Z",
		"service":                "CREAM-CE",
		"host":                   "cream01.example.foo",
		"endpoint_group":         "GROUP-03",
		"metric":                 "emi.cream.CREAMCE-JobCancel",
		"status":                 "CRITICAL",
		"time_integer":           10000,
		"previous_state":         "OK",
		"previous_timestamp":     "2015-05-01T00:00:00Z",
		"summary":                "Cream status is CRITICAL",
		"message":                "Cream job submission test failed",
		"actual_data":            "latency=15s",
		"threshold_rule_applied": "latency=1s;0:5;10:60",
		"original_status":        "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"monitoring_box":     "sensu",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T05:00:00Z",
		"service":            "CREAM-CE",
		"host":               "cream01.example.foo",
		"endpoint_group":     "GROUP-03",
		"metric":             "emi.cream.CREAMCE-JobCancel",
		"status":             "OK",
		"time_integer":       50000,
		"previous_state":     "CRITICAL",
		"previous_timestamp": "2015-05-01T01:00:00Z",
		"summary":            "Cream status is ok",
		"message":            "Cream job submission test return value of ok",
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
			"description": "report bbb",
			"created":     "2015-9-10 13:43:00",
			"updated":     "2015-10-11 13:43:00",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "GROUPS",
				"group": bson.M{
					"type": "SITES",
				},
			},
		},
		"profiles": []bson.M{
			{
				"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
				"type": "metric",
				"name": "metric.CRITICAL"},
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
	c = suite.cfg.MongoClient.Database(t2conf.Db).Collection("status_metrics")
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"monitoring_box":     "sensu",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T00:00:00Z",
		"service":            "someService",
		"host":               "someservice.example.foo",
		"endpoint_group":     "GROUP-01",
		"metric":             "someService-FileTransfer",
		"info":               bson.M{"URL": "http://host.example.foo"},
		"status":             "OK",
		"time_integer":       0,
		"previous_state":     "OK",
		"previous_timestamp": "2015-04-30T23:59:00Z",
		"summary":            "someService status is ok",
		"message":            "someService data upload test return value of ok",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":                 "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"monitoring_box":         "sensu",
		"date_integer":           20150501,
		"timestamp":              "2015-05-01T01:00:00Z",
		"service":                "someService",
		"host":                   "someservice.example.foo",
		"endpoint_group":         "GROUP-01",
		"metric":                 "someService-FileTransfer",
		"info":                   bson.M{"URL": "http://host.example.foo"},
		"status":                 "CRITICAL",
		"time_integer":           10000,
		"previous_state":         "OK",
		"previous_timestamp":     "2015-05-01T00:00:00Z",
		"summary":                "someService status is CRITICAL",
		"message":                "someService data upload test failed",
		"actual_data":            "latency=15s",
		"threshold_rule_applied": "latency=1s;0:5;10:60",
		"original_status":        "OK",
	})
	c.InsertOne(context.TODO(), bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"monitoring_box":     "sensu",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T05:00:00Z",
		"service":            "someService",
		"host":               "someservice.example.foo",
		"endpoint_group":     "GROUP-01",
		"metric":             "someService-FileTransfer",
		"info":               bson.M{"URL": "http://host.example.foo"},
		"status":             "OK",
		"time_integer":       50000,
		"previous_state":     "CRITICAL",
		"previous_timestamp": "2015-05-01T01:00:00Z",
		"summary":            "someService status is ok",
		"message":            "someService data upload test return value of ok",
	})

}

func (suite *StatusMetricsTestSuite) TestListStatusMetrics() {

	respJSON1 := `{
 "groups": [
  {
   "name": "GROUP-03",
   "type": "SITES",
   "service-types": [
    {
     "name": "CREAM-CE",
     "type": "service",
     "endpoints": [
      {
       "name": "cream01.example.foo",
       "metrics": [
        {
         "name": "emi.cream.CREAMCE-JobSubmit",
         "statuses": [
          {
           "timestamp": "2015-04-30T23:59:00Z",
           "value": "OK"
          },
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
          }
         ]
        }
       ]
      }
     ]
    }
   ]
  }
 ]
}`

	respJSON2 := `{
 "groups": [
  {
   "name": "GROUP-01",
   "type": "SITES",
   "service-types": [
    {
     "name": "someService",
     "type": "service",
     "endpoints": [
      {
       "name": "someservice.example.foo",
       "info": {
        "URL": "http://host.example.foo"
       },
       "metrics": [
        {
         "name": "someService-FileTransfer",
         "statuses": [
          {
           "timestamp": "2015-04-30T23:59:00Z",
           "value": "OK"
          },
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
          }
         ]
        }
       ]
      }
     ]
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

	respJSON3 := `{
 "groups": [
  {
   "name": "GROUP-01",
   "type": "SITES",
   "service-types": [
    {
     "name": "someService",
     "type": "service",
     "endpoints": [
      {
       "name": "someservice.example.foo",
       "info": {
        "URL": "http://host.example.foo"
       },
       "metrics": [
        {
         "name": "someService-FileTransfer",
         "statuses": [
          {
           "timestamp": "2015-04-30T23:59:00Z",
           "value": "OK"
          },
          {
           "timestamp": "2015-05-01T00:00:00Z",
           "value": "OK"
          },
          {
           "timestamp": "2015-05-01T01:00:00Z",
           "value": "CRITICAL",
           "actual_data": "latency=15s",
           "threshold_rule_applied": "latency=1s;0:5;10:60",
           "original_status": "OK"
          },
          {
           "timestamp": "2015-05-01T05:00:00Z",
           "value": "OK"
          }
         ]
        }
       ]
      }
     ]
    }
   ]
  }
 ]
}`

	fullurl1 := "/api/v5/status/Report_A/groups/GROUP-03" +
		"/service-types/CREAM-CE/endpoints/cream01.example.foo/metrics/emi.cream.CREAMCE-JobSubmit" +
		"?start-time=2015-05-01T00:00:00Z&end-time=2015-05-01T23:00:00Z"

	fullurl2 := "/api/v5/status/Report_B/groups/GROUP-01" +
		"/service-types/someService/endpoints/someservice.example.foo/metrics/someService-FileTransfer" +
		"?start-time=2015-05-01T00:00:00Z&end-time=2015-05-01T23:00:00Z"

	fullurl3 := "/api/v5/status/Report_B/groups/GROUP-01" +
		"/service-types/someService/endpoints/someservice.example.foo/metrics" +
		"?start-time=2015-05-01T00:00:00Z&end-time=2015-05-01T23:00:00Z"

	fullurl4 := "/api/v5/status/Report_B/groups/GROUP-01" +
		"/service-types/someService/endpoints/someservice.example.foo/metrics/someService-FileTransfer" +
		"?start-time=2015-05-01T00:00:00Z&end-time=2015-05-01T23:00:00Z&view=details"

	// 3. JSON REQUEST
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

	// 4. JSON REQUEST
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

	// 5. TENANT2 ALL JSON REQUEST
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
	request, _ = http.NewRequest("GET", fullurl1, strings.NewReader(""))
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

	// 7. TENANT2 JSON REQUEST DETAILS
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
	suite.Equal(respJSON3, response.Body.String(), "Response body mismatch")

}

func (suite *StatusMetricsTestSuite) TestLatestResults() {

	fullurl1 := "/api/v5/status/Report_A/groups/GROUP-039" +
		"/service-types/CREAM-CE/endpoints/cream01.example.foo/metrics/emi.cream.CREAMCE-JobSubmit" +
		"?start-time=2015-05-03T00:00:00Z&end-time=2015-05-03T23:00:00Z"

	fullurl2 := "/api/v5/status/Report_A/groups/GROUP-03" +
		"/service-types/CREAM-CE/endpoints/cream01.example.foo/metrics/emi.cream.CREAMCE-JobSubmit" +
		"?start-time=2015-05-02T00:00:00Z&end-time=2015-05-02T23:00:00Z"

	respJSON1 := `{
 "groups": [
  {
   "name": "GROUP-03",
   "type": "SITES",
   "service-types": [
    {
     "name": "CREAM-CE",
     "type": "service",
     "endpoints": [
      {
       "name": "cream01.example.foo",
       "metrics": [
        {
         "name": "emi.cream.CREAMCE-JobSubmit",
         "statuses": [
          {
           "timestamp": "2015-04-30T23:59:00Z",
           "value": "OK"
          },
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
          }
         ]
        }
       ]
      }
     ]
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

func (suite *StatusMetricsTestSuite) TestMultipleItems() {

	fullurl1 := "/api/v5/status/Report_A/groups/GROUP-03" +
		"/service-types/CREAM-CE/endpoints/cream01.example.foo/metrics" +
		"?start-time=2015-05-01T00:00:00Z&end-time=2015-05-01T23:00:00Z"

	respJSON1 := `{
 "groups": [
  {
   "name": "GROUP-03",
   "type": "SITES",
   "service-types": [
    {
     "name": "CREAM-CE",
     "type": "service",
     "endpoints": [
      {
       "name": "cream01.example.foo",
       "metrics": [
        {
         "name": "emi.cream.CREAMCE-JobCancel",
         "statuses": [
          {
           "timestamp": "2015-04-30T23:59:00Z",
           "value": "OK"
          },
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
          }
         ]
        },
        {
         "name": "emi.cream.CREAMCE-JobSubmit",
         "statuses": [
          {
           "timestamp": "2015-04-30T23:59:00Z",
           "value": "OK"
          },
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
          }
         ]
        }
       ]
      }
     ]
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

func (suite *StatusMetricsTestSuite) TestMultipleItemsDetails() {

	fullurl1 := "/api/v5/status/Report_A/groups/GROUP-03" +
		"/service-types/CREAM-CE/endpoints/cream01.example.foo/metrics" +
		"?start-time=2015-05-01T00:00:00Z&end-time=2015-05-01T23:00:00Z&view=details"

	respJSON1 := `{
 "groups": [
  {
   "name": "GROUP-03",
   "type": "SITES",
   "service-types": [
    {
     "name": "CREAM-CE",
     "type": "service",
     "endpoints": [
      {
       "name": "cream01.example.foo",
       "metrics": [
        {
         "name": "emi.cream.CREAMCE-JobCancel",
         "statuses": [
          {
           "timestamp": "2015-04-30T23:59:00Z",
           "value": "OK"
          },
          {
           "timestamp": "2015-05-01T00:00:00Z",
           "value": "OK"
          },
          {
           "timestamp": "2015-05-01T01:00:00Z",
           "value": "CRITICAL",
           "actual_data": "latency=15s",
           "threshold_rule_applied": "latency=1s;0:5;10:60",
           "original_status": "OK"
          },
          {
           "timestamp": "2015-05-01T05:00:00Z",
           "value": "OK"
          }
         ]
        },
        {
         "name": "emi.cream.CREAMCE-JobSubmit",
         "statuses": [
          {
           "timestamp": "2015-04-30T23:59:00Z",
           "value": "OK"
          },
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
          }
         ]
        }
       ]
      }
     ]
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

func (suite *StatusMetricsTestSuite) TestOptionsStatusMetrics() {
	request, _ := http.NewRequest("OPTIONS", "/api/v5/status/Report_A/groups/GROUP_A/service-types/service_a/endpoints/endpoint_a/metrics", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v5/status/Report_A/groups/GROUP_A/service-types/service_a/endpoints/endpoint_a/metrics/metric_a", strings.NewReader(""))

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
func (suite *StatusMetricsTestSuite) TearDownTest() {

	suite.cfg.MongoClient.Database("argotest_status_metrics").Drop(context.TODO())
	suite.cfg.MongoClient.Database("argotest_status_metrics_ta").Drop(context.TODO())
	suite.cfg.MongoClient.Database("argotest_status_metrics_tb").Drop(context.TODO())

}

// This is the first function called when go test is issued
func TestSuiteStatusMetrics(t *testing.T) {
	suite.Run(t, new(StatusMetricsTestSuite))
}
