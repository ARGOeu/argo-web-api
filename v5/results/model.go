package resultsV5

import (
	"encoding/json"
	"math"

	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/respond"
	"go.mongodb.org/mongo-driver/bson"
)

type list []interface{}

var customForm []string

func init() {
	customForm = []string{"20060102", "2006-01-02T15:04:05Z"} //{"Format that is returned by the database" , "Format that will be used in the generated report"}
}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

type basicQuery struct {
	Name         string                 `bson:"name"`
	Granularity  string                 `bson:"-"`
	Format       string                 `bson:"-"`
	StartTime    string                 `bson:"-"` // UTC time in W3C format
	EndTime      string                 `bson:"-"` // UTC time in W3C format
	Report       reports.MongoInterface `bson:"report"`
	StartTimeInt int                    `bson:"start_time"`
	EndTimeInt   int                    `bson:"end_time"`
	Vars         map[string]string      `bson:"-"`
}

// ServiceFlavor struct for formating json
type ServiceFlavor struct {
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	Availability []interface{} `json:"results"`
}

type endpointGroupResultQuery struct {
	basicQuery
	Group string `bson:"supergroup"`
}

// EndpointGroupInterface for mongodb object exchanging
type EndpointGroupInterface struct {
	Name         string  `bson:"name"`
	Report       string  `bson:"report"`
	Date         string  `bson:"date"`
	Type         string  `bson:"type"`
	Up           float64 `bson:"up"`
	Down         float64 `bson:"down"`
	Unknown      float64 `bson:"unknown"`
	Availability float64 `bson:"availability"`
	Reliability  float64 `bson:"reliability"`
	Weights      string  `bson:"weights"`
	SuperGroup   string  `bson:"supergroup"`
}

// SuperGroupInterface for mongodb object exchanging
type SuperGroupInterface struct {
	Report       string  `bson:"report"`
	Date         string  `bson:"date"`
	Type         string  `bson:"type"`
	Up           float64 `bson:"uptime"`
	Down         float64 `bson:"downtime"`
	Unknown      float64 `bson:"unknown"`
	Availability float64 `bson:"availability"`
	Reliability  float64 `bson:"reliability"`
	Weights      string  `bson:"weight"`
	SuperGroup   string  `bson:"supergroup"`
}

type RoundedFloat float64

func (rf RoundedFloat) MarshalJSON() ([]byte, error) {
	return json.Marshal(math.Round(float64(rf)*100) / 100)
}

type ShortResult struct {
	Timestamp    string       `json:"timestamp,omitempty"`
	Availability RoundedFloat `json:"availability"`
	Reliability  RoundedFloat `json:"reliability"`
}

type Result struct {
	Timestamp    string       `json:"timestamp,omitempty"`
	Availability RoundedFloat `json:"availability"`
	Reliability  RoundedFloat `json:"reliability"`
	Unknown      RoundedFloat `json:"unknown"`
	Uptime       RoundedFloat `json:"uptime"`
	Downtime     RoundedFloat `json:"downtime"`
}

// EndpointInterface for mongodb object exchanging
type EndpointInterface struct {
	Name         string            `bson:"name"`
	Report       string            `bson:"report"`
	Date         string            `bson:"date"`
	Type         string            `bson:"type"`
	Up           float64           `bson:"up"`
	Down         float64           `bson:"down"`
	Unknown      float64           `bson:"unknown"`
	Availability float64           `bson:"availability"`
	Reliability  float64           `bson:"reliability"`
	SuperGroup   string            `bson:"supergroup"`
	Service      string            `bson:"service"`
	Info         map[string]string `bson:"info"`
}

type serviceFlavorResultQuery struct {
	basicQuery
	EndpointGroup string `bson:"supergroup"`
}

// ServiceFlavorInterface for mongodb object exchanging
type ServiceFlavorInterface struct {
	Name         string  `bson:"name"`
	Report       string  `bson:"report"`
	Date         string  `bson:"date"`
	Type         string  `bson:"type"`
	Up           float64 `bson:"up"`
	Down         float64 `bson:"down"`
	Unknown      float64 `bson:"unknown"`
	Availability float64 `bson:"availability"`
	Reliability  float64 `bson:"reliability"`
	SuperGroup   string  `bson:"supergroup"`
}

// ServuceEndpointGroup struct listing included endpoints in a service for formating xjson
type ServiceEndpointGroup struct {
	Name      string        `json:"name"`
	Type      string        `json:"type"`
	Endpoints []interface{} `json:"endpoints"`
}

// ServiceFlavorGroup struct for formating json
type ServiceFlavorGroup struct {
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	ServiceFlavor []interface{} `json:"service-types"`
}

type endpointResultQuery struct {
	basicQuery
	EndpointGroup string `bson:"supergroup"`
	Service       string `bson:"service"`
}

// Endpoint A/R struct for formating json
type Endpoint struct {
	Name       string            `json:"name"`
	Service    string            `json:"service,omitempty"`
	SuperGroup string            `json:"supergroup,omitempty"`
	Type       string            `json:"type"`
	Info       map[string]string `json:"info,omitempty"`
	Results    []interface{}     `json:"results"`
}

// Group struct for formating xml/json
type Group struct {
	Name    string        `json:"name"`
	Type    string        `json:"type"`
	Results []interface{} `json:"results"`
}

// SuperGroup struct for formating xml/json
type SuperGroup struct {
	Name    string        `xml:"name,attr" json:"name"`
	Type    string        `xml:"type,attr" json:"type"`
	Groups  []interface{} `json:"groups,omitempty"`
	Results []interface{} `json:"results,omitempty"`
}

type pageRoot struct {
	Result    []interface{} `json:"results"`
	PageToken string        `json:"nextPageToken,omitempty"`
	PageSize  int           `json:"pageSize,omitempty"`
}

type root struct {
	Result []interface{} `json:"results"`
}

// errorMessage struct to hold the json/xml error response
type errorMessage struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// ErrorResponse shortcut to respond.ErrorResponse
type ErrorResponse respond.ErrorResponse

// DailyServiceFlavor query to aggregate daily SF results from mongoDB
func DailyServiceFlavor(filter bson.M) []bson.M {
	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"date":         bson.D{{Key: "$substr", Value: list{"$date", 0, 8}}},
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
			{Key: "supergroup", Value: 1},
			{Key: "name", Value: 1},
			{Key: "date", Value: 1}}}}

	return query
}

// MonthlyServiceFlavor query to aggregate daily SF results from mongoDB
func MonthlyServiceFlavor(filter bson.M) []bson.M {
	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"date":       bson.D{{Key: "$substr", Value: list{"$date", 0, 6}}},
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
			{Key: "supergroup", Value: 1},
			{Key: "name", Value: 1},
			{Key: "date", Value: 1}}}}
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
			{Key: "supergroup", Value: 1},
			{Key: "name", Value: 1}}}}
	return query
}
