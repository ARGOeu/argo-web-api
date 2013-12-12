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

package ngis

//we import the appropriate libraries

import (
	"encoding/xml"
	"fmt"
	"time"
)

//struct contains all information required to form an appropriate xml respnose

type MongoNgi struct {
	Date         string  "dt"
	Namespace    string  "ns"
	Profile      string  "p"
	Ngi          string  "n"
	Availability float64 "a"
	Reliability  float64 "r"
}

// a series of auxiliary structs that will
// help us form the xml response
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

func CreateXMLResponse(results []MongoNgi, customForm []string) ([]byte, error) {

	v := &Root{}

	prevProfile := ""
	prevNgi := ""
	ngi := &Ngi{}
	profile := &Profile{}
	// we iterate through the results struct array
	// keeping only the value of each row
	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], row.Date)
		//if new profile value does not match the previous profile value
		//we create a new profile in the xml
		if prevProfile != row.Profile {
			prevProfile = row.Profile
			profile = &Profile{
				Name:      row.Profile,
				Namespace: row.Namespace}
			v.Profile = append(v.Profile, profile)
			prevNgi = ""
		}
		//if new ngi does not match the previous ngi value
		//we create a new ngi entry in the xml
		if prevNgi != row.Ngi {
			prevNgi = row.Ngi
			ngi = &Ngi{
				Ngi: row.Ngi,
			}
			profile.Ngi = append(profile.Ngi, ngi)
		}
		//we append the new availability values
		ngi.Availability = append(ngi.Availability,
			&Availability{
				Timestamp:    timestamp.Format(customForm[1]),
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability)})
	}
	//we create the xml response and record the output and any possible errors
	//in the appropriate variables
	output, err := xml.MarshalIndent(v, " ", "  ")
	//we return the output
	return output, err
}
