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

package statusFlatEndpoints

import (
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
type StatusEndpointsTestSuite struct {
	suite.Suite
	cfg          config.Config
	router       *mux.Router
	confHandler  respond.ConfHandler
	tenantDbConf config.MongoConfig
}

// Setup the Test Environment
// This function runs before any test and setups the environment
// A test configuration object is instantiated using a reference
// to testdb: argo_test_details. Also here is are instantiated some expected
// xml response validation messages (authorization,crud responses).
// Also the testdb is seeded with tenants,reports,metric_profiles and status_metrics
func (suite *StatusEndpointsTestSuite) SetupTest() {

	log.SetOutput(ioutil.Discard)

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
    db = "argotest_flatendpoints"
`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(true).PathPrefix("/api/v2/status").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)

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
				"store":    "main",
				"server":   "localhost",
				"port":     27017,
				"database": "argotest_flatendpoints_egi",
				"username": "",
				"password": ""},
		},
		"users": []bson.M{
			bson.M{
				"name":    "egi_user",
				"email":   "egi_user@email.com",
				"roles":   []string{"viewer"},
				"api_key": "KEY1"},
		}})

	c.Insert(bson.M{
		"id": "6ac7d684-1f8e-4a02-a502-720e8f11e50d",
		"info": bson.M{
			"name":    "AVENGERS",
			"email":   "email@something2",
			"website": "www.gotg.com",
			"created": "2015-10-20 02:08:04",
			"updated": "2015-10-20 02:08:04"},
		"db_conf": []bson.M{
			bson.M{
				"store":    "main",
				"server":   "localhost",
				"port":     27017,
				"database": "argotest_flatendpoints_eudat",
				"username": "",
				"password": ""},
		},
		"users": []bson.M{
			bson.M{
				"name":    "eudat_user",
				"email":   "eudat_user@email.com",
				"roles":   []string{"viewer"},
				"api_key": "KEY2"},
		}})
	c = session.DB(suite.cfg.MongoDB.Db).C("roles")
	c.Insert(
		bson.M{
			"resource": "status.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "status.get",
			"roles":    []string{"editor", "viewer"},
		})
	// get dbconfiguration based on the tenant
	// Prepare the request object
	request, _ := http.NewRequest("GET", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY1")
	// authenticate user's api key and find corresponding tenant
	suite.tenantDbConf, _, err = authentication.AuthenticateTenant(request.Header, suite.cfg)

	// Now seed the report DEFINITIONS
	c = session.DB(suite.tenantDbConf.Db).C("reports")
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
					"type": "SITES",
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
	// seed the status detailed metric data
	c = session.DB(suite.tenantDbConf.Db).C("status_endpoints")
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"info":           bson.M{"Url": "http://example.foo/path/to/service"},
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream01.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "OK",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T01:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream01.afroditi.gr",
		"info":           bson.M{"Url": "http://example.foo/path/to/service"},
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "CRITICAL",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T05:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream01.afroditi.gr",
		"info":           bson.M{"Url": "http://example.foo/path/to/service"},
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "OK",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream02.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "OK",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T08:47:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream02.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "WARNING",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T12:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream02.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "OK",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream03.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "OK",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T04:40:00Z",
		"endpoint_group": "HG-03-AUTH",
		"service":        "CREAM-CE",
		"host":           "cream03.afroditi.gr",
		"metric":         "emi.cream.CREAMCE-JobSubmit",
		"status":         "UNKNOWN",
	})
	c.Insert(bson.M{
		"report":             "eba61a9e-22e9-4521-9e47-ecaa4a494364",
		"date_integer":       20150501,
		"timestamp":          "2015-05-01T06:00:00Z",
		"endpoint_group":     "HG-03-AUTH",
		"service":            "CREAM-CE",
		"host":               "cream03.afroditi.gr",
		"metric":             "emi.cream.CREAMCE-JobSubmit",
		"status":             "CRITICAL",
		"has_threshold_rule": true,
	})

	// get dbconfiguration based on the tenant
	// Prepare the request object
	request, _ = http.NewRequest("GET", "", strings.NewReader(""))
	// add the content-type header to application/json
	request.Header.Set("Content-Type", "application/json")
	// add the authentication token which is seeded in testdb
	request.Header.Set("x-api-key", "KEY2")
	// authenticate user's api key and find corresponding tenant
	suite.tenantDbConf, _, err = authentication.AuthenticateTenant(request.Header, suite.cfg)

	// Now seed the reports DEFINITIONS
	c = session.DB(suite.tenantDbConf.Db).C("reports")
	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"info": bson.M{
			"name":        "Report_B",
			"description": "report aaaaa",
			"created":     "2015-9-10 13:43:00",
			"updated":     "2015-10-11 13:43:00",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "EUDAT_GROUPS",
				"group": bson.M{
					"type": "EUDAT_SITES",
				},
			},
		},
		"profiles": []bson.M{
			bson.M{
				"id":   "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
				"type": "metric",
				"name": "eudat.CRITICAL"},
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

	// seed the status detailed metric data
	c = session.DB(suite.tenantDbConf.Db).C("status_endpoints")
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T00:00:00Z",
		"endpoint_group": "EL-01-AUTH",
		"service":        "srv.typeA",
		"host":           "host01.eudat.gr",
		"metric":         "typeA.metric.Memory",
		"status":         "OK",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T01:00:00Z",
		"endpoint_group": "EL-01-AUTH",
		"service":        "srv.typeA",
		"host":           "host01.eudat.gr",
		"metric":         "typeA.metric.Memory",
		"status":         "CRITICAL",
	})
	c.Insert(bson.M{
		"report":         "eba61a9e-22e9-4521-9e47-ecaa4a494365",
		"date_integer":   20150501,
		"timestamp":      "2015-05-01T05:00:00Z",
		"endpoint_group": "EL-01-AUTH",
		"service":        "srv.typeA",
		"host":           "host01.eudat.gr",
		"metric":         "typeA.metric.Memory",
		"status":         "OK",
	})

}

func (suite *StatusEndpointsTestSuite) TestFlatListStatusEndpoints() {

	type expReq struct {
		method string
		url    string
		code   int
		result string
		key    string
	}

	expReqs := []expReq{
		expReq{
			method: "GET",
			url: "/api/v2/status/Report_A/endpoints" +
				"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z",
			code: 200,
			key:  "KEY1",
			result: `{
 "endpoints": [
  {
   "name": "cream01.afroditi.gr",
   "service": "CREAM-CE",
   "supergroup": "HG-03-AUTH",
   "info": {
    "Url": "http://example.foo/path/to/service"
   },
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T01:00:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T05:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  },
  {
   "name": "cream02.afroditi.gr",
   "service": "CREAM-CE",
   "supergroup": "HG-03-AUTH",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T08:47:00Z",
     "value": "WARNING"
    },
    {
     "timestamp": "2015-05-01T12:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  },
  {
   "name": "cream03.afroditi.gr",
   "service": "CREAM-CE",
   "supergroup": "HG-03-AUTH",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T04:40:00Z",
     "value": "UNKNOWN"
    },
    {
     "timestamp": "2015-05-01T06:00:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "CRITICAL"
    }
   ]
  }
 ]
}`},
		expReq{
			method: "GET",
			url: "/api/v2/status/Report_B/endpoints" +
				"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z",
			code: 200,
			key:  "KEY2",
			result: `{
 "endpoints": [
  {
   "name": "host01.eudat.gr",
   "service": "srv.typeA",
   "supergroup": "EL-01-AUTH",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T01:00:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T05:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  }
 ]
}`},

		expReq{
			method: "GET",
			url: "/api/v2/status/Report_A/endpoints" +
				"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z&pageSize=3",
			code: 200,
			key:  "KEY1",
			result: `{
 "endpoints": [
  {
   "name": "cream01.afroditi.gr",
   "service": "CREAM-CE",
   "supergroup": "HG-03-AUTH",
   "info": {
    "Url": "http://example.foo/path/to/service"
   },
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T01:00:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T05:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  }
 ],
 "nextPageToken": "Mw==",
 "pageSize": 3
}`,
		},

		expReq{
			method: "GET",
			url: "/api/v2/status/Report_A/endpoints" +
				"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z&pageSize=3&nextPageToken=Mw==",
			code: 200,
			key:  "KEY1",
			result: `{
 "endpoints": [
  {
   "name": "cream02.afroditi.gr",
   "service": "CREAM-CE",
   "supergroup": "HG-03-AUTH",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T08:47:00Z",
     "value": "WARNING"
    },
    {
     "timestamp": "2015-05-01T12:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  }
 ],
 "nextPageToken": "Ng==",
 "pageSize": 3
}`,
		},

		expReq{
			method: "GET",
			url: "/api/v2/status/Report_A/endpoints" +
				"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z&pageSize=3&nextPageToken=Ng==",
			code: 200,
			key:  "KEY1",
			result: `{
 "endpoints": [
  {
   "name": "cream03.afroditi.gr",
   "service": "CREAM-CE",
   "supergroup": "HG-03-AUTH",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T04:40:00Z",
     "value": "UNKNOWN"
    },
    {
     "timestamp": "2015-05-01T06:00:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "CRITICAL"
    }
   ]
  }
 ],
 "pageSize": 3
}`,
		},

		expReq{
			method: "GET",
			url: "/api/v2/status/Report_A/endpoints" +
				"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z&pageSize=4",
			code: 200,
			key:  "KEY1",
			result: `{
 "endpoints": [
  {
   "name": "cream01.afroditi.gr",
   "service": "CREAM-CE",
   "supergroup": "HG-03-AUTH",
   "info": {
    "Url": "http://example.foo/path/to/service"
   },
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T01:00:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T05:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  },
  {
   "name": "cream02.afroditi.gr",
   "service": "CREAM-CE",
   "supergroup": "HG-03-AUTH",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    }
   ]
  }
 ],
 "nextPageToken": "NA==",
 "pageSize": 4
}`,
		},

		expReq{
			method: "GET",
			url: "/api/v2/status/Report_A/endpoints" +
				"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z&pageSize=4&nextPageToken=NA==",
			code: 200,
			key:  "KEY1",
			result: `{
 "endpoints": [
  {
   "name": "cream02.afroditi.gr",
   "service": "CREAM-CE",
   "supergroup": "HG-03-AUTH",
   "statuses": [
    {
     "timestamp": "2015-05-01T08:47:00Z",
     "value": "WARNING"
    },
    {
     "timestamp": "2015-05-01T12:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  },
  {
   "name": "cream03.afroditi.gr",
   "service": "CREAM-CE",
   "supergroup": "HG-03-AUTH",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T04:40:00Z",
     "value": "UNKNOWN"
    }
   ]
  }
 ],
 "nextPageToken": "OA==",
 "pageSize": 4
}`,
		},

		expReq{
			method: "GET",
			url: "/api/v2/status/Report_A/endpoints" +
				"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z&pageSize=4&nextPageToken=OA==",
			code: 200,
			key:  "KEY1",
			result: `{
 "endpoints": [
  {
   "name": "cream03.afroditi.gr",
   "service": "CREAM-CE",
   "supergroup": "HG-03-AUTH",
   "statuses": [
    {
     "timestamp": "2015-05-01T06:00:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "CRITICAL"
    }
   ]
  }
 ],
 "pageSize": 4
}`,
		},

		expReq{
			method: "GET",
			url: "/api/v2/status/Report_A/endpoints" +
				"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z&pageSize=4&nextPageToken=OA==&view=details",
			code: 200,
			key:  "KEY1",
			result: `{
 "endpoints": [
  {
   "name": "cream03.afroditi.gr",
   "service": "CREAM-CE",
   "supergroup": "HG-03-AUTH",
   "statuses": [
    {
     "timestamp": "2015-05-01T06:00:00Z",
     "value": "CRITICAL",
     "affected_by_threshold_rule": true
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "CRITICAL"
    }
   ]
  }
 ],
 "pageSize": 4
}`,
		},

		expReq{
			method: "GET",
			url: "/api/v2/status/Report_A/endpoints" +
				"?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:00:00Z&pageSize=4&nextPageToken=OA==",
			code: 401,
			key:  "KEYWRONG",
			result: `{
 "status": {
  "message": "Unauthorized",
  "code": "401",
  "details": "You need to provide a correct authentication token using the header 'x-api-key'"
 }
}`,
		},
	}

	for _, expReq := range expReqs {
		request, _ := http.NewRequest(expReq.method, expReq.url, strings.NewReader(""))
		request.Header.Set("x-api-key", expReq.key)
		request.Header.Set("Accept", "application/json")

		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		suite.Equal(expReq.code, response.Code, "Incorrect HTTP response code")
		// Compare the expected and actual xml response
		suite.Equal(expReq.result, response.Body.String(), "Response body mismatch")

	}

}

func (suite *StatusEndpointsTestSuite) TestOptionsStatusEndpoints() {
	request, _ := http.NewRequest("OPTIONS", "/api/v2/status/Report_A/endpoints", strings.NewReader(""))

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	headers := response.HeaderMap

	suite.Equal(200, code, "Error in response code")
	suite.Equal("", output, "Expected empty response body")
	suite.Equal("GET, OPTIONS", headers.Get("Allow"), "Error in Allow header response (supported resource verbs of resource)")
	suite.Equal("text/plain; charset=utf-8", headers.Get("Content-Type"), "Error in Content-Type header response")

}

// This function is actually called in the end of all tests
// and clears the test environment.
// Mainly it's purpose is to drop the testdb
func (suite *StatusEndpointsTestSuite) TearDownTest() {

	session, _ := mongo.OpenSession(suite.cfg.MongoDB)

	session.DB("argotest_flatendpoints").DropDatabase()
	session.DB("argotest_flatendpoints_eudat").DropDatabase()
	session.DB("argotest_flatendpoints_egi").DropDatabase()
}

// This is the first function called when go test is issued
func TestStatusEndpointsSuite(t *testing.T) {
	suite.Run(t, new(StatusEndpointsTestSuite))
}
