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

package topology

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gcfg.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type topologyTestSuite struct {
	suite.Suite
	cfg             config.Config
	router          *mux.Router
	confHandler     respond.ConfHandler
	tenantDbConf    config.MongoConfig
	tenantpassword  string
	tenantusername  string
	tenantstorename string
	clientkey       string
}

// Setup the Test Environment
func (suite *topologyTestSuite) SetupSuite() {

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
	db = "ARGO_test_topology_test"
	`

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

	suite.tenantDbConf.Db = "ARGO_test_topology_tenant"
	suite.tenantDbConf.Password = "h4shp4ss"
	suite.tenantDbConf.Username = "dbuser"
	suite.tenantDbConf.Store = "ar"
	suite.clientkey = "secretkey"

	// Create router and confhandler for test
	suite.confHandler = respond.ConfHandler{Config: suite.cfg}
	suite.router = mux.NewRouter().StrictSlash(false).PathPrefix("/api/v2/topology").Subrouter()
	HandleSubrouter(suite.router, &suite.confHandler)
}

// This function runs before any test and setups the environment
func (suite *topologyTestSuite) SetupTest() {

	log.SetOutput(ioutil.Discard)

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
		bson.M{"name": "Westeros",
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
					"roles":   []string{"viewer"},
				},
				{
					"name":    "King Joffrey",
					"email":   "g0dk1ng@kingslanding.gov",
					"api_key": "sansa <3",
					"roles":   []string{"viewer"},
				},
			}})
	c.Insert(
		bson.M{"name": "EGI",
			"db_conf": []bson.M{

				{
					"server":   "localhost",
					"port":     27017,
					"database": suite.tenantDbConf.Db,
					"username": suite.tenantDbConf.Username,
					"password": suite.tenantDbConf.Password,
				},
				{
					"server":   "localhost",
					"port":     27017,
					"database": "argo_wrong_db_serviceflavoravailability",
				},
			},
			"users": []bson.M{

				{
					"name":    "Joe Complex",
					"email":   "C.Joe@egi.eu",
					"api_key": suite.clientkey,
					"roles":   []string{"viewer"},
				},
				{
					"name":    "Josh Plain",
					"email":   "P.Josh@egi.eu",
					"api_key": "itsamysterytoyou",
					"roles":   []string{"viewer"},
				},
			}})

	c = session.DB(suite.cfg.MongoDB.Db).C("roles")
	c.Insert(
		bson.M{
			"resource": "topology_endpoints_report.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "topology_groups_report.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "topology_groups.delete",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "topology_groups.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "topology_groups.insert",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "topology_endpoints.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "topology_endpoints.insert",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "topology_endpoints.delete",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "topology_service_types.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "topology_service_types.insert",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "topology_service_types.delete",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "topology_stats.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "topology_tags.list",
			"roles":    []string{"editor", "viewer"},
		})
	c.Insert(
		bson.M{
			"resource": "results.get",
			"roles":    []string{"editor", "viewer"},
		})
	// Seed database with recomputations
	c = session.DB(suite.tenantDbConf.Db).C("service_ar")

	// Insert seed data
	c.Insert(
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "SF01",
			"supergroup":   "ST01",
			"up":           0.98264,
			"down":         0,
			"unknown":      0,
			"availability": 98.26389,
			"reliability":  98.26389,
			"tags": []bson.M{
				{
					"name":  "production",
					"value": "Y",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "SF02",
			"supergroup":   "ST01",
			"up":           0.96875,
			"down":         0,
			"unknown":      0,
			"availability": 96.875,
			"reliability":  96.875,
			"tags": []bson.M{
				{
					"name":  "production",
					"value": "Y",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "SF03",
			"supergroup":   "ST02",
			"up":           0.96875,
			"down":         0,
			"unknown":      0,
			"availability": 96.875,
			"reliability":  96.875,
			"tags": []bson.M{
				{
					"name":  "production",
					"value": "Y",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "SF01",
			"supergroup":   "ST01",
			"up":           0.53472,
			"down":         0.33333,
			"unknown":      0.01042,
			"availability": 54.03509,
			"reliability":  81.48148,
			"tags": []bson.M{
				{
					"name":  "production",
					"value": "Y",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "SF02",
			"supergroup":   "ST01",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 100,
			"reliability":  100,
			"tags": []bson.M{
				{
					"name":  "production",
					"value": "Y",
				},
			},
		})

	// Seed database with recomputations
	c = session.DB(suite.tenantDbConf.Db).C("endpoint_group_ar")

	// Insert seed data
	c.Insert(
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "ST01",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 66.7,
			"reliability":  54.6,
			"weight":       5634,
			"tags": []bson.M{
				{
					"name":  "",
					"value": "",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "ST02",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 70,
			"reliability":  45,
			"weight":       4356,
			"tags": []bson.M{
				{
					"name":  "",
					"value": "",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "ST01",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 100,
			"reliability":  100,
			"weight":       5634,
			"tags": []bson.M{
				{
					"name":  "",
					"value": "",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "ST04",
			"supergroup":   "GROUP_B",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 30,
			"reliability":  100,
			"weight":       5344,
			"tags": []bson.M{
				{
					"name":  "",
					"value": "",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "ST05",
			"supergroup":   "GROUP_B",
			"up":           0,
			"down":         0,
			"unknown":      1,
			"availability": 90,
			"reliability":  100,
			"weight":       5634,
			"tags": []bson.M{
				{
					"name":  "",
					"value": "",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150624,
			"name":         "ST05",
			"supergroup":   "GROUP_B",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 40,
			"reliability":  70,
			"weight":       5634,
			"tags": []bson.M{
				{
					"name":  "",
					"value": "",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150625,
			"name":         "ST05",
			"supergroup":   "GROUP_B",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 40,
			"reliability":  70,
			"weight":       5634,
			"tags": []bson.M{
				{
					"name":  "",
					"value": "",
				},
			},
		},
		bson.M{
			"report":       "eba61a9e-22e9-4521-9e47-ecaa4a49436",
			"date":         20150622,
			"name":         "ST02",
			"supergroup":   "GROUP_A",
			"up":           1,
			"down":         0,
			"unknown":      0,
			"availability": 43.5,
			"reliability":  56,
			"weight":       4356,
			"tags": []bson.M{
				{
					"name":  "",
					"value": "",
				},
			},
		})
	c = session.DB(suite.tenantDbConf.Db).C("reports")

	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a49436",
		"info": bson.M{
			"name":        "Critical",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "GROUP",
				"group": bson.M{
					"type": "SITE",
				},
			},
		},
		"profiles": []bson.M{
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"name":  "name1",
				"value": "value1"},
			{
				"name":  "name2",
				"value": "value2"},
		}})

	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a49435",
		"info": bson.M{
			"name":        "Critical2",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "NGIS",
				"group": bson.M{
					"type": "SITE",
				},
			},
		},
		"profiles": []bson.M{
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"name":  "name1",
				"value": "value1"},
			{
				"name":  "name2",
				"value": "value2"},
		}})

	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a49999",
		"info": bson.M{
			"name":        "CriticalExcludeGroup",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "NGIS",
				"group": bson.M{
					"type": "SITE",
				},
			},
		},
		"profiles": []bson.M{
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"name":  "name1",
				"value": "value1"},
			{
				"name":  "name2",
				"value": "value2"},
			{
				"name":    "subgroup",
				"value":   "~SITEB",
				"context": "argo.group.filter.fields",
			},
		}})

	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a4943x",
		"info": bson.M{
			"name":        "Critical3",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "ORG",
				"group": bson.M{
					"type": "SITE",
				},
			},
		},
		"profiles": []bson.M{
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"name":  "name1",
				"value": "value1"},
			{
				"name":  "name2",
				"value": "value2"},
		}})

	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a4943p",
		"info": bson.M{
			"name":        "Critical4",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "PROJECT",
				"group": bson.M{
					"type": "SITE",
				},
			},
		},
		"profiles": []bson.M{
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"name":  "name1",
				"value": "value1"},
			{
				"name":  "name2",
				"value": "value2"},
		}})
	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a4943d",
		"info": bson.M{
			"name":        "Critical5",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "PROJECT",
				"group": bson.M{
					"type": "SITE",
				},
			},
		},
		"profiles": []bson.M{
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"context": "argo.group.filter.tags",
				"name":    "infrastructure",
				"value":   "Devel"},
		}})

	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a4943z",
		"info": bson.M{
			"name":        "Critical6",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "PROJECT",
				"group": bson.M{
					"type": "SITE",
				},
			},
		},
		"profiles": []bson.M{
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"context": "argo.group.filter.tags",
				"name":    "certification",
				"value":   "Certified"},
		}})
	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a4943z",
		"info": bson.M{
			"name":        "Critical7",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "NGIS",
				"group": bson.M{
					"type": "SITES",
				},
			},
		},
		"profiles": []bson.M{
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"context": "argo.group.filter.fields",
				"name":    "group",
				"value":   "NGI0"},
		}})
	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a4943z",
		"info": bson.M{
			"name":        "Critical8",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "NGIS",
				"group": bson.M{
					"type": "SITES",
				},
			},
		},
		"profiles": []bson.M{
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"context": "argo.group.filter.fields",
				"name":    "group",
				"value":   "NGI0"},
			{
				"context": "argo.endpoint.filter.fields",
				"name":    "service",
				"value":   "service_1"},
		}})

	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a4943z",
		"info": bson.M{
			"name":        "Critical9",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "NGIS",
				"group": bson.M{
					"type": "SITES",
				},
			},
		},
		"profiles": []bson.M{
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"context": "argo.group.filter.fields",
				"name":    "group",
				"value":   "NGI0"},
			{
				"context": "argo.endpoint.filter.fields",
				"name":    "service",
				"value":   "service_1"},
			{
				"context": "argo.endpoint.filter.tags",
				"name":    "monitored",
				"value":   "YesNo"},
		}})
	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a4943z",
		"info": bson.M{
			"name":        "CriticalCombine",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "NGIS",
				"group": bson.M{
					"type": "SITES",
				},
			},
		},
		"profiles": []bson.M{
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"context": "argo.group.filter.fields",
				"name":    "group",
				"value":   "NGI0"},
			{
				"context": "argo.group.filter.fields",
				"name":    "group",
				"value":   "NGI0"},
			{
				"context": "argo.endpoint.filter.fields",
				"name":    "hostname",
				"value":   "host1.site_a.foo"},
			{
				"context": "argo.endpoint.filter.fields",
				"name":    "hostname",
				"value":   "host2.site_a.foo"},
		}})

	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a5553z",
		"info": bson.M{
			"name":        "CriticalScope",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "NGIS",
				"group": bson.M{
					"type": "SITES",
				},
			},
		},
		"profiles": []bson.M{
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"context": "argo.endpoint.filter.tags.array",
				"name":    "scope",
				"value":   "tier2, tier1"},
		}})

	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a5553z",
		"info": bson.M{
			"name":        "CriticalGT",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "ORG",
				"group": bson.M{
					"type": "SITES",
				},
			},
		},
		"profiles": []bson.M{
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"context": "whatever",
				"name":    "scope",
				"value":   "tier2, tier1"},
		}})

	c.Insert(bson.M{
		"id": "eba61a9e-22e9-4521-9e47-ecaa4a4956z",
		"info": bson.M{
			"name":        "CriticalExclude",
			"description": "lalalallala",
		},
		"topology_schema": bson.M{
			"group": bson.M{
				"type": "NGIS",
				"group": bson.M{
					"type": "SITES",
				},
			},
		},
		"profiles": []bson.M{
			{
				"type": "metric",
				"name": "ch.cern.SAM.ROC_CRITICAL"},
		},
		"filter_tags": []bson.M{
			{
				"context": "argo.group.filter.fields",
				"name":    "group",
				"value":   "NGI0"},
			{
				"context": "argo.group.filter.fields",
				"name":    "group",
				"value":   "NGI0"},
			{
				"context": "argo.endpoint.filter.fields",
				"name":    "hostname",
				"value":   "host1.site_a.foo"},
			{
				"context": "argo.endpoint.filter.fields",
				"name":    "service",
				"value":   "~service_2"},
			{
				"context": "argo.endpoint.filter.fields",
				"name":    "hostname",
				"value":   "host2.site_a.foo"},
		}})

	// Seed database with endpoint topology
	c = session.DB(suite.tenantDbConf.Db).C(endpointColName)
	c.EnsureIndexKey("-date_integer", "group")
	// Insert seed data

	c.Insert(
		bson.M{
			"date":         "2015-01-11",
			"date_integer": 20150111,
			"group":        "SITEA",
			"type":         "SITES",
			"hostname":     "host1.site_a.foo",
			"service":      "service_1",
			"tags":         bson.M{"production": "1", "monitored": "Yes"},
		},
		bson.M{
			"date":         "2015-01-11",
			"date_integer": 20150111,
			"group":        "SITEA",
			"type":         "SITES",
			"hostname":     "host2.site_a.foo",
			"service":      "service_2",
			"tags":         bson.M{"production": "0", "monitored": "Y"},
		},
		bson.M{
			"date":         "2015-01-11",
			"date_integer": 20150111,
			"group":        "SITEB",
			"type":         "SITES",
			"hostname":     "host1.site_b.foo",
			"service":      "service_1",
			"tags":         bson.M{"production": "Prod", "monitored": "YesNo"},
		},
		bson.M{
			"date":         "2015-01-11",
			"date_integer": 20150111,
			"group":        "SITEC",
			"type":         "SITES",
			"hostname":     "host1.site_c.foo",
			"service":      "service_3",
			"tags":         bson.M{"production": "Prod", "monitored": "No"},
		},
		bson.M{
			"date":         "2015-06-11",
			"date_integer": 20150611,
			"group":        "SITEA",
			"type":         "SITES",
			"hostname":     "host1.site_a.foo",
			"service":      "service_1",
			"tags":         bson.M{"production": "1", "monitored": "Yes"},
		},
		bson.M{
			"date":         "2015-06-11",
			"date_integer": 20150611,
			"group":        "SITEA",
			"type":         "SITES",
			"hostname":     "host2.site_a.foo",
			"service":      "service_2",
			"tags":         bson.M{"production": "0", "monitored": "Y", "scope": "GROUPC, GROUPD, GROUPE"},
		},
		bson.M{
			"date":         "2015-06-11",
			"date_integer": 20150611,
			"group":        "SITEB",
			"type":         "SITES",
			"hostname":     "host1.site_b.foo",
			"service":      "service_1",
			"tags":         bson.M{"production": "Prod", "monitored": "YesNo", "scope": "GROUPA, GROUPB"},
		},
		bson.M{
			"date":         "2015-06-22",
			"date_integer": 20150622,
			"group":        "SITEA",
			"type":         "SITES",
			"hostname":     "host1.site_a.foo",
			"service":      "service_1",
			"tags":         bson.M{"production": "1", "monitored": "1", "scope": "test,tier"},
		},
		bson.M{
			"date":         "2015-06-22",
			"date_integer": 20150622,
			"group":        "SITEA",
			"type":         "SITES",
			"hostname":     "host2.site_a.foo",
			"service":      "service_2",
			"tags":         bson.M{"production": "1", "monitored": "1", "scope": "tier1, lala"},
		},
		bson.M{
			"date":         "2015-06-22",
			"date_integer": 20150622,
			"group":        "SITEB",
			"type":         "SITES",
			"hostname":     "host1.site_b.foo",
			"service":      "service_1",
			"tags":         bson.M{"production": "1", "monitored": "1", "scope": "tier1, tier2, foo"},
		},
		bson.M{
			"date":         "2015-07-22",
			"date_integer": 20150722,
			"group":        "SITEA",
			"type":         "SITES",
			"hostname":     "host1.site_a.foo",
			"service":      "service_1",
			"tags":         bson.M{"production": "0", "monitored": "0"},
		},
		bson.M{
			"date":         "2015-07-22",
			"date_integer": 20150722,
			"group":        "SITEA",
			"type":         "SITES",
			"hostname":     "host2.site_a.foo",
			"service":      "service_2",
			"tags":         bson.M{"production": "0", "monitored": "0"},
		},
		bson.M{
			"date":         "2015-07-22",
			"date_integer": 20150722,
			"group":        "SITEA",
			"type":         "SITES",
			"hostname":     "host3.site_a.foo",
			"service":      "service_3",
			"tags":         bson.M{"production": "0", "monitored": "0", "scope": "TEST"},
		},
		bson.M{
			"date":         "2015-08-10",
			"date_integer": 20150810,
			"group":        "SITEA",
			"type":         "SITES",
			"hostname":     "host0.site_a.foo",
			"service":      "service_x",
			"tags":         bson.M{"production": "0", "monitored": "0"},
		},
		bson.M{
			"date":          "2015-08-10",
			"date_integer":  20150810,
			"group":         "SITEB",
			"type":          "SITES",
			"hostname":      "host0.site_b.foo",
			"service":       "service_x",
			"tags":          bson.M{"production": "0", "monitored": "0"},
			"notifications": bson.M{"contacts": []string{"contact01@email.example.foo", "contact02@email.example.foo"}, "enabled": true},
		},
		bson.M{
			"date":          "2021-01-11",
			"date_integer":  20210111,
			"group":         "SVORG",
			"type":          "SERVICEGROUPS",
			"hostname":      "host0.serv_org.foo",
			"service":       "service_x",
			"tags":          bson.M{"production": "0", "monitored": "0"},
			"notifications": bson.M{"contacts": []string{"contact01@email.example.foo", "contact02@email.example.foo"}, "enabled": true},
		},
		bson.M{
			"date":          "2021-01-11",
			"date_integer":  20210111,
			"group":         "SITEORG",
			"type":          "SITES",
			"hostname":      "host0.site_org.foo",
			"service":       "service_x",
			"tags":          bson.M{"production": "0", "monitored": "0"},
			"notifications": bson.M{"contacts": []string{"contact01@email.example.foo", "contact02@email.example.foo"}, "enabled": true},
		})
	// Seed database with group topology
	c = session.DB(suite.tenantDbConf.Db).C(groupColName)
	c.EnsureIndexKey("-date_integer", "group")
	// Insert seed data
	c.Insert(
		bson.M{
			"date":         "2015-01-11",
			"date_integer": 20150111,
			"group":        "NGI0",
			"type":         "NGIS",
			"subgroup":     "SITEA",
			"tags":         bson.M{"infrastructure": "devtest", "certification": "uncertified"},
		},
		bson.M{
			"date":         "2015-01-11",
			"date_integer": 20150111,
			"group":        "NGI0",
			"type":         "NGIS",
			"subgroup":     "SITEB",
			"tags":         bson.M{"infrastructure": "devel", "certification": "CertNot"},
		},
		bson.M{
			"date":         "2015-01-11",
			"date_integer": 20150111,
			"group":        "NGI1",
			"type":         "NGIS",
			"subgroup":     "SITEC",
			"tags":         bson.M{"infrastructure": "production", "certification": "Certified"},
		},
		bson.M{
			"date":         "2015-06-10",
			"date_integer": 20150610,
			"group":        "NGI0",
			"type":         "NGIS",
			"subgroup":     "SITE_01",
			"tags":         bson.M{"infrastructure": "devtest", "certification": "uncertified"},
		},
		bson.M{
			"date":         "2015-06-10",
			"date_integer": 20150610,
			"group":        "NGI0",
			"type":         "NGIS",
			"subgroup":     "SITE_02",
			"tags":         bson.M{"infrastructure": "devel", "certification": "CertNot"},
		},
		bson.M{
			"date":         "2015-06-10",
			"date_integer": 20150610,
			"group":        "NGI1",
			"type":         "NGIS",
			"subgroup":     "SITE_101",
			"tags":         bson.M{"infrastructure": "production", "certification": "Certified"},
		},
		bson.M{
			"date":         "2015-06-22",
			"date_integer": 20150622,
			"group":        "NGIA",
			"type":         "NGIS",
			"subgroup":     "SITEA",
			"tags":         bson.M{"infrastructure": "Production", "certification": "Certified"},
		},
		bson.M{
			"date":         "2015-06-22",
			"date_integer": 20150622,
			"group":        "NGIA",
			"type":         "NGIS",
			"subgroup":     "SITEB",
			"tags":         bson.M{"infrastructure": "Production", "certification": "Certified"},
		},
		bson.M{
			"date":         "2015-07-22",
			"date_integer": 20150722,
			"group":        "NGIA",
			"type":         "NGIS",
			"subgroup":     "SITEA",
			"tags":         bson.M{"infrastructure": "Production", "certification": "Certified"},
		},
		bson.M{
			"date":         "2015-07-22",
			"date_integer": 20150722,
			"group":        "NGIA",
			"type":         "NGIS",
			"subgroup":     "SITEB",
			"tags":         bson.M{"infrastructure": "Production", "certification": "Certified"},
		},
		bson.M{
			"date":         "2015-07-22",
			"date_integer": 20150722,
			"group":        "NGIX",
			"type":         "NGIS",
			"subgroup":     "SITEX",
			"tags":         bson.M{"infrastructure": "Production", "certification": "Certified"},
		},
		bson.M{
			"date":         "2015-08-11",
			"date_integer": 20150811,
			"group":        "NGIX",
			"type":         "NGIS",
			"subgroup":     "SITEX",
			"tags":         bson.M{"infrastructure": "Production", "certification": "Certified"},
		},
		bson.M{
			"date":         "2020-01-11",
			"date_integer": 20200111,
			"group":        "NGIX",
			"type":         "NGIS",
			"subgroup":     "SITEXYZ",
			"tags":         bson.M{"infrastructure": "Devel", "certification": "Uncertified"},
		},
		bson.M{
			"date":         "2020-01-11",
			"date_integer": 20200111,
			"group":        "NGIX",
			"type":         "NGIS",
			"subgroup":     "SITEXZ",
			"tags":         bson.M{"infrastructure": "Production", "certification": "Certified"},
		},
		bson.M{
			"date":         "2020-01-11",
			"date_integer": 20200111,
			"group":        "NGIX",
			"type":         "NGIS",
			"subgroup":     "SITEX",
			"tags":         bson.M{"infrastructure": "Production", "certification": "Uncertified"},
		},
		bson.M{
			"date":         "2021-01-11",
			"date_integer": 20210111,
			"group":        "ORGB",
			"type":         "ORG",
			"subgroup":     "SITEORG",
			"tags":         bson.M{"infrastructure": "Production", "certification": "Uncertified"},
		},
		bson.M{
			"date":         "2021-01-11",
			"date_integer": 20210111,
			"group":        "ORGB",
			"type":         "ORG",
			"subgroup":     "SVORG",
			"tags":         bson.M{"infrastructure": "Production", "certification": "Uncertified"},
		},
		bson.M{
			"date":          "2020-01-11",
			"date_integer":  20200111,
			"group":         "PR01",
			"type":          "PROJECT",
			"subgroup":      "SITEPROJECT",
			"tags":          bson.M{"infrastructure": "Devel", "certification": "Certified"},
			"notifications": bson.M{"contacts": []string{"contact01@email.example.foo", "contact02@email.example.foo"}, "enabled": true},
		})

	// Seed database with group topology
	c = session.DB(suite.tenantDbConf.Db).C(serviceTypeColName)
	c.EnsureIndexKey("-date_integer", "name")
	// Insert seed data
	c.Insert(
		bson.M{
			"date":         "2015-01-11",
			"date_integer": 20150111,
			"name":         "DB",
			"title":        "Database Service",
			"description":  "A Database type of Service",
		},
		bson.M{
			"date":         "2015-01-11",
			"date_integer": 20150111,
			"name":         "API",
			"title":        "Web API Service",
			"description":  "An API type of Service",
		},
		bson.M{
			"date":         "2015-04-12",
			"date_integer": 20150412,
			"name":         "DB",
			"title":        "Database Service",
			"description":  "A Database type of Service",
		},
		bson.M{
			"date":         "2015-04-12",
			"date_integer": 20150412,
			"name":         "API",
			"title":        "Web API Service",
			"description":  "An API type of Service",
		},
		bson.M{
			"date":         "2015-04-12",
			"date_integer": 20150412,
			"name":         "STORAGE",
			"title":        "Data Storage Service",
			"description":  "A Storage type of Service",
		},
		bson.M{
			"date":         "2015-06-13",
			"date_integer": 20150613,
			"name":         "STORAGE",
			"title":        "Data Storage Service",
			"description":  "A Storage type of Service",
			"tags":         []string{"poem"},
		})

}

func (suite *topologyTestSuite) TestBadDate() {

	badDate := `{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "Bad Request",
   "code": "400",
   "details": "date parameter value: 2020-02 is not in the valid form of YYYY-MM-DD"
  }
 ]
}`

	type reqHeader struct {
		Method string
		Path   string
		Data   string
	}

	requests := []reqHeader{
		{Method: "GET", Path: "/api/v2/topology/endpoints?date=2020-02", Data: ""},
		{Method: "GET", Path: "/api/v2/topology/endpoints/by_report/Critical?date=2020-02", Data: ""},
		{Method: "POST", Path: "/api/v2/topology/endpoints?date=2020-02", Data: ""},
		{Method: "DELETE", Path: "/api/v2/topology/endpoints?date=2020-02", Data: ""},
		{Method: "GET", Path: "/api/v2/topology/groups?date=2020-02", Data: ""},
		{Method: "GET", Path: "/api/v2/topology/groups/by_report/Critical?date=2020-02", Data: ""},
		{Method: "POST", Path: "/api/v2/topology/groups?date=2020-02", Data: ""},
		{Method: "DELETE", Path: "/api/v2/topology/groups?date=2020-02", Data: ""},
		{Method: "GET", Path: "/api/v2/topology/stats/Critical?date=2020-02", Data: ""},
	}

	for _, r := range requests {
		request, _ := http.NewRequest(r.Method, r.Path, strings.NewReader(r.Data))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		code := response.Code
		output := response.Body.String()

		// Check that we must have a 200 ok code
		suite.Equal(400, code, "Internal Server Error")
		// Compare the expected and actual json response
		suite.Equal(badDate, output, "Response body mismatch")

	}

}

func (suite *topologyTestSuite) TestCreateEndpointGroupTopology() {

	expJSON := `{
 "message": "Topology of 3 endpoints created for date: 2019-03-03",
 "code": "201"
}`

	expJSON2 := `{
 "message": "Topology already exists for date: 2019-03-03, please either update it or delete it first!",
 "code": "409"
}`

	jsonInput := `	[
	 {"group": "SITE_A", "hostname": "host1.site-a.foo", "type": "SITES", "service": "a.service.foo", "tags": {"scope": "TENANT", "production": "1", "monitored": "1"}},
	 {"group": "SITE_A", "hostname": "host2.site-b.foo", "type": "SITES", "service": "b.service.foo", "tags": {"scope": "TENANT", "production": "1", "monitored": "1"}},
	 {"group": "SITE_B", "hostname": "host1.site-a.foo", "type": "SITES", "service": "c.service.foo", "tags": {"scope": "TENANT", "production": "1", "monitored": "1"}}
	]`
	request, _ := http.NewRequest("POST", "/api/v2/topology/endpoints?date=2019-03-03", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(201, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(expJSON, output, "Creation failed")

	// Now test inserting again it should create a conflict

	request2, _ := http.NewRequest("POST", "/api/v2/topology/endpoints?date=2019-03-03", strings.NewReader(jsonInput))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)
	code2 := response2.Code
	output2 := response2.Body.String()

	// Check that we must have a 409 conflict code
	suite.Equal(409, code2, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(expJSON2, output2, "Creation failed")

}

func (suite *topologyTestSuite) TestCreateGroupTopology() {

	expJSON := `{
 "message": "Topology of 3 groups created for date: 2019-03-03",
 "code": "201"
}`

	expJSON2 := `{
 "message": "Topology already exists for date: 2019-03-03, please either update it or delete it first!",
 "code": "409"
}`

	jsonInput := `	[
	 {"group": "NGIA", "type": "NGIS", "service": "SITEA", "tags": {"scope": "FEDERATION", "infrastructure": "Production", "certification": "Certified"}},
	 {"group": "NGIA", "type": "NGIS", "service": "SITEB", "tags": {"scope": "FEDERATION", "infrastructure": "Production", "certification": "Certified"}},
	 {"group": "NGIZ", "type": "NGIS", "service": "SITEZ", "tags": {"scope": "FEDERATION", "infrastructure": "Production", "certification": "Certified"}}
	]`
	request, _ := http.NewRequest("POST", "/api/v2/topology/groups?date=2019-03-03", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(201, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(expJSON, output, "Creation failed")

	// Now test inserting again it should create a conflict

	request2, _ := http.NewRequest("POST", "/api/v2/topology/groups?date=2019-03-03", strings.NewReader(jsonInput))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)
	code2 := response2.Code
	output2 := response2.Body.String()

	// Check that we must have a 409 conflict code
	suite.Equal(409, code2, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(expJSON2, output2, "Creation failed")

}

func (suite *topologyTestSuite) TestCreateServiceTypeTopology() {

	expJSON := `{
 "message": "Topology of 3 service types created for date: 2019-03-03",
 "code": "201"
}`

	expJSON2 := `{
 "message": "Topology list of service types already exists for date: 2019-03-03, please either update it or delete it first!",
 "code": "409"
}`

	expJSON3 := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2019-03-03",
   "name": "Service_A",
   "title": "Service A",
   "description": "a short descritpion of service type a"
  },
  {
   "date": "2019-03-03",
   "name": "Service_B",
   "title": "Service B",
   "description": "a short descritpion of service type b"
  },
  {
   "date": "2019-03-03",
   "name": "Service_C",
   "title": "Service C",
   "description": "a short descritpion of service type c"
  }
 ]
}`

	jsonInput := `[
  {
    "name": "Service_A",
	"title": "Service A",
    "description": "a short descritpion of service type a"
  },
  {
    "name": "Service_B",
	"title": "Service B",
    "description": "a short descritpion of service type b"
  },
  {
    "name": "Service_C",
	"title": "Service C",
    "description": "a short descritpion of service type c"
  }
]`
	request, _ := http.NewRequest("POST", "/api/v2/topology/service-types?date=2019-03-03", strings.NewReader(jsonInput))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()
	// Check that we must have a 200 ok code
	suite.Equal(201, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(expJSON, output, "Creation failed")

	// Now test inserting again it should create a conflict

	request2, _ := http.NewRequest("POST", "/api/v2/topology/service-types?date=2019-03-03", strings.NewReader(jsonInput))
	request2.Header.Set("x-api-key", suite.clientkey)
	request2.Header.Set("Accept", "application/json")
	response2 := httptest.NewRecorder()

	suite.router.ServeHTTP(response2, request2)
	code2 := response2.Code
	output2 := response2.Body.String()

	// Check that we must have a 409 conflict code
	suite.Equal(409, code2, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(expJSON2, output2, "Creation failed")

	// Now test retrieving the documents

	request3, _ := http.NewRequest("GET", "/api/v2/topology/service-types?date=2019-03-03", strings.NewReader(jsonInput))
	request3.Header.Set("x-api-key", suite.clientkey)
	request3.Header.Set("Accept", "application/json")
	response3 := httptest.NewRecorder()

	suite.router.ServeHTTP(response3, request3)
	code3 := response3.Code
	output3 := response3.Body.String()

	// Check that we must have a 409 conflict code
	suite.Equal(200, code3, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(expJSON3, output3, "Creation failed")

}

func (suite *topologyTestSuite) TestListFilterEndpointTags() {

	type TestReq struct {
		Path string
		Code int
		JSON string
	}

	expected := []TestReq{
		{
			Path: "/api/v2/topology/endpoints?date=2015-06-12&tags=monitored:Y*",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-11",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_1",
   "hostname": "host1.site_a.foo",
   "tags": {
    "monitored": "Yes",
    "production": "1"
   }
  },
  {
   "date": "2015-06-11",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_2",
   "hostname": "host2.site_a.foo",
   "tags": {
    "monitored": "Y",
    "production": "0",
    "scope": "GROUPC, GROUPD, GROUPE"
   }
  },
  {
   "date": "2015-06-11",
   "group": "SITEB",
   "type": "SITES",
   "service": "service_1",
   "hostname": "host1.site_b.foo",
   "tags": {
    "monitored": "YesNo",
    "production": "Prod",
    "scope": "GROUPA, GROUPB"
   }
  }
 ]
}`},
		{
			Path: "/api/v2/topology/endpoints?date=2015-06-12&tags=production:1*",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-11",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_1",
   "hostname": "host1.site_a.foo",
   "tags": {
    "monitored": "Yes",
    "production": "1"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/endpoints?date=2015-06-12&tags=monitored:Yes,monitored:Y",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-11",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_1",
   "hostname": "host1.site_a.foo",
   "tags": {
    "monitored": "Yes",
    "production": "1"
   }
  },
  {
   "date": "2015-06-11",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_2",
   "hostname": "host2.site_a.foo",
   "tags": {
    "monitored": "Y",
    "production": "0",
    "scope": "GROUPC, GROUPD, GROUPE"
   }
  }
 ]
}`},
	}

	for _, exp := range expected {
		request, _ := http.NewRequest("GET", exp.Path, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		code := response.Code
		output := response.Body.String()

		// Check that we must have a 200 ok code
		suite.Equal(exp.Code, code, "Response Code Mismatch on call:"+exp.Path)
		// Compare the expected and actual json response
		suite.Equal(exp.JSON, output, "Response body mismatch on call:"+exp.Path)
	}
}
func (suite *topologyTestSuite) TestListFilterGroupTags() {

	type TestReq struct {
		Path string
		Code int
		JSON string
	}

	expected := []TestReq{
		{
			Path: "/api/v2/topology/groups?&date=2015-06-12&tags=infrastructure:production",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-10",
   "group": "NGI1",
   "type": "NGIS",
   "subgroup": "SITE_101",
   "tags": {
    "certification": "Certified",
    "infrastructure": "production"
   }
  }
 ]
}`},
		{
			Path: "/api/v2/topology/groups?&date=2015-06-12&tags=infrastructure:dev*",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-10",
   "group": "NGI0",
   "type": "NGIS",
   "subgroup": "SITE_01",
   "tags": {
    "certification": "uncertified",
    "infrastructure": "devtest"
   }
  },
  {
   "date": "2015-06-10",
   "group": "NGI0",
   "type": "NGIS",
   "subgroup": "SITE_02",
   "tags": {
    "certification": "CertNot",
    "infrastructure": "devel"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/groups?&date=2015-06-12&tags=infrastructure:*test",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-10",
   "group": "NGI0",
   "type": "NGIS",
   "subgroup": "SITE_01",
   "tags": {
    "certification": "uncertified",
    "infrastructure": "devtest"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/groups?&date=2015-06-12&tags=certification:Cert*",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-10",
   "group": "NGI0",
   "type": "NGIS",
   "subgroup": "SITE_02",
   "tags": {
    "certification": "CertNot",
    "infrastructure": "devel"
   }
  },
  {
   "date": "2015-06-10",
   "group": "NGI1",
   "type": "NGIS",
   "subgroup": "SITE_101",
   "tags": {
    "certification": "Certified",
    "infrastructure": "production"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/groups?&date=2015-06-12&tags=infrastructure:devel,infrastructure:production",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-10",
   "group": "NGI0",
   "type": "NGIS",
   "subgroup": "SITE_02",
   "tags": {
    "certification": "CertNot",
    "infrastructure": "devel"
   }
  },
  {
   "date": "2015-06-10",
   "group": "NGI1",
   "type": "NGIS",
   "subgroup": "SITE_101",
   "tags": {
    "certification": "Certified",
    "infrastructure": "production"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/groups?&date=2015-06-12&tags=infrastructure:~production",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-10",
   "group": "NGI0",
   "type": "NGIS",
   "subgroup": "SITE_01",
   "tags": {
    "certification": "uncertified",
    "infrastructure": "devtest"
   }
  },
  {
   "date": "2015-06-10",
   "group": "NGI0",
   "type": "NGIS",
   "subgroup": "SITE_02",
   "tags": {
    "certification": "CertNot",
    "infrastructure": "devel"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/groups?&date=2015-06-12&subgroup=~SITE_02",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-10",
   "group": "NGI0",
   "type": "NGIS",
   "subgroup": "SITE_01",
   "tags": {
    "certification": "uncertified",
    "infrastructure": "devtest"
   }
  },
  {
   "date": "2015-06-10",
   "group": "NGI1",
   "type": "NGIS",
   "subgroup": "SITE_101",
   "tags": {
    "certification": "Certified",
    "infrastructure": "production"
   }
  }
 ]
}`},
	}

	for _, exp := range expected {
		request, _ := http.NewRequest("GET", exp.Path, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		code := response.Code
		output := response.Body.String()

		// Check that we must have a 200 ok code
		suite.Equal(exp.Code, code, "Response Code Mismatch on call:"+exp.Path)
		// Compare the expected and actual json response
		suite.Equal(exp.JSON, output, "Response body mismatch on call:"+exp.Path)

	}
}

