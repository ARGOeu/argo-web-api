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

package reports

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/ARGOeu/argo-web-api/respond"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoInterface is used as an interface to Marshal and Unmarshal from different formats
type MongoInterface struct {
	ID           string        `bson:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	Tenant       string        `json:"tenant" xml:"tenant"`
	Weight       string        `bson:"weight,omitempty"  json:"weight,omitempty"  xml:"weight,omitempty"`
	Disabled     bool          `bson:"disabled" json:"disabled" xml:"disabled"`
	Info         Info          `bson:"info" json:"info" xml:"info"`
	Computations *Computations `bson:"computations" json:"computations" xml:"-"`
	Thresholds   *Thresholds   `bson:"thresholds,omitempty" json:"thresholds,omitempty" xml:"thresholds"`
	Topology     Topology      `bson:"topology_schema" json:"topology_schema" xml:"topology_schema"`
	Profiles     []Profile     `bson:"profiles" json:"profiles" xml:"profiles"`
	Tags         []Tag         `bson:"filter_tags" json:"filter_tags" xml:"filter_tags"`
}

// Info contains info about a report and is used inside the main MongoInterface struct
type Info struct {
	Name        string `bson:"name,omitempty" json:"name" xml:"name"`
	Description string `bson:"description,omitempty" json:"description" xml:"description"`
	Created     string `bson:"created,omitempty" json:"created,omitempty" xml:"created,omitempty"`
	Updated     string `bson:"updated,omitempty" json:"updated,omitempty" xml:"updated,omitempty"`
}

// Computations struct holds information about what needs to be computed in this report
// e.g. a/r, status and specific trends
type Computations struct {
	AR     bool     `bson:"ar" json:"ar"`
	Status bool     `bson:"status" json:"status"`
	Trends []string `bson:"trends" json:"trends"`
}

// Thresholds contains information about the percentage thresholds used to color report scores
type Thresholds struct {
	Availabilty float32 `bson:"availability" json:"availability" xml:"availability"`
	Reliability float32 `bson:"reliability" json:"reliability" xml:"reliability"`
	Uptime      float32 `bson:"uptime" json:"uptime" xml:"uptime"`
	Unknown     float32 `bson:"unknown" json:"unknown" xml:"unknown"`
	Downtime    float32 `bson:"downtime" json:"downtime" xml:"downtime"`
}

// Topology contains the topology used in this report and is used inside the main MongoInterface struct
type Topology struct {
	// Nesting int            `bson:"nesting" json:"nesting" xml:"nesting"`
	Group *TopologyLevel `bson:"group,omitempty" json:"group,omitempty" xml:"group,omitempty"`
}

// TopologyLevel is used to create the multiple nesting levels for the Topology struct
type TopologyLevel struct {
	Type  string         `bson:"type" json:"type" xml:"type"`
	Group *TopologyLevel `bson:"group,omitempty" json:"group,omitempty" xml:"group,omitempty"`
}

// Profile holds info about the profiles included in a report definition
type Profile struct {
	XMLName xml.Name `bson:"-"          json:"-"     xml:"profile"`
	ID      string   `bson:"id"         json:"id"    xml:"id,attr"`
	Name    string   `bson:"name"       json:"name"  xml:"name,attr"`
	Type    string   `bson:"type"       json:"type"  xml:"type,attr"`
}

// Tag holds info about the tags used in filtering in a report definition
type Tag struct {
	XMLName xml.Name `bson:",omitempty" json:"-"     xml:"tag"`
	Name    string   `bson:"name"       json:"name"  xml:"name,attr"`
	Value   string   `bson:"value"      json:"value" xml:"value,attr"`
	Context string   `bson:"context"    json:"context" xml:"context,attr"`
}

// Message struct for xml message response
type Message struct {
	XMLName xml.Name `xml:"root"`
	Message string
}

// RootXML struct to represent the root of the xml document
type RootXML struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Reports interface{}
}

// GetEndpointGroupType retrieves the deepest type nested inside the group hierarchy
func (report MongoInterface) GetEndpointGroupType() string {
	currentObject := report.Topology.Group
	for currentObject.Group != nil {
		currentObject = currentObject.Group
	}
	return currentObject.Type
}

// GetGroupType retrieves the first type nested inside the group hierarchy
func (report MongoInterface) GetGroupType() string {
	return report.Topology.Group.Type
}

func defaultThresholds() Thresholds {
	return Thresholds{
		Availabilty: 80.00,
		Reliability: 85.00,
		Uptime:      80.00,
		Unknown:     10.00,
		Downtime:    10.00,
	}
}

// DetermineGroupType looks into a report struct topology group pointers and determines
// whether a given groupType is a lesser_group or group or does not exist in the report.
func (report MongoInterface) DetermineGroupType(groupType string) string {
	nestinglevel := 1
	currentObject := report.Topology.Group
	found := false
	for currentObject.Group != nil {
		nestinglevel++
		if currentObject.Type == groupType {
			found = true
		}
		currentObject = currentObject.Group
	}
	if currentObject.Type == groupType {
		return "endpoint"
	} else if found {
		return "group"
	}
	return ""
}

var validators = map[string]string{
	"metric":      "metric_profiles",
	"aggregation": "aggregation_profiles",
	"operations":  "operations_profiles",
	"thresholds":  "thresholds_profiles",
	"weights":     "weights",
}

// ValidateProfiles ensures that the profiles in a report actually exist in the database and
// corrects possible name inconsistencies
func (report *MongoInterface) ValidateProfiles(db *mongo.Database) []respond.ErrorResponse {
	errs := []respond.ErrorResponse{}
	for idx, element := range report.Profiles {
		var result Profile
		colName := validators[element.Type]
		if colName != "" {
			err := db.Collection(colName).FindOne(context.TODO(), bson.M{"id": element.ID}).Decode(&result)
			if err != nil {
				errs = append(errs,
					respond.ErrorResponse{
						Message: "Profile id not found",
						Code:    "422",
						Details: fmt.Sprintf("No profile in %s was found with id %s", colName, element.ID),
					})
				continue
			}
			report.Profiles[idx].Name = result.Name
		} else {
			errs = append(errs,
				respond.ErrorResponse{
					Message: "Profile type invalid",
					Code:    "422",
					Details: fmt.Sprintf("Profile type %s is invalid", element.Type),
				})
			continue
		}
	}
	return errs
}

func (report *MongoInterface) ValidateTrends() []respond.ErrorResponse {
	errs := []respond.ErrorResponse{}
	for _, trendName := range report.Computations.Trends {
		if trendName != "flapping" && trendName != "status" && trendName != "tags" {
			errs = append(errs,
				respond.ErrorResponse{
					Message: "Invalid Trend Name",
					Code:    "422",
					Details: fmt.Sprintf("Trends with the name:%s doesn't exist", trendName),
				})
		}

	}
	return errs
}

// GetMetricProfile is a function that takes a report struc element
// and returns the name of the metric profile (if exists)
func GetMetricProfile(input MongoInterface) (string, error) {
	for _, element := range input.Profiles {
		if element.Type == "metric" {
			return element.Name, nil
		}
	}
	return "", errors.New("unable to find metric profile with specified name")
}

// searchName is used to create a simple query object based on name
func searchName(name string) bson.M {
	query := bson.M{
		"info.name": name,
	}

	return query
}

// genDefaultComp is used to generate a default computations json object for the report
func genDefaultComp() *Computations {
	return &Computations{AR: true, Status: true, Trends: []string{"flapping", "status", "tags"}}
}
