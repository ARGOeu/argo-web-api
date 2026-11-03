package resultsV5

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (query *basicQuery) Validate() []ErrorResponse {

	errs := []ErrorResponse{}
	query.Granularity = strings.ToLower(query.Granularity)
	if query.Granularity == "" {
		query.Granularity = "daily"
	} else if query.Granularity != "daily" && query.Granularity != "monthly" && query.Granularity != "custom" {
		errs = append(errs, ErrorResponse{
			Message: "Wrong Granularity",
			Code:    "400",
			Details: fmt.Sprintf("%s is not acceptesd as granularity parameter, please provide either daily, monthly or custom", query.Granularity),
		})
	}

	if query.StartTime == "" && query.EndTime == "" {
		errs = append(errs, ErrorResponse{
			Message: "No time span set",
			Code:    "400",
			Details: "Please use start-time and/or end-time url parameters to set the prefered time span",
		})
	} else {
		if query.StartTime != "" {
			ts, tserr := time.Parse(zuluForm, query.StartTime)
			if tserr != nil {
				errs = append(errs, ErrorResponse{
					Message: "start-time parsing error",
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
					Message: "end-time parsing error",
					Code:    "400",
					Details: fmt.Sprintf("Error parsing date string %s please use zulu format like %s", query.EndTime, zuluForm),
				})
			}
			query.EndTimeInt, _ = strconv.Atoi(te.Format(ymdForm))
		}
	}

	return errs
}
