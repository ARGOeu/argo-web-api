package statusV5

import (
	"encoding/json"
	"encoding/xml"
	"strconv"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/respond"
)

func createGroupView(report reports.MongoInterface, results []GroupDataOutput, input GroupInputParams, endDate string, details bool) ([]byte, error) {

	// calculate part of the timestamp that closes the timeline of each item
	var extraTS string

	tsNow := time.Now().UTC()
	today := tsNow.Format("2006-01-02")

	if strings.Split(endDate, "T")[0] == today {
		extraTS = "T" + strings.Split(tsNow.Format(zuluForm), "T")[1]
	} else {
		extraTS = "T23:59:59Z"
	}

	var output []byte
	err := error(nil)

	docRoot := &rootGroupOUT{}

	if len(results) == 0 {
		if strings.EqualFold(input.format, "application/json") {
			output, err = json.MarshalIndent(docRoot, " ", "  ")
		} else {
			output, err = xml.MarshalIndent(docRoot, " ", "  ")
		}
		return output, err
	}

	prevGroup := ""

	var ppGroup *groupOUT

	for _, row := range results {

		if row.Group != prevGroup && row.Group != "" {
			// close the status timeline by adding a new status item at 23:59 or at current time
			if ppGroup != nil {
				eStatus := &statusOUT{}
				latestStatus := ppGroup.Statuses[len(ppGroup.Statuses)-1]
				eStatus.Timestamp = strings.Split(latestStatus.Timestamp, "T")[0] + extraTS
				eStatus.Value = latestStatus.Value
				ppGroup.Statuses = append(ppGroup.Statuses, eStatus)
			}

			group := &groupOUT{}
			group.Name = row.Group
			group.GroupType = report.Topology.Group.Group.Type
			docRoot.Groups = append(docRoot.Groups, group)
			prevGroup = row.Group
			ppGroup = group
		}

		status := &statusOUT{}
		status.Timestamp = row.Timestamp
		status.Value = row.Status
		if details {
			status.AffectedByThresholdRule = row.HasThresholdRule
		}

		ppGroup.Statuses = append(ppGroup.Statuses, status)

	}

	// close the status timeline of the last item by adding a new status item at 23:59 or at current time
	if ppGroup != nil {
		eStatus := &statusOUT{}
		latestStatus := ppGroup.Statuses[len(ppGroup.Statuses)-1]
		eStatus.Timestamp = strings.Split(latestStatus.Timestamp, "T")[0] + extraTS
		eStatus.Value = latestStatus.Value
		ppGroup.Statuses = append(ppGroup.Statuses, eStatus)
	}

	output, err = respond.MarshalContent(docRoot, input.format, "", " ")
	return output, err

}

