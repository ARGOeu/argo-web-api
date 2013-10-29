package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"ngis"
	"services"
	"sites"
	"strconv"
	"strings"
	"time"
)

type MongoSite struct {
	SiteScope     string  "ss"
	Scope         string  "sc"
	Date          string  "dt"
	Namespace     string  "ns"
	Profile       string  "p"
	Production    string  "pr"
	Monitored     string  "m"
	Ngi           string  "n"
	Site          string  "s"
	Infastructure string  "i"
	CertStatus    string  "cs"
	Availability  float64 "a"
	Reliability   float64 "r"
}

type MongoNgi struct {
	Date         string  "dt"
	Namespace    string  "ns"
	Profile      string  "p"
	Ngi          string  "n"
	Availability float64 "a"
	Reliability  float64 "r"
}

type list []interface{}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

func ServiceAvailabilityInProfile(w http.ResponseWriter, r *http.Request) string {

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

	customForm := []string{"20060102", "2006-01-02T15:04:05Z"}

	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	out, found := httpcache.Get("service_endpoint " + fmt.Sprint(input))
	if !found {
		session, err := mgo.Dial(cfg.MongoDB.Host + ":" + fmt.Sprint(cfg.MongoDB.Port))
		if err != nil {
			panic(err)
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(cfg.MongoDB.Db).C("timelines")
		results := []services.Timeline{}
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
			q["ns"] = bson.M{"$in": input.service_flavour}
		}

		if len(input.service_hostname) > 0 {
			q["h"] = bson.M{"$in": input.service_hostname}
		}
		query := []bson.M{{"$match": q}, {"$sort": bson.D{{"p", 1}, {"h", 1}, {"sf", 1}, {"d", 1}}}}
		err = c.Pipe(query).All(&results)

		//err = c.Find(q).Sort("p", "h", "sf").All(&results)
		if err != nil {
			return ("<root><error>" + err.Error() + "</error></root>")
		}

		//rootfmt.Println(results)
		output, err := services.CreateXMLResponse(results, customForm)
		httpcache.Set("service_endpoint "+fmt.Sprint(input), mystring(output))
		return string(output)

	} else {
		return fmt.Sprint(out)
	}
}

func SitesAvailabilityInProfile(w http.ResponseWriter, r *http.Request) string {

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
	customForm := []string{"20060102", "2006-01-02T15:04:05Z"} //{"Format that is returned by the database" , "Format that will be used in the generated report"}

	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	out, found := httpcache.Get("sites " + fmt.Sprint(input))
	if !found {

		session, err := mgo.Dial("127.0.0.1")
		if err != nil {
			panic(err)
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB("AR").C("sites")
		results := []sites.MongoSite{}
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

		if len(input.availabilityperiod) == 0 || strings.ToLower(input.availabilityperiod) == "daily" {
			customForm[0] = "20060102"
			customForm[1] = "2006-01-02"
			err = c.Pipe([]bson.M{{"$match": q}, {"$project": bson.M{"dt": bson.M{"$substr": list{"$dt", 0, 8}}, "i": 1, "sc": 1, "ss": 1, "n": 1, "pr": 1, "m": 1, "cs": 1, "ns": 1, "s": 1, "p": 1, "a": 1, "r": 1}}, {"$sort": bson.D{{"p", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}).All(&results)
			//fmt.Println(len(results))

		} else if strings.ToLower(input.availabilityperiod) == "monthly" {
			customForm[0] = "200601"
			customForm[1] = "2006-01"
			query := []bson.M{{"$match": bson.M{"a": bson.M{"$gte": 0}, "r": bson.M{"$gte": 0}, "i": "Production", "cs": "Certified", "pr": "Y", "m": "Y", "dt": bson.M{"$gte": tsYMD, "$lte": teYMD}, "p": bson.M{"$in": input.profile_name}}}, {"$group": bson.M{"_id": bson.M{"dt": bson.M{"$substr": list{"$dt", 0, 6}}, "i": "$i", "sc": "$sc", "ss": "$ss", "n": "$n", "pr": "$pr", "m": "$m", "cs": "$cs", "ns": "$ns", "s": "$s", "p": "$p"}, "avgup": bson.M{"$avg": "$up"}, "avgu": bson.M{"$avg": "$u"}, "avgd": bson.M{"$avg": "$d"}}}, {"$project": bson.M{"dt": "$_id.dt", "i": "$_id.i", "sc": "$_id.sc", "ss": "$_id.ss", "n": "$_id.n", "pr": "$_id.pr", "m": "$_id.m", "cs": "$_id.cs", "ns": "$_id.ns", "s": "$_id.s", "p": "$_id.p", "avgup": 1, "avgu": 1, "avgd": 1, "a": bson.M{"$multiply": list{bson.M{"$divide": list{"$avgup", bson.M{"$subtract": list{1.00000001, "$avgu"}}}}, 100}}, "r": bson.M{"$multiply": list{bson.M{"$divide": list{"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgu"}}, "$avgd"}}}}, 100}}}}, {"$sort": bson.D{{"ns", 1}, {"p", 1}, {"n", 1}, {"c", 1}, {"dt", 1}}}}

			pipe := c.Pipe(query)
			err = pipe.All(&results)
			fmt.Println(query)
		}

		if err != nil {
			return ("<root><error>" + err.Error() + "</error></root>")
		}

		fmt.Println(len(results))
		output, err := sites.CreateXMLResponse(results, customForm)
		httpcache.Set("sites "+fmt.Sprint(input), mystring(output))
		return string(output)
	} else {
		return fmt.Sprint(out)
	}

}

