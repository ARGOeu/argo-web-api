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
	"net/http"
	"testing"
)

type VOTestSuite struct {
	suite.Suite
	cfg                       config.Config
	expectedOneDayOneVOXML    string
	expectedTwoDaysOneVOXML   string
	expectedOneMonthOneVOXML  string
	expectedOneMonthTwoVOsXML string
	expectedMonthly           string
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
	seed := []bson.M{
		bson.M{"dt": 20140101, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140102, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140103, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140104, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140105, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140106, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140107, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140108, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140109, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140110, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140111, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140112, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140113, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140114, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140115, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140116, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140117, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140118, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140119, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140120, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140121, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140122, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140123, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140124, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140125, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140126, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140127, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140128, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140129, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140130, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140131, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140101, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140102, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140103, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140104, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140105, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140106, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140107, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140108, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140109, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140110, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140111, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140112, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140113, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140114, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140115, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140116, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140117, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140118, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140119, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140120, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140121, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140122, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140123, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140124, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140125, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140126, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140127, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140128, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140129, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140130, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140131, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140201, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140202, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140203, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140204, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140205, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140206, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140207, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140208, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140209, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140210, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140211, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140212, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140213, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140214, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140215, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140216, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140217, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140218, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140219, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140220, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140221, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140222, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140223, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140224, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140225, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140226, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140227, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140228, "v": "test-ops1", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140201, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140202, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140203, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140204, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140205, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140206, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140207, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140208, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140209, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140210, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140211, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140212, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140213, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140214, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140215, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140216, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140217, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140218, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140219, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140220, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140221, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140222, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140223, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140224, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140225, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140226, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140227, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
		bson.M{"dt": 20140228, "v": "test-ops2", "p": "ch.cern.sam.ROC_CRITICAL", "ap": "test-ap1", "a": 100, "r": 100, "up": 0.99306, "u": 0.00694, "d": 0},
	}

	session, _ := mongo.OpenSession(suite.cfg)

	_ = mongo.InsertMultiple(session, suite.cfg.MongoDB.Db, "voreports", seed)

	//SEED END

	//EXPECTED OUTPUT DEF
	suite.expectedOneDayOneVOXML = ` <root>
   <Profile name="test-ap1">
     <Vo VO="test-ops1">
       <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
     </Vo>
   </Profile>
 </root>`

	suite.expectedTwoDaysOneVOXML = ` <root>
   <Profile name="test-ap1">
     <Vo VO="test-ops1">
       <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-02" availability="100" reliability="100"></Availability>
     </Vo>
   </Profile>
 </root>`

	suite.expectedOneMonthOneVOXML = ` <root>
   <Profile name="test-ap1">
     <Vo VO="test-ops1">
       <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-02" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-03" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-04" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-05" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-06" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-07" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-08" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-09" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-10" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-11" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-12" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-13" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-14" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-15" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-16" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-17" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-18" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-19" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-20" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-21" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-22" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-23" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-24" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-25" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-26" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-27" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-28" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-29" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-30" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-31" availability="100" reliability="100"></Availability>
     </Vo>
   </Profile>
 </root>`

	suite.expectedOneMonthTwoVOsXML = ` <root>
   <Profile name="test-ap1">
     <Vo VO="test-ops1">
       <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-02" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-03" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-04" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-05" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-06" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-07" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-08" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-09" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-10" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-11" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-12" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-13" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-14" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-15" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-16" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-17" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-18" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-19" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-20" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-21" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-22" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-23" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-24" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-25" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-26" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-27" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-28" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-29" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-30" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-31" availability="100" reliability="100"></Availability>
     </Vo>
     <Vo VO="test-ops2">
       <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-02" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-03" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-04" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-05" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-06" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-07" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-08" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-09" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-10" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-11" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-12" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-13" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-14" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-15" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-16" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-17" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-18" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-19" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-20" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-21" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-22" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-23" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-24" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-25" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-26" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-27" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-28" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-29" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-30" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-31" availability="100" reliability="100"></Availability>
     </Vo>
   </Profile>
 </root>`

	suite.expectedMonthly = ` <root>
   <Profile name="test-ap1">
     <Vo VO="test-ops1">
       <Availability timestamp="2014-01" availability="99.9999989930115" reliability="99.9999989930115"></Availability>
       <Availability timestamp="2014-02" availability="99.9999989930115" reliability="99.9999989930115"></Availability>
     </Vo>
     <Vo VO="test-ops2">
       <Availability timestamp="2014-01" availability="99.9999989930115" reliability="99.9999989930115"></Availability>
       <Availability timestamp="2014-02" availability="99.9999989930115" reliability="99.9999989930115"></Availability>
     </Vo>
   </Profile>
 </root>`

	//EXPECTED OUTPUT DEF

	mongo.CloseSession(session)

}

func (suite *VOTestSuite) TearDownTest() {

	session, _ := mongo.OpenSession(suite.cfg)

	_ = mongo.DropDatabase(session, suite.cfg.MongoDB.Db)

}

func (suite *VOTestSuite) TestOneDayOneVOXML() {

	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=vo&start_time=2014-01-01T10:00:00Z&end_time=2014-01-01T10:00:00Z&granularity=daily&format=XML&group_name=test-ops1", nil)

	code, _, output, _ := List(request, suite.cfg)

	suite.NotEqual(code, 500, "Internal Server Error")
	suite.Equal(string(output), suite.expectedOneDayOneVOXML, "Response body mismatch")

}

func (suite *VOTestSuite) TestTwoDaysOneVOXML() {

	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=vo&start_time=2014-01-01T10:00:00Z&end_time=2014-01-02T10:00:00Z&granularity=daily&format=XML&group_name=test-ops1", nil)

	code, _, output, _ := List(request, suite.cfg)

	suite.NotEqual(code, 500, "Internal Server Error")
	suite.Equal(string(output), suite.expectedTwoDaysOneVOXML, "Response body mismatch")

}

func (suite *VOTestSuite) TestOneMonthOneVOXML() {

	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=vo&start_time=2014-01-01T10:00:00Z&end_time=2014-01-31T10:00:00Z&granularity=daily&format=XML&group_name=test-ops1", nil)

	code, _, output, _ := List(request, suite.cfg)

	suite.NotEqual(code, 500, "Internal Server Error")
	suite.Equal(string(output), suite.expectedOneMonthOneVOXML, "Response body mismatch")

}

func (suite *VOTestSuite) TestOneMonthTwoVOsXML() {

	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=vo&start_time=2014-01-01T10:00:00Z&end_time=2014-01-31T10:00:00Z&granularity=daily&format=XML", nil)

	code, _, output, _ := List(request, suite.cfg)

	suite.NotEqual(code, 500, "Internal Server Error")
	suite.Equal(string(output), suite.expectedOneMonthTwoVOsXML, "Response body mismatch")

}

func (suite *VOTestSuite) TestMonthy() {

	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=vo&start_time=2014-01-01T10:00:00Z&end_time=2014-02-28T10:00:00Z&granularity=monthly&format=XML", nil)

	code, _, output, _ := List(request, suite.cfg)

	suite.NotEqual(code, 500, "Internal Server Error")
	suite.Equal(string(output), suite.expectedMonthly, "Response body mismatch")

}

func TestVOTestSuite(t *testing.T) {
	suite.Run(t, new(VOTestSuite))
}
