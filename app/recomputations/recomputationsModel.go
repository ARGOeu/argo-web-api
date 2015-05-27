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

package recomputations

import (
	"encoding/xml"

	"labix.org/v2/mgo/bson"
)

// RecomputationsInputOutput struct used to retrieve recomputations from db
type RecomputationsInputOutput struct {
	StartTime   string   `bson:"start_time"`
	EndTime     string   `bson:"end_time"`
	Reason      string   `bson:"reason"`
	Group       string   `bson:"group"`
	ExcludeSite []string `bson:"exclude"`
	Status      string   `bson:"status"`
	Timestamp   string   `bson:"timestamp"`
	//Exclude_sf		[]string
	//Exclude_end_point []string
}

// Exclude struct to represent the excluded sites for a recomputation
type Exclude struct {
	XMLName xml.Name `xml:"Exclude" json:"-"`
	Site    string   `xml:"site,attr" json:"site"`
}

// Request struct to represent a request in xml/json
type Request struct {
	XMLName   xml.Name `xml:"Request" json:"-"`
	StartTime string   `xml:"start_time,attr" json:"start_time"`
	EndTime   string   `xml:"end_time,attr" json:"end_time"`
	Reason    string   `xml:"reason,attr" json:"reason"`
	Group     string   `xml:"group,attr" json:"group"`
	Status    string   `xml:"status,attr" json:"status"`
	Timestamp string   `xml:"timestamp,attr" json:"timestamp"`
	Exclude   []*Exclude
}

// Root struct to represent the root of the xml/json document
type Root struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Request []*Request
}

// Message struct to use when outputing an error in xml format
type Message struct {
	XMLName xml.Name `xml:"root"`
	Message string
}

func insertQuery(input RecomputationsInputOutput) bson.M {

	query := bson.M{
		"start_time": input.StartTime,
		"end_time":   input.EndTime,
		"reason":     input.Reason,
		"group":      input.Group,
		"excluded":   input.ExcludeSite,
		"status":     input.Status,
		"timestamp":  input.Timestamp,
	}

	return query
}
