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

import "time"

type RecomputationsInputOutput struct {
	startTime   string    `bson:"st"`
	endTime     string    `bson:"et"`
	Reason      string    `bson:"r"`
	ngiName     string    `bson:"n"`
	excludeSite []string  `bson:"es"`
	status      string    `bson:"s"`
	timestamp   time.Time `bson:"t"`
	//Exclude_sf		[]string
	//Exclude_end_point []string
}

type Request struct {
	XMLName     xml.Name `xml:"Request" json:"-"`
	startTime   string   `xml:"start_time,attr" json:"start_time"`
	endTime     string   `xml:"end_time,attr" json:"end_time"`
	reason      string   `xml:"reason,attr" json:"reason"`
	ngiName     string   `xml:"ngi_name",attr json:"ngi_name"`
	excludeSite string   `xml:"exclude_site,attr" json:"exclude_site"`
	status      string   `xml:"status, attr" json:"status"`
	timestamp   string   `xml:"timestamp,attr" json:"timestamp"`
}

type Root struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Request []*Profile
}

func insertQuery(input RecomputationsInputOutput) {

	query := bson.M{
		"st": input.Start_time,
		"et": input.End_time,
		"r":  input.Reason,
		"n":  input.Ngi_name,
		"es": input.Exclude_site,
		"s":  input.Status,
		"t":  input.Timestamp,
	}

	return query

}
