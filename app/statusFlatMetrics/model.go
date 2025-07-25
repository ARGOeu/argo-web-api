/*
 * Copyright (c) 2015 GRNET S.A.
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

package statusFlatMetrics

import "encoding/xml"

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

// InputParams struct holds as input all the url params of the request
type InputParams struct {
	startTime int // UTC time in W3C format
	endTime   int
	report    string
	groupType string
	group     string
	service   string
	hostname  string
	metric    string
	format    string
}

// DataOutput struct holds the queried data from datastore
type DataOutput struct {
	Report         string            `bson:"report"`
	Timestamp      string            `bson:"timestamp"`
	EndpointGroup  string            `bson:"endpoint_group"`
	Service        string            `bson:"service"`
	Hostname       string            `bson:"host"`
	Status         string            `bson:"status"`
	Metric         string            `bson:"metric"`
	DateInt        int               `bson:"date_integer"`
	Info           map[string]string `bson:"info"`
	ActualData     string            `bson:"actual_data"`
	RuleApplied    string            `bson:"threshold_rule_applied"`
	OriginalStatus string            `bson:"original_status"`
	PrevTimestamp  string            `bson:"previous_timestamp"`
	PrevStatus     string            `bson:"previous_state"`
}

// json/xml response related structs

type rootPagedOUT struct {
	XMLName   xml.Name       `xml:"root" json:"-"`
	Endpoints []*endpointOUT `json:"endpoint_metrics"`
	PageToken string         `json:"nextPageToken,omitempty"`
	PageSize  int            `json:"pageSize,omitempty"`
}

type endpointOUT struct {
	XMLName    xml.Name          `xml:"endpoint_metric" json:"-"`
	Name       string            `xml:"name,attr" json:"name"`
	Service    string            `xml:"service,attr,omitempty" json:"service,omitempty"`
	SuperGroup string            `xml:"supergroup,attr,omitempty" json:"supergroup,omitempty"`
	Metric     string            `xml:"metric,attr,omitempty" json:"metric,omitempty"`
	Info       map[string]string `xml:"-" json:"info,omitempty"`
	Statuses   []*statusOUT      `json:"statuses"`
}

type statusOUT struct {
	XMLName        xml.Name `xml:"status" json:"-"`
	Timestamp      string   `xml:"timestamp,attr" json:"timestamp"`
	Value          string   `xml:"value,attr" json:"value"`
	ActualData     string   `xml:"-" json:"actual_data,omitempty"`
	RuleApplied    string   `xml:"-" json:"threshold_rule_applied,omitempty"`
	OriginalStatus string   `xml:"-" json:"original_status,omitempty"`
}
