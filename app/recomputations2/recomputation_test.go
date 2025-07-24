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

package recomputations2

import (
	"bytes"
	"context"
	"encoding/json"
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
type RecomputationsProfileTestSuite struct {
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
func (suite *RecomputationsProfileTestSuite) SetupSuite() {

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
	    db = "AR_test_recomputations2"
	    `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.respRecomputationsCreated = " <root>\n" +
		"   <Message>A recalculation request has been filed</Message>\n </root>"

	suite.respUnauthorized = "Unauthorized"
	suite.tenantDbConf = config.MongoConfig{
		Host:     "localhost",
		Port:     27017,
		Db:       "AR_test_recomputations2_tenant",
		Password: "h4shp4ss",
		Username: "dbuser",
		Store:    "ar",
	}
	suite.clientkey = "mysecretcombination"

	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v2").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)

}

// This function runs before any test and setups the environment
func (suite *RecomputationsProfileTestSuite) SetupTest() {

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
					"name":    "John Snow",
					"email":   "J.Snow@brothers.wall",
					"api_key": "wh1t3_w@lk3rs",
					"roles":   []string{"editor"},
				},
				{
					"name":    "King Joffrey",
					"email":   "g0dk1ng@kingslanding.gov",
					"api_key": "sansa <3",
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
					"name":    "Joe Complex",
					"email":   "C.Joe@egi.eu",
					"api_key": suite.clientkey,
					"roles":   []string{"editor"},
				},
				{
					"name":    "Josh Plain",
					"email":   "P.Josh@egi.eu",
					"api_key": "VIEWERKEY",
					"roles":   []string{"viewer"},
				},
			}})

	c = suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("roles")
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "recomputations.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "recomputations.get",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "recomputations.submit",
			"roles":    []string{"editor"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "recomputations.delete",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "recomputations.update",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "recomputations.changeStatus",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "recomputations.resetStatus",
			"roles":    []string{"editor", "viewer"},
		})
	// Seed database with recomputations
	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("recomputations")
	c.InsertOne(context.TODO(),
		MongoInterface{
			ID:             "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			RequesterName:  "John Snow",
			RequesterEmail: "jsnow@wall.com",
			StartTime:      "2015-03-10T12:00:00Z",
			EndTime:        "2015-03-30T23:00:00Z",
			Reason:         "reasons",
			Report:         "EGI_Critical",
			Exclude:        []string{"SITE1", "SITE3"},
			Status:         "pending",
			Timestamp:      "2015-04-01T14:58:40Z",
			History:        []HistoryItem{{Status: "pending", Timestamp: "2015-04-01T14:58:40Z"}},
		},
	)
	c.InsertOne(context.TODO(),
		MongoInterface{
			ID:             "6ac7d684-1f8e-4a02-a502-720e8f11e50a",
			RequesterName:  "Arya Stark",
			RequesterEmail: "astark@shadowguild.com",
			StartTime:      "2015-01-10T12:00:00Z",
			EndTime:        "2015-01-30T23:00:00Z",
			Reason:         "power cuts",
			Report:         "EGI_Critical",
			Exclude:        []string{"SITE2", "SITE4"},
			Status:         "running",
			Timestamp:      "2015-02-01T14:58:40Z",
			History: []HistoryItem{{Status: "pending", Timestamp: "2015-02-01T14:58:40Z"},
				{Status: "running", Timestamp: "2015-02-01T16:58:40Z"}},
		},
	)

	c.InsertOne(context.TODO(),
		MongoInterface{
			ID:             "6ac7d684-1f8e-4a02-a502-720e8f11e777",
			RequesterName:  "John Doe",
			RequesterEmail: "johndoe@example.com",
			StartTime:      "2022-01-10T12:00:00Z",
			EndTime:        "2022-01-10T23:00:00Z",
			Reason:         "issue with metric checks",
			Report:         "EGI_Critical",
			Exclude:        []string{},
			Status:         "pending",
			Timestamp:      "2022-01-11T23:58:40Z",
			ExcludeMetrics: []ExcludedMetric{{Metric: "check-1"}, {Hostname: "host1.example.com", Metric: "check-2"}, {Group: "Affected-Site", Metric: "check-3"}},
			History:        []HistoryItem{{Status: "pending", Timestamp: "2015-02-01T14:58:40Z"}},
		})

}

func (suite *RecomputationsProfileTestSuite) TestListOneRecomputations() {
	request, _ := http.NewRequest("GET", "/api/v2/recomputations/6ac7d684-1f8e-4a02-a502-720e8f11e50a", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	recomputationRequestsJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": {
  "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50a",
  "requester_name": "Arya Stark",
  "requester_email": "astark@shadowguild.com",
  "reason": "power cuts",
  "start_time": "2015-01-10T12:00:00Z",
  "end_time": "2015-01-30T23:00:00Z",
  "report": "EGI_Critical",
  "exclude": [
   "SITE2",
   "SITE4"
  ],
  "status": "running",
  "timestamp": "2015-02-01T14:58:40Z",
  "history": [
   {
    "status": "pending",
    "timestamp": "2015-02-01T14:58:40Z"
   },
   {
    "status": "running",
    "timestamp": "2015-02-01T16:58:40Z"
   }
  ]
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(recomputationRequestsJSON, output, "Response body mismatch")

}

func (suite *RecomputationsProfileTestSuite) TestListOneRecomputationNotFound() {
	request, _ := http.NewRequest("GET", "/api/v2/recomputations/wrong_id", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	recomputationRequestsJSON := `{
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
	suite.Equal(404, code)
	// Compare the expected and actual json response
	suite.Equal(recomputationRequestsJSON, output)
}

func (suite *RecomputationsProfileTestSuite) TestListErrorRecomputations() {
	suite.router.Methods("GET").Handler(suite.confHandler.Respond(List))
	request, _ := http.NewRequest("GET", "/api/v2/recomputations", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/jaason")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	recomputationRequestsJSON := `{
 "status": {
  "message": "Not Acceptable Content Type",
  "code": "406",
  "details": "Accept header provided did not contain any valid content types. Acceptable content types are 'application/xml' and 'application/json'"
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(406, code, "Should be not acceptable")
	// Compare the expected and actual xml response
	suite.Equal(recomputationRequestsJSON, output, "Response body mismatch")
}

func (suite *RecomputationsProfileTestSuite) TestListRecomputations() {
	request, _ := http.NewRequest("GET", "/api/v2/recomputations", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	recomputationRequestsJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50a",
   "requester_name": "Arya Stark",
   "requester_email": "astark@shadowguild.com",
   "reason": "power cuts",
   "start_time": "2015-01-10T12:00:00Z",
   "end_time": "2015-01-30T23:00:00Z",
   "report": "EGI_Critical",
   "exclude": [
    "SITE2",
    "SITE4"
   ],
   "status": "running",
   "timestamp": "2015-02-01T14:58:40Z",
   "history": [
    {
     "status": "pending",
     "timestamp": "2015-02-01T14:58:40Z"
    },
    {
     "status": "running",
     "timestamp": "2015-02-01T16:58:40Z"
    }
   ]
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "requester_name": "John Snow",
   "requester_email": "jsnow@wall.com",
   "reason": "reasons",
   "start_time": "2015-03-10T12:00:00Z",
   "end_time": "2015-03-30T23:00:00Z",
   "report": "EGI_Critical",
   "exclude": [
    "SITE1",
    "SITE3"
   ],
   "status": "pending",
   "timestamp": "2015-04-01T14:58:40Z",
   "history": [
    {
     "status": "pending",
     "timestamp": "2015-04-01T14:58:40Z"
    }
   ]
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e777",
   "requester_name": "John Doe",
   "requester_email": "johndoe@example.com",
   "reason": "issue with metric checks",
   "start_time": "2022-01-10T12:00:00Z",
   "end_time": "2022-01-10T23:00:00Z",
   "report": "EGI_Critical",
   "exclude": [],
   "status": "pending",
   "timestamp": "2022-01-11T23:58:40Z",
   "history": [
    {
     "status": "pending",
     "timestamp": "2015-02-01T14:58:40Z"
    }
   ],
   "exclude_metrics": [
    {
     "metric": "check-1"
    },
    {
     "metric": "check-2",
     "timestamp": "host1.example.com"
    },
    {
     "metric": "check-3",
     "group": "Affected-Site"
    }
   ]
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(recomputationRequestsJSON, output, "Response body mismatch")

}

func (suite *RecomputationsProfileTestSuite) TestListRecomputationsWithPeriodA() {

	urls := []string{
		"/api/v2/recomputations?date=2015-01-10",
		"/api/v2/recomputations?date=2015-01-11",
		"/api/v2/recomputations?date=2015-01-14",
		"/api/v2/recomputations?date=2015-01-16",
		"/api/v2/recomputations?date=2015-01-20",
		"/api/v2/recomputations?date=2015-01-29",
		"/api/v2/recomputations?date=2015-01-30",
	}

	recomputationRequestsJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50a",
   "requester_name": "Arya Stark",
   "requester_email": "astark@shadowguild.com",
   "reason": "power cuts",
   "start_time": "2015-01-10T12:00:00Z",
   "end_time": "2015-01-30T23:00:00Z",
   "report": "EGI_Critical",
   "exclude": [
    "SITE2",
    "SITE4"
   ],
   "status": "running",
   "timestamp": "2015-02-01T14:58:40Z",
   "history": [
    {
     "status": "pending",
     "timestamp": "2015-02-01T14:58:40Z"
    },
    {
     "status": "running",
     "timestamp": "2015-02-01T16:58:40Z"
    }
   ]
  }
 ]
}`

	for _, url := range urls {
		request, _ := http.NewRequest("GET", url, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		code := response.Code
		output := response.Body.String()

		// Check that we must have a 200 ok code
		suite.Equal(200, code, "Internal Server Error")
		// Compare the expected and actual xml response
		suite.Equal(recomputationRequestsJSON, output, "Response body mismatch")

	}

}

func (suite *RecomputationsProfileTestSuite) TestListRecomputationsWithPeriodB() {

	urls := []string{
		"/api/v2/recomputations?date=2015-03-10",
		"/api/v2/recomputations?date=2015-03-11",
		"/api/v2/recomputations?date=2015-03-14",
		"/api/v2/recomputations?date=2015-03-16",
		"/api/v2/recomputations?date=2015-03-20",
		"/api/v2/recomputations?date=2015-03-29",
		"/api/v2/recomputations?date=2015-03-30",
	}

	recomputationRequestsJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "requester_name": "John Snow",
   "requester_email": "jsnow@wall.com",
   "reason": "reasons",
   "start_time": "2015-03-10T12:00:00Z",
   "end_time": "2015-03-30T23:00:00Z",
   "report": "EGI_Critical",
   "exclude": [
    "SITE1",
    "SITE3"
   ],
   "status": "pending",
   "timestamp": "2015-04-01T14:58:40Z",
   "history": [
    {
     "status": "pending",
     "timestamp": "2015-04-01T14:58:40Z"
    }
   ]
  }
 ]
}`

	for _, url := range urls {
		request, _ := http.NewRequest("GET", url, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		code := response.Code
		output := response.Body.String()

		// Check that we must have a 200 ok code
		suite.Equal(200, code, "Internal Server Error")
		// Compare the expected and actual xml response
		suite.Equal(recomputationRequestsJSON, output, "Response body mismatch")

	}

}

func (suite *RecomputationsProfileTestSuite) TestListRecomputationsWithPeriodC() {
	urls := []string{
		"/api/v2/recomputations?date=2015-01-01",
		"/api/v2/recomputations?date=2015-01-09",
		"/api/v2/recomputations?date=2015-02-14",
		"/api/v2/recomputations?date=2015-02-16",
		"/api/v2/recomputations?date=2015-04-20",
		"/api/v2/recomputations?date=2015-04-29",
		"/api/v2/recomputations?date=2015-05-30",
		"/api/v2/recomputations?date=2015-01-30&report=Foo",
		"/api/v2/recomputations?date=2015-03-15&report=Bar",
	}

	recomputationRequestsJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": []
}`

	for _, url := range urls {
		request, _ := http.NewRequest("GET", url, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		code := response.Code
		output := response.Body.String()

		// Check that we must have a 200 ok code
		suite.Equal(200, code, "Internal Server Error")
		// Compare the expected and actual xml response
		suite.Equal(recomputationRequestsJSON, output, "Response body mismatch")

	}

}

func (suite *RecomputationsProfileTestSuite) TestListRecomputationsWithPeriodError() {
	request, _ := http.NewRequest("GET", "/api/v2/recomputations?date=2020-03-303", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	recomputationRequestsJSON := `{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "Bad Request",
   "code": "400",
   "details": "date argument should be in the YYYY-MM-DD format"
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(400, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(recomputationRequestsJSON, output, "Response body mismatch")

}

func (suite *RecomputationsProfileTestSuite) TestSubmitRecomputations() {
	submission := IncomingRecomputation{
		StartTime:      "2015-01-10T12:00:00Z",
		EndTime:        "2015-01-30T23:00:00Z",
		RequesterName:  "Joe Complexz",
		RequesterEmail: "C.Joecz@egi.eu",
		Reason:         "Ups failure",
		Report:         "EGI_Critical",
		Exclude:        []string{"SITE5", "SITE8"},
	}
	jsonsubmission, _ := json.Marshal(submission)

	request, _ := http.NewRequest("POST", "https://argo-web-api.grnet.gr:443/api/v2/recomputations", bytes.NewBuffer(jsonsubmission))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	recomputationRequestsJSON := `{
 "status": {
  "message": "Recomputations successfully created",
  "code": "201"
 },
 "data": {
  "id": ".+",
  "links": {
   "self": "https://argo-web-api.grnet.gr:443/api/v2/recomputations/.+"
  }
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(202, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Regexp(recomputationRequestsJSON, output, "Response body mismatch")

	var results []MongoInterface
	rCol := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection(recomputationsColl)
	findOptions := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})
	cursor, err := rCol.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	suite.Equal(len(results), 4)
	suite.Equal("2015-01-10T12:00:00Z", results[3].StartTime)
	suite.Equal("2015-01-30T23:00:00Z", results[3].EndTime)
	suite.Equal("Joe Complexz", results[3].RequesterName)

}

func (suite *RecomputationsProfileTestSuite) TestSubmitRecomΑppliedStatuses() {
	jsonRecomp := `
	{
  "id": "56db4f1a-f331-46ca-b0fd-4555b4aa1cfc",
  "requester_name": "requester",
  "requester_email": "request@request.gr",
  "reason": "testing_compute_engine",
  "start_time": "2022-01-12T00:00:00Z",
  "end_time": "2022-01-15T00:00:00Z",
  "exclude": [
    "SITE-A",
    "SITE-B"
  ],
  "applied_status_changes": [
    {
      "metric": "Metric_X",
      "service": "Service_X",
      "state": "CRITICAL"
    },
    {
      "metric": "Metric_XX",
      "service": "Service_X",
      "state": "OK"
    },
    {
      "metric": "Metric_XXX",
      "hostname": "Hostname_XXX",
      "state": "EXCLUDED"
    },
    {
      "metric": "Metric_XXXX",
      "group": "Group_XXXX",
      "state": "WARNING"
    }
  ],
  "status": "running",
  "timestamp": "2022-01-14 14:58:40",
  "exclude_monitoring_source": [
    {
      "host": "monA",
      "start_time": "2022-01-13T00:00:00Z",
      "end_time": "2022-01-13T23:59:59Z"
    }
  ]
}`

	request, _ := http.NewRequest("POST", "https://argo-web-api.grnet.gr:443/api/v2/recomputations", bytes.NewBuffer([]byte(jsonRecomp)))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	recomputationRequestsJSON := `{
 "status": {
  "message": "Recomputations successfully created",
  "code": "201"
 },
 "data": {
  "id": ".+",
  "links": {
   "self": "https://argo-web-api.grnet.gr:443/api/v2/recomputations/.+"
  }
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(202, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Regexp(recomputationRequestsJSON, output, "Response body mismatch")

	var results []MongoInterface
	rCol := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection(recomputationsColl)
	findOptions := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})
	cursor, err := rCol.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	suite.Equal(len(results), 4)
	suite.Equal("2022-01-12T00:00:00Z", results[3].StartTime)
	suite.Equal("2022-01-15T00:00:00Z", results[3].EndTime)
	suite.Equal("requester", results[3].RequesterName)
	// check first applied status change item
	suite.Equal("", results[3].AppliedStatusChanges[0].Group)
	suite.Equal("", results[3].AppliedStatusChanges[0].Hostname)
	suite.Equal("Service_X", results[3].AppliedStatusChanges[0].Service)
	suite.Equal("Metric_X", results[3].AppliedStatusChanges[0].Metric)
	suite.Equal("CRITICAL", results[3].AppliedStatusChanges[0].State)
	// check second applied status change item
	suite.Equal("", results[3].AppliedStatusChanges[1].Group)
	suite.Equal("", results[3].AppliedStatusChanges[1].Hostname)
	suite.Equal("Service_X", results[3].AppliedStatusChanges[1].Service)
	suite.Equal("Metric_XX", results[3].AppliedStatusChanges[1].Metric)
	suite.Equal("OK", results[3].AppliedStatusChanges[1].State)
	// check second applied status change item
	suite.Equal("", results[3].AppliedStatusChanges[2].Group)
	suite.Equal("Hostname_XXX", results[3].AppliedStatusChanges[2].Hostname)
	suite.Equal("", results[3].AppliedStatusChanges[2].Service)
	suite.Equal("Metric_XXX", results[3].AppliedStatusChanges[2].Metric)
	suite.Equal("OK", results[3].AppliedStatusChanges[1].State)
	// check fourth applied status change item
	suite.Equal("Group_XXXX", results[3].AppliedStatusChanges[3].Group)
	suite.Equal("", results[3].AppliedStatusChanges[3].Hostname)
	suite.Equal("", results[3].AppliedStatusChanges[3].Service)
	suite.Equal("Metric_XXXX", results[3].AppliedStatusChanges[3].Metric)
	suite.Equal("WARNING", results[3].AppliedStatusChanges[3].State)
}

func (suite *RecomputationsProfileTestSuite) TestSubmitRecomputation2() {
	submission := IncomingRecomputation{
		StartTime:        "2015-01-10T12:00:00Z",
		EndTime:          "2015-01-30T23:00:00Z",
		RequesterName:    "Joe Complexz",
		RequesterEmail:   "C.Joecz@egi.eu",
		Reason:           "Ups failure",
		Report:           "EGI_Critical",
		ExcludeMonSource: []ExcludeMonSource{{Host: "MON1", StartTime: "2021-05-06T00:00:00Z", EndTime: "2021-05-06T05:00:00Z"}},
	}
	jsonsubmission, _ := json.Marshal(submission)

	request, _ := http.NewRequest("POST", "https://argo-web-api.grnet.gr:443/api/v2/recomputations", bytes.NewBuffer(jsonsubmission))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	recomputationRequestsJSON := `{
 "status": {
  "message": "Recomputations successfully created",
  "code": "201"
 },
 "data": {
  "id": ".+",
  "links": {
   "self": "https://argo-web-api.grnet.gr:443/api/v2/recomputations/.+"
  }
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(202, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Regexp(recomputationRequestsJSON, output, "Response body mismatch")

	var results []MongoInterface
	rCol := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection(recomputationsColl)
	findOptions := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})
	cursor, err := rCol.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	suite.Equal(len(results), 4)
	suite.Equal("2021-05-06T00:00:00Z", results[3].ExcludeMonSource[0].StartTime)
	suite.Equal("2021-05-06T05:00:00Z", results[3].ExcludeMonSource[0].EndTime)
	suite.Equal("MON1", results[3].ExcludeMonSource[0].Host)

}

func (suite *RecomputationsProfileTestSuite) TestSubmitRecomputationBadJSON() {

	// malformed json
	jsonSubmission := []byte("{{")

	request, _ := http.NewRequest("POST", "https://argo-web-api.grnet.gr:443/api/v2/recomputations", bytes.NewBuffer(jsonSubmission))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	recomputationRequestsJSON := `{
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
	suite.Equal(400, code)
	// Compare the expected and actual json response
	suite.Equal(recomputationRequestsJSON, output)
}

func (suite *RecomputationsProfileTestSuite) TestChangeAndResetStatus() {

	// malformed json
	jsonSubmission := []byte("{\"status\": \"approved\"}")
	jsonSubmission2 := []byte("{\"status\": \"running\"}")
	jsonSubmissionX := []byte("{\"status\": \"whatever\"}")
	jsonSubmission3 := []byte("{\"status\": \"done\"}")

	request, _ := http.NewRequest("POST", "https://argo-web-api.grnet.gr:443/api/v2/recomputations/6ac7d684-1f8e-4a02-a502-720e8f11e50b/status", bytes.NewBuffer(jsonSubmission))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	request2, _ := http.NewRequest("POST", "https://argo-web-api.grnet.gr:443/api/v2/recomputations/6ac7d684-1f8e-4a02-a502-720e8f11e50b/status", bytes.NewBuffer(jsonSubmission2))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)

	requestX, _ := http.NewRequest("POST", "https://argo-web-api.grnet.gr:443/api/v2/recomputations/6ac7d684-1f8e-4a02-a502-720e8f11e50b/status", bytes.NewBuffer(jsonSubmissionX))
	requestX.Header.Set("x-api-key", suite.clientkey)
	requestX.Header.Set("Accept", "application/json")
	responseX := httptest.NewRecorder()

	suite.router.ServeHTTP(responseX, requestX)

	request3, _ := http.NewRequest("POST", "https://argo-web-api.grnet.gr:443/api/v2/recomputations/6ac7d684-1f8e-4a02-a502-720e8f11e50b/status", bytes.NewBuffer(jsonSubmission3))
	request3.Header.Set("x-api-key", suite.clientkey)
	request3.Header.Set("Accept", "application/json")
	response3 := httptest.NewRecorder()

	suite.router.ServeHTTP(response3, request3)

	code1 := response.Code
	output1 := response.Body.String()

	code2 := response2.Code
	output2 := response2.Body.String()

	codeX := responseX.Code
	outputX := responseX.Body.String()

	code3 := response3.Code
	output3 := response3.Body.String()

	recomputationRequestsJSON1 := `{
 "status": {
  "message": "Recomputation status updated successfully to: approved",
  "code": "200"
 }
}`

	recomputationRequestsJSON2 := `{
 "status": {
  "message": "Recomputation status updated successfully to: running",
  "code": "200"
 }
}`

	recomputationRequestsJSONX := `{
 "status": {
  "message": "Conflict",
  "code": "409"
 },
 "errors": [
  {
   "message": "Conflict",
   "code": "409",
   "details": "status should be among values: \"pending\",\"approved\",\"rejected\",\"running\",\"done\""
  }
 ]
}`

	recomputationRequestsJSON3 := `{
 "status": {
  "message": "Recomputation status updated successfully to: done",
  "code": "200"
 }
}`

	suite.Equal(200, code1)
	// Compare the expected and actual json response
	suite.Equal(recomputationRequestsJSON1, output1)

	suite.Equal(200, code2)
	// Compare the expected and actual json response
	suite.Equal(recomputationRequestsJSON2, output2)

	suite.Equal(409, codeX)
	// Compare the expected and actual json response
	suite.Equal(recomputationRequestsJSONX, outputX)

	suite.Equal(200, code3)
	// Compare the expected and actual json response
	suite.Equal(recomputationRequestsJSON3, output3)

	// check that final recomputation contains correct history

	request4, _ := http.NewRequest("GET", "/api/v2/recomputations/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
	request4.Header.Set("x-api-key", suite.clientkey)
	request4.Header.Set("Accept", "application/json")
	response4 := httptest.NewRecorder()

	suite.router.ServeHTTP(response4, request4)

	code4 := response4.Code
	output4 := response4.Body.String()

	recomputationRequestsJSON4 := `\{
 "status": \{
  "message": "Success",
  "code": "200"
 \},
 "data": \{
  "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
  "requester_name": "John Snow",
  "requester_email": "jsnow@wall.com",
  "reason": "reasons",
  "start_time": "2015-03-10T12:00:00Z",
  "end_time": "2015-03-30T23:00:00Z",
  "report": "EGI_Critical",
  "exclude": \[
   "SITE1",
   "SITE3"
  \],
  "status": "done",
  "timestamp": "2015-04-01T14:58:40Z",
  "history": \[
   \{
    "status": "pending",
    "timestamp": ".*"
   \},
   \{
    "status": "approved",
    "timestamp": ".*"
   \},
   \{
    "status": "running",
    "timestamp": ".*"
   \},
   \{
    "status": "done",
    "timestamp": ".*"
   \}
  \]
 \}
\}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code4, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Regexp(recomputationRequestsJSON4, output4, "Response body mismatch")

	// Now reset status to see if history is cleared
	request5, _ := http.NewRequest("DELETE", "https://argo-web-api.grnet.gr:443/api/v2/recomputations/6ac7d684-1f8e-4a02-a502-720e8f11e50b/status", nil)
	request5.Header.Set("x-api-key", suite.clientkey)
	request5.Header.Set("Accept", "application/json")
	response5 := httptest.NewRecorder()

	suite.router.ServeHTTP(response5, request5)

	code5 := response5.Code
	output5 := response5.Body.String()

	recomputationRequestsJSON5 := `{
 "status": {
  "message": "Recomputation status reset successfully to: pending",
  "code": "200"
 }
}`

	suite.Equal(200, code5)
	suite.Equal(recomputationRequestsJSON5, output5)

	// Get again the recomputation json
	request6, _ := http.NewRequest("GET", "/api/v2/recomputations/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
	request6.Header.Set("x-api-key", suite.clientkey)
	request6.Header.Set("Accept", "application/json")
	response6 := httptest.NewRecorder()

	suite.router.ServeHTTP(response6, request6)
	recomputationRequestsJSON6 := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": {
  "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
  "requester_name": "John Snow",
  "requester_email": "jsnow@wall.com",
  "reason": "reasons",
  "start_time": "2015-03-10T12:00:00Z",
  "end_time": "2015-03-30T23:00:00Z",
  "report": "EGI_Critical",
  "exclude": [
   "SITE1",
   "SITE3"
  ],
  "status": "pending",
  "timestamp": "2015-04-01T14:58:40Z",
  "history": [
   {
    "status": "pending",
    "timestamp": "2015-04-01T14:58:40Z"
   }
  ]
 }
}`

	code6 := response6.Code
	output6 := response6.Body.String()

	suite.Equal(200, code6)
	suite.Equal(recomputationRequestsJSON6, output6)

}

// TearDownTest to tear down every test
func (suite *RecomputationsProfileTestSuite) TearDownTest() {

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

func (suite *RecomputationsProfileTestSuite) TestSubmitForbidViewer() {

	jsonInput := `{
  "name": "test",
  }`

	jsonOutput := `{
 "status": {
  "message": "Access to the resource is Forbidden",
  "code": "403"
 }
}`

	request, _ := http.NewRequest("POST", "/api/v2/recomputations", strings.NewReader(jsonInput))
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

func (suite *RecomputationsProfileTestSuite) TestDeleteRecomputation() {

	request, _ := http.NewRequest("DELETE", "https://argo-web-api.grnet.gr:443/api/v2/recomputations/6ac7d684-1f8e-4a02-a502-720e8f11e50b", nil)
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	jsonOutput := `{
 "status": {
  "message": "Recomputation Successfully Deleted",
  "code": "200"
 }
}`
	suite.Equal(200, code)
	// Compare the expected and actual json response
	suite.Equal(jsonOutput, output)

	// make sure the recomputation was really deleted from the store

	rCol := suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection(recomputationsColl)
	queryResult := rCol.FindOne(context.TODO(), bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b"})

	suite.Equal(queryResult.Err(), mongo.ErrNoDocuments)

}

func (suite *RecomputationsProfileTestSuite) TestDeleteRecomputationNotFound() {

	request, _ := http.NewRequest("DELETE", "https://argo-web-api.grnet.gr:443/api/v2/recomputations/unknown_id", nil)
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

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
	suite.Equal(404, code)
	// Compare the expected and actual json response
	suite.Equal(jsonOutput, output)
}

func (suite *RecomputationsProfileTestSuite) TestUpdateRecomputation() {
	submission := IncomingRecomputation{
		StartTime: "2015-01-10T12:00:00Z",
		EndTime:   "2015-01-30T23:00:00Z",
		Reason:    "Ups failure",
		Report:    "EGI_Critical",
		Exclude:   []string{"SITE5", "SITE8"},
	}
	jsonsubmission, _ := json.Marshal(submission)

	request, _ := http.NewRequest("PUT", "https://argo-web-api.grnet.gr:443/api/v2/recomputations/6ac7d684-1f8e-4a02-a502-720e8f11e50b", bytes.NewBuffer(jsonsubmission))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	jsonOutput := `{
 "status": {
  "message": "Recomputation updated successfully",
  "code": "200"
 }
}`

	suite.Equal(200, code)
	suite.Equal(jsonOutput, output)
}

func (suite *RecomputationsProfileTestSuite) TestUpdateRecomputationNotFound() {

	jsonSubmission := []byte("{}")

	request, _ := http.NewRequest("PUT", "https://argo-web-api.grnet.gr:443/api/v2/recomputations/unknown_id", bytes.NewBuffer(jsonSubmission))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

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
	suite.Equal(404, code)
	// Compare the expected and actual json response
	suite.Equal(jsonOutput, output)
}

func (suite *RecomputationsProfileTestSuite) TestUpdateRecomputationBadJSON() {

	// malformed json
	jsonSubmission := []byte("{{")

	request, _ := http.NewRequest("PUT", "https://argo-web-api.grnet.gr:443/api/v2/recomputations/6ac7d684-1f8e-4a02-a502-720e8f11e50b", bytes.NewBuffer(jsonSubmission))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	recomputationRequestsJSON := `{
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
	suite.Equal(400, code)
	// Compare the expected and actual json response
	suite.Equal(recomputationRequestsJSON, output)
}

// TearDownTest to tear down every test
func (suite *RecomputationsProfileTestSuite) TearDownSuite() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
}

func TestRecompuptationsTestSuite(t *testing.T) {
	suite.Run(t, new(RecomputationsProfileTestSuite))
}
