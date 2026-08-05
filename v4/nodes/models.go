package nodes

import (
	"encoding/json"
	"math"

	"github.com/ARGOeu/argo-web-api/app/reports"
)

type list []interface{}

var customForm []string

func init() {
	customForm = []string{"20060102", "2006-01-02T15:04:05Z"} //{"Format that is returned by the database" , "Format that will be used in the generated report"}
}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"
const dtForm = "2006-01-02"

type RoundedFloat float64

func (rf RoundedFloat) MarshalJSON() ([]byte, error) {
	return json.Marshal(math.Round(float64(rf)*100) / 100)
}

type basicQuery struct {
	Name         string                 `bson:"name"`
	Granularity  string                 `bson:"-"`
	Format       string                 `bson:"-"`
	StartTime    string                 `bson:"-"` // UTC time in W3C format
	EndTime      string                 `bson:"-"` // UTC time in W3C format
	Date         string                 `bson:"-"`
	StartDate    string                 `bson:"-"`
	EndDate      string                 `bson:"-"`
	Report       reports.MongoInterface `bson:"report"`
	StartTimeInt int                    `bson:"start_time"`
	EndTimeInt   int                    `bson:"end_time"`
	Vars         map[string]string      `bson:"-"`
}

// GroupInterface used to hold mongodb group information such as SITES, SERVICEGROUPS etc
type GroupInterface struct {
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

type InputStatus struct {
	startTime int
	endTime   int
	report    string
	groupType string
	group     string
	format    string
	ID        string
}

type GroupStatusData struct {
	Report           string `bson:"report"`
	Timestamp        string `bson:"timestamp"`
	Group            string `bson:"endpoint_group"`
	Status           string `bson:"status"`
	DateInteger      int    `bson:"date_integer"`
	HasThresholdRule bool   `bson:"has_threshold_rule"`
}

type NodeReport struct {
	ReportId string `bson:"report_id"`
}

type Metrics struct {
	Date         string       `json:"date,omitempty"`
	Availability RoundedFloat `json:"availability"`
	Reliability  RoundedFloat `json:"reliability"`
	Uptime       RoundedFloat `json:"uptime,omitempty"`
}

type Summary struct {
	Date         string `json:"date,omitempty"`
	Availability string `json:"availability"`
	Uptime       string `json:"uptime,omitempty"`
}

type Availability struct {
	Date         string `json:"date,omitempty"`
	Availability string `json:"availability"`
}

type Uptime struct {
	Date   string `json:"date,omitempty"`
	Uptime string `json:"uptime,omitempty"`
}

type Status struct {
	Name    string         `json:"name"`
	Results []StatusResult `json:"results"`
}

type StatusResult struct {
	Timestamp               string `json:"timestamp"`
	Value                   string `json:"value"`
	AffectedByThresholdRule bool   `json:"affected_by_threshold_rule,omitempty"`
}

type Results[T Availability | Uptime | Summary | Metrics] struct {
	Name    string `json:"name"`
	Results []T    `json:"results"`
}

type Data[T Availability | Uptime | Summary | Metrics] struct {
	Data []Results[T] `json:"data"`
}

type StatusData struct {
	Data []Status `json:"data"`
}

// errorMessage struct to hold the json/xml error response
type errorMessage struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
