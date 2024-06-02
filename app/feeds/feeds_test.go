package feeds

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

// This is a util. suite struct used in tests (see pkg "testify")
type FeedsTestSuite struct {
	suite.Suite
	cfg              config.Config
	router           *mux.Router
	confHandler      respond.ConfHandler
	tenantDbConf     config.MongoConfig
	clientkey        string
	respUnauthorized string
}

// Setup the Test Environment
func (suite *FeedsTestSuite) SetupSuite() {
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
	 db = "AR_test_feeds"
	 `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.respUnauthorized = "Unauthorized"
	suite.tenantDbConf = config.MongoConfig{
		Host:     "localhost",
		Port:     27017,
		Db:       "AR_test_feeds_tenant",
		Password: "pass",
		Username: "dbuser",
		Store:    "ar",
	}
	suite.clientkey = "123456"

	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v2/feeds").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)
}

// This function runs before any test and setups the environment
func (suite *FeedsTestSuite) SetupTest() {

	log.SetOutput(io.Discard)

	// Seed database with tenants
	c := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("tenants")
	c.InsertOne(context.TODO(),
		bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"info": bson.M{
				"name":    "TENANT_A",
				"email":   "email@something2",
				"website": "tenant-b.example.com",
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
				"name":    "TENANT_B",
				"email":   "email@something2",
				"website": "tenanta.example.com",
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
			"resource": "feeds.topo.get",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "feeds.topo.update",
			"roles":    []string{"editor"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "feeds.weights.get",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "feeds.weights.update",
			"roles":    []string{"editor"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "feeds.data.get",
			"roles":    []string{"editor", "viewer"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "feeds.data.update",
			"roles":    []string{"editor"},
		})

	// Seed database with topology feeds
	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("feeds_topology")
	c.InsertOne(context.TODO(),
		bson.M{
			"type":          "gocdb",
			"feed_url":      "https://somewhere.foo.bar/topology/feed",
			"paginated":     "true",
			"fetch_type":    []string{"item1", "item2"},
			"uid_endpoints": "endpointA",
		})

	// Seed database with weights feeds
	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("feeds_weights")
	c.InsertOne(context.TODO(),
		bson.M{
			"type":        "vapor",
			"feed_url":    "https://somewhere.foo.bar/weight/feed",
			"weight_type": "hepspec2006 cpu",
			"group_type":  "SITES",
		})

	// Seed database with weights feeds
	c = suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Collection("feeds_data")
	c.InsertOne(context.TODO(),
		bson.M{
			"tenants": []string{"6ac7d684-1f8e-4a02-a502-720e8f11e50c", "6ac7d684-1f8e-4a02-a502-720e8f11e50d"},
		})

}

func (suite *FeedsTestSuite) TestUpdateFeedData() {

	jsonInput := `
  {
  "tenants": ["TENANT_B"]
  }
`

	jsonOutput := `{
 "status": {
  "message": "Feeds resource succesfully updated",
  "code": "200"
 },
 "data": [
  {
   "tenants": [
    "TENANT_B"
   ]
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/feeds/data", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *FeedsTestSuite) TestUpdateFeedDataNotFound() {

	jsonInput := `
  {
  "tenants": ["TENANT_C"]
  }
`

	jsonOutput := `{
 "status": {
  "message": "Tenant TENANT_C not found",
  "code": "404"
 }
}`

	request, _ := http.NewRequest("PUT", "/api/v2/feeds/data", strings.NewReader(jsonInput))
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

func (suite *FeedsTestSuite) TestUpdateFeedWeights() {

	jsonInput := `
  {
   "type": "vapor",
   "feed_url": "https://somewhere2.foo.bar/weights/feed",
   "weight_type": "hepspec2006 memory",
   "group_type": "SITES"
  }
`

	jsonOutput := `{
 "status": {
  "message": "Feeds resource succesfully updated",
  "code": "200"
 },
 "data": [
  {
   "type": "vapor",
   "feed_url": "https://somewhere2.foo.bar/weights/feed",
   "weight_type": "hepspec2006 memory",
   "group_type": "SITES"
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/feeds/weights", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *FeedsTestSuite) TestListData() {

	request, _ := http.NewRequest("GET", "/api/v2/feeds/data", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	feedsTopoJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "tenants": [
    "TENANT_A",
    "TENANT_B"
   ]
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(feedsTopoJSON, output, "Response body mismatch")

}

func (suite *FeedsTestSuite) TestListWeights() {

	request, _ := http.NewRequest("GET", "/api/v2/feeds/weights", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	feedsTopoJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "type": "vapor",
   "feed_url": "https://somewhere.foo.bar/weight/feed",
   "weight_type": "hepspec2006 cpu",
   "group_type": "SITES"
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(feedsTopoJSON, output, "Response body mismatch")

}

func (suite *FeedsTestSuite) TestUpdateFeedTopo() {

	jsonInput := `
  {
   "type": "gocdb",
   "feed_url": "https://somewhere2.foo.bar/topology/feed",
   "paginated": "true",
   "fetch_type": [
    "item4",
    "item5"
   ],
   "uid_endpoints": "endpointA"
  }
`

	jsonOutput := `{
 "status": {
  "message": "Feeds resource succesfully updated",
  "code": "200"
 },
 "data": [
  {
   "type": "gocdb",
   "feed_url": "https://somewhere2.foo.bar/topology/feed",
   "paginated": "true",
   "fetch_type": [
    "item4",
    "item5"
   ],
   "uid_endpoints": "endpointA"
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/feeds/topology", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(jsonOutput, output, "Response body mismatch")

}

func (suite *FeedsTestSuite) TestListTopo() {

	request, _ := http.NewRequest("GET", "/api/v2/feeds/topology", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	feedsTopoJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "type": "gocdb",
   "feed_url": "https://somewhere.foo.bar/topology/feed",
   "paginated": "true",
   "fetch_type": [
    "item1",
    "item2"
   ],
   "uid_endpoints": "endpointA"
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(feedsTopoJSON, output, "Response body mismatch")

}

// TearDownTest do things after each test ends
func (suite *FeedsTestSuite) TearDownTest() {

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

// TearDownSuite do things after suite ends
func (suite *FeedsTestSuite) TearDownSuite() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
}

func TestSuiteFeeds(t *testing.T) {
	suite.Run(t, new(FeedsTestSuite))
}
