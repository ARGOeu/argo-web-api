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

package statusDetail

import (
	"net/http"
	"strings"
	"testing"

	"code.google.com/p/gcfg"
	"github.com/argoeu/argo-web-api/utils/authentication"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// This is a util. suite struct used in tests (see pkg "testify")
type StatusDetailTestSuite struct {
	suite.Suite
	cfg              config.Config
	tenantDbConf     config.MongoConfig
	respUnauthorized string
	respBadJSON      string
}

// Setup the Test Environment
// This function runs before any test and setups the environment
// A test configuration object is instantiated using a reference
// to testdb: argo_test_details. Also here is are instantiated some expected
// xml response validation messages (authorization,crud responses).
// Also the testdb is seeded with tenants,jobs,metric_profiles and status_metrics
func (suite *StatusDetailTestSuite) SetupTest() {

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
    db = "argotest_detail"
`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.respBadJSON = " <root>\n" +
		"   <Message>Malformated json input data</Message>\n </root>"

	suite.respUnauthorized = "Unauthorized"

	// Connect to mongo testdb
	session, _ := mongo.OpenSession(suite.cfg.MongoDB)

	// Add authentication token to mongo testdb
	seedAuth := bson.M{"api_key": "S3CR3T"}
	_ = mongo.Insert(session, suite.cfg.MongoDB.Db, "authentication", seedAuth)

	// seed mongo
	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// seed a tenant to use
	c := session.DB(suite.cfg.MongoDB.Db).C("tenants")
	c.Insert(bson.M{
		"name": "EGI",
		"db_conf": []bson.M{
			bson.M{
				"store":    "main",
				"server":   "localhost",
				"port":     27017,
				"database": "argotest_detail_egi",
				"username": "",
				"password": ""},
		},
		"users": []bson.M{
			bson.M{
				"name":    "egi_user",
				"email":   "egi_user@email.com",
				"api_key": "KEY1"},
		}})

	c.Insert(bson.M{
		"name": "EUDAT",
		"db_conf": []bson.M{
			bson.M{
				"store":    "main",
				"server":   "localhost",
				"port":     27017,
				"database": "argotest_detail_eudat",
				"username": "",
				"password": ""},
		},
		"users": []bson.M{
			bson.M{
				"name":    "eudat_user",
				"email":   "eudat_user@email.com",
				"api_key": "KEY2"},
		}})

	// get dbconfiguration based on the tenant
	// Prepare the request object
	request, _ := http.NewRequest("GET", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY1")
	// authenticate user's api key and find corresponding tenant
	suite.tenantDbConf, err = authentication.AuthenticateTenant(request.Header, suite.cfg)

	// Now seed the job DEFINITIONS
	c = session.DB(suite.tenantDbConf.Db).C("jobs")
	c.Insert(bson.M{
		"name":            "ROC_CRITICAL",
		"tenant":          "EGI",
		"endpoint_group":  "SITES",
		"group_of_groups": "NGI",
		"profiles": []bson.M{
			bson.M{
				"name":  "metric",
				"value": "ch.cern.sam.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			bson.M{
				"name":  "name1",
				"value": "value1"},
			bson.M{
				"name":  "name2",
				"value": "value2"},
		}})

	// Now seed the metric_profiles
	c = session.DB(suite.tenantDbConf.Db).C("metric_profiles")
	c.Insert(
		bson.M{
			"name": "ch.cern.SAM.ROC_CRITICAL",
			"services": []bson.M{
				bson.M{"service": "CREAM-CE",
					"metrics": []string{
						"emi.cream.CREAMCE-JobSubmit",
						"emi.wn.WN-Bi",
						"emi.wn.WN-Csh",
						"emi.wn.WN-SoftVer"},
				},
				bson.M{"service": "SRMv2",
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

	// seed the status detailed metric data
	c = session.DB(suite.tenantDbConf.Db).C("status_metric")
	c.Insert(bson.M{
		"job":                 "ROC_CRITICAL",
		"date_int":            20150501,
		"timestamp":           "2015-05-01T00:00:00Z",
		"supergroup":          "NGI_GRNET",
		"endpoint_group":      "HG-03-AUTH",
		"group_type":          "NGI",
		"endpoint_group_type": "SITES",
		"service":             "CREAM-CE",
		"hostname":            "cream01.afroditi.gr",
		"metric":              "emi.cream.CREAMCE-JobSubmit",
		"status":              "OK",
		"time_int":            0,
		"prev_status":         "OK",
		"prev_timestamp":      "2015-04-30T23:59:00Z",
	})
	c.Insert(bson.M{
		"job":                 "ROC_CRITICAL",
		"date_int":            20150501,
		"timestamp":           "2015-05-01T01:00:00Z",
		"supergroup":          "NGI_GRNET",
		"endpoint_group":      "HG-03-AUTH",
		"group_type":          "NGI",
		"endpoint_group_type": "SITES",
		"service":             "CREAM-CE",
		"hostname":            "cream01.afroditi.gr",
		"metric":              "emi.cream.CREAMCE-JobSubmit",
		"status":              "CRITICAL",
		"time_int":            10000,
		"prev_status":         "OK",
		"prev_timestamp":      "2015-05-01T00:00:00Z",
	})
	c.Insert(bson.M{
		"job":                 "ROC_CRITICAL",
		"date_int":            20150501,
		"timestamp":           "2015-05-01T05:00:00Z",
		"supergroup":          "NGI_GRNET",
		"endpoint_group":      "HG-03-AUTH",
		"group_type":          "NGI",
		"endpoint_group_type": "SITES",
		"service":             "CREAM-CE",
		"hostname":            "cream01.afroditi.gr",
		"metric":              "emi.cream.CREAMCE-JobSubmit",
		"status":              "OK",
		"time_int":            50000,
		"prev_status":         "CRITICAL",
		"prev_timestamp":      "2015-05-01T01:00:00Z",
	})

	// get dbconfiguration based on the tenant
	// Prepare the request object
	request, _ = http.NewRequest("GET", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY2")
	// authenticate user's api key and find corresponding tenant
	suite.tenantDbConf, err = authentication.AuthenticateTenant(request.Header, suite.cfg)

	// Now seed the job DEFINITIONS
	c = session.DB(suite.tenantDbConf.Db).C("jobs")
	c.Insert(bson.M{
		"name":            "EUDAT_CRITICAL",
		"tenant":          "EUDAT",
		"endpoint_group":  "EUDAT_SITE",
		"group_of_groups": "EUDAT_GROUP",
		"profiles": []bson.M{
			bson.M{
				"name":  "metric",
				"value": "eudat.CRITICAL"},
		},
		"filter_tags": []bson.M{
			bson.M{
				"name":  "name1",
				"value": "value1"},
			bson.M{
				"name":  "name2",
				"value": "value2"},
		}})

	// Now seed the metric_profiles
	c = session.DB(suite.tenantDbConf.Db).C("metric_profiles")
	c.Insert(
		bson.M{
			"name": "eudat.CRITICAL",
			"services": []bson.M{
				bson.M{"service": "srv.typeA",
					"metrics": []string{
						"typeA.metric.Memory",
						"typeA.metric.Disk"},
				},
				bson.M{"service": "srv.typeB",
					"metrics": []string{
						"typeB.metric.Memory",
						"typeB.metric.Disk"},
				},
			},
		})

	// seed the status detailed metric data
	c = session.DB(suite.tenantDbConf.Db).C("status_metric")
	c.Insert(bson.M{
		"job":                 "EUDAT_CRITICAL",
		"date_int":            20150501,
		"timestamp":           "2015-05-01T00:00:00Z",
		"supergroup":          "EUDAT_EL",
		"endpoint_group":      "EL-01-AUTH",
		"group_type":          "EUDAT_GROUP",
		"endpoint_group_type": "EUDAT_SITE",
		"service":             "srv.typeA",
		"hostname":            "host01.eudat.gr",
		"metric":              "typeA.metric.Memory",
		"status":              "OK",
		"time_int":            0,
		"prev_status":         "OK",
		"prev_timestamp":      "2015-04-30T23:59:00Z",
	})
	c.Insert(bson.M{
		"job":                 "EUDAT_CRITICAL",
		"date_int":            20150501,
		"timestamp":           "2015-05-01T01:00:00Z",
		"supergroup":          "EUDAT_EL",
		"endpoint_group":      "EL-01-AUTH",
		"group_type":          "EUDAT_GROUP",
		"endpoint_group_type": "EUDAT_SITE",
		"service":             "srv.typeA",
		"hostname":            "host01.eudat.gr",
		"metric":              "typeA.metric.Memory",
		"status":              "CRITICAL",
		"time_int":            10000,
		"prev_status":         "OK",
		"prev_timestamp":      "2015-05-01T00:00:00Z",
	})
	c.Insert(bson.M{
		"job":                 "EUDAT_CRITICAL",
		"date_int":            20150501,
		"timestamp":           "2015-05-01T05:00:00Z",
		"supergroup":          "EUDAT_EL",
		"endpoint_group":      "EL-01-AUTH",
		"group_type":          "EUDAT_GROUP",
		"endpoint_group_type": "EUDAT_SITE",
		"service":             "srv.typeA",
		"hostname":            "host01.eudat.gr",
		"metric":              "typeA.metric.Memory",
		"status":              "OK",
		"time_int":            50000,
		"prev_status":         "CRITICAL",
		"prev_timestamp":      "2015-05-01T01:00:00Z",
	})

}

func (suite *StatusDetailTestSuite) TestReadStatusDetail() {
	respXML := ` <root>
   <job name="EUDAT_CRITICAL">
     <group name="EUDAT_EL" type="EUDAT_GROUP">
       <group name="EL-01-AUTH" type="EUDAT_SITE">
         <group name="srv.typeA" type="service_type">
           <host name="host01.eudat.gr">
             <metric name="typeA.metric.Memory">
               <status timestamp="2015-04-30T23:59:00Z" status="OK"></status>
               <status timestamp="2015-05-01T01:00:00Z" status="CRITICAL"></status>
               <status timestamp="2015-05-01T05:00:00Z" status="OK"></status>
             </metric>
           </host>
         </group>
       </group>
     </group>
   </job>
 </root>`
	fullurl := "/api/v1/status/metrics/timeline/EUDAT_EL?" +
		"group_type=EUDAT_GROUP&start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z&job=EUDAT_CRITICAL"
	// Prepare the request object
	request, _ := http.NewRequest("GET", fullurl, strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY2")
	// Pass request to controller calling List() handler method
	code, _, output, _ := List(request, suite.cfg)
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Equal(respXML, string(output), "Response body mismatch")

}

// This function is actually called in the end of all tests
// and clears the test environment.
// Mainly it's purpose is to drop the testdb
func (suite *StatusDetailTestSuite) TearDownTest() {

	session, _ := mongo.OpenSession(suite.cfg.MongoDB)

	session.DB("argotest_detail").DropDatabase()
	session.DB("argotest_detail_eudat").DropDatabase()
	session.DB("argotest_detail_egi").DropDatabase()
}

// This is the first function called when go test is issued
func TestJobsSuite(t *testing.T) {
	suite.Run(t, new(StatusDetailTestSuite))
}
