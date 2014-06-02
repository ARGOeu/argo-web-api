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

type RecomputationsInputOutput struct {
	StartTime   string   `bson:"st"`
	EndTime     string   `bson:"et"`
	Reason      string   `bson:"r"`
	NgiName     string   `bson:"n"`
	ExcludeSite []string `bson:"es"`
	Status      string   `bson:"s"`
	Timestamp   string   `bson:"t"`
	//Exclude_sf		[]string
	//Exclude_end_point []string
}

type Exclude struct {
	XMLName xml.Name `xml:"Exclude" json:"-"`
	Site    string   `xml:"site,attr" json:"site"`
}

type Request struct {
	XMLName   xml.Name `xml:"Request" json:"-"`
	StartTime string   `xml:"start_time,attr" json:"start_time"`
	EndTime   string   `xml:"end_time,attr" json:"end_time"`
	Reason    string   `xml:"reason,attr" json:"reason"`
	NgiName   string   `xml:"ngi_name,attr" json:"ngi_name"`
	Status    string   `xml:"status,attr" json:"status"`
	Timestamp string   `xml:"timestamp,attr" json:"timestamp"`
	Exclude   []*Exclude
}

type Root struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Request []*Request
}

func insertQuery(input RecomputationsInputOutput) bson.M {

	query := bson.M{
		"st": input.StartTime,
		"et": input.EndTime,
		"r":  input.Reason,
		"n":  input.NgiName,
		"es": input.ExcludeSite,
		"s":  input.Status,
		"t":  input.Timestamp,
	}

	return query

}