func (suite *topologyTestSuite) TestListTopologyTags() {

	type TestReq struct {
		Path string
		Code int
		JSON string
	}

	expected := []TestReq{
		{
			Path: "/api/v2/topology/tags?&date=2015-06-12",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "name": "endpoints",
   "values": [
    {
     "name": "monitored",
     "values": [
      "Y",
      "Yes",
      "YesNo"
     ]
    },
    {
     "name": "production",
     "values": [
      "0",
      "1",
      "Prod"
     ]
    },
    {
     "name": "scope",
     "values": [
      "GROUPA",
      "GROUPB",
      "GROUPC",
      "GROUPD",
      "GROUPE"
     ]
    }
   ]
  },
  {
   "name": "groups",
   "values": [
    {
     "name": "certification",
     "values": [
      "CertNot",
      "Certified",
      "uncertified"
     ]
    },
    {
     "name": "infrastructure",
     "values": [
      "devel",
      "devtest",
      "production"
     ]
    }
   ]
  }
 ]
}`},
	}

	for _, exp := range expected {
		request, _ := http.NewRequest("GET", exp.Path, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		code := response.Code
		output := response.Body.String()

		// Check that we must have a 200 ok code
		suite.Equal(exp.Code, code, "Response Code Mismatch on call:"+exp.Path)
		// Compare the expected and actual json response

		suite.Equal(exp.JSON, output, "Response body mismatch on call:"+exp.Path)

	}
}

func (suite *topologyTestSuite) TestListFilterGroupsByReport() {

	type TestReq struct {
		Path string
		Code int
		JSON string
	}

	expected := []TestReq{
		{
			Path: "/api/v2/topology/groups/by_report/Critical2?date=2020-02-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2020-01-11",
   "group": "NGIX",
   "type": "NGIS",
   "subgroup": "SITEXYZ",
   "tags": {
    "certification": "Uncertified",
    "infrastructure": "Devel"
   }
  },
  {
   "date": "2020-01-11",
   "group": "NGIX",
   "type": "NGIS",
   "subgroup": "SITEXZ",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  },
  {
   "date": "2020-01-11",
   "group": "NGIX",
   "type": "NGIS",
   "subgroup": "SITEX",
   "tags": {
    "certification": "Uncertified",
    "infrastructure": "Production"
   }
  }
 ]
}`},
		{
			Path: "/api/v2/topology/groups/by_report/Critical3?date=2021-02-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2021-01-11",
   "group": "ORGB",
   "type": "ORG",
   "subgroup": "SITEORG",
   "tags": {
    "certification": "Uncertified",
    "infrastructure": "Production"
   }
  },
  {
   "date": "2021-01-11",
   "group": "ORGB",
   "type": "ORG",
   "subgroup": "SVORG",
   "tags": {
    "certification": "Uncertified",
    "infrastructure": "Production"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/groups/by_report/Critical4?date=2020-02-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2020-01-11",
   "group": "PR01",
   "type": "PROJECT",
   "subgroup": "SITEPROJECT",
   "notifications": {
    "contacts": [
     "contact01@email.example.foo",
     "contact02@email.example.foo"
    ],
    "enabled": true
   },
   "tags": {
    "certification": "Certified",
    "infrastructure": "Devel"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/groups/by_report/Critical5?date=2020-02-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2020-01-11",
   "group": "PR01",
   "type": "PROJECT",
   "subgroup": "SITEPROJECT",
   "notifications": {
    "contacts": [
     "contact01@email.example.foo",
     "contact02@email.example.foo"
    ],
    "enabled": true
   },
   "tags": {
    "certification": "Certified",
    "infrastructure": "Devel"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/groups/by_report/Critical6?date=2020-02-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2020-01-11",
   "group": "PR01",
   "type": "PROJECT",
   "subgroup": "SITEPROJECT",
   "notifications": {
    "contacts": [
     "contact01@email.example.foo",
     "contact02@email.example.foo"
    ],
    "enabled": true
   },
   "tags": {
    "certification": "Certified",
    "infrastructure": "Devel"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/groups/by_report/Critical7?date=2015-01-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-01-11",
   "group": "NGI0",
   "type": "NGIS",
   "subgroup": "SITEA",
   "tags": {
    "certification": "uncertified",
    "infrastructure": "devtest"
   }
  },
  {
   "date": "2015-01-11",
   "group": "NGI0",
   "type": "NGIS",
   "subgroup": "SITEB",
   "tags": {
    "certification": "CertNot",
    "infrastructure": "devel"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/groups/by_report/CriticalCombine?date=2015-01-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-01-11",
   "group": "NGI0",
   "type": "NGIS",
   "subgroup": "SITEA",
   "tags": {
    "certification": "uncertified",
    "infrastructure": "devtest"
   }
  },
  {
   "date": "2015-01-11",
   "group": "NGI0",
   "type": "NGIS",
   "subgroup": "SITEB",
   "tags": {
    "certification": "CertNot",
    "infrastructure": "devel"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/groups/by_report/CriticalScope?date=2015-06-22",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-22",
   "group": "NGIA",
   "type": "NGIS",
   "subgroup": "SITEA",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  },
  {
   "date": "2015-06-22",
   "group": "NGIA",
   "type": "NGIS",
   "subgroup": "SITEB",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/groups/by_report/CriticalExcludeGroup?date=2015-06-22",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-22",
   "group": "NGIA",
   "type": "NGIS",
   "subgroup": "SITEA",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  }
 ]
}`},
	}

	for _, exp := range expected {
		request, _ := http.NewRequest("GET", exp.Path, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		code := response.Code
		output := response.Body.String()

		// Check that we must have a 200 ok code
		suite.Equal(exp.Code, code, "Response Code Mismatch on call:"+exp.Path)
		// Compare the expected and actual json response
		suite.Equal(exp.JSON, output, "Response body mismatch on call:"+exp.Path)

	}
}

func (suite *topologyTestSuite) TestListFilterEndpointsByReport() {

	type TestReq struct {
		Path string
		Code int
		JSON string
	}

	expected := []TestReq{
		{
			Path: "/api/v2/topology/endpoints/by_report/Critical7?date=2015-01-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-01-11",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_1",
   "hostname": "host1.site_a.foo",
   "tags": {
    "monitored": "Yes",
    "production": "1"
   }
  },
  {
   "date": "2015-01-11",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_2",
   "hostname": "host2.site_a.foo",
   "tags": {
    "monitored": "Y",
    "production": "0"
   }
  },
  {
   "date": "2015-01-11",
   "group": "SITEB",
   "type": "SITES",
   "service": "service_1",
   "hostname": "host1.site_b.foo",
   "tags": {
    "monitored": "YesNo",
    "production": "Prod"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/endpoints/by_report/Critical8?date=2015-01-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-01-11",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_1",
   "hostname": "host1.site_a.foo",
   "tags": {
    "monitored": "Yes",
    "production": "1"
   }
  },
  {
   "date": "2015-01-11",
   "group": "SITEB",
   "type": "SITES",
   "service": "service_1",
   "hostname": "host1.site_b.foo",
   "tags": {
    "monitored": "YesNo",
    "production": "Prod"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/endpoints/by_report/Critical9?date=2015-01-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-01-11",
   "group": "SITEB",
   "type": "SITES",
   "service": "service_1",
   "hostname": "host1.site_b.foo",
   "tags": {
    "monitored": "YesNo",
    "production": "Prod"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/endpoints/by_report/CriticalScope?date=2015-06-22",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-22",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_2",
   "hostname": "host2.site_a.foo",
   "tags": {
    "monitored": "1",
    "production": "1",
    "scope": "tier1, lala"
   }
  },
  {
   "date": "2015-06-22",
   "group": "SITEB",
   "type": "SITES",
   "service": "service_1",
   "hostname": "host1.site_b.foo",
   "tags": {
    "monitored": "1",
    "production": "1",
    "scope": "tier1, tier2, foo"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/endpoints/by_report/CriticalGT?date=2021-01-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2021-01-11",
   "group": "SITEORG",
   "type": "SITES",
   "service": "service_x",
   "hostname": "host0.site_org.foo",
   "notifications": {
    "contacts": [
     "contact01@email.example.foo",
     "contact02@email.example.foo"
    ],
    "enabled": true
   },
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/endpoints/by_report/CriticalCombine?date=2015-01-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-01-11",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_1",
   "hostname": "host1.site_a.foo",
   "tags": {
    "monitored": "Yes",
    "production": "1"
   }
  },
  {
   "date": "2015-01-11",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_2",
   "hostname": "host2.site_a.foo",
   "tags": {
    "monitored": "Y",
    "production": "0"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/endpoints/by_report/CriticalExclude?date=2015-01-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-01-11",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_1",
   "hostname": "host1.site_a.foo",
   "tags": {
    "monitored": "Yes",
    "production": "1"
   }
  }
 ]
}`},
	}

	for _, exp := range expected {
		request, _ := http.NewRequest("GET", exp.Path, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		code := response.Code
		output := response.Body.String()

		// Check that we must have a 200 ok code
		suite.Equal(exp.Code, code, "Response Code Mismatch on call:"+exp.Path)
		// Compare the expected and actual json response
		suite.Equal(exp.JSON, output, "Response body mismatch on call:"+exp.Path)
	}
}

func (suite *topologyTestSuite) TestListFilterGroups() {

	type TestReq struct {
		Path string
		Code int
		JSON string
	}

	expected := []TestReq{
		{
			Path: "/api/v2/topology/groups?date=2015-06-30&group=NGIA",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-22",
   "group": "NGIA",
   "type": "NGIS",
   "subgroup": "SITEA",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  },
  {
   "date": "2015-06-22",
   "group": "NGIA",
   "type": "NGIS",
   "subgroup": "SITEB",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  }
 ]
}`},
		{
			Path: "/api/v2/topology/groups?date=2015-06-30&group=NGIA&subgroup=SITEB",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-22",
   "group": "NGIA",
   "type": "NGIS",
   "subgroup": "SITEB",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  }
 ]
}`},
		{
			Path: "/api/v2/topology/groups?date=2015-06-30&subgroup=*A",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-22",
   "group": "NGIA",
   "type": "NGIS",
   "subgroup": "SITEA",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/groups?date=2015-06-30&subgroup=~*A",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-22",
   "group": "NGIA",
   "type": "NGIS",
   "subgroup": "SITEB",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  }
 ]
}`},

		{
			Path: "/api/v2/topology/groups?date=2015-06-30&subgroup=A*",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": []
}`},

		{
			Path: "/api/v2/topology/groups?date=2015-06-30&subgroup=*A&subgroup=*B",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-22",
   "group": "NGIA",
   "type": "NGIS",
   "subgroup": "SITEA",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  },
  {
   "date": "2015-06-22",
   "group": "NGIA",
   "type": "NGIS",
   "subgroup": "SITEB",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  }
 ]
}`},
	}

	for _, exp := range expected {
		request, _ := http.NewRequest("GET", exp.Path, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		code := response.Code
		output := response.Body.String()

		// Check that we must have a 200 ok code
		suite.Equal(exp.Code, code, "Response Code Mismatch on call:"+exp.Path)
		// Compare the expected and actual json response
		suite.Equal(exp.JSON, output, "Response body mismatch on call:"+exp.Path)

	}
}

func (suite *topologyTestSuite) TestListServiceTypes() {

	type TestReq struct {
		Path string
		Code int
		JSON string
	}

	expected := []TestReq{
		{
			Path: "/api/v2/topology/service-types?date=2015-01-30",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-01-11",
   "name": "API",
   "title": "Web API Service",
   "description": "An API type of Service"
  },
  {
   "date": "2015-01-11",
   "name": "DB",
   "title": "Database Service",
   "description": "A Database type of Service"
  }
 ]
}`},

		{
			Path: "/api/v2/topology/service-types?date=2015-04-28",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-04-12",
   "name": "API",
   "title": "Web API Service",
   "description": "An API type of Service"
  },
  {
   "date": "2015-04-12",
   "name": "DB",
   "title": "Database Service",
   "description": "A Database type of Service"
  },
  {
   "date": "2015-04-12",
   "name": "STORAGE",
   "title": "Data Storage Service",
   "description": "A Storage type of Service"
  }
 ]
}`},

		{
			Path: "/api/v2/topology/service-types?date=2015-06-30",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-13",
   "name": "STORAGE",
   "title": "Data Storage Service",
   "description": "A Storage type of Service",
   "tags": [
    "poem"
   ]
  }
 ]
}`},
	}

	for _, exp := range expected {
		request, _ := http.NewRequest("GET", exp.Path, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		code := response.Code
		output := response.Body.String()

		// Check that we must have a 200 ok code
		suite.Equal(exp.Code, code, "Response Code Mismatch on call:"+exp.Path)
		// Compare the expected and actual json response
		if !(suite.Equal(exp.JSON, output, "Response body mismatch on call:"+exp.Path)) {
			fmt.Println(output)
		}
	}
}

