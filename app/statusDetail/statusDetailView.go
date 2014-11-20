/*
 * Copyright (c) 2014 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
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
	prevSite := ""
	prevRoc := ""
	
	ngi := &Group{}
	ngi.Type = "ngi"
	ngi.Name = results[0].Roc

	var pp_Host *Host
	var pp_Metric *Metric
	var pp_Site *Group
	var pp_Roc *Group

	for _, row := range results {

		if row.Roc != prevRoc

		if row.Site != prevSite && row.Site != "" {
			site := &Group{}
			site.Name = row.Site
			site.Type = "site"
			ngi.Groups = append(ngi.Groups, site)
			prevSite = row.Site
			pp_Site = site
		}

		if row.Hostname != prevHostname && row.Hostname != "" {
			host := &Host{} //create new host
			host.Name = row.Hostname
			pp_Site.Hosts = append(pp_Site.Hosts, host)
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
	docRoot.Profile = profile

	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err

}
