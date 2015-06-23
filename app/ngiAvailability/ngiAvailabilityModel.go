package ngiAvailability

import (
	"encoding/xml"
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

type Availability struct {
	XMLName      xml.Name `xml:"Availability" json:"-"`
	Timestamp    string   `xml:"timestamp,attr" json:"timestamp"`
	Availability string   `xml:"availability,attr" json:"availability"`
	Reliability  string   `xml:"reliability,attr" json:"reliability"`
}

type SuperGroup struct {
	XMLName      xml.Name `xml:"SuperGroup" json:"-"`
	SuperGroup   string   `xml:"SuperGroup,attr" json:"SuperGroup"`
	Availability []*Availability
}

type Job struct {
	XMLName    xml.Name `xml:"Profile" json:"-"`
	Name       string   `xml:"name,attr" json:"name"`
	SuperGroup []*SuperGroup
}

type Root struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Job     []*Job
}

type ApiSuperGroupAvailabilityInProfileInput struct {
	// mandatory values
	Start_time string // UTC time in W3C format
	End_time   string // UTC time in W3C format
	//Availability_profile string //availability profile
	Job string //unique id that represents the current job which produces the ar result.
	// optional values
	Granularity string //availability period; possible values: `DAILY`, MONTHLY`
	//Infrastructure string   //infrastructure name
	//Production     string   //production or not
	//Monitored      string   //yes or no
	//Certification  string   //certification status
	format     string   // default XML; possible values are: XML, JSON
	Group_name []string // site name; may appear more than once
}

type ApiSuperGroupAvailabilityInProfileOutput struct {
	Date         string  `bson:"date"`
	Job          string  `bson:"job"`
	SuperGroup   string  `bson:"supergroup"`
	Availability float64 `bson:"availability"`
	Reliability  float64 `bson:"reliability"`
}

type list []interface{}

var CustomForm []string

func init() {
	CustomForm = []string{"20060102", "2006-01-02"} //{"Format that is returned by the database" , "Format that will be used in the generated report"}
}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

func prepareFilter(input ApiSuperGroupAvailabilityInProfileInput) bson.M {

	ts, _ := time.Parse(zuluForm, input.Start_time)
	te, _ := time.Parse(zuluForm, input.End_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	filter := bson.M{
		"date": bson.M{"$gte": tsYMD, "$lte": teYMD},
		"job":  input.Job,
		//"ap": input.Availability_profile,
	}

	if len(input.Group_name) > 0 {
		filter["supergroup"] = bson.M{"$in": input.Group_name}
	}

	// filter["i"] = input.Infrastructure
	// filter["cs"] = input.Certification
	// filter["pr"] = input.Production
	// filter["m"] = input.Monitored
	//
	// filter["sc"] = "EGI"
	// filter["ss"] = "EGI"

	return filter
}

func Daily(input ApiSuperGroupAvailabilityInProfileInput) []bson.M {
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
		{"$project": bson.M{"date": 1, "availability": 1, "reliability": 1, "job": 1, "supergroup": 1, "weights": bson.M{"$add": list{"$weights", 1}}}},
		{"$group": bson.M{"_id": bson.M{"date": bson.D{{"$substr", list{"$date", 0, 8}}}, "supergroup": "$supergroup", "job": "$job"},
			"availability": bson.M{"$sum": bson.M{"$multiply": list{"$availability", "$weights"}}}, "reliability": bson.M{"$sum": bson.M{"$multiply": list{"$reliability", "$weights"}}}, "weights": bson.M{"$sum": "$weights"}}},
		{"$project": bson.M{"date": "$_id.date", "supergroup": "$_id.supergroup", "job": "$_id.job", "availability": bson.M{"$divide": list{"$availability", "$weights"}},
			"reliability": bson.M{"$divide": list{"$reliability", "$weights"}}}},
		{"$sort": bson.D{{"job", 1}, {"supergroup", 1}, {"site", 1}, {"date", 1}}}}

	//query := []bson.M{{"$match": q}, {"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "n": "$n", "ns": "$ns", "p": "$p"}, "a": bson.M{"$sum": bson.M{"$multiply": list{"$a", "$hs"}}}, 		"r": bson.M{"$sum": bson.M{"$multiply": list{"$r", "$hs"}}}, "hs": bson.M{"$sum": "$hs"}}}, {"$match": bson.M{"hs": bson.M{"$gt": 0}}}, {"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ns": "$_id.ns", "p": 		"$_id.p", "a": bson.M{"$divide": list{"$a", "$hs"}}, "r": bson.M{"$divide": list{"$r", "$hs"}}}}, {"$sort": bson.D{{"p", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}

	return query
}

func Monthly(input ApiSuperGroupAvailabilityInProfileInput) []bson.M {
	filter := prepareFilter(input)
	//PROBABLY THIS LEADS TO THE SAME BUG WE RAN INTO WITH SITES. MUST BE INVESTIGATED!!!!!!!!!!!!
	filter["availability"] = bson.M{"$gte": 0}
	filter["reliability"] = bson.M{"$gte": 0}

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
		{"$match": filter}, {"$project": bson.M{"date": 1, "availability": 1, "reliability": 1, "job": 1, "supergroup": 1, "weights": bson.M{"$add": list{"$weights", 1}}}},
		{"$group": bson.M{"_id": bson.M{"date": bson.D{{"$substr", list{"$date", 0, 8}}}, "supergroup": "$supergroup", "job": "$job"}, "availability": bson.M{"$sum": bson.M{"$multiply": list{"$availability", "$weights"}}},
			"reliability": bson.M{"$sum": bson.M{"$multiply": list{"$reliability", "$weights"}}}, "weights": bson.M{"$sum": "$weights"}}}, {"$match": bson.M{"weights": bson.M{"$gt": 0}}},
		{"$project": bson.M{"date": "$_id.date", "supergroup": "$_id.supergroup", "job": "$_id.job", "availability": bson.M{"$divide": list{"$availability", "$weights"}}, "reliability": bson.M{"$divide": list{"$reliability", "$weights"}}}},
		{"$group": bson.M{"_id": bson.M{"date": bson.D{{"$substr", list{"$date", 0, 6}}}, "supergroup": "$supergroup", "job": "$job"}, "availability": bson.M{"$avg": "$availability"},
			"reliability": bson.M{"$avg": "$reliability"}}}, {"$project": bson.M{"date": "$_id.date", "supergroup": "$_id.supergroup", "job": "$_id.job", "availability": 1, "reliability": 1}},
		{"$sort": bson.D{{"job", 1}, {"supergroup", 1}, {"date", 1}}}}

	return query
}
