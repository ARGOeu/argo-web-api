package feeds

import (
	"encoding/json"
	"strconv"

	"github.com/ARGOeu/argo-web-api/respond"
)

// createTopoListView constructs the list response template and exports it as json
func createTopoListView(results []Topo, msg string, code int) ([]byte, error) {

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

// createWeightsListView constructs the list response template and exports it as json
func createWeightsListView(results []Weights, msg string, code int) ([]byte, error) {

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
