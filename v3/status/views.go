package status

import (
	"encoding/json"
	"encoding/xml"
	"strconv"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/respond"
)

func createCombinedView(resGroups []GroupData, resEndpoints []EndpointData, input InputParams, endDate string, details bool) ([]byte, error) {

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

	docRoot := &rootOUT{}

	if len(resGroups) == 0 {
		output, err = json.MarshalIndent(docRoot, " ", "  ")
		return output, err
	}

	// make index map to keep track of different groups
	indexGroup := make(map[string]*groupOUT)

	for _, row := range resGroups {
		// prepare status information to be added
		status := &statusOUT{}
		status.Timestamp = row.Timestamp
		status.Value = row.Status
		if details {
			status.AffectedByThresholdRule = row.HasThresholdRule
		}

		// check if item has an already created group
		if ptrGroup, ok := indexGroup[row.Group]; ok {
			ptrGroup.Statuses = append(ptrGroup.Statuses, status)
		} else {
			newGroup := &groupOUT{}
			newGroup.GroupType = input.groupType
			newGroup.Name = row.Group
			newGroup.Statuses = make([]*statusOUT, 0)
			newGroup.Endpoints = make([]*endpointOUT, 0)
			newGroup.Statuses = append(newGroup.Statuses, status)
			indexGroup[row.Group] = newGroup
		}

	}

	// make index map to keep track of different endpoint
	indexEndp := make(map[string]*endpointOUT)

	for _, row := range resEndpoints {
		// prepare status information to be added
		status := &statusOUT{}
		status.Timestamp = row.Timestamp
		status.Value = row.Status
		if details {
			status.AffectedByThresholdRule = row.HasThresholdRule
		}

		// check if item has an already created group
		if ptrEndp, ok := indexEndp[row.Hostname+row.Service]; ok {
			ptrEndp.Statuses = append(ptrEndp.Statuses, status)
		} else {
			newEndp := &endpointOUT{}
			newEndp.Info = row.Info
			newEndp.Name = row.Hostname
			newEndp.Service = row.Service
			newEndp.SuperGroup = row.EndpointGroup
			newEndp.Statuses = make([]*statusOUT, 0)
			newEndp.Statuses = append(newEndp.Statuses, status)
			indexEndp[row.Hostname+row.Service] = newEndp
		}

	}

	// repeat over group items and add them to the response root

	for _, value := range indexEndp {
		// check if endpoint supergroup is indexed in groups
		if ptrGroup, ok := indexGroup[value.SuperGroup]; ok {

			// add extra status that closes the timeline
			// get last status of the existing timeline
			lastStatus := value.Statuses[len(value.Statuses)-1]
			extraStatus := &statusOUT{}
			extraStatus.Timestamp = strings.Split(lastStatus.Timestamp, "T")[0] + extraTS
			extraStatus.Value = lastStatus.Value
			value.Statuses = append(value.Statuses, extraStatus)
			ptrGroup.Endpoints = append(ptrGroup.Endpoints, value)

		}

	}

	// repeat over group items and add them to the response root
	for _, value := range indexGroup {
		// add extra status that closes the timeline
		// get last status of the existing timeline
		lastStatus := value.Statuses[len(value.Statuses)-1]
		extraStatus := &statusOUT{}
		extraStatus.Timestamp = strings.Split(lastStatus.Timestamp, "T")[0] + extraTS
		extraStatus.Value = lastStatus.Value
		value.Statuses = append(value.Statuses, extraStatus)
		docRoot.Groups = append(docRoot.Groups, value)
	}

	output, err = respond.MarshalContent(docRoot, input.format, "", " ")
	return output, err

}

func createMessageOUT(message string, code int, format string) ([]byte, error) {

	output := []byte("message placeholder")
	err := error(nil)
	docRoot := &messageOUT{}

	docRoot.Message = message
	docRoot.Code = strconv.Itoa(code)
	output, err = respond.MarshalContent(docRoot, format, "", " ")
	return output, err
}

func createErrorMessage(message string, code int, format string) ([]byte, error) {

	output := []byte("message placeholder")
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
