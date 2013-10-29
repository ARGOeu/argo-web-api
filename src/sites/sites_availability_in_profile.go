package sites

import (
	"encoding/xml"
	"fmt"
	"time"
)

type MongoSite struct {
	SiteScope     string  "ss"
	Scope         string  "sc"
	Date          string  "dt"
	Namespace     string  "ns"
	Profile       string  "p"
	Production    string  "pr"
	Monitored     string  "m"
	Ngi           string  "n"
	Site          string  "s"
	Infastructure string  "i"
	CertStatus    string  "cs"
	Availability  float64 "a"
	Reliability   float64 "r"
}
type Availability struct {
	XMLName      xml.Name `xml:"Availability"`
	Timestamp    string   `xml:"timestamp,attr"`
	Availability string   `xml:"availability,attr"`
	Reliability  string   `xml:"reliability,attr"`
}

type Site struct {
	Site          string `xml:"site,attr"`
	Ngi           string `xml:"NGI,attr"`
	Infastructure string `xml:"infastructure,attr"`
	Scope         string `xml:"scope,attr"`
	SiteScope     string `xml:"site_scope,attr"`
	Production    string `xml:"production,attr"`
	Monitored     string `xml:"monitored,attr"`
	CertStatus    string `xml:"certification_status,attr"`
	Availability  []*Availability
}

type Profile struct {
	XMLName   xml.Name `xml:"Profile"`
	Name      string   `xml:"name,attr"`
	Namespace string   `xml:"namespace,attr"`
	Site      []*Site
}

type Root struct {
	XMLName xml.Name `xml:"root"`
	Profile []*Profile
}

func CreateXMLResponse(results []MongoSite, customForm []string) ([]byte, error) {

	v := &Root{}

	prevProfile := ""
	prevSite := ""
	site := &Site{}
	profile := &Profile{}
	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], fmt.Sprint(row.Date))

		if prevProfile != row.Profile {
			prevProfile = row.Profile
			profile = &Profile{
				Name:      row.Profile,
				Namespace: row.Namespace}
			v.Profile = append(v.Profile, profile)
			prevSite = ""
		}

		if prevSite != row.Site {
			prevSite = row.Site
			site = &Site{
				Site:          row.Site,
				Ngi:           row.Ngi,
				Infastructure: row.Infastructure,
				Scope:         row.Scope,
				SiteScope:     row.SiteScope,
				Production:    row.Production,
				Monitored:     row.Monitored,
				CertStatus:    row.CertStatus}
			profile.Site = append(profile.Site, site)
		}
		site.Availability = append(site.Availability,
			&Availability{
				Timestamp:    timestamp.Format(customForm[1]),
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability)})
		fmt.Println(v)
	}

	output, err := xml.MarshalIndent(v, " ", "  ")

	return output, err
}
