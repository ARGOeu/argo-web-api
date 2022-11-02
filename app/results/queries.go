/*
 * Copyright (c) 2015 GRNET S.A.
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
 * or implied, of GRNET S.A.
 *
 */

package results

import (
	"gopkg.in/mgo.v2/bson"
)

// FlatDailyEndpoint query to aggregate daily endpoint a/r results from mongoDB
func FlatDailyEndpoint(filter bson.M, limit int, skip int) []bson.M {

	query := []bson.M{
		{"$match": filter},
		{"$project": bson.M{
			"_id":          1,
			"date":         bson.D{{"$substr", list{"$date", 0, 8}}},
			"name":         1,
			"availability": 1,
			"reliability":  1,
			"unknown":      1,
			"up":           1,
			"down":         1,
			"supergroup":   1,
			"service":      1,
			"info":         1,
			"report":       1}},
		{"$sort": bson.D{
			{"name", 1},
			{"service", 1},
			{"supergroup", 1},
			{"date", 1},
		}}}

	if limit > 0 {
		query = append(query, bson.M{"$skip": skip})
		query = append(query, bson.M{"$limit": limit + 1})

	}

	return query
}

// DailyEndpoint query to aggregate daily endpoint a/r results from mongoDB
func DailyEndpoint(filter bson.M) []bson.M {
	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"id":           "$_id",
				"date":         bson.D{{"$substr", list{"$date", 0, 8}}},
				"name":         "$name",
				"supergroup":   "$supergroup",
				"service":      "$service",
				"availability": "$availability",
				"reliability":  "$reliability",
				"unknown":      "$unknown",
				"up":           "$up",
				"down":         "$down",
				"report":       "$report"},
			"info": bson.M{"$first": "$info"},
		}},

		{"$project": bson.M{
			"_id":          "$_id.id",
			"date":         "$_id.date",
			"name":         "$_id.name",
			"availability": "$_id.availability",
			"reliability":  "$_id.reliability",
			"unknown":      "$_id.unknown",
			"up":           "$_id.up",
			"down":         "$_id.down",
			"supergroup":   "$_id.supergroup",
			"service":      "$_id.service",
			"info":         "$info",
			"report":       "$_id.report"}},
		{"$sort": bson.D{
			{"supergroup", 1},
			{"service", 1},
			{"name", 1},
			{"date", 1}}}}

	return query
}

// FlatMonthlyEndpoint query to aggregate monthly a/r results from mongoDB
func FlatMonthlyEndpoint(filter bson.M, limit int, skip int) []bson.M {
	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"date":       bson.D{{"$substr", list{"$date", 0, 6}}},
				"name":       "$name",
				"supergroup": "$supergroup",
				"service":    "$service",
				"report":     "$report"},
			"avgup":      bson.M{"$avg": "$up"},
			"avgunknown": bson.M{"$avg": "$unknown"},
			"avgdown":    bson.M{"$avg": "$down"},
			"info":       bson.M{"$first": "$info"}}},
		{"$project": bson.M{
			"date":       "$_id.date",
			"name":       "$_id.name",
			"supergroup": "$_id.supergroup",
			"service":    "$_id.service",
			"report":     "$_id.report",
			"info":       "$info",
			"unknown":    "$avgunknown",
			"up":         "$avgup",
			"down":       "$avgdown",
			"availability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{1.00000001, "$avgunknown"}}}},
					100}},
			"reliability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgunknown"}}, "$avgdown"}}}},
					100}}}},
		{"$sort": bson.D{
			{"name", 1},
			{"service", 1},
			{"supergroup", 1},
			{"date", 1}}}}

	if limit > 0 {
		query = append(query, bson.M{"$skip": skip})
		query = append(query, bson.M{"$limit": limit + 1})

	}

	return query
}

