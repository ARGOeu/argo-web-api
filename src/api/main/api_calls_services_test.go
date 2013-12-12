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
	"api/services"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceAvailabilityInProfileQueryWithTwoServiceHostnames(t *testing.T) {
	cfg.Server.Cache = false
	xmlStruct := services.Root{}
	request, _ := http.NewRequest("GET", "/api/v1/service_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T23:00:00Z&end_time=2013-08-01T23:59:00Z&type=HOURLY&output=XML&profile_name=ROC_CRITICAL&service_hostname=sbdii.afroditi.hellasgrid.gr&service_hostname=mon.kallisto.hellasgrid.gr", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(ServiceAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile[0].Service) != 2 {
		t.Error("Expected to find 2 services, but instead found", len(xmlStruct.Profile[0].Service))
	}
}

func TestServiceAvailabilityInProfileQueryWithOneServiceHostnameForOneDay(t *testing.T) {
	cfg.Server.Cache = false
	xmlStruct := services.Root{}
	request, _ := http.NewRequest("GET", "/api/v1/service_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T23:00:00Z&end_time=2013-08-01T23:59:00Z&type=HOURLY&output=XML&profile_name=ROC_CRITICAL&service_hostname=sbdii.afroditi.hellasgrid.gr", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(ServiceAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile[0].Service[0].Availability) != 24 {
		t.Error("Expected to find 24 availabilities, but instead found", len(xmlStruct.Profile[0].Service[0].Availability))
	}
}

func TestServiceAvailabilityInProfileQueryWithOneServiceHostnameForThreeDays(t *testing.T) {
	cfg.Server.Cache = false
	xmlStruct := services.Root{}
	request, _ := http.NewRequest("GET", "/api/v1/service_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T23:00:00Z&end_time=2013-08-03T23:59:00Z&type=HOURLY&output=XML&profile_name=ROC_CRITICAL&service_hostname=sbdii.afroditi.hellasgrid.gr", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(ServiceAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile[0].Service[0].Availability) != 72 {
		t.Error("Expected to find 72 availabilities, but instead found", len(xmlStruct.Profile[0].Service[0].Availability))
	}
}

func TestServiceAvailabilityInProfileQueryWithTwoServiceHostnameForThreeDays(t *testing.T) {
	cfg.Server.Cache = false
	xmlStruct := services.Root{}
	request, _ := http.NewRequest("GET", "/api/v1/service_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T23:00:00Z&end_time=2013-08-03T23:59:00Z&type=HOURLY&output=XML&profile_name=ROC_CRITICAL&service_hostname=sbdii.afroditi.hellasgrid.gr&service_hostname=mon.kallisto.hellasgrid.gr", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(ServiceAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile[0].Service[0].Availability)+len(xmlStruct.Profile[0].Service[1].Availability) != 144 {
		t.Error("Expected to find 144 availabilities, but instead found", len(xmlStruct.Profile[0].Service[0].Availability)+len(xmlStruct.Profile[0].Service[1].Availability))
	}
}

func TestServiceAvailabilityInProfileQueryWithTwoProfiles(t *testing.T) {
	t.Skip("Skipping test as we have data only for the ROC CRITICAL profile")
	cfg.Server.Cache = false
	xmlStruct := services.Root{}
	request, _ := http.NewRequest("GET", "/api/v1/service_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T23:00:00Z&end_time=2013-08-03T23:59:00Z&type=HOURLY&output=XML&profile_name=ROC_CRITICAL&profile_name=ROC", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(ServiceAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile) != 2 {
		t.Error("Expected to find 2 Profile, but instead found", len(xmlStruct.Profile[0].Service[0].Availability))
	}
}
