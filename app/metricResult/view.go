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
	"fmt"
	"strings"
)

func createMetricResultView(result metricResultOutput, format string) ([]byte, error) {

	docRoot := &root{}

	// Exit here in case result is empty
	if result.Hostname == "" {
		output, err := xml.MarshalIndent(docRoot, "", "")
		return output, err
	}

	hostname := &HostXML{
		Name: result.Hostname,
	}
	docRoot.Result = append(docRoot.Result, hostname)

	metric := &MetricXML{
		Name: result.Metric,
	}
	hostname.Metrics = append(hostname.Metrics, metric)

	// we append the detailed results
	metric.Details = append(metric.Details,
		&StatusXML{
			Timestamp:    fmt.Sprintf("%s", result.Timestamp),
			Value:        fmt.Sprintf("%s", result.Status),
			Summary:      fmt.Sprintf("%s", result.Summary),
			Message:      fmt.Sprintf("%s", result.Message),
		})

	if strings.ToLower(format) == "application/json" {
		return json.MarshalIndent(docRoot, " ", "  ")
	} else {
		return xml.MarshalIndent(docRoot, " ", "  ")
	}

}


