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

package statusEndpoints

import "encoding/xml"

type StatusEndpointsInput struct {
	start_time   string // UTC time in W3C format
	end_time     string
	job          string
	hostname     string
	service_type string
}

type StatusEndpointsOutput struct {
	Timestamp          string `bson:"timestamp"`
	Service            string `bson:"service"`
	Hostname           string `bson:"hostname"`
	Status             string `bson:"status"`
	Time_int           int    `bson:"time_int"`
	Job                string `bson:"job"`
}

type ReadRoot struct {
	XMLName xml.Name `xml:"root"`
	Job     *JobXML
}

type JobXML struct {
	XMLName   xml.Name `xml:"job"`
	Name      string   `xml:"name,attr"`
	Endpoints []*EndpointXML
}

type EndpointXML struct {
	XMLName  xml.Name `xml:"endpoint"`
	Hostname string   `xml:"hostname,attr"`
	Service  string   `xml:"service,attr"`
	Timeline []*StatusXML
}

type StatusXML struct {
	XMLName   xml.Name `xml:"status"`
	Timestamp string   `xml:"timestamp,attr"`
	Status    string   `xml:"status,attr"`
}
