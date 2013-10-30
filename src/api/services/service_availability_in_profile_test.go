package services

import (
	"testing"
)

func TestCreateXMLResponse(t *testing.T) {
	var v []byte

	var tl Timeline
	customForm := []string{"20060102", "2006-01-02T15:04:05Z"}
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
	v, _ = CreateXMLResponse(results, customForm)
	if string(v) != expected_response {
		t.Error("XML response is not correct", string(v))
	}
}
