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

package voAvailability

import (
	"encoding/xml"
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

type Availability struct {
	XMLName      xml.Name `xml:"Availability" json:"-"`
	Timestamp    string   `xml:"timestamp,attr" json:"timestamp"`
	Availability string   `xml:"availability,attr" json:"availability"`
	Reliability  string   `xml:"reliability,attr" json:"reliability"`
}

type Vo struct {
	XMLName      xml.Name `xml:"Vo" json:"-"`
	Vo           string   `xml:"VO,attr" json:"VO"`
	Availability []*Availability
}

type Profile struct {
	XMLName xml.Name `xml:"Profile" json:"-"`
	Name    string   `xml:"name,attr" json:"name"`
	Vo      []*Vo
}

type Root struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Profile []*Profile
}

type ApiVoAvailabilityInProfileInput struct {
	// mandatory values
	start_time           string // UTC time in W3C format
	end_time             string // UTC time in W3C format
	availability_profile string //availability profile
	granularity          string // availability period; possible values: `DAILY`  `MONTHLY`
	// optional values
	format     string   // default XML; possible values are: XML, JSON
	group_name []string // site name; may appear more than once
}

type ApiVoAvailabilityInProfileOutput struct {
	Date         string  `bson:"dt"`
	Profile      string  `bson:"ap"`
	Vo           string  `bson:"v"`
	Availability float64 `bson:"a"`
	Reliability  float64 `bson:"r"`
}

type list []interface{}

var customForm []string

func init() {
	customForm = []string{"20060102", "2006-01-02"} //{"Format that is returned by the database" , "Format that will be used in the generated report"}
}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

func prepareFilter(input ApiVoAvailabilityInProfileInput) bson.M {

	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	filter := bson.M{
		"ap": input.availability_profile,
		"dt": bson.M{"$gte": tsYMD, "$lte": teYMD},
	}

	if len(input.group_name) > 0 {
		filter["v"] = bson.M{"$in": input.group_name}
	}

	return filter
}

func Daily(input ApiVoAvailabilityInProfileInput) []bson.M {

	filter := prepareFilter(input)

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "ap": "$ap", "v": "$v", "a": "$a", "r": "$r"}}},
		{"$project": bson.M{"dt": "$_id.dt", "v": "$_id.v", "ap": "$_id.ap", "a": "$_id.a", "r": "$_id.r"}},
		{"$sort": bson.D{{"ap", 1}, {"v", 1}, {"dt", 1}}}}

	return query
}

func Monthly(input ApiVoAvailabilityInProfileInput) []bson.M {
	filter := prepareFilter(input)

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 6}}}, "ap": "$ap", "v": "$v"},
			"avgup": bson.M{"$avg": "$up"}, "avgu": bson.M{"$avg": "$u"}, "avgd": bson.M{"$avg": "$d"}}},
		{"$project": bson.M{"dt": "$_id.dt", "v": "$_id.v", "ap": "$_id.ap",
			"a": bson.M{"$multiply": list{bson.M{"$divide": list{"$avgup", bson.M{"$subtract": list{1.00000001, "$avgu"}}}}, 100}},
			"r": bson.M{"$multiply": list{bson.M{"$divide": list{"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgu"}}, "$avgd"}}}}, 100}}}},
		{"$sort": bson.D{{"ap", 1}, {"v", 1}, {"dt", 1}}}}

	return query
}
