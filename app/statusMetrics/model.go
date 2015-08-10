package statusMetrics

import "encoding/xml"

// InputParams struct holds as input all the url params of the request
type InputParams struct {
	startTime string // UTC time in W3C format
	endTime   string
	report    string
	groupType string
	group     string
	service   string
	hostname  string
	metric    string
}

// DataOutput struct holds the queried data from datastore
type DataOutput struct {
	Report        string `bson:"report"`
	Timestamp     string `bson:"timestamp"`
	EndpointGroup string `bson:"endpoint_group"`
	Service       string `bson:"service"`
	Hostname      string `bson:"hostname"`
	Metric        string `bson:"metric"`
	Status        string `bson:"status"`
	DateInt       string `bson:"date_int"`
	PrevTimestamp string `bson:"prev_timestamp"`
	PrevStatus    string `bson:"prev_status"`
}

// xml response related structs

type rootXML struct {
	XMLName        xml.Name `xml:"root"`
	EndpointGroups []*endpointGroupXML
}

type endpointGroupXML struct {
	XMLName   xml.Name `xml:"Group"`
	Name      string   `xml:"name,attr"`
	GroupType string   `xml:"type,attr"`
	Services  []*serviceXML
}

type serviceXML struct {
	XMLName   xml.Name `xml:"Group"`
	Name      string   `xml:"name,attr"`
	GroupType string   `xml:"type,attr"`
	Endpoints []*endpointXML
}

type endpointXML struct {
	XMLName   xml.Name `xml:"Group"`
	Name      string   `xml:"name,attr"`
	GroupType string   `xml:"type,attr"`
	Metrics   []*metricXML
}

type metricXML struct {
	XMLName   xml.Name `xml:"Group"`
	Name      string   `xml:"name,attr"`
	GroupType string   `xml:"type,attr"`
	Statuses  []*statusXML
}

type statusXML struct {
	XMLName   xml.Name `xml:"status"`
	Timestamp string   `xml:"timestamp,attr"`
	Value     string   `xml:"value,attr"`
}

// Message struct to hold the xml response
type message struct {
	XMLName xml.Name `xml:"root"`
	Message string
}