// FlatCustomEndpoint query to aggregate a/r results for a custom period of time from mongoDB
func FlatCustomEndpoint(filter bson.M, limit int, skip int) []bson.M {
	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"name":       "$name",
				"supergroup": "$supergroup",
				"service":    "$service",
				"report":     "$report"},
			"avgup":      bson.M{"$avg": "$up"},
			"avgunknown": bson.M{"$avg": "$unknown"},
			"avgdown":    bson.M{"$avg": "$down"},
			"info":       bson.M{"$first": "$info"}}},
		{"$project": bson.M{
			"name":       "$_id.name",
			"supergroup": "$_id.supergroup",
			"service":    "$_id.service",
			"report":     "$_id.report",
			"info":       "$info",
			"unknown":    "$avgunknown",
			"up":         "$avgup",
			"down":       "$avgdown",
			"availability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{1.00000001, "$avgunknown"}}}},
					100}},
			"reliability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgunknown"}}, "$avgdown"}}}},
					100}}}},
		{"$sort": bson.D{
			{"name", 1},
			{"service", 1},
			{"supergroup", 1}}}}

	if limit > 0 {
		query = append(query, bson.M{"$skip": skip})
		query = append(query, bson.M{"$limit": limit + 1})

	}

	return query
}

// MonthlyEndpoint query to aggregate monthly a/r results from mongoDB
func MonthlyEndpoint(filter bson.M) []bson.M {
	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"date":       bson.D{{"$substr", list{"$date", 0, 6}}},
				"name":       "$name",
				"supergroup": "$supergroup",
				"service":    "$service",
				"report":     "$report"},
			"avgup":      bson.M{"$avg": "$up"},
			"avgunknown": bson.M{"$avg": "$unknown"},
			"avgdown":    bson.M{"$avg": "$down"},
			"info":       bson.M{"$first": "$info"}}},

		{"$project": bson.M{
			"date":       "$_id.date",
			"name":       "$_id.name",
			"supergroup": "$_id.supergroup",
			"service":    "$_id.service",
			"report":     "$_id.report",
			"info":       "$info",
			"unknown":    "$avgunknown",
			"up":         "$avgup",
			"down":       "$avgdown",
			"availability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{1.00000001, "$avgunknown"}}}},
					100}},
			"reliability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgunknown"}}, "$avgdown"}}}},
					100}}}},
		{"$sort": bson.D{
			{"supergroup", 1},
			{"name", 1},
			{"date", 1}}}}
	return query
}

// CustomEndpoint query to aggregate a/r results over a custom period of time from mongoDB
func CustomEndpoint(filter bson.M) []bson.M {
	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"name":       "$name",
				"supergroup": "$supergroup",
				"service":    "$service",
				"report":     "$report"},
			"avgup":      bson.M{"$avg": "$up"},
			"avgunknown": bson.M{"$avg": "$unknown"},
			"avgdown":    bson.M{"$avg": "$down"},
			"info":       bson.M{"$first": "$info"}}},

		{"$project": bson.M{
			"name":       "$_id.name",
			"supergroup": "$_id.supergroup",
			"service":    "$_id.service",
			"report":     "$_id.report",
			"info":       "$info",
			"unknown":    "$avgunknown",
			"up":         "$avgup",
			"down":       "$avgdown",
			"availability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{1.00000001, "$avgunknown"}}}},
					100}},
			"reliability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgunknown"}}, "$avgdown"}}}},
					100}}}},
		{"$sort": bson.D{
			{"supergroup", 1},
			{"name", 1}}}}
	return query
}

// DailyServiceFlavor query to aggregate daily SF results from mongoDB
func DailyServiceFlavor(filter bson.M) []bson.M {
	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"date":         bson.D{{"$substr", list{"$date", 0, 8}}},
				"name":         "$name",
				"supergroup":   "$supergroup",
				"availability": "$availability",
				"reliability":  "$reliability",
				"unknown":      "$unknown",
				"up":           "$up",
				"down":         "$down",
				"report":       "$report"}}},
		{"$project": bson.M{
			"date":         "$_id.date",
			"name":         "$_id.name",
			"availability": "$_id.availability",
			"reliability":  "$_id.reliability",
			"unknown":      "$_id.unknown",
			"up":           "$_id.up",
			"down":         "$_id.down",
			"supergroup":   "$_id.supergroup",
			"report":       "$_id.report"}},
		{"$sort": bson.D{
			{"supergroup", 1},
			{"name", 1},
			{"date", 1}}}}

	return query
}

