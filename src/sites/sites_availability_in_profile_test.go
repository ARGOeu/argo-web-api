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
