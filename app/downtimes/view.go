package downtimes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ARGOeu/argo-web-api/respond"
)

func createRefView(inserted Downtimes, msg string, code int, r *http.Request) ([]byte, error) {
	docRoot := &respond.ResponseMessage{
		Status: respond.StatusResponse{
			Message: msg,
			Code:    strconv.Itoa(code),
		},
		Data: SelfReference{
			ID:    inserted.ID,
			Links: Links{Self: "https://" + r.Host + r.URL.Path + "/" + inserted.ID},
		},
	}

	output, err := json.MarshalIndent(docRoot, "", " ")
	return output, err
}

// createListView constructs the list response template and exports it as json
func createListView(results []Downtimes, msg string, code int) ([]byte, error) {

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
