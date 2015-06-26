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

package endpointGroupAvailability

import (
	"encoding/xml"
	"strconv"
	"time"

	"labix.org/v2/mgo/bson"
)

type MongoInterface struct {
	Name         string  `bson:"name"`
	Job          string  `bson:"job"`
	Date         string  `bson:"date"`
	Type         string  `bson:"type"`
	Up           float64 `bson:"uptime"`
	Down         float64 `bson:"downtime"`
	Unknown      float64 `bson:"unknown"`
	Availability float64 `bson:"availability"`
	Reliability  float64 `bson:"reliability"`
	Weights      string  `bson:"weights"`
	SuperGroup   string  `bson:"supergroup"`
}

type Availability struct {
	XMLName      xml.Name `xml:"Availability" json:"-"`
	Timestamp    string   `xml:"timestamp,attr" json:"timestamp"`
	Availability string   `xml:"availability,attr" json:"availability"`
	Reliability  string   `xml:"reliability,attr" json:"reliability"`
}

type EndpointGroup struct {
	XMLName      xml.Name `xml:"EndpointGroup" json:"-"`
	Name         string   `xml:"name,attr" json:"name"`
	SuperGroup   string   `xml:"SuperGroup,attr" json:"supergroup"`
	Availability []*Availability
}

type Job struct {
	XMLName       xml.Name `xml:"Job" json:"-"`
	Name          string   `xml:"name,attr" json:"name"`
	EndpointGroup []*EndpointGroup
}

type Root struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Job     []*Job
}

type EndpointGroupAvailabilityInput struct {
	// mandatory values
	StartTime string // UTC time in W3C format
	EndTime   string // UTC time in W3C format
	Job       string //Job name
	// optional values
	Granularity    string   //availability period; possible values: `DAILY`, MONTHLY`
	Infrastructure string   //infrastructure name
	Production     string   //production or not
	Monitored      string   //yes or no
	Certification  string   //certification status
	Format         string   // default XML; possible values are: XML, JSON
	GroupName      []string // endpointGroup name; may appear more than once
	SuperGroup     []string
}

type list []interface{}

var customForm []string

func init() {
	customForm = []string{"20060102", "2006-01-02T15:04:05Z"} //{"Format that is returned by the database" , "Format that will be used in the generated report"}
}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

func prepareFilter(input EndpointGroupAvailabilityInput) bson.M {
	ts, _ := time.Parse(zuluForm, input.StartTime)
	te, _ := time.Parse(zuluForm, input.EndTime)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	// Construct the query to mongodb based on the input
	filter := bson.M{
		"date": bson.M{"$gte": tsYMD, "$lte": teYMD},
		"job":  input.Job,
	}

	if len(input.GroupName) > 0 {
		filter["name"] = bson.M{"$in": input.GroupName}
	}

	if len(input.SuperGroup) > 0 {
		filter["supergroup"] = bson.M{"$in": input.SuperGroup}
	}

	return filter
}

func Daily(input EndpointGroupAvailabilityInput) []bson.M {
	filter := prepareFilter(input)

	// Mongo aggregation pipeline
	// Select all the records that match q
	// Project to select just the first 8 digits of the date YYYYMMDD
	// Sort by profile->supergroup->endpointGroup->datetime
	query := []bson.M{
		{"$match": filter},
		{"$project": bson.M{
			"date":         bson.M{"$substr": list{"$date", 0, 8}},
			"availability": 1,
			"reliability":  1,
			"job":          1,
			"supergroup":   1,
			"name":         1}},
		{"$sort": bson.D{
			{"job", 1},
			{"supergroup", 1},
			{"name", 1},
			{"date", 1}}}}

	return query
}

func Monthly(input EndpointGroupAvailabilityInput) []bson.M {

	filter := prepareFilter(input)

	// Mongo aggregation pipeline
	// Select all the records that match q
	// Group them by the first six digits of their date (YYYYMM), their supergroup, their endpointGroup, their profile, etc...
	// from that group find the average of the uptime, u, downtime
	// Project the result to a better format and do this computation
	// availability = (avgup/(1.00000001 - avgu))*100
	// reliability = (avgup/((1.00000001 - avgu)-avgd))*100
	// Sort the results by namespace->profile->supergroup->endpointGroup->datetime

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"date":       bson.M{"$substr": list{"$date", 0, 6}},
				"name":       "$name",
				"supergroup": "$supergroup",
				"job":        "$job"},
			"avguptime": bson.M{"$avg": "$uptime"},
			"avgunkown": bson.M{"$avg": "$unknown"},
			"avgdown":   bson.M{"$avg": "$downtime"}}},
		{"$project": bson.M{
			"date":       "$_id.date",
			"name":       "$_id.name",
			"job":        "$_id.job",
			"supergroup": "$_id.supergroup",
			"avguptime":  1,
			"avgunkown":  1,
			"avgdown":    1,
			"availability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avguptime", bson.M{"$subtract": list{1.00000001, "$avgunkown"}}}},
					100}},
			"reliability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avguptime", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgunkown"}}, "$avgdown"}}}},
					100}}}},
		{"$sort": bson.D{
			{"job", 1},
			{"supergroup", 1},
			{"name", 1},
			{"date", 1}}}}

	return query
}
