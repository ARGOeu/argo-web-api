/*
 * Copyright (c) 2022 GRNET S.A.
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

package ar

import (
	"github.com/ARGOeu/argo-web-api/app/reports"
)

type list []interface{}

var customForm []string

func init() {
	customForm = []string{"20060102", "2006-01-02T15:04:05Z"} //{"Format that is returned by the database" , "Format that will be used in the generated report"}
}

type root struct {
	Result []*SuperGroup `json:"results"`
}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

type basicQuery struct {
	Name         string                 `bson:"name"`
	Granularity  string                 `bson:"-"`
	Format       string                 `bson:"-"`
	StartTime    string                 `bson:"-"` // UTC time in W3C format
	EndTime      string                 `bson:"-"` // UTC time in W3C format
	Report       reports.MongoInterface `bson:"report"`
	StartTimeInt int                    `bson:"start_time"`
	EndTimeInt   int                    `bson:"end_time"`
	Vars         map[string]string      `bson:"-"`
}

type GroupResultQuery struct {
	basicQuery
	Group string `bson:"supergroup"`
}

// GroupInterface used to hold mongodb group information such as SITES, SERVICEGROUPS etc
type GroupInterface struct {
	Name         string  `bson:"name"`
	Report       string  `bson:"report"`
	Date         string  `bson:"date"`
	Type         string  `bson:"type"`
	Up           float64 `bson:"up"`
	Down         float64 `bson:"down"`
	Unknown      float64 `bson:"unknown"`
	Availability float64 `bson:"availability"`
	Reliability  float64 `bson:"reliability"`
	Weights      string  `bson:"weights"`
	SuperGroup   string  `bson:"supergroup"`
}

// GroupInterface used to hold mongodb supergroup information such as NGIs, PROJETS etc
type SuperGroupInterface struct {
	Report       string  `bson:"report"`
	Date         string  `bson:"date"`
	Type         string  `bson:"type"`
	Up           float64 `bson:"uptime"`
	Down         float64 `bson:"downtime"`
	Unknown      float64 `bson:"unknown"`
	Availability float64 `bson:"availability"`
	Reliability  float64 `bson:"reliability"`
	Weights      string  `bson:"weight"`
	SuperGroup   string  `bson:"supergroup"`
}

// EndpointInterface used to hold mongodb group information about endpoints
type EndpointInterface struct {
	Name         string            `bson:"name"`
	Report       string            `bson:"report"`
	Date         string            `bson:"date"`
	Type         string            `bson:"type"`
	Up           float64           `bson:"up"`
	Down         float64           `bson:"down"`
	Unknown      float64           `bson:"unknown"`
	Availability float64           `bson:"availability"`
	Reliability  float64           `bson:"reliability"`
	Service      string            `bson:"service"`
	SuperGroup   string            `bson:"supergroup"`
	Info         map[string]string `bson:"info"`
}

//Availability struct for formating json
type Availability struct {
	Date         string `json:"date,omitempty"`
	Availability string `json:"availability"`
	Reliability  string `json:"reliability"`
	Unknown      string `json:"unknown,omitempty"`
	Uptime       string `json:"uptime,omitempty"`
	Downtime     string `json:"downtime,omitempty"`
}

// Group struct for formating json
type Group struct {
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Availability []Availability `json:"results"`
}

// SuperGroup struct for formating json
type SuperGroup struct {
	Name    string         `json:"name"`
	Type    string         `json:"type"`
	Results []Availability `json:"results,omitempty"`
	Groups  []*Group       `json:"groups,omitempty"`
}

// idOUT holds a/r information about endpoints with specific id
type idOUT struct {
	ID        string      `json:"id"`
	Endpoints []*Endpoint `json:"endpoints,omitempty"`
}

// endpoint holds a/r information about a specific endpoint
type Endpoint struct {
	Name       string            `json:"name"`
	Service    string            `json:"service"`
	Supergroup string            `json:"supergroup"`
	Info       map[string]string `json:"info"`
	Results    []Availability    `json:"results,omitempty"`
}

// errorMessage struct to hold the json/xml error response
type errorMessage struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
