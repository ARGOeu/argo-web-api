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

package trends

import "encoding/xml"

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "2006-01-02"

// MetricData holds flapping information about metrics
type MetricData struct {
	EndpointGroup string `bson:"group" json:"endpoint_group"`
	Service       string `bson:"service" json:"service"`
	Endpoint      string `bson:"endpoint" json:"endpoint"`
	Metric        string `bson:"metric" json:"metric"`
	Flapping      int    `bson:"flipflop" json:"flapping"`
}

// EndpointData holds flapping information about endpoints
type EndpointData struct {
	EndpointGroup string `bson:"group" json:"endpoint_group"`
	Service       string `bson:"service" json:"service"`
	Endpoint      string `bson:"endpoint" json:"endpoint"`
	Flapping      int    `bson:"flipflop" json:"flapping"`
}

// ServiceData holds flapping information about services
type ServiceData struct {
	EndpointGroup string `bson:"group" json:"endpoint_group"`
	Service       string `bson:"service" json:"service"`
	Flapping      int    `bson:"flipflop" json:"flapping"`
}

// EndpointGroupData holds flapping information about endpoint groups
type EndpointGroupData struct {
	EndpointGroup string `bson:"group" json:"endpoint_group"`
	Flapping      int    `bson:"flipflop" json:"flapping"`
}

// Message struct to hold the json/xml response
type messageOUT struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Message string   `xml:"message" json:"message"`
	Code    string   `xml:"code,omitempty" json:"code,omitempty"`
}
