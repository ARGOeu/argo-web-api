package services

import (
	"encoding/xml"
	"strconv"
	"strings"
	"time"
)

type Timeline struct {
	Profile       string "p"
	ServiceFlavor string "sf"
	Host          string "h"
	Timeline      string "tm"
	VO            string "vo"
	Date          int    "d"
	Namespace     string "ns"
}

func CreateXMLResponse(results []Timeline, customForm []string) ([]byte, error) {
	type Availability struct {
		XMLName      xml.Name `xml:"Availability"`
		Timestamp    string   `xml:"timestamp,attr"`
		Availability string   `xml:"availability,attr"`
		Reliability  string   `xml:"reliability,attr"`
		Maintenance  string   `xml:"maintenance,attr"`
	}

	type Service struct {
		Hostname       string `xml:"hostname,attr"`
		Service_Type   string `xml:"type,attr"`
		Service_Flavor string `xml:"flavor,attr"`
		Availability   []*Availability
	}

	type Profile struct {
		XMLName   xml.Name `xml:"Profile"`
		Name      string   `xml:"name,attr"`
		Namespace string   `xml:"namespace,attr"`
		VO        string   `xml:"defined_by_vo_name,attr"`
		Service   []*Service
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Profile []*Profile
	}

	v := &Root{}

	prevProfile := ""
	prevService := ""
	service := &Service{}
	profile := &Profile{}
	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], strconv.Itoa(row.Date))
		timeline := strings.Split(strings.Trim(row.Timeline, "[\n]"), ", ")

		if prevProfile != row.Profile {
			prevProfile = row.Profile
			profile = &Profile{
				Name:      row.Profile,
				Namespace: row.Namespace,
				VO:        row.VO}
			v.Profile = append(v.Profile, profile)
			prevService = ""
		}

		if prevService != row.Host+row.ServiceFlavor {
			prevService = row.Host + row.ServiceFlavor
			service = &Service{
				Hostname:       row.Host,
				Service_Type:   row.ServiceFlavor,
				Service_Flavor: row.ServiceFlavor}
			profile.Service = append(profile.Service, service)
		}

		for _, timeslot := range timeline {
			ar := strings.Split(timeslot, ":")
			if len(ar) != 3 {
				return []byte("<root><error>500: Internal server error (Malformed timeslot)</error></root>"), nil
			}

			service.Availability = append(service.Availability,
				&Availability{
					Timestamp:    timestamp.Format(customForm[1]),
					Availability: ar[0],
					Reliability:  ar[1],
					Maintenance:  ar[2]})
			timestamp = timestamp.Add(time.Duration(60*60) * time.Second)
		}

	}

	output, err := xml.MarshalIndent(v, " ", "  ")

	return output, err
}
