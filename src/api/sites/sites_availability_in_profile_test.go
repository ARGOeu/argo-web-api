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
 * Test file for site availability calls.
 * Using the testing support provided by golang
 * we define the expected response and compare it 
 * to the actual repsonse.
*/



package sites

import (
	"testing"
)

func TestCreateSiteXMLResponse(t *testing.T) {

	var v []byte
	var st MongoSite
	expected_response := ` <root>
   <Profile name="ROC_CRITICAL" namespace="ch.cern.sam">
     <Site site="JINR-LCG2" NGI="Russia" infastructure="Production" scope="EGI" site_scope="EGI" production="Y" monitored="Y" certification_status="Certified">
       <Availability timestamp="2013-08" availability="99.9" reliability="99.9"></Availability>
     </Site>
   </Profile>
 </root>`
	st.Site = "JINR-LCG2"
	st.Ngi = "Russia"
	st.Infastructure = "Production"
	st.Scope = "EGI"
	st.SiteScope = "EGI"
	st.Production = "Y"
	st.Monitored = "Y"
	st.CertStatus = "Certified"
	st.Date = "201308"
	st.Namespace = "ch.cern.sam"
	st.Profile = "ROC_CRITICAL"
	st.Availability = 99.9
	st.Reliability = 99.9
	results := []MongoSite{st}
	customForm := []string{"200601", "2006-01"}
	v, _ = CreateXMLResponse(results, customForm)
	if string(v) != expected_response {
		t.Error("XML response is not correct", string(v))
	}
}