func (suite *topologyTestSuite) TestListFilterEndpoints() {

	type TestReq struct {
		Path string
		Code int
		JSON string
	}

	expected := []TestReq{
		{
			Path: "/api/v2/topology/endpoints?group=SITEA&date=2020-02-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-08-10",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_x",
   "hostname": "host0.site_a.foo",
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  }
 ]
}`},
		{
			Path: "/api/v2/topology/endpoints?group=B*&date=2020-02-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": []
}`},
		{
			Path: "/api/v2/topology/endpoints?service=serv*&date=2020-02-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-08-10",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_x",
   "hostname": "host0.site_a.foo",
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  },
  {
   "date": "2015-08-10",
   "group": "SITEB",
   "type": "SITES",
   "service": "service_x",
   "hostname": "host0.site_b.foo",
   "notifications": {
    "contacts": [
     "contact01@email.example.foo",
     "contact02@email.example.foo"
    ],
    "enabled": true
   },
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  }
 ]
}`},
		{
			Path: "/api/v2/topology/endpoints?service=serv*&hostname=*site_b*&date=2020-02-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-08-10",
   "group": "SITEB",
   "type": "SITES",
   "service": "service_x",
   "hostname": "host0.site_b.foo",
   "notifications": {
    "contacts": [
     "contact01@email.example.foo",
     "contact02@email.example.foo"
    ],
    "enabled": true
   },
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  }
 ]
}`},
		{
			Path: "/api/v2/topology/endpoints?service=serv*&hostname=*.foo*&date=2020-02-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-08-10",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_x",
   "hostname": "host0.site_a.foo",
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  },
  {
   "date": "2015-08-10",
   "group": "SITEB",
   "type": "SITES",
   "service": "service_x",
   "hostname": "host0.site_b.foo",
   "notifications": {
    "contacts": [
     "contact01@email.example.foo",
     "contact02@email.example.foo"
    ],
    "enabled": true
   },
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  }
 ]
}`},
		{
			Path: "/api/v2/topology/endpoints?service=serv*&hostname=*.foo*&group=*B&date=2020-02-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-08-10",
   "group": "SITEB",
   "type": "SITES",
   "service": "service_x",
   "hostname": "host0.site_b.foo",
   "notifications": {
    "contacts": [
     "contact01@email.example.foo",
     "contact02@email.example.foo"
    ],
    "enabled": true
   },
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  }
 ]
}`},
		{
			Path: "/api/v2/topology/endpoints?service=serv*&hostname=*.foo*&group=*B&type=SITES&date=2020-02-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-08-10",
   "group": "SITEB",
   "type": "SITES",
   "service": "service_x",
   "hostname": "host0.site_b.foo",
   "notifications": {
    "contacts": [
     "contact01@email.example.foo",
     "contact02@email.example.foo"
    ],
    "enabled": true
   },
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  }
 ]
}`},
		{
			Path: "/api/v2/topology/endpoints?service=serv*&hostname=*.foo*&group=*B&type=BITES&date=2020-02-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": []
}`},

		{
			Path: "/api/v2/topology/endpoints?group=*A&group=*B&date=2020-02-11",
			Code: 200,
			JSON: `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-08-10",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_x",
   "hostname": "host0.site_a.foo",
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  },
  {
   "date": "2015-08-10",
   "group": "SITEB",
   "type": "SITES",
   "service": "service_x",
   "hostname": "host0.site_b.foo",
   "notifications": {
    "contacts": [
     "contact01@email.example.foo",
     "contact02@email.example.foo"
    ],
    "enabled": true
   },
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  }
 ]
}`},
	}

	for _, exp := range expected {
		request, _ := http.NewRequest("GET", exp.Path, strings.NewReader(""))
		request.Header.Set("x-api-key", suite.clientkey)
		request.Header.Set("Accept", "application/json")
		response := httptest.NewRecorder()

		suite.router.ServeHTTP(response, request)

		code := response.Code
		output := response.Body.String()

		// Check that we must have a 200 ok code
		suite.Equal(exp.Code, code, "Response Code Mismatch on call:"+exp.Path)
		// Compare the expected and actual json response
		suite.Equal(exp.JSON, output, "Response body mismatch on call:"+exp.Path)

	}
}

