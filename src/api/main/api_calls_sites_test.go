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

package main

import (
	"api/sites"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSiteAvailabilityInProfileQueryWithTwoSitesForOneDay(t *testing.T) {
	t.Skip("Skipping test because we dont have a filter for Sites")
	cfg.Server.Cache = false
	xmlStruct := sites.Root{}
	request, _ := http.NewRequest("GET", "/api/v1/sites_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T00:00:00Z&end_time=2013-08-02T23:59:00Z&type=DAILY&output=XML&profile_name=ROC_CRITICAL&service_hostname=sbdii.afroditi.hellasgrid.gr&service_hostname=mon.kallisto.hellasgrid.gr", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(SitesAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile[0].Site) != 2 {
		t.Error("Expected to find 2 services, but instead found", len(xmlStruct.Profile[0].Site))
	}
}

func TestSiteAvailabilityInProfileQueryWithTwoProfiles(t *testing.T) {
	t.Skip("Skipping test as we have data only for the ROC CRITICAL profile")
	cfg.Server.Cache = false
	xmlStruct := sites.Root{}
	request, _ := http.NewRequest("GET", "/api/v1/sites_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T23:00:00Z&end_time=2013-08-03T23:59:00Z&type=DAILY&output=XML&profile_name=ROC_CRITICAL&profile_name=ROC", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(SitesAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile) != 2 {
		t.Error("Expected to find 2 Profile, but instead found", len(xmlStruct.Profile[0].Site[0].Availability))
	}
}
