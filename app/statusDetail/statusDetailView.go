package statusDetail

import "encoding/xml"

//import "fmt"

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
	prevMetric := ""

	ngi := &Group{}
	ngi.Type = "ngi"
	ngi.Name = results[0].Roc

	site := &Group{}
	site.Type = "site"
	site.Name = results[0].Site

	var pp_Host *Host
	var pp_Metric *Metric
	for _, row := range results {

		if row.Hostname != prevHostname && row.Hostname != "" {
			host := &Host{} //create new host
			host.Name = row.Hostname
			site.Hosts = append(site.Hosts, host)
			prevHostname = row.Hostname
			pp_Host = host
		}

		if row.Metric != prevMetric {

			metric := &Metric{}
			metric.Name = row.Metric
			pp_Host.Metrics = append(pp_Host.Metrics, metric)
			prevMetric = row.Metric
			pp_Metric = metric
		}

		status := &Status{}
		status.Timestamp = row.Timestamp
		status.Status = row.Status
		pp_Metric.Timeline = append(pp_Metric.Timeline, status)
	}

	profile.Groups = append(profile.Groups, vo)
	vo.Groups = append(vo.Groups, ngi)
	ngi.Groups = append(ngi.Groups, site)
	docRoot.Profile = profile

	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err

}
