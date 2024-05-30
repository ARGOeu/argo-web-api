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

package metricProfiles

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
type MetricProfilesTestSuite struct {
	suite.Suite
	cfg              config.Config
	router           *mux.Router
	confHandler      respond.ConfHandler
	tenantDbConf     config.MongoConfig
	clientkey        string
	respUnauthorized string
}

// Setup the Test Environment
func (suite *MetricProfilesTestSuite) SetupSuite() {
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
	db = "AR_test_mprof"
	`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.respUnauthorized = "Unauthorized"
	suite.tenantDbConf = config.MongoConfig{
		Host:     "localhost",
		Port:     27017,
		Db:       "AR_test_metric_profiles_tenant",
		Password: "pass",
		Username: "dbuser",
		Store:    "ar",
	}
	suite.clientkey = "123456"

	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v2").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)
}

func (suite *MetricProfilesTestSuite) TestBadDate() {

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
		{Method: "GET", Path: "/api/v2/metric_profiles?date=2020-02", Data: ""},
		{Method: "GET", Path: "/api/v2/metric_profiles/some-uuid?date=2020-02", Data: ""},
		{Method: "POST", Path: "/api/v2/metric_profiles?date=2020-02", Data: ""},
		{Method: "PUT", Path: "/api/v2/metric_profiles/some-id?date=2020-02", Data: ""},
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

// This function runs before any test and setups the environment
func (suite *MetricProfilesTestSuite) SetupTest() {

	log.SetOutput(io.Discard)

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
			"resource": "metricProfiles.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "metricProfiles.get",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "metricProfiles.create",
			"roles":    []string{"editor"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "metricProfiles.delete",
			"roles":    []string{"editor"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "metricProfiles.update",
			"roles":    []string{"editor"},
		})

	// Seed database with metric profiles
	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("metric_profiles")

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
			"name":         "ch.cern.SAM.ROC_CRITICAL",
			"date_integer": 20191004,
			"date":         "2019-10-04",
			"description":  "critical profile",
			"services": []bson.M{
				{"service": "CREAM-CE",
					"metrics": []string{
						"emi.cream.CREAMCE-JobSubmit",
						"emi.wn.WN-Bi",
						"emi.wn.WN-Csh",
						"emi.wn.WN-SoftVer"},
				},
				{"service": "SRMv2",
					"metrics": []string{"hr.srce.SRM2-CertLifetime",
						"org.sam.SRM-Del",
						"org.sam.SRM-Get",
						"org.sam.SRM-GetSURLs",
						"org.sam.SRM-GetTURLs",
						"org.sam.SRM-Ls",
						"org.sam.SRM-LsDir",
						"org.sam.SRM-Put"},
				},
			},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			"name":         "ch.cern.SAM.ROC_CRITICAL",
			"date_integer": 20191104,
			"date":         "2019-11-04",
			"description":  "critical profile",
			"services": []bson.M{
				{"service": "CREAM-CE2",
					"metrics": []string{
						"emi.cream.CREAMCE-JobSubmit",
						"emi.wn.WN-Bi",
						"emi.wn.WN-Csh",
						"emi.wn.WN-SoftVer"},
				},
				{"service": "SRMv3",
					"metrics": []string{"hr.srce.SRM2-CertLifetime",
						"org.sam.SRM-Del",
						"org.sam.SRM-Get",
						"org.sam.SRM-GetSURLs",
						"org.sam.SRM-GetTURLs",
						"org.sam.SRM-Ls",
						"org.sam.SRM-LsDir",
						"org.sam.SRM-Put"},
				},
			},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"name":         "ch.cern.SAM.ROC",
			"date_integer": 20190504,
			"date":         "2019-05-04",
			"services": []bson.M{
				{"service": "CREAM-CE",
					"metrics": []string{
						"emi.cream.CREAMCE-JobSubmit",
						"emi.wn.WN-Bi",
						"emi.wn.WN-Csh",
						"hr.srce.CADist-Check",
						"hr.srce.CREAMCE-CertLifetime",
						"emi.wn.WN-SoftVer"},
				},
				{"service": "SRMv2",
					"metrics": []string{"hr.srce.SRM2-CertLifetime",
						"org.sam.SRM-Del",
						"org.sam.SRM-Get",
						"org.sam.SRM-GetSURLs",
						"org.sam.SRM-GetTURLs",
						"org.sam.SRM-Ls",
						"org.sam.SRM-LsDir",
						"org.sam.SRM-Put"},
				},
			},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"name":         "ch.cern.SAM.ROC",
			"date_integer": 20190604,
			"date":         "2019-06-04",
			"services": []bson.M{
				{"service": "CREAM-CE2",
					"metrics": []string{
						"emi.cream.CREAMCE-JobSubmit",
						"emi.wn.WN-Bi",
						"emi.wn.WN-Csh",
						"hr.srce.CADist-Check",
						"hr.srce.CREAMCE-CertLifetime",
						"emi.wn.WN-SoftVer"},
				},
				{"service": "SRMv3",
					"metrics": []string{"hr.srce.SRM2-CertLifetime",
						"org.sam.SRM-Del",
						"org.sam.SRM-Get",
						"org.sam.SRM-GetSURLs",
						"org.sam.SRM-GetTURLs",
						"org.sam.SRM-Ls",
						"org.sam.SRM-LsDir",
						"org.sam.SRM-Put"},
				},
			},
		})

}

func (suite *MetricProfilesTestSuite) TestList() {

	request, _ := http.NewRequest("GET", "/api/v2/metric_profiles", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-11-04",
   "name": "ch.cern.SAM.ROC_CRITICAL",
   "description": "critical profile",
   "services": [
    {
     "service": "CREAM-CE2",
     "metrics": [
      "emi.cream.CREAMCE-JobSubmit",
      "emi.wn.WN-Bi",
      "emi.wn.WN-Csh",
      "emi.wn.WN-SoftVer"
     ]
    },
    {
     "service": "SRMv3",
     "metrics": [
      "hr.srce.SRM2-CertLifetime",
      "org.sam.SRM-Del",
      "org.sam.SRM-Get",
      "org.sam.SRM-GetSURLs",
      "org.sam.SRM-GetTURLs",
      "org.sam.SRM-Ls",
      "org.sam.SRM-LsDir",
      "org.sam.SRM-Put"
     ]
    }
   ]
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "date": "2019-06-04",
   "name": "ch.cern.SAM.ROC",
   "description": "",
   "services": [
    {
     "service": "CREAM-CE2",
     "metrics": [
      "emi.cream.CREAMCE-JobSubmit",
      "emi.wn.WN-Bi",
      "emi.wn.WN-Csh",
      "hr.srce.CADist-Check",
      "hr.srce.CREAMCE-CertLifetime",
      "emi.wn.WN-SoftVer"
     ]
    },
    {
     "service": "SRMv3",
     "metrics": [
      "hr.srce.SRM2-CertLifetime",
      "org.sam.SRM-Del",
      "org.sam.SRM-Get",
      "org.sam.SRM-GetSURLs",
      "org.sam.SRM-GetTURLs",
      "org.sam.SRM-Ls",
      "org.sam.SRM-LsDir",
      "org.sam.SRM-Put"
     ]
    }
   ]
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricProfileJSON, output, "Response body mismatch")

}

func (suite *MetricProfilesTestSuite) TestListQueryName() {

	request, _ := http.NewRequest("GET", "/api/v2/metric_profiles?name=ch.cern.SAM.ROC", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "date": "2019-06-04",
   "name": "ch.cern.SAM.ROC",
   "description": "",
   "services": [
    {
     "service": "CREAM-CE2",
     "metrics": [
      "emi.cream.CREAMCE-JobSubmit",
      "emi.wn.WN-Bi",
      "emi.wn.WN-Csh",
      "hr.srce.CADist-Check",
      "hr.srce.CREAMCE-CertLifetime",
      "emi.wn.WN-SoftVer"
     ]
    },
    {
     "service": "SRMv3",
     "metrics": [
      "hr.srce.SRM2-CertLifetime",
      "org.sam.SRM-Del",
      "org.sam.SRM-Get",
      "org.sam.SRM-GetSURLs",
      "org.sam.SRM-GetTURLs",
      "org.sam.SRM-Ls",
      "org.sam.SRM-LsDir",
      "org.sam.SRM-Put"
     ]
    }
   ]
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricProfileJSON, output, "Response body mismatch")

}

func (suite *MetricProfilesTestSuite) TestListOneNotFound() {

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

	request, _ := http.NewRequest("GET", "/api/v2/metric_profiles/wrong-id", strings.NewReader(jsonInput))
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

func (suite *MetricProfilesTestSuite) TestListOne() {

	request, _ := http.NewRequest("GET", "/api/v2/metric_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-11-04",
   "name": "ch.cern.SAM.ROC_CRITICAL",
   "description": "critical profile",
   "services": [
    {
     "service": "CREAM-CE2",
     "metrics": [
      "emi.cream.CREAMCE-JobSubmit",
      "emi.wn.WN-Bi",
      "emi.wn.WN-Csh",
      "emi.wn.WN-SoftVer"
     ]
    },
    {
     "service": "SRMv3",
     "metrics": [
      "hr.srce.SRM2-CertLifetime",
      "org.sam.SRM-Del",
      "org.sam.SRM-Get",
      "org.sam.SRM-GetSURLs",
      "org.sam.SRM-GetTURLs",
      "org.sam.SRM-Ls",
      "org.sam.SRM-LsDir",
      "org.sam.SRM-Put"
     ]
    }
   ]
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricProfileJSON, output, "Response body mismatch")

}

func (suite *MetricProfilesTestSuite) TestCreateBadJson() {

	jsonInput := `{
  "name": "test_profile",
  "services": [
    {
      "service": "Service-A",
      "metrics": [
        "metric.A.1",
        "metric.A.2",
        "metric.A.3",
        "metric.A.4"
      ]
    },
    {
      "service": "Service-B",
      "metrics": [
        "metric.B.1",
        "metric.B.2"
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

	request, _ := http.NewRequest("POST", "/api/v2/metric_profiles", strings.NewReader(jsonInput))
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

func (suite *MetricProfilesTestSuite) TestCreate() {

	jsonInput := `{
  "name": "test_profile",
  "services": [
    {
      "service": "Service-A",
      "metrics": [
        "metric.A.1",
        "metric.A.2",
        "metric.A.3",
        "metric.A.4"
      ]
    },
    {
      "service": "Service-B",
      "metrics": [
        "metric.B.1",
        "metric.B.2"
      ]
    }
  ]
}`

	jsonOutput := `{
 "status": {
  "message": "Metric Profile successfully created",
  "code": "201"
 },
 "data": {
  "id": "{{id}}",
  "links": {
   "self": "https:///api/v2/metric_profiles/{{id}}"
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
   "date": "2019-11-08",
   "name": "test_profile",
   "description": "",
   "services": [
    {
     "service": "Service-A",
     "metrics": [
      "metric.A.1",
      "metric.A.2",
      "metric.A.3",
      "metric.A.4"
     ]
    },
    {
     "service": "Service-B",
     "metrics": [
      "metric.B.1",
      "metric.B.2"
     ]
    }
   ]
  }
 ]
}`

	request, _ := http.NewRequest("POST", "/api/v2/metric_profiles?date=2019-11-08", strings.NewReader(jsonInput))
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
	var result MetricProfile
	c := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("metric_profiles")
	c.FindOne(context.TODO(), bson.M{"name": "test_profile"}).Decode(&result)
	id := result.ID

	// Apply id to output template and check
	suite.Equal(strings.Replace(jsonOutput, "{{id}}", id, 2), output, "Response body mismatch")

	// Check that actually the item has been created
	// Call List one with the specific id
	request2, _ := http.NewRequest("GET", "/api/v2/metric_profiles/"+id, strings.NewReader(jsonInput))
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

func (suite *MetricProfilesTestSuite) TestCreateNameAlreadyExists() {

	jsonInput := `{
  "name": "ch.cern.SAM.ROC_CRITICAL",
  "services": [
    {
      "service": "Service-A",
      "metrics": [
        "metric.A.1",
        "metric.A.2",
        "metric.A.3",
        "metric.A.4"
      ]
    },
    {
      "service": "Service-B",
      "metrics": [
        "metric.B.1",
        "metric.B.2"
      ]
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
   "details": "Metric profile with the same name already exists"
  }
 ]
}`

	request, _ := http.NewRequest("POST", "/api/v2/metric_profiles", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(409, code)
	suite.Equal(jsonOutput, output)

}

func (suite *MetricProfilesTestSuite) TestUpdateNameAlreadyExists() {

	jsonInput := `{
  "name": "ch.cern.SAM.ROC_CRITICAL",
  "services": [
    {
      "service": "Service-A",
      "metrics": [
        "metric.A.1",
        "metric.A.2",
        "metric.A.3",
        "metric.A.4"
      ]
    },
    {
      "service": "Service-B",
      "metrics": [
        "metric.B.1",
        "metric.B.2"
      ]
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
   "details": "Metric profile with the same name already exists"
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/metric_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(409, code)
	suite.Equal(jsonOutput, output)

}

func (suite *MetricProfilesTestSuite) TestUpdateBadJson() {

	jsonInput := `{
  "name": "test_profile",
  "services": [
    {
      "service": "Service-A",
      "metrics": [
        "metric.A.1",
        "metric.A.2",
        "metric.A.3",
        "metric.A.4"
      ]
    },
    {
      "service": "Service-B",
      "metrics": [
        "metric.B.1",
        "metric.B.2"
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

	request, _ := http.NewRequest("PUT", "/api/v2/metric_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
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

func (suite *MetricProfilesTestSuite) TestUpdateNotFound() {

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

	request, _ := http.NewRequest("PUT", "/api/v2/metric_profiles/wrong-id", strings.NewReader(jsonInput))
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

func (suite *MetricProfilesTestSuite) TestUpdate() {

	jsonInput := `{
  "name": "test_profile",
  "description": "just for testing",
  "services": [
    {
      "service": "Service-AX",
      "metrics": [
        "metric.A.1",
        "metric.A.2",
        "metric.A.3",
        "metric.A.4"
      ]
    },
    {
      "service": "Service-BX",
      "metrics": [
        "metric.B.1",
        "metric.B.2"
      ]
    }
  ]
}`

	jsonOutput := `{
 "status": {
  "message": "Metric Profile successfully updated (new history snapshot)",
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
   "date": "2019-07-08",
   "name": "test_profile",
   "description": "just for testing",
   "services": [
    {
     "service": "Service-AX",
     "metrics": [
      "metric.A.1",
      "metric.A.2",
      "metric.A.3",
      "metric.A.4"
     ]
    },
    {
     "service": "Service-BX",
     "metrics": [
      "metric.B.1",
      "metric.B.2"
     ]
    }
   ]
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/metric_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c?date=2019-07-08", strings.NewReader(jsonInput))
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
	request2, _ := http.NewRequest("GET", "/api/v2/metric_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50c", strings.NewReader(jsonInput))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)

	code2 := response2.Code
	output2 := response2.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(200, code2, "Internal Server Error")
	// Compare the expected and actual json response

	res := []MetricProfile{}

	col := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("metric_profiles")
	cursor, _ := col.Find(context.TODO(), bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c"})

	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &res)

	suite.Equal(jsonUpdated, output2, "Response body mismatch")

}

func (suite *MetricProfilesTestSuite) TestListEmpty() {

	c := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("metric_profiles")
	c.Drop(context.TODO())

	request, _ := http.NewRequest("GET", "/api/v2/metric_profiles", strings.NewReader(""))
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

func (suite *MetricProfilesTestSuite) TestDeleteNotFound() {

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

	request, _ := http.NewRequest("DELETE", "/api/v2/metric_profiles/wrong-id", strings.NewReader(jsonInput))
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

func (suite *MetricProfilesTestSuite) TestDelete() {

	request, _ := http.NewRequest("DELETE", "/api/v2/metric_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricProfileJSON := `{
 "status": {
  "message": "Metric Profile Successfully Deleted",
  "code": "200"
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricProfileJSON, output, "Response body mismatch")

	// check that the element has actually been Deleted
	// connect to mongodb

	// try to retrieve the item
	c := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("metric_profiles")
	queryResult := c.FindOne(context.TODO(), bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"})

	suite.NotEqual(queryResult.Err(), nil, "No not found error")
	suite.Equal(queryResult.Err(), mongo.ErrNoDocuments, "No not found error")
}

func (suite *MetricProfilesTestSuite) TestOptionsMetricProfiles() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/metric_profiles", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.Result().Header

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, POST, DELETE, PUT, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

	request, _ = http.NewRequest("OPTIONS", "/api/v2/metric_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))

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

func (suite *MetricProfilesTestSuite) TestCreateForbidViewer() {

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

	request, _ := http.NewRequest("POST", "/api/v2/metric_profiles", strings.NewReader(jsonInput))
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

func (suite *MetricProfilesTestSuite) TestUpdateForbidViewer() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Access to the resource is Forbidden",
  "code": "403"
 }
}`

	request, _ := http.NewRequest("PUT", "/api/v2/metric_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e007", strings.NewReader(jsonInput))
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

func (suite *MetricProfilesTestSuite) TestDeleteForbidViewer() {

	request, _ := http.NewRequest("DELETE", "/api/v2/metric_profiles/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
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

// TearDownTest things to do after each test
func (suite *MetricProfilesTestSuite) TearDownTest() {

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

// TearDownSuite things to do after suite ends
func (suite *MetricProfilesTestSuite) TearDownSuite() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
}

func TestSuiteMetricProfiles(t *testing.T) {
	suite.Run(t, new(MetricProfilesTestSuite))
}