func NgiAvailabilityInProfile(w http.ResponseWriter, r *http.Request) string {

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

	out, found := httpcache.Get("ngi " + fmt.Sprint(input))
	if !found {

		session, err := mgo.Dial("127.0.0.1")
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
			// TODO: We do not have the ngi name in the timeline
		}

		fmt.Println(input)

		if len(input.availabilityperiod) == 0 || strings.ToLower(input.availabilityperiod) == "daily" {
			customForm[0] = "20060102"
			customForm[1] = "2006-01-02"
			query := []bson.M{{"$match": q}, {"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "n": "$n", "ns": "$ns", "p": "$p"}, "a": bson.M{"$sum": bson.M{"$multiply": list{"$a", "$hs"}}}, "r": bson.M{"$sum": bson.M{"$multiply": list{"$r", "$hs"}}}, "hs": bson.M{"$sum": "$hs"}}}, {"$match": bson.M{"hs": bson.M{"$gt": 0}}}, {"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ns": "$_id.ns", "p": "$_id.p", "a": bson.M{"$divide": list{"$a", "$hs"}}, "r": bson.M{"$divide": list{"$r", "$hs"}}}}, {"$sort": bson.D{{"p", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}
			err = c.Pipe(query).All(&results)
			//err = c.Find(q).Sort("p", "n", "s", "dt").All(&results)
			//fmt.Println(q)
			fmt.Println(query)

		} else if strings.ToLower(input.availabilityperiod) == "monthly" {
			customForm[0] = "200601"
			customForm[1] = "2006-01"
			q["a"] = bson.M{"$gte": 0}
			q["r"] = bson.M{"$gte": 0}

			query := []bson.M{{"$match": q}, {"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "n": "$n", "ns": "$ns", "p": "$p"}, "a": bson.M{"$sum": bson.M{"$multiply": list{"$a", "$hs"}}}, "r": bson.M{"$sum": bson.M{"$multiply": list{"$r", "$hs"}}}, "hs": bson.M{"$sum": "$hs"}}}, {"$match": bson.M{"hs": bson.M{"$gt": 0}}}, {"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ns": "$_id.ns", "p": "$_id.p", "a": bson.M{"$divide": list{"$a", "$hs"}}, "r": bson.M{"$divide": list{"$r", "$hs"}}}}, {"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 6}}}, "n": "$n", "ns": "$ns", "p": "$p"}, "a": bson.M{"$avg": "$a"}, "r": bson.M{"$avg": "$r"}}}, {"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ns": "$_id.ns", "p": "$_id.p", "a": 1, "r": 1}}, {"$sort": bson.D{{"ns", 1}, {"p", 1}, {"n", 1}, {"dt", 1}}}}

			pipe := c.Pipe(query)
			err = pipe.All(&results)
			fmt.Println(query)
		}

		if err != nil {
			return ("<root><error>" + err.Error() + "</error></root>")
		}

		//fmt.Println(results)
		output, err := ngis.CreateXMLResponse(results, customForm)
		httpcache.Set("ngi "+fmt.Sprint(input), mystring(output))
		return string(output)
	} else {
		return fmt.Sprint(out)
	}
}
