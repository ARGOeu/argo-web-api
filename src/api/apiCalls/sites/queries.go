/*
 * Copyright (c) 2013 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
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

package sites

import (
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

type list []interface{}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

func prepareFilter(input ApiSiteAvailabilityInProfileInput) bson.M {
	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	// Construct the query to mongodb based on the input
	filter := bson.M{
		"dt": bson.M{"$gte": tsYMD, "$lte": teYMD},
		"ap": input.availability_profile,
	}

	if len(input.group_name) > 0 {
		filter["s"] = bson.M{"$in": input.group_name}
	}

	filter["i"] = input.infrastructure
	filter["cs"] = input.certification
	filter["pr"] = input.production
	filter["m"] = input.monitored

	//TODO: Remove hardcoded filtering and add this
	// as a parameter
	filter["sc"] = "EGI"
	filter["ss"] = "EGI"

	return filter
}

func Daily(input ApiSiteAvailabilityInProfileInput) []bson.M {
	filter := prepareFilter(input)

	// Mongo aggregation pipeline
	// Select all the records that match q
	// Project to select just the first 8 digits of the date YYYYMMDD
	// Sort by profile->ngi->site->datetime
	query := []bson.M{
		{"$match": filter},
		{"$project": bson.M{"dt": bson.M{"$substr": list{"$dt", 0, 8}}, "ap":1 ,"i": 1, "sc": 1, "ss": 1, "n": 1, "pr": 1, "m": 1, "cs": 1, "s": 1, "a": 1, "r": 1}},
		{"$sort": bson.D{{"ap", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}

	return query
}

func Monthly(input ApiSiteAvailabilityInProfileInput) []bson.M {

	filter := prepareFilter(input)

	// Mongo aggregation pipeline
	// Select all the records that match q
	// Group them by the first six digits of their date (YYYYMM), their ngi, their site, their profile, etc...
	// from that group find the average of the uptime, u, downtime
	// Project the result to a better format and do this computation
	// availability = (avgup/(1.00000001 - avgu))*100
	// reliability = (avgup/((1.00000001 - avgu)-avgd))*100
	// Sort the results by namespace->profile->ngi->site->datetime

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": bson.M{"dt": bson.M{"$substr": list{"$dt", 0, 6}}, "i": "$i", "n": "$n", "pr": "$pr", "m": "$m", "cs": "$cs", "s": "$s", "ap":"$ap"},
			"avgup": bson.M{"$avg": "$up"}, "avgu": bson.M{"$avg": "$u"}, "avgd": bson.M{"$avg": "$d"}}},
		{"$project": bson.M{"dt": "$_id.dt", "i": "$_id.i", "n": "$_id.n", "pr": "$_id.pr", "m": "$_id.m", "cs": "$_id.cs", "s": "$_id.s", "ap": "$_id.ap", "avgup": 1, "avgu": 1, "avgd": 1,
			"a": bson.M{"$multiply": list{bson.M{"$divide": list{"$avgup", bson.M{"$subtract": list{1.00000001, "$avgu"}}}}, 100}},
			"r": bson.M{"$multiply": list{bson.M{"$divide": list{"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgu"}}, "$avgd"}}}}, 100}}}},
		{"$sort": bson.D{{"ap", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}

	return query
}
