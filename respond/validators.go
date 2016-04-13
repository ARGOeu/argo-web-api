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
	"net/http"
	"net/url"
	"time"
)

func validateDate(dateStr string) error {
	_, err := time.Parse(zuluForm, dateStr)
	if err != nil {
		return err
	}
	return nil
}

func ValidateMetricParams(queries url.Values) []ErrorResponse {

	var errs []ErrorResponse

	if queries["exec_time"] == nil {
		errs = append(errs, ErrorResponse{
			Message: "exec_time not set",
			Code:    fmt.Sprintf("%d", http.StatusBadRequest),
			Details: fmt.Sprintf("Please use exec_time url parameter in zulu format (like %s) to indicate the exact probe execution time", zuluForm),
		})
	} else {
		execDate := queries.Get("exec_time")
		errExec := validateDate(execDate)
		if errExec != nil {
			errs = append(errs, ErrorResponse{
				Message: "exec_time parsing error",
				Code:    fmt.Sprintf("%d", http.StatusBadRequest),
				Details: fmt.Sprintf("Error parsing date string %s please use zulu format like %s", execDate, zuluForm),
			})
		}
	}

	return errs

}

func ValidateResultsParams(queries url.Values) []ErrorResponse {

	var errs []ErrorResponse

	if queries["end_time"] == nil {
		if queries["start_time"] == nil {
			errs = append(errs, ErrorResponse{
				Message: "No time span set",
				Code:    fmt.Sprintf("%d", http.StatusBadRequest),
				Details: "Please use start_time and end_time url parameters to set the prefered time span",
			})
		} else {
			errs = append(errs, ErrorResponse{
				Message: "end_time not set",
				Code:    fmt.Sprintf("%d", http.StatusBadRequest),
				Details: fmt.Sprintf("Please use end_time url parameter in zulu format (like %s) to indicate the query end time", zuluForm),
			})
		}
	} else {
		endDate := queries.Get("end_time")
		errEnd := validateDate(endDate)
		if errEnd != nil {
			errs = append(errs, ErrorResponse{
				Message: "end_time parsing error",
				Code:    fmt.Sprintf("%d", http.StatusBadRequest),
				Details: fmt.Sprintf("Error parsing date string %s please use zulu format like %s", endDate, zuluForm),
			})
		}
	}
	if queries["start_time"] == nil {
		errs = append(errs, ErrorResponse{
			Message: "start_time not set",
			Code:    fmt.Sprintf("%d", http.StatusBadRequest),
			Details: fmt.Sprintf("Please use start_time url parameter in zulu format (like %s) to indicate the query start time", zuluForm),
		})
	} else {
		startDate := queries.Get("start_time")
		errStart := validateDate(startDate)
		if errStart != nil {
			errs = append(errs, ErrorResponse{
				Message: "start_time parsing error",
				Code:    fmt.Sprintf("%d", http.StatusBadRequest),
				Details: fmt.Sprintf("Error parsing date string %s please use zulu format like %s", startDate, zuluForm),
			})
		}
	}

	if queries["granularity"] != nil {
		granularity := queries["granularity"][0]
		if granularity != "daily" && granularity != "monthly" {
			errs = append(errs, ErrorResponse{
				Message: "Wrong Granularity",
				Code:    fmt.Sprintf("%d", http.StatusBadRequest),
				Details: fmt.Sprintf("%s is not accepted as granularity parameter, please provide either daily or monthly", granularity),
			})
		}
	}

	return errs

}

func ValidateStatusParams(queries url.Values) []ErrorResponse {

	var errs []ErrorResponse

	if queries["end_time"] == nil {
		if queries["start_time"] == nil {
			errs = append(errs, ErrorResponse{
				Message: "No time span set",
				Code:    fmt.Sprintf("%d", http.StatusBadRequest),
				Details: "Please use start_time and end_time url parameters to set the prefered time span",
			})
		} else {
			errs = append(errs, ErrorResponse{
				Message: "end_time not set",
				Code:    fmt.Sprintf("%d", http.StatusBadRequest),
				Details: fmt.Sprintf("Please use end_time url parameter in zulu format (like %s) to indicate the query end time", zuluForm),
			})
		}
	} else {
		endDate := queries.Get("end_time")
		errEnd := validateDate(endDate)
		if errEnd != nil {
			errs = append(errs, ErrorResponse{
				Message: "end_time parsing error",
				Code:    fmt.Sprintf("%d", http.StatusBadRequest),
				Details: fmt.Sprintf("Error parsing date string %s please use zulu format like %s", endDate, zuluForm),
			})
		}
	}
	if queries["start_time"] == nil {
		errs = append(errs, ErrorResponse{
			Message: "start_time not set",
			Code:    fmt.Sprintf("%d", http.StatusBadRequest),
			Details: fmt.Sprintf("Please use start_time url parameter in zulu format (like %s) to indicate the query start time", zuluForm),
		})
	} else {
		startDate := queries.Get("start_time")
		errStart := validateDate(startDate)
		if errStart != nil {
			errs = append(errs, ErrorResponse{
				Message: "start_time parsing error",
				Code:    fmt.Sprintf("%d", http.StatusBadRequest),
				Details: fmt.Sprintf("Error parsing date string %s please use zulu format like %s", startDate, zuluForm),
			})
		}
	}

	return errs

}

// ValidateAcceptHeader parses the accept header to determine the content type
func ValidateAcceptHeader(accept string) ErrorResponse {

	// err := ErrorResponse{}

	if accept == acceptedContentTypes[0] {
		return (ErrorResponse{})
	}
	if accept == acceptedContentTypes[1] {
		return (ErrorResponse{})
	}

	err := ErrorResponse{
		Message: "Not Acceptable Content Type",
		Code:    fmt.Sprintf("%d", http.StatusNotAcceptable),
		Details: fmt.Sprintf("Accept header provided did not contain any valid content types. Acceptable content types are '%s' and '%s'", acceptedContentTypes[0], acceptedContentTypes[1]),
	}
	return err

}
