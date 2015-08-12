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

package statusMetrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"code.google.com/p/gcfg"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/authentication"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// This is a util. suite struct used in tests (see pkg "testify")
type StatusMetricsTestSuite struct {
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
func (suite *StatusMetricsTestSuite) SetupTest() {

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
    db = "argotest_metrics"
`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(true).PathPrefix("/api/v2/status").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)

	// Connect to mongo testdb
	session, _ := mongo.OpenSession(suite.cfg.MongoDB)

	// Add authentication token to mongo testdb
	seedAuth := bson.M{"api_key": "S3CR3T"}
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "authentication", seedAuth)

	// seed mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// seed a tenant to use
	c := session.DB(suite.cfg.MongoDB.Db).C("tenants")
	c.Insert(bson.M{
		"name": "EGI",
		"db_conf": []bson.M{
			bson.M{
				"store":    "main",
				"server":   "localhost",
				"port":     27017,
				"database": "argotest_metrics_egi",
				"username": "",
				"password": ""},
		},
		"users": []bson.M{
			bson.M{
				"name":    "egi_user",
				"email":   "egi_user@email.com",
				"api_key": "KEY1"},
		}})

	c.Insert(bson.M{
		"name": "EUDAT",
		"db_conf": []bson.M{
			bson.M{
				"store":    "main",
				"server":   "localhost",
				"port":     27017,
				"database": "argotest_metrics_eudat",
				"username": "",
				"password": ""},
		},
		"users": []bson.M{
			bson.M{
				"name":    "eudat_user",
				"email":   "eudat_user@email.com",
				"api_key": "KEY2"},
		}})

	// get dbconfiguration based on the tenant
	// Prepare the request object
	request, _ := http.NewRequest("GET", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY1")
	// authenticate user's api key and find corresponding tenant
	suite.tenantDbConf, err = authentication.AuthenticateTenant(request.Header, suite.cfg)

	// Now seed the report DEFINITIONS
	c = session.DB(suite.tenantDbConf.Db).C("reports")
	c.Insert(bson.M{
		"name":            "ROC_CRITICAL",
		"tenant":          "EGI",
		"endpoint_group":  "SITES",
		"group_of_groups": "NGI",
		"profiles": []bson.M{
			bson.M{
				"name":  "metric",
				"value": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			bson.M{
				"name":  "name1",
				"value": "value1"},
			bson.M{
				"name":  "name2",
				"value": "value2"},
		}})

	// seed the status detailed metric data
	c = session.DB(suite.tenantDbConf.Db).C("status_metrics")
	c.Insert(bson.M{
		"report":         "ROC_CRITICAL",
		"date_int":       20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"hostname":       "cream01.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "OK",
		"prev_timestamp": "2015-04-30T22:00:00Z",
		"prev_status":    "WARNING",
	})
	c.Insert(bson.M{
		"report":         "ROC_CRITICAL",
		"date_int":       20150501,
		"timestamp":      "2015-05-01T01:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"hostname":       "cream01.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "CRITICAL",
		"prev_timestamp": "2015-05-01T00:00:00Z",
		"prev_status":    "CRITICAL",
	})
	c.Insert(bson.M{
		"report":         "ROC_CRITICAL",
		"date_int":       20150501,
		"timestamp":      "2015-05-01T05:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"hostname":       "cream01.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "OK",
		"prev_timestamp": "2015-05-01T01:00:00Z",
		"prev_status":    "CRITICAL",
	})

	// get dbconfiguration based on the tenant
	// Prepare the request object
	request, _ = http.NewRequest("GET", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY2")
	// authenticate user's api key and find corresponding tenant
	suite.tenantDbConf, err = authentication.AuthenticateTenant(request.Header, suite.cfg)

	// Now seed the reports DEFINITIONS
	c = session.DB(suite.tenantDbConf.Db).C("reports")
	c.Insert(bson.M{
		"name":            "EUDAT_CRITICAL",
		"tenant":          "EUDAT",
		"endpoint_group":  "EUDAT_SITES",
		"group_of_groups": "EUDAT_GROUPS",
		"profiles": []bson.M{
			bson.M{
				"name":  "metric",
				"value": "eudat.CRITICAL"},
		},
		"filter_tags": []bson.M{
			bson.M{
				"name":  "name1",
				"value": "value1"},
			bson.M{
				"name":  "name2",
				"value": "value2"},
		}})

	// seed the status detailed metric data
	c = session.DB(suite.tenantDbConf.Db).C("status_metrics")
	c.Insert(bson.M{
		"report":         "EUDAT_CRITICAL",
		"date_int":       20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "EL-01-AUTH",
		"service":        "srv.typeA",
		"hostname":       "host01.eudat.gr",
		"metric":         "typeA.metric.Memory",
		"status":         "OK",
		"prev_timestamp": "2015-04-30T22:00:00Z",
		"prev_status":    "WARNING",
	})
	c.Insert(bson.M{
		"report":         "EUDAT_CRITICAL",
		"date_int":       20150501,
		"timestamp":      "2015-05-01T01:00:00Z",
		"endpoint_group": "EL-01-AUTH",
		"service":        "srv.typeA",
		"hostname":       "host01.eudat.gr",
		"metric":         "typeA.metric.Memory",
		"status":         "CRITICAL",
		"prev_timestamp": "2015-05-01T00:00:00Z",
		"prev_status":    "OK",
	})
	c.Insert(bson.M{
		"report":         "EUDAT_CRITICAL",
		"date_int":       20150501,
		"timestamp":      "2015-05-01T05:00:00Z",
		"endpoint_group": "EL-01-AUTH",
		"service":        "srv.typeA",
		"hostname":       "host01.eudat.gr",
		"metric":         "typeA.metric.Memory",
		"status":         "OK",
		"prev_timestamp": "2015-05-01T01:00:00Z",
		"prev_status":    "WARNING",
	})

}

func (suite *StatusMetricsTestSuite) TestListStatusMetrics() {
	respXML1 := ` <root>
   <group name="HG-03-AUTH" type="SITES">
     <group name="CREAM-CE" type="service">
       <group name="cream01.afroditi.gr" type="endpoint">
         <group name="emi.cream.CREAMCE-JobSubmit" type="metric">
           <status timestamp="2015-04-30T22:00:00Z" value="WARNING"></status>
           <status timestamp="2015-05-01T00:00:00Z" value="OK"></status>
           <status timestamp="2015-05-01T01:00:00Z" value="CRITICAL"></status>
           <status timestamp="2015-05-01T05:00:00Z" value="OK"></status>
         </group>
       </group>
     </group>
   </group>
 </root>`

	respXML2 := ` <root>
   <group name="EL-01-AUTH" type="EUDAT_SITES">
     <group name="srv.typeA" type="service">
       <group name="host01.eudat.gr" type="endpoint">
         <group name="typeA.metric.Memory" type="metric">
           <status timestamp="2015-04-30T22:00:00Z" value="WARNING"></status>
           <status timestamp="2015-05-01T00:00:00Z" value="OK"></status>
           <status timestamp="2015-05-01T01:00:00Z" value="CRITICAL"></status>
           <status timestamp="2015-05-01T05:00:00Z" value="OK"></status>
         </group>
       </group>
     </group>
   </group>
 </root>`

	respJSON1 := `{
   "endpoint_groups": [
     {
       "name": "HG-03-AUTH",
       "type": "SITES",
       "services": [
         {
           "name": "CREAM-CE",
           "type": "service",
           "endpoints": [
             {
               "name": "cream01.afroditi.gr",
               "type": "endpoint",
               "metrics": [
                 {
                   "name": "emi.cream.CREAMCE-JobSubmit",
                   "type": "metric",
                   "statuses": [
                     {
                       "timestamp": "2015-04-30T22:00:00Z",
                       "value": "WARNING"
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
   "endpoint_groups": [
     {
       "name": "EL-01-AUTH",
       "type": "EUDAT_SITES",
       "services": [
         {
           "name": "srv.typeA",
           "type": "service",
           "endpoints": [
             {
               "name": "host01.eudat.gr",
               "type": "endpoint",
               "metrics": [
                 {
                   "name": "typeA.metric.Memory",
                   "type": "metric",
                   "statuses": [
                     {
                       "timestamp": "2015-04-30T22:00:00Z",
                       "value": "WARNING"
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

	fullurl1 := "/api/v2/status/ROC_CRITICAL/SITES/HG-03-AUTH" +
		"/services/CREAM-CE/endpoints/cream01.afroditi.gr/metrics/emi.cream.CREAMCE-JobSubmit" +
		"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z"

	fullurl2 := "/api/v2/status/EUDAT_CRITICAL/EUDAT_SITES/EL-01-AUTH" +
		"/services/srv.typeA/endpoints/host01.eudat.gr/metrics/typeA.metric.Memory" +
		"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z"

	// 1. EGI XML REQUEST
	// init the response placeholder
	response := httptest.NewRecorder()
	// Prepare the request object for fist tenant
	request, _ := http.NewRequest("GET", fullurl1, strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json")
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
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json")
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
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json")
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
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json")
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

}

// This function is actually called in the end of all tests
// and clears the test environment.
// Mainly it's purpose is to drop the testdb
func (suite *StatusMetricsTestSuite) TearDownTest() {

	session, _ := mongo.OpenSession(suite.cfg.MongoDB)

	session.DB("argotest_metrics").DropDatabase()
	session.DB("argotest_metrics_eudat").DropDatabase()
	session.DB("argotest_metrics_egi").DropDatabase()
}

// This is the first function called when go test is issued
func TestStatusMetricsSuite(t *testing.T) {
	suite.Run(t, new(StatusMetricsTestSuite))
}
