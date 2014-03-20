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

package availabilityProfiles

import "encoding/xml"

type Content struct {
	XMLName       xml.Name `xml:Content`
	ServiceFlavor string   `xml:"ServiceFlavor,attr"`
	Grouping      string   `xml:"Grouping,attr"`
}

type Profile struct {
	XMLName xml.Name `xml:"Profile"`
	Name    string   `xml:"name,attr"`
	Content []*Content
}

type ReadRoot struct {
	XMLName xml.Name `xml:"root"`
	Profile []*Profile
}

type Message struct {
	XMLName xml.Name `xml:"root"`
	Message string	 
}


func readXML(results []ApiAPOutput) ([]byte, error) {
	prevProfile := ""
	docRoot := &ReadRoot{}
	profile := &Profile{}
	for _, row := range results {
		if prevProfile != row.Name {
			prevProfile = row.Name
			profile = &Profile{
				Name: row.Name,
			}
			docRoot.Profile = append(docRoot.Profile, profile)
		}
		profile.Content = append(profile.Content,
			&Content{
				ServiceFlavor: row.ServiceFlavor,
				Grouping:      row.Grouping})
	}
	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err
}

func messageXML(answer string) ([]byte, error){
	docRoot := &Message{}
	docRoot.Message = answer
	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err
}
