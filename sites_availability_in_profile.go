package main

import (
	"encoding/xml"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
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

type list []interface{}

func createSiteXMLResponse(results []MongoSite, customForm []string) ([]byte, error) {

	type Availability struct {
		XMLName      xml.Name `xml:"Availability"`
		Timestamp    string   `xml:"timestamp,attr"`
		Availability string   `xml:"availability,attr"`
		Reliability  string   `xml:"reliability,attr"`
	}

	type Site struct {
		Site          string `xml:"site,attr"`
		Ngi           string `xml:"NGI,attr"`
		Infastructure string `xml:"infastructure,attr"`
		Scope         string `xml:"scope,attr"`
		SiteScope     string `xml:"site_scope,attr"`
		Production    string `xml:"production,attr"`
		Monitored     string `xml:"monitored,attr"`
		CertStatus    string `xml:"certification_status,attr"`
		Availability  []*Availability
	}

	type Profile struct {
		XMLName   xml.Name `xml:"Profile"`
		Name      string   `xml:"name,attr"`
		Namespace string   `xml:"namespace,attr"`
		Site      []*Site
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Profile []*Profile
	}

	v := &Root{}

	prevProfile := ""
	prevSite := ""
	site := &Site{}
	profile := &Profile{}
	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], row.Date)

		if prevProfile != row.Profile {
			prevProfile = row.Profile
			profile = &Profile{
				Name:      row.Profile,
				Namespace: row.Namespace}
			v.Profile = append(v.Profile, profile)
			prevSite = ""
		}

		if prevSite != row.Site {
			prevSite = row.Site
			site = &Site{
				Site:          row.Site,
				Ngi:           row.Ngi,
				Infastructure: row.Infastructure,
				Scope:         row.Scope,
				SiteScope:     row.SiteScope,
				Production:    row.Production,
				Monitored:     row.Monitored,
				CertStatus:    row.CertStatus}
			profile.Site = append(profile.Site, site)
		}
		site.Availability = append(site.Availability,
			&Availability{
				Timestamp:    timestamp.Format(customForm[1]),
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability)})
		fmt.Println(v)
	}

	output, err := xml.MarshalIndent(v, " ", "  ")

	return output, err
}

//const zuluForm = "2006-01-02T15:04:05Z"
//const ymdForm = "20060102"

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
		results := []MongoSite{}
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
		output, err := createSiteXMLResponse(results, customForm)
		httpcache.Set("sites "+fmt.Sprint(input), mystring(output))
		return string(output)
	} else {
		return fmt.Sprint(out)
	}

}
