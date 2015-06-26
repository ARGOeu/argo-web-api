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

package statusMsg

import (
	"encoding/xml"

	"github.com/argoeu/argo-web-api/app/metricProfiles"
)

func createView(results []MsgOutput, input MsgInput, metricDetail metricProfiles.MongoInterface) ([]byte, error) {

	docRoot := &ReadRoot{}

	if len(results) == 0 {
		output, err := xml.MarshalIndent(docRoot, " ", "  ")
		return output, err
	}

	job := &JobXML{}
	job.Name = input.job

	prevHostname := ""
	prevMetric := ""
	prevEndpointGroup := ""
	prevGroup := ""
	prevService := ""

	var ppHost *HostXML
	var ppMetric *MetricXML
	var ppEndpointGroup *GroupXML
	var ppGroup *GroupXML
	var ppService *GroupXML

	for _, row := range results {

		// Filter row of metric result based on metric profile (check metric name and service type)
		if filterByProfile(row.Service, row.Metric, metricDetail) == 1 {

			continue
		}

		if row.Group != prevGroup && row.Group != "" {
			group := &GroupXML{}
			group.Name = row.Group
			group.Type = row.GroupType
			job.Groups = append(job.Groups, group)
			prevGroup = group.Name
			ppGroup = group
		}

		if row.EndpointGroup != prevEndpointGroup && row.EndpointGroup != "" {
			endpointGroup := &GroupXML{}
			endpointGroup.Name = row.EndpointGroup
			endpointGroup.Type = row.EndpointGroupType
			ppGroup.Groups = append(ppGroup.Groups, endpointGroup)
			prevEndpointGroup = row.EndpointGroup
			ppEndpointGroup = endpointGroup
		}

		if row.Service != prevService && row.Service != "" {
			service := &GroupXML{}
			service.Name = row.Service
			service.Type = "service_type"
			ppEndpointGroup.Groups = append(ppEndpointGroup.Groups, service)

			prevService = row.Service
			ppService = service
		}

		if row.Hostname != prevHostname && row.Hostname != "" {
			host := &HostXML{} //create new host
			host.Name = row.Hostname
			ppService.Hosts = append(ppService.Hosts, host)
			prevHostname = row.Hostname
			ppHost = host
		}

		if row.Metric != prevMetric {

			metric := &MetricXML{}

			metric.Name = row.Metric
			ppHost.Metrics = append(ppHost.Metrics, metric)
			prevMetric = row.Metric
			ppMetric = metric

			status := &StatusXML{}
			status.Timestamp = input.execTime
			status.Summary = row.Summary
			status.Message = row.Message
			ppMetric.Timeline = append(ppMetric.Timeline, status)

		} else {
			status := &StatusXML{}
			status.Timestamp = row.Timestamp
			status.Status = row.Status
			status.Summary = row.Summary
			status.Message = row.Message
			ppMetric.Timeline = append(ppMetric.Timeline, status)
		}

	}

	docRoot.Job = job

	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err

}

func messageXML(answer string) ([]byte, error) {
	docRoot := &Message{}
	docRoot.Message = answer
	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err
}

func filterByProfile(service string, metric string, metricDetail metricProfiles.MongoInterface) int {

	for _, serviceItem := range metricDetail.Services {

		if serviceItem.Service == service {

			for _, metricItem := range serviceItem.Metrics {

				if metricItem == metric {
					return 0
				}
			}
		}
	}

	return 1

}
