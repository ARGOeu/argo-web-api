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
	"api/services"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strconv"
	"time"
)


//Reply to requests about service_availability_in_profile
func ServiceAvailabilityInProfile(w http.ResponseWriter, r *http.Request) []byte {

	// This is the input we will receive from the API

	type ApiServiceAvailabilityInProfileInput struct {
		// mandatory values
		start_time          string   // UTC time in W3C format
		end_time            string   // UTC time in W3C format
		vo_name             []string // may appear more than once. (eg: ops)
		profile_name        []string // may appear more than once. (eg: CMS_CRITICAL)
		group_type          []string // may appear more than once. (eg: CMS_Site)
		availability_period string   // availability period; possible values: 'HOURLY', 'DAILY', 'WEEKLY', 'MONTHLY'
		// optional values
		output           string   // default XML; possible values are: XML, JSON
		namespace        []string // profile namespace; may appear more than once. (eg: ch.cern.sam)
		group_name       []string // site name; may appear more than once
		service_flavour  []string // service flavour name; may appear more than once. (eg: SRMv2)
		service_hostname []string // service hostname; may appear more than once. (eg: ce202.cern.ch)
	}

	// Parse the request into the input
	urlValues := r.URL.Query()
	input := ApiServiceAvailabilityInProfileInput{
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues["vo_name"],
		urlValues["profile_name"],
		urlValues["group_type"],
		urlValues.Get("type"),
		urlValues.Get("output"),
		urlValues["namespace"],
		urlValues["group_name"],
		urlValues["service_flavour"],
		urlValues["service_hostname"],
	}

	// Parse the date range of the query
	customForm := []string{"20060102", "2006-01-02T15:04:05Z"}
	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	// If caching is enabled search the cache for matches
	if cfg.Server.Cache == true {
		out, found := httpcache.Get("service_endpoint " + fmt.Sprint(input))
		if found {
			return []byte(fmt.Sprint(out))
		}
	}

	// Create a mongodb session
	session, err := mgo.Dial(cfg.MongoDB.Host + ":" + fmt.Sprint(cfg.MongoDB.Port))
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(cfg.MongoDB.Db).C("timelines")
	results := []services.Timeline{}

	// Construct the query to mongodb based on the input
	q := bson.M{
		"d":  bson.M{"$gte": tsYMD, "$lte": teYMD},
		"vo": bson.M{"$in": input.vo_name},
		"p":  bson.M{"$in": input.profile_name},
	}

	if len(input.namespace) > 0 {
		q["ns"] = bson.M{"$in": input.namespace}
	}

	if len(input.group_name) > 0 {
		// TODO: We do not have the site name in the timeline
	}

	if len(input.service_flavour) > 0 {
		q["sf"] = bson.M{"$in": input.service_flavour}
	}

	if len(input.service_hostname) > 0 {
		q["h"] = bson.M{"$in": input.service_hostname}
	}
	query := []bson.M{{"$match": q}, {"$sort": bson.D{{"p", 1}, {"h", 1}, {"sf", 1}, {"d", 1}}}}
	err = c.Pipe(query).All(&results)

	//err = c.Find(q).Sort("p", "h", "sf").All(&results)
	if err != nil {
		return []byte("<root><error>" + err.Error() + "</error></root>")
	}

	//rootfmt.Println(results)
	output, err := services.CreateXMLResponse(results, customForm)

	//if caching is enabled save the result to the cache
	if cfg.Server.Cache == true && len(results) > 0 {
		httpcache.Set("service_endpoint "+fmt.Sprint(input), mystring(output))
	}
	return output
}
