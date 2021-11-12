package metrics

import (
	"encoding/json"
	"strconv"

	"github.com/ARGOeu/argo-web-api/respond"
)

// createMetricsListView constructs the list response template and exports it as json
func createMetricsListView(results []Metric, msg string, code int) ([]byte, error) {

	docRoot := &respond.ResponseMessage{
		Status: respond.StatusResponse{
			Message: msg,
			Code:    strconv.Itoa(code),
		},
	}
	docRoot.Data = results

	output, err := json.MarshalIndent(docRoot, "", " ")
	return output, err

}

// createMsgView constructs a simple message response without data
func createMsgView(msg string, code int) ([]byte, error) {
	docRoot := &respond.ResponseMessage{
		Status: respond.StatusResponse{
			Message: msg,
			Code:    strconv.Itoa(code),
		},
	}

	output, err := json.MarshalIndent(docRoot, "", " ")
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
