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

package results

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/app/reports"
)

func createEndpointResultView(results []EndpointInterface, report reports.MongoInterface, format string, custom bool) ([]byte, error) {
	docRoot := &root{}

	prevServiceFlavorGroup := ""
	prevServiceFlavor := ""
	prevEndpoint := ""
	serviceFlavorGroup := &ServiceFlavorGroup{}
	serviceEndpointGroup := &ServiceEndpointGroup{}
	endpoint := &Endpoint{}

	// we iterate through the results struct array
	// keeping only the value of each row
	for _, row := range results {

		timestamp, _ := time.Parse(customForm[0], fmt.Sprint(row.Date))
		//if new superGroup value does not match the previous superGroup value
		//we create a new superGroup in the xml
		if prevServiceFlavorGroup != row.SuperGroup {
			prevServiceFlavorGroup = row.SuperGroup
			serviceFlavorGroup = &ServiceFlavorGroup{
				Name: row.SuperGroup,
				Type: report.GetEndpointGroupType(), // Endpoint groups are parents of SFs
			}
			docRoot.Result = append(docRoot.Result, serviceFlavorGroup)
			prevServiceFlavor = ""
		}
		//if new service flavor does not match the previous service value
		//we create a new service flavor entry in the xml/json output
		if prevServiceFlavor != row.Service {
			prevServiceFlavor = row.Service
			serviceEndpointGroup = &ServiceEndpointGroup{
				Name: row.Service,
				Type: fmt.Sprintf("service"),
			}
			serviceFlavorGroup.ServiceFlavor = append(serviceFlavorGroup.ServiceFlavor, serviceEndpointGroup)
			prevEndpoint = ""
		}

		if prevEndpoint != row.Name {
			prevEndpoint = row.Name
			endpoint = &Endpoint{
				Name: row.Name,
				Type: fmt.Sprintf("endpoint"),
				Info: row.Info,
			}
			serviceEndpointGroup.Endpoints = append(serviceEndpointGroup.Endpoints, endpoint)
		}
		//we append the new availability values
		prepDate := timestamp.Format(customForm[1])
		if custom {
			prepDate = ""
		}
		endpoint.Availability = append(endpoint.Availability,
			&Availability{
				Timestamp:    prepDate,
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability),
				Unknown:      fmt.Sprintf("%g", row.Unknown),
				Uptime:       fmt.Sprintf("%g", row.Up),
				Downtime:     fmt.Sprintf("%g", row.Down),
			})
	}

	if strings.ToLower(format) == "application/json" {
		return json.MarshalIndent(docRoot, " ", "  ")
	}
	return xml.MarshalIndent(docRoot, " ", "  ")

}

func createFlatEndpointResultView(results []EndpointInterface, report reports.MongoInterface, format string, limit int, skip int, custom bool) ([]byte, error) {

	docRoot := &pageRoot{}

	prevEndpoint := ""
	prevService := ""
	prevSuperGroup := ""
	endpoint := &Endpoint{}

	endloop := len(results)

	if len(results) > limit && limit > 0 {
		endloop = len(results) - 1

	}

	for i := 0; i < endloop; i++ {
		row := results[i]

		timestamp, _ := time.Parse(customForm[0], fmt.Sprint(row.Date))

		if prevEndpoint != row.Name || prevService != row.Service || prevSuperGroup != row.SuperGroup {
			prevEndpoint = row.Name
			prevService = row.Service
			prevSuperGroup = row.SuperGroup
			endpoint = &Endpoint{
				Name:       row.Name,
				Type:       fmt.Sprintf("endpoint"),
				Service:    row.Service,
				SuperGroup: row.SuperGroup,
				Info:       row.Info,
			}
			docRoot.Result = append(docRoot.Result, endpoint)

		}
		//we append the new availability values
		prepDate := timestamp.Format(customForm[1])
		if custom {
			prepDate = ""
		}
		endpoint.Availability = append(endpoint.Availability,
			&Availability{
				Timestamp:    prepDate,
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability),
				Unknown:      fmt.Sprintf("%g", row.Unknown),
				Uptime:       fmt.Sprintf("%g", row.Up),
				Downtime:     fmt.Sprintf("%g", row.Down),
			})

	}

	//docRoot.Result = docRoot.Result[:len(docRoot.Result)-1]

	if limit > 0 {
		if len(results) > limit {
			docRoot.PageToken = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(skip + limit)))
		}
		docRoot.PageSize = limit
	}

	if strings.ToLower(format) == "application/json" {
		return json.MarshalIndent(docRoot, " ", "  ")
	}
	return xml.MarshalIndent(docRoot, " ", "  ")

}

