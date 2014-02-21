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

package main

import (
	"api/ngis"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NgiAvailabilityInProfile(w http.ResponseWriter, r *http.Request) []byte {

	// This is the input we will receive from the API

	type ApiNgiAvailabilityInProfileInput struct {
		// mandatory values
		start_time         string   // UTC time in W3C format
		end_time           string   // UTC time in W3C format
		profile_name       []string // may appear more than once. (eg: CMS_CRITICAL)
		group_type         []string // may appear more than once. (eg: CMS_Site)
		availabilityperiod string   // availability period; possible values: `HOURLY`, `DAILY`, `WEEKLY`, `MONTHLY`
		// optional values
		output     string   // default XML; possible values are: XML, JSON
		namespace  []string // profile namespace; may appear more than once. (eg: ch.cern.sam)
		group_name []string // ngi name; may appear more than once
	}

	urlValues := r.URL.Query()

	input := ApiNgiAvailabilityInProfileInput{
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues["profile_name"],
		urlValues["group_type"],
		urlValues.Get("type"),
		urlValues.Get("output"),
		urlValues["namespace"],
		urlValues["group_name"],
	}
	customForm := []string{"20060102", "2006-01-02"} //{"Format that is returned by the database" , "Format that will be used in the generated report"}

	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	if cfg.Server.Cache == true {
		out, found := httpcache.Get("ngi " + fmt.Sprint(input))
		if found {
			return []byte(fmt.Sprint(out))
		}
	}
	session, err := mgo.Dial(cfg.MongoDB.Host + ":" + fmt.Sprint(cfg.MongoDB.Port))
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("AR").C("sites")
	results := []ngis.MongoNgi{}
	q := bson.M{
		"dt": bson.M{"$gte": tsYMD, "$lte": teYMD},
		"p":  bson.M{"$in": input.profile_name},
	}

	if len(input.namespace) > 0 {
		q["ns"] = bson.M{"$in": input.namespace}
	}

	if len(input.group_name) > 0 {
		q["n"] = bson.M{"$in": input.group_name}
		// TODO: We do not have the ngi name in the timeline
	}

	fmt.Println(input)

	if len(input.availabilityperiod) == 0 || strings.ToLower(input.availabilityperiod) == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
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
		query := []bson.M{{"$match": q}, {"$project": bson.M{"dt": 1, "a": 1, "r": 1, "p": 1, "ns": 1, "n": 1, "hs": bson.M{"$add": list{"$hs", 1}}}}, {"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "n": "$n", "ns": "$ns", "p": "$p"}, "a": bson.M{"$sum": bson.M{"$multiply": list{"$a", "$hs"}}}, "r": bson.M{"$sum": bson.M{"$multiply": list{"$r", "$hs"}}}, "hs": bson.M{"$sum": "$hs"}}}, {"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ns": "$_id.ns", "p": "$_id.p", "a": bson.M{"$divide": list{"$a", "$hs"}}, "r": bson.M{"$divide": list{"$r", "$hs"}}}}, {"$sort": bson.D{{"p", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}
		//query := []bson.M{{"$match": q}, {"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "n": "$n", "ns": "$ns", "p": "$p"}, "a": bson.M{"$sum": bson.M{"$multiply": list{"$a", "$hs"}}}, "r": bson.M{"$sum": bson.M{"$multiply": list{"$r", "$hs"}}}, "hs": bson.M{"$sum": "$hs"}}}, {"$match": bson.M{"hs": bson.M{"$gt": 0}}}, {"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ns": "$_id.ns", "p": "$_id.p", "a": bson.M{"$divide": list{"$a", "$hs"}}, "r": bson.M{"$divide": list{"$r", "$hs"}}}}, {"$sort": bson.D{{"p", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}
		err = c.Pipe(query).All(&results)
	} else if strings.ToLower(input.availabilityperiod) == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		q["a"] = bson.M{"$gte": 0}
		q["r"] = bson.M{"$gte": 0}

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
		query := []bson.M{{"$match": q}, {"$project": bson.M{"dt": 1, "a": 1, "r": 1, "p": 1, "ns": 1, "n": 1, "hs": bson.M{"$add": list{"$hs", 1}}}}, {"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "n": "$n", "ns": "$ns", "p": "$p"}, "a": bson.M{"$sum": bson.M{"$multiply": list{"$a", "$hs"}}}, "r": bson.M{"$sum": bson.M{"$multiply": list{"$r", "$hs"}}}, "hs": bson.M{"$sum": "$hs"}}}, {"$match": bson.M{"hs": bson.M{"$gt": 0}}}, {"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ns": "$_id.ns", "p": "$_id.p", "a": bson.M{"$divide": list{"$a", "$hs"}}, "r": bson.M{"$divide": list{"$r", "$hs"}}}}, {"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 6}}}, "n": "$n", "ns": "$ns", "p": "$p"}, "a": bson.M{"$avg": "$a"}, "r": bson.M{"$avg": "$r"}}}, {"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ns": "$_id.ns", "p": "$_id.p", "a": 1, "r": 1}}, {"$sort": bson.D{{"ns", 1}, {"p", 1}, {"n", 1}, {"dt", 1}}}}

		pipe := c.Pipe(query)
		err = pipe.All(&results)
		//fmt.Println(query)
	}

	if err != nil {
		return []byte("<root><error>" + err.Error() + "</error></root>")
	}

	//fmt.Println(results)
	output, err := ngis.CreateXMLResponse(results, customForm)
	if cfg.Server.Cache == true && len(results) > 0 {
		httpcache.Set("ngis "+fmt.Sprint(input), mystring(output))
	}

	return output
}
