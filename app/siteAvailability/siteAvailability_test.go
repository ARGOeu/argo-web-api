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

package siteAvailability

import (
	"code.google.com/p/gcfg"
	"github.com/argoeu/ar-web-api/utils/config"
	"github.com/argoeu/ar-web-api/utils/mongo"
	"github.com/stretchr/testify/suite"
	"labix.org/v2/mgo/bson"
	"net/http"
	"testing"
)

type SiteTestSuite struct {
	suite.Suite
	cfg                       config.Config
	expectedOneDayOneSiteXML    string
	expectedTwoDaysOneSiteXML   string
	expectedOneMonthOneSiteXML  string
	expectedOneMonthTwoSitesXML string
	expectedMonthly           string
}

func (suite *SiteTestSuite) SetupTest() {

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
	
	seed := []bson.M{
		bson.M{"dt" : 20140101, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140102, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140103, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140104, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140105, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140106, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140107, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140108, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140109, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140110, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140111, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140112, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140113, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140114, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140115, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140116, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140117, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140118, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140119, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140120, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140121, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140122, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140123, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140124, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140125, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140126, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140127, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140128, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140129, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140130, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140131, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140201, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140202, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140203, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140204, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140205, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140206, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140207, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140208, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140209, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140210, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140211, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140212, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140213, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140214, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140215, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140216, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140217, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140218, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140219, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140220, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140221, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140222, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140223, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140224, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140225, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140226, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140227, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },
		bson.M{"dt" : 20140228, "s" : "GR-01-AUTH", "p" : "ch.cern.sam.ROC_CRITICAL", "pr" : "Y", "m" : "Y", "sc" : "EGI", "n" : "NGI_GRNET", "i" : "Production", "cs" : "Certified", "ss" : "EGI", "ap" : "test-ap1", "a" : 100, "r" : 100, "up" : 0, "u" : 0, "d" : 0, "hs" : 1 },		
	}
	
	session, _ := mongo.OpenSession(suite.cfg)

	_ = mongo.InsertMultiple(session, suite.cfg.MongoDB.Db, "sites", seed)


	suite.expectedOneDayOneSiteXML = ` <root>
   <Profile name="test-ap1">
     <Site site="GR-01-AUTH" NGI="NGI_GRNET" infastructure="Production" scope="EGI" site_scope="EGI" production="Y" monitored="Y" certification_status="Certified">
       <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
     </Site>
   </Profile>
 </root>`

// const expectedOneDayOneSiteJson = `{
//    "Profile": [
//      {
//        "name": "test-ap1",
//        "Site": [
//          {
//            "site": "GR-01-AUTH",
//            "Ngi": "NGI_GRNET",
//            "infrastructure": "Production",
//            "scope": "EGI",
//            "site_scope": "EGI",
//            "production": "Y",
//            "monitored": "Y",
//            "certification_status": "Certified",
//            "Availability": [
//              {
//                "timestamp": "2014-01-01",
//                "availability": "100",
//                "reliability": "100"
//              }
//            ]
//          }
//        ]
//      }
//    ]
//  }`

	suite.expectedTwoDaysOneSiteXML = ` <root>
   <Profile name="test-ap1">
     <Site site="GR-01-AUTH" NGI="NGI_GRNET" infastructure="Production" scope="EGI" site_scope="EGI" production="Y" monitored="Y" certification_status="Certified">
       <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
       <Availability timestamp="2014-01-02" availability="100" reliability="100"></Availability>
     </Site>
   </Profile>
 </root>`

// const expectedTwoDaysOneSiteJson = `{
//    "Profile": [
//      {
//        "name": "test-ap1",
//        "Site": [
//          {
//            "site": "GR-01-AUTH",
//            "Ngi": "NGI_GRNET",
//            "infrastructure": "Production",
//            "scope": "EGI",
//            "site_scope": "EGI",
//            "production": "Y",
//            "monitored": "Y",
//            "certification_status": "Certified",
//            "Availability": [
//              {
//                "timestamp": "2014-01-01",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-02",
//                "availability": "100",
//                "reliability": "100"
//              }
//            ]
//          }
//        ]
//      }
//    ]
//  }`

	suite.expectedOneMonthOneSiteXML = ` <root>
   <Profile name="test-ap1">
     <Site site="GR-01-AUTH" NGI="NGI_GRNET" infastructure="Production" scope="EGI" site_scope="EGI" production="Y" monitored="Y" certification_status="Certified">
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
     </Site>
   </Profile>
 </root>`

// const expectedOneMonthOneSiteJson = `{
//    "Profile": [
//      {
//        "name": "test-ap1",
//        "Site": [
//          {
//            "site": "GR-01-AUTH",
//            "Ngi": "NGI_GRNET",
//            "infrastructure": "Production",
//            "scope": "EGI",
//            "site_scope": "EGI",
//            "production": "Y",
//            "monitored": "Y",
//            "certification_status": "Certified",
//            "Availability": [
//              {
//                "timestamp": "2014-01-01",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-02",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-03",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-04",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-05",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-06",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-07",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-08",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-09",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-10",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-11",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-12",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-13",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-14",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-15",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-16",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-17",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-18",
//                "availability": "50",
//                "reliability": "50"
//              },
//              {
//                "timestamp": "2014-01-19",
//                "availability": "98.611",
//                "reliability": "98.611"
//              },
//              {
//                "timestamp": "2014-01-20",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-21",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-22",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-23",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-24",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-25",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-26",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-27",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-28",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-29",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-30",
//                "availability": "95.833",
//                "reliability": "95.833"
//              },
//              {
//                "timestamp": "2014-01-31",
//                "availability": "100",
//                "reliability": "100"
//              }
//            ]
//          }
//        ]
//      }
//    ]
//  }`

	suite.expectedOneMonthTwoSitesXML = ` <root>
   <Profile name="test-ap1">
     <Site site="GR-01-AUTH" NGI="NGI_GRNET" infastructure="Production" scope="EGI" site_scope="EGI" production="Y" monitored="Y" certification_status="Certified">
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
     </Site>
     <Site site="HG-03-AUTH" NGI="NGI_GRNET" infastructure="Production" scope="EGI" site_scope="EGI" production="Y" monitored="Y" certification_status="Certified">
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
     </Site>
   </Profile>
 </root>`

// const expectedOneMonthTwoSitesJson = `{
//    "Profile": [
//      {
//        "name": "test-ap1",
//        "Site": [
//          {
//            "site": "GR-01-AUTH",
//            "Ngi": "NGI_GRNET",
//            "infrastructure": "Production",
//            "scope": "EGI",
//            "site_scope": "EGI",
//            "production": "Y",
//            "monitored": "Y",
//            "certification_status": "Certified",
//            "Availability": [
//              {
//                "timestamp": "2014-01-01",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-02",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-03",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-04",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-05",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-06",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-07",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-08",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-09",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-10",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-11",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-12",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-13",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-14",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-15",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-16",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-17",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-18",
//                "availability": "50",
//                "reliability": "50"
//              },
//              {
//                "timestamp": "2014-01-19",
//                "availability": "98.611",
//                "reliability": "98.611"
//              },
//              {
//                "timestamp": "2014-01-20",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-21",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-22",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-23",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-24",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-25",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-26",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-27",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-28",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-29",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-30",
//                "availability": "95.833",
//                "reliability": "95.833"
//              },
//              {
//                "timestamp": "2014-01-31",
//                "availability": "100",
//                "reliability": "100"
//              }
//            ]
//          },
//          {
//            "site": "HG-03-AUTH",
//            "Ngi": "NGI_GRNET",
//            "infrastructure": "Production",
//            "scope": "EGI",
//            "site_scope": "EGI",
//            "production": "Y",
//            "monitored": "Y",
//            "certification_status": "Certified",
//            "Availability": [
//              {
//                "timestamp": "2014-01-01",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-02",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-03",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-04",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-05",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-06",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-07",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-08",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-09",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-10",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-11",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-12",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-13",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-14",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-15",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-16",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-17",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-18",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-19",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-20",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-21",
//                "availability": "-1",
//                "reliability": "-1"
//              },
//              {
//                "timestamp": "2014-01-22",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-23",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-24",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-25",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-26",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-27",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-28",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-29",
//                "availability": "100",
//                "reliability": "100"
//              },
//              {
//                "timestamp": "2014-01-30",
//                "availability": "91.667",
//                "reliability": "91.667"
//              },
//              {
//                "timestamp": "2014-01-31",
//                "availability": "100",
//                "reliability": "100"
//              }
//            ]
//          }
//        ]
//      }
//    ]
//  }`
}

