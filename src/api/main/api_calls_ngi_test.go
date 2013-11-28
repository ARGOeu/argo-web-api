package main

import (
	"api/ngis"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNgiAvailabilityInProfileQueryWithTwoNgisForTwoDay(t *testing.T) {
	cfg.Server.Cache = false
	xmlStruct := ngis.Root{}
	request, _ := http.NewRequest("GET", "/api/v1/ngi_availability_in_profile?profile?start_time=2013-08-01T00:00:00Z&end_time=2013-08-02T23:59:00Z&group_name=Russia&group_name=NGI_GRNET&type=DAILY&output=XML&profile_name=ROC_CRITICAL", nil)
	response := httptest.NewRecorder()
	err := xml.Unmarshal([]byte(NgiAvailabilityInProfile(response, request)), &xmlStruct)
	if err != nil {
		t.Error("Error parsing XML file: %v", err)
	} else if len(xmlStruct.Profile[0].Ngi) != 2 {
		t.Error("Expected to find 2 Ngis, but instead found", len(xmlStruct.Profile[0].Ngi))
	} else if len(xmlStruct.Profile[0].Ngi[0].Availability) != 2 {
		t.Error("Expected to find 2 Availabilities, but instead found", len(xmlStruct.Profile[0].Ngi[0].Availability))
	}
}
