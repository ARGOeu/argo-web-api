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

package voAvailability

import (
	"code.google.com/p/gcfg"
	"github.com/argoeu/ar-web-api/utils/config"
	"github.com/argoeu/ar-web-api/utils/mongo"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo/bson"
	//"labix.org/v2/mgo"
	"net/http"
	"testing"
)

type VOTestSuite struct {
	suite.Suite
	cfg                    config.Config
	expectedOneDayOneVOXML string
}

func (suite *VOTestSuite) SetupTest() {

	const defaultConfig = `
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
    db = "AR_test"
`
	
	_ = gcfg.ReadStringInto(&suite.cfg, defaultConfig)
	
	//SEED
	seed := bson.M{ "dt" : 20140101, "v" : "ops", "p" : "ch.cern.sam.ROC_CRITICAL", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0.99306, "u" : 0.00694, "d" : 0 }
	session, _ := mongo.OpenSession(suite.cfg)
    _ = mongo.Insert(session, suite.cfg.MongoDB.Db, "voreports", seed)

	suite.expectedOneDayOneVOXML = ` <root>
   <Profile name="test-ap1">
     <Vo VO="ops">
       <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
     </Vo>
   </Profile>
 </root>`
  
	mongo.CloseSession(session)
  

}

func (suite *VOTestSuite) TearDownTest() {

	session, _ := mongo.OpenSession(suite.cfg)

	_ = session.DB(suite.cfg.MongoDB.Db).DropDatabase()

}

func (suite *VOTestSuite) TestOneDayOneVOXML() {

	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=vo&start_time=2014-01-01T10:00:00Z&end_time=2014-01-01T10:00:00Z&granularity=daily&format=XML", nil)

	code, _, output, _ := List(request, suite.cfg)
	
	suite.NotEqual(code,500,"Internal Server Error")
	suite.Equal(string(output), suite.expectedOneDayOneVOXML, "Response body mismatch")

}

func TestVOTestSuite(t *testing.T) {
	suite.Run(t, new(VOTestSuite))
}
