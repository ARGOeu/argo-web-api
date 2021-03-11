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

package reports

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/authentication"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gcfg.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// This is a util. suite struct used in tests (see pkg "testify")
type ReportTestSuite struct {
	suite.Suite
	cfg                config.Config
	tenantDbConf       config.MongoConfig
	router             *mux.Router
	confHandler        respond.ConfHandler
	respReportCreated  string
	respReportUpdated  string
	respReportDeleted  string
	respReportNotFound string
	respUnauthorized   string
	respBadJSON        string
	respNameConflict   string
}

// Setup the Test Environment
func (suite *ReportTestSuite) SetupSuite() {

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
	    db = "argo_test_reports2"
	`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.respReportCreated = `{
 "status": {
  "message": "Successfully Created Report",
  "code": "201"
 },
 "data": {
  "id": ".+-.+-.+-.+-.+",
  "links": {
   "self": "https://myapi.test.com/api/v2/reports/.+-.+-.+-.+-.+"
  }
 }
}`

	suite.respReportUpdated = `{
 "status": {
  "message": "Report was successfully updated",
  "code": "200"
 }
}`

	suite.respReportDeleted = `{
 "status": {
  "message": "Report was successfully deleted",
  "code": "200"
 }
}`

	suite.respReportNotFound = `{
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

	suite.respNameConflict = `{
 "status": {
  "message": "Conflict",
  "code": "409"
 },
 "errors": [
  {
   "message": "Conflict",
   "code": "409",
   "details": "Report with the same name already exists"
  }
 ]
}`

	suite.respBadJSON = `{
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

	suite.respUnauthorized = `{
 "status": {
  "message": "Unauthorized",
  "code": "401",
  "details": "You need to provide a correct authentication token using the header 'x-api-key'"
 }
}`
	suite.confHandler = respond.ConfHandler{
		Config: suite.cfg,
	}
	suite.router = mux.NewRouter().StrictSlash(true).PathPrefix("/api/v2").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)
}

// This function runs before any test and setups the environment
// A test configuration object is instantiated using a reference
// to testdb: argo_test_reports. Also here is are instantiated some expected
// xml response validation messages (authorization,crud responses).
// Also the testdb is seeded with two reports
func (suite *ReportTestSuite) SetupTest() {

	log.SetOutput(ioutil.Discard)

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
		"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
		"info": bson.M{
			"name":    "GUARDIANS",
			"email":   "email@something2",
			"website": "www.gotg.com",
			"created": "2015-10-20 02:08:04",
			"updated": "2015-10-20 02:08:04"},
		"db_conf": []bson.M{
			bson.M{
				"store":    "ar",
				"server":   "localhost",
				"port":     27017,
				"database": "argo_test_tenant_reports2_db1",
				"username": "admin",
				"password": "3NCRYPT3D"},
			bson.M{
				"store":    "status",
				"server":   "b.mongodb.org",
				"port":     27017,
				"database": "reports_db_tenant",
				"username": "admin",
				"password": "3NCRYPT3D"},
		},
		"users": []bson.M{
			bson.M{
				"name":    "cap",
				"email":   "cap@email.com",
				"roles":   []string{"editor"},
				"api_key": "C4PK3Y"},
			bson.M{
				"name":    "thor",
				"email":   "thor@email.com",
				"roles":   []string{"viewer"},
				"api_key": "VIEWERKEY"},
		}})

	c = session.DB(suite.cfg.MongoDB.Db).C("roles")
	c.Insert(
		bson.M{
			"resource": "reports.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "reports.get",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "reports.create",
			"roles":    []string{"editor"},
		})
	c.Insert(
		bson.M{
			"resource": "reports.delete",
			"roles":    []string{"editor"},
		})
	c.Insert(
		bson.M{
			"resource": "reports.update",
			"roles":    []string{"editor"},
		})

	// get dbconfiguration based on the tenant
	// Prepare the request object
	request, _ := http.NewRequest("GET", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json;")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")
	// authenticate user's api key and find corresponding tenant
	suite.tenantDbConf, _, err = authentication.AuthenticateTenant(request.Header, suite.cfg)

	c = session.DB(suite.tenantDbConf.Db).C("metric_profiles")

	c.Insert(
		bson.M{
			"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
			"name": "profile1",
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
			"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
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

	c = session.DB(suite.tenantDbConf.Db).C("aggregation_profiles")
	c.Insert(
		bson.M{
			"id":                "6ac7d684-1f8e-4a02-a502-720e8f11e50bq",
			"name":              "profile3",
			"namespace":         "test",
			"endpoint_group":    "sites",
			"metric_operation":  "AND",
			"profile_operation": "AND",
			"metric_profile": bson.M{
				"name": "roc.critical",
				"id":   "5637d684-1f8e-4a02-a502-720e8f11e432",
			},
			"groups": []bson.M{
				bson.M{"name": "compute",
					"operation": "OR",
					"services": []bson.M{
						bson.M{
							"name":      "CREAM-CE",
							"operation": "AND",
						},
						bson.M{
							"name":      "ARC-CE",
							"operation": "AND",
						},
					}},
				bson.M{"name": "storage",
					"operation": "OR",
					"services": []bson.M{
						bson.M{
							"name":      "SRMv2",
							"operation": "AND",
						},
						bson.M{
							"name":      "SRM",
							"operation": "AND",
						},
					}},
			}})
	c.Insert(
		bson.M{
			"id":                "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
			"name":              "cloud",
			"namespace":         "test",
			"endpoint_group":    "sites",
			"metric_operation":  "AND",
			"profile_operation": "AND",
			"metric_profile": bson.M{
				"name": "roc.critical",
				"id":   "5637d684-1f8e-4a02-a502-720e8f11e432",
			},
			"groups": []bson.M{
				bson.M{"name": "compute",
					"operation": "OR",
					"services": []bson.M{
						bson.M{
							"name":      "SERVICEA",
							"operation": "AND",
						},
						bson.M{
							"name":      "SERVICEB",
							"operation": "AND",
						},
					}},
				bson.M{"name": "images",
					"operation": "OR",
					"services": []bson.M{
						bson.M{
							"name":      "SERVICEC",
							"operation": "AND",
						},
						bson.M{
							"name":      "SERVICED",
							"operation": "AND",
						},
					}},
			}})

	c = session.DB(suite.tenantDbConf.Db).C("operations_profiles")
	c.Insert(
		bson.M{
			"id":               "6ac7d684-1f8e-4a02-a502-720e8f11e523",
			"name":             "profile2",
			"available_states": []string{"A,B,C"},
			"defaults": bson.M{
				"missing": "A",
				"down":    "B",
				"unknown": "C"},
			"operations": []bson.M{
				bson.M{
					"name": "AND",
					"truth_table": []bson.M{
						bson.M{
							"a": "A",
							"b": "B",
							"x": "B",
						},
						bson.M{
							"a": "A",
							"b": "C",
							"x": "C",
						},
						bson.M{
							"a": "B",
							"b": "C",
							"x": "C",
						}}},
				bson.M{
					"name": "OR",
					"truth_table": []bson.M{
						bson.M{
							"a": "A",
							"b": "B",
							"x": "A",
						},
						bson.M{
							"a": "A",
							"b": "C",
							"x": "A",
						},
						bson.M{
							"a": "B",
							"b": "C",
							"x": "B",
						}}},
			}})

	// Seed database with weights
	c = session.DB(suite.tenantDbConf.Db).C("weights")
	c.EnsureIndexKey("-date_integer", "id")
	c.Insert(
		bson.M{
			"id":           "6ac7d684-1f8e-4a02-a502-720e8f11e533",
			"name":         "Critical",
			"date_integer": 20191004,
			"date":         "2019-10-04",
			"weight_type":  "hepsepc",
			"group_type":   "SITES",
			"groups": []bson.M{
				bson.M{"name": "SITE-A", "value": 1673},
				bson.M{"name": "SITE-B", "value": 1234},
				bson.M{"name": "SITE-C", "value": 523},
				bson.M{"name": "SITE-D", "value": 2},
			},
		})

	// Now seed the report DEFINITIONS
	c = session.DB(suite.tenantDbConf.Db).C(reportsColl)
	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"info": bson.M{
			"name":        "Report_A",
			"description": "report aaaaa",
			"created":     "2015-9-10 13:43:00",
			"updated":     "2015-10-11 13:43:00",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "NGI",
				"group": bson.M{
					"type": "SITE",
				},
			},
		},
		"profiles": []bson.M{
			bson.M{
				"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
				"type": "metric",
				"name": "profile1"},
			bson.M{
				"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e523",
				"type": "operations",
				"name": "profile2"},
			bson.M{
				"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e50q",
				"type": "aggregation",
				"name": "profile3"},
		},
		"filter_tags": []bson.M{
			bson.M{
				"name":  "name1",
				"value": "value1"},
			bson.M{
				"name":  "name2",
				"value": "value2"},
		}})

	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a494360",
		"info": bson.M{
			"name":        "Report_B",
			"description": "report bbb",
			"created":     "2015-10-08 13:43:00",
			"updated":     "2015-10-09 13:43:00",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "ARCHIPELAGO",
				"group": bson.M{
					"type": "ISLAND",
				},
			},
		},
		"profiles": []bson.M{
			bson.M{
				"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
				"type": "metric",
				"name": "profile1"},
			bson.M{
				"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e523",
				"type": "operations",
				"name": "profile2"},
			bson.M{
				"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e50q",
				"type": "aggregation",
				"name": "profile3"},
		},
		"filter_tags": []bson.M{
			bson.M{
				"name":  "name1",
				"value": "value1"},
			bson.M{
				"name":  "name2",
				"value": "value2"},
		}})
}

func (suite *ReportTestSuite) TestCreateReportWrongProfileType() {

	// create json input data for the request
	postData := `{
    "info": {
        "name": "Foo_Report",
        "description": "olalala"
    },
    "topology_schema": {
        "group": {
            "type": "ngi",
            "group": {
                "type": "site"
            }
        }
    },
	"profiles": [
        {
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
            "type": "metric_wrong"
            
        },
		{
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e523",
            "type": "operations"
         
        },
        {
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50bq",
            "type": "aggregation"
            
        },
		{
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e533",
			"type": "weight_wrong"
		}
    ],
    "filter_tags": [
        {
            "name": "production",
            "value": "Y"
        },
        {
            "name": "monitored",
            "value": "Y"
        }
    ]
}`

	responseJSON := `{
 "status": {
  "message": "Unprocessable Entity",
  "code": "422"
 },
 "errors": [
  {
   "message": "Profile type invalid",
   "code": "422",
   "details": "Profile type metric_wrong is invalid"
  },
  {
   "message": "Profile type invalid",
   "code": "422",
   "details": "Profile type weight_wrong is invalid"
  }
 ]
}`

	// Prepare the request object
	request, _ := http.NewRequest("POST", "https://myapi.test.com/api/v2/reports", strings.NewReader(postData))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	response := httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	// Check that we must have a 422 error code
	suite.Equal(422, code, "Incorrect Error Code")
	// Compare the expected and actual xml response
	suite.Equal(responseJSON, output, "Response body mismatch")
}

// TestCreateReport function implements testing the http POST create report request.
// Request requires admin authentication and gets as input a json body containing
// all the available information to be added to the datastore
// After the operation succeeds is double-checked
// that the newly created report is correctly retrieved
func (suite *ReportTestSuite) TestCreateReport() {

	// create json input data for the request
	postData := `{
    "info": {
        "name": "Foo_Report",
        "description": "olalala"
    },
    "topology_schema": {
        "group": {
            "type": "ngi",
            "group": {
                "type": "site"
            }
        }
    },
	"profiles": [
        {
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
            "type": "metric"
            
        },
		{
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e523",
            "type": "operations"
         
        },
        {
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50bq",
            "type": "aggregation"
            
        },
		{
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e533",
			"type": "weights"
		}
    ],
    "filter_tags": [
        {
            "name": "production",
            "value": "Y"
        },
        {
            "name": "monitored",
            "value": "Y"
        }
    ]
}`

	// Prepare the request object
	request, _ := http.NewRequest("POST", "https://myapi.test.com/api/v2/reports", strings.NewReader(postData))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	response := httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(201, code, "Incorrect Error Code")
	suite.Regexp(suite.respReportCreated, output, "Response body mismatch")

	// Double check that you read the newly inserted profile
	// Create a string literal of the expected xml Response

	responseMessage := respond.ResponseMessage{}
	json.Unmarshal([]byte(output), &responseMessage)
	newUUID := responseMessage.Data.(map[string]interface{})["id"].(string)
	// Prepare the request object using report name as urlvar in url path
	request, _ = http.NewRequest("GET", "/api/v2/reports/"+newUUID, strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	response = httptest.NewRecorder()

	// Pass request to controller calling List() handler method
	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()
	responseJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": \[
  {
   "id": ".*-.*-.*-.*-.*",
   "tenant": "GUARDIANS",
   "disabled": false,
   "info": {
    "name": "Foo_Report",
    "description": "olalala",
    "created": ".*",
    "updated": ".*"
   },
   "thresholds": {
    "availability": 80,
    "reliability": 85,
    "uptime": 80,
    "unknown": 10,
    "downtime": 10
   },
   "topology_schema": {
    "group": {
     "type": "ngi",
     "group": {
      "type": "site"
     }
    }
   },
   "profiles": \[
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
     "name": "profile1",
     "type": "metric"
    },
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e523",
     "name": "profile2",
     "type": "operations"
    },
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50bq",
     "name": "profile3",
     "type": "aggregation"
    },
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e533",
     "name": "Critical",
     "type": "weights"
    }
   \],
   "filter_tags": \[
    {
     "name": "production",
     "value": "Y",
     "context": ""
    },
    {
     "name": "monitored",
     "value": "Y",
     "context": ""
    }
   \]
  }
 \]
}`

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Incorrect Error Code")
	// Compare the expected and actual xml response
	suite.Regexp(responseJSON, output, "Response body mismatch")
}

// TestUpdateReport function implements testing the http PUT update report request.
// Request requires admin authentication and gets as input the name of the
// report to be updated and a json body with the update.
// After the operation succeeds is double-checked
// that the specific report has been updated
func (suite *ReportTestSuite) TestUpdateReport() {

	// create json input data for the request
	postData := `{
    "info": {
        "name": "newname",
        "description": "newdescription",
        "created": "shouldnotchange",
        "updated": "shouldnotchange"
    },
    "topology_schema": {
        "group": {
            "type": "ngi",
            "group": {
                "type": "fight"
            }
        }
	},
	"thresholds": {
		"availability" : 60.33,
		"reliability"  : 40.55
	},
    "profiles": [
        {
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
            "type": "metric",
            "name": "profile1"
        },
		{
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e523",
            "type": "operations",
            "name": "profile2"
        },
        {
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50bq",
            "type": "aggregation",
            "name": "profile3"
        }
    ],
    "filter_tags": [
        {
            "name": "production",
            "value": "N"
        },
        {
            "name": "monitored",
            "value": "N"
        }
    ]
}`
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "/api/v2/reports/eba61a9e-22e9-4521-9e47-ecaa4a494364", strings.NewReader(postData))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")
	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)

	// Execute the request in the controller
	code := response.Code
	output := response.Body.String()

	suite.Equal(200, code, "Incorrect Error Code")
	suite.Equal(suite.respReportUpdated, output, "Response body mismatch")

	// Double check that you read the newly inserted profile
	// Create a string literal of the expected xml Response
	respondJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": \[
  {
   "id": "eba61a9e-22e9-4521-9e47-ecaa4a494364",
   "tenant": "GUARDIANS",
   "disabled": false,
   "info": {
    "name": "newname",
    "description": "newdescription",
    "created": "2015-9-10 13:43:00",
    "updated": ".*"
   },
   "thresholds": {
    "availability": 60.33,
    "reliability": 40.55,
    "uptime": 0,
    "unknown": 0,
    "downtime": 0
   },
   "topology_schema": {
    "group": {
     "type": "ngi",
     "group": {
      "type": "fight"
     }
    }
   },
   "profiles": \[
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
     "name": "profile1",
     "type": "metric"
    },
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e523",
     "name": "profile2",
     "type": "operations"
    },
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50bq",
     "name": "profile3",
     "type": "aggregation"
    }
   \],
   "filter_tags": \[
    {
     "name": "production",
     "value": "N",
     "context": ""
    },
    {
     "name": "monitored",
     "value": "N",
     "context": ""
    }
   \]
  }
 \]
}`

	// Prepare the request object using report name as urlvar in url path
	request, _ = http.NewRequest("GET", "/api/v2/reports/eba61a9e-22e9-4521-9e47-ecaa4a494364", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")
	// Pass request to controller calling List() handler method
	response = httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)

	// Execute the request in the controller
	code = response.Code
	output = response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Incorrect Error Code")
	// Compare the expected and actual xml response
	suite.Regexp(respondJSON, output, "Response body mismatch")
}

func (suite *ReportTestSuite) TestWrongUUIDUpdateReport() {

	// create json input data for the request
	postData := `{
    "info": {
        "name": "newname",
        "description": "newdescription",
        "created": "shouldnotchange",
        "updated": "shouldnotchange"
    },
    "topology_schema": {
        "group": {
            "type": "ngi",
            "group": {
                "type": "fight"
            }
        }
    },
    "profiles": [
        {
			"id": "6ac7d684-1f8e-4a02-a302-720e8f11770b",
            "type": "metric",
            "name": "profile1"
        },
		{
			"id": "6ac7d684-1f8e-4a02-a502-7258f11e523",
            "type": "operations",
            "name": "profile2"
        },
        {
			"id": "6ac7d684-1f8e-4a02-a502-720e8f1250bq",
            "type": "aggregation",
            "name": "profile3"
        }
    ],
    "filter_tags": [
        {
            "name": "production",
            "value": "N"
        },
        {
            "name": "monitored",
            "value": "N"
        }
    ]
}`
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "/api/v2/reports/eba61a9e-22e9-4521-9e47-ecaa4a494364", strings.NewReader(postData))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")
	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)

	// Execute the request in the controller
	code := response.Code
	output := response.Body.String()
	respReportWrongProfileUUIDJSON := `{
 "status": {
  "message": "Unprocessable Entity",
  "code": "422"
 },
 "errors": [
  {
   "message": "Profile id not found",
   "code": "422",
   "details": "No profile in metric_profiles was found with id 6ac7d684-1f8e-4a02-a302-720e8f11770b"
  },
  {
   "message": "Profile id not found",
   "code": "422",
   "details": "No profile in operations_profiles was found with id 6ac7d684-1f8e-4a02-a502-7258f11e523"
  },
  {
   "message": "Profile id not found",
   "code": "422",
   "details": "No profile in aggregation_profiles was found with id 6ac7d684-1f8e-4a02-a502-720e8f1250bq"
  }
 ]
}`

	suite.Equal(422, code, "Incorrect Error Code")
	suite.Equal(respReportWrongProfileUUIDJSON, output, "Response body mismatch")
}

// TestDeleteReport function implements testing the http DELETE report request.
// Request requires admin authentication and gets as input the name of the
// report to be deleted. After the operation succeeds is double-checked
// that the deleted report is actually missing from the datastore
func (suite *ReportTestSuite) TestDeleteReport() {

	// Prepare the request object
	request, _ := http.NewRequest("DELETE", "/api/v2/reports/eba61a9e-22e9-4521-9e47-ecaa4a494364", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	// Execute the request in the controller
	response := httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(200, code, "Incorrect Error Code")
	suite.Equal(suite.respReportDeleted, output, "Response body mismatch")

	// Double check that the report is actually removed when you try
	// to retrieve it's information by name
	// Prepare the request object using report name as urlvar in url path
	request, _ = http.NewRequest("GET", "/api/v2/reports/eba61a9e-22e9-4521-9e47-ecaa4a494364", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")
	// Pass request to controller calling List() handler method
	response = httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(404, code, "Incorrect error code")
	// Compare the expected and actual xml response
	suite.Equal(suite.respReportNotFound, output, "Response body mismatch")
}

// TestReadOneReport function implements the testing
// of the get request which retrieves information
// about a specific report (using it's name as input)
func (suite *ReportTestSuite) TestReadOneReport() {

	// Create a string literal of the expected xml Response
	respondJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "eba61a9e-22e9-4521-9e47-ecaa4a494364",
   "tenant": "GUARDIANS",
   "disabled": false,
   "info": {
    "name": "Report_A",
    "description": "report aaaaa",
    "created": "2015-9-10 13:43:00",
    "updated": "2015-10-11 13:43:00"
   },
   "topology_schema": {
    "group": {
     "type": "NGI",
     "group": {
      "type": "SITE"
     }
    }
   },
   "profiles": [
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
     "name": "profile1",
     "type": "metric"
    },
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e523",
     "name": "profile2",
     "type": "operations"
    },
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50q",
     "name": "profile3",
     "type": "aggregation"
    }
   ],
   "filter_tags": [
    {
     "name": "name1",
     "value": "value1",
     "context": ""
    },
    {
     "name": "name2",
     "value": "value2",
     "context": ""
    }
   ]
  }
 ]
}`

	// Prepare the request object using report name as urlvar in url path
	request, _ := http.NewRequest("GET", "/api/v2/reports/eba61a9e-22e9-4521-9e47-ecaa4a494364", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")
	response := httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Incorrect Error Code")

	// Compare the expected and actual xml response
	suite.Equal(respondJSON, output, "Response body mismatch")

}

