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

package trends

import (
	"encoding/json"
	"strconv"

	"github.com/ARGOeu/argo-web-api/respond"
)

// createStatusEgroupListView constructs the list response template and exports it as json
func createStatusEgroupListView(results []StatusGroupEgroupData, msg string, code int) ([]byte, error) {

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

// createTagMetricDataListView constructs the list response template and exports it as json
func createTagMetricDataListView(results []TagMetricData, msg string, code int) ([]byte, error) {

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

// createStatusServiceListView constructs the list response template and exports it as json
func createStatusServiceListView(results []StatusGroupServiceData, msg string, code int) ([]byte, error) {

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

// createStatusEndpointListView constructs the list response template and exports it as json
func createStatusEndpointListView(results []StatusGroupEndpointData, msg string, code int) ([]byte, error) {

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

// createMetricListView constructs the list response template and exports it as json
func createStatusMetricListView(results []StatusGroupMetricData, msg string, code int) ([]byte, error) {

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

// createMetricListView constructs the list response template and exports it as json
func createMetricListView(results []MetricData, msg string, code int) ([]byte, error) {

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

// createServiceListView constructs the list response template and exports it as json
func createServiceListView(results []ServiceData, msg string, code int) ([]byte, error) {

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

// createEndpointGroupListView constructs the list response template and exports it as json
func createEndpointGroupListView(results []EndpointGroupData, msg string, code int) ([]byte, error) {

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

// createStatusMonthEgroupListView constructs the list response template and exports it as json
func createStatusMonthEgroupListView(results []StatusMonthEgroupData, msg string, code int) ([]byte, error) {

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

// createStatusMonthServiceListView constructs the list response template and exports it as json
func createStatusMonthServiceListView(results []StatusMonthServiceData, msg string, code int) ([]byte, error) {

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

// createStatusMonthEndpointListView constructs the list response template and exports it as json
func createStatusMonthEndpointListView(results []StatusMonthEndpointData, msg string, code int) ([]byte, error) {

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

// createStatusMonthMetricListView constructs the list response template and exports it as json
func createStatusMonthMetricListView(results []StatusMonthMetricData, msg string, code int) ([]byte, error) {

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

// createMonthMetricListView constructs the list response template and exports it as json
func createMonthMetricListView(results []MonthMetricData, msg string, code int) ([]byte, error) {

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

// createMonthEndpointListView constructs the list response template and exports it as json
func createMonthEndpointListView(results []MonthEndpointData, msg string, code int) ([]byte, error) {

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

// createMonthEndpointListView constructs the list response template and exports it as json
func createMonthServiceListView(results []MonthServiceData, msg string, code int) ([]byte, error) {

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

// createMonthEndpointGroupListView constructs the list response template and exports it as json
func createMonthEndpointGroupListView(results []MonthEndpointGroupData, msg string, code int) ([]byte, error) {

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
