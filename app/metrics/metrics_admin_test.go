package metrics

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
type MetricsAdminTestSuite struct {
	suite.Suite
	cfg              config.Config
	router           *mux.Router
	confHandler      respond.ConfHandler
	tenantDbConf     config.MongoConfig
	clientkey        string
	respUnauthorized string
}

// Setup the Test Environment
func (suite *MetricsAdminTestSuite) SetupSuite() {
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
	 db = "AR_test_metrics"
	 `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	client := store.GetMongoClient(suite.cfg.MongoDB)
	suite.cfg.MongoClient = client

	suite.respUnauthorized = "Unauthorized"
	suite.tenantDbConf = config.MongoConfig{
		Host:     "localhost",
		Port:     27017,
		Db:       "AR_test_metrics_tenant",
		Password: "pass",
		Username: "dbuser",
		Store:    "ar",
	}
	suite.clientkey = "SUPERADMINTOKEN"

	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v2/admin").Subrouter()
	HandleAdminSubrouter(suite.router, &suite.confHandler)
}

// This function runs before any test and setups the environment
func (suite *MetricsAdminTestSuite) SetupTest() {

	log.SetOutput(io.Discard)

	authCol := suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("authentication")

	seedSuperAdmin := bson.M{"api_key": suite.clientkey}
	seedAuth := bson.M{"api_key": "S3CR3T"}
	seedResAuth := bson.M{"api_key": "R3STRICT3D", "restricted": true}
	seedResAdminUI := bson.M{"api_key": "ADM1NU1", "super_admin_ui": true}
	authCol.InsertOne(context.TODO(), seedSuperAdmin)
	authCol.InsertOne(context.TODO(), seedAuth)
	authCol.InsertOne(context.TODO(), seedResAuth)
	authCol.InsertOne(context.TODO(), seedResAdminUI)

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
					"api_key": "TESTKEY",
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
			"resource": "metrics_admin.get",
			"roles":    []string{"super_admin"},
		})
	c.InsertOne(context.TODO(),
		bson.M{
			"resource": "metrics_admin.update",
			"roles":    []string{"super_admin"},
		})

	// Seed database with metrics
	c = suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Collection("monitoring_metrics")
	c.InsertOne(context.TODO(),
		bson.M{
			"name": "test_metric_1",
			"tags": []string{"network", "internal"},
		})

	c.InsertOne(context.TODO(),
		bson.M{
			"name": "test_metric_2",
			"tags": []string{"disk", "agent"},
		})

	c.InsertOne(context.TODO(),
		bson.M{
			"name": "test_metric_3",
			"tags": []string{"aai"},
		})

}

func (suite *MetricsAdminTestSuite) TestAdminUpdateMetricData() {

	jsonInput := `
  [{"name":"metric1","tags":["tag1", "tag2"]}]
`

	jsonOutput := `{
 "status": {
  "message": "Metrics resource succesfully updated",
  "code": "200"
 },
 "data": [
  {
   "name": "metric1",
   "tags": [
    "tag1",
    "tag2"
   ]
  }
 ]
}`

	request, _ := http.NewRequest("PUT", "/api/v2/admin/metrics", strings.NewReader(jsonInput))
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

func (suite *MetricsAdminTestSuite) TestAdminListMetrics() {

	request, _ := http.NewRequest("GET", "/api/v2/admin/metrics", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	metricsJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "name": "test_metric_1",
   "tags": [
    "network",
    "internal"
   ]
  },
  {
   "name": "test_metric_2",
   "tags": [
    "disk",
    "agent"
   ]
  },
  {
   "name": "test_metric_3",
   "tags": [
    "aai"
   ]
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(metricsJSON, output, "Response body mismatch")

}

// TearDownTest to tear down every test
func (suite *MetricsAdminTestSuite) TearDownTest() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
}

// TearDownTest to tear down every test
func (suite *MetricsAdminTestSuite) TearDownSuite() {

	suite.cfg.MongoClient.Database(suite.cfg.MongoDB.Db).Drop(context.TODO())
	suite.cfg.MongoClient.Database(suite.tenantDbConf.Db).Drop(context.TODO())
}

func TestSuiteMetricsAdmin(t *testing.T) {
	suite.Run(t, new(MetricsAdminTestSuite))
}
