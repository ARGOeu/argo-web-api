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

package metricProfiles

import (
	"net/http"
	"strings"
	"testing"

	"code.google.com/p/gcfg"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type MetricProfilesTestSuite struct {
	suite.Suite
	cfg                       config.Config
	router                    mux.Router
	tenantDbConf              config.MongoConfig
	clientkey                 string
	respRecomputationsCreated string
	respUnauthorized          string
}

// SetupTest adds the required entries in the database and
// give the required values to the MetricProfilesTestSuite struct
func (suite *MetricProfilesTestSuite) SetupTest() {

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
    db = "AR_test_core_metric_profiles"
    `

	suite.respUnauthorized = "Unauthorized"
	suite.clientkey = "mysecretcombination"
	suite.tenantDbConf.Db = "argo_egi_test_metric_profiles"
	suite.tenantDbConf.Password = "h4shp4ss"
	suite.tenantDbConf.Username = "johndoe"
	suite.tenantDbConf.Store = "ar"

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	session, err := mongo.OpenSession(suite.cfg.MongoDB)

	if err != nil {
		panic(err)
	}
	defer mongo.CloseSession(session)

	// Seed database with tenants
	//TODO: move tests to
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
	c.Insert(
		bson.M{
			"name": "ch.cern.SAM.ROC",
			"services": []bson.M{
				bson.M{"service": "CREAM-CE",
					"metrics": []string{
						"emi.cream.CREAMCE-JobSubmit",
						"emi.wn.WN-Bi",
						"emi.wn.WN-Csh",
						"hr.srce.CADist-Check",
						"hr.srce.CREAMCE-CertLifetime",
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

}

//TestListMetricProfiles tests the correct formatting when listing Metric Profiles
func (suite *MetricProfilesTestSuite) TestListMetricProfiles() {

	request, _ := http.NewRequest("GET", "/api/v1/metric_profiles", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)

	code, _, output, _ := List(request, suite.cfg)

	metricProfileRequestXML := `<root>
 <MetricProfiles id=".*" name="ch.cern.SAM.ROC">
  <services service="CREAM-CE">
   <metrics>emi.cream.CREAMCE-JobSubmit</metrics>
   <metrics>emi.wn.WN-Bi</metrics>
   <metrics>emi.wn.WN-Csh</metrics>
   <metrics>hr.srce.CADist-Check</metrics>
   <metrics>hr.srce.CREAMCE-CertLifetime</metrics>
   <metrics>emi.wn.WN-SoftVer</metrics>
  </services>
  <services service="SRMv2">
   <metrics>hr.srce.SRM2-CertLifetime</metrics>
   <metrics>org.sam.SRM-Del</metrics>
   <metrics>org.sam.SRM-Get</metrics>
   <metrics>org.sam.SRM-GetSURLs</metrics>
   <metrics>org.sam.SRM-GetTURLs</metrics>
   <metrics>org.sam.SRM-Ls</metrics>
   <metrics>org.sam.SRM-LsDir</metrics>
   <metrics>org.sam.SRM-Put</metrics>
  </services>
 </MetricProfiles>
 <MetricProfiles id=".*" name="ch.cern.SAM.ROC_CRITICAL">
  <services service="CREAM-CE">
   <metrics>emi.cream.CREAMCE-JobSubmit</metrics>
   <metrics>emi.wn.WN-Bi</metrics>
   <metrics>emi.wn.WN-Csh</metrics>
   <metrics>emi.wn.WN-SoftVer</metrics>
  </services>
  <services service="SRMv2">
   <metrics>hr.srce.SRM2-CertLifetime</metrics>
   <metrics>org.sam.SRM-Del</metrics>
   <metrics>org.sam.SRM-Get</metrics>
   <metrics>org.sam.SRM-GetSURLs</metrics>
   <metrics>org.sam.SRM-GetTURLs</metrics>
   <metrics>org.sam.SRM-Ls</metrics>
   <metrics>org.sam.SRM-LsDir</metrics>
   <metrics>org.sam.SRM-Put</metrics>
  </services>
 </MetricProfiles>
</root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Regexp(metricProfileRequestXML, string(output), "Response body mismatch")
}

// TestCreateMetricProfiles tests the Create method of the metricProfiles package
func (suite *MetricProfilesTestSuite) TestCreateMetricProfiles() {

	postData := `
	{
	"name" : "ch.cern.BOB.ROCK_AND_ROLL",
	"services" : [
		{ "service" : "CREAM-CE",
		  "metrics" : ["emi.cream.CREAMCE-JobSubmit", "emi.wn.WN-Bi", "emi.wn.WN-Cs"]
		},
		{
		  "service" : "SRMv2",
		  "metrics" : ["org.sam.SRM-Del","org.sam.SRM-Get","org.sam.SRM-GetSURLs"]
		}
	]}`

	request, _ := http.NewRequest("POST", "/api/v1/metric_profiles", strings.NewReader(postData))
	request.Header.Set("x-api-key", suite.clientkey)

	code, _, output, _ := Create(request, suite.cfg)

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	suite.Equal("Metric profile successfully inserted", string(output), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v1/metric_profiles", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)

	code, _, output, _ = List(request, suite.cfg)

	metricProfileRequestXML := `<root>
 <MetricProfiles id=".*" name="ch.cern.BOB.ROCK_AND_ROLL">
  <services service="CREAM-CE">
   <metrics>emi.cream.CREAMCE-JobSubmit</metrics>
   <metrics>emi.wn.WN-Bi</metrics>
   <metrics>emi.wn.WN-Cs</metrics>
  </services>
  <services service="SRMv2">
   <metrics>org.sam.SRM-Del</metrics>
   <metrics>org.sam.SRM-Get</metrics>
   <metrics>org.sam.SRM-GetSURLs</metrics>
  </services>
 </MetricProfiles>
 <MetricProfiles id=".*" name="ch.cern.SAM.ROC">
  <services service="CREAM-CE">
   <metrics>emi.cream.CREAMCE-JobSubmit</metrics>
   <metrics>emi.wn.WN-Bi</metrics>
   <metrics>emi.wn.WN-Csh</metrics>
   <metrics>hr.srce.CADist-Check</metrics>
   <metrics>hr.srce.CREAMCE-CertLifetime</metrics>
   <metrics>emi.wn.WN-SoftVer</metrics>
  </services>
  <services service="SRMv2">
   <metrics>hr.srce.SRM2-CertLifetime</metrics>
   <metrics>org.sam.SRM-Del</metrics>
   <metrics>org.sam.SRM-Get</metrics>
   <metrics>org.sam.SRM-GetSURLs</metrics>
   <metrics>org.sam.SRM-GetTURLs</metrics>
   <metrics>org.sam.SRM-Ls</metrics>
   <metrics>org.sam.SRM-LsDir</metrics>
   <metrics>org.sam.SRM-Put</metrics>
  </services>
 </MetricProfiles>
 <MetricProfiles id=".*" name="ch.cern.SAM.ROC_CRITICAL">
  <services service="CREAM-CE">
   <metrics>emi.cream.CREAMCE-JobSubmit</metrics>
   <metrics>emi.wn.WN-Bi</metrics>
   <metrics>emi.wn.WN-Csh</metrics>
   <metrics>emi.wn.WN-SoftVer</metrics>
  </services>
  <services service="SRMv2">
   <metrics>hr.srce.SRM2-CertLifetime</metrics>
   <metrics>org.sam.SRM-Del</metrics>
   <metrics>org.sam.SRM-Get</metrics>
   <metrics>org.sam.SRM-GetSURLs</metrics>
   <metrics>org.sam.SRM-GetTURLs</metrics>
   <metrics>org.sam.SRM-Ls</metrics>
   <metrics>org.sam.SRM-LsDir</metrics>
   <metrics>org.sam.SRM-Put</metrics>
  </services>
 </MetricProfiles>
</root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Regexp(metricProfileRequestXML, string(output), "Response body mismatch")

}

// TestUpdateMetricProfiles test the Update function of the metricProfiles package
func (suite *MetricProfilesTestSuite) TestUpdateMetricProfiles() {

	putData := `
	{
	"name" : "ch.cern.BOB.ROCK_AND_ROLL",
	"services" : [
		{ "service" : "CREAM-CE",
		  "metrics" : ["emi.cream.CREAMCE-JobSubmit", "emi.wn.WN-Bi", "emi.wn.WN-Cs"]
		},
		{
		  "service" : "SRMv2",
		  "metrics" : ["org.sam.SRM-Del","org.sam.SRM-Get","org.sam.SRM-GetSURLs"]
		}
	]}`

	session, err := mongo.OpenSession(suite.cfg.MongoDB)

	if err != nil {
		panic(err)
	}
	defer mongo.CloseSession(session)

	result := MongoInterface{}
	c := session.DB(suite.tenantDbConf.Db).C("metric_profiles")
	c.Find(bson.M{}).One(&result)

	request, _ := http.NewRequest("PUT", "/api/v1/metric_profiles/"+result.ID.Hex(), strings.NewReader(putData))
	request.Header.Set("x-api-key", suite.clientkey)
	context.Set(request, "id", result.ID.Hex())
	code, _, output, _ := Update(request, suite.cfg)

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	suite.Equal("Metric profile successfully updated", string(output), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v1/metric_profiles", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)

	code, _, output, _ = List(request, suite.cfg)

	metricProfileRequestXML := `<root>
 <MetricProfiles id=".*" name="ch.cern.BOB.ROCK_AND_ROLL">
  <services service="CREAM-CE">
   <metrics>emi.cream.CREAMCE-JobSubmit</metrics>
   <metrics>emi.wn.WN-Bi</metrics>
   <metrics>emi.wn.WN-Cs</metrics>
  </services>
  <services service="SRMv2">
   <metrics>org.sam.SRM-Del</metrics>
   <metrics>org.sam.SRM-Get</metrics>
   <metrics>org.sam.SRM-GetSURLs</metrics>
  </services>
 </MetricProfiles>
 <MetricProfiles id=".*" name="ch.cern.SAM.ROC">
  <services service="CREAM-CE">
   <metrics>emi.cream.CREAMCE-JobSubmit</metrics>
   <metrics>emi.wn.WN-Bi</metrics>
   <metrics>emi.wn.WN-Csh</metrics>
   <metrics>hr.srce.CADist-Check</metrics>
   <metrics>hr.srce.CREAMCE-CertLifetime</metrics>
   <metrics>emi.wn.WN-SoftVer</metrics>
  </services>
  <services service="SRMv2">
   <metrics>hr.srce.SRM2-CertLifetime</metrics>
   <metrics>org.sam.SRM-Del</metrics>
   <metrics>org.sam.SRM-Get</metrics>
   <metrics>org.sam.SRM-GetSURLs</metrics>
   <metrics>org.sam.SRM-GetTURLs</metrics>
   <metrics>org.sam.SRM-Ls</metrics>
   <metrics>org.sam.SRM-LsDir</metrics>
   <metrics>org.sam.SRM-Put</metrics>
  </services>
 </MetricProfiles>
</root>`

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Regexp(metricProfileRequestXML, string(output), "Response body mismatch")

}

// TestDeleteMetricProfiles test the Delete function of the metricProfiles package
func (suite *MetricProfilesTestSuite) TestDeleteMetricProfiles() {

	session, err := mongo.OpenSession(suite.cfg.MongoDB)

	if err != nil {
		panic(err)
	}
	defer mongo.CloseSession(session)

	result := MongoInterface{}
	c := session.DB(suite.tenantDbConf.Db).C("metric_profiles")
	c.Find(bson.M{}).One(&result)

	request, _ := http.NewRequest("DELETE", "/api/v1/metric_profiles/"+result.ID.Hex(), strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	// context.Set(request, "id", result.ID.Hex())
	code, _, output, _ := Delete(request, suite.cfg)

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	suite.Equal("Metric profile successfully removed", string(output), "Response body mismatch")

	request, _ = http.NewRequest("GET", "/api/v1/metric_profiles", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)

	code, _, output, _ = List(request, suite.cfg)

	metricProfileRequestXML := `<root>
 <MetricProfiles id=".*" name="ch.cern.SAM.ROC">
  <services service="CREAM-CE">
   <metrics>emi.cream.CREAMCE-JobSubmit</metrics>
   <metrics>emi.wn.WN-Bi</metrics>
   <metrics>emi.wn.WN-Csh</metrics>
   <metrics>hr.srce.CADist-Check</metrics>
   <metrics>hr.srce.CREAMCE-CertLifetime</metrics>
   <metrics>emi.wn.WN-SoftVer</metrics>
  </services>
  <services service="SRMv2">
   <metrics>hr.srce.SRM2-CertLifetime</metrics>
   <metrics>org.sam.SRM-Del</metrics>
   <metrics>org.sam.SRM-Get</metrics>
   <metrics>org.sam.SRM-GetSURLs</metrics>
   <metrics>org.sam.SRM-GetTURLs</metrics>
   <metrics>org.sam.SRM-Ls</metrics>
   <metrics>org.sam.SRM-LsDir</metrics>
   <metrics>org.sam.SRM-Put</metrics>
  </services>
 </MetricProfiles>
</root>`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual xml response
	suite.Regexp(metricProfileRequestXML, string(output), "Response body mismatch")

}

//TearDownTest to tear down every test
func (suite *MetricProfilesTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()

}

func TestRecompuptationsTestSuite(t *testing.T) {
	suite.Run(t, new(MetricProfilesTestSuite))
}
