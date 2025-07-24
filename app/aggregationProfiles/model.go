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

package aggregationProfiles

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AggProfile struct to retrieve and insert aggregation profiles in mongo
type AggProfile struct {
	ID            string        `bson:"id" json:"id"`
	DateInt       int           `bson:"date_integer" json:"-"`
	Date          string        `bson:"date" json:"date"`
	Name          string        `bson:"name" json:"name"`
	Namespace     string        `bson:"namespace" json:"namespace"`
	EndpointGroup string        `bson:"endpoint_group" json:"endpoint_group"`
	MetricOp      string        `bson:"metric_operation" json:"metric_operation"`
	ProfileOp     string        `bson:"profile_operation" json:"profile_operation"`
	MetricProf    MetricProfile `bson:"metric_profile" json:"metric_profile"`
	Groups        []Group       `bson:"groups" json:"groups"`
}

// MetricProfile is just a reference struct holding the name and the uuid of the profile
type MetricProfile struct {
	Name string `bson:"name" json:"name"`
	ID   string `bson:"id" json:"id"`
}

// Group struct to represent groupings
type Group struct {
	Name     string    `bson:"name" json:"name"`
	Op       string    `bson:"operation" json:"operation"`
	Services []Service `bson:"services" json:"services"`
}

// Service struct hold information about service operations
type Service struct {
	Name string `bson:"name" json:"name"`
	Op   string `bson:"operation" json:"operation"`
}

// SelfReference to hold links and uuid
type SelfReference struct {
	ID    string `json:"id" bson:"id,omitempty"`
	Links Links  `json:"links"`
}

// Links struct to hold links
type Links struct {
	Self string `json:"self"`
}

// validateID validates the metric profile id
func (mp *MetricProfile) validateID(col *mongo.Collection) error {
	result := MetricProfile{}
	err := col.FindOne(context.TODO(), bson.M{"id": mp.ID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err := errors.New("cannot validate")
			return err
		}

		return err
	}

	mp.Name = result.Name

	return nil

}
