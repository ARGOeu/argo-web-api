package metrics

import (
	"io/ioutil"
	"log"
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
type MetricsAdminTestSuite struct {
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

	log.SetOutput(ioutil.Discard)

	// seed mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	seedSuperAdmin := bson.M{"api_key": suite.clientkey}
	seedAuth := bson.M{"api_key": "S3CR3T"}
	seedResAuth := bson.M{"api_key": "R3STRICT3D", "restricted": true}
	seedResAdminUI := bson.M{"api_key": "ADM1NU1", "super_admin_ui": true}
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "authentication", seedSuperAdmin)
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "authentication", seedAuth)
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "authentication", seedResAuth)
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "authentication", seedResAdminUI)

	// Seed database with tenants
	c := session.DB(suite.cfg.MongoDB.Db).C("tenants")
	c.Insert(
		bson.M{"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"info": bson.M{
				"name":    "TENANT_A",
				"email":   "email@something2",
				"website": "tenant-b.example.com",
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
				"name":    "TENANT_B",
				"email":   "email@something2",
				"website": "tenanta.example.com",
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
					"api_key": "TESTKEY",
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
			"resource": "metrics_admin.get",
			"roles":    []string{"super_admin"},
		})
	c.Insert(
		bson.M{
			"resource": "metrics_admin.update",
			"roles":    []string{"super_admin"},
		})

	// Seed database with metrics
	c = session.DB(suite.cfg.MongoDB.Db).C("monitoring_metrics")
	c.Insert(
		bson.M{
			"name": "test_metric_1",
			"tags": []string{"network", "internal"},
		})

	c.Insert(
		bson.M{
			"name": "test_metric_2",
			"tags": []string{"disk", "agent"},
		})

	c.Insert(
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

//TearDownTest to tear down every test
func (suite *MetricsAdminTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

//TearDownTest to tear down every test
func (suite *MetricsAdminTestSuite) TearDownSuite() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

func TestMetricsAdminTestSuite(t *testing.T) {
	suite.Run(t, new(MetricsAdminTestSuite))
}