func (suite *SiteTestSuite) TearDownTest() {

	session, _ := mongo.OpenSession(suite.cfg)

	_ = mongo.DropDatabase(session, suite.cfg.MongoDB.Db)

}


func (suite *SiteTestSuite) TestOneDayOneSiteXML() {

	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=site&start_time=2014-01-01T10:00:00Z&end_time=2014-01-01T10:00:00Z&group_name=GR-01-AUTH", nil)

	code, _, output, _ := List(request, suite.cfg)

	suite.NotEqual(code, 500, "Internal Server Error")
	suite.Equal(string(output), suite.expectedOneDayOneSiteXML, "Response body mismatch")

}

func (suite *SiteTestSuite) TestTwoDaysOneSiteXML() {

	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=site&start_time=2014-01-01T10:00:00Z&end_time=2014-01-02T10:00:00Z&group_name=GR-01-AUTH", nil)

	code, _, output, _ := List(request, suite.cfg)

	suite.NotEqual(code, 500, "Internal Server Error")
	suite.Equal(string(output), suite.expectedTwoDaysOneSiteXML, "Response body mismatch")

}

func (suite *SiteTestSuite) TestOneMonthOneSiteXML() {

	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=site&start_time=2014-01-01T10:00:00Z&end_time=2014-01-31T10:00:00Z&group_name=GR-01-AUTH", nil)

	code, _, output, _ := List(request, suite.cfg)

	suite.NotEqual(code, 500, "Internal Server Error")
	suite.Equal(string(output), suite.expectedOneMonthOneSiteXML, "Response body mismatch")

}

