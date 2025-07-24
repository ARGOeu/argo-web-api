package weights

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
type WeightsTestSuite struct {
	suite.Suite
	cfg              config.Config
	router           *mux.Router
	confHandler      respond.ConfHandler
	tenantDbConf     config.MongoConfig
	clientkey        string
	respUnauthorized string
}

// Setup the Test Environment
func (suite *WeightsTestSuite) SetupSuite() {
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
	 db = "AR_test_weights"
	 `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.respUnauthorized = "Unauthorized"
	suite.tenantDbConf = config.MongoConfig{
		Host:     "localhost",
		Port:     27017,
		Db:       "AR_test_weights_tenant",
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
func (suite *WeightsTestSuite) SetupTest() {

	log.SetOutput(io.Discard)

	// seed mongo
	client := suite.cfg.MongoClient

	// Seed database with tenants
	c := client.Database(suite.cfg.MongoDB.Db).Collection("tenants")
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

	c = client.Database(suite.cfg.MongoDB.Db).Collection("roles")
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "weights.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "weights.get",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "weights.create",
			"roles":    []string{"editor"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "weights.delete",
			"roles":    []string{"editor"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "weights.update",
			"roles":    []string{"editor"},
		})

	// Seed database with weights
	c = client.Database(suite.tenantDbConf.Db).Collection("weights")
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
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			"name":         "Critical",
			"date_integer": 20191004,
			"date":         "2019-10-04",
			"weight_type":  "hepsepc",
			"group_type":   "SITES",
			"groups": []bson.M{
				{"name": "SITE-A", "value": 1673},
				{"name": "SITE-B", "value": 1234},
				{"name": "SITE-C", "value": 523},
				{"name": "SITE-D", "value": 2},
			},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			"name":         "Critical",
			"date_integer": 20191023,
			"date":         "2019-10-23",
			"weight_type":  "hepsepc",
			"group_type":   "SITES",
			"groups": []bson.M{
				{"name": "SITE-A", "value": 3373},
				{"name": "SITE-B", "value": 1434},
				{"name": "SITE-C", "value": 623},
				{"name": "SITE-D", "value": 7},
			},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"name":         "NonCritical",
			"date_integer": 20191004,
			"date":         "2019-10-04",
			"weight_type":  "hepsepc",
			"group_type":   "SERVICEGROUPS",
			"groups": []bson.M{
				{"name": "SVGROUP-A", "value": 334},
				{"name": "SVGROUP-B", "value": 588},
			},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"name":         "NonCritical",
			"date_integer": 20191022,
			"date":         "2019-10-22",
			"weight_type":  "hepsepc",
			"group_type":   "SERVICEGROUPS",
			"groups": []bson.M{
				{"name": "SVGROUP-A", "value": 400},
				{"name": "SVGROUP-B", "value": 188},
			},
		})

	c.InsertOne(context.TODO(),
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"name":         "NonCritical",
			"date_integer": 20191104,
			"date":         "2019-11-04",
			"weight_type":  "hepsepc",
			"group_type":   "SERVICEGROUPS",
			"groups": []bson.M{
				{"name": "SVGROUP-A", "value": 634},
				{"name": "SVGROUP-B", "value": 888},
			},
		})

}

func (suite *WeightsTestSuite) TestBadDate() {

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
		{Method: "GET", Path: "/api/v2/weights?date=2020-02", Data: ""},
		{Method: "GET", Path: "/api/v2/weights/some-uuid?date=2020-02", Data: ""},
		{Method: "POST", Path: "/api/v2/weights?date=2020-02", Data: ""},
		{Method: "PUT", Path: "/api/v2/weights/some-id?date=2020-02", Data: ""},
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

func (suite *WeightsTestSuite) TestCreateBadJson() {

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

	request, _ := http.NewRequest("POST", "/api/v2/weights", strings.NewReader(jsonInput))
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

func (suite *WeightsTestSuite) TestListQueryName() {

	request, _ := http.NewRequest("GET", "/api/v2/weights?name=NonCritical", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	profileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "date": "2019-11-04",
   "name": "NonCritical",
   "weight_type": "hepsepc",
   "group_type": "SERVICEGROUPS",
   "groups": [
    {
     "name": "SVGROUP-A",
     "value": 634
    },
    {
     "name": "SVGROUP-B",
     "value": 888
    }
   ]
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(profileJSON, output, "Response body mismatch")

}

func (suite *WeightsTestSuite) TestCreate() {

	jsonInput := `{
   "name": "weight_set3",
   "weight_type": "hepspec2",
   "group_type": "SITES",
   "groups": [
	 { "name": "site-a" , "value": 336 },
	 { "name": "site-b" , "value": 343 },
	 { "name": "site-c" , "value": 553 },
	 { "name": "site-d" , "value": 435 },
	 { "name": "site-e" , "value": 3.33 },
	 { "name": "site-f" , "value": 323.3 }
   ]
 }`

	jsonOutput := `{
 "status": {
  "message": "Weights resource succesfully created",
  "code": "201"
 },
 "data": {
  "id": "{{id}}",
  "links": {
   "self": "https:///api/v2/weights/{{id}}"
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
   "date": "2019-11-20",
   "name": "weight_set3",
   "weight_type": "hepspec2",
   "group_type": "SITES",
   "groups": [
    {
     "name": "site-a",
     "value": 336
    },
    {
     "name": "site-b",
     "value": 343
    },
    {
     "name": "site-c",
     "value": 553
    },
    {
     "name": "site-d",
     "value": 435
    },
    {
     "name": "site-e",
     "value": 3.33
    },
    {
     "name": "site-f",
     "value": 323.3
    }
   ]
  }
 ]
}`

	request, _ := http.NewRequest("POST", "/api/v2/weights?date=2019-11-20", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(201, code, "Internal Server Error")
	// Compare the expected and actual json response

	// Retrieve id from database
	var result map[string]interface{}
	c := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("weights")
	c.FindOne(context.TODO(), bson.M{"name": "weight_set3"}).Decode(&result)
	id := result["id"].(string)

	// Apply id to output template and check
	suite.Equal(strings.Replace(jsonOutput, "{{id}}", id, 2), output, "Response body mismatch")

	// Apply id to output template and check
	suite.Equal(strings.Replace(jsonOutput, "{{id}}", id, 2), output, "Response body mismatch")

	// Check that actually the item has been created
	// Call List one with the specific id
	request2, _ := http.NewRequest("GET", "/api/v2/weights/"+id, strings.NewReader(jsonInput))
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

func (suite *WeightsTestSuite) TestList() {

	request, _ := http.NewRequest("GET", "/api/v2/weights", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	weightsJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-10-23",
   "name": "Critical",
   "weight_type": "hepsepc",
   "group_type": "SITES",
   "groups": [
    {
     "name": "SITE-A",
     "value": 3373
    },
    {
     "name": "SITE-B",
     "value": 1434
    },
    {
     "name": "SITE-C",
     "value": 623
    },
    {
     "name": "SITE-D",
     "value": 7
    }
   ]
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "date": "2019-11-04",
   "name": "NonCritical",
   "weight_type": "hepsepc",
   "group_type": "SERVICEGROUPS",
   "groups": [
    {
     "name": "SVGROUP-A",
     "value": 634
    },
    {
     "name": "SVGROUP-B",
     "value": 888
    }
   ]
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(weightsJSON, output, "Response body mismatch")

}

func (suite *WeightsTestSuite) TestListPast() {

	request, _ := http.NewRequest("GET", "/api/v2/weights?date=2019-10-20", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	weightsJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-10-04",
   "name": "Critical",
   "weight_type": "hepsepc",
   "group_type": "SITES",
   "groups": [
    {
     "name": "SITE-A",
     "value": 1673
    },
    {
     "name": "SITE-B",
     "value": 1234
    },
    {
     "name": "SITE-C",
     "value": 523
    },
    {
     "name": "SITE-D",
     "value": 2
    }
   ]
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "date": "2019-10-04",
   "name": "NonCritical",
   "weight_type": "hepsepc",
   "group_type": "SERVICEGROUPS",
   "groups": [
    {
     "name": "SVGROUP-A",
     "value": 334
    },
    {
     "name": "SVGROUP-B",
     "value": 588
    }
   ]
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(weightsJSON, output, "Response body mismatch")

}

func (suite *WeightsTestSuite) TestOptionsWeights() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/weights", strings.NewReader(""))

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

func (suite *WeightsTestSuite) TestListOneNotFound() {

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

	request, _ := http.NewRequest("GET", "/api/v2/weights/wrong-id", strings.NewReader(jsonInput))
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

func (suite *WeightsTestSuite) TestListOne() {

	request, _ := http.NewRequest("GET", "/api/v2/weights/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	weightsJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-10-23",
   "name": "Critical",
   "weight_type": "hepsepc",
   "group_type": "SITES",
   "groups": [
    {
     "name": "SITE-A",
     "value": 3373
    },
    {
     "name": "SITE-B",
     "value": 1434
    },
    {
     "name": "SITE-C",
     "value": 623
    },
    {
     "name": "SITE-D",
     "value": 7
    }
   ]
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(weightsJSON, output, "Response body mismatch")

}

func (suite *WeightsTestSuite) TestUpdateNameAlreadyExists() {

	jsonInput := `{
   "name": "Critical",
   "weight_type": "hepsepc5",
   "group_type": "SITES",
   "groups": [
    {
     "name": "SITE-A",
     "value": 16733
    },
    {
     "name": "SITE-B",
     "value": 12345
    },
    {
     "name": "SITE-C",
     "value": 5233
    },
    {
     "name": "SITE-D",
     "value": 23
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
   "details": "Weights resource with the same name already exists"
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/weights/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(409, code)
	suite.Equal(jsonOutput, output)

}

func (suite *WeightsTestSuite) TestUpdateBadJson() {

	jsonInput := `{
   "name": "Critical",
   "weight_type": "hepsepc5",
   "group_type": "SITES",
   "groups": [
    {
     "name": "SITE-A",
     "value": 16733
    },
    {
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

	request, _ := http.NewRequest("PUT", "/api/v2/weights/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
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

func (suite *WeightsTestSuite) TestUpdateNotFound() {

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

	request, _ := http.NewRequest("PUT", "/api/v2/weights/wrong-id", strings.NewReader(jsonInput))
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

func (suite *WeightsTestSuite) TestUpdate() {

	jsonInput := `{
   "name": "NonCritical",
   "weight_type": "hepsepc5",
   "group_type": "SITES",
   "groups": [
    {
     "name": "SITE-A",
     "value": 16733
    },
    {
     "name": "SITE-B",
     "value": 12345
    }
   ]
}`

	jsonOutput := `{
 "status": {
  "message": "Weights resource successfully updated (new history snapshot)",
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
   "date": "2019-12-04",
   "name": "NonCritical",
   "weight_type": "hepsepc5",
   "group_type": "SITES",
   "groups": [
    {
     "name": "SITE-A",
     "value": 16733
    },
    {
     "name": "SITE-B",
     "value": 12345
    }
   ]
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/weights/6ac7d684-1f8e-4a02-a502-720e8f11e50c?date=2019-12-04", strings.NewReader(jsonInput))
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
	request2, _ := http.NewRequest("GET", "/api/v2/weights/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
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

func (suite *WeightsTestSuite) TestDeleteNotFound() {

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

	request, _ := http.NewRequest("DELETE", "/api/v2/weights/wrong-id", strings.NewReader(jsonInput))
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

func (suite *WeightsTestSuite) TestDelete() {

	request, _ := http.NewRequest("DELETE", "/api/v2/weights/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Weights resource successfully deleted",
  "code": "200"
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricProfileJSON, output, "Response body mismatch")

	// try to retrieve item
	var result map[string]interface{}
	c := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("weights")
	queryResult := c.FindOne(context.TODO(), bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"}).Decode(&result)

	suite.Equal(queryResult.Error(), mongo.ErrNoDocuments.Error(), "No not found error")
}

func (suite *WeightsTestSuite) TestCreateForbidViewer() {

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

	request, _ := http.NewRequest("POST", "/api/v2/weights", strings.NewReader(jsonInput))
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

func (suite *WeightsTestSuite) TestUpdateForbidViewer() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Access to the resource is Forbidden",
  "code": "403"
 }
}`

	request, _ := http.NewRequest("PUT", "/api/v2/weights/6ac7d684-1f8e-4a02-a502-720e8f11e007", strings.NewReader(jsonInput))
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

func (suite *WeightsTestSuite) TestDeleteForbidViewer() {

	request, _ := http.NewRequest("DELETE", "/api/v2/weights/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
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
func (suite *WeightsTestSuite) TearDownTest() {

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
func (suite *WeightsTestSuite) TearDownSuite() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())

}

func TestSuiteWeights(t *testing.T) {
	suite.Run(t, new(WeightsTestSuite))
}
