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
 * or implied, of either GRNET S.A., SRCE or IN2P3 CNRS Computing
 * Centre
 *
 * The work represented by this source file is partially funded by
 * the EGI-InSPIRE project through the European Commission's 7th
 * Framework Programme (contract # INFSO-RI-261323)
 */

package respond

import (
	"fmt"
	"strconv"
	"time"
)

func validateDate(dateStr string) (int, error) {
	parsedTime, err := time.Parse(zuluForm, dateStr)
	if err != nil {
		return 0, err
	}
	ymdTime, err := strconv.Atoi(parsedTime.Format(ymdForm))
	return ymdTime, err
}

func ValidateDateRange(dateStart string, dateEnd string) (int, int, []ErrorResponse) {
	errs := []ErrorResponse{}
	var parsedStart, parsedEnd int = 0, 0

	if dateStart == "" && dateEnd == "" {
		errs = append(errs, ErrorResponse{
			Message: "No time span set",
			Code:    "400",
			Details: "Please use start_time and/or end_time url parameters to set the prefered time span",
		})
	} else {
		if dateStart != "" {
			parsedStartWG, errStart := validateDate(dateStart)
			parsedStart = parsedStartWG
			if errStart != nil {
				errs = append(errs, ErrorResponse{
					Message: "start_time parsing error",
					Code:    "400",
					Details: fmt.Sprintf("Error parsing date string %s please use zulu format like %s", dateStart, zuluForm),
				})
			}
		}
		if dateEnd != "" {
			parsedEndWG, errEnd := validateDate(dateEnd)
			parsedEnd = parsedEndWG
			if errEnd != nil {
				errs = append(errs, ErrorResponse{
					Message: "end_time parsing error",
					Code:    "400",
					Details: fmt.Sprintf("Error parsing date string %s please use zulu format like %s", dateEnd, zuluForm),
				})
			}
		}
	}

	return parsedStart, parsedEnd, errs
}
