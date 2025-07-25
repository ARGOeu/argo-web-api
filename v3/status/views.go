package status

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/respond"
)

func createCombinedView(resGroups []GroupData, resEndpoints []EndpointData, input InputParams, endDate string, details bool, latest bool) ([]byte, error) {

	// calculate part of the timestamp that closes the timeline of each item
	var extraTS string

	tsNow := time.Now().UTC()
	today := tsNow.Format("2006-01-02")

	if endDate == today {
		extraTS = "T" + strings.Split(tsNow.Format(zuluForm), "T")[1]
	} else {
		extraTS = "T23:59:59Z"
	}

	var output []byte
	err := error(nil)

	docRoot := &rootOUT{}

	if len(resGroups) == 0 {
		output, err = json.MarshalIndent(docRoot, " ", "  ")
		return output, err
	}

	// make index map to keep track of different groups
	indexGroup := make(map[string]*groupOUT)
	keysGroup := make([]string, 0)

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
			// add key to keysGroup array to be used in sorted traversal
			keysGroup = append(keysGroup, row.Group)
		}

	}

	// make index map to keep track of different endpoint
	indexEndp := make(map[string]*endpointOUT)
	keysEndp := make([]string, 0)

	for _, row := range resEndpoints {
		// prepare status information to be added
		status := &statusOUT{}
		status.Timestamp = row.Timestamp
		status.Value = row.Status
		if details {
			status.AffectedByThresholdRule = row.HasThresholdRule
		}

		// check if item has an already created group
		if ptrEndp, ok := indexEndp[row.EndpointGroup+row.Service+row.Hostname]; ok {
			ptrEndp.Statuses = append(ptrEndp.Statuses, status)
		} else {
			newEndp := &endpointOUT{}
			newEndp.Info = row.Info
			newEndp.Name = row.Hostname
			newEndp.Service = row.Service
			newEndp.SuperGroup = row.EndpointGroup
			newEndp.Statuses = make([]*statusOUT, 0)
			newEndp.Statuses = append(newEndp.Statuses, status)
			indexEndp[row.EndpointGroup+row.Service+row.Hostname] = newEndp
			// add key to keysEndp to be used in sorted traversal
			keysEndp = append(keysEndp, row.EndpointGroup+row.Service+row.Hostname)
		}

	}

	// sort keys
	sort.Strings(keysGroup)
	sort.Strings(keysEndp)

	// repeat over group items and add them to the response root

	for _, key := range keysEndp {
		value := indexEndp[key]
		// check if endpoint supergroup is indexed in groups
		if ptrGroup, ok := indexGroup[value.SuperGroup]; ok {

			// add extra status that closes the timeline
			// get last status of the existing timeline
			lastStatus := value.Statuses[len(value.Statuses)-1]
			extraStatus := &statusOUT{}
			extraStatus.Timestamp = strings.Split(lastStatus.Timestamp, "T")[0] + extraTS
			extraStatus.Value = lastStatus.Value
			if latest {
				value.Statuses = nil
			}
			value.Statuses = append(value.Statuses, extraStatus)

			ptrGroup.Endpoints = append(ptrGroup.Endpoints, value)
			// clear supergroup value from endpoint it self so as not to be printed in json output
			value.SuperGroup = ""

		}

	}

	// repeat over group items and add them to the response root
	for _, key := range keysGroup {
		value := indexGroup[key]
		// add extra status that closes the timeline
		// get last status of the existing timeline
		lastStatus := value.Statuses[len(value.Statuses)-1]
		extraStatus := &statusOUT{}
		extraStatus.Timestamp = strings.Split(lastStatus.Timestamp, "T")[0] + extraTS
		extraStatus.Value = lastStatus.Value
		if latest {
			value.Statuses = nil
		}
		value.Statuses = append(value.Statuses, extraStatus)
		docRoot.Groups = append(docRoot.Groups, value)
	}

	output, err = respond.MarshalContent(docRoot, input.format, "", " ")
	return output, err

}

func createViewByID(resEndpoints []EndpointData, input InputParams, endDate string, details bool, latest bool) ([]byte, error) {

	// calculate part of the timestamp that closes the timeline of each item
	var extraTS string

	tsNow := time.Now().UTC()
	today := tsNow.Format("2006-01-02")

	if endDate == today {
		extraTS = "T" + strings.Split(tsNow.Format(zuluForm), "T")[1]
	} else {
		extraTS = "T23:59:59Z"
	}

	var output []byte
	err := error(nil)

	docID := &idOUT{}
	docID.ID = input.ID
	docID.Endpoints = make([]*endpointOUT, 0)

	// make index map to keep track of different endpoint
	indexEndp := make(map[string]*endpointOUT)
	keysEndp := make([]string, 0)

	for _, row := range resEndpoints {
		// prepare status information to be added
		status := &statusOUT{}
		status.Timestamp = row.Timestamp
		status.Value = row.Status
		if details {
			status.AffectedByThresholdRule = row.HasThresholdRule
		}

		if ptrEndp, ok := indexEndp[row.EndpointGroup+row.Service+row.Hostname]; ok {
			ptrEndp.Statuses = append(ptrEndp.Statuses, status)
		} else {
			newEndp := &endpointOUT{}
			newEndp.Info = row.Info
			newEndp.Name = row.Hostname
			newEndp.Service = row.Service
			newEndp.SuperGroup = row.EndpointGroup
			newEndp.Statuses = make([]*statusOUT, 0)
			newEndp.Statuses = append(newEndp.Statuses, status)
			indexEndp[row.EndpointGroup+row.Service+row.Hostname] = newEndp
			// add key to keysEndp to be used in sorted traversal
			keysEndp = append(keysEndp, row.EndpointGroup+row.Service+row.Hostname)
		}

	}

	sort.Strings(keysEndp)

	// repeat over group items and add them to the response root

	for _, key := range keysEndp {
		value := indexEndp[key]
		// check if endpoint supergroup is indexed in groups

		// add extra status that closes the timeline
		// get last status of the existing timeline
		lastStatus := value.Statuses[len(value.Statuses)-1]
		extraStatus := &statusOUT{}
		extraStatus.Timestamp = strings.Split(lastStatus.Timestamp, "T")[0] + extraTS
		extraStatus.Value = lastStatus.Value
		if latest {
			value.Statuses = nil
		}
		value.Statuses = append(value.Statuses, extraStatus)
		docID.Endpoints = append(docID.Endpoints, value)

	}

	output, err = respond.MarshalContent(docID, input.format, "", " ")
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

func createErrorMessage(message string, code int) ([]byte, error) {

	output := []byte("message placeholder")
	err := error(nil)
	docRoot := &errorMessage{}

	docRoot.Message = message
	docRoot.Code = code

	output, err = json.MarshalIndent(docRoot, " ", "  ")

	return output, err
}
