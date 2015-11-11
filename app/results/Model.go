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

package results

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/respond"
	"gopkg.in/mgo.v2"
)

type list []interface{}

var customForm []string

func init() {
	customForm = []string{"20060102", "2006-01-02T15:04:05Z"} //{"Format that is returned by the database" , "Format that will be used in the generated report"}
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

type serviceFlavorResultQuery struct {
	basicQuery
	EndpointGroup string `bson:"supergroup"`
}

type endpointGroupResultQuery struct {
	basicQuery
	Group string `bson:"supergroup"`
}

// ReportInterface for mongodb object exchanging
// type ReportInterface struct {
// 	Name              string `bson:"name"`
// 	Tenant            string `bson:"tenant"`
// 	EndpointGroupType string `bson:"endpoint_group"`
// 	SuperGroupType    string `bson:"group_of_groups"`
// }

// ServiceFlavorInterface for mongodb object exchanging
type ServiceFlavorInterface struct {
	Name         string  `bson:"name"`
	Report       string  `bson:"report"`
	Date         string  `bson:"date"`
	Type         string  `bson:"type"`
	Up           float64 `bson:"up"`
	Down         float64 `bson:"down"`
	Unknown      float64 `bson:"unknown"`
	Availability float64 `bson:"availability"`
	Reliability  float64 `bson:"reliability"`
	SuperGroup   string  `bson:"supergroup"`
}

// EndpointGroupInterface for mongodb object exchanging
type EndpointGroupInterface struct {
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

// SuperGroupInterface for mongodb object exchanging
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

//Availability struct for formating xml/json
type Availability struct {
	XMLName      xml.Name `xml:"results" json:"-"`
	Timestamp    string   `xml:"timestamp,attr" json:"timestamp"`
	Availability string   `xml:"availability,attr" json:"availability"`
	Reliability  string   `xml:"reliability,attr" json:"reliability"`
	Unknown      string   `xml:"unknown,attr,omitempty" json:"unknown,omitempty"`
	Uptime       string   `xml:"uptime,attr,omitempty" json:"uptime,omitempty"`
	Downtime     string   `xml:"downtime,attr,omitempty" json:"downtime,omitempty"`
}

// ServiceFlavor struct for formating xml/json
type ServiceFlavor struct {
	XMLName      xml.Name      `xml:"group" json:"-"`
	Name         string        `xml:"name,attr" json:"name"`
	Type         string        `xml:"type,attr" json:"type"`
	Availability []interface{} `json:"results"`
}

// Group struct for formating xml/json
type Group struct {
	XMLName      xml.Name      `xml:"group" json:"-"`
	Name         string        `xml:"name,attr" json:"name"`
	Type         string        `xml:"type,attr" json:"type"`
	Availability []interface{} `json:"results"`
}

// ServiceFlavorGroup struct for formating xml/json
type ServiceFlavorGroup struct {
	XMLName       xml.Name      `xml:"group" json:"-"`
	Name          string        `xml:"name,attr" json:"name"`
	Type          string        `xml:"type,attr" json:"type"`
	ServiceFlavor []interface{} `json:"serviceflavors"`
}

// SuperGroup struct for formating xml/json
type SuperGroup struct {
	XMLName   xml.Name      `xml:"group" json:"-"`
	Name      string        `xml:"name,attr" json:"name"`
	Type      string        `xml:"type,attr" json:"type"`
	Endpoints []interface{} `json:"endpoints,omitempty"`
	Results   []interface{} `json:"results,omitempty"`
}

type root struct {
	XMLName xml.Name      `xml:"root" json:"-"`
	Result  []interface{} `json:"root"`
}

// errorMessage struct to hold the json/xml error response
type errorMessage struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Message string   `xml:"message" json:"message"`
}

// ErrorResponse shortcut to respond.ErrorResponse
type ErrorResponse respond.ErrorResponse

func (query *basicQuery) Validate(db *mgo.Database) []ErrorResponse {
	errs := []ErrorResponse{}
	query.Granularity = strings.ToLower(query.Granularity)
	if query.Granularity == "" {
		query.Granularity = "daily"
	} else if query.Granularity != "daily" && query.Granularity != "monthly" {
		errs = append(errs, ErrorResponse{
			Message: "Wrong Granularity",
			Code:    "400",
			Details: fmt.Sprintf("%s is not accepted as granularity parameter, please provide either daily or monthly", query.Granularity),
		})
	}

	if query.StartTime == "" && query.EndTime == "" {
		errs = append(errs, ErrorResponse{
			Message: "No time span set",
			Code:    "400",
			Details: "Please use start_time and/or end_time url parameters to set the prefered time span",
		})
	} else {
		if query.StartTime != "" {
			ts, tserr := time.Parse(zuluForm, query.StartTime)
			if tserr != nil {
				errs = append(errs, ErrorResponse{
					Message: "start_time parsing error",
					Code:    "400",
					Details: fmt.Sprintf("Error parsing date string %s please use zulu format like %s", query.StartTime, zuluForm),
				})
			}
			query.StartTimeInt, _ = strconv.Atoi(ts.Format(ymdForm))
		}
		if query.EndTime != "" {
			te, teerr := time.Parse(zuluForm, query.EndTime)
			if teerr != nil {
				errs = append(errs, ErrorResponse{
					Message: "end_time parsing error",
					Code:    "400",
					Details: fmt.Sprintf("Error parsing date string %s please use zulu format like %s", query.EndTime, zuluForm),
				})
			}
			query.EndTimeInt, _ = strconv.Atoi(te.Format(ymdForm))
		}
	}

	if query.Report.DetermineGroupType(query.Vars["group_type"]) == "endpoint" {
		query.Vars["lgroup_type"] = query.Vars["group_type"]
		query.Vars["lgroup_name"] = query.Vars["group_name"]
		query.Vars["group_type"] = ""
		query.Vars["group_name"] = ""

	} else if query.Vars["group_type"] != "" && query.Report.DetermineGroupType(query.Vars["group_type"]) != "group" {
		errs = append(errs, ErrorResponse{
			Message: "Group type not in report",
			Code:    "400",
			Details: fmt.Sprintf("Group type %s not present in report %s.", query.Vars["group_type"], query.Vars["report_name"]),
		})
	}

	_, exists := query.Vars["lgroup_type"]

	if exists && query.Report.DetermineGroupType(query.Vars["lgroup_type"]) != "endpoint" {
		errs = append(errs, ErrorResponse{
			Message: "Endpoint Group type not in report",
			Code:    "400",
			Details: fmt.Sprintf("Endpoint Group type %s not present in report %s. Try using %s instead",
				query.Vars["lgroup_type"], query.Vars["report_name"], query.Report.GetEndpointGroupType()),
		})
	}

	return errs
}
