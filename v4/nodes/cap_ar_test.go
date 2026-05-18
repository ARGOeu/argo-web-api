package nodes

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/store"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/gcfg.v1"
)

type CapAvailabilityTestSuite struct {
	suite.Suite
	cfg             config.Config
	router          *mux.Router
	confHandler     respond.ConfHandler
	tenantDbConf    config.MongoConfig
	tenantpassword  string
	tenantusername  string
	tenantstorename string
	clientkey       string
}

// Setup the Test Environment
func (suite *CapAvailabilityTestSuite) SetupSuite() {

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
		 db = "ARGO_test_node_ar"
		 `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.tenantDbConf.Db = "ARGO_test_arV3_tenant"
	suite.tenantDbConf.Password = "h4shp4ss"
	suite.tenantDbConf.Username = "dbuser"
	suite.tenantDbConf.Store = "ar"
	suite.clientkey = "secretkey"

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v4").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)
}

// This function runs before any test and setups the environment
func (suite *CapAvailabilityTestSuite) SetupTest() {

	log.SetOutput(io.Discard)

	// Seed database with tenants
	//TODO: move tests to
	c := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("tenants")
	c.InsertOne(context.TODO(),
		bson.M{"id": "NODEA-id",
			"info": bson.M{"name": "NODEA"},
			"db_conf": []bson.M{
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Tenant1",
				},
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Tenant2",
				},
			},
			"node": true,
			"users": []bson.M{
				{
					"name":    "bob",
					"email":   "bob@foo",
					"api_key": "random-1",
					"roles":   []string{"viewer"},
				},
				{
					"name":    "alice",
					"email":   "alice@foo",
					"api_key": "random-2",
					"roles":   []string{"viewer"},
				},
			}})
	c.InsertOne(context.TODO(),
		bson.M{"id": "NODEB-id",
			"info": bson.M{"name": "NODEB"},
			"db_conf": []bson.M{
				{
					"server":   "localhost",
					"port":     27017,
					"database": suite.tenantDbConf.Db,
					"username": suite.tenantDbConf.Username,
					"password": suite.tenantDbConf.Password,
				},
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_wrong_db_endpointgrouavailability",
				},
			},
			"node": true,
			"users": []bson.M{
				{
					"name":    "john",
					"email":   "john@foo",
					"api_key": suite.clientkey,
					"roles":   []string{"viewer"},
				},
				{
					"name":    "george",
					"email":   "george@foo",
					"api_key": "random-2",
					"roles":   []string{"viewer"},
				},
			}})

	c = suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("roles")
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "v4.nodes.summary.item",
			"roles":    []string{"super_admin", "admin", "editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "v4.nodes.summary",
			"roles":    []string{"super_admin", "admin", "editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "v4.nodes.availability.item",
			"roles":    []string{"super_admin", "admin", "editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "v4.nodes.availability",
			"roles":    []string{"super_admin", "admin", "editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "v4.nodes.uptime.item",
			"roles":    []string{"super_admin", "admin", "editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "v4.nodes.uptime",
			"roles":    []string{"super_admin", "admin", "editor", "viewer"},
		})

	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection(nodeReportColName)
	c.InsertOne(context.TODO(), bson.M{
		"_id":       "node-report",
		"report_id": "eba61a9e-22e9-4521-9e47-ecaa4a49436",
	})

	// Seed tenant database with data
	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("endpoint_group_ar")

	// Insert seed data
	c.InsertOne(context.TODO(),
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "ST01",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 66.7,
			"reliability":  54.6,
			"weight":       5634,
			"tags": []bson.M{
				{
					"name":  "",
					"value": "",
				},
			},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "ST02",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 70,
			"reliability":  45,
			"weight":       4356,
			"tags": []bson.M{
				bson.M{
					"name":  "",
					"value": "",
				},
			},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150623,
			"name":         "ST01",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 100,
			"reliability":  100,
			"weight":       5634,
			"tags": []bson.M{
				{
					"name":  "",
					"value": "",
				},
			},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150623,
			"name":         "ST02",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 43.5,
			"reliability":  56,
			"weight":       4356,
			"tags": []bson.M{
				{
					"name":  "",
					"value": "",
				},
			},
		})

	// Seed endpoint data
	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("endpoint_ar")

	// Insert seed data
	c.InsertOne(context.TODO(),
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "host01",
			"service":      "service01",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 66.7,
			"reliability":  54.6,
			"weight":       5634,
			"tags": []bson.M{
				{
					"name":  "",
					"value": "",
				},
			},
			"info": bson.M{"ID": "special-queue"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150623,
			"name":         "host01",
			"service":      "service01",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 100,
			"reliability":  100,
			"weight":       5634,
			"tags": []bson.M{
				{
					"name":  "",
					"value": "",
				},
			},
			"info": bson.M{"ID": "special-queue"},
		})

	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("reports")

	c.InsertOne(context.TODO(), bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a49436",
		"info": bson.M{
			"name":        "Report_A",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "GROUP",
				"group": bson.M{
					"type": "SITES",
				},
			},
		},
		"profiles": []bson.M{
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"name":  "name1",
				"value": "value1"},
			{
				"name":  "name2",
				"value": "value2"},
		}})

}

func (suite *CapAvailabilityTestSuite) TestSummary() {

	request, _ := http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/summary?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	endpointGroupAvailabilityA := `{
   "data": [
     {
       "name": "ST01",
       "results": [
         {
           "date": "2015-06-22",
           "availability": "66.7",
           "uptime": "1"
         },
         {
           "date": "2015-06-23",
           "availability": "100",
           "uptime": "1"
         }
       ]
     },
     {
       "name": "ST02",
       "results": [
         {
           "date": "2015-06-22",
           "availability": "70",
           "uptime": "1"
         },
         {
           "date": "2015-06-23",
           "availability": "43.5",
           "uptime": "1"
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointGroupAvailabilityA, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/summary?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	monthlyAvailJSON := `{
   "data": [
     {
       "name": "ST01",
       "results": [
         {
           "date": "2015-06",
           "availability": "99.99999900000002",
           "uptime": "1"
         }
       ]
     },
     {
       "name": "ST02",
       "results": [
         {
           "date": "2015-06",
           "availability": "99.99999900000002",
           "uptime": "1"
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(monthlyAvailJSON, response.Body.String(), "Response body mismatch")

	// check by start and end date
	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/summary?start_date=2015-06-20&end_date=2015-06-23", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointGroupAvailabilityA, response.Body.String(), "Response body mismatch")

	oneItemJSON := `{
   "data": [
     {
       "name": "ST01",
       "results": [
         {
           "date": "2015-06-22",
           "availability": "66.7",
           "uptime": "1"
         }
       ]
     }
   ]
 }`

	// check by date
	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/summary/ST01?date=2015-06-22", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(oneItemJSON, response.Body.String(), "Response body mismatch")

	oneDateJSON := `{
   "data": [
     {
       "name": "ST01",
       "results": [
         {
           "date": "2015-06-22",
           "availability": "66.7",
           "uptime": "1"
         }
       ]
     },
     {
       "name": "ST02",
       "results": [
         {
           "date": "2015-06-22",
           "availability": "70",
           "uptime": "1"
         }
       ]
     }
   ]
 }`

	// check by date
	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/summary?date=2015-06-22", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(oneDateJSON, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/summary?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", "AWRONGKEY")
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	unauthorizedresponse := `{
 "status": {
  "message": "Unauthorized",
  "code": "401",
  "details": "You need to provide a correct authentication token using the header 'x-api-key'"
 }
}`

	// Check that we must have a 401 Unauthorized code
	suite.Equal(401, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(unauthorizedresponse, response.Body.String(), "Response body mismatch")

}

func (suite *CapAvailabilityTestSuite) TestAvailability() {

	request, _ := http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/availability?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	endpointGroupAvailabilityA := `{
   "data": [
     {
       "name": "ST01",
       "results": [
         {
           "date": "2015-06-22",
           "availability": "66.7"
         },
         {
           "date": "2015-06-23",
           "availability": "100"
         }
       ]
     },
     {
       "name": "ST02",
       "results": [
         {
           "date": "2015-06-22",
           "availability": "70"
         },
         {
           "date": "2015-06-23",
           "availability": "43.5"
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointGroupAvailabilityA, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/availability?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	monthlyAvailJSON := `{
   "data": [
     {
       "name": "ST01",
       "results": [
         {
           "date": "2015-06",
           "availability": "99.99999900000002"
         }
       ]
     },
     {
       "name": "ST02",
       "results": [
         {
           "date": "2015-06",
           "availability": "99.99999900000002"
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(monthlyAvailJSON, response.Body.String(), "Response body mismatch")

	// check by start and end date
	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/availability?start_date=2015-06-20&end_date=2015-06-23", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointGroupAvailabilityA, response.Body.String(), "Response body mismatch")

	oneItemJSON := `{
   "data": [
     {
       "name": "ST02",
       "results": [
         {
           "date": "2015-06-22",
           "availability": "70"
         }
       ]
     }
   ]
 }`

	// check by date
	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/availability/ST02?date=2015-06-22", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(oneItemJSON, response.Body.String(), "Response body mismatch")

	oneDateJSON := `{
   "data": [
     {
       "name": "ST01",
       "results": [
         {
           "date": "2015-06-22",
           "availability": "66.7"
         }
       ]
     },
     {
       "name": "ST02",
       "results": [
         {
           "date": "2015-06-22",
           "availability": "70"
         }
       ]
     }
   ]
 }`

	// check by date
	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/availability?date=2015-06-22", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(oneDateJSON, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/availability?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", "AWRONGKEY")
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	unauthorizedresponse := `{
 "status": {
  "message": "Unauthorized",
  "code": "401",
  "details": "You need to provide a correct authentication token using the header 'x-api-key'"
 }
}`

	// Check that we must have a 401 Unauthorized code
	suite.Equal(401, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(unauthorizedresponse, response.Body.String(), "Response body mismatch")

}

func (suite *CapAvailabilityTestSuite) TestAvailabilityOptions() {

	request, _ := http.NewRequest("OPTIONS", "/api/v4/nodes/NODEB/capabilities/availability", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

}

func (suite *CapAvailabilityTestSuite) TestAvailabilityItemOptions() {

	request, _ := http.NewRequest("OPTIONS", "/api/v4/nodes/NODEB/capabilities/availability/ST01", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

}

func (suite *CapAvailabilityTestSuite) TestUptime() {

	request, _ := http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/uptime?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	endpointGroupUptimeA := `{
   "data": [
     {
       "name": "ST01",
       "results": [
         {
           "date": "2015-06-22",
           "uptime": "1"
         },
         {
           "date": "2015-06-23",
           "uptime": "1"
         }
       ]
     },
     {
       "name": "ST02",
       "results": [
         {
           "date": "2015-06-22",
           "uptime": "1"
         },
         {
           "date": "2015-06-23",
           "uptime": "1"
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointGroupUptimeA, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/uptime?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	monthlyAvailJSON := `{
   "data": [
     {
       "name": "ST01",
       "results": [
         {
           "date": "2015-06",
           "uptime": "1"
         }
       ]
     },
     {
       "name": "ST02",
       "results": [
         {
           "date": "2015-06",
           "uptime": "1"
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(monthlyAvailJSON, response.Body.String(), "Response body mismatch")

	// check by start and end date
	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/uptime?start_date=2015-06-20&end_date=2015-06-23", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointGroupUptimeA, response.Body.String(), "Response body mismatch")

	oneItemJSON := `{
   "data": [
     {
       "name": "ST01",
       "results": [
         {
           "date": "2015-06-22",
           "uptime": "1"
         }
       ]
     }
   ]
 }`

	// check by date
	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/uptime/ST01?date=2015-06-22", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(oneItemJSON, response.Body.String(), "Response body mismatch")

	oneDateJSON := `{
   "data": [
     {
       "name": "ST01",
       "results": [
         {
           "date": "2015-06-22",
           "uptime": "1"
         }
       ]
     },
     {
       "name": "ST02",
       "results": [
         {
           "date": "2015-06-22",
           "uptime": "1"
         }
       ]
     }
   ]
 }`

	// check by date
	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/uptime?date=2015-06-22", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(oneDateJSON, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v4/nodes/NODEB/capabilities/uptime?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", "AWRONGKEY")
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	unauthorizedresponse := `{
 "status": {
  "message": "Unauthorized",
  "code": "401",
  "details": "You need to provide a correct authentication token using the header 'x-api-key'"
 }
}`

	// Check that we must have a 401 Unauthorized code
	suite.Equal(401, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(unauthorizedresponse, response.Body.String(), "Response body mismatch")

}

func (suite *CapAvailabilityTestSuite) TestUptimeOptions() {

	request, _ := http.NewRequest("OPTIONS", "/api/v4/nodes/NODEB/capabilities/uptime", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

}

func (suite *CapAvailabilityTestSuite) TestUptimeItemOptions() {

	request, _ := http.NewRequest("OPTIONS", "/api/v4/nodes/NODEB/capabilities/uptime/ST01", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

}

// TearDownTest to tear down every test
func (suite *CapAvailabilityTestSuite) TearDownTest() {

	mainDB := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db)
	cols, err := mainDB.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		panic(err)
	}

	for _, col := range cols {
		mainDB.Collection(col).Drop(context.TODO())
	}

	tenantDB := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db)
	cols, err = tenantDB.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		panic(err)
	}

	for _, col := range cols {
		tenantDB.Collection(col).Drop(context.TODO())
	}

}

// TearDownTest to tear down every test
func (suite *CapAvailabilityTestSuite) TearDownSuite() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
}

// CapTestSuiteAR is responsible for calling the tests
func CapTestSuiteAR(t *testing.T) {
	suite.Run(t, new(CapAvailabilityTestSuite))
}
