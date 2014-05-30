/*
 * Copyright (c) 2014 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
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

import (
	"encoding/xml"
	"labix.org/v2/mgo/bson"
)

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

//Struct for inserting data into DB
type AvailabilityProfileInput struct {
	Name      string     `json:"name"`
	Namespace string     `json:"namespace"`
	Groups    [][]string `json:"groups"`
	Poems     []string   `json:"poems"`
}

//Struct for searching based on name and namespace combination
type AvailabilityProfileSearch struct {
	Name      []string
	Namespace []string
}

//Struct for record retrieval
type AvailabilityProfileOutput struct {
	ID        string
	Name      string     `bson:"name"`
	Namespace string     `bson:"namespace"`
	Groups    [][]string `bson:"groups"`
	Poems     []string   `bson:"poems"`
}

func prepareFilter(input AvailabilityProfileSearch) bson.M {

	filter := bson.M{
		"name":      bson.M{"$in": input.Name},
		"namespace": bson.M{"$in": input.Namespace},
	}

	return filter
}

func createOne(input AvailabilityProfileInput) bson.M {
	query := bson.M{
		"name":      input.Name,
		"namespace": input.Namespace,
		"groups":    input.Groups,
		"poems":     input.Poems,
	}
	return query
}

func readOne(input AvailabilityProfileSearch) bson.M {
	filter := prepareFilter(input)
	return filter
}