// MonthlyServiceFlavor query to aggregate daily SF results from mongoDB
func MonthlyServiceFlavor(filter bson.M) []bson.M {
	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"date":       bson.D{{"$substr", list{"$date", 0, 6}}},
				"name":       "$name",
				"supergroup": "$supergroup",
				"report":     "$report"},
			"avgup":      bson.M{"$avg": "$up"},
			"avgunknown": bson.M{"$avg": "$unknown"},
			"avgdown":    bson.M{"$avg": "$down"}}},
		{"$project": bson.M{
			"date":       "$_id.date",
			"name":       "$_id.name",
			"supergroup": "$_id.supergroup",
			"report":     "$_id.report",
			"unknown":    "$avgunknown",
			"up":         "$avgup",
			"down":       "$avgdown",
			"availability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{1.00000001, "$avgunknown"}}}},
					100}},
			"reliability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgunknown"}}, "$avgdown"}}}},
					100}}}},
		{"$sort": bson.D{
			{"supergroup", 1},
			{"name", 1},
			{"date", 1}}}}
	return query
}

// CustomServiceFlavor query to aggregate SF results over a custom period of time from mongoDB
func CustomServiceFlavor(filter bson.M) []bson.M {
	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"name":       "$name",
				"supergroup": "$supergroup",
				"report":     "$report"},
			"avgup":      bson.M{"$avg": "$up"},
			"avgunknown": bson.M{"$avg": "$unknown"},
			"avgdown":    bson.M{"$avg": "$down"}}},
		{"$project": bson.M{
			"name":       "$_id.name",
			"supergroup": "$_id.supergroup",
			"report":     "$_id.report",
			"unknown":    "$avgunknown",
			"up":         "$avgup",
			"down":       "$avgdown",
			"availability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{1.00000001, "$avgunknown"}}}},
					100}},
			"reliability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgunknown"}}, "$avgdown"}}}},
					100}}}},
		{"$sort": bson.D{
			{"supergroup", 1},
			{"name", 1}}}}
	return query
}

// DailyEndpointGroup query to aggregate daily results from mongodb
func DailyEndpointGroup(filter bson.M) []bson.M {
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
			"unknown":      1,
			"up":           1,
			"down":         1,
			"report":       1,
			"supergroup":   1,
			"name":         1}},
		{"$sort": bson.D{
			{"report", 1},
			{"supergroup", 1},
			{"name", 1},
			{"date", 1}}}}

	return query
}

// MonthlyEndpointGroup query to aggregate monthly results from mongodb
func MonthlyEndpointGroup(filter bson.M) []bson.M {

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
				"report":     "$report"},
			"avgup":      bson.M{"$avg": "$up"},
			"avgunknown": bson.M{"$avg": "$unknown"},
			"avgdown":    bson.M{"$avg": "$down"}}},
		{"$project": bson.M{
			"date":       "$_id.date",
			"name":       "$_id.name",
			"report":     "$_id.report",
			"supergroup": "$_id.supergroup",
			"unknown":    "$avgunknown",
			"up":         "$avgup",
			"down":       "$avgdown",
			"avgup":      1,
			"avgunknown": 1,
			"avgdown":    1,
			"availability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{1.00000001, "$avgunknown"}}}},
					100}},
			"reliability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgunknown"}}, "$avgdown"}}}},
					100}}}},
		{"$sort": bson.D{
			{"report", 1},
			{"supergroup", 1},
			{"name", 1},
			{"date", 1}}}}

	return query
}

