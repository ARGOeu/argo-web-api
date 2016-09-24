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

import "encoding/xml"

type metricResultQuery struct {
	EndpointName string `bson:"hostname"`
	MetricName   string `bson:"metric_name"`
	ExecTime     string `bson:"exec_time"` // UTC time in W3C format
}

// metricResultOutput structure holds mongo results
type metricResultOutput struct {
	Timestamp string `bson:"timestamp"`
	Hostname  string `bson:"host"`
	Metric    string `bson:"metric"`
	Status    string `bson:"status"`
	Summary   string `bson:"summary"`
	Message   string `bson:"message"`
}

// HostXML struct used as xml block
type HostXML struct {
	XMLName xml.Name `xml:"host" json:"-"`
	Name    string   `xml:"name,attr"`
	Metrics []*MetricXML
}

// MetricXML struct used as xml block
type MetricXML struct {
	XMLName xml.Name `xml:"metric" json:"-"`
	Name    string   `xml:"name,attr"`
	Details []*StatusXML
}

// StatusXML struct used as xml block
type StatusXML struct {
	XMLName   xml.Name `xml:"status" json:"-"`
	Timestamp string   `xml:"timestamp,attr"`
	Value     string   `xml:"value,attr"`
	Summary   string   `xml:"summary"`
	Message   string   `xml:"message"`
}

type root struct {
	XMLName xml.Name      `xml:"root" json:"-"`
	Result  []interface{} `json:"root"`
}
