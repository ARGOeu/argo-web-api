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

func createView(results []DataOutput, input InputParams, metricDetail []MetricDetailOutput) ([]byte, error) {

	docRoot := &ReadRoot{}

	if len(results) == 0 {
		output, err := xml.MarshalIndent(docRoot, " ", "  ")
		return output, err
	}

	profile := &Profile{}
	profile.Name = input.profile
	vo := &Group{}
	vo.Type = "vo"
	vo.Name = input.vo

	prevHostname := ""
	prevMetric := ""
	prevSite := ""
	prevRoc := ""
	prevService := ""

	var ppHost *Host
	var ppMetric *Metric
	var ppSite *Group
	var ppRoc *Group
	var ppService *Group

	for _, row := range results {

		if filterByProfile(row.Service, row.Metric, metricDetail) == 1 {
			continue
		}

		if row.Roc != prevRoc && row.Roc != "" {
			roc := &Group{}
			roc.Name = row.Roc
			roc.Type = "ngi"
			vo.Groups = append(vo.Groups, roc)
			prevRoc = roc.Name
			ppRoc = roc
		}

		if row.Site != prevSite && row.Site != "" {
			site := &Group{}
			site.Name = row.Site
			site.Type = "site"
			ppRoc.Groups = append(ppRoc.Groups, site)
			prevSite = row.Site
			ppSite = site
		}

		if row.Service != prevService && row.Service != "" {
			service := &Group{}
			service.Name = row.Service
			service.Type = "service_type"
			ppSite.Groups = append(ppSite.Groups, service)

			prevService = row.Service
			ppService = service
		}

		if row.Hostname != prevHostname && row.Hostname != "" {
			host := &Host{} //create new host
			host.Name = row.Hostname
			ppService.Hosts = append(ppService.Hosts, host)
			prevHostname = row.Hostname
			ppHost = host
		}

		if row.Metric != prevMetric {

			metric := &Metric{}
			//Add the prev status as the firstone

			metric.Name = row.Metric
			ppHost.Metrics = append(ppHost.Metrics, metric)
			prevMetric = row.Metric
			ppMetric = metric

			status := &Status{}
			status.Timestamp = row.PTimestamp
			status.Status = row.PStatus
			ppMetric.Timeline = append(ppMetric.Timeline, status)

		} else {
			status := &Status{}
			status.Timestamp = row.Timestamp
			status.Status = row.Status
			ppMetric.Timeline = append(ppMetric.Timeline, status)
		}

	}

	profile.Groups = append(profile.Groups, vo)
	docRoot.Profile = profile

	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err

}

func filterByProfile(stype string, metric string, metricDetail []MetricDetailOutput) int {

	for _, item := range metricDetail {

		if item.Metric == metric {
			if item.Service == stype {
				return 0
			}
		}
	}

	return 1

}
