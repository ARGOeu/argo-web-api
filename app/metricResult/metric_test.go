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

package metricResult

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
type metricResultTestSuite struct {
	suite.Suite
	cfg          config.Config
	router       *mux.Router
	confHandler  respond.ConfHandler
	tenantDbConf config.MongoConfig
	clientkey    string
}

// Setup the Test Environment
// This function runs before any test and setups the environment
// A test configuration object is instantiated using a reference
// to testdb: argo_test_details. Also here is are instantiated some expected
// xml response validation messages (authorization,crud responses).
// Also the testdb is seeded with tenants,reports,metric_profiles and status_metrics
func (suite *metricResultTestSuite) SetupTest() {

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
    db = "ARGO_test_metric_result"
`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.tenantDbConf.Db = "ARGO_test_metric_result"
	suite.tenantDbConf.Password = "h4shp4ss"
	suite.tenantDbConf.Username = "dbuser"
	suite.tenantDbConf.Store = "ar"
	suite.clientkey = "KEY1"

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(true).PathPrefix("/api/v2/metric_result").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)

	// Connect to mongo testdb
	coreCol := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("authentication")

	// Add authentication token to mongo testdb
	seedAuth := bson.M{"api_key": "S3CR3T"}
	coreCol.InsertOne(context.TODO(), seedAuth)

	// seed a tenant to use

	c := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("tenants")
	c.InsertOne(context.TODO(), bson.M{
		"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
		"info": bson.M{
			"name":    "EGI",
			"email":   "email@something2",
			"website": "www.gotg.com",
			"created": "2015-10-20 02:08:04",
			"updated": "2015-10-20 02:08:04"},
		"db_conf": []bson.M{
			{
				"store":    "main",
				"server":   "localhost",
				"port":     27017,
				"database": "ARGO_test_metric_result_egi",
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

	c = suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("roles")
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "metricResult.get",
			"roles":    []string{"editor", "viewer"},
		})

	// get dbconfiguration based on the tenant
	// Prepare the request object
	request, _ := http.NewRequest("GET", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY1")
	// authenticate user's api key and find corresponding tenant
	tenantDbConf1, _, _ := authentication.AuthenticateTenant(request.Header, suite.cfg)

	// seed the status detailed metric data
	c = suite.cfg.MongoClient.Database(tenantDbConf1.Db).Collection("status_metrics")
	c.InsertOne(context.TODO(), bson.M{
		"monitoring_box":         "nagios3.hellasgrid.gr",
		"date_integer":           20150501,
		"timestamp":              "2015-05-01T00:00:00Z",
		"service":                "CREAM-CE",
		"host":                   "cream01.afroditi.gr",
		"metric":                 "emi.cream.CREAMCE-JobSubmit",
		"status":                 "OK",
		"time_integer":           0,
		"previous_state":         "OK",
		"previous_timestamp":     "2015-04-30T23:59:00Z",
		"summary":                "Cream status is ok",
		"info":                   bson.M{"URL": "http://creamce.example.com"},
		"message":                "Cream job submission test return value of ok",
		"actual_data":            "latency=15s",
		"threshold_rule_applied": "latency=15s;20:30;30:40",
		"original_status":        "WARNING",
	})
	c.InsertOne(context.TODO(), bson.M{
		"monitoring_box":     "nagios3.hellasgrid.gr",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T01:00:00Z",
		"service":            "CREAM-CE",
		"host":               "cream01.afroditi.gr",
		"metric":             "emi.cream.CREAMCE-JobSubmit",
		"status":             "CRITICAL",
		"time_integer":       10000,
		"previous_state":     "OK",
		"previous_timestamp": "2015-05-01T00:00:00Z",
		"info":               bson.M{"URL": "http://creamce.example.com"},
		"summary":            "Cream status is CRITICAL",
		"message":            "Cream job submission test failed",
	})
	c.InsertOne(context.TODO(), bson.M{
		"monitoring_box":     "nagios3.hellasgrid.gr",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T05:00:00Z",
		"service":            "CREAM-CE",
		"host":               "cream01.afroditi.gr",
		"metric":             "emi.cream.CREAMCE-JobSubmit",
		"info":               bson.M{"URL": "http://creamce.example.com"},
		"status":             "OK",
		"time_integer":       50000,
		"previous_state":     "CRITICAL",
		"previous_timestamp": "2015-05-01T01:00:00Z",
		"summary":            "Cream status is ok",
		"message":            "Cream job submission test return value of ok",
	})

	c.InsertOne(context.TODO(), bson.M{
		"monitoring_box":     "nagios3.hellasgrid.gr",
		"date_integer":       20200501,
		"timestamp":          "2020-05-01T05:00:00Z",
		"service":            "CREAM-CE",
		"host":               "cream01.afroditi.gr",
		"metric":             "emi.cream.CREAMCE-JobSubmit",
		"info":               bson.M{"URL": "http://creamce.example.com"},
		"status":             "OK",
		"time_integer":       50000,
		"previous_state":     "CRITICAL",
		"previous_timestamp": "2020-05-01T01:00:00Z",
		"summary":            "Cream status is ok",
		"message":            "Cream job submission test return value of ok",
	})

	c.InsertOne(context.TODO(), bson.M{
		"monitoring_box":     "nagios3.hellasgrid.gr",
		"date_integer":       20200501,
		"timestamp":          "2020-05-01T05:00:00Z",
		"service":            "CREAM-CE2",
		"host":               "cream01.afroditi.gr",
		"metric":             "emi.cream.CREAMCE-JobSubmit",
		"info":               bson.M{"URL": "http://creamce.example.com"},
		"status":             "OK",
		"time_integer":       50000,
		"previous_state":     "CRITICAL",
		"previous_timestamp": "2020-05-01T01:00:00Z",
		"summary":            "Cream status is ok",
		"message":            "Cream job submission test return value of ok",
	})

}

func (suite *metricResultTestSuite) TestReadStatusDetail() {

	respXML := ` <root>
   <host name="cream01.afroditi.gr">
     <metric name="emi.cream.CREAMCE-JobSubmit" service="CREAM-CE">
       <status timestamp="2015-05-01T01:00:00Z" value="CRITICAL">
         <summary>Cream status is CRITICAL</summary>
         <message>Cream job submission test failed</message>
       </status>
     </metric>
   </host>
 </root>`

	respJSON := `{
   "root": [
     {
       "Name": "cream01.afroditi.gr",
       "info": {
         "URL": "http://creamce.example.com"
       },
       "Metrics": [
         {
           "Name": "emi.cream.CREAMCE-JobSubmit",
           "Service": "CREAM-CE",
           "Details": [
             {
               "Timestamp": "2015-05-01T01:00:00Z",
               "Value": "CRITICAL",
               "Summary": "Cream status is CRITICAL",
               "Message": "Cream job submission test failed"
             }
           ]
         }
       ]
     }
   ]
 }`

	respJSON2 := `{
   "root": [
     {
       "Name": "cream01.afroditi.gr",
       "info": {
         "URL": "http://creamce.example.com"
       },
       "Metrics": [
         {
           "Name": "emi.cream.CREAMCE-JobSubmit",
           "Service": "CREAM-CE",
           "Details": [
             {
               "Timestamp": "2015-05-01T00:00:00Z",
               "Value": "OK",
               "Summary": "Cream status is ok",
               "Message": "Cream job submission test return value of ok",
               "actual_data": "latency=15s",
               "threshold_rule_applied": "latency=15s;20:30;30:40",
               "original_status": "WARNING"
             }
           ]
         }
       ]
     }
   ]
 }`

	respNotFound := `{
   "message": "Metric not found!",
   "code": 404
 }`

	request, _ := http.NewRequest("GET", "/api/v2/metric_result/cream01.afroditi.gr/emi.cream.CREAMCE-JobSubmit?exec_time=2015-05-01T01:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")
	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(respXML, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v2/metric_result/cream01.afroditi.gr/emi.cream.CREAMCE-JobSubmit?exec_time=2015-05-01T01:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(respJSON, response.Body.String(), "Response body mismatch")

	// Check returned xml when no results are available for a given timestamp
	request, _ = http.NewRequest("GET", "/api/v2/metric_result/cream01.afroditi.gr/emi.cream.CREAMCE-JobSubmit?exec_time=2015-05-01T01:01:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(404, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(respNotFound, response.Body.String(), "Response body mismatch")

	// Check returned xml when no results are available for a given timestamp
	request, _ = http.NewRequest("GET", "/api/v2/metric_result/cream01.afroditi.gr/emi.cream.CREAMCE-JobSubmit?exec_time=2015-05-01T00:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(respJSON2, response.Body.String(), "Response body mismatch")

}

func (suite *metricResultTestSuite) TestMultipleMetricResults() {

	respJSON := `{
   "root": [
     {
       "Name": "cream01.afroditi.gr",
       "info": {
         "URL": "http://creamce.example.com"
       },
       "Metrics": [
         {
           "Name": "emi.cream.CREAMCE-JobSubmit",
           "Service": "CREAM-CE",
           "Details": [
             {
               "Timestamp": "2015-05-01T00:00:00Z",
               "Value": "OK",
               "Summary": "Cream status is ok",
               "Message": "Cream job submission test return value of ok",
               "actual_data": "latency=15s",
               "threshold_rule_applied": "latency=15s;20:30;30:40",
               "original_status": "WARNING"
             },
             {
               "Timestamp": "2015-05-01T01:00:00Z",
               "Value": "CRITICAL",
               "Summary": "Cream status is CRITICAL",
               "Message": "Cream job submission test failed"
             },
             {
               "Timestamp": "2015-05-01T05:00:00Z",
               "Value": "OK",
               "Summary": "Cream status is ok",
               "Message": "Cream job submission test return value of ok"
             }
           ]
         }
       ]
     }
   ]
 }`

	request, _ := http.NewRequest("GET", "/api/v2/metric_result/cream01.afroditi.gr?exec_time=2015-05-01T00:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(respJSON, response.Body.String(), "Response body mismatch")

}

func (suite *metricResultTestSuite) TestFilters() {

	type TestReq struct {
		Path string
		Code int
		JSON string
	}
	expected := []TestReq{
		{
			Path: "/api/v2/metric_result/cream01.afroditi.gr?exec_time=2020-05-01T05:00:00Z&service=CREAM-CE",
			Code: 200,
			JSON: `{
   "root": [
     {
       "Name": "cream01.afroditi.gr",
       "info": {
         "URL": "http://creamce.example.com"
       },
       "Metrics": [
         {
           "Name": "emi.cream.CREAMCE-JobSubmit",
           "Service": "CREAM-CE",
           "Details": [
             {
               "Timestamp": "2020-05-01T05:00:00Z",
               "Value": "OK",
               "Summary": "Cream status is ok",
               "Message": "Cream job submission test return value of ok"
             }
           ]
         }
       ]
     }
   ]
 }`},
		{
			Path: "/api/v2/metric_result/cream01.afroditi.gr?exec_time=2020-05-01T05:00:00Z&service=CREAM-CE2",
			Code: 200,
			JSON: `{
   "root": [
     {
       "Name": "cream01.afroditi.gr",
       "info": {
         "URL": "http://creamce.example.com"
       },
       "Metrics": [
         {
           "Name": "emi.cream.CREAMCE-JobSubmit",
           "Service": "CREAM-CE2",
           "Details": [
             {
               "Timestamp": "2020-05-01T05:00:00Z",
               "Value": "OK",
               "Summary": "Cream status is ok",
               "Message": "Cream job submission test return value of ok"
             }
           ]
         }
       ]
     }
   ]
 }`},

		{
			Path: "/api/v2/metric_result/cream01.afroditi.gr/emi.cream.CREAMCE-JobSubmit?exec_time=2020-05-01T05:00:00Z&service=CREAM-CE2",
			Code: 200,
			JSON: `{
   "root": [
     {
       "Name": "cream01.afroditi.gr",
       "info": {
         "URL": "http://creamce.example.com"
       },
       "Metrics": [
         {
           "Name": "emi.cream.CREAMCE-JobSubmit",
           "Service": "CREAM-CE2",
           "Details": [
             {
               "Timestamp": "2020-05-01T05:00:00Z",
               "Value": "OK",
               "Summary": "Cream status is ok",
               "Message": "Cream job submission test return value of ok"
             }
           ]
         }
       ]
     }
   ]
 }`},
	}

	for _, exp := range expected {
		request, _ := http.NewRequest("GET", exp.Path, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		code := response.Code
		output := response.Body.String()

		// Check that we must have a 200 ok code
		suite.Equal(exp.Code, code, "Response Code Mismatch on call:"+exp.Path)
		// Compare the expected and actual json response
		suite.Equal(exp.JSON, output, "Response body mismatch on call:"+exp.Path)
	}

}

func (suite *metricResultTestSuite) TestMultipleMetricResultsNonOK() {

	respJSON := `{
   "root": [
     {
       "Name": "cream01.afroditi.gr",
       "info": {
         "URL": "http://creamce.example.com"
       },
       "Metrics": [
         {
           "Name": "emi.cream.CREAMCE-JobSubmit",
           "Service": "CREAM-CE",
           "Details": [
             {
               "Timestamp": "2015-05-01T01:00:00Z",
               "Value": "CRITICAL",
               "Summary": "Cream status is CRITICAL",
               "Message": "Cream job submission test failed"
             }
           ]
         }
       ]
     }
   ]
 }`

	request, _ := http.NewRequest("GET", "/api/v2/metric_result/cream01.afroditi.gr?exec_time=2015-05-01T00:00:00Z&filter=non-ok", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(respJSON, response.Body.String(), "Response body mismatch")

}

func (suite *metricResultTestSuite) TestMultipleMetricResultsCritical() {

	respJSON := `{
   "root": [
     {
       "Name": "cream01.afroditi.gr",
       "info": {
         "URL": "http://creamce.example.com"
       },
       "Metrics": [
         {
           "Name": "emi.cream.CREAMCE-JobSubmit",
           "Service": "CREAM-CE",
           "Details": [
             {
               "Timestamp": "2015-05-01T01:00:00Z",
               "Value": "CRITICAL",
               "Summary": "Cream status is CRITICAL",
               "Message": "Cream job submission test failed"
             }
           ]
         }
       ]
     }
   ]
 }`

	request, _ := http.NewRequest("GET", "/api/v2/metric_result/cream01.afroditi.gr?exec_time=2015-05-01T00:00:00Z&filter=critical", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(respJSON, response.Body.String(), "Response body mismatch")

}

func (suite *metricResultTestSuite) TestMultipleMetricResultsWarning() {

	respJSON := `{
   "root": []
 }`

	request, _ := http.NewRequest("GET", "/api/v2/metric_result/cream01.afroditi.gr?exec_time=2015-05-01T00:00:00Z&filter=warning", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(respJSON, response.Body.String(), "Response body mismatch")

}

func (suite *metricResultTestSuite) TestMultipleMetricResultsOK() {

	respJSON := `{
   "root": [
     {
       "Name": "cream01.afroditi.gr",
       "info": {
         "URL": "http://creamce.example.com"
       },
       "Metrics": [
         {
           "Name": "emi.cream.CREAMCE-JobSubmit",
           "Service": "CREAM-CE",
           "Details": [
             {
               "Timestamp": "2015-05-01T00:00:00Z",
               "Value": "OK",
               "Summary": "Cream status is ok",
               "Message": "Cream job submission test return value of ok",
               "actual_data": "latency=15s",
               "threshold_rule_applied": "latency=15s;20:30;30:40",
               "original_status": "WARNING"
             },
             {
               "Timestamp": "2015-05-01T05:00:00Z",
               "Value": "OK",
               "Summary": "Cream status is ok",
               "Message": "Cream job submission test return value of ok"
             }
           ]
         }
       ]
     }
   ]
 }`

	request, _ := http.NewRequest("GET", "/api/v2/metric_result/cream01.afroditi.gr?exec_time=2015-05-01T00:00:00Z&filter=ok", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(respJSON, response.Body.String(), "Response body mismatch")

}

// TestOptionsMetricResult is used to test the OPTIONS response
func (suite *metricResultTestSuite) TestOptionsMetricResult() {

	request, _ := http.NewRequest("OPTIONS", "/api/v2/metric_result/cream01.afroditi.gr/emi.cream.CREAMCE-JobSubmit", strings.NewReader(""))
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
func (suite *metricResultTestSuite) TearDownTest() {

	suite.cfg.MongoClient.Database("ARGO_test_metric_result").Drop(context.TODO())
	suite.cfg.MongoClient.Database("ARGO_test_metric_result_egi").Drop(context.TODO())

}

// This is the first function called when go test is issued
func TestSuiteMetricResult(t *testing.T) {
	suite.Run(t, new(metricResultTestSuite))
}
