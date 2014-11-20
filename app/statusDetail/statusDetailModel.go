package statusDetail

import "encoding/xml"

type StatusDetailInput struct {
	start_time string // UTC time in W3C format
	end_time   string
	vo         string
	profile    string
	group_type string
	group      string
}

type StatusDetailOutput struct {
	Timestamp string `bson:"ts"`
	Roc       string `bson:"roc"`
	Site      string `bson:"si"`
	Service   string `bson:"sf"`
	Hostname  string `bson:"sh"`
	Metric    string `bson:"mn"`
	Status    string `bson:"s"`
	Time_int  int    `bson:"ti"`
}

type ReadRoot struct {
	XMLName xml.Name `xml:"root"`
	Profile *Profile
}

type Profile struct {
	XMLName xml.Name `xml:"profile"`
	Name    string   `xml:"name,attr"`
	Groups  []*Group
}

type Group struct {
	XMLName xml.Name `xml:"group"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Groups  []*Group
	Hosts   []*Host
}

type Host struct {
	XMLName xml.Name `xml:"host"`
	Name    string   `xml:"name,attr"`
	Metrics []*Metric
}

type Metric struct {
	XMLName  xml.Name `xml:"metric"`
	Name     string   `xml:"name,attr"`
	Timeline []*Status
}

type Status struct {
	XMLName   xml.Name `xml:"status"`
	Timestamp string
	Status    string
}
