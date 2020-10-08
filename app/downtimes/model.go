package downtimes

//Downtimes holds a list of downtimes for endpoints for a specific day
type Downtimes struct {
	DateInt   int        `bson:"date_integer" json:"-"`
	Date      string     `bson:"date" json:"date"`
	Endpoints []Downtime `bson:"endpoints" json:"endpoints"`
}

//Downtime holds downtime information for a specific host
type Downtime struct {
	HostName  string `bson:"hostname" json:"hostname"`
	Service   string `bson:"service" json:"service"`
	StartTime string `bson:"start_time" json:"start_time"`
	EndTime   string `bson:"end_time" json:"end_time"`
}

// Links struct to hold links
type Links struct {
	Self string `json:"self"`
}