// CustomEndpointGroup query to aggregate results over a custom period from mongodb
func CustomEndpointGroup(filter bson.M) []bson.M {

	// Mongo aggregation pipeline
	// Select all the records that match q
	// Group them by their supergroup, their endpointGroup, their profile, etc...
	// from that group find the average of the uptime, u, downtime
	// Project the result to a better format and do this computation
	// availability = (avgup/(1.00000001 - avgu))*100
	// reliability = (avgup/((1.00000001 - avgu)-avgd))*100
	// Sort the results by namespace->profile->supergroup->endpointGroup

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"name":       "$name",
				"supergroup": "$supergroup",
				"report":     "$report"},
			"avgup":      bson.M{"$avg": "$up"},
			"avgunknown": bson.M{"$avg": "$unknown"},
			"avgdown":    bson.M{"$avg": "$down"}}},
		{"$project": bson.M{
			"name":       "$_id.name",
			"report":     "$_id.report",
			"supergroup": "$_id.supergroup",
			"unknown":    "$avgunknown",
			"up":         "$avgup",
			"down":       "$avgdown",
			"avgup":      1,
			"avgunknown": 1,
			"avgdown":    1,
			"availability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{1.00000001, "$avgunknown"}}}},
					100}},
			"reliability": bson.M{
				"$multiply": list{
					bson.M{"$divide": list{
						"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgunknown"}}, "$avgdown"}}}},
					100}}}},
		{"$sort": bson.D{
			{"report", 1},
			{"supergroup", 1},
			{"name", 1}}}}

	return query
}

// DailySuperGroup function to build the MongoDB aggregation query for daily calculations
func DailySuperGroup(filter bson.M) []bson.M {
	// The following aggregation query consists of 5 grand steps
	// 1. Match   : records for the specific date and report and supergroup(optional)
	// 2. Project : all necessary fields (date,availability,reliability,report) etc but also
	//              if avail >= 0 set an availability-weigh = weight + 1, else = 0
	//							if rel >=0 set a reliability-weight = weight + 1, else = 0
	//              keep also weight = weight + 1 (to compensate for zero values)
	//
	//              Keeping two extra weights (a/r) has the following result:
	//               - If an item has undef availab. then it will have an weightAv=0 and will not affect sums
	//                    for eg. avg_daily_supergroup_availability = (av1*w1 + av2*w2 + undefAv3*0) / (w1 + w1 + 0)
	//               - If an item has undef reliab. then it will have an weightRel=0 and will not affect sums
	//                    for eg. avg_daily_supergroup_reliability = (rel1*w2 + rel2*w2 + undefRel3*0) / (w1 + w1 + 0)
	//
	// 3. Group   : by supergroup and day and calculate the sum of weighted daily availabilites (and reliabilities also)
	//              - availability(weighted_sum) = av1*w1 + av2*w2 + undefAv3*0 etc...
	//              - reliability(weighted_sum) = rel1*w1 + rel2*w2 + undefRel3*0 etc...
	//
	// 4. Match   : assertion step - keep only items that have a valid weight > 0
	// 5. Project : the previous results and try to find the weighted average of daily avail. and reliability by:
	//              - divide the previous sum of weighted availabilities by the total weightAv
	//                SPECIAL CASE: If total weightAv remains : 0 that means that total daily supergroup avail = undef
	//                              so instead of a numeric value, add a "nan" string (will not be counted in monthly average)
	//              - divide the previous sum of weighted availabilities by the total weightAv
	//								SPECIAL CASE: If total weightRem remains : 0 that means that total daily supergroup rel = undef
	//                              so instead of a numeric value, add a "nan" string (will not be counted in monthly average)
	// 6. Project : the relevant fields to form the appropriate final response (date,supergroup,report,avail,rel)
	// 7. Sort    : the final results by report, supergroup and then date
	query := []bson.M{
		{"$match": filter},
		{"$project": bson.M{
			"date":         1,
			"availability": 1,
			"reliability":  1,
			"report":       1,
			"supergroup":   1,
			"weightAv":     bson.M{"$cond": list{bson.M{"$gte": list{"$availability", 0}}, bson.M{"$add": list{"$weight", 1}}, 0}},
			"weightRel":    bson.M{"$cond": list{bson.M{"$gte": list{"$reliability", 0}}, bson.M{"$add": list{"$weight", 1}}, 0}},
			"weight": bson.M{
				"$add": list{"$weight", 1}}},
		},
		{"$group": bson.M{
			"_id": bson.M{
				"date":       bson.D{{"$substr", list{"$date", 0, 8}}},
				"supergroup": "$supergroup",
				"report":     "$report"},
			"availability": bson.M{"$sum": bson.M{"$multiply": list{"$availability", "$weightAv"}}},
			"reliability":  bson.M{"$sum": bson.M{"$multiply": list{"$reliability", "$weightRel"}}},
			"weightAv":     bson.M{"$sum": "$weightAv"},
			"weightRel":    bson.M{"$sum": "$weightRel"},
			"weight":       bson.M{"$sum": "$weight"}},
		},
		{"$match": bson.M{
			"weight": bson.M{"$gt": 0}},
		},
		{"$project": bson.M{
			"date":         "$_id.date",
			"supergroup":   "$_id.supergroup",
			"report":       "$_id.report",
			"availability": bson.M{"$cond": list{bson.M{"$gt": list{"$weightAv", 0}}, bson.M{"$divide": list{"$availability", "$weightAv"}}, "nan"}},
			"reliability":  bson.M{"$cond": list{bson.M{"$gt": list{"$weightRel", 0}}, bson.M{"$divide": list{"$reliability", "$weightRel"}}, "nan"}}},
		},
		{"$project": bson.M{
			"date":         "$_id.date",
			"supergroup":   "$_id.supergroup",
			"report":       "$_id.report",
			"availability": 1,
			"reliability":  1},
		},
		{"$sort": bson.D{
			{"report", 1},
			{"supergroup", 1},
			{"date", 1}},
		}}

	return query
}

