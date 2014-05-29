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

package serviceFlavorAvailability

import (
	"encoding/xml"
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

// a series of auxiliary structs that will
// help us form the xml response
type Availability struct {
	XMLName      xml.Name `xml:"Availability" json:"-"`
	Timestamp    string   `xml:"timestamp,attr", json:"timestamp"`
	Availability string   `xml:"availability,attr" json:"availability"`
	Reliability  string   `xml:"reliability,attr" json:"reliability"`
}

type SF struct {
	XMLName      xml.Name `xml:"Flavor" json:"-"`
	SF           string   `xml:"Flavor,attr" json:"Flavor"`
	Availability []*Availability
}

type Site struct {
	XMLName xml.Name `xml:"Site" json:"-"`
	Site    string   `xml:"Site,attr" json:"Site"`
	SF      []*SF
}

type Profile struct {
	XMLName xml.Name `xml:"Profile" json:"-"`
	Name    string   `xml:"name,attr" json:"name"`
	Site    []*Site
}

type Root struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Profile []*Profile
}

type ApiSFAvailabilityInProfileInput struct {
	// mandatory values
	start_time  string // UTC time in W3C format
	end_time    string // UTC time in W3C format
	profile     string
	granularity string // availability period; possible values: `HOURLY`, `DAILY`, `WEEKLY`, `MONTHLY`
	format      string
	flavor      []string // sf name; may appear more than once
	site        []string // egi site
}

type ApiSFAvailabilityInProfileOutput struct {
	Date         string  "dt"
	SF           string  "sf"
	Site         string  "s"
	Profile      string  "p"
	Availability float64 "a"
	Reliability  float64 "r"
}

type list []interface{}

var customForm []string

func init() {
	customForm = []string{"20060102", "2006-01-02"} //{"Format that is returned by the database" , "Format that will be used in the generated report"}
}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

func prepareFilter(input ApiSFAvailabilityInProfileInput) bson.M {

	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	filter := bson.M{
		"p":  input.profile,
		"dt": bson.M{"$gte": tsYMD, "$lte": teYMD},
	}

	if len(input.flavor) > 0 {
		filter["sf"] = bson.M{"$in": input.flavor}
	}

	if len(input.site) > 0 {
		filter["s"] = bson.M{"$in": input.site}
	}

	return filter
}

func Daily(input ApiSFAvailabilityInProfileInput) []bson.M {

	filter := prepareFilter(input)

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "sf": "$sf", "s": "$s", "a": "$a", "r": "$r", "p": "$p"}}},
		{"$project": bson.M{"dt": "$_id.dt", "sf": "$_id.sf", "a": "$_id.a", "r": "$_id.r", "s": "$_id.s", "p": "$_id.p"}},
		{"$sort": bson.D{{"s", 1}, {"sf", 1}, {"dt", 1}}}}

	return query
}

func Monthly(input ApiSFAvailabilityInProfileInput) []bson.M {

	filter := prepareFilter(input)

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 6}}}, "s": "$s", "p": "$p", "sf": "$sf"}, "avgup": bson.M{"$avg": "$up"}, "avgu": bson.M{"$avg": "$u"}, "avgd": bson.M{"$avg": "$d"}}},
		{"$project": bson.M{"dt": "$_id.dt", "sf": "$_id.sf", "s": "$_id.s", "p": "$_id.p", "a": bson.M{"$multiply": list{bson.M{"$divide": list{"$avgup", bson.M{"$subtract": list{1.00000001, "$avgu"}}}}, 100}},
			"r": bson.M{"$multiply": list{bson.M{"$divide": list{"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgu"}}, "$avgd"}}}}, 100}}}},
		{"$sort": bson.D{{"s", 1}, {"sf", 1}, {"dt", 1}}}}

	return query
}
