package groupGroupsAvailability

import (
	"encoding/xml"
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

// Availability struct to represent the availability-reliability results
type Availability struct {
	XMLName      xml.Name `xml:"Availability" json:"-"`
	Timestamp    string   `xml:"timestamp,attr" json:"timestamp"`
	Availability string   `xml:"availability,attr" json:"availability"`
	Reliability  string   `xml:"reliability,attr" json:"reliability"`
}

// SuperGroup struct to hold the availability-reliability results for each group
type SuperGroup struct {
	XMLName      xml.Name `xml:"SuperGroup" json:"-"`
	SuperGroup   string   `xml:"name,attr" json:"name"`
	Availability []*Availability
}

// Job struct to hold all SuperGroups related with this job
type Job struct {
	XMLName    xml.Name `xml:"Job" json:"-"`
	Name       string   `xml:"name,attr" json:"name"`
	SuperGroup []*SuperGroup
}

// Root struct to represent the root of the XML document
type Root struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Job     []*Job
}

// ApiSuperGroupAvailabilityInProfileInput struct to represent the api call input parameters
type ApiSuperGroupAvailabilityInProfileInput struct {
	// mandatory values
	StartTime string // UTC time in W3C format
	EndTime   string // UTC time in W3C format
	Job       string //unique id that represents the current job which produces the ar result.
	// optional values
	Granularity string   //availability period; possible values: `DAILY`, MONTHLY`
	format      string   // default XML; possible values are: XML, JSON
	GroupName   []string // site name; may appear more than once
}

// ApiSuperGroupAvailabilityInProfileOutput to represent db data retrieval
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

	ts, _ := time.Parse(zuluForm, input.StartTime)
	te, _ := time.Parse(zuluForm, input.EndTime)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	filter := bson.M{
		"date": bson.M{"$gte": tsYMD, "$lte": teYMD},
		"job":  input.Job,
	}

	if len(input.GroupName) > 0 {
		filter["supergroup"] = bson.M{"$in": input.GroupName}
	}

	return filter
}

// Daily function to build the MongoDB aggregation query for daily calculations
func Daily(input ApiSuperGroupAvailabilityInProfileInput) []bson.M {
	filter := prepareFilter(input)
	// Mongo aggregation pipeline
	// Select all the records that match q
	// Project the results to add 1 to every weights to avoid having 0 as a weights
	// Group them by the first 8 digits of datetime (YYYYMMDD) and each group find
	// availability = sum(availability*weights)
	// reliability = sum(reliability*weights)
	// weights = sum(weights)
	// Project to a better format and do these computations
	// availability = availability/weights
	// reliability = reliability/weights
	// Sort by job->supergroup->name->datetime
	query := []bson.M{
		{"$match": filter},
		{"$project": bson.M{"date": 1, "availability": 1, "reliability": 1, "job": 1, "supergroup": 1, "weights": bson.M{"$add": list{"$weights", 1}}}},
		{"$group": bson.M{"_id": bson.M{"date": bson.D{{"$substr", list{"$date", 0, 8}}}, "supergroup": "$supergroup", "job": "$job"},
			"availability": bson.M{"$sum": bson.M{"$multiply": list{"$availability", "$weights"}}}, "reliability": bson.M{"$sum": bson.M{"$multiply": list{"$reliability", "$weights"}}}, "weights": bson.M{"$sum": "$weights"}}},
		{"$project": bson.M{"date": "$_id.date", "supergroup": "$_id.supergroup", "job": "$_id.job", "availability": bson.M{"$divide": list{"$availability", "$weights"}},
			"reliability": bson.M{"$divide": list{"$reliability", "$weights"}}}},
		{"$sort": bson.D{{"job", 1}, {"supergroup", 1}, {"name", 1}, {"date", 1}}}}

	//query := []bson.M{{"$match": q}, {"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "n": "$n", "ns": "$ns", "p": "$p"}, "a": bson.M{"$sum": bson.M{"$multiply": list{"$a", "$hs"}}}, 		"r": bson.M{"$sum": bson.M{"$multiply": list{"$r", "$hs"}}}, "hs": bson.M{"$sum": "$hs"}}}, {"$match": bson.M{"hs": bson.M{"$gt": 0}}}, {"$project": bson.M{"dt": "$_id.dt", "n": "$_id.n", "ns": "$_id.ns", "p": 		"$_id.p", "a": bson.M{"$divide": list{"$a", "$hs"}}, "r": bson.M{"$divide": list{"$r", "$hs"}}}}, {"$sort": bson.D{{"p", 1}, {"n", 1}, {"s", 1}, {"dt", 1}}}}
	return query
}

// Monthly function to build the MongoDB aggregation query for monthly calculations
func Monthly(input ApiSuperGroupAvailabilityInProfileInput) []bson.M {
	filter := prepareFilter(input)
	//PROBABLY THIS LEADS TO THE SAME BUG WE RAN INTO WITH SITES. MUST BE INVESTIGATED!!!!!!!!!!!!
	filter["availability"] = bson.M{"$gte": 0}
	filter["reliability"] = bson.M{"$gte": 0}

	// Mongo aggregation pipeline
	// Select all the records that match q
	// Project the results to add 1 to every weights to avoid having 0 as a weights
	// Group them by the first 8 digits of datetime (YYYYMMDD) and each group find
	// availability = sum(availability*weights)
	// reliability = sum(reliability*weights)
	// weights = sum(weights)
	// Project to a better format and do these computations
	// availability = availability/weights
	// reliability = reliability/weights
	// Group by the first 6 digits of the datetime (YYYYMM) and by ngi,site,profile and for each group find
	// availability = average(availability)
	// reliability = average(reliability)
	// Project the results to a better format
	// Sort by namespace->job->supergroup->datetime

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
