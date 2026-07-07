package resultsV5

import (
	"context"
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

type SuperGroupAvailabilityTestSuite struct {
	suite.Suite
	cfg          config.Config
	router       *mux.Router
	confHandler  respond.ConfHandler
	tenantDbConf config.MongoConfig
	clientkey    string
}

// Setup the Test Environment
func (suite *SuperGroupAvailabilityTestSuite) SetupSuite() {

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
	    db = "ARGO_test_SuperGroup_availability"
	    `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.tenantDbConf.Db = "ARGO_test_SuperGroup_availability_tenant"
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
func (suite *SuperGroupAvailabilityTestSuite) SetupTest() {

	// Seed database with tenants
	//TODO: move tests to
	c := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("tenants")
	c.InsertOne(context.TODO(),
		bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"info": bson.M{
				"name":    "TENANT-AA",
				"email":   "email@tenant-a.foo",
				"website": "www.tenant-a.foo",
				"created": "2015-10-20 02:08:04",
				"updated": "2015-10-20 02:08:04"},
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
					"name":    "Bob",
					"email":   "bob@tenant-a.foo",
					"api_key": "bob_key",
					"roles":   []string{"viewer"},
				},
				{
					"name":    "Alice",
					"email":   "alice@tenant-a.foo",
					"api_key": "alice_key",
					"roles":   []string{"viewer"},
				},
			}})
	c.InsertOne(context.TODO(),
		bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50d",
			"info": bson.M{
				"name":    "TENANT-BB",
				"email":   "email@tenant-b.foo",
				"website": "www.tenant-b.foo",
				"created": "2015-10-20 02:08:04",
				"updated": "2015-10-20 02:08:04"},
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
			"users": []bson.M{
				{
					"name":    "bob",
					"email":   "bob@tenant-b.foo",
					"api_key": suite.clientkey,
					"roles":   []string{"viewer"},
				},
				{
					"name":    "alice",
					"email":   "alice@tenant-b.foo",
					"api_key": "alice_secret_2",
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
	// Seed database with recomputations
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
			"name":         "ST04",
			"supergroup":   "GROUP_B",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 30,
			"reliability":  100,
			"weight":       5344,
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
			"name":         "ST05",
			"supergroup":   "GROUP_B",
			"up":           0,
			"down":         0,
			"unknown":      1,
			"availability": 90,
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
			"date":         20150624,
			"name":         "ST05",
			"supergroup":   "GROUP_B",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 40,
			"reliability":  70,
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
			"date":         20150625,
			"name":         "ST05",
			"supergroup":   "GROUP_B",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 40,
			"reliability":  70,
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
					"type": "SITE",
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

// TestListSuperGroupAvailability test if daily results are returned correctly
func (suite *SuperGroupAvailabilityTestSuite) TestListSuperGroupAvailability() {

	request, _ := http.NewRequest("GET", "/api/v5/results/Report_A/supergroups/GROUP_A?start-time=2015-06-20T12:00:00Z&end-time=2015-06-26T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	SuperGrouAvailabilityJSON := `{
   "results": [
     {
       "name": "GROUP_A",
       "type": "GROUP",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": 68.14,
           "reliability": 50.41
         },
         {
           "timestamp": "2015-06-23",
           "availability": 75.36,
           "reliability": 80.81
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(SuperGrouAvailabilityJSON, response.Body.String(), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v5/results/Report_A/supergroups/GROUP_A?start-time=2015-06-20T12:00:00Z&end-time=2015-06-26T23:00:00Z", strings.NewReader(""))
	request.Header.Set("x-api-key", "AWRONGKEY")
	request.Header.Set("Accept", "application/xml")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	// Check that we must have a 401 Unauthorized code
	suite.Equal(401, response.Code, "Incorrect HTTP response code")

}

// TestListAllSuperGroupAvailability test if daily results are returned correctly
func (suite *SuperGroupAvailabilityTestSuite) TestListAllSuperGroupAvailability() {

	request, _ := http.NewRequest("GET", "/api/v5/results/Report_A/supergroups?start-time=2015-06-20T12:00:00Z&end-time=2015-06-26T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	output := response.Body.String()

	SuperGroupAvailabilityJSON := `{
   "results": [
     {
       "name": "GROUP_A",
       "type": "GROUP",
       "results": [
         {
           "timestamp": "2015-06-22",
           "availability": 68.14,
           "reliability": 50.41
         },
         {
           "timestamp": "2015-06-23",
           "availability": 75.36,
           "reliability": 80.81
         }
       ]
     },
     {
       "name": "GROUP_B",
       "type": "GROUP",
       "results": [
         {
           "timestamp": "2015-06-23",
           "availability": 60.79,
           "reliability": 100
         },
         {
           "timestamp": "2015-06-24",
           "availability": 40,
           "reliability": 70
         },
         {
           "timestamp": "2015-06-25",
           "availability": 40,
           "reliability": 70
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(SuperGroupAvailabilityJSON, output, "Response body mismatch")

}

// TestListAllSuperGroupAvailabilityMonthly test if results are returned correctly for a monthly period
func (suite *SuperGroupAvailabilityTestSuite) TestListAllSuperGroupAvailabilityMonthly() {

	request, _ := http.NewRequest("GET", "/api/v5/results/Report_A/supergroups?start-time=2015-06-20T12:00:00Z&end-time=2015-06-26T23:00:00Z&granularity=monthly", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	SuperGroupAvailabilityJSON := `{
   "results": [
     {
       "name": "GROUP_A",
       "type": "GROUP",
       "results": [
         {
           "timestamp": "2015-06",
           "availability": 71.75,
           "reliability": 65.61
         }
       ]
     },
     {
       "name": "GROUP_B",
       "type": "GROUP",
       "results": [
         {
           "timestamp": "2015-06",
           "availability": 46.93,
           "reliability": 80
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(SuperGroupAvailabilityJSON, response.Body.String(), "Response body mismatch")

}

// TestListAllSuperGroupAvailabilityCustom test if results are returned correctly for a custom period
func (suite *SuperGroupAvailabilityTestSuite) TestListAllSuperGroupAvailabilityCustom() {

	request, _ := http.NewRequest("GET", "/api/v5/results/Report_A/supergroups?start-time=2015-06-20T12:00:00Z&end-time=2015-06-26T23:00:00Z&granularity=custom", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	SuperGroupAvailabilityJSON := `{
   "results": [
     {
       "name": "GROUP_A",
       "type": "GROUP",
       "results": [
         {
           "availability": 71.75,
           "reliability": 65.61
         }
       ]
     },
     {
       "name": "GROUP_B",
       "type": "GROUP",
       "results": [
         {
           "availability": 46.93,
           "reliability": 80
         }
       ]
     }
   ]
 }`

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(SuperGroupAvailabilityJSON, response.Body.String(), "Response body mismatch")

}

// TestListSuperGroupAvailabilityErrors tests if errors/exceptions are returned correctly
func (suite *SuperGroupAvailabilityTestSuite) TestListSuperGroupAvailabilityErrors() {

	typeError1JSON := `{
   "message": "No results found for given query",
   "code": 404
 }`

	request, _ := http.NewRequest("GET", "/api/v5/results/Report_A/supergroups/lala?start-time=2025-06-22T00:00:00Z&end-time=2025-06-23T23:59:59Z", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	//request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	output := response.Body.String()

	// Check that we must have a 404 not found code
	suite.Equal(404, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(typeError1JSON, output, "Response body mismatch")

}

func (suite *SuperGroupAvailabilityTestSuite) TestOptionsSuperGroup() {
	request, _ := http.NewRequest("OPTIONS", "/api/v5/results/Report_A/supergroups", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v5/results/Report_A/supergroups/GROUP_A", strings.NewReader(""))

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

// TestStrictSlashSuperGroupResults test if not found responses are returned correctly
func (suite *SuperGroupAvailabilityTestSuite) TestStrictSlashSuperGroupResults() {

	request, _ := http.NewRequest("GET", "/api/v2/results/Report_A/supergroups/?start-time=2015-06-20T12:00:00Z&end-time=2015-06-26T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")
	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	suite.Equal(404, response.Code, "Incorrect HTTP response code")

	request, _ = http.NewRequest("GET", "/api/v2/results/Report_A/supergroups/GROUP_A/?start-time=2015-06-20T12:00:00Z&end-time=2015-06-26T23:00:00Z&granularity=daily", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/xml")
	response = httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)
	suite.Equal(404, response.Code, "Incorrect HTTP response code")

}

// TearDownTest to tear down every test
func (suite *SuperGroupAvailabilityTestSuite) TearDownTest() {

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
func (suite *SuperGroupAvailabilityTestSuite) TearDownSuite() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
}

// TestRecompuptationsTestSuite is responsible for calling the tests
func TestSuperGroupResultsTestSuite(t *testing.T) {
	suite.Run(t, new(SuperGroupAvailabilityTestSuite))
}
