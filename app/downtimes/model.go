package downtimes

//Downtimes holds a list of downtimes for endpoints for a specific day
type Downtimes struct {
	ID        string     `bson:"id" json:"id"`
	Name      string     `bson:"name" json:"name"`
	Endpoints []Downtime `bson:"endpoints" json:"endpoints"`
}

//Downtime holds downtime information for a specific host
type Downtime struct {
	HostName  string `bson:"hostname" json:"hostname"`
	Service   string `bson:"service" json:"service"`
	StartTime string `bson:"start_time" json:"start_time"`
	EndTime   string `bson:"end_time" json:"end_time"`
}

// SelfReference to hold links and id
type SelfReference struct {
	ID    string `json:"id" bson:"id,omitempty"`
	Links Links  `json:"links"`
}

// Links struct to hold links
type Links struct {
	Self string `json:"self"`
}
