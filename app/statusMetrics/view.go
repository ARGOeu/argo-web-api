/*
 * Copyright (c) 2015 GRNET S.A.
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
 * or implied, of GRNET S.A.
 *
 */

package statusMetrics

import (
	"encoding/json"
	"encoding/xml"
	"strings"
)

func createView(results []DataOutput, input InputParams) ([]byte, error) {

	output := []byte("reponse output")
	err := error(nil)

	docRoot := &rootOUT{}

	if len(results) == 0 {
		if strings.EqualFold(input.format, "application/json") {
			output, err = json.MarshalIndent(docRoot, " ", "  ")
		} else {
			output, err = xml.MarshalIndent(docRoot, " ", "  ")
		}
		return output, err
	}

	prevHostname := ""
	prevMetric := ""
	prevEndpointGroup := ""
	prevService := ""

	var ppHost *endpointOUT
	var ppMetric *metricOUT
	var ppEndpointGroup *endpointGroupOUT
	var ppService *serviceOUT

	for _, row := range results {

		if row.EndpointGroup != prevEndpointGroup && row.EndpointGroup != "" {
			endpointGroup := &endpointGroupOUT{}
			endpointGroup.Name = row.EndpointGroup
			endpointGroup.GroupType = input.groupType
			docRoot.EndpointGroups = append(docRoot.EndpointGroups, endpointGroup)
			prevEndpointGroup = row.EndpointGroup
			ppEndpointGroup = endpointGroup
		}

		if row.Service != prevService && row.Service != "" {
			service := &serviceOUT{}
			service.Name = row.Service
			service.GroupType = "service"
			ppEndpointGroup.Services = append(ppEndpointGroup.Services, service)

			prevService = row.Service
			ppService = service
		}

		if row.Hostname != prevHostname && row.Hostname != "" {
			host := &endpointOUT{} //create new host
			host.Name = row.Hostname
			host.GroupType = "endpoint"
			ppService.Endpoints = append(ppService.Endpoints, host)
			prevHostname = row.Hostname
			ppHost = host
		}

		if row.Metric != prevMetric {

			metric := &metricOUT{}
			//Add the prev status as the firstone

			metric.Name = row.Metric
			metric.GroupType = "metric"
			ppHost.Metrics = append(ppHost.Metrics, metric)
			prevMetric = row.Metric
			ppMetric = metric

			prevStatus := &statusOUT{}
			prevStatus.Timestamp = row.PrevTimestamp
			prevStatus.Value = row.PrevStatus
			ppMetric.Statuses = append(ppMetric.Statuses, prevStatus)

		}

		status := &statusOUT{}
		status.Timestamp = row.Timestamp
		status.Value = row.Status
		ppMetric.Statuses = append(ppMetric.Statuses, status)

	}

	if strings.EqualFold(input.format, "application/json") {
		output, err = json.MarshalIndent(docRoot, " ", "  ")
	} else {
		output, err = xml.MarshalIndent(docRoot, " ", "  ")
	}
	return output, err

}

func createMessageOUT(message string, format string) ([]byte, error) {

	output := []byte("message placeholder")
	err := error(nil)
	docRoot := &messageOUT{}

	docRoot.Message = message
	if strings.EqualFold(format, "application/json") {
		output, err = json.MarshalIndent(docRoot, " ", "  ")
	} else {
		output, err = xml.MarshalIndent(docRoot, " ", "  ")
	}
	return output, err
}
