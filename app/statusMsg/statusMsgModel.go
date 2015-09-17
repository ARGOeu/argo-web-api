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

package statusMsg

import "encoding/xml"

// MsgInput structure holds url input params
type MsgInput struct {
	execTime string // UTC time in W3C format
	report   string
	host     string
	service  string
	metric   string
}

// MsgOutput structure holds mongo results
type MsgOutput struct {
	Timestamp         string `bson:"timestamp"`
	Group             string `bson:"supergroup"`
	EndpointGroup     string `bson:"endpoint_group"`
	GroupType         string `bson:"group_type"`
	EndpointGroupType string `bson:"endpoint_group_type"`
	Service           string `bson:"service"`
	Hostname          string `bson:"hostname"`
	Metric            string `bson:"metric"`
	Status            string `bson:"status"`
	TimeInt           int    `bson:"time_int"`
	PrevStatus        string `bson:"prev_status"`
	Summary           string `bson:"summary"`
	Message           string `bson:"message"`
}

// ReadRoot struct used as xml block
type ReadRoot struct {
	XMLName xml.Name `xml:"root"`
	Report  *ReportXML
}

// ReportXML struct used as xml block
type ReportXML struct {
	XMLName xml.Name `xml:"report"`
	Name    string   `xml:"name,attr"`
	Groups  []*GroupXML
}

// GroupXML struct used as xml block
type GroupXML struct {
	XMLName xml.Name `xml:"group"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Groups  []*GroupXML
	Hosts   []*HostXML
}

// HostXML struct used as xml block
type HostXML struct {
	XMLName xml.Name `xml:"endpoint"`
	Name    string   `xml:"name,attr"`
	Metrics []*MetricXML
}

// MetricXML struct used as xml block
type MetricXML struct {
	XMLName  xml.Name `xml:"metric"`
	Name     string   `xml:"name,attr"`
	Timeline []*StatusXML
}

// StatusXML struct used as xml block
type StatusXML struct {
	XMLName   xml.Name `xml:"status"`
	Timestamp string   `xml:"timestamp,attr"`
	Status    string   `xml:"status,attr"`
	Summary   string   `xml:"summary"`
	Message   string   `xml:"message"`
}

// Message struct to hold the xml response
type Message struct {
	XMLName xml.Name `xml:"root"`
	Message string
}
