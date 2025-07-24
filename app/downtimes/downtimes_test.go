package downtimes

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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/gcfg.v1"
)

// This is a util. suite struct used in tests (see pkg "testify")
type DowntimesTestSuite struct {
	suite.Suite
	cfg              config.Config
	router           *mux.Router
	confHandler      respond.ConfHandler
	tenantDbConf     config.MongoConfig
	clientkey        string
	respUnauthorized string
}

// Setup the Test Environment
func (suite *DowntimesTestSuite) SetupSuite() {
	const testConfig = `
	 [server]
	 bindip = ""
	 port = 8080
	 maxprocs = 4
	 cache = false
	 lrucache = 700000000
	 gzip = true
	 reqsizelimit = 1073741824
 
	 [mongodb]
	 host = "127.0.0.1"
	 port = 27017
	 db = "AR_test_downtimes"
	 `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.respUnauthorized = "Unauthorized"
	suite.tenantDbConf = config.MongoConfig{
		Host:     "localhost",
		Port:     27017,
		Db:       "AR_test_downtimes_tenant",
		Password: "pass",
		Username: "dbuser",
		Store:    "ar",
	}
	suite.clientkey = "123456"

	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v2").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)
}

// This function runs before any test and setups the environment
func (suite *DowntimesTestSuite) SetupTest() {

	log.SetOutput(io.Discard)

	// Seed database with tenants
	c := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("tenants")
	c.InsertOne(context.TODO(),
		bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"info": bson.M{
				"name":    "GUARDIANS",
				"email":   "email@something2",
				"website": "www.gotg.com",
				"created": "2015-10-20 02:08:04",
				"updated": "2015-10-20 02:08:04"},
			"db_conf": []bson.M{
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_FOO",
				},
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_FOO",
				},
			},
			"users": []bson.M{
				{
					"name":    "user1",
					"email":   "user1@email.com",
					"api_key": "USER1KEY",
					"roles":   []string{"editor"},
				},
				{
					"name":    "user2",
					"email":   "user2@email.com",
					"api_key": "USER2KEY",
					"roles":   []string{"editor"},
				},
			}})
	c.InsertOne(context.TODO(),
		bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50d",
			"info": bson.M{
				"name":    "AVENGERS",
				"email":   "email@something2",
				"website": "www.gotg.com",
				"created": "2015-10-20 02:08:04",
				"updated": "2015-10-20 02:08:04"},
			"db_conf": []bson.M{

				{
					// "store":    "ar",
					"server":   suite.tenantDbConf.Host,
					"port":     suite.tenantDbConf.Port,
					"database": suite.tenantDbConf.Db,
					"username": suite.tenantDbConf.Username,
					"password": suite.tenantDbConf.Password,
				},
				{
					"server":   suite.tenantDbConf.Host,
					"port":     suite.tenantDbConf.Port,
					"database": suite.tenantDbConf.Db,
				},
			},
			"users": []bson.M{

				{
					"name":    "user3",
					"email":   "user3@email.com",
					"api_key": suite.clientkey,
					"roles":   []string{"editor"},
				},
				{
					"name":    "user4",
					"email":   "user4@email.com",
					"api_key": "VIEWERKEY",
					"roles":   []string{"viewer"},
				},
			}})

	c = suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("roles")
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "downtimes.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "downtimes.get",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "downtimes.create",
			"roles":    []string{"editor"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "downtimes.delete",
			"roles":    []string{"editor"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "downtimes.update",
			"roles":    []string{"editor"},
		})

	// Seed database with downtimes
	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("downtimes")
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "date_integer", Value: -1},
			{Key: "id", Value: 1},
		},
		Options: options.Index().SetUnique(false),
	}
	c.Indexes().CreateOne(context.TODO(), indexModel)
	c.InsertOne(context.TODO(),
		bson.M{

			"date_integer": 20191011,
			"date":         "2019-10-11",
			"name":         "Critical",
			"endpoints": []bson.M{
				{"hostname": "host-A", "service": "service-A", "start_time": "2019-10-11T04:00:33Z", "end_time": "2019-10-11T15:33:00Z"},
				{"hostname": "host-B", "service": "service-B", "start_time": "2019-10-11T12:00:33Z", "end_time": "2019-10-11T12:33:00Z"},
				{"hostname": "host-C",
					"service":     "service-C",
					"start_time":  "2019-10-11T20:00:33Z",
					"end_time":    "2019-10-11T22:15:00Z",
					"description": "a simple description"},
			},
		})
	c.InsertOne(context.TODO(),
		bson.M{

			"date_integer": 20191012,
			"date":         "2019-10-12",
			"name":         "Critical",
			"endpoints": []bson.M{
				{"hostname": "host-A", "service": "service-A", "start_time": "2019-10-12T04:00:33Z", "end_time": "2019-10-12T15:33:00Z",
					"classification": "unscheduled", "severity": "warning"},
				{"hostname": "host-B", "service": "service-B", "start_time": "2019-10-12T12:00:33Z", "end_time": "2019-10-12T12:33:00Z",
					"classification": "unscheduled", "severity": "outage"},
				{"hostname": "host-C", "service": "service-C", "start_time": "2019-10-12T20:00:33Z", "end_time": "2019-10-12T22:15:00Z",
					"classification": "scheduled", "severity": "warning"},
			},
		})
	c.InsertOne(context.TODO(),
		bson.M{

			"date_integer": 20191013,
			"date":         "2019-10-13",
			"name":         "Critical",
			"endpoints": []bson.M{
				{"hostname": "host-A", "service": "service-A", "start_time": "2019-10-13T04:00:33Z", "end_time": "2019-10-13T15:33:00Z"},
				{"hostname": "host-B", "service": "service-B", "start_time": "2019-10-13T12:00:33Z", "end_time": "2019-10-13T12:33:00Z"},
				{"hostname": "host-C", "service": "service-C", "start_time": "2019-10-13T20:00:33Z", "end_time": "2019-10-13T22:15:00Z"},
			},
		})

}

