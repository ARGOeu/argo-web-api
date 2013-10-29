package ngis

import (
	"encoding/xml"
	"fmt"
	"time"
)

type MongoNgi struct {
	Date         string  "dt"
	Namespace    string  "ns"
	Profile      string  "p"
	Ngi          string  "n"
	Availability float64 "a"
	Reliability  float64 "r"
}

func CreateXMLResponse(results []MongoNgi, customForm []string) ([]byte, error) {

	type Availability struct {
		XMLName      xml.Name `xml:"Availability"`
		Timestamp    string   `xml:"timestamp,attr"`
		Availability string   `xml:"availability,attr"`
		Reliability  string   `xml:"reliability,attr"`
	}

	type Ngi struct {
		Ngi          string `xml:"NGI,attr"`
		Availability []*Availability
	}

	type Profile struct {
		XMLName   xml.Name `xml:"Profile"`
		Name      string   `xml:"name,attr"`
		Namespace string   `xml:"namespace,attr"`
		Ngi       []*Ngi
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Profile []*Profile
	}

	v := &Root{}

	prevProfile := ""
	prevNgi := ""
	ngi := &Ngi{}
	profile := &Profile{}
	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], row.Date)

		if prevProfile != row.Profile {
			prevProfile = row.Profile
			profile = &Profile{
				Name:      row.Profile,
				Namespace: row.Namespace}
			v.Profile = append(v.Profile, profile)
			prevNgi = ""
		}

		if prevNgi != row.Ngi {
			prevNgi = row.Ngi
			ngi = &Ngi{
				Ngi: row.Ngi,
			}
			profile.Ngi = append(profile.Ngi, ngi)
		}
		ngi.Availability = append(ngi.Availability,
			&Availability{
				Timestamp:    timestamp.Format(customForm[1]),
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability)})
	}

	output, err := xml.MarshalIndent(v, " ", "  ")

	return output, err
}