// TestReadReport function implements the testing
// of the get request which retrieves information
// about all available reports
func (suite *ReportTestSuite) TestReadReports() {

	// Create a string literal of the expected xml Response
	respondJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "eba61a9e-22e9-4521-9e47-ecaa4a494360",
   "tenant": "GUARDIANS",
   "disabled": false,
   "info": {
    "name": "Report_B",
    "description": "report bbb",
    "created": "2015-10-08 13:43:00",
    "updated": "2015-10-09 13:43:00"
   },
   "topology_schema": {
    "group": {
     "type": "ARCHIPELAGO",
     "group": {
      "type": "ISLAND"
     }
    }
   },
   "profiles": [
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
     "name": "profile1",
     "type": "metric"
    },
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e523",
     "name": "profile2",
     "type": "operations"
    },
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50q",
     "name": "profile3",
     "type": "aggregation"
    }
   ],
   "filter_tags": [
    {
     "name": "name1",
     "value": "value1",
     "context": ""
    },
    {
     "name": "name2",
     "value": "value2",
     "context": ""
    }
   ]
  },
  {
   "id": "eba61a9e-22e9-4521-9e47-ecaa4a494364",
   "tenant": "GUARDIANS",
   "disabled": false,
   "info": {
    "name": "Report_A",
    "description": "report aaaaa",
    "created": "2015-9-10 13:43:00",
    "updated": "2015-10-11 13:43:00"
   },
   "topology_schema": {
    "group": {
     "type": "NGI",
     "group": {
      "type": "SITE"
     }
    }
   },
   "profiles": [
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
     "name": "profile1",
     "type": "metric"
    },
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e523",
     "name": "profile2",
     "type": "operations"
    },
    {
     "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50q",
     "name": "profile3",
     "type": "aggregation"
    }
   ],
   "filter_tags": [
    {
     "name": "name1",
     "value": "value1",
     "context": ""
    },
    {
     "name": "name2",
     "value": "value2",
     "context": ""
    }
   ]
  }
 ]
}`

	// Prepare the request object
	request, _ := http.NewRequest("GET", "/api/v2/reports", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")
	response := httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Incorrect Error Code")
	// Compare the expected and actual xml response
	suite.Equal(respondJSON, string(output), "Response body mismatch")

}

// TestCreateUnauthorized function tests calling the create report request (POST) and
// providing a wrong api-key. The response should be unauthorized
func (suite *ReportTestSuite) TestCreateUnauthorized() {
	// Prepare the request object (use id2 for path)
	request, _ := http.NewRequest("POST", "/api/v2/reports", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "F00T0K3N")

	response := httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(401, code, "Incorrect Error Code")
	suite.Equal(suite.respUnauthorized, output, "Response body mismatch")
}

// TestUpdateUnauthorized function tests calling the update report request (PUT)
// and providing  a wrong api-key. The response should be unauthorized
func (suite *ReportTestSuite) TestUpdateUnauthorized() {
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "/api/v2/reports/eba61a9e-22e9-4521-9e47-ecaa4a494360", strings.NewReader("{}"))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "F00T0K3N")

	response := httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(401, code, "Incorrect Error Code")
	suite.Equal(suite.respUnauthorized, output, "Response body mismatch")
}

// TestDeleteUnauthorized function tests calling the remove report request (DELETE)
// and providing a wrong api-key. The response should be unauthorized
func (suite *ReportTestSuite) TestDeleteUnauthorized() {
	// Prepare the request object
	request, _ := http.NewRequest("DELETE", "/api/v2/reports/eba61a9e-22e9-4521-9e47-ecaa4a494360", strings.NewReader("{}"))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "F00T0K3N")

	response := httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(401, code, "Incorrect Error Code")
	suite.Equal(suite.respUnauthorized, string(output), "Response body mismatch")
}

// TestCreateBadJson tests calling the create report request (POST) and providing
// bad json input. The response should be malformed json
func (suite *ReportTestSuite) TestCreateBadJson() {
	// Prepare the request object
	request, _ := http.NewRequest("POST", "/api/v2/reports", strings.NewReader("{ bad json"))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	response := httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(400, code, "Incorrect Error Code")
	suite.Equal(suite.respBadJSON, output, "Response body mismatch")
}

// TestUpdateBadJson tests calling the update report request (PUT) and providing
// bad json input. The response should be malformed json
func (suite *ReportTestSuite) TestUpdateBadJson() {
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "/api/v2/reports/Re[prt_A", strings.NewReader("{ bad json"))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	response := httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(400, code, "Incorrect Error Code")
	suite.Equal(suite.respBadJSON, output, "Response body mismatch")
}

// TestListOneNotFound tests calling the http (GET) report info request
// and provide a non-existing report name. The response should be report not found
func (suite *ReportTestSuite) TestListOneNotFound() {
	// Prepare the request object
	request, _ := http.NewRequest("GET", "/api/v2/reports/WRONG-UUID-123-124123", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	response := httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	suite.Equal(404, code, "Incorrect Error Code")
	suite.Equal(suite.respReportNotFound, output, "Response body mismatch")
}

// TestUpdateNotFound tests calling the http (PUT) update report equest
// and provide a non-existing report name. The response should be report not found
func (suite *ReportTestSuite) TestUpdateNotFound() {
	// Prepare the request object
	request, _ := http.NewRequest("PUT", "/api/v2/reports/WRONG-UUID-123-124123", strings.NewReader("{}"))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	response := httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(404, code, "Incorrect Error Code")
	suite.Equal(suite.respReportNotFound, output, "Response body mismatch")
}

// TestDeleteNotFound tests calling the http (PUT) update report request
// and provide a non-existing report name. The response should be report not found
func (suite *ReportTestSuite) TestDeleteNotFound() {
	// Prepare the request object
	request, _ := http.NewRequest("DELETE", "/api/v2/reports/WRONG-UUID-123-124123", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	response := httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(404, code, "Incorrect Error Code")
	suite.Equal(suite.respReportNotFound, output, "Response body mismatch")
}

func (suite *ReportTestSuite) TestOptionsReports() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/reports", strings.NewReader(""))
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, POST, DELETE, PUT, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")
}

func (suite *ReportTestSuite) TestCreateForbidViewer() {

	jsonInput := `{
  "name": "test_report",
  }`

	jsonOutput := `{
 "status": {
  "message": "Access to the resource is Forbidden",
  "code": "403"
 }
}`

	request, _ := http.NewRequest("POST", "/api/v2/reports", strings.NewReader(jsonInput))
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

func (suite *ReportTestSuite) TestUpdateForbidViewer() {

	jsonInput := `{}`

	jsonOutput := `{
 "status": {
  "message": "Access to the resource is Forbidden",
  "code": "403"
 }
}`

	request, _ := http.NewRequest("PUT", "/api/v2/reports/6ac7d684-1f8e-4a02-a502-720e8f11e007", strings.NewReader(jsonInput))
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

func (suite *ReportTestSuite) TestDeleteForbidViewer() {

	request, _ := http.NewRequest("DELETE", "/api/v2/reports/6ac7d684-1f8e-4a02-a502-720e8f11e50b", strings.NewReader(""))
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

// TestCreateReportAlreadyExistingName tests the case where the given report has a defined name that already exists
func (suite *ReportTestSuite) TestCreateReportAlreadyExistingName() {

	// create json input data for the request
	postData := `{
    "info": {
        "name": "Report_A",
        "description": "olalala"
    },
    "topology_schema": {
        "group": {
            "type": "ngi",
            "group": {
                "type": "site"
            }
        }
    },
	"profiles": [
        {
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
            "type": "metric",
            "name": "profile1"
        },
		{
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e523",
            "type": "operations",
            "name": "profile2"
        },
        {
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50bq",
            "type": "aggregation",
            "name": "profile3"
        }
    ],
    "filter_tags": [
        {
            "name": "production",
            "value": "Y"
        },
        {
            "name": "monitored",
            "value": "Y"
        }
    ]
}`

	// Prepare the request object
	request, _ := http.NewRequest("POST", "https://myapi.test.com/api/v2/reports", strings.NewReader(postData))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	response := httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(409, code, "Incorrect Error Code")
	suite.Equal(suite.respNameConflict, output, "Response body mismatch")

}

// TestCreateReportAlreadyExistingName tests the case where the given report has a defined name that already exists
func (suite *ReportTestSuite) TestUpdateReportAlreadyExistingName() {

	// create json input data for the request
	postData := `{
    "info": {
        "name": "Report_A",
        "description": "olalala"
    },
    "topology_schema": {
        "group": {
            "type": "ngi",
            "group": {
                "type": "site"
            }
        }
    },
	"profiles": [
        {
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
            "type": "metric",
            "name": "profile1"
        },
		{
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e523",
            "type": "operations",
            "name": "profile2"
        },
        {
			"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50bq",
            "type": "aggregation",
            "name": "profile3"
        }
    ],
    "filter_tags": [
        {
            "name": "production",
            "value": "Y"
        },
        {
            "name": "monitored",
            "value": "Y"
        }
    ]
}`

	// Prepare the request object
	request, _ := http.NewRequest("PUT", "https://myapi.test.com/api/v2/reports/eba61a9e-22e9-4521-9e47-ecaa4a494360", strings.NewReader(postData))
	// add the content-type header to application/json
	request.Header.Set("Accept", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "C4PK3Y")

	response := httptest.NewRecorder()

	// Execute the request in the controller
	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	suite.Equal(409, code, "Incorrect Error Code")
	suite.Equal(suite.respNameConflict, output, "Response body mismatch")

}

// This function is actually called in the end of all tests
// and clears the test environment.
// Mainly it's purpose is to drop the testdb
func (suite *ReportTestSuite) TearDownTest() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}

	tenantDB := session.DB(suite.tenantDbConf.Db)
	mainDB := session.DB(suite.cfg.MongoDB.Db)

	cols, err := tenantDB.CollectionNames()
	for _, col := range cols {
		tenantDB.C(col).RemoveAll(nil)
	}

	cols, err = mainDB.CollectionNames()
	for _, col := range cols {
		mainDB.C(col).RemoveAll(nil)
	}

}

//TearDownTest to tear down every test
func (suite *ReportTestSuite) TearDownSuite() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

// This is the first function called when go test is issued
func TestReportSuite(t *testing.T) {
	suite.Run(t, new(ReportTestSuite))
}