func (suite *DowntimesTestSuite) TestCreateBadJson() {

	jsonInput := `{
   "weight_type":"hepsec",
   "group_type": "SITES",
   "groups": [
	 {
	   "name": "SITE-A",
		"value": 33.33
	 `

	jsonOutput := `{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "Bad Request",
   "code": "400",
   "details": "Request Body contains malformed JSON, thus rendering the Request Bad"
  }
 ]
}`

	request, _ := http.NewRequest("POST", "/api/v2/downtimes", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(400, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *DowntimesTestSuite) TestBadDate() {

	badDate := `{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "Bad Request",
   "code": "400",
   "details": "date parameter value: 2020-02 is not in the valid form of YYYY-MM-DD"
  }
 ]
}`

	type reqHeader struct {
		Method string
		Path   string
		Data   string
	}

	requests := []reqHeader{
		{Method: "GET", Path: "/api/v2/downtimes?date=2020-02", Data: ""},
		{Method: "POST", Path: "/api/v2/downtimes?date=2020-02", Data: ""},
	}

	for _, r := range requests {
		request, _ := http.NewRequest(r.Method, r.Path, strings.NewReader(r.Data))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		code := response.Code
		output := response.Body.String()

		// Check that we must have a 200 ok code
		suite.Equal(400, code, "Internal Server Error")
		// Compare the expected and actual json response
		suite.Equal(badDate, output, "Response body mismatch")

	}

}
func (suite *DowntimesTestSuite) TestCreate() {

	jsonInput := `{
   "name": "downtimes_set",
   "endpoints": [
	{ "hostname":"new-host-foo",
	  "service":"service-new-foo",
	  "start_time":"2019-10-11T23:10:00Z",
	  "end_time":"2019-10-11T23:20:00Z",
	  "description": "this downtime has severity and classification fields defined",
	  "severity": "warning",
	  "classification": "unscheduled"
	},
	{"hostname":"new-host-bar","service":"service-new-bar","start_time":"2019-10-11T23:40:00Z","end_time":"2019-10-11T23:50:00Z"}
  ]
 }`

	jsonOutput := `{
 "status": {
  "message": "Downtimes set created for date: 2019-11-29",
  "code": "201"
 }
}`

	jsonCreated := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2019-11-29",
   "endpoints": [
    {
     "hostname": "new-host-foo",
     "service": "service-new-foo",
     "start_time": "2019-10-11T23:10:00Z",
     "end_time": "2019-10-11T23:20:00Z",
     "description": "this downtime has severity and classification fields defined",
     "classification": "unscheduled",
     "severity": "warning"
    },
    {
     "hostname": "new-host-bar",
     "service": "service-new-bar",
     "start_time": "2019-10-11T23:40:00Z",
     "end_time": "2019-10-11T23:50:00Z"
    }
   ]
  }
 ]
}`

	request, _ := http.NewRequest("POST", "/api/v2/downtimes?date=2019-11-29", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(201, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(jsonOutput, output, "Resonse Body Mismatch")

	// Call List one with the specific id
	request2, _ := http.NewRequest("GET", "/api/v2/downtimes?date=2019-11-29", strings.NewReader(jsonInput))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)

	code2 := response2.Code
	output2 := response2.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(200, code2, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(jsonCreated, output2, "Response body mismatch")

}

func (suite *DowntimesTestSuite) TestListNotFound() {

	request, _ := http.NewRequest("GET", "/api/v2/downtimes?date=2020-05-05", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	downtimesJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2020-05-05",
   "endpoints": []
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(downtimesJSON, output, "Response body mismatch")

}

func (suite *DowntimesTestSuite) TestListPast() {

	request, _ := http.NewRequest("GET", "/api/v2/downtimes?date=2019-10-11", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	downtimesJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2019-10-11",
   "endpoints": [
    {
     "hostname": "host-A",
     "service": "service-A",
     "start_time": "2019-10-11T04:00:33Z",
     "end_time": "2019-10-11T15:33:00Z"
    },
    {
     "hostname": "host-B",
     "service": "service-B",
     "start_time": "2019-10-11T12:00:33Z",
     "end_time": "2019-10-11T12:33:00Z"
    },
    {
     "hostname": "host-C",
     "service": "service-C",
     "start_time": "2019-10-11T20:00:33Z",
     "end_time": "2019-10-11T22:15:00Z",
     "description": "a simple description"
    }
   ]
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(downtimesJSON, output, "Response body mismatch")

}

func (suite *DowntimesTestSuite) TestListFilter1() {

	request, _ := http.NewRequest("GET", "/api/v2/downtimes?date=2019-10-12", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	downtimesJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2019-10-12",
   "endpoints": [
    {
     "hostname": "host-A",
     "service": "service-A",
     "start_time": "2019-10-12T04:00:33Z",
     "end_time": "2019-10-12T15:33:00Z",
     "classification": "unscheduled",
     "severity": "warning"
    },
    {
     "hostname": "host-B",
     "service": "service-B",
     "start_time": "2019-10-12T12:00:33Z",
     "end_time": "2019-10-12T12:33:00Z",
     "classification": "unscheduled",
     "severity": "outage"
    },
    {
     "hostname": "host-C",
     "service": "service-C",
     "start_time": "2019-10-12T20:00:33Z",
     "end_time": "2019-10-12T22:15:00Z",
     "classification": "scheduled",
     "severity": "warning"
    }
   ]
  }
 ]
}`

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(downtimesJSON, output, "Response body mismatch")

}

