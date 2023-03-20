/*
 * Copyright (c) 2018 GRNET S.A.
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
 * or implied, of GRNET S.A.
 *
 */

package topology

import (
	"encoding/xml"

	"github.com/ARGOeu/argo-web-api/utils/config"
)

// MongoInterface to retrieve and insert metricProfiles in mongo
// type MongoInterface struct {
// 	ID       string    `bson:"id" json:"id"`
// 	Name     string    `bson:"name" json:"name"`
// 	Services []Serv `bson:"services" json:"services"`
// }

// Topology struct to represent topology statistics
type Topology struct {
	GroupCount    int      `json:"group_count"`
	GroupType     string   `json:"group_type"`
	GroupList     []string `json:"group_list"`
	EndGroupCount int      `json:"endpoint_group_count"`
	EndGroupType  string   `json:"endpoint_group_type"`
	EndGroupList  []string `json:"endpoint_group_list"`
	ServiceCount  int      `json:"service_count"`
	ServiceList   []string `json:"service_list"`
}

// Message struct to hold the json/xml response
type messageOUT struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Message string   `xml:"message" json:"message"`
	Code    string   `xml:"code,omitempty" json:"code,omitempty"`
}

type fltrEndpoint struct {
	Group     []string
	GroupType []string
	Service   []string
	Hostname  []string
	Tags      string
}

type fltrGroup struct {
	Group     []string
	GroupType []string
	Subgroup  []string
	Tags      string
}

// Endpoint includes information on endpoint group topology
type Endpoint struct {
	Date          string            `bson:"date" json:"date"`
	DateInt       int               `bson:"date_integer" json:"-"`
	Group         string            `bson:"group" json:"group"`
	GroupType     string            `bson:"type" json:"type"`
	Service       string            `bson:"service" json:"service"`
	Hostname      string            `bson:"hostname" json:"hostname"`
	Notifications *Notifications    `bson:"notifications" json:"notifications,omitempty"`
	Tags          map[string]string `bson:"tags" json:"tags"`
}

// Group includes information on  of group group topology
type Group struct {
	Date          string            `bson:"date" json:"date"`
	DateInt       int               `bson:"date_integer" json:"-"`
	Group         string            `bson:"group" json:"group"`
	GroupType     string            `bson:"type" json:"type"`
	Subgroup      string            `bson:"subgroup" json:"subgroup"`
	Notifications *Notifications    `bson:"notifications" json:"notifications,omitempty"`
	Tags          map[string]string `bson:"tags" json:"tags"`
}

// ServiceType includes information about an available service type
type ServiceType struct {
	Date        string   `bson:"date" json:"date"`
	DateInt     int      `bson:"date_integer" json:"-"`
	Name        string   `bson:"name" json:"name"`
	Title       string   `bson:"title" json:"title"`
	Description string   `bson:"description" json:"description"`
	Tags        []string `bson:"tags" json:"tags,omitempty"`
	Tenant      string   `json:"tenant,omitempty"`
}

// Notifications holds notification information about topology items
type Notifications struct {
	Contacts []string `bson:"contacts" json:"contacts,omitempty"`
	Enabled  bool     `bson:"enabled" json:"enabled,omitempty"`
}

type TenantDB struct {
	Tenant string
	Config config.MongoConfig
}

// TagInfo groups all tags for a topology type
type TagInfo struct {
	Name   string      `json:"name"`
	Values []TagValues `json:"values"`
}

// TagValues holds each value appearing for each tag
type TagValues struct {
	Name   string   `bson:"_id" json:"name"`
	Values []string `bson:"v" json:"values"`
}
