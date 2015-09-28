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

package recomputations2

import "encoding/xml"

type IncomingRequest struct {
	Data []IncomingRecomputation `xml:"data" json:"data"`
}

type IncomingRecomputation struct {
	StartTime string   `xml:"start_time,attr" json:"start_time" bson:"start_time,omitempty"`
	EndTime   string   `xml:"end_time,attr" json:"end_time" bson:"end_time,omitempty"`
	Reason    string   `xml:"reason,attr" json:"reason" bson:"reason,omitempty"`
	Report    string   `xml:"report,attr" json:"report" bson:"report,omitempty"`
	Exclude   []string `xml:"excluded" json:"excluded" bson:"excluded,omitempty"`
}

type MongoInterface struct {
	RequesterName  string   `bson:"requester_name" xml:"requester_name" json:"requester_name"`
	RequesterEmail string   `bson:"requester_email" xml:"requester_email" json:"requester_email"`
	Reason         string   `bson:"reason" xml:"reason" json:"reason"`
	StartTime      string   `bson:"start_time" xml:"start_time" json:"start_time"`
	EndTime        string   `bson:"end_time" xml:"end_time" json:"end_time"`
	Report         string   `bson:"report" xml:"report" json:"report"`
	Exclude        []string `bson:"exclude" xml:"exclude" json:"exclude"`
	Status         string   `bson:"status" xml:"status" json:"status"`
	Timestamp      string   `bson:"timestamp" xml:"timestamp" json:"timestamp"`
}

type Exclude struct {
	XMLName xml.Name `xml:"exclude" json:"-"`
	Group   string   `xml:"name,attr" json:"name"`
}

type root struct {
	XMLName xml.Name    `xml:"root" json:"-"`
	Result  interface{} `json:"root"`
}

type Message struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Message string   `json:"message" xml:"message"`
	Status  string   `json:"status" xml:"status"`
}

// errorMessage struct to hold the json/xml error response
type errorMessage struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Message string   `xml:"message" json:"message"`
}
