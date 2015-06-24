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

type SuperGroup struct {
	XMLName    xml.Name `xml:"SuperGroup" json:"-"`
	SuperGroup string   `xml:"name,attr"  json:"name"`
	SF         []*SF
}

type Job struct{
	XMLName     xml.Name `xml:"Job" json:"-"`
	Name        string   `xml:"name,attr" json:"name"`
	SuperGroup  []*SuperGroup
}

type Root struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Job     []*Job
}

type ApiSFAvailabilityInProfileInput struct {
	// mandatory values
	start_time     string // UTC time in W3C format
	end_time       string // UTC time in W3C format
	job            string
	granularity    string // availability period; possible values: `DAILY`, `MONTHLY`
	format         string
	flavor         []string // sf name; may appear more than once
	supergroup     []string // name of group
}

type ApiSFAvailabilityInProfileOutput struct {
	Date         string  `bson:"date"`
	SF           string  `bson:"name"`
	SuperGroup   string  `bson:"supergroup"`
	Job          string  `bson:"job"`
	Availability float64 `bson:"availability"`
	Reliability  float64 `bson:"reliability"`
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
		"job":  input.job,
		"date": bson.M{"$gte": tsYMD, "$lte": teYMD},
	}

	if len(input.flavor) > 0 {
		filter["name"] = bson.M{"$in": input.flavor}
	}

	if len(input.supergroup) > 0 {
		filter["supergroup"] = bson.M{"$in": input.supergroup}
	}

	return filter
}

func Daily(input ApiSFAvailabilityInProfileInput) []bson.M {

	filter := prepareFilter(input)

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": bson.M{"date": bson.D{{"$substr", list{"$date", 0, 8}}}, "name": "$name", "supergroup": "$supergroup", "availability": "$availability", "reliability": "$reliability", "job": "$job"}}},
		{"$project": bson.M{"date": "$_id.date", "name": "$_id.name", "availability": "$_id.availability", "reliability": "$_id.reliability", "supergroup": "$_id.supergroup", "job": "$_id.job"}},
		{"$sort": bson.D{{"supergroup", 1}, {"name", 1}, {"date", 1}}}}

	return query
}

func Monthly(input ApiSFAvailabilityInProfileInput) []bson.M {

	filter := prepareFilter(input)

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": bson.M{"date": bson.D{{"$substr", list{"$date", 0, 6}}}, "name": "$name", "supergroup": "$supergroup", "job": "$job"}, "avgup": bson.M{"$avg": "$up"}, "avgunknown": bson.M{"$avg": "$unknown"}, "avgdown": bson.M{"$avg": "$down"}}},
		{"$project": bson.M{"date": "$_id.date", "name": "$_id.name", "supergroup": "$_id.supergroup", "job": "$_id.job", "availability": bson.M{"$multiply": list{bson.M{"$divide": list{"$avgup", bson.M{"$subtract": list{1.00000001, "$avgunknown"}}}}, 100.0}},
			"reliability": bson.M{"$multiply": list{bson.M{"$divide": list{"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgunknown"}}, "$avgdown"}}}}, 100}}}},
		{"$sort": bson.D{{"supergroup", 1}, {"name", 1}, {"date", 1}}}}

	return query
}
