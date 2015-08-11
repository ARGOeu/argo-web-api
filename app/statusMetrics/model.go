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

package statusMetrics

import "encoding/xml"

// InputParams struct holds as input all the url params of the request
type InputParams struct {
	startTime string // UTC time in W3C format
	endTime   string
	report    string
	groupType string
	group     string
	service   string
	hostname  string
	metric    string
}

// DataOutput struct holds the queried data from datastore
type DataOutput struct {
	Report        string `bson:"report"`
	Timestamp     string `bson:"timestamp"`
	EndpointGroup string `bson:"endpoint_group"`
	Service       string `bson:"service"`
	Hostname      string `bson:"hostname"`
	Metric        string `bson:"metric"`
	Status        string `bson:"status"`
	DateInt       string `bson:"date_int"`
	PrevTimestamp string `bson:"prev_timestamp"`
	PrevStatus    string `bson:"prev_status"`
}

// xml response related structs

type rootXML struct {
	XMLName        xml.Name `xml:"root"`
	EndpointGroups []*endpointGroupXML
}

type endpointGroupXML struct {
	XMLName   xml.Name `xml:"Group"`
	Name      string   `xml:"name,attr"`
	GroupType string   `xml:"type,attr"`
	Services  []*serviceXML
}

type serviceXML struct {
	XMLName   xml.Name `xml:"Group"`
	Name      string   `xml:"name,attr"`
	GroupType string   `xml:"type,attr"`
	Endpoints []*endpointXML
}

type endpointXML struct {
	XMLName   xml.Name `xml:"Group"`
	Name      string   `xml:"name,attr"`
	GroupType string   `xml:"type,attr"`
	Metrics   []*metricXML
}

type metricXML struct {
	XMLName   xml.Name `xml:"Group"`
	Name      string   `xml:"name,attr"`
	GroupType string   `xml:"type,attr"`
	Statuses  []*statusXML
}

type statusXML struct {
	XMLName   xml.Name `xml:"status"`
	Timestamp string   `xml:"timestamp,attr"`
	Value     string   `xml:"value,attr"`
}

// Message struct to hold the xml response
type message struct {
	XMLName xml.Name `xml:"root"`
	Message string
}
