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

package jobs

import (
	"encoding/xml"

	"labix.org/v2/mgo/bson"
)

// Job structure holds information for a defined computational job
type Job struct {
	XMLName       xml.Name     `bson:",omitempty"      json:"-"               xml:"job"`
	Name          string       `bson:"name"            json:"name"            xml:"name,attr"`
	Tenant        string       `bson:"tenant"          json:"tenant"          xml:"tenant,attr"`
	EndpointGroup string       `bson:"endpoint_group"  json:"endpoint_group"  xml:"endpoint_group,attr"`
	GroupOfGroups string       `bson:"group_of_groups" json:"group_of_groups" xml:"group_of_groups,attr"`
	Profiles      []JobProfile `bson:"profiles"        json:"profiles"        xml:"profiles>profile"`
	FilterTags    []JobTag     `bson:"filter_tags"     json:"filter_tags"     xml:"filter_tags>tag"`
}

// JobProfile holds info about the profiles included in a job definition
type JobProfile struct {
	XMLName xml.Name `bson:",omitempty" json:"-"     xml:"profile"`
	Name    string   `bson:"name"       json:"name"  xml:"name,attr"`
	Value   string   `bson:"value"      json:"value" xml:"value,attr"`
}

// JobTag holds info about the tags used in filtering in a job definition
type JobTag struct {
	XMLName xml.Name `bson:",omitempty" json:"-"     xml:"tag"`
	Name    string   `bson:"name"       json:"name"  xml:"name,attr"`
	Value   string   `bson:"value"      json:"value" xml:"value,attr"`
}

// Message struct for xml message response
type Message struct {
	XMLName xml.Name `xml:"root"`
	Message string
}

// RootXML struct to represent the root of the xml document
type RootXML struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Jobs    *[]Job
}

// createJob is used to create a new job definition
func createJob(input Job) bson.M {
	query := bson.M{
		"name":            input.Name,
		"tenant":          input.Tenant,
		"endpoint_group":  input.EndpointGroup,
		"group_of_groups": input.GroupOfGroups,
		"profiles":        input.Profiles,
		"filter_tags":     input.FilterTags,
	}
	return query
}

// searchName is used to create a simple query object based on name
func searchName(name string) bson.M {
	query := bson.M{
		"name": name,
	}

	return query
}
