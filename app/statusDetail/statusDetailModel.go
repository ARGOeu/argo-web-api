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

package statusDetail

import "encoding/xml"

// InputParams struct holds as input all the url params of the request
type InputParams struct {
	startTime string // UTC time in W3C format
	endTime   string
	vo        string
	profile   string
	groupType string
	group     string
}

// DataOutput struct holds the queried data from datastore
type DataOutput struct {
	Timestamp  string `bson:"ts"`
	Roc        string `bson:"roc"`
	Site       string `bson:"site"`
	Service    string `bson:"srv"`
	Hostname   string `bson:"h"`
	Metric     string `bson:"m"`
	Status     string `bson:"s"`
	TimeInt    int    `bson:"ti"`
	PStatus    string `bson:"ps"`
	PTimestamp string `bson:"pts"`
}

// MetricDetailOutput struct holds metric profile data
// from secondary collection
type MetricDetailOutput struct {
	Service string `bson:"s"`
	Metric  string `bson:"m"`
}

// ReadRoot struct used as xml block
type ReadRoot struct {
	XMLName xml.Name `xml:"root"`
	Profile *Profile
}

// Profile struct used as xml block
type Profile struct {
	XMLName xml.Name `xml:"profile"`
	Name    string   `xml:"name,attr"`
	Groups  []*Group
}

// Group struct used as xml block
type Group struct {
	XMLName xml.Name `xml:"group"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Groups  []*Group
	Hosts   []*Host
}

// Host struct used as xml block
type Host struct {
	XMLName xml.Name `xml:"host"`
	Name    string   `xml:"name,attr"`
	Metrics []*Metric
}

// Metric struct used as xml block
type Metric struct {
	XMLName  xml.Name `xml:"metric"`
	Name     string   `xml:"name,attr"`
	Timeline []*Status
}

// Status struct used as xml block
type Status struct {
	XMLName   xml.Name `xml:"status"`
	Timestamp string   `xml:"timestamp,attr"`
	Status    string   `xml:"status,attr"`
}
