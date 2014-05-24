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

type Group struct {
	XMLName       xml.Name
	ServiceFlavor string `xml:"service_flavor,attr"`
}

type Or struct {
	XMLName xml.Name `xml:"OR"`
	Group   []*Group
}

type And struct {
	XMLName xml.Name `xml:"AND"`
	Or      []*Or
}

type Profile struct {
	XMLName   xml.Name `xml:"profile"`
	ID        string   `xml:"id,attr"`
	Name      string   `xml:"name,attr"`
	Namespace string   `xml:"namespace,attr"`
	Poem      string   `xml:"poems,attr"`
	And       *And
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
	docRoot := &ReadRoot{}
	for _, row := range results {
		profile := &Profile{
			ID:        row.ID,
			Name:      row.Name,
			Namespace: row.Namespace,
			Poem:      row.Poems[0],
		}
		and := &And{}
		docRoot.Profile = append(docRoot.Profile, profile)
		for _, group := range row.Groups {
			or := &Or{}
			for _, sf := range group {
				group := &Group{
					ServiceFlavor: sf,
				}
				or.Group = append(or.Group, group)
			}
			and.Or = append(and.Or, or)
		}
		profile.And = and
	}
	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err
}

func messageXML(answer string) ([]byte, error) {
	docRoot := &Message{}
	docRoot.Message = answer
	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err
}