// MonthlySuperGroup function to build the MongoDB aggregation query for monthly calculations
func MonthlySuperGroup(filter bson.M) []bson.M {

	// The following aggregation query consists of 5 grand steps
	// 1. Match   : records for the specific date and report and supergroup(optional)
	// 2. Project : all necessary fields (date,availability,reliability,report) etc but also
	//              if avail >= 0 set an availability-weigh = weight + 1, else = 0
	//							if rel >=0 set a reliability-weight = weight + 1, else = 0
	//              keep also weight = weight + 1 (to compensate for zero values)
	//
	//              Keeping two extra weights (a/r) has the following result:
	//               - If an item has undef availab. then it will have an weightAv=0 and will not affect sums
	//                    for eg. avg_daily_supergroup_availability = (av1*w1 + av2*w2 + undefAv3*0) / (w1 + w1 + 0)
	//               - If an item has undef reliab. then it will have an weightRel=0 and will not affect sums
	//                    for eg. avg_daily_supergroup_reliability = (rel1*w2 + rel2*w2 + undefRel3*0) / (w1 + w1 + 0)
	//
	// 3. Group   : by supergroup and day and calculate the sum of weighted daily availabilites (and reliabilities also)
	//              - availability(weighted_sum) = av1*w1 + av2*w2 + undefAv3*0 etc...
	//              - reliability(weighted_sum) = rel1*w1 + rel2*w2 + undefRel3*0 etc...
	//
	// 4. Match   : assertion step - keep only items that have a valid weight > 0
	// 5. Project : the previous results and try to find the weighted average of daily avail. and reliability by:
	//              - divide the previous sum of weighted availabilities by the total weightAv
	//                SPECIAL CASE: If total weightAv remains : 0 that means that total daily supergroup avail = undef
	//                              so instead of a numeric value, add a "nan" string (will not be counted in monthly average)
	//              - divide the previous sum of weighted availabilities by the total weightAv
	//								SPECIAL CASE: If total weightRem remains : 0 that means that total daily supergroup rel = undef
	//                              so instead of a numeric value, add a "nan" string (will not be counted in monthly average)
	// 6. Group   : by supergroup and report to calculate ar for custom period
	//							- custom period availability avg = avg(daily_availabilities) ~ but items with "nan" values will be neglected
	//						  - custom period reliability avg = avg(daily_reliabilities) ~ but items with "nan" values will be neglected
	//
	// 7. Project : the relevant fields to form the appropriate final response (supergroup,report,avail,rel)
	// 8. Sort    : the final results by report, supergroup

	query := []bson.M{
		{"$match": filter},
		{"$project": bson.M{
			"date":         1,
			"availability": 1,
			"reliability":  1,
			"supergroup":   1,
			"weightAv":     bson.M{"$cond": list{bson.M{"$gte": list{"$availability", 0}}, bson.M{"$add": list{"$weight", 1}}, 0}},
			"weightRel":    bson.M{"$cond": list{bson.M{"$gte": list{"$reliability", 0}}, bson.M{"$add": list{"$weight", 1}}, 0}},
			"weight": bson.M{
				"$add": list{"$weight", 1}}},
		},
		{"$group": bson.M{
			"_id": bson.M{
				"date":       bson.D{{"$substr", list{"$date", 0, 8}}},
				"supergroup": "$supergroup",
				"report":     "$report"},
			"availability": bson.M{"$sum": bson.M{"$multiply": list{"$availability", "$weightAv"}}},
			"reliability":  bson.M{"$sum": bson.M{"$multiply": list{"$reliability", "$weightRel"}}},
			"weightAv":     bson.M{"$sum": "$weightAv"},
			"weightRel":    bson.M{"$sum": "$weightRel"},
			"weight":       bson.M{"$sum": "$weight"}},
		},
		{"$match": bson.M{
			"weight": bson.M{"$gt": 0}},
		},
		{"$project": bson.M{
			"date":         "$_id.date",
			"supergroup":   "$_id.supergroup",
			"report":       "$_id.report",
			"availability": bson.M{"$cond": list{bson.M{"$gt": list{"$weightAv", 0}}, bson.M{"$divide": list{"$availability", "$weightAv"}}, "nan"}},
			"reliability":  bson.M{"$cond": list{bson.M{"$gt": list{"$weightRel", 0}}, bson.M{"$divide": list{"$reliability", "$weightRel"}}, "nan"}}},
		},
		{"$group": bson.M{
			"_id": bson.M{
				"supergroup": "$supergroup", "report": "$report"},
			"availability": bson.M{"$avg": "$availability"},
			"reliability":  bson.M{"$avg": "$reliability"}},
		},
		{"$project": bson.M{
			"supergroup":   "$_id.supergroup",
			"report":       "$_id.report",
			"availability": 1,
			"reliability":  1},
		},
		{"$sort": bson.D{
			{"report", 1},
			{"supergroup", 1}},
		}}

	return query
}

