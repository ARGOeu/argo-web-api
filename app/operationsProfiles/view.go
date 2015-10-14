/*
 * Copyright (c) 2014 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
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

package operationsProfiles

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ARGOeu/argo-web-api/respond"
)

// createListView constructs the list response template and exports it as json
func createListView(results []OpsProfile, msg string, code int) ([]byte, error) {

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

// createListView constructs self-reference response and exports it as json
func createRefView(inserted OpsProfile, msg string, code int, r *http.Request) ([]byte, error) {
	docRoot := &respond.ResponseMessage{
		Status: respond.StatusResponse{
			Message: msg,
			Code:    strconv.Itoa(code),
		},
		Data: SelfReference{
			UUID:  inserted.UUID,
			Links: Links{Self: "https://" + r.Host + r.URL.Path + "/" + inserted.UUID},
		},
	}

	output, err := json.MarshalIndent(docRoot, "", " ")
	return output, err
}

// createErrView constructs a simple message response without data
func createErrView(msg string, code int, errList []string) ([]byte, error) {

	var errRespond []respond.ErrorResponse

	for _, item := range errList {
		temp := respond.ErrorResponse{"Validation Failed", strconv.Itoa(code), item}
		errRespond = append(errRespond, temp)
	}

	docRoot := &respond.ResponseMessage{
		Status: respond.StatusResponse{
			Message: msg,
			Code:    strconv.Itoa(code),
		},
		Errors: errRespond,
	}

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
