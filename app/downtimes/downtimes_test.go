package downtimes

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gcfg.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// This is a util. suite struct used in tests (see pkg "testify")
type DowntimesTestSuite struct {
	suite.Suite
	cfg                       config.Config
	router                    *mux.Router
	confHandler               respond.ConfHandler
	tenantDbConf              config.MongoConfig
	clientkey                 string
	respRecomputationsCreated string
	respUnauthorized          string
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

	log.SetOutput(ioutil.Discard)

	// seed mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Seed database with tenants
	c := session.DB(suite.cfg.MongoDB.Db).C("tenants")
	c.Insert(
		bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"info": bson.M{
				"name":    "GUARDIANS",
				"email":   "email@something2",
				"website": "www.gotg.com",
				"created": "2015-10-20 02:08:04",
				"updated": "2015-10-20 02:08:04"},
			"db_conf": []bson.M{

				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_FOO",
				},
				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_FOO",
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "user1",
					"email":   "user1@email.com",
					"api_key": "USER1KEY",
					"roles":   []string{"editor"},
				},
				bson.M{
					"name":    "user2",
					"email":   "user2@email.com",
					"api_key": "USER2KEY",
					"roles":   []string{"editor"},
				},
			}})
	c.Insert(
		bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50d",
			"info": bson.M{
				"name":    "AVENGERS",
				"email":   "email@something2",
				"website": "www.gotg.com",
				"created": "2015-10-20 02:08:04",
				"updated": "2015-10-20 02:08:04"},
			"db_conf": []bson.M{

				bson.M{
					// "store":    "ar",
					"server":   suite.tenantDbConf.Host,
					"port":     suite.tenantDbConf.Port,
					"database": suite.tenantDbConf.Db,
					"username": suite.tenantDbConf.Username,
					"password": suite.tenantDbConf.Password,
				},
				bson.M{
					"server":   suite.tenantDbConf.Host,
					"port":     suite.tenantDbConf.Port,
					"database": suite.tenantDbConf.Db,
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "user3",
					"email":   "user3@email.com",
					"api_key": suite.clientkey,
					"roles":   []string{"editor"},
				},
				bson.M{
					"name":    "user4",
					"email":   "user4@email.com",
					"api_key": "VIEWERKEY",
					"roles":   []string{"viewer"},
				},
			}})

	c = session.DB(suite.cfg.MongoDB.Db).C("roles")
	c.Insert(
		bson.M{
			"resource": "downtimes.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "downtimes.get",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "downtimes.create",
			"roles":    []string{"editor"},
		})
	c.Insert(
		bson.M{
			"resource": "downtimes.delete",
			"roles":    []string{"editor"},
		})
	c.Insert(
		bson.M{
			"resource": "downtimes.update",
			"roles":    []string{"editor"},
		})

	// Seed database with downtimes
	c = session.DB(suite.tenantDbConf.Db).C("downtimes")
	c.EnsureIndexKey("-date_integer", "id")
	c.Insert(
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			"date_integer": 20191011,
			"date":         "2019-10-11",
			"name":         "Critical",
			"endpoints": []bson.M{
				bson.M{"hostname": "host-A", "service": "service-A", "start_time": "2019-10-11T04:00:33Z", "end_time": "2019-10-11T15:33:00Z"},
				bson.M{"hostname": "host-B", "service": "service-B", "start_time": "2019-10-11T12:00:33Z", "end_time": "2019-10-11T12:33:00Z"},
				bson.M{"hostname": "host-C", "service": "service-C", "start_time": "2019-10-11T20:00:33Z", "end_time": "2019-10-11T22:15:00Z"},
			},
		})
	c.Insert(
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			"date_integer": 20191012,
			"date":         "2019-10-12",
			"name":         "Critical",
			"endpoints": []bson.M{
				bson.M{"hostname": "host-A", "service": "service-A", "start_time": "2019-10-12T04:00:33Z", "end_time": "2019-10-12T15:33:00Z"},
				bson.M{"hostname": "host-B", "service": "service-B", "start_time": "2019-10-12T12:00:33Z", "end_time": "2019-10-12T12:33:00Z"},
				bson.M{"hostname": "host-C", "service": "service-C", "start_time": "2019-10-12T20:00:33Z", "end_time": "2019-10-12T22:15:00Z"},
			},
		})
	c.Insert(
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			"date_integer": 20191013,
			"date":         "2019-10-13",
			"name":         "Critical",
			"endpoints": []bson.M{
				bson.M{"hostname": "host-A", "service": "service-A", "start_time": "2019-10-13T04:00:33Z", "end_time": "2019-10-13T15:33:00Z"},
				bson.M{"hostname": "host-B", "service": "service-B", "start_time": "2019-10-13T12:00:33Z", "end_time": "2019-10-13T12:33:00Z"},
				bson.M{"hostname": "host-C", "service": "service-C", "start_time": "2019-10-13T20:00:33Z", "end_time": "2019-10-13T22:15:00Z"},
			},
		})
	c.Insert(
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"date_integer": 20191011,
			"date":         "2019-10-11",
			"name":         "NonCritical",
			"endpoints": []bson.M{
				bson.M{"hostname": "host-01", "service": "service-01", "start_time": "2019-10-11T02:00:33Z", "end_time": "2019-10-11T23:33:00Z"},
				bson.M{"hostname": "host-02", "service": "service-02", "start_time": "2019-10-11T16:00:33Z", "end_time": "2019-10-11T16:45:00Z"},
			},
		})
	c.Insert(
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"date_integer": 20191012,
			"date":         "2019-10-12",
			"name":         "NonCritical",
			"endpoints":    []bson.M{},
		})
	c.Insert(
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"date_integer": 20191013,
			"date":         "2019-10-13",
			"name":         "NonCritical",
			"endpoints": []bson.M{
				bson.M{"hostname": "host-01", "service": "service-01", "start_time": "2019-10-13T02:00:33Z", "end_time": "2019-10-13T23:33:00Z"},
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
		reqHeader{Method: "GET", Path: "/api/v2/downtimes?date=2020-02", Data: ""},
		reqHeader{Method: "GET", Path: "/api/v2/downtimes/some-uuid?date=2020-02", Data: ""},
		reqHeader{Method: "POST", Path: "/api/v2/downtimes?date=2020-02", Data: ""},
		reqHeader{Method: "PUT", Path: "/api/v2/downtimes/some-id?date=2020-02", Data: ""},
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
	{"hostname":"new-host-foo","service":"service-new-foo","start_time":"2019-10-11T23:10:00Z","end_time":"2019-10-11T23:20:00Z"},
	{"hostname":"new-host-bar","service":"service-new-bar","start_time":"2019-10-11T23:40:00Z","end_time":"2019-10-11T23:50:00Z"}
  ]
 }`

	jsonOutput := `{
 "status": {
  "message": "Downtimes resource succesfully created",
  "code": "201"
 },
 "data": {
  "id": "{{id}}",
  "links": {
   "self": "https:///api/v2/downtimes/{{id}}"
  }
 }
}`

	jsonCreated := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "{{id}}",
   "date": "2019-11-29",
   "name": "downtimes_set",
   "endpoints": [
    {
     "hostname": "new-host-foo",
     "service": "service-new-foo",
     "start_time": "2019-10-11T23:10:00Z",
     "end_time": "2019-10-11T23:20:00Z"
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

	// Grab id from mongodb
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	defer session.Close()
	if err != nil {
		panic(err)
	}
	// Retrieve id from database
	var result map[string]interface{}
	c := session.DB(suite.tenantDbConf.Db).C("downtimes")
	c.Find(bson.M{"name": "downtimes_set"}).One(&result)
	id := result["id"].(string)

	// Apply id to output template and check
	suite.Equal(strings.Replace(jsonOutput, "{{id}}", id, 2), output, "Response body mismatch")

	// Check that actually the item has been created
	// Call List one with the specific id
	request2, _ := http.NewRequest("GET", "/api/v2/downtimes/"+id, strings.NewReader(jsonInput))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)

	code2 := response2.Code
	output2 := response2.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(200, code2, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(strings.Replace(jsonCreated, "{{id}}", id, 2), output2, "Response body mismatch")

}

func (suite *DowntimesTestSuite) TestList() {

	request, _ := http.NewRequest("GET", "/api/v2/downtimes", strings.NewReader(""))
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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-10-13",
   "name": "Critical",
   "endpoints": [
    {
     "hostname": "host-A",
     "service": "service-A",
     "start_time": "2019-10-13T04:00:33Z",
     "end_time": "2019-10-13T15:33:00Z"
    },
    {
     "hostname": "host-B",
     "service": "service-B",
     "start_time": "2019-10-13T12:00:33Z",
     "end_time": "2019-10-13T12:33:00Z"
    },
    {
     "hostname": "host-C",
     "service": "service-C",
     "start_time": "2019-10-13T20:00:33Z",
     "end_time": "2019-10-13T22:15:00Z"
    }
   ]
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "date": "2019-10-13",
   "name": "NonCritical",
   "endpoints": [
    {
     "hostname": "host-01",
     "service": "service-01",
     "start_time": "2019-10-13T02:00:33Z",
     "end_time": "2019-10-13T23:33:00Z"
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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-10-11",
   "name": "Critical",
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
     "end_time": "2019-10-11T22:15:00Z"
    }
   ]
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "date": "2019-10-11",
   "name": "NonCritical",
   "endpoints": [
    {
     "hostname": "host-01",
     "service": "service-01",
     "start_time": "2019-10-11T02:00:33Z",
     "end_time": "2019-10-11T23:33:00Z"
    },
    {
     "hostname": "host-02",
     "service": "service-02",
     "start_time": "2019-10-11T16:00:33Z",
     "end_time": "2019-10-11T16:45:00Z"
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

func (suite *DowntimesTestSuite) TestOptionsdowntimes() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/downtimes", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, POST, DELETE, PUT, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	headers = response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, POST, DELETE, PUT, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

}

//TearDownTest to tear down every test
func (suite *DowntimesTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}

	tenantDB := session.DB(suite.tenantDbConf.Db)
	mainDB := session.DB(suite.cfg.MongoDB.Db)

	cols, err := tenantDB.CollectionNames()
	for _, col := range cols {
		tenantDB.C(col).RemoveAll(nil)
	}

	cols, err = mainDB.CollectionNames()
	for _, col := range cols {
		mainDB.C(col).RemoveAll(nil)
	}

}

func (suite *DowntimesTestSuite) TestListOneNotFound() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Not Found",
  "code": "404"
 },
 "errors": [
  {
   "message": "Not Found",
   "code": "404",
   "details": "item with the specific ID was not found on the server"
  }
 ]
}`

	request, _ := http.NewRequest("GET", "/api/v2/downtimes/wrong-id", strings.NewReader(jsonInput))
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