func createServiceView(report reports.MongoInterface, results []ServiceDataOutput, input ServiceInputParams, endDate string, details bool) ([]byte, error) {

	// calculate part of the timestamp that closes the timeline of each item
	var extraTS string

	tsNow := time.Now().UTC()
	today := tsNow.Format("2006-01-02")

	if strings.Split(endDate, "T")[0] == today {
		extraTS = "T" + strings.Split(tsNow.Format(zuluForm), "T")[1]
	} else {
		extraTS = "T23:59:59Z"
	}

	output := []byte("reponse output")
	err := error(nil)

	docRoot := &rootServiceOUT{}

	if len(results) == 0 {
		if strings.EqualFold(input.format, "application/json") {
			output, err = json.MarshalIndent(docRoot, " ", "  ")
		} else {
			output, err = xml.MarshalIndent(docRoot, " ", "  ")
		}
		return output, err
	}

	prevGroup := ""
	prevService := ""

	var ppEndpointGroup *serviceGroupOUT
	var ppService *serviceOUT

	for _, row := range results {

		if row.Group != prevGroup && row.Group != "" {
			endpointGroup := &serviceGroupOUT{}
			endpointGroup.Name = row.Group
			endpointGroup.GroupType = report.Topology.Group.Group.Type
			docRoot.Groups = append(docRoot.Groups, endpointGroup)
			prevGroup = row.Group
			ppEndpointGroup = endpointGroup
		}

		if row.Service != prevService && row.Service != "" {

			// close the status timeline of item by adding a new status item at 23:59 or at current time
			if ppService != nil {
				eStatus := &statusOUT{}
				latestStatus := ppService.Statuses[len(ppService.Statuses)-1]
				eStatus.Timestamp = strings.Split(latestStatus.Timestamp, "T")[0] + extraTS
				eStatus.Value = latestStatus.Value
				ppService.Statuses = append(ppService.Statuses, eStatus)
			}

			service := &serviceOUT{}
			service.Name = row.Service
			service.GroupType = "service"
			ppEndpointGroup.Services = append(ppEndpointGroup.Services, service)

			prevService = row.Service
			ppService = service
		}

		status := &statusOUT{}
		status.Timestamp = row.Timestamp
		status.Value = row.Status
		if details {
			status.AffectedByThresholdRule = row.HasThresholdRule
		}

		ppService.Statuses = append(ppService.Statuses, status)

	}

	// close the status timeline of the last item by adding a new status item at 23:59 or at current time
	if ppService != nil {
		eStatus := &statusOUT{}
		latestStatus := ppService.Statuses[len(ppService.Statuses)-1]
		eStatus.Timestamp = strings.Split(latestStatus.Timestamp, "T")[0] + extraTS
		eStatus.Value = latestStatus.Value
		ppService.Statuses = append(ppService.Statuses, eStatus)
	}

	output, err = respond.MarshalContent(docRoot, input.format, "", " ")
	return output, err

}

func createEndpointView(report reports.MongoInterface, results []EndpointDataOutput, input EndpointInputParams, endDate string, details bool) ([]byte, error) {

	// calculate part of the timestamp that closes the timeline of each item
	var extraTS string

	tsNow := time.Now().UTC()
	today := tsNow.Format("2006-01-02")

	if strings.Split(endDate, "T")[0] == today {
		extraTS = "T" + strings.Split(tsNow.Format(zuluForm), "T")[1]
	} else {
		extraTS = "T23:59:59Z"
	}

	var output []byte
	err := error(nil)

	docRoot := &rootEndpointOUT{}

	if len(results) == 0 {
		if strings.EqualFold(input.format, "application/json") {
			output, err = json.MarshalIndent(docRoot, " ", "  ")
		} else {
			output, err = xml.MarshalIndent(docRoot, " ", "  ")
		}
		return output, err
	}

	prevHostname := ""
	prevGroup := ""
	prevService := ""

	var ppHost *endpointOUT
	var ppGroup *endpointGroupOUT
	var ppService *endpointServiceOUT

	for _, row := range results {

		if row.Group != prevGroup && row.Group != "" {
			endpointGroup := &endpointGroupOUT{}
			endpointGroup.Name = row.Group
			endpointGroup.GroupType = report.Topology.Group.Group.Type
			docRoot.Groups = append(docRoot.Groups, endpointGroup)
			prevGroup = row.Group
			ppGroup = endpointGroup
		}

		if row.Service != prevService && row.Service != "" {
			service := &endpointServiceOUT{}
			service.Name = row.Service
			service.GroupType = "service"
			ppGroup.Services = append(ppGroup.Services, service)

			prevService = row.Service
			ppService = service
		}

		if row.Hostname != prevHostname && row.Hostname != "" {
			// close the status timeline of item by adding a new status item at 23:59 or at current time
			if ppHost != nil {
				eStatus := &statusOUT{}
				latestStatus := ppHost.Statuses[len(ppHost.Statuses)-1]
				eStatus.Timestamp = strings.Split(latestStatus.Timestamp, "T")[0] + extraTS
				eStatus.Value = latestStatus.Value
				ppHost.Statuses = append(ppHost.Statuses, eStatus)
			}

			host := &endpointOUT{} //create new host
			host.Name = row.Hostname
			host.Info = row.Info
			ppService.Endpoints = append(ppService.Endpoints, host)
			prevHostname = row.Hostname
			ppHost = host
		}

		status := &statusOUT{}
		status.Timestamp = row.Timestamp
		status.Value = row.Status
		if details {
			status.AffectedByThresholdRule = row.HasThresholdRule
		}

		ppHost.Statuses = append(ppHost.Statuses, status)

	}
	// close the status timeline of the last item by adding a new status item at 23:59 or at current time
	if ppHost != nil {
		eStatus := &statusOUT{}
		latestStatus := ppHost.Statuses[len(ppHost.Statuses)-1]
		eStatus.Timestamp = strings.Split(latestStatus.Timestamp, "T")[0] + extraTS
		eStatus.Value = latestStatus.Value
		ppHost.Statuses = append(ppHost.Statuses, eStatus)
	}

	output, err = respond.MarshalContent(docRoot, input.format, "", " ")
	return output, err

}

