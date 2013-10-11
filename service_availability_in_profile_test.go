package main

import (
	"encoding/xml"
	"github.com/makistsan/go-lru-cache"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateXMLResponse(t *testing.T) {
	var v []byte
	var tl Timeline
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

	results := []Timeline{tl}
	v, _ = createXMLResponse(results)
	if string(v) != expected_response {
		t.Error("XML response is not correct", string(v))
	}
}

type Availability struct {
	XMLName      xml.Name `xml:"Availability"`
	Timestamp    string   `xml:"timestamp,attr"`
	Availability string   `xml:"availability,attr"`
	Reliability  string   `xml:"reliability,attr"`
	Maintenance  string   `xml:"maintenance,attr"`
}
type Service struct {
	Hostname       string `xml:"hostname,attr"`
	Service_Type   string `xml:"type,attr"`
	Service_Flavor string `xml:"flavor,attr"`
	Availability   []*Availability
}
type Profile struct {
	Service []*Service
}
type Root struct {
	XMLName xml.Name `xml:"root"`
	Profile []*Profile
}
