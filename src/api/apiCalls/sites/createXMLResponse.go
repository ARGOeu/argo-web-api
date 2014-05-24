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

package sites

import (
	"encoding/xml"
	"fmt"
	"time"
)

// a series of auxiliary structs that will
// help us form the xml response

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
	XMLName xml.Name `xml:"Profile"`
	Name    string   `xml:"name,attr"`
	Site    []*Site
}

type Root struct {
	XMLName xml.Name `xml:"root"`
	Profile []*Profile
}

func CreateXMLResponse(results []ApiSiteAvailabilityInProfileOutput) ([]byte, error) {

	docRoot := &Root{}

	prevProfile := ""
	prevSite := ""
	site := &Site{}
	profile := &Profile{}

	// we iterate through the results struct array
	// keeping only the value of each row

	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], fmt.Sprint(row.Date))
		//if new profile value does not match the previous profile value
		//we create a new profile in the xml
		if prevProfile != row.Profile {
			prevProfile = row.Profile
			profile = &Profile{
				Name: row.Profile}
			docRoot.Profile = append(docRoot.Profile, profile)
			prevSite = ""
		}
		//if new site does not match the previous service value
		//we create a new site entry in the xml
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
		//we append the new availability values
		site.Availability = append(site.Availability,
			&Availability{
				Timestamp:    timestamp.Format(customForm[1]),
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability)})
	}
	//we create the xml response and record the output and any possible errors
	//in the appropriate variables
	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	//we return the output
	return output, err
}
