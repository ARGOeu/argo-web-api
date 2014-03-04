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


package ngis

import (
	"labix.org/v2/mgo/bson"
	"time"
	"strconv"
)

type list []interface{}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"


func prepareFilter(input ApiNgiAvailabilityInProfileInput) bson.M{
	
	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))
	
	filter := bson.M{
		"dt": bson.M{"$gte": tsYMD, "$lte": teYMD},
		"p":  bson.M{"$in": input.profile_name},
	}

	if len(input.namespace) > 0 {
		filter["ns"] = bson.M{"$in": input.namespace}
	}

	if len(input.group_name) > 0 {
		filter["n"] = bson.M{"$in": input.group_name}
		// TODO: We do not have the ngi name in the timeline
	}
	return filter
}

func Daily(input ApiNgiAvailabilityInProfileInput) []bson.M{
	filter:=prepareFilter(input)
	// Mongo aggregation pipeline
	// Select all the records that match q
	// Project the results to add 1 to every hepspec(hs) to avoid having 0 as a hepspec
	// Group them by the first 8 digits of datetime (YYYYMMDD) and each group find
	// a = sum(a*hs)
	// r = sum(r*hs)
	// hs = sum(hs)
	// Project to a better format and do these computations
	// a = a/hs
	// r = r/hs
	// Sort by profile->ngi->site->datetime
	query := []bson.M{{"$match": filter}, {"$project": bson.M{"dt": 1, "a": 1, "r": 1, "p": 1, "ns": 1, "n": 1, "hs": bson.M{"$add": list{"$hs", 1}}}}, 
	{"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "n": "$n", "ns": "$ns", "p": "$p"},
	"a": bson.M{"$sum": bson.M{"$multiply": list{"$a", "$hs"}}}, "r": bson.M{"$sum": bson.M{"$multiply": list{"$r", "$hs"}}}, "hs": bson.M{"$sum": "$hs"}}}, 
	{"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ns": "$_id.ns", "p": "$_id.p", "a": bson.M{"$divide": list{"$a", "$hs"}}, 
	"r": bson.M{"$divide": list{"$r", "$hs"}}}}, {"$sort": bson.D{{"p", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}
	//query := []bson.M{{"$match": q}, {"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "n": "$n", "ns": "$ns", "p": "$p"}, "a": bson.M{"$sum": bson.M{"$multiply": list{"$a", "$hs"}}}, 		"r": bson.M{"$sum": bson.M{"$multiply": list{"$r", "$hs"}}}, "hs": bson.M{"$sum": "$hs"}}}, {"$match": bson.M{"hs": bson.M{"$gt": 0}}}, {"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ns": "$_id.ns", "p": 		"$_id.p", "a": bson.M{"$divide": list{"$a", "$hs"}}, "r": bson.M{"$divide": list{"$r", "$hs"}}}}, {"$sort": bson.D{{"p", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}
	
	return query
}

func Monthly(input ApiNgiAvailabilityInProfileInput) []bson.M{
	filter:=prepareFilter(input)
	//PROBABLY THIS LEADS TO THE SAME BUG WE RAN INTO WITH SITES. MUST BE INVESTIGATED!!!!!!!!!!!!
	filter["a"] = bson.M{"$gte": 0}
	filter["r"] = bson.M{"$gte": 0}
	
	// Mongo aggregation pipeline
	// Select all the records that match q
	// Project the results to add 1 to every hepspec(hs) to avoid having 0 as a hepspec
	// Group them by the first 8 digits of datetime (YYYYMMDD) and each group find
	// a = sum(a*hs)
	// r = sum(r*hs)
	// hs = sum(hs)
	// Project to a better format and do these computations
	// a = a/hs
	// r = r/hs
	// Group by the first 6 digits of the datetime (YYYYMM) and by ngi,site,profile and for each group find
	// a = average(a)
	// r = average(r)
	// Project the results to a better format
	// Sort by namespace->profile->ngi->datetime
	
	query := []bson.M{{"$match": filter}, {"$project": bson.M{"dt": 1, "a": 1, "r": 1, "p": 1, "ns": 1, "n": 1, "hs": bson.M{"$add": list{"$hs", 1}}}}, 
	{"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "n": "$n", "ns": "$ns", "p": "$p"}, "a": bson.M{"$sum": bson.M{"$multiply": list{"$a", "$hs"}}},
	"r": bson.M{"$sum": bson.M{"$multiply": list{"$r", "$hs"}}}, "hs": bson.M{"$sum": "$hs"}}}, {"$match": bson.M{"hs": bson.M{"$gt": 0}}}, 
	{"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ns": "$_id.ns", "p": "$_id.p", "a": bson.M{"$divide": list{"$a", "$hs"}}, "r": bson.M{"$divide": list{"$r", "$hs"}}}}, 
	{"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 6}}}, "n": "$n", "ns": "$ns", "p": "$p"}, "a": bson.M{"$avg": "$a"}, 
	"r": bson.M{"$avg": "$r"}}}, {"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ns": "$_id.ns", "p": "$_id.p", "a": 1, "r": 1}}, {"$sort": bson.D{{"ns", 1}, {"p", 1}, {"n", 1}, {"dt", 1}}}}
	
	return query
}