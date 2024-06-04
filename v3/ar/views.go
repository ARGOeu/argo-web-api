/*
 * Copyright (c) 2022 GRNET S.A.
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

package ar

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/app/reports"
)

func createGroupResult(results []GroupInterface, report reports.MongoInterface, custom bool) []*SuperGroup {

	result := []*SuperGroup{}

	prevSuperGroup := ""
	prevEndpointGroup := ""
	group := &Group{}
	superGroup := &SuperGroup{}

	// we iterate through the results struct array
	// keeping only the value of each row

	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], fmt.Sprint(row.Date))
		//if new superGroup value does not match the previous superGroup value
		//we create a new superGroup in the xml
		if prevSuperGroup != row.SuperGroup {
			prevSuperGroup = row.SuperGroup
			superGroup = &SuperGroup{
				Name: row.SuperGroup,
				Type: report.GetGroupType(),
			}
			result = append(result, superGroup)
			prevEndpointGroup = ""
		}
		//if new endpointGroup does not match the previous service value
		//we create a new endpointGroup entry in the xml
		if prevEndpointGroup != row.Name {
			prevEndpointGroup = row.Name
			group = &Group{
				Name: row.Name,
				Type: report.GetEndpointGroupType(),
			}
			superGroup.Groups = append(superGroup.Groups, group)
		}
		//we append the new availability values
		// if custom period is selected delete timestamps so no dates appear on a/r results
		prepDate := timestamp.Format(customForm[1])
		if custom {
			prepDate = ""
		}
		group.Availability = append(group.Availability,
			Availability{
				Date:         prepDate,
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability),
				Unknown:      fmt.Sprintf("%g", row.Unknown),
				Uptime:       fmt.Sprintf("%g", row.Up),
				Downtime:     fmt.Sprintf("%g", row.Down),
			})
	}

	return result

}

func createSuperGroupResult(results []SuperGroupInterface, report reports.MongoInterface, custom bool) []*SuperGroup {

	result := []*SuperGroup{}

	prevSuperGroup := ""
	superGroup := &SuperGroup{}

	// we iterate through the results struct array
	// keeping only the value of each row
	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], row.Date)

		//if new superGroup does not match the previous superGroup value
		//we create a new superGroup entry in the xml

		if prevSuperGroup != row.SuperGroup {
			prevSuperGroup = row.SuperGroup
			superGroup = &SuperGroup{
				Name: row.SuperGroup,
				Type: report.GetGroupType(),
			}
			result = append(result, superGroup)
		}

		//we append the new availability values
		// if custom period is selected delete timestamps so no dates appear on a/r results
		prepDate := timestamp.Format(customForm[1])
		if custom {
			prepDate = ""
		}
		superGroup.Results = append(superGroup.Results,
			Availability{
				Date:         prepDate,
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability)})

	}

	return result

}

func createEndpointResult(results []EndpointInterface, id string, custom bool) ([]byte, error) {

	docID := &idOUT{}
	docID.ID = id
	docID.Endpoints = make([]*Endpoint, 0)

	prevEndpoint := ""
	endpoint := &Endpoint{}

	// we iterate through the results struct array
	// keeping only the value of each row
	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], row.Date)

		//if new superGroup does not match the previous superGroup value
		//we create a new superGroup entry in the xml

		if prevEndpoint != row.Name {
			prevEndpoint = row.Name
			endpoint = &Endpoint{
				Name:       row.Name,
				Service:    row.Service,
				Supergroup: row.SuperGroup,
				Info:       row.Info,
			}
			docID.Endpoints = append(docID.Endpoints, endpoint)
		}

		// if custom period is selected delete timestamps so no dates appear on a/r results
		prepDate := timestamp.Format(customForm[1])
		if custom {
			prepDate = ""
		}
		endpoint.Results = append(endpoint.Results,

			Availability{
				Date:         prepDate,
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability),
				Unknown:      fmt.Sprintf("%g", row.Unknown),
				Uptime:       fmt.Sprintf("%g", row.Up),
				Downtime:     fmt.Sprintf("%g", row.Down),
			})

	}

	return json.MarshalIndent(docID, " ", "  ")

}

func createResultView(resultsSuperGroups []SuperGroupInterface, resultsGroups []GroupInterface, report reports.MongoInterface, custom bool) ([]byte, error) {
	docRoot := &root{}

	groupResult := createGroupResult(resultsGroups, report, custom)
	superGroupResult := createSuperGroupResult(resultsSuperGroups, report, custom)

	docRoot.Result = groupResult
	for i := range docRoot.Result {
		sname := docRoot.Result[i].Name
		for j := range superGroupResult {
			if superGroupResult[j].Name == sname {
				docRoot.Result[i].Results = superGroupResult[j].Results
				break
			}
		}
	}

	return json.MarshalIndent(docRoot, " ", "  ")

}

func createErrorMessage(message string, code int, format string) ([]byte, error) {

	var output []byte
	err := error(nil)
	docRoot := &errorMessage{}

	docRoot.Message = message
	docRoot.Code = code
	if strings.EqualFold(format, "application/json") {
		output, err = json.MarshalIndent(docRoot, " ", "  ")
	} else {
		output, err = xml.MarshalIndent(docRoot, " ", "  ")
	}
	return output, err
}
