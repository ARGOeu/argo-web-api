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
	StartTime string   `xml:"start_time,attr" json:"start_time" bson:"start_time"`
	EndTime   string   `xml:"end_time,attr" json:"end_time" bson:"end_time"`
	Reason    string   `xml:"reason,attr" json:"reason" bson:"reason"`
	Group     string   `xml:"group,attr" json:"group" bson:"group"`
	SubGroups []string `xml:"subgroups" json:"subgroups" bson:"subgroups"`
}

type GroupInRequest struct {
	Type      string   `xml:"type,attr" json:"type"`
	Group     string   `xml:"group,attr" json:"group"`
	SubTypes  string   `xml:"sub_type,attr" json:"subtype"`
	SubGroups []string `xml:"sub_groups" json:"sub_groups"`
}

type MongoInterface struct {
	RequesterName  string   `bson:"requester_name"`
	RequesterEmail string   `bson:"requester_email"`
	Reason         string   `bson:"reason"`
	StartTime      string   `bson:"start_time"`
	EndTime        string   `bson:"end_time"`
	Report         string   `bson:"report"`
	Exclude        []string `bson:"exclude"`
	Status         string   `bson:"status"`
	Timestamp      string   `bson:"timestamp"`
}

type Exclude struct {
	XMLName xml.Name `xml:"subgroup" json:"-"`
	Group   string   `xml:"name,attr" json:"name"`
}