func (suite *topologyTestSuite) TestListEndpoints() {

	request, _ := http.NewRequest("GET", "/api/v2/topology/endpoints", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	profileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2021-01-11",
   "group": "SITEORG",
   "type": "SITES",
   "service": "service_x",
   "hostname": "host0.site_org.foo",
   "notifications": {
    "contacts": [
     "contact01@email.example.foo",
     "contact02@email.example.foo"
    ],
    "enabled": true
   },
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  },
  {
   "date": "2021-01-11",
   "group": "SVORG",
   "type": "SERVICEGROUPS",
   "service": "service_x",
   "hostname": "host0.serv_org.foo",
   "notifications": {
    "contacts": [
     "contact01@email.example.foo",
     "contact02@email.example.foo"
    ],
    "enabled": true
   },
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(profileJSON, output, "Response body mismatch")

}

func (suite *topologyTestSuite) TestDeleteEndpoints() {

	request, _ := http.NewRequest("DELETE", "/api/v2/topology/endpoints?date=2015-08-10", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	output1JSON := `{
 "message": "Topology of 2 endpoints deleted for date: 2015-08-10",
 "code": "200"
}`

	output2JSON := `{
 "status": {
  "message": "Not Found",
  "code": "404"
 },
 "errors": [
  {
   "message": "Not Found",
   "code": "404",
   "details": "Specific query returned no items"
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(output1JSON, output, "Response body mismatch")

	request, _ = http.NewRequest("DELETE", "/api/v2/topology/endpoints?date=2015-08-10", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()

	// Check that we must have a 404 not found
	suite.Equal(404, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(output2JSON, output, "Response body mismatch")

}

func (suite *topologyTestSuite) TestDeleteServiceTypes() {

	request, _ := http.NewRequest("DELETE", "/api/v2/topology/service-types?date=2015-04-12", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	output1JSON := `{
 "message": "Topology of 3 service types deleted for date: 2015-04-12",
 "code": "200"
}`

	output2JSON := `{
 "status": {
  "message": "Not Found",
  "code": "404"
 },
 "errors": [
  {
   "message": "Not Found",
   "code": "404",
   "details": "Specific query returned no items"
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(output1JSON, output, "Response body mismatch")

	request, _ = http.NewRequest("DELETE", "/api/v2/topology/service-types?date=2015-04-12", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()

	// Check that we must have a 404 not found
	suite.Equal(404, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(output2JSON, output, "Response body mismatch")

}

func (suite *topologyTestSuite) TestListEndpoints2() {

	request, _ := http.NewRequest("GET", "/api/v2/topology/endpoints?date=2015-06-30", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	profileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-22",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_1",
   "hostname": "host1.site_a.foo",
   "tags": {
    "monitored": "1",
    "production": "1",
    "scope": "test,tier"
   }
  },
  {
   "date": "2015-06-22",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_2",
   "hostname": "host2.site_a.foo",
   "tags": {
    "monitored": "1",
    "production": "1",
    "scope": "tier1, lala"
   }
  },
  {
   "date": "2015-06-22",
   "group": "SITEB",
   "type": "SITES",
   "service": "service_1",
   "hostname": "host1.site_b.foo",
   "tags": {
    "monitored": "1",
    "production": "1",
    "scope": "tier1, tier2, foo"
   }
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(profileJSON, output, "Response body mismatch")

}

func (suite *topologyTestSuite) TestListEndpoints3() {

	request, _ := http.NewRequest("GET", "/api/v2/topology/endpoints?date=2015-07-30", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	profileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-07-22",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_1",
   "hostname": "host1.site_a.foo",
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  },
  {
   "date": "2015-07-22",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_2",
   "hostname": "host2.site_a.foo",
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  },
  {
   "date": "2015-07-22",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_3",
   "hostname": "host3.site_a.foo",
   "tags": {
    "monitored": "0",
    "production": "0",
    "scope": "TEST"
   }
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(profileJSON, output, "Response body mismatch")
}

func (suite *topologyTestSuite) TestListEndpoints4() {

	request, _ := http.NewRequest("GET", "/api/v2/topology/endpoints?date=2015-08-15", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	profileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-08-10",
   "group": "SITEA",
   "type": "SITES",
   "service": "service_x",
   "hostname": "host0.site_a.foo",
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  },
  {
   "date": "2015-08-10",
   "group": "SITEB",
   "type": "SITES",
   "service": "service_x",
   "hostname": "host0.site_b.foo",
   "notifications": {
    "contacts": [
     "contact01@email.example.foo",
     "contact02@email.example.foo"
    ],
    "enabled": true
   },
   "tags": {
    "monitored": "0",
    "production": "0"
   }
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(profileJSON, output, "Response body mismatch")

}

func (suite *topologyTestSuite) TestListEndpoints5() {

	request, _ := http.NewRequest("GET", "/api/v2/topology/endpoints?date=2015-01-01", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	profileJSON := `{
 "status": {
  "message": "Not Found",
  "code": "404"
 },
 "errors": [
  {
   "message": "Not Found",
   "code": "404",
   "details": "Specific query returned no items"
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(404, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(profileJSON, output, "Response body mismatch")
}

func (suite *topologyTestSuite) TestListGroups() {

	request, _ := http.NewRequest("GET", "/api/v2/topology/groups?date=2020-02-01", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	profileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2020-01-11",
   "group": "NGIX",
   "type": "NGIS",
   "subgroup": "SITEXYZ",
   "tags": {
    "certification": "Uncertified",
    "infrastructure": "Devel"
   }
  },
  {
   "date": "2020-01-11",
   "group": "NGIX",
   "type": "NGIS",
   "subgroup": "SITEXZ",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  },
  {
   "date": "2020-01-11",
   "group": "NGIX",
   "type": "NGIS",
   "subgroup": "SITEX",
   "tags": {
    "certification": "Uncertified",
    "infrastructure": "Production"
   }
  },
  {
   "date": "2020-01-11",
   "group": "PR01",
   "type": "PROJECT",
   "subgroup": "SITEPROJECT",
   "notifications": {
    "contacts": [
     "contact01@email.example.foo",
     "contact02@email.example.foo"
    ],
    "enabled": true
   },
   "tags": {
    "certification": "Certified",
    "infrastructure": "Devel"
   }
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(profileJSON, output, "Response body mismatch")

}

func (suite *topologyTestSuite) TestListGroups2() {

	request, _ := http.NewRequest("GET", "/api/v2/topology/groups?date=2015-06-30", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	profileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-06-22",
   "group": "NGIA",
   "type": "NGIS",
   "subgroup": "SITEA",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  },
  {
   "date": "2015-06-22",
   "group": "NGIA",
   "type": "NGIS",
   "subgroup": "SITEB",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(profileJSON, output, "Response body mismatch")

}

func (suite *topologyTestSuite) TestListGroups3() {

	request, _ := http.NewRequest("GET", "/api/v2/topology/groups?date=2015-07-30", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	profileJSON := `{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "date": "2015-07-22",
   "group": "NGIA",
   "type": "NGIS",
   "subgroup": "SITEA",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  },
  {
   "date": "2015-07-22",
   "group": "NGIA",
   "type": "NGIS",
   "subgroup": "SITEB",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  },
  {
   "date": "2015-07-22",
   "group": "NGIX",
   "type": "NGIS",
   "subgroup": "SITEX",
   "tags": {
    "certification": "Certified",
    "infrastructure": "Production"
   }
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(profileJSON, output, "Response body mismatch")

}

func (suite *topologyTestSuite) TestListGroups5() {

	request, _ := http.NewRequest("GET", "/api/v2/topology/groups?date=2015-01-01", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	profileJSON := `{
 "status": {
  "message": "Not Found",
  "code": "404"
 },
 "errors": [
  {
   "message": "Not Found",
   "code": "404",
   "details": "Specific query returned no items"
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(404, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(profileJSON, output, "Response body mismatch")
}

func (suite *topologyTestSuite) TestDeleteGroups() {

	request, _ := http.NewRequest("DELETE", "/api/v2/topology/groups?date=2015-07-22", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code := response.Code
	output := response.Body.String()

	output1JSON := `{
 "message": "Topology of 3 groups deleted for date: 2015-07-22",
 "code": "200"
}`

	output2JSON := `{
 "status": {
  "message": "Not Found",
  "code": "404"
 },
 "errors": [
  {
   "message": "Not Found",
   "code": "404",
   "details": "Specific query returned no items"
  }
 ]
}`
	// Check that we must have a 200 ok code
	suite.Equal(200, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(output1JSON, output, "Response body mismatch")

	request, _ = http.NewRequest("DELETE", "/api/v2/topology/groups?date=2015-08-10", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")
	response = httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)

	code = response.Code
	output = response.Body.String()

	// Check that we must have a 404 not found
	suite.Equal(404, code, "Internal Server Error")
	// Compare the expected and actual json response
	suite.Equal(output2JSON, output, "Response body mismatch")

}

// TestListTopologyStats
func (suite *topologyTestSuite) TestListTopologyStats() {

	expJSON := `{
 "status": {
  "message": "application/json",
  "code": "200"
 },
 "data": {
  "group_count": 2,
  "group_type": "GROUP",
  "group_list": [
   "GROUP_A",
   "GROUP_B"
  ],
  "endpoint_group_count": 4,
  "endpoint_group_type": "SITE",
  "endpoint_group_list": [
   "ST01",
   "ST02",
   "ST04",
   "ST05"
  ],
  "service_count": 3,
  "service_list": [
   "SF01",
   "SF02",
   "SF03"
  ]
 }
}`

	request, _ := http.NewRequest("GET", "/api/v2/topology/stats/Critical?date=2015-06-22", strings.NewReader(""))
	request.Header.Set("x-api-key", suite.clientkey)
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	suite.router.ServeHTTP(response, request)
	responseBody := response.Body.String()

	// Check that we must have a 200 ok code
	suite.Equal(200, response.Code, "Incorrect HTTP response code")
	// Compare the expected and actual xml response
	suite.Equal(expJSON, responseBody, "Response body mismatch")

}

//TearDownTest to tear down every test
func (suite *topologyTestSuite) TearDownTest() {

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
func (suite *topologyTestSuite) TearDownSuite() {

	session, err := mgo.Dial(suite.cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	session.DB(suite.tenantDbConf.Db).DropDatabase()
	session.DB(suite.cfg.MongoDB.Db).DropDatabase()
}

// TestTopologyTestSuite is responsible for calling the tests
func TestTopologyTestSuite(t *testing.T) {
	suite.Run(t, new(topologyTestSuite))
}
