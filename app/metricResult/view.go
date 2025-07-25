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

package metricResult

import (
	"encoding/json"
	"encoding/xml"
	"strings"
)

func createMultipleMetricResultsView(results []metricResultOutput, format string) ([]byte, error) {

	docRoot := &root{}

	// Exit here in case result is empty
	if len(results) == 0 {

		docRoot.Result = []*HostXML{}
		if strings.ToLower(format) == "application/json" {
			return json.MarshalIndent(docRoot, " ", "  ")
		}
		return xml.MarshalIndent(docRoot, "", "")
	}

	hostname := &HostXML{
		Name: results[0].Hostname,
		Info: results[0].Info,
	}
	docRoot.Result = append(docRoot.Result, hostname)

	prevMetric := ""
	prevService := ""

	metric := &MetricXML{
		Name:    "",
		Service: "",
	}

	for i, result := range results {
		if i == 0 {
			prevMetric = result.Metric
			prevService = result.Service

			metric.Name = result.Metric
			metric.Service = result.Service

		}

		if result.Metric != prevMetric || result.Service != prevService {
			hostname.Metrics = append(hostname.Metrics, metric)

			metric = &MetricXML{
				Name:    result.Metric,
				Service: result.Service,
			}
		}

		// we append the detailed results
		metric.Details = append(metric.Details,
			&StatusXML{
				Timestamp:      result.Timestamp,
				Value:          result.Status,
				Summary:        result.Summary,
				Message:        result.Message,
				ActualData:     result.ActualData,
				RuleApplied:    result.RuleApplied,
				OriginalStatus: result.OriginalStatus,
			})

		prevMetric = result.Metric
		prevService = result.Service

	}

	hostname.Metrics = append(hostname.Metrics, metric)

	if strings.ToLower(format) == "application/json" {
		return json.MarshalIndent(docRoot, " ", "  ")
	}

	return xml.MarshalIndent(docRoot, " ", "  ")

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
