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

package statusEndpointGroups

import "encoding/xml"

type MongoInterface struct {
	Job            string `bson:"job"         xml:"-" json:"-"`
	SuperGroup     string `bson:"supergroup"           xml:"-" json:"-"`
	EndpointGroup  string `bson:"name"  xml:"-" json:"-"`
	Timestamp      string `bson:"timestamp"       xml:"-" json:"-"`
	Status         string `bson:"status"          xml:"-" json:"-"`
	PreviousStatus string `bson:"previous_status" xml:"-" json:"-"`
	DateInteger    int    `bson:"date_integer"    xml:"-" json:"-"`
	TimeInteger    int    `bson:"time_integer"    xml:"-" json:"-"`
}

type StatusEndpointGroupInput struct {
	Start      string // UTC time in W3C format
	End        string
	Job        string
	Type       string
	Name       string
	SuperGroup string
}
type Root struct {
	XMLName xml.Name `xml:"root"`
	Jobs    []Job
}

type Job struct {
	XMLName       xml.Name        `xml:"Job" json:"-"`
	Name          string          `xml:"name,attr" json:"name" bson:"job"`
	EndpointGroup []EndpointGroup `bson:"endpointgroup"`
}

type EndpointGroup struct {
	XMLName xml.Name `xml:"EndpointGroup"`
	Name    string   `xml:"name,attr" bson:"name"`
	Status  []Status `bson:"statuses"`
}

type Status struct {
	XMLName        xml.Name `xml:"Status"`
	Timestamp      string   `xml:"timestamp,attr" bson:"timestamp"`
	Status         string   `xml:"Status,attr" bson:"status"`
	PreviousStatus string   `xml:"PreviousStatus,attr" bson:"previous_status"`
}
