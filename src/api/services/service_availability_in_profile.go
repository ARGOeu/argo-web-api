/*
 * Copyright (c) 2013 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the
 * License. You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an "AS
 * IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language
 * governing permissions and limitations under the License.
 *
 * The views and conclusions contained in the software and
 * documentation are those of the authors and should not be
 * interpreted as representing official policies, either expressed
 * or implied, of either GRNET S.A., SRCE or IN2P3 CNRS Computing
 * Centre
 *
 * The work represented by this source file is partially funded by
 * the EGI-InSPIRE project through the European Commission's 7th
 * Framework Programme (contract # INFSO-RI-261323)
 */

package services

//we import the appropriate libraries

import (
	"encoding/xml"
	"strconv"
	"strings"
	"time"
)

//struct contains all information required to form an appropriate xml respnose

type Timeline struct {
	Profile       string "p"
	ServiceFlavor string "sf"
	Host          string "h"
	Timeline      string "tm"
	VO            string "vo"
	Date          int    "d"
	Namespace     string "ns"
}

// a series of auxiliary structs that will
// help us form the xml response

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

func CreateXMLResponse(results []Timeline, customForm []string) ([]byte, error) {
	docRoot := &Root{}

	prevProfile := ""
	prevService := ""
	service := &Service{}
	profile := &Profile{}
	// we iterate through the results struct array
	// keeping only the value of each row
	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], strconv.Itoa(row.Date))
		timeline := strings.Split(strings.Trim(row.Timeline, "[\n]"), ", ")
		//if new profile value does not match the previous profile value
		//we create a new profile in the xml
		if prevProfile != row.Profile {
			prevProfile = row.Profile
			profile = &Profile{
				Name:      row.Profile,
				Namespace: row.Namespace,
				VO:        row.VO}
			docRoot.Profile = append(docRoot.Profile, profile)
			prevService = ""
		}
		//if new service does not match the previous service value
		//we create a new service entry in the xml
		if prevService != row.Host+row.ServiceFlavor {
			prevService = row.Host + row.ServiceFlavor
			service = &Service{
				Hostname:       row.Host,
				Service_Type:   row.ServiceFlavor,
				Service_Flavor: row.ServiceFlavor}
			profile.Service = append(profile.Service, service)
		}
		//we append the new availability values checking for errors
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
	//we create the xml response and record the output and any possible errors
	//in the appropriate variables
	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	//we return the output
	return output, err
}
