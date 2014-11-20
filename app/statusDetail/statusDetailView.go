package statusDetail

import "encoding/xml"

func createView(results []StatusDetailOutput) ([]byte, error) {

	docRoot := &ReadRoot{}

	if len(results) == 0 {
		output, err := xml.MarshalIndent(docRoot, " ", "  ")
		return output, err
	}

	profile := &Profile{}
	profile.Name = "ch.cern.sam.ROC_CRITICAL"

	vo := &Group{}
	vo.Type = "vo"
	vo.Name = "ops"

	prevHostname := ""

	ngi := &Group{}
	ngi.Type = "ngi"
	ngi.Name = results[0].Roc

	site := &Group{}
	site.Type = "site"
	site.Name = results[0].Site

	for _, row := range results {
		if row.Hostname != prevHostname && row.Hostname != "" {
			host := &Host{} //create new host
			host.Name = row.Hostname
			site.Hosts = append(site.Hosts, host)
			prevHostname = row.Hostname
		}
	}

	profile.Groups = append(profile.Groups, vo)
	vo.Groups = append(vo.Groups, ngi)
	ngi.Groups = append(ngi.Groups, site)
	docRoot.Profile = profile

	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err

}
