package ngiAvailability

import (
	"encoding/xml"
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

type Availability struct {
	XMLName      xml.Name `xml:"availability" json:"-"`
	Timestamp    string   `xml:"timestamp,attr" json:"timestamp"`
	Availability string   `xml:"availability,attr" json:"availability"`
	Reliability  string   `xml:"reliability,attr" json:"reliability"`
}

type Ngi struct {
	XMLName      xml.Name `xml:"ngi" json:"-"`
	Ngi          string   `xml:"NGI,attr" json:"NGI"`
	Availability []*Availability
}

type Profile struct {
	XMLName xml.Name `xml:"Profile" json:"-"`
	Name    string   `xml:"name,attr" json:"name"`
	Ngi     []*Ngi
}

type Root struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Profile []*Profile
}

type ApiNgiAvailabilityInProfileInput struct {
	// mandatory values
	Start_time           string // UTC time in W3C format
	End_time             string // UTC time in W3C format
	Availability_profile string //availability profile
	// optional values
	Granularity    string   //availability period; possible values: `DAILY`, MONTHLY`
	Infrastructure string   //infrastructure name
	Production     string   //production or not
	Monitored      string   //yes or no
	Certification  string   //certification status
	format         string   // default XML; possible values are: XML, JSON
	Group_name     []string // site name; may appear more than once
}

type ApiNgiAvailabilityInProfileOutput struct {
	Date         string  `bson:"dt"`
	Profile      string  `bson:"ap"`
	Ngi          string  `bson:"n"`	
	Availability float64 `bson:"a"`
	Reliability  float64 `bson:"r"`
}

type list []interface{}

var CustomForm []string

func init() {
	CustomForm = []string{"20060102", "2006-01-02"} //{"Format that is returned by the database" , "Format that will be used in the generated report"}
}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

func prepareFilter(input ApiNgiAvailabilityInProfileInput) bson.M {

	ts, _ := time.Parse(zuluForm, input.Start_time)
	te, _ := time.Parse(zuluForm, input.End_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	filter := bson.M{
		"dt": bson.M{"$gte": tsYMD, "$lte": teYMD},
		"ap": input.Availability_profile,
	}

	if len(input.Group_name) > 0 {
		filter["n"] = bson.M{"$in": input.Group_name}
	}

	filter["i"] = input.Infrastructure
	filter["cs"] = input.Certification
	filter["pr"] = input.Production
	filter["m"] = input.Monitored

	//TODO: Remove hardcoded filtering and add this
	// as a parameter
	filter["sc"] = "EGI"
	filter["ss"] = "EGI"

	return filter
}

func Daily(input ApiNgiAvailabilityInProfileInput) []bson.M {
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
	// Sort by profile->ngi->site->datetime
	query := []bson.M{
		{"$match": filter},
		{"$project": bson.M{"dt": 1, "a": 1, "r": 1, "ap": 1, "n": 1, "hs": bson.M{"$add": list{"$hs", 1}}}},
		{"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "n": "$n", "ap": "$ap"},
			"a": bson.M{"$sum": bson.M{"$multiply": list{"$a", "$hs"}}}, "r": bson.M{"$sum": bson.M{"$multiply": list{"$r", "$hs"}}}, "hs": bson.M{"$sum": "$hs"}}},
		{"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ap": "$_id.ap", "a": bson.M{"$divide": list{"$a", "$hs"}},
			"r": bson.M{"$divide": list{"$r", "$hs"}}}},
		{"$sort": bson.D{{"ap", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}

	//query := []bson.M{{"$match": q}, {"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "n": "$n", "ns": "$ns", "p": "$p"}, "a": bson.M{"$sum": bson.M{"$multiply": list{"$a", "$hs"}}}, 		"r": bson.M{"$sum": bson.M{"$multiply": list{"$r", "$hs"}}}, "hs": bson.M{"$sum": "$hs"}}}, {"$match": bson.M{"hs": bson.M{"$gt": 0}}}, {"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ns": "$_id.ns", "p": 		"$_id.p", "a": bson.M{"$divide": list{"$a", "$hs"}}, "r": bson.M{"$divide": list{"$r", "$hs"}}}}, {"$sort": bson.D{{"p", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}

	return query
}

func Monthly(input ApiNgiAvailabilityInProfileInput) []bson.M {
	filter := prepareFilter(input)
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

	query := []bson.M{
		{"$match": filter}, {"$project": bson.M{"dt": 1, "a": 1, "r": 1, "ap": 1, "n": 1, "hs": bson.M{"$add": list{"$hs", 1}}}},
		{"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "n": "$n", "ap": "$ap"}, "a": bson.M{"$sum": bson.M{"$multiply": list{"$a", "$hs"}}},
			"r": bson.M{"$sum": bson.M{"$multiply": list{"$r", "$hs"}}}, "hs": bson.M{"$sum": "$hs"}}}, {"$match": bson.M{"hs": bson.M{"$gt": 0}}},
		{"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ap": "$_id.ap", "a": bson.M{"$divide": list{"$a", "$hs"}}, "r": bson.M{"$divide": list{"$r", "$hs"}}}},
		{"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 6}}}, "n": "$n", "ap": "$ap"}, "a": bson.M{"$avg": "$a"},
			"r": bson.M{"$avg": "$r"}}}, {"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ap": "$_id.ap", "a": 1, "r": 1}},
		{"$sort": bson.D{{"ap", 1}, {"n", 1}, {"dt", 1}}}}

	return query
}
