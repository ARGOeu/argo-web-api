package main

import (
	"api/sites"
	"encoding/xml"
	"github.com/makistsan/go-lru-cache"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSiteAvailabilityInProfileQueryWithTwoSiteHostnames(t *testing.T) {
	httpcache = cache.NewLRUCache(uint64(700000000))
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

func TestSiteAvailabilityInProfileQueryWithOneSiteHostnameForOneDay(t *testing.T) {
	httpcache = cache.NewLRUCache(uint64(700000000))
	xmlStruct := sites.Root{}
	request, _ := http.NewRequest("GET", "/api/v1/sites_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T00:00:00Z&end_time=2013-08-01T23:59:00Z&type=DAILY&output=XML&profile_name=ROC_CRITICAL&service_hostname=sbdii.afroditi.hellasgrid.gr", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(SitesAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile[0].Site[0].Availability) != 1 {
		t.Error("Expected to find 24 availabilities, but instead found", len(xmlStruct.Profile[0].Site[0].Availability))
	}
}

func TestSiteAvailabilityInProfileQueryWithOneSiteHostnameForThreeDays(t *testing.T) {
	httpcache = cache.NewLRUCache(uint64(700000000))
	xmlStruct := sites.Root{}
	request, _ := http.NewRequest("GET", "/api/v1/sites_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T23:00:00Z&end_time=2013-08-03T23:59:00Z&type=DAILY&output=XML&profile_name=ROC_CRITICAL&service_hostname=sbdii.afroditi.hellasgrid.gr", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(SitesAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile[0].Site[0].Availability) != 3 {
		t.Error("Expected to find 72 availabilities, but instead found", len(xmlStruct.Profile[0].Site[0].Availability))
	}
}

func TestSiteAvailabilityInProfileQueryWithTwoSiteHostnameForThreeDays(t *testing.T) {
	httpcache = cache.NewLRUCache(uint64(700000000))
	xmlStruct := sites.Root{}
	request, _ := http.NewRequest("GET", "/api/v1/sites_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T23:00:00Z&end_time=2013-08-03T23:59:00Z&type=DAILY&output=XML&profile_name=ROC_CRITICAL&service_hostname=sbdii.afroditi.hellasgrid.gr&service_hostname=mon.kallisto.hellasgrid.gr", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(SitesAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile[0].Site[0].Availability)+len(xmlStruct.Profile[0].Site[1].Availability) != 6 {
		t.Error("Expected to find 144 availabilities, but instead found", len(xmlStruct.Profile[0].Site[0].Availability)+len(xmlStruct.Profile[0].Site[1].Availability))
	}
}

func TestSiteAvailabilityInProfileQueryWithTwoProfiles(t *testing.T) {
	t.Skip("Skipping test as we have data only for the ROC CRITICAL profile")
	httpcache = cache.NewLRUCache(uint64(700000000))
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
