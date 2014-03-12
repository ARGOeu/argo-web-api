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

package vos

import (
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

type list []interface{}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

func prepareFilter(input ApiVoAvailabilityInProfileInput) bson.M {

	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	filter := bson.M{
		"v": input.vo_name,
		"p":bson.M{"$in": input.profile_name},
		"d": bson.M{"$gte": tsYMD, "$lte": teYMD},
	}	
	
	return filter
}

func Daily(input ApiVoAvailabilityInProfileInput) []bson.M {
	
	filter := prepareFilter(input)
	
	query := []bson.M{
		{"$match": filter}, 
		{"$group": bson.M{"_id": bson.M{"d": bson.D{{"$substr", list{"$d", 0, 8}}}, "p": "$p", "v": "$v", "a": "$a","r": "$r"}}},
		{"$project": bson.M{"d": "$_id.d", "v": "$_id.v", "p": "$_id.p","a":"$_id.a","r":"$_id.r"}},
		{"$sort": bson.D{{"p", 1}, {"d", 1}}}}

	
	return query
}

func Monthly(input ApiVoAvailabilityInProfileInput) []bson.M {
	filter := prepareFilter(input)
	
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

	query := []bson.M{
		{"$match": filter}, 
		{"$group": bson.M{"_id": bson.M{"d": bson.D{{"$substr", list{"$d", 0, 6}}}, "p": "$p", "v": "$v"}, "a": bson.M{"$avg": "$a"},"r": bson.M{"$avg": "$r"}}},
		{"$project": bson.M{"d": "$_id.d", "v": "$_id.v", "p": "$_id.p","a":"$a","r":"$r"}},
		{"$sort": bson.D{{"p", 1}, {"d", 1}}}}

	return query
}
