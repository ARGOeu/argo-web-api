/*
 * Copyright (c) 2018 GRNET S.A.
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

package thresholdsProfiles

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
type ThresholdsProfilesTestSuite struct {
	suite.Suite
	cfg              config.Config
	router           *mux.Router
	confHandler      respond.ConfHandler
	tenantDbConf     config.MongoConfig
	clientkey        string
	respUnauthorized string
}

// Setup the Test Environment
func (suite *ThresholdsProfilesTestSuite) SetupSuite() {

	log.SetOutput(io.Discard)

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
	    db = "AR_test_op_profiles"
	    `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.respUnauthorized = "Unauthorized"
	suite.tenantDbConf = config.MongoConfig{
		Host:     "localhost",
		Port:     27017,
		Db:       "AR_test_thresholds_profiles_tenant",
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
func (suite *ThresholdsProfilesTestSuite) SetupTest() {

	// Seed database with tenants
	//TODO: move tests to
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
			"resource": "thresholdsProfiles.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "thresholdsProfiles.get",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "thresholdsProfiles.create",
			"roles":    []string{"editor"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "thresholdsProfiles.delete",
			"roles":    []string{"editor"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "thresholdsProfiles.update",
			"roles":    []string{"editor"},
		})

	// Seed database with thresholds profiles
	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("thresholds_profiles")

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "date_integer", Value: -1},
			{Key: "id", Value: 1},
		},
		Options: options.Index().SetUnique(false), // Set this according to your requirements
	}

	c.Indexes().CreateOne(context.TODO(), indexModel)
	c.InsertOne(context.TODO(),
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			"date_integer": 20190105,
			"date":         "2019-01-05",
			"name":         "thr01",
			"rules": []bson.M{{
				"host":       "hostFoo",
				"metric":     "metricA",
				"thresholds": "entries=1;3;2:0;10",
			}},
		})

	c.InsertOne(context.TODO(),
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			"date_integer": 20191004,
			"date":         "2019-10-04",
			"name":         "thr01",
			"rules": []bson.M{{
				"host":       "hostFoo",
				"metric":     "metricA",
				"thresholds": "freshnesss=1s;10;9:;0;25 entries=1;3;2:0;10",
			}},
		})

	c.InsertOne(context.TODO(),
		bson.M{
			"id":           "6ac7d222-1f8e-4a02-a502-720e8f11e50b",
			"date_integer": 20191004,
			"date":         "2019-10-04",
			"name":         "thr02",
			"rules": []bson.M{{
				"host":       "hostFoo",
				"metric":     "metricA",
				"thresholds": "freshness=1s;10;9:;0;25 entries=1;3;2:0;10",
			}},
		})

	c.InsertOne(context.TODO(),
		bson.M{
			"id":           "6ac7d555-1f8e-4a02-a502-720e8f11e50b",
			"date_integer": 20191004,
			"date":         "2019-10-04",
			"name":         "thr03",
			"rules": []bson.M{{
				"host":       "hostFoo",
				"metric":     "metricA",
				"thresholds": "freshness=1s;10;9:;0;25 entries=1;3;2:0;10",
			}},
		})
}

func (suite *ThresholdsProfilesTestSuite) TestBadDate() {

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
		{Method: "GET", Path: "/api/v2/thresholds_profiles?date=2020-02", Data: ""},
		{Method: "GET", Path: "/api/v2/thresholds_profiles/some-uuid?date=2020-02", Data: ""},
		{Method: "POST", Path: "/api/v2/thresholds_profiles?date=2020-02", Data: ""},
		{Method: "PUT", Path: "/api/v2/thresholds_profiles/some-id?date=2020-02", Data: ""},
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
func (suite *ThresholdsProfilesTestSuite) TestList() {

	request, _ := http.NewRequest("GET", "/api/v2/thresholds_profiles", strings.NewReader(""))
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
   "id": "6ac7d222-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-10-04",
   "name": "thr02",
   "rules": [
    {
     "host": "hostFoo",
     "metric": "metricA",
     "thresholds": "freshness=1s;10;9:;0;25 entries=1;3;2:0;10"
    }
   ]
  },
  {
   "id": "6ac7d555-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-10-04",
   "name": "thr03",
   "rules": [
    {
     "host": "hostFoo",
     "metric": "metricA",
     "thresholds": "freshness=1s;10;9:;0;25 entries=1;3;2:0;10"
    }
   ]
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-10-04",
   "name": "thr01",
   "rules": [
    {
     "host": "hostFoo",
     "metric": "metricA",
     "thresholds": "freshnesss=1s;10;9:;0;25 entries=1;3;2:0;10"
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

func (suite *ThresholdsProfilesTestSuite) TestListQueryName() {

	request, _ := http.NewRequest("GET", "/api/v2/thresholds_profiles?name=thr02", strings.NewReader(""))
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
   "id": "6ac7d222-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-10-04",
   "name": "thr02",
   "rules": [
    {
     "host": "hostFoo",
     "metric": "metricA",
     "thresholds": "freshness=1s;10;9:;0;25 entries=1;3;2:0;10"
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

func (suite *ThresholdsProfilesTestSuite) TestListOneNotFound() {

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

	request, _ := http.NewRequest("GET", "/api/v2/thresholds_profiles/wrong-id", strings.NewReader(jsonInput))
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

func (suite *ThresholdsProfilesTestSuite) TestListOne() {

	request, _ := http.NewRequest("GET", "/api/v2/thresholds_profiles/6ac7d222-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
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
   "id": "6ac7d222-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-10-04",
   "name": "thr02",
   "rules": [
    {
     "host": "hostFoo",
     "metric": "metricA",
     "thresholds": "freshness=1s;10;9:;0;25 entries=1;3;2:0;10"
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

func (suite *ThresholdsProfilesTestSuite) TestCreateBadJson() {

	jsonInput := `{
  "name": "test_profile",
  "namespace [
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

	request, _ := http.NewRequest("POST", "/api/v2/thresholds_profiles", strings.NewReader(jsonInput))
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

func (suite *ThresholdsProfilesTestSuite) TestCreate() {

	jsonInput := `{
  "name" : "thr04",
  "rules": [
    {
      "metric": "metricB",
      "thresholds": "time=1s;10;9:30;0;30"
    }
  ]
}`

	jsonOutput := `{
 "status": {
  "message": "Thresholds Profile successfully created",
  "code": "201"
 },
 "data": {
  "id": "{{ID}}",
  "links": {
   "self": "https:///api/v2/thresholds_profiles/{{ID}}"
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
   "id": "{{ID}}",
   "date": "2019-11-11",
   "name": "thr04",
   "rules": [
    {
     "metric": "metricB",
     "thresholds": "time=1s;10;9:30;0;30"
    }
   ]
  }
 ]
}`

	request, _ := http.NewRequest("POST", "/api/v2/thresholds_profiles?date=2019-11-11", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	// Check that we must have a 201 ok code
	suite.Equal(201, code, "Internal Server Error")
	// Compare the expected and actual json response

	// Grab ID from mongodb

	// Retrieve id from database
	var result ThresholdsProfile
	c := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("thresholds_profiles")
	c.FindOne(context.TODO(), bson.M{"name": "thr04"}).Decode(&result)

	id := result.ID

	// Apply id to output template and check
	suite.Equal(strings.Replace(jsonOutput, "{{ID}}", id, 2), output, "Response body mismatch")

	// Check that actually the item has been created
	// Call List one with the specific ID
	request2, _ := http.NewRequest("GET", "/api/v2/thresholds_profiles/"+id, strings.NewReader(jsonInput))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)

	code2 := response2.Code
	output2 := response2.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(200, code2, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(strings.Replace(jsonCreated, "{{ID}}", id, 1), output2, "Response body mismatch")
}

func (suite *ThresholdsProfilesTestSuite) TestUpdateBadJson() {

	jsonInput := `{
   "name": "yolo",
   "namespace": "testin
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

	request, _ := http.NewRequest("PUT", "/api/v2/thresholds_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
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

func (suite *ThresholdsProfilesTestSuite) TestUpdateNotFound() {

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

	request, _ := http.NewRequest("PUT", "/api/v2/thresholds_profiles/wrong-id", strings.NewReader(jsonInput))
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

func (suite *ThresholdsProfilesTestSuite) TestUpdateSameName() {

	jsonInput := `{
  "name" : "thr02",
  "rules": [
    {
      "metric": "metricB",
      "thresholds": "time=2s;10;9:30;0;30"
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
   "details": "Thresholds profile with the same name already exists"
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/thresholds_profiles/6ac7d555-1f8e-4a02-a502-720e8f11e50b?date=2019-11-12", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(409, code, "Internal Server Error")
	// Compare the expected and actual json response

	// Apply id to output template and check
	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *ThresholdsProfilesTestSuite) TestUpdate() {

	jsonInput := `{
  "name" : "thr04",
  "rules": [
    {
      "metric": "metricB",
      "thresholds": "time=2s;10;9:30;0;30"
    }
  ]
}`

	jsonOutput := `{
 "status": {
  "message": "Thresholds Profile successfully updated",
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
   "id": "6ac7d555-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-11-12",
   "name": "thr04",
   "rules": [
    {
     "metric": "metricB",
     "thresholds": "time=2s;10;9:30;0;30"
    }
   ]
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/thresholds_profiles/6ac7d555-1f8e-4a02-a502-720e8f11e50b?date=2019-11-12", strings.NewReader(jsonInput))
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
	request2, _ := http.NewRequest("GET", "/api/v2/thresholds_profiles/6ac7d555-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(jsonInput))
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

func (suite *ThresholdsProfilesTestSuite) TestDeleteNotFound() {

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

	request, _ := http.NewRequest("DELETE", "/api/v2/thresholds_profiles/wrong-id", strings.NewReader(jsonInput))
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

func (suite *ThresholdsProfilesTestSuite) TestDelete() {

	request, _ := http.NewRequest("DELETE", "/api/v2/thresholds_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Thresholds Profile Successfully Deleted",
  "code": "200"
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricProfileJSON, output, "Response body mismatch")

	// try to retrieve item
	var result map[string]interface{}
	c := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("thresholds_profiles")
	err := c.FindOne(context.TODO(), bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"}).Decode(&result)

	suite.NotEqual(err, nil, "No not found error")
	suite.Equal(err.Error(), mongo.ErrNoDocuments.Error(), "No not found error")
}

func (suite *ThresholdsProfilesTestSuite) TestListEmpty() {

	c := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("thresholds_profiles")
	c.Drop(context.TODO())

	request, _ := http.NewRequest("GET", "/api/v2/thresholds_profiles", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	emptyList := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": []
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(emptyList, output, "Response body mismatch")

}

func (suite *ThresholdsProfilesTestSuite) TestOptionsOperationsProfiles() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/thresholds_profiles", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, POST, DELETE, PUT, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v2/thresholds_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))

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

func (suite *ThresholdsProfilesTestSuite) TestCreateForbidViewer() {

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

	request, _ := http.NewRequest("POST", "/api/v2/thresholds_profiles", strings.NewReader(jsonInput))
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

func (suite *ThresholdsProfilesTestSuite) TestUpdateForbidViewer() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Access to the resource is Forbidden",
  "code": "403"
 }
}`

	request, _ := http.NewRequest("PUT", "/api/v2/thresholds_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e007", strings.NewReader(jsonInput))
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

func (suite *ThresholdsProfilesTestSuite) TestDeleteForbidViewer() {

	request, _ := http.NewRequest("DELETE", "/api/v2/thresholds_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
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

func (suite *ThresholdsProfilesTestSuite) TestInvalidCreate() {

	jsonInput := `{
			"name":"test-invalid-01",
			"rules":[
				{"thresholds":"bad01=33;33s"},
				{"thresholds":"good01=33s;33 good02=1s;~:10;9:;-20;30"},
				{"thresholds":"bad02=33sbad03=1s;~~:10;9:;-20;30"},
				{"thresholds":"33;33 bad04=33s;33 -20;30"},
				{"thresholds":"good01=2KB;0:3;2:10;0;20 good02=1c;~:10;9:30;-30;30"}
			]
		}`

	jsonOutput := `{
 "status": {
  "message": "Validation Error",
  "code": "422"
 },
 "errors": [
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Invalid threshold: bad01=33;33s"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Invalid threshold: bad02=33sbad03=1s;~~:10;9:;-20;30"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Invalid threshold: 33;33 bad04=33s;33 -20;30"
  }
 ]
}`

	request, _ := http.NewRequest("POST", "/api/v2/thresholds_profiles", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(422, code, "Internal Server Error")

	// Apply id to output template and check
	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *ThresholdsProfilesTestSuite) TestInvalidUpdate() {

	jsonInput := `{
			"name":"test-invalid-01",
			"rules":[
				{"thresholds":"bad01=33;33s"},
				{"thresholds":"good01=33s;33 good02=1s;~:10;9:;-20;30"},
				{"thresholds":"bad02=33sbad03=1s;~~:10;9:;-20;30"},
				{"thresholds":"33;33 bad04=33s;33 -20;30"},
				{"thresholds":"good01=2KB;0:3;2:10;0;20 good02=1c;~:10;9:30;-30;30"}
			]
		}`

	jsonOutput := `{
 "status": {
  "message": "Validation Error",
  "code": "422"
 },
 "errors": [
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Invalid threshold: bad01=33;33s"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Invalid threshold: bad02=33sbad03=1s;~~:10;9:;-20;30"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Invalid threshold: 33;33 bad04=33s;33 -20;30"
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/thresholds_profiles/6ac7d555-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(422, code, "Internal Server Error")
	// Compare the expected and actual json response

	// Apply id to output template and check
	suite.Equal(jsonOutput, output, "Response body mismatch")

}

// TearDownTest to tear down every test
func (suite *ThresholdsProfilesTestSuite) TearDownTest() {

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
func (suite *ThresholdsProfilesTestSuite) TearDownSuite() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
}

func TestSuiteThresholdsProfiles(t *testing.T) {
	suite.Run(t, new(ThresholdsProfilesTestSuite))
}
