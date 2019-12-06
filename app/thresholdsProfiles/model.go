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

package thresholdsProfiles

import (
	"log"
	"regexp"
)

// datastore collection name that contains threshold profile records
const aggColName = "threshold_profiles"

// ThresholdsProfile struct that holds information about threshold rules
type ThresholdsProfile struct {
	ID      string `bson:"id" json:"id"`
	DateInt int    `bson:"date_integer" json:"-"`
	Date    string `bson:"date" json:"date"`
	Name    string `bson:"name" json:"name"`
	Rules   []Rule `bson:"rules" json:"rules"`
}

// Rule represents a thresholds rule that must be applied to a metric and then optionally to a host and/or endpoint group
type Rule struct {
	EndpointGroup string `bson:"endpoint_group" json:"endpoint_group,omitempty"`
	Host          string `bson:"host" json:"host,omitempty"`
	Metric        string `bson:"metric" json:"metric"`
	Thresholds    string `bson:"thresholds" json:"thresholds"`
}

// SelfReference to hold links and id
type SelfReference struct {
	ID    string `json:"id" bson:"id,omitempty"`
	Links Links  `json:"links"`
}

// Links struct to hold links
type Links struct {
	Self string `json:"self"`
}

// ValidateRule takes a thresholds rule declaration and validates it against it standar format
func validateRule(rule string) bool {
	// Threshold rule validation regexp
	r, err := regexp.Compile(`^([a-zA-Z][\w]+=-?\d+(\.\d+)*(s|us|ms|%|B|KB|MB|TB|c)?(;(-?\d+(\.\d+)*|(~|-?\d+(\.\d+)*):(-?\d+(\.\d+)*)?)?){0,2}(;(-?\d+(\.\d+)*)?){0,2}(;|\s)*)+$`)
	if err != nil {
		log.Println(err)
		return false
	}

	return r.MatchString(rule)
}

// Validate validates all threshold rules of a thresholds profile
func (tprof *ThresholdsProfile) Validate() []string {

	var errList []string

	for _, rule := range tprof.Rules {
		if !validateRule(rule.Thresholds) {
			errList = append(errList, "Invalid threshold: "+rule.Thresholds)
		}
	}

	return errList

}
