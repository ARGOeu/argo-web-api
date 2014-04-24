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

package serviceFlavors

import (
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

type list []interface{}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

func prepareFilter(input ApiSFAvailabilityInProfileInput) bson.M {

	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	filter := bson.M{
		"p" : input.profile,
		"dt": bson.M{"$gte": tsYMD, "$lte": teYMD},
	}	
	
	if len(input.group_name) > 0 {
		filter["sf"] = bson.M{"$in": input.group_name}
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
		{"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "sf": "$sf", "s":"$s", "a":"$a", "r":"$r","p":"$p"}}},
		{"$project": bson.M{"dt": "$_id.dt", "sf": "$_id.sf","a":"$_id.a","r":"$_id.r","s":"$_id.s","p":"$_id.p"}},
		{"$sort": bson.D{ {"s",1}, {"sf", 1}, {"dt", 1}}}}

	
	return query
}

func Monthly(input ApiSFAvailabilityInProfileInput) []bson.M {
	
	filter := prepareFilter(input)

	query := []bson.M{
		{"$match": filter}, 
		{"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 6}}}, "s":"$s", "p": "$p", "sf": "$sf"}, "a": bson.M{"$avg": "$a"},"r": bson.M{"$avg": "$r"}}},
		{"$project": bson.M{"dt": "$_id.dt", "sf": "$_id.sf", "a":"$a","r":"$r","s":"$_id.s","p":"$_id.p"}},
		{"$sort": bson.D{{"s",1},{"sf", 1}, {"dt", 1}}}}

	return query
}
