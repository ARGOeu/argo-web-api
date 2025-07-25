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

package statusEndpointGroups

import (
	"encoding/json"
	"encoding/xml"
	"strconv"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/respond"
)

func createView(results []DataOutput, input InputParams, endDate string, details bool) ([]byte, error) {

	// calculate part of the timestamp that closes the timeline of each item
	var extraTS string

	tsNow := time.Now().UTC()
	today := tsNow.Format("2006-01-02")

	if strings.Split(endDate, "T")[0] == today {
		extraTS = "T" + strings.Split(tsNow.Format(zuluForm), "T")[1]
	} else {
		extraTS = "T23:59:59Z"
	}

	var output []byte
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

	prevEndpointGroup := ""

	var ppEndpointGroup *endpointGroupOUT

	for _, row := range results {

		if row.EndpointGroup != prevEndpointGroup && row.EndpointGroup != "" {
			// close the status timeline by adding a new status item at 23:59 or at current time
			if ppEndpointGroup != nil {
				eStatus := &statusOUT{}
				latestStatus := ppEndpointGroup.Statuses[len(ppEndpointGroup.Statuses)-1]
				eStatus.Timestamp = strings.Split(latestStatus.Timestamp, "T")[0] + extraTS
				eStatus.Value = latestStatus.Value
				ppEndpointGroup.Statuses = append(ppEndpointGroup.Statuses, eStatus)
			}

			endpointGroup := &endpointGroupOUT{}
			endpointGroup.Name = row.EndpointGroup
			endpointGroup.GroupType = input.groupType
			docRoot.EndpointGroups = append(docRoot.EndpointGroups, endpointGroup)
			prevEndpointGroup = row.EndpointGroup
			ppEndpointGroup = endpointGroup
		}

		status := &statusOUT{}
		status.Timestamp = row.Timestamp
		status.Value = row.Status
		if details {
			status.AffectedByThresholdRule = row.HasThresholdRule
		}

		ppEndpointGroup.Statuses = append(ppEndpointGroup.Statuses, status)

	}

	// close the status timeline of the last item by adding a new status item at 23:59 or at current time
	if ppEndpointGroup != nil {
		eStatus := &statusOUT{}
		latestStatus := ppEndpointGroup.Statuses[len(ppEndpointGroup.Statuses)-1]
		eStatus.Timestamp = strings.Split(latestStatus.Timestamp, "T")[0] + extraTS
		eStatus.Value = latestStatus.Value
		ppEndpointGroup.Statuses = append(ppEndpointGroup.Statuses, eStatus)
	}

	output, err = respond.MarshalContent(docRoot, input.format, "", " ")
	return output, err

}

func createMessageOUT(message string, code int, format string) ([]byte, error) {

	var output []byte
	err := error(nil)
	docRoot := &messageOUT{}

	docRoot.Message = message
	docRoot.Code = strconv.Itoa(code)
	output, err = respond.MarshalContent(docRoot, format, "", " ")
	return output, err
}
