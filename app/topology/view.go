/*
 * Copyright (c) 2018 GRNET S.A.
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
 * or implied, of GRNET S.A.
 *
 */

package topology

import (
	"encoding/json"
	"strconv"

	"github.com/ARGOeu/argo-web-api/respond"
)

func createTopoView(results Topology, msg string, code int) ([]byte, error) {

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

func createListEndpoint(results []Endpoint, msg string, code int) ([]byte, error) {

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

func createListTags(results []TagInfo, msg string, code int) ([]byte, error) {

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

func createListGroup(results []Group, msg string, code int) ([]byte, error) {

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

func createMessageOUT(message string, code int, format string) ([]byte, error) {

	output := []byte("message placeholder")
	err := error(nil)
	docRoot := &messageOUT{}

	docRoot.Message = message
	docRoot.Code = strconv.Itoa(code)
	output, err = respond.MarshalContent(docRoot, format, "", " ")
	return output, err
}
