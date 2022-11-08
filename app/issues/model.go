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

package issues

import "encoding/xml"

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

// EndpointData struct holds the queried endpoint data from datastore
type EndpointData struct {
	Timestamp     string            `bson:"timestamp" json:"timestamp"`
	EndpointGroup string            `bson:"endpoint_group" json:"endpoint_group"`
	Service       string            `bson:"service" json:"service"`
	Hostname      string            `bson:"host" json:"endpoint"`
	Status        string            `bson:"status" json:"status"`
	DateInt       string            `bson:"date_integer" json:"-"`
	Info          map[string]string `bson:"info" json:"info,omitempty"`
}

// GroupMetrics hold issues with metrics in a specific group
type GroupMetrics struct {
	Service  string            `bson:"service" json:"service,omitempty"`
	Hostname string            `bson:"host" json:"hostname"`
	Metric   string            `bson:"metric" json:"metric,omitempty"`
	Status   string            `bson:"status" json:"status,omitempty"`
	Info     map[string]string `bson:"info" json:"info,omitempty"`
}

// Message struct to hold the json/xml response
type messageOUT struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Message string   `xml:"message" json:"message"`
	Code    string   `xml:"code,omitempty" json:"code,omitempty"`
}