func (suite *DowntimesTestSuite) TestListOne() {

	request, _ := http.NewRequest("GET", "/api/v2/downtimes/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-10-13",
   "name": "Critical",
   "endpoints": [
    {
     "hostname": "host-A",
     "service": "service-A",
     "start_time": "2019-10-13T04:00:33Z",
     "end_time": "2019-10-13T15:33:00Z"
    },
    {
     "hostname": "host-B",
     "service": "service-B",
     "start_time": "2019-10-13T12:00:33Z",
     "end_time": "2019-10-13T12:33:00Z"
    },
    {
     "hostname": "host-C",
     "service": "service-C",
     "start_time": "2019-10-13T20:00:33Z",
     "end_time": "2019-10-13T22:15:00Z"
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

func (suite *DowntimesTestSuite) TestUpdateNameAlreadyExists() {

	jsonInput := `{
  
   "name": "Critical",
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
     "end_time": "2019-10-11T22:15:00Z"
    }
   ]
  }`

	jsonOutput := `{
 "status": {
  "message": "Conflict",
  "code": "409"
 },
 "errors": [
  {
   "message": "Conflict",
   "code": "409",
   "details": "Downtimes resource with the same name already exists"
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/downtimes/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(409, code)
	suite.Equal(jsonOutput, output)

}

func (suite *DowntimesTestSuite) TestUpdateBadJson() {

	jsonInput := `{
		"name": "downtimes_set",
		"endpoints": [
	
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

	request, _ := http.NewRequest("PUT", "/api/v2/downtimes/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
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

func (suite *DowntimesTestSuite) TestUpdateNotFound() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Not Found",
  "code": "404"
 },
 "errors": [
  {
   "message": "Not Found",
   "code": "404",
   "details": "item with the specific ID was not found on the server"
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/downtimes/wrong-id", strings.NewReader(jsonInput))
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

func (suite *DowntimesTestSuite) TestUpdate() {

	jsonInput := `{
   "name": "downtimes_set",
   "endpoints": [
	{"hostname":"updated-host-foo","service":"service-new-foo","start_time":"2019-11-30T23:10:00Z","end_time":"2019-11-30T23:25:00Z"},
	{"hostname":"updated-host-bar","service":"service-new-bar","start_time":"2019-11-30T23:40:00Z","end_time":"2019-11-30T23:55:00Z"}
  ]
 }`

	jsonOutput := `{
 "status": {
  "message": "Downtimes resource successfully updated (new history snapshot)",
  "code": "200"
 }
}`

	jsonUpdated := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "date": "2019-11-30",
   "name": "downtimes_set",
   "endpoints": [
    {
     "hostname": "updated-host-foo",
     "service": "service-new-foo",
     "start_time": "2019-11-30T23:10:00Z",
     "end_time": "2019-11-30T23:25:00Z"
    },
    {
     "hostname": "updated-host-bar",
     "service": "service-new-bar",
     "start_time": "2019-11-30T23:40:00Z",
     "end_time": "2019-11-30T23:55:00Z"
    }
   ]
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/downtimes/6ac7d684-1f8e-4a02-a502-720e8f11e50c?date=2019-11-30", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response

	// Apply id to output template and check
	suite.Equal(jsonOutput, output, "Response body mismatch")

	// Check that the item has actually updated
	// run a list specific
	request2, _ := http.NewRequest("GET", "/api/v2/downtimes/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)

	code2 := response2.Code
	output2 := response2.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(200, code2, "Internal Server Error")
	// Compare the expected and actual json response

	suite.Equal(jsonUpdated, output2, "Response body mismatch")

}

func (suite *DowntimesTestSuite) TestDeleteNotFound() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Not Found",
  "code": "404"
 },
 "errors": [
  {
   "message": "Not Found",
   "code": "404",
   "details": "item with the specific ID was not found on the server"
  }
 ]
}`

	request, _ := http.NewRequest("DELETE", "/api/v2/downtimes/wrong-id", strings.NewReader(jsonInput))
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

	request, _ := http.NewRequest("DELETE", "/api/v2/downtimes/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Downtimes resource successfully deleted",
  "code": "200"
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricProfileJSON, output, "Response body mismatch")

	// check that the element has actually been Deleted
	// connect to mongodb
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	defer session.Close()
	if err != nil {
		panic(err)
	}
	// try to retrieve item
	var result map[string]interface{}
	c := session.DB(suite.tenantDbConf.Db).C("downtimes")
	err = c.Find(bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"}).One(&result)

	suite.NotEqual(err, nil, "No not found error")
	suite.Equal(err.Error(), "not found", "No not found error")
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

func (suite *DowntimesTestSuite) TestUpdateForbidViewer() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Access to the resource is Forbidden",
  "code": "403"
 }
}`

	request, _ := http.NewRequest("PUT", "/api/v2/downtimes/6ac7d684-1f8e-4a02-a502-720e8f11e007", strings.NewReader(jsonInput))
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

	request, _ := http.NewRequest("DELETE", "/api/v2/downtimes/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
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

//TearDownTest to tear down every test
func (suite *DowntimesTestSuite) TearDownSuite() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

func TestDowntimesTestSuite(t *testing.T) {
	suite.Run(t, new(DowntimesTestSuite))
}
