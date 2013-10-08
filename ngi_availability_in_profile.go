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

type MongoNgi struct {
	Date         string  "dt"
	Namespace    string  "ns"
	Profile      string  "p"
	Ngi          string  "n"
	Availability float64 "a"
	Reliability  float64 "r"
}

//type list []interface{}

func createNgiXMLResponse(results []MongoNgi, customForm []string) ([]byte, error) {

	type Availability struct {
		XMLName      xml.Name `xml:"Availability"`
		Timestamp    string   `xml:"timestamp,attr"`
		Availability string   `xml:"availability,attr"`
		Reliability  string   `xml:"reliability,attr"`
	}

	type Ngi struct {
		Ngi          string `xml:"NGI,attr"`
		Availability []*Availability
	}

	type Profile struct {
		XMLName   xml.Name `xml:"Profile"`
		Name      string   `xml:"name,attr"`
		Namespace string   `xml:"namespace,attr"`
		Ngi       []*Ngi
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Profile []*Profile
	}

	v := &Root{}

	prevProfile := ""
	prevNgi := ""
	ngi := &Ngi{}
	profile := &Profile{}
	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], row.Date)

		if prevProfile != row.Profile {
			prevProfile = row.Profile
			profile = &Profile{
				Name:      row.Profile,
				Namespace: row.Namespace}
			v.Profile = append(v.Profile, profile)
			prevNgi = ""
		}

		if prevNgi != row.Ngi {
			prevNgi = row.Ngi
			ngi = &Ngi{
				Ngi: row.Ngi,
			}
			profile.Ngi = append(profile.Ngi, ngi)
		}
		ngi.Availability = append(ngi.Availability,
			&Availability{
				Timestamp:    timestamp.Format(customForm[1]),
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability)})
	}

	output, err := xml.MarshalIndent(v, " ", "  ")

	return output, err
}

//const zuluForm = "2006-01-02T15:04:05Z"
//const ymdForm = "20060102"

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
		results := []MongoNgi{}
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
		output, err := createNgiXMLResponse(results, customForm)
		httpcache.Set("ngi "+fmt.Sprint(input), mystring(output))
		return string(output)
	} else {
		return fmt.Sprint(out)
	}
}
