/*
 * Copyright (c) 2014 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
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
 * or implied, of either GRNET S.A., SRCE or IN2P3 CNRS Computing
 * Centre
 *
 * The work represented by this source file is partially funded by
 * the EGI-InSPIRE project through the European Commission's 7th
 * Framework Programme (contract # INFSO-RI-261323)
 */

package statusEndpointGroups

import (
	"net/http"
	"strings"
	"testing"

	"code.google.com/p/gcfg"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type StatusEndpointGroupTestSuite struct {
	suite.Suite
	cfg                       config.Config
	router                    mux.Router
	tenantDbConf              config.MongoConfig
	clientkey                 string
	respRecomputationsCreated string
	respUnauthorized          string
}

// SetupTest adds the required entries in the database and
// give the required values to the StatusEndpointGroupTestSuite struct
func (suite *StatusEndpointGroupTestSuite) SetupTest() {

	const testConfig = `
    [server]
    bindip = ""
    port = 8080
    maxprocs = 4
    cache = false
    lrucache = 700000000
    gzip = true
    [mongodb]
    host = "127.0.0.1"
    port = 27017
    db = "AR_test_core_status_site"
    `

	suite.respUnauthorized = "Unauthorized"
	suite.clientkey = "mysecretcombination"
	suite.tenantDbConf = config.MongoConfig{
		Host:     "127.0.0.1",
		Port:     27017,
		Db:       "argo_egi_test_status_site",
		Username: "johndoe",
		Password: "h4shp4ss",
		Store:    "ar",
	}

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	session, err := mongo.OpenSession(suite.cfg.MongoDB)

	if err != nil {
		panic(err)
	}
	defer mongo.CloseSession(session)

	c := session.DB(suite.cfg.MongoDB.Db).C("tenants")
	c.Insert(
		bson.M{"name": "Westeros",
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
		bson.M{"name": "EGI",
			"db_conf": []bson.M{

				bson.M{
					// "store":    "ar",
					"server":   "localhost",
					"port":     27017,
					"database": suite.tenantDbConf.Db,
					"username": suite.tenantDbConf.Username,
					"password": suite.tenantDbConf.Password,
				},
				bson.M{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_egi_metric_data",
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

	c = session.DB(suite.tenantDbConf.Db).C("status_endpointgroups")
	c.Insert(
		bson.M{
			"job":             "JOB_A",
			"supergroup":      "NGI_GRNET",
			"name":            "GR-01-AUTH",
			"timestamp":       "2015-01-06T00:00:41Z",
			"status":          "OK",
			"previous_status": "OK",
			"date_integer":    20150106,
			"time_integer":    41,
		})
	c.Insert(
		bson.M{
			"job":             "JOB_A",
			"supergroup":      "NGI_GRNET",
			"name":            "GR-01-AUTH",
			"timestamp":       "2015-01-06T00:05:00Z",
			"status":          "CRITICAL",
			"previous_status": "OK",
			"date_integer":    20150106,
			"time_integer":    500,
		})
	c.Insert(
		bson.M{
			"job":             "JOB_A",
			"supergroup":      "NGI_GRNET",
			"name":            "GR-01-AUTH",
			"timestamp":       "2015-01-06T00:12:00Z",
			"status":          "OK",
			"previous_status": "CRITICAL",
			"date_integer":    20150106,
			"time_integer":    1200,
		})
}

//TestListStatusEndpointGroup tests the correct formatting when listing Sites' statuses
func (suite *StatusEndpointGroupTestSuite) TestListStatusEndpointGroup() {
	query := "?start_time=2015-01-06T00:00:00Z&end_time=2015-01-06T23:59:59Z&job=JOB_A&supergroup_name=NGI_GRNET"
	request, _ := http.NewRequest("GET", "/api/v1/status/sites/timeline/GR-01-AUTH"+query, strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)

	code, _, output, _ := List(request, suite.cfg)
	statusEndpointGroupRequestXML := ` <root>
   <Job name="JOB_A">
     <EndpointGroup name="GR-01-AUTH">
       <Status timestamp="2015-01-06T00:00:41Z" Status="OK" PreviousStatus="OK"></Status>
       <Status timestamp="2015-01-06T00:05:00Z" Status="CRITICAL" PreviousStatus="OK"></Status>
       <Status timestamp="2015-01-06T00:12:00Z" Status="OK" PreviousStatus="CRITICAL"></Status>
     </EndpointGroup>
   </Job>
 </root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Regexp(statusEndpointGroupRequestXML, string(output), "Response body mismatch")
}

//TearDownTest to tear down every test
func (suite *StatusEndpointGroupTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()

}

func TestRecompuptationsTestSuite(t *testing.T) {
	suite.Run(t, new(StatusEndpointGroupTestSuite))
}