// func (suite *SiteTestSuite) TestOneMonthTwoSitesXML() {
//
// 	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=site&start_time=2014-01-01T10:00:00Z&end_time=2014-01-31T10:00:00Z", nil)
//
// 	code, _, output, _ := List(request, suite.cfg)
//
// 	suite.NotEqual(code, 500, "Internal Server Error")
// 	suite.Equal(string(output), suite.expectedOneMonthTwoSitesXML, "Response body mismatch")
//
// }


func TestSiteTestSuite(t *testing.T) {
	suite.Run(t, new(SiteTestSuite))
}


// func TestOneDayOneSiteJson(t *testing.T) {
//
// 	cfg := prepare()
//
// 	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=site&start_time=2014-01-01T10:00:00Z&end_time=2014-01-01T10:00:00Z&group_name=GR-01-AUTH&format=json", nil)
//
// 	code, _, output, err := List(request, cfg)
//
// 	if code != http.StatusOK {
// 		t.Error("Error", err)
// 	} else if string(output) != expectedOneDayOneSiteJson {
// 		t.Error("Response body mismatch")
// 	}
// }
//
// //Tests for one day one egi site and one service flavor with xml formated output
// func TestTwoDaysOneSiteXML(t *testing.T) {
//
// 	cfg := prepare()
//
// 	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=site&start_time=2014-01-01T10:00:00Z&end_time=2014-01-02T10:00:00Z&group_name=GR-01-AUTH", nil)
//
// 	code, _, output, err := List(request, cfg)
//
// 	if code != http.StatusOK {
// 		t.Error("Error", err)
// 	} else if string(output) != expectedTwoDaysOneSiteXML {
// 		t.Error("Response body mismatch")
// 	}
// }
//
// //Tests for one day one egi site and one service flavor with xml formated output
// func TestTwoDaysOneSiteJson(t *testing.T) {
//
// 	cfg := prepare()
//
// 	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=site&start_time=2014-01-01T10:00:00Z&end_time=2014-01-02T10:00:00Z&group_name=GR-01-AUTH&format=json", nil)
//
// 	code, _, output, err := List(request, cfg)
//
// 	if code != http.StatusOK {
// 		t.Error("Error", err)
// 	} else if string(output) != expectedTwoDaysOneSiteJson {
// 		t.Error("Response body mismatch")
// 	}
// }
//
// func TestOneMonthOneSiteXML(t *testing.T) {
//
// 	cfg := prepare()
//
// 	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=site&start_time=2014-01-01T10:00:00Z&end_time=2014-01-31T10:00:00Z&group_name=GR-01-AUTH", nil)
//
// 	code, _, output, err := List(request, cfg)
//
// 	if code != http.StatusOK {
// 		t.Error("Error", err)
// 	} else if string(output) != expectedOneMonthOneSiteXML {
// 		t.Error("Response body mismatch")
// 	}
// }
//
// func TestOneMonthOneSiteJson(t *testing.T) {
//
// 	cfg := prepare()
//
// 	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=site&start_time=2014-01-01T10:00:00Z&end_time=2014-01-31T10:00:00Z&group_name=GR-01-AUTH&format=json", nil)
//
// 	code, _, output, err := List(request, cfg)
//
// 	if code != http.StatusOK {
// 		t.Error("Error", err)
// 	} else if string(output) != expectedOneMonthOneSiteJson {
// 		t.Error("Response body mismatch")
// 	}
// }
//
// func TestOneMonthTwoSitesXML(t *testing.T) {
//
// 	cfg := prepare()
//
// 	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=site&start_time=2014-01-01T10:00:00Z&end_time=2014-01-31T10:00:00Z&group_name=GR-01-AUTH&group_name=HG-03-AUTH", nil)
//
// 	code, _, output, err := List(request, cfg)
//
// 	if code != http.StatusOK {
// 		t.Error("Error", err)
// 	} else if string(output) != expectedOneMonthTwoSitesXML {
// 		t.Error("Response body mismatch")
// 	}
// }
//
// func TestOneMonthTwoSitesJson(t *testing.T) {
//
// 	cfg := prepare()
//
// 	request, _ := http.NewRequest("GET", "?availability_profile=test-ap1&group_type=site&start_time=2014-01-01T10:00:00Z&end_time=2014-01-31T10:00:00Z&group_name=GR-01-AUTH&group_name=HG-03-AUTH&format=json", nil)
//
// 	code, _, output, err := List(request, cfg)
//
// 	if code != http.StatusOK {
// 		t.Error("Error", err)
// 	} else if string(output) != expectedOneMonthTwoSitesJson {
// 		t.Error("Response body mismatch")
// 	}
// }