func (suite *DowntimesTestSuite) TestListFilter2() {

	request, _ := http.NewRequest("GET", "/api/v2/downtimes?date=2019-10-12&severity=outage", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	downtimesJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2019-10-12",
   "endpoints": [
    {
     "hostname": "host-B",
     "service": "service-B",
     "start_time": "2019-10-12T12:00:33Z",
     "end_time": "2019-10-12T12:33:00Z",
     "classification": "unscheduled",
     "severity": "outage"
    }
   ]
  }
 ]
}`

	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(downtimesJSON, output, "Response body mismatch")

}

func (suite *DowntimesTestSuite) TestListFilter3() {

	request, _ := http.NewRequest("GET", "/api/v2/downtimes?date=2019-10-12&classification=scheduled", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	downtimesJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2019-10-12",
   "endpoints": [
    {
     "hostname": "host-C",
     "service": "service-C",
     "start_time": "2019-10-12T20:00:33Z",
     "end_time": "2019-10-12T22:15:00Z",
     "classification": "scheduled",
     "severity": "warning"
    }
   ]
  }
 ]
}`

	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(downtimesJSON, output, "Response body mismatch")

}

func (suite *DowntimesTestSuite) TestListFilter4() {

	request, _ := http.NewRequest("GET", "/api/v2/downtimes?date=2019-10-12&classification=unscheduled&severity=warning", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	downtimesJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2019-10-12",
   "endpoints": [
    {
     "hostname": "host-A",
     "service": "service-A",
     "start_time": "2019-10-12T04:00:33Z",
     "end_time": "2019-10-12T15:33:00Z",
     "classification": "unscheduled",
     "severity": "warning"
    }
   ]
  }
 ]
}`

	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(downtimesJSON, output, "Response body mismatch")

}

func (suite *DowntimesTestSuite) TestOptionsdowntimes() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/downtimes", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, POST, DELETE, PUT, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	headers = response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, POST, DELETE, PUT, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

}

// TearDownTest to tear down every test
func (suite *DowntimesTestSuite) TearDownTest() {

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

func (suite *DowntimesTestSuite) TestDeleteNotFound() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Downtimes dataset not found for date: 2020-02-11",
  "code": "404"
 }
}`

	request, _ := http.NewRequest("DELETE", "/api/v2/downtimes?date=2020-02-11", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(404, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *DowntimesTestSuite) TestDelete() {

	request, _ := http.NewRequest("DELETE", "/api/v2/downtimes?date=2019-10-11", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Downtimes set deleted for date: 2019-10-11",
  "code": "200"
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricProfileJSON, output, "Response body mismatch")

	// check that the element has actually been Deleted

	// try to retrieve item
	c := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("downtimes")
	queryResult := c.FindOne(context.TODO(), bson.M{"date_integer": 20191011})
	suite.NotEqual(queryResult.Err(), nil, "No not found error")
	suite.Equal(queryResult.Err(), mongo.ErrNoDocuments, "No not found error")
}

func (suite *DowntimesTestSuite) TestCreateDateConflict() {

	jsonInput := `{
		"name": "downtimes_set",
		"endpoints": [
		 {"hostname":"new-host-foo","service":"service-new-foo","start_time":"2019-10-11T23:10:00Z","end_time":"2019-10-11T23:20:00Z", "classification": "severe"},
		 {"hostname":"new-host-bar","service":"service-new-bar","start_time":"2019-10-11T23:40:00Z","end_time":"2019-10-11T23:50:00Z"},
		 {"hostname":"new-host-bar","service":"service-new-bar","start_time":"2019-10-11T23:40:00Z","end_time":"2019-10-11T23:50:00Z", "description": "simple"}
	   ]
	  }`

	jsonOutput := `{
 "status": {
  "message": "Downtimes already exists for date: 2019-10-11, please either update it or delete it first!",
  "code": "409"
 }
}`

	request, _ := http.NewRequest("POST", "/api/v2/downtimes?date=2019-10-11", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 409 code
	suite.Equal(409, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")
}

func (suite *DowntimesTestSuite) TestCreateForbidViewer() {

	jsonInput := `{
   "name": "test_profile",
   "namespace [
	 `

	jsonOutput := `{
 "status": {
  "message": "Access to the resource is Forbidden",
  "code": "403"
 }
}`

	request, _ := http.NewRequest("POST", "/api/v2/downtimes", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", "VIEWERKEY")
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(403, code, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonOutput, output, "Response body mismatch")
}

func (suite *DowntimesTestSuite) TestDeleteForbidViewer() {

	request, _ := http.NewRequest("DELETE", "/api/v2/downtimes?date=2019-10-11", strings.NewReader(""))
	request.Header.Set("x-api-key", "VIEWERKEY")
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Access to the resource is Forbidden",
  "code": "403"
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(403, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricProfileJSON, output, "Response body mismatch")
}

// TearDownTest to tear down every test
func (suite *DowntimesTestSuite) TearDownSuite() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
}

func TestSuiteDowntimes(t *testing.T) {
	suite.Run(t, new(DowntimesTestSuite))
}
