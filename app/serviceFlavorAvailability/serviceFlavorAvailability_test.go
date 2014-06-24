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

package serviceFlavorAvailability

import (
	"code.google.com/p/gcfg"
	"github.com/argoeu/ar-web-api/utils/config"
	"net/http"
	"testing"
)

//Default configuration has to be copied inside the source code in order for the test to be autonomous 
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
    db = "AR"
`
//DECLARATIONS OF EXPECTED OUTPUTS
const expectedOneDayOneSiteOneFlavorXML = ` <root>
   <Profile name="ch.cern.sam.ROC_CRITICAL">
     <Site Site="HG-03-AUTH">
       <Flavor Flavor="CREAM-CE">
         <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
       </Flavor>
     </Site>
   </Profile>
 </root>`

const expectedOneDayOneSiteOneFlavorJson = `{
   "Profile": [
     {
       "name": "ch.cern.sam.ROC_CRITICAL",
       "Site": [
         {
           "Site": "HG-03-AUTH",
           "SF": [
             {
               "Flavor": "CREAM-CE",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             }
           ]
         }
       ]
     }
   ]
 }`

const expectedOneDayOneSiteAllFlavorsXML = ` <root>
   <Profile name="ch.cern.sam.ROC_CRITICAL">
     <Site Site="HG-03-AUTH">
       <Flavor Flavor="CREAM-CE">
         <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
       </Flavor>
       <Flavor Flavor="SRMv2">
         <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
       </Flavor>
       <Flavor Flavor="Site-BDII">
         <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
       </Flavor>
     </Site>
   </Profile>
 </root>`

const expectedOneDayOneSiteAllFlavorsJson = `{
   "Profile": [
     {
       "name": "ch.cern.sam.ROC_CRITICAL",
       "Site": [
         {
           "Site": "HG-03-AUTH",
           "SF": [
             {
               "Flavor": "CREAM-CE",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             },
             {
               "Flavor": "SRMv2",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             },
             {
               "Flavor": "Site-BDII",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             }
           ]
         }
       ]
     }
   ]
 }`

const expectedOneDayTwoSitesAllFlavorsXML = ` <root>
   <Profile name="ch.cern.sam.ROC_CRITICAL">
     <Site Site="GR-01-AUTH">
       <Flavor Flavor="CREAM-CE">
         <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
       </Flavor>
       <Flavor Flavor="SRMv2">
         <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
       </Flavor>
       <Flavor Flavor="Site-BDII">
         <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
       </Flavor>
     </Site>
     <Site Site="HG-03-AUTH">
       <Flavor Flavor="CREAM-CE">
         <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
       </Flavor>
       <Flavor Flavor="SRMv2">
         <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
       </Flavor>
       <Flavor Flavor="Site-BDII">
         <Availability timestamp="2014-01-01" availability="100" reliability="100"></Availability>
       </Flavor>
     </Site>
   </Profile>
 </root>`

const expectedOneDayTwoSitesAllFlavorsJson = `{
   "Profile": [
     {
       "name": "ch.cern.sam.ROC_CRITICAL",
       "Site": [
         {
           "Site": "GR-01-AUTH",
           "SF": [
             {
               "Flavor": "CREAM-CE",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             },
             {
               "Flavor": "SRMv2",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             },
             {
               "Flavor": "Site-BDII",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             }
           ]
         },
         {
           "Site": "HG-03-AUTH",
           "SF": [
             {
               "Flavor": "CREAM-CE",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             },
             {
               "Flavor": "SRMv2",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             },
             {
               "Flavor": "Site-BDII",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             }
           ]
         }
       ]
     }
   ]
 }`

const expectedMultipleDaysTwoSitesAllFlavorsXML = ` <root>
   <Profile name="ch.cern.sam.ROC_CRITICAL">
     <Site Site="GR-01-AUTH">
       <Flavor Flavor="CREAM-CE">
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
       </Flavor>
       <Flavor Flavor="SRMv2">
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
       </Flavor>
       <Flavor Flavor="Site-BDII">
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
       </Flavor>
     </Site>
     <Site Site="HG-03-AUTH">
       <Flavor Flavor="CREAM-CE">
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
       </Flavor>
       <Flavor Flavor="SRMv2">
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
       </Flavor>
       <Flavor Flavor="Site-BDII">
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
       </Flavor>
     </Site>
   </Profile>
 </root>`

const expectedMultipleDaysTwoSitesAllFlavorsJson = `{
   "Profile": [
     {
       "name": "ch.cern.sam.ROC_CRITICAL",
       "Site": [
         {
           "Site": "GR-01-AUTH",
           "SF": [
             {
               "Flavor": "CREAM-CE",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-02",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-03",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-04",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-05",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-06",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-07",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-08",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-09",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-10",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             },
             {
               "Flavor": "SRMv2",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-02",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-03",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-04",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-05",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-06",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-07",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-08",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-09",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-10",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             },
             {
               "Flavor": "Site-BDII",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-02",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-03",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-04",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-05",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-06",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-07",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-08",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-09",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-10",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             }
           ]
         },
         {
           "Site": "HG-03-AUTH",
           "SF": [
             {
               "Flavor": "CREAM-CE",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-02",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-03",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-04",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-05",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-06",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-07",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-08",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-09",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-10",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             },
             {
               "Flavor": "SRMv2",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-02",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-03",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-04",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-05",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-06",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-07",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-08",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-09",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-10",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             },
             {
               "Flavor": "Site-BDII",
               "Availability": [
                 {
                   "Timestamp": "2014-01-01",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-02",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-03",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-04",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-05",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-06",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-07",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-08",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-09",
                   "availability": "100",
                   "reliability": "100"
                 },
                 {
                   "Timestamp": "2014-01-10",
                   "availability": "100",
                   "reliability": "100"
                 }
               ]
             }
           ]
         }
       ]
     }
   ]
 }`
 
// EXPECTED OUTPUTS END


//Preparing configuration struct
func prepare() config.Config {

	var cfg config.Config

	_ = gcfg.ReadStringInto(&cfg, defaultConfig)

	return cfg

}

//Tests for one day one egi site and one service flavor with xml formated output
func TestOneDayOneSiteOneFlavorXML(t *testing.T) {

	cfg := prepare()

	request, _ := http.NewRequest("GET", "?&start_time=2014-01-01T10:00:00Z&end_time=2014-01-01T10:00:00Z&granularity=daily&profile=ch.cern.sam.ROC_CRITICAL&site=HG-03-AUTH&flavor=CREAM-CE", nil)

	code, _, output, err := List(request, cfg)

	if code != http.StatusOK {
		t.Error("Error", err)
	} else if string(output) != expectedOneDayOneSiteOneFlavorXML {
		t.Error("Response body mismatch")
	}
}

//Tests for one day one egi site and one service flavor with json formated output
func TestOneDayOneSiteOneFlavorJson(t *testing.T) {

	cfg := prepare()

	request, _ := http.NewRequest("GET", "?&start_time=2014-01-01T10:00:00Z&end_time=2014-01-01T10:00:00Z&granularity=daily&profile=ch.cern.sam.ROC_CRITICAL&site=HG-03-AUTH&flavor=CREAM-CE&format=json", nil)

	code, _, output, err := List(request, cfg)

	if code != http.StatusOK {
		t.Error("Error", err)
	} else if string(output) != expectedOneDayOneSiteOneFlavorJson {
		t.Error("Response body mismatch")
	}
}

//Tests for one day one egi site and all service flavors with xml formated output
func TestOneDayOneSiteAllFlavorsXML(t *testing.T) {

	cfg := prepare()

	request, _ := http.NewRequest("GET", "?&start_time=2014-01-01T10:00:00Z&end_time=2014-01-01T10:00:00Z&granularity=daily&profile=ch.cern.sam.ROC_CRITICAL&site=HG-03-AUTH", nil)

	code, _, output, err := List(request, cfg)

	if code != http.StatusOK {
		t.Error("Error", err)
	} else if string(output) != expectedOneDayOneSiteAllFlavorsXML {
		t.Error("Response body mismatch")
	}
}

//Tests for one day one egi site and all service flavors with json formated output
func TestOneDayOneSiteAllFlavorsJson(t *testing.T) {

	cfg := prepare()

	request, _ := http.NewRequest("GET", "?&start_time=2014-01-01T10:00:00Z&end_time=2014-01-01T10:00:00Z&granularity=daily&profile=ch.cern.sam.ROC_CRITICAL&site=HG-03-AUTH&format=json", nil)

	code, _, output, err := List(request, cfg)

	if code != http.StatusOK {
		t.Error("Error", err)
	} else if string(output) != expectedOneDayOneSiteAllFlavorsJson {
		t.Error("Response body mismatch")
	}
}

//Tests for one day two egi sites and all service flavors with xml formated output
func TestOneDayTwoSitesAllFlavorsXML(t *testing.T) {

	cfg := prepare()

	request, _ := http.NewRequest("GET", "?&start_time=2014-01-01T10:00:00Z&end_time=2014-01-01T10:00:00Z&granularity=daily&profile=ch.cern.sam.ROC_CRITICAL&site=HG-03-AUTH&site=GR-01-AUTH", nil)

	code, _, output, err := List(request, cfg)

	if code != http.StatusOK {
		t.Error("Error", err)
	} else if string(output) != expectedOneDayTwoSitesAllFlavorsXML {
		t.Error("Response body mismatch")
	}
}

//Tests for one day two egi sites and all service flavors with json formated output
func TestOneDayTwoSitesAllFlavorsJson(t *testing.T) {

	cfg := prepare()

	request, _ := http.NewRequest("GET", "?&start_time=2014-01-01T10:00:00Z&end_time=2014-01-01T10:00:00Z&granularity=daily&profile=ch.cern.sam.ROC_CRITICAL&site=HG-03-AUTH&site=GR-01-AUTH", nil)

	code, _, output, err := List(request, cfg)

	if code != http.StatusOK {
		t.Error("Error", err)
	} else if string(output) != expectedOneDayTwoSitesAllFlavorsXML {
		t.Error("Response body mismatch")
	}
}

//Tests for multiple day two egi sites and all service flavors with xml formated output
func TestMultipleDaysTwoSitesAllFlavorsXML(t *testing.T) {

	cfg := prepare()

	request, _ := http.NewRequest("GET", "?&start_time=2014-01-01T10:00:00Z&end_time=2014-01-10T10:00:00Z&granularity=daily&profile=ch.cern.sam.ROC_CRITICAL&site=HG-03-AUTH&site=GR-01-AUTH", nil)

	code, _, output, err := List(request, cfg)

	if code != http.StatusOK {
		t.Error("Error", err)
	} else if string(output) != expectedMultipleDaysTwoSitesAllFlavorsXML {
		t.Error("Response body mismatch")
	}
}

//Tests for multiple day two egi sites and all service flavors with json formated output
func TestMultipleDaysTwoSitesAllFlavorsJson(t *testing.T) {

	cfg := prepare()

	request, _ := http.NewRequest("GET", "?&start_time=2014-01-01T10:00:00Z&end_time=2014-01-10T10:00:00Z&granularity=daily&profile=ch.cern.sam.ROC_CRITICAL&site=HG-03-AUTH&site=GR-01-AUTH&format=json", nil)

	code, _, output, err := List(request, cfg)

	if code != http.StatusOK {
		t.Error("Error", err)
	} else if string(output) != expectedMultipleDaysTwoSitesAllFlavorsJson {
		t.Error("Response body mismatch")
	}
}