func createServiceFlavorResultView(results []ServiceFlavorInterface, report reports.MongoInterface, format string, custom bool) ([]byte, error) {

	docRoot := &root{}

	prevServiceFlavorGroup := ""
	prevServiceFlavor := ""
	serviceFlavor := &ServiceFlavor{}
	serviceFlavorGroup := &ServiceFlavorGroup{}

	// we iterate through the results struct array
	// keeping only the value of each row
	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], fmt.Sprint(row.Date))
		//if new superGroup value does not match the previous superGroup value
		//we create a new superGroup in the xml
		if prevServiceFlavorGroup != row.SuperGroup {
			prevServiceFlavorGroup = row.SuperGroup
			serviceFlavorGroup = &ServiceFlavorGroup{
				Name: row.SuperGroup,
				Type: report.GetEndpointGroupType(), // Endpoint groups are parents of SFs
			}
			docRoot.Result = append(docRoot.Result, serviceFlavorGroup)
			prevServiceFlavor = ""
		}
		//if new service flavor does not match the previous service value
		//we create a new service flavor entry in the xml/json output
		if prevServiceFlavor != row.Name {
			prevServiceFlavor = row.Name
			serviceFlavor = &ServiceFlavor{
				Name: row.Name,
				Type: fmt.Sprintf("service"),
			}
			serviceFlavorGroup.ServiceFlavor = append(serviceFlavorGroup.ServiceFlavor, serviceFlavor)
		}
		//we append the new availability values
		prepDate := timestamp.Format(customForm[1])
		if custom {
			prepDate = ""
		}
		serviceFlavor.Availability = append(serviceFlavor.Availability,
			&Availability{
				Timestamp:    prepDate,
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability),
				Unknown:      fmt.Sprintf("%g", row.Unknown),
				Uptime:       fmt.Sprintf("%g", row.Up),
				Downtime:     fmt.Sprintf("%g", row.Down),
			})
	}

	if strings.ToLower(format) == "application/json" {
		return json.MarshalIndent(docRoot, " ", "  ")
	}
	return xml.MarshalIndent(docRoot, " ", "  ")

}

func createEndpointGroupResultView(results []EndpointGroupInterface, report reports.MongoInterface, format string, custom bool) ([]byte, error) {

	docRoot := &root{}

	prevSuperGroup := ""
	prevEndpointGroup := ""
	endpointGroup := &Group{}
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
			docRoot.Result = append(docRoot.Result, superGroup)
			prevEndpointGroup = ""
		}
		//if new endpointGroup does not match the previous service value
		//we create a new endpointGroup entry in the xml
		if prevEndpointGroup != row.Name {
			prevEndpointGroup = row.Name
			endpointGroup = &Group{
				Name: row.Name,
				Type: report.GetEndpointGroupType(),
			}
			superGroup.Endpoints = append(superGroup.Endpoints, endpointGroup)
		}
		//we append the new availability values
		prepDate := timestamp.Format(customForm[1])
		if custom {
			prepDate = ""
		}
		endpointGroup.Availability = append(endpointGroup.Availability,
			&Availability{
				Timestamp:    prepDate,
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability),
				Unknown:      fmt.Sprintf("%g", row.Unknown),
				Uptime:       fmt.Sprintf("%g", row.Up),
				Downtime:     fmt.Sprintf("%g", row.Down),
			})
	}
	if strings.ToLower(format) == "application/json" {
		return json.MarshalIndent(docRoot, " ", "  ")
	}
	return xml.MarshalIndent(docRoot, " ", "  ")

}

func createSuperGroupView(results []SuperGroupInterface, report reports.MongoInterface, format string, custom bool) ([]byte, error) {

	docRoot := &root{}

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
			docRoot.Result = append(docRoot.Result, superGroup)
		}
		//we append the new availability values
		prepDate := timestamp.Format(customForm[1])
		if custom {
			prepDate = ""
		}
		superGroup.Results = append(superGroup.Results,
			&Availability{
				Timestamp:    prepDate,
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability)})
	}

	if strings.ToLower(format) == "application/json" {
		return json.MarshalIndent(docRoot, " ", "  ")
	}
	return xml.MarshalIndent(docRoot, " ", "  ")

}

func createErrorMessage(message string, code int, format string) ([]byte, error) {

	output := []byte("message placeholder")
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
