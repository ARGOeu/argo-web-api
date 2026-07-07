package resultsV5

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

type endpointGroupAvailabilityTestSuite struct {
	suite.Suite
	cfg          config.Config
	router       *mux.Router
	confHandler  respond.ConfHandler
	tenantDbConf config.MongoConfig
	clientkey    string
}

// Setup the Test Environment
func (suite *endpointGroupAvailabilityTestSuite) SetupSuite() {

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
	    db = "ARGO_test_endpointGroup_availability"
	    `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.tenantDbConf.Db = "ARGO_test_endpointGroup_availability_tenant"
	suite.tenantDbConf.Password = "h4shp4ss"
	suite.tenantDbConf.Username = "dbuser"
	suite.tenantDbConf.Store = "ar"
	suite.clientkey = "secretkey"

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v5/results").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)
}

// This function runs before any test and setups the environment
func (suite *endpointGroupAvailabilityTestSuite) SetupTest() {

	log.SetOutput(io.Discard)

	// Seed database with tenants
	//TODO: move tests to
	c := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("tenants")
	c.InsertOne(context.TODO(),
		bson.M{"name": "TENANT-AA",
			"db_conf": []bson.M{
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_test1",
				},
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_test2",
				},
			},
			"users": []bson.M{
				{
					"name":    "alice",
					"email":   "alice@tenant-a.foo",
					"api_key": "alice_key",
					"roles":   []string{"viewer"},
				},
				{
					"name":    "bob",
					"email":   "bob@tenant-a.foo",
					"api_key": "bob_key",
					"roles":   []string{"viewer"},
				},
			}})
	c.InsertOne(context.TODO(),
		bson.M{"name": "TENANT-BB",
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
					"database": "tenant_b_db",
				},
			},
			"users": []bson.M{
				{
					"name":    "alice",
					"email":   "alice@tenant-b.foo",
					"api_key": suite.clientkey,
					"roles":   []string{"viewer"},
				},
				{
					"name":    "bob",
					"email":   "bob@tenant-b.foo",
					"api_key": "bob_key",
					"roles":   []string{"viewer"},
				},
			}})

	c = suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("roles")
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "results.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "results.get",
			"roles":    []string{"editor", "viewer"},
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

// TestEndpointGroupAvailability test if daily results are returned correctly
func (suite *endpointGroupAvailabilityTestSuite) TestListGroupAvailability() {

	request, _ := http.NewRequest("GET", "/api/v5/results/Report_A/groups/ST01?start-time=2015-06-20T12:00:00Z&end-time=2015-06-23T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	endpointGroupAvailabilityJSON := `{
   "results": [
     {
       "name": "GROUP_A",
       "type": "GROUP",
       "groups": [
         {
           "name": "ST01",
           "type": "SITES",
           "results": [
             {
               "timestamp": "2015-06-22",
               "availability": 66.7,
               "reliability": 54.6,
               "unknown": 0,
               "uptime": 1,
               "downtime": 0
             },
             {
               "timestamp": "2015-06-23",
               "availability": 100,
               "reliability": 100,
               "unknown": 0,
               "uptime": 1,
               "downtime": 0
             }
           ]
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointGroupAvailabilityJSON, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v5/results/Report_A/groups/ST01?start-time=2015-06-20T12:00:00Z&end-time=2015-06-23T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", "AWRONGKEY")
	request.Header.Set("Accept", "application/json")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	unauthorizedresponseXML := `{
 "status": {
  "message": "Unauthorized",
  "code": "401",
  "details": "You need to provide a correct authentication token using the header 'x-api-key'"
 }
}`

	// Check that we must have a 401 Unauthorized code
	suite.Equal(401, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(unauthorizedresponseXML, response.Body.String(), "Response body mismatch")

}

// TestListAllEndpointGroupAvailability test if daily results are returned correctly
func (suite *endpointGroupAvailabilityTestSuite) TestListAllEndpointGroupAvailability() {

	request, _ := http.NewRequest("GET", "/api/v5/results/Report_A/groups?start-time=2015-06-20T12:00:00Z&end-time=2015-06-23T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	endpointGroupAvailabilityJSON := `{
   "results": [
     {
       "name": "GROUP_A",
       "type": "GROUP",
       "groups": [
         {
           "name": "ST01",
           "type": "SITES",
           "results": [
             {
               "timestamp": "2015-06-22",
               "availability": 66.7,
               "reliability": 54.6,
               "unknown": 0,
               "uptime": 1,
               "downtime": 0
             },
             {
               "timestamp": "2015-06-23",
               "availability": 100,
               "reliability": 100,
               "unknown": 0,
               "uptime": 1,
               "downtime": 0
             }
           ]
         },
         {
           "name": "ST02",
           "type": "SITES",
           "results": [
             {
               "timestamp": "2015-06-22",
               "availability": 70,
               "reliability": 45,
               "unknown": 0,
               "uptime": 1,
               "downtime": 0
             },
             {
               "timestamp": "2015-06-23",
               "availability": 43.5,
               "reliability": 56,
               "unknown": 0,
               "uptime": 1,
               "downtime": 0
             }
           ]
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointGroupAvailabilityJSON, response.Body.String(), "Response body mismatch")

}

// TestOptionsEndpointGroup tests responses in case the OPTIONS http verb is used
func (suite *endpointGroupAvailabilityTestSuite) TestOptionsEndpointGroup() {

	request, _ := http.NewRequest("OPTIONS", "/api/v5/results/Report_A/groups", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v5/results/Report_A/groups/ST01", strings.NewReader(""))

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	headers = response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v5/results/Report_A/groups", strings.NewReader(""))

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	headers = response.Result().Header

}

// TestListAllEndpointGroupAvailabilityCustom tests if a/r results are returned correctly over a custom period
func (suite *endpointGroupAvailabilityTestSuite) TestListAllEndpointGroupAvailabilityCustom() {

	request, _ := http.NewRequest("GET", "/api/v5/results/Report_A/groups?start-time=2015-06-20T12:00:00Z&end-time=2015-06-23T23:00:00Z&granularity=custom", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	endpointGroupAvailabilityJSON := `{
   "results": [
     {
       "name": "GROUP_A",
       "type": "GROUP",
       "groups": [
         {
           "name": "ST01",
           "type": "SITES",
           "results": [
             {
               "availability": 100,
               "reliability": 100,
               "unknown": 0,
               "uptime": 1,
               "downtime": 0
             }
           ]
         },
         {
           "name": "ST02",
           "type": "SITES",
           "results": [
             {
               "availability": 100,
               "reliability": 100,
               "unknown": 0,
               "uptime": 1,
               "downtime": 0
             }
           ]
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(endpointGroupAvailabilityJSON, response.Body.String(), "Response body mismatch")

}

// TearDownTest to tear down every test
func (suite *endpointGroupAvailabilityTestSuite) TearDownTest() {

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
func (suite *endpointGroupAvailabilityTestSuite) TearDownSuite() {
	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
}

// TestSuiteResultsEndpointGroup is responsible for calling the tests
func TestSuiteResultsEndpointGroup(t *testing.T) {
	suite.Run(t, new(endpointGroupAvailabilityTestSuite))
}
