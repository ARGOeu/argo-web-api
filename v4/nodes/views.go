package nodes

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/respond"
)

func createErrorMessage(message string, code int, format string) ([]byte, error) {

	var output []byte
	err := error(nil)
	docRoot := &errorMessage{}

	docRoot.Message = message
	docRoot.Code = code
	if strings.EqualFold(format, "application/json") {
		output, err = json.MarshalIndent(docRoot, " ", "  ")
	} else {
		output, err = xml.MarshalIndent(docRoot, " ", "  ")
	}
	return output, err
}

func createSummaryView(results []GroupInterface) ([]byte, error) {

	docRoot := Data[Summary]{
		Data: []Results[Summary]{},
	}

	prevEndpoint := ""

	for i := 0; i < len(results); i++ {
		row := results[i]
		timestamp, _ := time.Parse(customForm[0], fmt.Sprint(row.Date))

		if prevEndpoint != row.Name {
			prevEndpoint = row.Name
			docRoot.Data = append(docRoot.Data, Results[Summary]{
				Name:    row.Name,
				Results: []Summary{},
			})
		}

		currIdx := len(docRoot.Data) - 1
		prepDate := timestamp.Format(customForm[1])

		docRoot.Data[currIdx].Results = append(docRoot.Data[currIdx].Results,
			Summary{
				Date:         prepDate,
				Availability: fmt.Sprintf("%g", row.Availability),
				Uptime:       fmt.Sprintf("%g", row.Up),
			})
	}
	return json.MarshalIndent(docRoot, " ", "  ")

}

func createAvailabilityView(results []GroupInterface) ([]byte, error) {

	docRoot := Data[Availability]{
		Data: []Results[Availability]{},
	}

	prevEndpoint := ""

	for i := 0; i < len(results); i++ {
		row := results[i]
		timestamp, _ := time.Parse(customForm[0], fmt.Sprint(row.Date))

		if prevEndpoint != row.Name {
			prevEndpoint = row.Name
			docRoot.Data = append(docRoot.Data, Results[Availability]{
				Name:    row.Name,
				Results: []Availability{},
			})
		}

		currIdx := len(docRoot.Data) - 1
		prepDate := timestamp.Format(customForm[1])

		docRoot.Data[currIdx].Results = append(docRoot.Data[currIdx].Results,
			Availability{
				Date:         prepDate,
				Availability: fmt.Sprintf("%g", row.Availability),
			})
	}
	return json.MarshalIndent(docRoot, " ", "  ")

}

func createUptimeView(results []GroupInterface) ([]byte, error) {

	docRoot := Data[Uptime]{
		Data: []Results[Uptime]{},
	}

	prevEndpoint := ""

	for i := 0; i < len(results); i++ {
		row := results[i]
		timestamp, _ := time.Parse(customForm[0], fmt.Sprint(row.Date))

		if prevEndpoint != row.Name {
			prevEndpoint = row.Name
			docRoot.Data = append(docRoot.Data, Results[Uptime]{
				Name:    row.Name,
				Results: []Uptime{},
			})
		}

		currIdx := len(docRoot.Data) - 1
		prepDate := timestamp.Format(customForm[1])

		docRoot.Data[currIdx].Results = append(docRoot.Data[currIdx].Results,
			Uptime{
				Date:   prepDate,
				Uptime: fmt.Sprintf("%g", row.Up),
			})
	}
	return json.MarshalIndent(docRoot, " ", "  ")

}

func createStatusView(results []GroupStatusData, input InputStatus, endDate string, details bool, latest bool) ([]byte, error) {

	var extraTS string
	tsNow := time.Now().UTC()
	today := tsNow.Format("2006-01-02")

	if endDate == today {
		extraTS = "T" + strings.Split(tsNow.Format(zuluForm), "T")[1]
	} else {
		extraTS = "T23:59:59Z"
	}

	docRoot := StatusData{
		Data: []Status{},
	}

	if len(results) == 0 {
		return respond.MarshalContent(docRoot, input.format, "", " ")
	}

	prevEndpointGroup := ""
	var ppEndpointGroup *Status

	for _, row := range results {
		if row.Group != prevEndpointGroup && row.Group != "" {

			if !latest && ppEndpointGroup != nil && len(ppEndpointGroup.Results) > 0 {
				latestStatus := ppEndpointGroup.Results[len(ppEndpointGroup.Results)-1]
				ppEndpointGroup.Results = append(ppEndpointGroup.Results, StatusResult{
					Timestamp: strings.Split(latestStatus.Timestamp, "T")[0] + extraTS,
					Value:     latestStatus.Value,
				})
			}

			docRoot.Data = append(docRoot.Data, Status{
				Name:    row.Group,
				Results: []StatusResult{},
			})
			prevEndpointGroup = row.Group
			ppEndpointGroup = &docRoot.Data[len(docRoot.Data)-1]
		}

		status := StatusResult{
			Timestamp: row.Timestamp,
			Value:     row.Status,
		}
		if details {
			status.AffectedByThresholdRule = row.HasThresholdRule
		}

		if latest {
			if endDate == today {
				status.Timestamp = strings.Split(status.Timestamp, "T")[0] + extraTS
			}
			ppEndpointGroup.Results = []StatusResult{status}
		} else {
			ppEndpointGroup.Results = append(ppEndpointGroup.Results, status)
		}
	}

	if !latest && ppEndpointGroup != nil && len(ppEndpointGroup.Results) > 0 {
		latestStatus := ppEndpointGroup.Results[len(ppEndpointGroup.Results)-1]
		ppEndpointGroup.Results = append(ppEndpointGroup.Results, StatusResult{
			Timestamp: strings.Split(latestStatus.Timestamp, "T")[0] + extraTS,
			Value:     latestStatus.Value,
		})
	}

	return respond.MarshalContent(docRoot, input.format, "", " ")

}
