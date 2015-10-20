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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gcfg.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
// This function runs before any test and setups the environment
func (suite *RecomputationsProfileTestSuite) SetupTest() {

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

	suite.confHandler = respond.ConfHandler{suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v2").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)

	// seed mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Seed database with tenants
	//TODO: move tests to
	c := session.DB(suite.cfg.MongoDB.Db).C("tenants")
	c.Insert(
		bson.M{"uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
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
					"database": "argo_Westeros1",
				},
				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_Westeros2",
				},
			},
			"users": []bson.M{

				bson.M{
					"name":    "John Snow",
					"email":   "J.Snow@brothers.wall",
					"api_key": "wh1t3_w@lk3rs",
				},
				bson.M{
					"name":    "King Joffrey",
					"email":   "g0dk1ng@kingslanding.gov",
					"api_key": "sansa <3",
				},
			}})
	c.Insert(
		bson.M{"uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50d",
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
					"name":    "Joe Complex",
					"email":   "C.Joe@egi.eu",
					"api_key": suite.clientkey,
				},
				bson.M{
					"name":    "Josh Plain",
					"email":   "P.Josh@egi.eu",
					"api_key": "itsamysterytoyou",
				},
			}})
	// Seed database with recomputations
	c = session.DB(suite.tenantDbConf.Db).C("recomputations")
	c.Insert(
		MongoInterface{
			UUID:           "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			RequesterName:  "John Snow",
			RequesterEmail: "jsnow@wall.com",
			StartTime:      "2015-03-10T12:00:00Z",
			EndTime:        "2015-03-30T23:00:00Z",
			Reason:         "reasons",
			Report:         "EGI_Critical",
			Exclude:        []string{"SITE1", "SITE3"},
			Status:         "pending",
			Timestamp:      "2015-04-01 14:58:40",
		},
	)
	c.Insert(
		MongoInterface{
			UUID:           "6ac7d684-1f8e-4a02-a502-720e8f11e50a",
			RequesterName:  "Arya Stark",
			RequesterEmail: "astark@shadowguild.com",
			StartTime:      "2015-01-10T12:00:00Z",
			EndTime:        "2015-01-30T23:00:00Z",
			Reason:         "power cuts",
			Report:         "EGI_Critical",
			Exclude:        []string{"SITE2", "SITE4"},
			Status:         "running",
			Timestamp:      "2015-02-01 14:58:40",
		},
	)

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
  "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50a",
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
  "timestamp": "2015-02-01 14:58:40"
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(recomputationRequestsJSON, output, "Response body mismatch")
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
   "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50a",
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
   "timestamp": "2015-02-01 14:58:40"
  },
  {
   "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
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
   "timestamp": "2015-04-01 14:58:40"
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(recomputationRequestsJSON, output, "Response body mismatch")
}

func (suite *RecomputationsProfileTestSuite) TestSubmitRecomputations() {
	submission := IncomingRecomputation{
		StartTime: "2015-01-10T12:00:00Z",
		EndTime:   "2015-01-30T23:00:00Z",
		Reason:    "Ups failure",
		Report:    "EGI_Critical",
		Exclude:   []string{"SITE5", "SITE8"},
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
  "uuid": ".+",
  "links": {
   "self": "https://argo-web-api.grnet.gr:443/api/v2/recomputations/.+"
  }
 }
}`
	// Check that we must have a 200 ok code
	suite.Equal(202, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Regexp(recomputationRequestsJSON, output, "Response body mismatch")

	dbDumpJson := `\[
 \{
  "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50a",
  "requester_name": "Arya Stark",
  "requester_email": "astark@shadowguild.com",
  "reason": "power cuts",
  "start_time": "2015-01-10T12:00:00Z",
  "end_time": "2015-01-30T23:00:00Z",
  "report": "EGI_Critical",
  "exclude": \[
   "SITE2",
   "SITE4"
  \],
  "status": "running",
  "timestamp": "2015-02-01 14:58:40"
 \},
 \{
  "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
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
  "status": "pending",
  "timestamp": "2015-04-01 14:58:40"
 \},
 \{
  "uuid": ".+-.+-.+-.+-.+",
  "requester_name": "Joe Complex",
  "requester_email": "C.Joe@egi.eu",
  "reason": "Ups failure",
  "start_time": "2015-01-10T12:00:00Z",
  "end_time": "2015-01-30T23:00:00Z",
  "report": "EGI_Critical",
  "exclude": \[
   "SITE5",
   "SITE8"
  \],
  "status": "pending",
  "timestamp": ".*"
 \}
\]`

	session, _ := mongo.OpenSession(suite.tenantDbConf)
	defer mongo.CloseSession(session)
	var results []MongoInterface
	mongo.Find(session, suite.tenantDbConf.Db, recomputationsColl, nil, "timestamp", &results)
	json, _ := json.MarshalIndent(results, "", " ")
	suite.Regexp(dbDumpJson, string(json), "Database contents were not expected")

}

//TearDownTest to tear down every test
func (suite *RecomputationsProfileTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

func TestRecompuptationsTestSuite(t *testing.T) {
	suite.Run(t, new(RecomputationsProfileTestSuite))
}
