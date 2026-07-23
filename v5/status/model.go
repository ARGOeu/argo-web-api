package statusV5

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

// ServiceInputParams struct holds as input all the url params of the request
type ServiceInputParams struct {
	startTime int // UTC time in W3C format
	endTime   int
	report    string
	group     string
	service   string
	format    string
}

// ServiceDataOutput struct holds the queried data from datastore
type ServiceDataOutput struct {
	Report           string `bson:"report"`
	Timestamp        string `bson:"timestamp"`
	Group            string `bson:"endpoint_group"`
	Service          string `bson:"service"`
	Status           string `bson:"status"`
	DateInteger      int    `bson:"date_integer"`
	HasThresholdRule bool   `bson:"has_threshold_rule"`
}

// InputParams struct holds as input all the url params of the request
type GroupInputParams struct {
	startTime int // UTC time in W3C format
	endTime   int
	report    string
	group     string
	format    string
}

// errorMessage struct to hold the json/xml error response
type errorMessage struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// DataOutput struct holds the queried data from datastore
type GroupDataOutput struct {
	Report           string `bson:"report"`
	Timestamp        string `bson:"timestamp"`
	Group            string `bson:"endpoint_group"`
	Status           string `bson:"status"`
	DateInteger      int    `bson:"date_integer"`
	HasThresholdRule bool   `bson:"has_threshold_rule"`
}

// EndpointInputParams struct holds as input all the url params of the request
type EndpointInputParams struct {
	startTime int // UTC time in W3C format
	endTime   int
	report    string
	group     string
	service   string
	hostname  string
	format    string
}

// EndpointDataOutput struct holds the queried data from datastore
type EndpointDataOutput struct {
	Report           string            `bson:"report"`
	Timestamp        string            `bson:"timestamp"`
	Group            string            `bson:"endpoint_group"`
	Service          string            `bson:"service"`
	Hostname         string            `bson:"host"`
	Status           string            `bson:"status"`
	DateInt          int               `bson:"date_integer"`
	HasThresholdRule bool              `bson:"has_threshold_rule"`
	Info             map[string]string `bson:"info"`
}

// MetricInputParams struct holds as input all the url params of the request
type MetricInputParams struct {
	startTime int // UTC time in W3C format
	endTime   int
	report    string
	group     string
	service   string
	hostname  string
	metric    string
	format    string
}

// MetricDataOutput struct holds the queried data from datastore
type MetricDataOutput struct {
	Timestamp      string            `bson:"timestamp"`
	Group          string            `bson:"endpoint_group"`
	Service        string            `bson:"service"`
	Hostname       string            `bson:"host"`
	Metric         string            `bson:"metric"`
	Status         string            `bson:"status"`
	DateInt        int               `bson:"date_int"`
	PrevTimestamp  string            `bson:"previous_timestamp"`
	PrevStatus     string            `bson:"previous_state"`
	ActualData     string            `bson:"actual_data"`
	RuleApplied    string            `bson:"threshold_rule_applied"`
	OriginalStatus string            `bson:"original_status"`
	Info           map[string]string `bson:"info"`
}

// json response related structs

type rootGroupOUT struct {
	Groups []*groupOUT `json:"groups"`
}

type groupOUT struct {
	Name      string       `json:"name"`
	GroupType string       `json:"type"`
	Statuses  []*statusOUT `json:"statuses"`
}

type rootServiceOUT struct {
	Groups []*serviceGroupOUT `json:"groups"`
}

type serviceGroupOUT struct {
	Name      string        `json:"name"`
	GroupType string        `json:"type"`
	Services  []*serviceOUT `json:"service-types"`
}

type serviceOUT struct {
	Name      string       `json:"name"`
	GroupType string       `json:"type"`
	Statuses  []*statusOUT `json:"statuses"`
}

type rootEndpointOUT struct {
	Groups []*endpointGroupOUT `json:"groups"`
}

type endpointGroupOUT struct {
	Name      string                `json:"name"`
	GroupType string                `json:"type"`
	Services  []*endpointServiceOUT `json:"service-types"`
}

type endpointServiceOUT struct {
	Name      string         `json:"name"`
	GroupType string         `json:"type"`
	Endpoints []*endpointOUT `json:"endpoints"`
}

type endpointOUT struct {
	Name       string            `json:"name"`
	Service    string            `json:"service,omitempty"`
	SuperGroup string            `json:"supergroup,omitempty"`
	Info       map[string]string `json:"info,omitempty"`
	Statuses   []*statusOUT      `json:"statuses"`
}

type rootMetricOUT struct {
	Groups []*metricGroupOUT `json:"groups"`
}

type metricGroupOUT struct {
	Name      string              `json:"name"`
	GroupType string              `json:"type"`
	Services  []*metricServiceOUT `json:"service-types"`
}

type metricServiceOUT struct {
	Name      string               `json:"name"`
	GroupType string               `json:"type"`
	Endpoints []*metricEndpointOUT `json:"endpoints"`
}

type metricEndpointOUT struct {
	Name    string            `json:"name"`
	Info    map[string]string `json:"info,omitempty"`
	Metrics []*metricOUT      `json:"metrics"`
}

type metricOUT struct {
	Name     string             `json:"name"`
	Statuses []*metricStatusOUT `json:"statuses"`
}

type metricStatusOUT struct {
	Timestamp      string `json:"timestamp"`
	Value          string `json:"value"`
	ActualData     string `json:"actual_data,omitempty"`
	RuleApplied    string `json:"threshold_rule_applied,omitempty"`
	OriginalStatus string `json:"original_status,omitempty"`
}

type statusOUT struct {
	Timestamp               string `json:"timestamp"`
	Value                   string `json:"value"`
	AffectedByThresholdRule bool   `json:"affected_by_threshold_rule,omitempty"`
}

// Message struct to hold the json/xml response
type messageOUT struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}
