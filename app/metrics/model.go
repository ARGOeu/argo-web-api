package metrics

import "encoding/xml"

//Metric holds basic information about a metric profile and it's tags
type Metric struct {
	Name string   `bson:"name" json:"name"`
	Tags []string `bson:"tags" json:"tags"`
}

// Message struct to hold the json/xml response
type messageOUT struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Message string   `xml:"message" json:"message"`
	Code    string   `xml:"code,omitempty" json:"code,omitempty"`
}
