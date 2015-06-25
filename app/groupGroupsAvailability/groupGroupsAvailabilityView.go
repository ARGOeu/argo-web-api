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

package groupGroupsAvailability

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

func createView(results []ApiSuperGroupAvailabilityInProfileOutput, format string) ([]byte, error) {

	docRoot := &Root{}

	prevJob := ""
	prevSuperGroup := ""
	superGroup := &SuperGroup{}
	job := &Job{}
	// we iterate through the results struct array
	// keeping only the value of each row
	for _, row := range results {
		timestamp, _ := time.Parse(CustomForm[0], row.Date)
		//if new job value does not match the previous job value
		//we create a new job in the xml
		if prevJob != row.Job {
			prevJob = row.Job
			job = &Job{
				Name: row.Job,
			}
			docRoot.Job = append(docRoot.Job, job)
			prevSuperGroup = ""
		}
		//if new superGroup does not match the previous superGroup value
		//we create a new superGroup entry in the xml
		if prevSuperGroup != row.SuperGroup {
			prevSuperGroup = row.SuperGroup
			superGroup = &SuperGroup{
				SuperGroup: row.SuperGroup,
			}
			job.SuperGroup = append(job.SuperGroup, superGroup)
		}
		//we append the new availability values
		superGroup.Availability = append(superGroup.Availability,
			&Availability{
				Timestamp:    timestamp.Format(CustomForm[1]),
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability)})
	}

	if strings.ToLower(format) == "json" {
		return json.MarshalIndent(docRoot, " ", "  ")
	} else {
		return xml.MarshalIndent(docRoot, " ", "  ")
	}
}
