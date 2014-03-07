/*
 * Copyright (c) 2013 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
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

/*
 * Test file for service availability calls.
 * Using the testing support provided by golang
 * we define the expected response and compare it
 * to the actual repsonse.
 */

package services

import (
	"testing"
)

func TestCreateXMLResponse(t *testing.T) {
	var v []byte

	var tl ApiServiceAvailabilityInProfileOutput
	customForm = []string{"20060102", "2006-01-02T15:04:05Z"}
	expected_response := ` <root>
   <Profile name="ROC_CRITICAL" namespace="ch.cern.sam" defined_by_vo_name="ops">
     <Service hostname="lcgbdii.phy.bris.ac.uk" type="CREAM-CE" flavor="CREAM-CE">
       <Availability timestamp="2013-08-01T00:00:00Z" availability="1" reliability="1" maintenance="-1"></Availability>
       <Availability timestamp="2013-08-01T01:00:00Z" availability="1" reliability="1" maintenance="-1"></Availability>
     </Service>
   </Profile>
 </root>`
	tl.Date = 20130801
	tl.Host = "lcgbdii.phy.bris.ac.uk"
	tl.Namespace = "ch.cern.sam"
	tl.Profile = "ROC_CRITICAL"
	tl.ServiceFlavor = "CREAM-CE"
	tl.Timeline = "[1:1:-1, 1:1:-1]\n"
	tl.VO = "ops"

	results := []ApiServiceAvailabilityInProfileOutput{tl}
	v, _ = CreateXMLResponse(results)
	if string(v) != expected_response {
		t.Error("XML response is not correct", string(v))
	}
}