func createMetricView(report reports.MongoInterface, results []MetricDataOutput, input MetricInputParams, details bool) ([]byte, error) {

	var output []byte
	err := error(nil)

	docRoot := &rootMetricOUT{}

	if len(results) == 0 {
		if strings.EqualFold(input.format, "application/json") {
			output, err = json.MarshalIndent(docRoot, " ", "  ")
		} else {
			output, err = xml.MarshalIndent(docRoot, " ", "  ")
		}
		return output, err
	}

	prevHostname := ""
	prevMetric := ""
	prevGroup := ""
	prevService := ""

	var ppHost *metricEndpointOUT
	var ppMetric *metricOUT
	var ppGroup *metricGroupOUT
	var ppService *metricServiceOUT

	for _, row := range results {

		if row.Group != prevGroup && row.Group != "" {
			group := &metricGroupOUT{}
			group.Name = row.Group
			group.GroupType = report.Topology.Group.Group.Type
			docRoot.Groups = append(docRoot.Groups, group)
			prevGroup = row.Group
			ppGroup = group
		}

		if row.Service != prevService && row.Service != "" {
			service := &metricServiceOUT{}
			service.Name = row.Service
			service.GroupType = "service"
			ppGroup.Services = append(ppGroup.Services, service)

			prevService = row.Service
			ppService = service
		}

		if row.Hostname != prevHostname && row.Hostname != "" {
			host := &metricEndpointOUT{} //create new host
			host.Name = row.Hostname
			host.Info = row.Info
			ppService.Endpoints = append(ppService.Endpoints, host)
			prevHostname = row.Hostname
			ppHost = host
		}

		if row.Metric != prevMetric {

			metric := &metricOUT{}
			//Add the prev status as the firstone

			metric.Name = row.Metric
			ppHost.Metrics = append(ppHost.Metrics, metric)
			prevMetric = row.Metric
			ppMetric = metric

			prevStatus := &metricStatusOUT{}
			prevStatus.Timestamp = row.PrevTimestamp
			prevStatus.Value = row.PrevStatus
			ppMetric.Statuses = append(ppMetric.Statuses, prevStatus)

		}

		status := &metricStatusOUT{}
		status.Timestamp = row.Timestamp
		status.Value = row.Status
		if details {
			status.ActualData = row.ActualData
			status.OriginalStatus = row.OriginalStatus
			status.RuleApplied = row.RuleApplied
		}
		ppMetric.Statuses = append(ppMetric.Statuses, status)

	}

	output, err = respond.MarshalContent(docRoot, input.format, "", " ")
	return output, err

}

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

func createMessageOUT(message string, code int, format string) ([]byte, error) {

	var output []byte
	err := error(nil)
	docRoot := &messageOUT{}

	docRoot.Message = message
	docRoot.Code = strconv.Itoa(code)
	output, err = respond.MarshalContent(docRoot, format, "", " ")
	return output, err
}
