package main

import (
	"encoding/xml"
	"github.com/makistsan/go-lru-cache"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSiteAvailabilityInProfileQueryWithOneSiteHostname(t *testing.T) {
	type Site struct {
		Site          string `xml:"site,attr"`
		Ngi           string `xml:"NGI,attr"`
		Infastructure string `xml:"infastructure,attr"`
		Scope         string `xml:"scope,attr"`
		SiteScope     string `xml:"site_scope,attr"`
		Production    string `xml:"production,attr"`
		Monitored     string `xml:"monitored,attr"`
		CertStatus    string `xml:"certification_status,attr"`
		Availability  []*Availability
	}

	type Profile struct {
		XMLName   xml.Name `xml:"Profile"`
		Name      string   `xml:"name,attr"`
		Namespace string   `xml:"namespace,attr"`
		Site      []*Site
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Profile []*Profile
	}

	httpcache = cache.NewLRUCache(uint64(700000000))
	xmlStruct := Root{}
	request, _ := http.NewRequest("GET", "/api/v1/service_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T23:00:00Z&end_time=2013-08-01T23:59:00Z&type=HOURLY&output=XML&profile_name=ROC_CRITICAL&service_hostname=sbdii.afroditi.hellasgrid.gr", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(SitesAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile[0].Site) != 1 {
		t.Error("Expected to find 1 services, but instead found", len(xmlStruct.Profile[0].Site))
	}
}

func TestSiteAvailabilityInProfileQueryWithTwoSiteHostnames(t *testing.T) {
	type Site struct {
		Site          string `xml:"site,attr"`
		Ngi           string `xml:"NGI,attr"`
		Infastructure string `xml:"infastructure,attr"`
		Scope         string `xml:"scope,attr"`
		SiteScope     string `xml:"site_scope,attr"`
		Production    string `xml:"production,attr"`
		Monitored     string `xml:"monitored,attr"`
		CertStatus    string `xml:"certification_status,attr"`
		Availability  []*Availability
	}

	type Profile struct {
		XMLName   xml.Name `xml:"Profile"`
		Name      string   `xml:"name,attr"`
		Namespace string   `xml:"namespace,attr"`
		Site      []*Site
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Profile []*Profile
	}

	httpcache = cache.NewLRUCache(uint64(700000000))
	xmlStruct := Root{}
	request, _ := http.NewRequest("GET", "/api/v1/service_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T23:00:00Z&end_time=2013-08-01T23:59:00Z&type=HOURLY&output=XML&profile_name=ROC_CRITICAL&service_hostname=sbdii.afroditi.hellasgrid.gr&service_hostname=mon.kallisto.hellasgrid.gr", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(SitesAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile[0].Site) != 2 {
		t.Error("Expected to find 2 services, but instead found", len(xmlStruct.Profile[0].Site))
	}
}

func TestSiteAvailabilityInProfileQueryWithOneSiteHostnameForOneDay(t *testing.T) {
	type Site struct {
		Site          string `xml:"site,attr"`
		Ngi           string `xml:"NGI,attr"`
		Infastructure string `xml:"infastructure,attr"`
		Scope         string `xml:"scope,attr"`
		SiteScope     string `xml:"site_scope,attr"`
		Production    string `xml:"production,attr"`
		Monitored     string `xml:"monitored,attr"`
		CertStatus    string `xml:"certification_status,attr"`
		Availability  []*Availability
	}

	type Profile struct {
		XMLName   xml.Name `xml:"Profile"`
		Name      string   `xml:"name,attr"`
		Namespace string   `xml:"namespace,attr"`
		Site      []*Site
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Profile []*Profile
	}

	httpcache = cache.NewLRUCache(uint64(700000000))
	xmlStruct := Root{}
	request, _ := http.NewRequest("GET", "/api/v1/service_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T23:00:00Z&end_time=2013-08-01T23:59:00Z&type=HOURLY&output=XML&profile_name=ROC_CRITICAL&service_hostname=sbdii.afroditi.hellasgrid.gr", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(SitesAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile[0].Site[0].Availability) != 24 {
		t.Error("Expected to find 24 availabilities, but instead found", len(xmlStruct.Profile[0].Site[0].Availability))
	}
}

func TestSiteAvailabilityInProfileQueryWithOneSiteHostnameForThreeDays(t *testing.T) {
	type Site struct {
		Site          string `xml:"site,attr"`
		Ngi           string `xml:"NGI,attr"`
		Infastructure string `xml:"infastructure,attr"`
		Scope         string `xml:"scope,attr"`
		SiteScope     string `xml:"site_scope,attr"`
		Production    string `xml:"production,attr"`
		Monitored     string `xml:"monitored,attr"`
		CertStatus    string `xml:"certification_status,attr"`
		Availability  []*Availability
	}

	type Profile struct {
		XMLName   xml.Name `xml:"Profile"`
		Name      string   `xml:"name,attr"`
		Namespace string   `xml:"namespace,attr"`
		Site      []*Site
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Profile []*Profile
	}

	httpcache = cache.NewLRUCache(uint64(700000000))
	xmlStruct := Root{}
	request, _ := http.NewRequest("GET", "/api/v1/service_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T23:00:00Z&end_time=2013-08-03T23:59:00Z&type=HOURLY&output=XML&profile_name=ROC_CRITICAL&service_hostname=sbdii.afroditi.hellasgrid.gr", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(SitesAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile[0].Site[0].Availability) != 72 {
		t.Error("Expected to find 72 availabilities, but instead found", len(xmlStruct.Profile[0].Site[0].Availability))
	}
}

func TestSiteAvailabilityInProfileQueryWithTwoSiteHostnameForThreeDays(t *testing.T) {
	type Site struct {
		Site          string `xml:"site,attr"`
		Ngi           string `xml:"NGI,attr"`
		Infastructure string `xml:"infastructure,attr"`
		Scope         string `xml:"scope,attr"`
		SiteScope     string `xml:"site_scope,attr"`
		Production    string `xml:"production,attr"`
		Monitored     string `xml:"monitored,attr"`
		CertStatus    string `xml:"certification_status,attr"`
		Availability  []*Availability
	}

	type Profile struct {
		XMLName   xml.Name `xml:"Profile"`
		Name      string   `xml:"name,attr"`
		Namespace string   `xml:"namespace,attr"`
		Site      []*Site
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Profile []*Profile
	}

	httpcache = cache.NewLRUCache(uint64(700000000))
	xmlStruct := Root{}
	request, _ := http.NewRequest("GET", "/api/v1/service_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T23:00:00Z&end_time=2013-08-03T23:59:00Z&type=HOURLY&output=XML&profile_name=ROC_CRITICAL&service_hostname=sbdii.afroditi.hellasgrid.gr&service_hostname=mon.kallisto.hellasgrid.gr", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(SitesAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile[0].Site[0].Availability)+len(xmlStruct.Profile[0].Site[1].Availability) != 144 {
		t.Error("Expected to find 144 availabilities, but instead found", len(xmlStruct.Profile[0].Site[0].Availability)+len(xmlStruct.Profile[0].Site[1].Availability))
	}
}

func TestSiteAvailabilityInProfileQueryWithTwoProfiles(t *testing.T) {
	type Site struct {
		Site          string `xml:"site,attr"`
		Ngi           string `xml:"NGI,attr"`
		Infastructure string `xml:"infastructure,attr"`
		Scope         string `xml:"scope,attr"`
		SiteScope     string `xml:"site_scope,attr"`
		Production    string `xml:"production,attr"`
		Monitored     string `xml:"monitored,attr"`
		CertStatus    string `xml:"certification_status,attr"`
		Availability  []*Availability
	}

	type Profile struct {
		XMLName   xml.Name `xml:"Profile"`
		Name      string   `xml:"name,attr"`
		Namespace string   `xml:"namespace,attr"`
		Site      []*Site
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Profile []*Profile
	}

	t.Skip("Skipping test as we have data only for the ROC CRITICAL profile")
	httpcache = cache.NewLRUCache(uint64(700000000))
	xmlStruct := Root{}
	request, _ := http.NewRequest("GET", "/api/v1/service_availability_in_profile?vo_name=ops&group_type=Site&start_time=2013-08-01T23:00:00Z&end_time=2013-08-03T23:59:00Z&type=HOURLY&output=XML&profile_name=ROC_CRITICAL&profile_name=ROC", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(SitesAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile) != 2 {
		t.Error("Expected to find 2 Profile, but instead found", len(xmlStruct.Profile[0].Site[0].Availability))
	}
}
