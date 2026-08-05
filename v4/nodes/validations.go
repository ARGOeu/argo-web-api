package nodes

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/respond"
)

// ErrorResponse shortcut to respond.ErrorResponse
type ErrorResponse respond.ErrorResponse

func (query *basicQuery) Validate() []ErrorResponse {
	errs := []ErrorResponse{}
	query.Granularity = strings.ToLower(query.Granularity)
	if query.Granularity == "" {
		query.Granularity = "daily"
	} else if query.Granularity != "daily" && query.Granularity != "monthly" && query.Granularity != "custom" {
		errs = append(errs, ErrorResponse{
			Message: "Wrong Granularity",
			Code:    "400",
			Details: fmt.Sprintf("%s is not accepted as granularity parameter, please provide either daily, monthly or custom", query.Granularity),
		})
	}

	if query.Date != "" {
		ts, tserr := time.Parse(dtForm, query.Date)
		if tserr != nil {
			errs = append(errs, ErrorResponse{
				Message: "date parsing error",
				Code:    "400",
				Details: fmt.Sprintf("Error parsing date string %s please use YYYY-MM-DD format like %s", query.StartTime, dtForm),
			})
			return errs
		}
		query.StartTime = time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
		query.EndTime = time.Date(ts.Year(), ts.Month(), ts.Day(), 23, 59, 59, 0, time.UTC).Format(time.RFC3339)
	}

	if query.StartDate != "" || query.EndDate != "" {
		if !(query.StartDate != "" && query.EndDate != "") {
			errs = append(errs, ErrorResponse{
				Message: "date parsing error",
				Code:    "400",
				Details: "You should speficy both start_date and end_date",
			})
			return errs
		}
		sd, tserr := time.Parse(dtForm, query.StartDate)
		if tserr != nil {
			errs = append(errs, ErrorResponse{
				Message: "date parsing error",
				Code:    "400",
				Details: fmt.Sprintf("Error parsing start_date string %s please use YYYY-MM-DD format like %s", query.StartTime, dtForm),
			})
		}
		ed, tserr := time.Parse(dtForm, query.EndDate)
		if ed.Before(sd) {
			errs = append(errs, ErrorResponse{
				Message: "date parsing error",
				Code:    "400",
				Details: "start_date should be a date before end_date",
			})
			return errs
		}
		if tserr != nil {
			errs = append(errs, ErrorResponse{
				Message: "date parsing error",
				Code:    "400",
				Details: fmt.Sprintf("Error parsing end_date string %s please use YYYY-MM-DD format like %s", query.StartTime, dtForm),
			})
		}

		query.StartTime = time.Date(sd.Year(), sd.Month(), sd.Day(), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
		query.EndTime = time.Date(ed.Year(), ed.Month(), ed.Day(), 23, 59, 59, 0, time.UTC).Format(time.RFC3339)
	}

	if query.StartTime == "" || query.EndTime == "" {
		now := time.Now().UTC()
		query.StartTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
		query.EndTime = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC).Format(time.RFC3339)

	}
	if query.StartTime != "" && query.EndTime != "" {
		ts, tserr := time.Parse(zuluForm, query.StartTime)
		if tserr != nil {
			errs = append(errs, ErrorResponse{
				Message: "start_time parsing error",
				Code:    "400",
				Details: fmt.Sprintf("Error parsing date string %s please use zulu format like %s", query.StartTime, zuluForm),
			})
			return errs
		}
		query.StartTimeInt, _ = strconv.Atoi(ts.Format(ymdForm))

		te, teerr := time.Parse(zuluForm, query.EndTime)
		if teerr != nil {
			errs = append(errs, ErrorResponse{
				Message: "end_time parsing error",
				Code:    "400",
				Details: fmt.Sprintf("Error parsing date string %s please use zulu format like %s", query.EndTime, zuluForm),
			})
			return errs
		}
		query.EndTimeInt, _ = strconv.Atoi(te.Format(ymdForm))

		if te.Before(ts) {
			errs = append(errs, ErrorResponse{
				Message: "time parsing error",
				Code:    "400",
				Details: "start_time should be a timestamp before end_time",
			})
			return errs
		}
	}

	return errs
}

