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
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/respond"
	"gopkg.in/mgo.v2"
)

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

	if query.StartTime == "" || query.EndTime == "" {
		errs = append(errs, ErrorResponse{
			Message: "No time span set",
			Code:    "400",
			Details: "Please use start_time and end_time url parameters to set the prefered time span",
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

	return errs
}