// CustomSuperGroup function to build the MongoDB aggregation query for custom period aggregation
func CustomSuperGroup(filter bson.M) []bson.M {

	// The following aggregation query consists of 5 grand steps
	// 1. Match   : records for the specific date and report and supergroup(optional)
	// 2. Project : all necessary fields (date,availability,reliability,report) etc but also
	//              if avail >= 0 set an availability-weigh = weight + 1, else = 0
	//							if rel >=0 set a reliability-weight = weight + 1, else = 0
	//              keep also weight = weight + 1 (to compensate for zero values)
	//
	//              Keeping two extra weights (a/r) has the following result:
	//               - If an item has undef availab. then it will have an weightAv=0 and will not affect sums
	//                    for eg. avg_daily_supergroup_availability = (av1*w1 + av2*w2 + undefAv3*0) / (w1 + w1 + 0)
	//               - If an item has undef reliab. then it will have an weightRel=0 and will not affect sums
	//                    for eg. avg_daily_supergroup_reliability = (rel1*w2 + rel2*w2 + undefRel3*0) / (w1 + w1 + 0)
	//
	// 3. Group   : by supergroup and day and calculate the sum of weighted daily availabilites (and reliabilities also)
	//              - availability(weighted_sum) = av1*w1 + av2*w2 + undefAv3*0 etc...
	//              - reliability(weighted_sum) = rel1*w1 + rel2*w2 + undefRel3*0 etc...
	//
	// 4. Match   : assertion step - keep only items that have a valid weight > 0
	// 5. Project : the previous results and try to find the weighted average of daily avail. and reliability by:
	//              - divide the previous sum of weighted availabilities by the total weightAv
	//                SPECIAL CASE: If total weightAv remains : 0 that means that total daily supergroup avail = undef
	//                              so instead of a numeric value, add a "nan" string (will not be counted in monthly average)
	//              - divide the previous sum of weighted availabilities by the total weightAv
	//								SPECIAL CASE: If total weightRem remains : 0 that means that total daily supergroup rel = undef
	//                              so instead of a numeric value, add a "nan" string (will not be counted in monthly average)
	// 6. Group   : by first date part (month, eg: 201608) to calculate monthly average avail and rel.
	//							- monthly availability avg = avg(daily_availabilities) ~ but items with "nan" values will be neglected
	//						  - monthly reliability avg = avg(daily_reliabilities) ~ but items with "nan" values will be neglected
	//
	// 7. Project : the relevant fields to form the appropriate final response (date,supergroup,report,avail,rel)
	// 8. Sort    : the final results by report, supergroup and then date

	query := []bson.M{
		{"$match": filter},
		{"$project": bson.M{
			"date":         1,
			"availability": 1,
			"reliability":  1,
			"report":       1,
			"supergroup":   1,
			"weightAv":     bson.M{"$cond": list{bson.M{"$gte": list{"$availability", 0}}, bson.M{"$add": list{"$weight", 1}}, 0}},
			"weightRel":    bson.M{"$cond": list{bson.M{"$gte": list{"$reliability", 0}}, bson.M{"$add": list{"$weight", 1}}, 0}},
			"weight": bson.M{
				"$add": list{"$weight", 1}}},
		},
		{"$group": bson.M{
			"_id": bson.M{
				"date":       bson.D{{"$substr", list{"$date", 0, 8}}},
				"supergroup": "$supergroup",
				"report":     "$report"},
			"availability": bson.M{"$sum": bson.M{"$multiply": list{"$availability", "$weightAv"}}},
			"reliability":  bson.M{"$sum": bson.M{"$multiply": list{"$reliability", "$weightRel"}}},
			"weightAv":     bson.M{"$sum": "$weightAv"},
			"weightRel":    bson.M{"$sum": "$weightRel"},
			"weight":       bson.M{"$sum": "$weight"}},
		},
		{"$match": bson.M{
			"weight": bson.M{"$gt": 0}},
		},
		{"$project": bson.M{
			"date":         "$_id.date",
			"supergroup":   "$_id.supergroup",
			"report":       "$_id.report",
			"availability": bson.M{"$cond": list{bson.M{"$gt": list{"$weightAv", 0}}, bson.M{"$divide": list{"$availability", "$weightAv"}}, "nan"}},
			"reliability":  bson.M{"$cond": list{bson.M{"$gt": list{"$weightRel", 0}}, bson.M{"$divide": list{"$reliability", "$weightRel"}}, "nan"}}},
		},
		{"$group": bson.M{
			"_id": bson.M{
				"date":       bson.D{{"$substr", list{"$date", 0, 6}}},
				"supergroup": "$supergroup", "report": "$report"},
			"availability": bson.M{"$avg": "$availability"},
			"reliability":  bson.M{"$avg": "$reliability"}},
		},
		{"$project": bson.M{
			"date":         "$_id.date",
			"supergroup":   "$_id.supergroup",
			"report":       "$_id.report",
			"availability": 1,
			"reliability":  1},
		},
		{"$sort": bson.D{
			{"report", 1},
			{"supergroup", 1},
			{"date", 1}},
		}}

	return query
}