func (query *basicQuery) ValidateV2() []ErrorResponse {
	errs := []ErrorResponse{}
	query.Granularity = strings.ToLower(query.Granularity)
	if query.Granularity == "" {
		query.Granularity = "daily"
	} else if query.Granularity != "daily" && query.Granularity != "monthly" && query.Granularity != "custom" {
		errs = append(errs, ErrorResponse{
			Message: "Wrong Granularity",
			Code:    "400",
			Details: fmt.Sprintf("%s is not accepted as granularity parameter, please provide either daily, monthly or custom", query.Granularity),
		})
	}

	if query.Date != "" {
		ts, tserr := time.Parse(dtForm, query.Date)
		if tserr != nil {
			errs = append(errs, ErrorResponse{
				Message: "date parsing error",
				Code:    "400",
				Details: fmt.Sprintf("Error parsing date string %s please use YYYY-MM-DD format like %s", query.StartTime, dtForm),
			})
			return errs
		}
		query.StartTime = time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
		query.EndTime = time.Date(ts.Year(), ts.Month(), ts.Day(), 23, 59, 59, 0, time.UTC).Format(time.RFC3339)
	}

	if query.StartDate != "" || query.EndDate != "" {
		if !(query.StartDate != "" && query.EndDate != "") {
			errs = append(errs, ErrorResponse{
				Message: "date parsing error",
				Code:    "400",
				Details: "You should speficy both start-date and end-date",
			})
			return errs
		}
		sd, tserr := time.Parse(dtForm, query.StartDate)
		if tserr != nil {
			errs = append(errs, ErrorResponse{
				Message: "date parsing error",
				Code:    "400",
				Details: fmt.Sprintf("Error parsing start-date string %s please use YYYY-MM-DD format like %s", query.StartTime, dtForm),
			})
		}
		ed, tserr := time.Parse(dtForm, query.EndDate)
		if ed.Before(sd) {
			errs = append(errs, ErrorResponse{
				Message: "date parsing error",
				Code:    "400",
				Details: "start-date should be a date before end-date",
			})
			return errs
		}
		if tserr != nil {
			errs = append(errs, ErrorResponse{
				Message: "date parsing error",
				Code:    "400",
				Details: fmt.Sprintf("Error parsing end-date string %s please use YYYY-MM-DD format like %s", query.StartTime, dtForm),
			})
		}

		query.StartTime = time.Date(sd.Year(), sd.Month(), sd.Day(), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
		query.EndTime = time.Date(ed.Year(), ed.Month(), ed.Day(), 23, 59, 59, 0, time.UTC).Format(time.RFC3339)
	}

	if query.StartTime == "" || query.EndTime == "" {
		now := time.Now().UTC()
		query.StartTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
		query.EndTime = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC).Format(time.RFC3339)

	}
	if query.StartTime != "" && query.EndTime != "" {
		ts, tserr := time.Parse(zuluForm, query.StartTime)
		if tserr != nil {
			errs = append(errs, ErrorResponse{
				Message: "start-time parsing error",
				Code:    "400",
				Details: fmt.Sprintf("Error parsing date string %s please use zulu format like %s", query.StartTime, zuluForm),
			})
			return errs
		}
		query.StartTimeInt, _ = strconv.Atoi(ts.Format(ymdForm))

		te, teerr := time.Parse(zuluForm, query.EndTime)
		if teerr != nil {
			errs = append(errs, ErrorResponse{
				Message: "end-time parsing error",
				Code:    "400",
				Details: fmt.Sprintf("Error parsing date string %s please use zulu format like %s", query.EndTime, zuluForm),
			})
			return errs
		}
		query.EndTimeInt, _ = strconv.Atoi(te.Format(ymdForm))

		if te.Before(ts) {
			errs = append(errs, ErrorResponse{
				Message: "time parsing error",
				Code:    "400",
				Details: "start-time should be a timestamp before end-time",
			})
			return errs
		}
	}

	return errs
}
