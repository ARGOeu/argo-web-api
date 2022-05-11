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

package status

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

// InputParams struct holds as input all the url params of the request
type InputParams struct {
	startTime int // UTC time in W3C format
	endTime   int
	report    string
	groupType string
	group     string
	format    string
}

// GroupData struct holds the grouped queried data from datastore
type GroupData struct {
	Report           string `bson:"report"`
	Timestamp        string `bson:"timestamp"`
	Group            string `bson:"endpoint_group"`
	Status           string `bson:"status"`
	DateInteger      string `bson:"date_integer"`
	HasThresholdRule bool   `bson:"has_threshold_rule"`
}

// json response related structs

type rootOUT struct {
	Groups []*groupOUT `json:"groups"`
}

type groupOUT struct {
	Name      string         `json:"name"`
	GroupType string         `json:"type"`
	Statuses  []*statusOUT   `json:"statuses"`
	Endpoints []*endpointOUT `json:"endpoints"`
}

type statusOUT struct {
	Timestamp               string `json:"timestamp"`
	Value                   string `json:"value"`
	AffectedByThresholdRule bool   `json:"affected_by_threshold_rule,omitempty"`
}

// Message struct to hold the json response
type messageOUT struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// EndpointData struct holds the queried data from datastore
type EndpointData struct {
	Report           string            `bson:"report"`
	Timestamp        string            `bson:"timestamp"`
	EndpointGroup    string            `bson:"endpoint_group"`
	Service          string            `bson:"service"`
	Hostname         string            `bson:"host"`
	Status           string            `bson:"status"`
	DateInt          string            `bson:"date_integer"`
	HasThresholdRule bool              `bson:"has_threshold_rule"`
	Info             map[string]string `bson:"info"`
}

type endpointOUT struct {
	Name       string            `json:"hostname"`
	Service    string            `json:"service,omitempty"`
	SuperGroup string            `json:"-"`
	Info       map[string]string `json:"info,omitempty"`
	Statuses   []*statusOUT      `json:"statuses"`
}

// errorMessage struct to hold the json/xml error response
type errorMessage struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
