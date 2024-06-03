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

package statusFlatMetrics

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/ARGOeu/argo-web-api/respond"
)

func createFlatView(results []DataOutput, input InputParams, limit int, skip int, details bool) ([]byte, error) {

	var output []byte
	err := error(nil)

	docRoot := &rootPagedOUT{}

	if len(results) == 0 {
		var elist = make([]*endpointOUT, 0)
		docRoot.Endpoints = elist
		if strings.EqualFold(input.format, "application/json") {
			output, err = json.MarshalIndent(docRoot, " ", "  ")
		} else {
			output, err = xml.MarshalIndent(docRoot, " ", "  ")
		}
		return output, err
	}

	prevHostname := ""
	prevEndpointGroup := ""
	prevService := ""

	var ppHost *endpointOUT

	endloop := len(results)

	if len(results) > limit && limit > 0 {
		endloop = len(results) - 1

	}

	for i := 0; i < endloop; i++ {

		row := results[i]

		if (row.Hostname != prevHostname && row.Hostname != "") ||
			(row.EndpointGroup != prevEndpointGroup && row.EndpointGroup != "") ||
			(row.Service != prevService && row.Service != "") {

			host := &endpointOUT{} //create new host
			host.Name = row.Hostname
			host.Info = row.Info
			host.Service = row.Service
			host.SuperGroup = row.EndpointGroup
			host.Metric = row.Metric
			docRoot.Endpoints = append(docRoot.Endpoints, host)
			prevHostname = row.Hostname
			prevEndpointGroup = row.EndpointGroup
			prevService = row.Service
			ppHost = host

			if row.Timestamp[0:10] != row.PrevTimestamp[0:10] {
				prevStatus := &statusOUT{}
				prevStatus.Timestamp = row.PrevTimestamp
				prevStatus.Value = row.PrevStatus
				ppHost.Statuses = append(ppHost.Statuses, prevStatus)
			}

		}

		status := &statusOUT{}
		status.Timestamp = row.Timestamp
		status.Value = row.Status
		if details {
			status.ActualData = row.ActualData
			status.OriginalStatus = row.OriginalStatus
			status.RuleApplied = row.RuleApplied
		}
		ppHost.Statuses = append(ppHost.Statuses, status)

	}

	if limit > 0 {
		if len(results) > limit {
			docRoot.PageToken = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(skip + limit)))

		}
		docRoot.PageSize = limit
	}

	output, err = respond.MarshalContent(docRoot, input.format, "", " ")
	return output, err

}
