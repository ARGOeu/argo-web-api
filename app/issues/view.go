/*
 * Copyright (c) 2015 GRNET S.A.
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

package issues

import (
	"encoding/json"
	"strconv"

	"github.com/ARGOeu/argo-web-api/respond"
)

// createEndpointListView constructs the list response template and exports it as json
func createEndpointListView(results []EndpointData, msg string, code int) ([]byte, error) {

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

// createGroupMetricsView constructs the json response for listing metric issues of a group
func createGroupMetricsView(results []GroupMetrics, msg string, code int) ([]byte, error) {

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
