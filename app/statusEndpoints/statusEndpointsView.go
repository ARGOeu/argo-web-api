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

package statusEndpoints

import "encoding/xml"

func createView(results []StatusEndpointsOutput, input StatusEndpointsInput) ([]byte, error) {

	docRoot := &ReadRoot{}

	if len(results) == 0 {
		output, err := xml.MarshalIndent(docRoot, " ", "  ")
		return output, err
	}

	job := &JobXML{}
	job.Name = input.job

	endpoint := &EndpointXML{}
	endpoint.Hostname = input.hostname
	endpoint.Service = input.service_type
	job.Endpoints = append(job.Endpoints, endpoint)

	for _, row := range results {
		if row.Job != input.job {
			continue
		}
		if row.Hostname != input.hostname {
			continue
		}
		if row.Service != input.service_type {
			continue
		}
		status := &StatusXML{}
		status.Timestamp = row.Timestamp
		status.Status = row.Status
		endpoint.Timeline = append(endpoint.Timeline, status)
	}

	docRoot.Job = job

	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err

}
