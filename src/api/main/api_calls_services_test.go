package main

import (
	"api/services"
	"encoding/xml"
	"github.com/makistsan/go-lru-cache"
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
