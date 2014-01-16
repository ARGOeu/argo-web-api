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
	"api/sites"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func SitesAvailabilityInProfile(w http.ResponseWriter, r *http.Request) []byte {

	// This is the input we will receive from the API

	type ApiSiteAvailabilityInProfileInput struct {
		// mandatory values
		start_time         string   // UTC time in W3C format
		end_time           string   // UTC time in W3C format
		profile_name       []string // may appear more than once. (eg: CMS_CRITICAL)
		group_type         []string // may appear more than once. (eg: CMS_Site)
		availabilityperiod string   // availability period; possible values: `HOURLY`, `DAILY`, `WEEKLY`, `MONTHLY`
		// optional values
		output     string   // default XML; possible values are: XML, JSON
		namespace  []string // profile namespace; may appear more than once. (eg: ch.cern.sam)
		group_name []string // site name; may appear more than once
	}

	// Parse the request into the input
	urlValues := r.URL.Query()
	input := ApiSiteAvailabilityInProfileInput{
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues["profile_name"],
		urlValues["group_type"],
		urlValues.Get("type"),
		urlValues.Get("output"),
		urlValues["namespace"],
		urlValues["group_name"],
	}

	// Parse the date range of the query
	customForm := []string{"20060102", "2006-01-02T15:04:05Z"} //{"Format that is returned by the database" , "Format that will be used in the generated report"}
	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	// If caching is enabled search the cache for matches
	if cfg.Server.Cache == true {
		out, found := httpcache.Get("sites " + fmt.Sprint(input))
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
	results := []sites.MongoSite{}

	// Construct the query to mongodb based on the input
	q := bson.M{
		"dt": bson.M{"$gte": tsYMD, "$lte": teYMD},
		"p":  bson.M{"$in": input.profile_name},
	}

	if len(input.namespace) > 0 {
		q["ns"] = bson.M{"$in": input.namespace}
	}

	if len(input.group_name) > 0 {
		// TODO: We do not have the site name in the timeline
	}

	q["i"] = "Production"
	q["cs"] = "Certified"
	q["pr"] = "Y"
	q["m"] = "Y"

	// Select the granularity of the search daily/monthly
	if len(input.availabilityperiod) == 0 || strings.ToLower(input.availabilityperiod) == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		// Mongo aggregation pipeline
		// Select all the records that match q
		// Project to select just the first 8 digits of the date YYYYMMDD
		// Sort by profile->ngi->site->datetime
		err = c.Pipe([]bson.M{{"$match": q}, {"$project": bson.M{"dt": bson.M{"$substr": list{"$dt", 0, 8}}, "i": 1, "sc": 1, "ss": 1, "n": 1, "pr": 1, "m": 1, "cs": 1, "ns": 1, "s": 1, "p": 1, "a": 1, "r": 1}}, {"$sort": bson.D{{"p", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}).All(&results)
		//fmt.Println(len(results))

	} else if strings.ToLower(input.availabilityperiod) == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		q["a"] = bson.M{"$gte": 0}
		q["r"] = bson.M{"$gte": 0}

		// Mongo aggregation pipeline
		// Select all the records that match q
		// Group them by the first six digits of their date (YYYYMM), their ngi, their site, their profile, etc...
		// from that group find the average of the uptime, u, downtime
		// Project the result to a better format and do this computation
		// availability = (avgup/(1.00000001 - avgu))*100
		// reliability = (avgup/((1.00000001 - avgu)-avgd))*100
		// Sort the results by namespace->profile->ngi->site->datetime
		query := []bson.M{{"$match": q}, {"$group": bson.M{"_id": bson.M{"dt": bson.M{"$substr": list{"$dt", 0, 6}}, "i": "$i", "sc": "$sc", "ss": "$ss", "n": "$n", "pr": "$pr", "m": "$m", "cs": "$cs", "ns": "$ns", "s": "$s", "p": "$p"}, "avgup": bson.M{"$avg": "$up"}, "avgu": bson.M{"$avg": "$u"}, "avgd": bson.M{"$avg": "$d"}}}, {"$project": bson.M{"dt": "$_id.dt", "i": "$_id.i", "sc": "$_id.sc", "ss": "$_id.ss", "n": "$_id.n", "pr": "$_id.pr", "m": "$_id.m", "cs": "$_id.cs", "ns": "$_id.ns", "s": "$_id.s", "p": "$_id.p", "avgup": 1, "avgu": 1, "avgd": 1, "a": bson.M{"$multiply": list{bson.M{"$divide": list{"$avgup", bson.M{"$subtract": list{1.00000001, "$avgu"}}}}, 100}}, "r": bson.M{"$multiply": list{bson.M{"$divide": list{"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgu"}}, "$avgd"}}}}, 100}}}}, {"$sort": bson.D{{"ns", 1}, {"p", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}

		pipe := c.Pipe(query)
		err = pipe.All(&results)
		fmt.Println(query)
	}

	if err != nil {
		return []byte("<root><error>" + err.Error() + "</error></root>")
	}

	fmt.Println(len(results))
	output, err := sites.CreateXMLResponse(results, customForm)
	if cfg.Server.Cache == true && len(results) > 0 {
		httpcache.Set("sites "+fmt.Sprint(input), mystring(output))
	}

	return output

}