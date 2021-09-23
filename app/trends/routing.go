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
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/gorilla/mux"
)

// HandleSubrouter contains the different paths to follow during subrouting
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {

	s = respond.PrepAppRoutes(s, confhandler, appRoutesV2)
}

var appRoutesV2 = []respond.AppRoutes{

	{"trends.flapping_metrics", "GET", "/{report_name}/flapping/metrics", ListFlappingMetrics},
	{"trends.flapping_endpoints", "GET", "/{report_name}/flapping/endpoints", ListFlappingEndpoints},
	{"trends.flapping_services", "GET", "/{report_name}/flapping/services", ListFlappingServices},
	{"trends.flapping_endpoint_groups", "GET", "/{report_name}/flapping/groups", ListFlappingEndpointGroups},
	{"trends.status_metrics", "GET", "/{report_name}/status/metrics", ListStatusMetrics},
	{"trends.options", "OPTIONS", "/{report_name}/flapping/metrics", Options},
	{"trends.options", "OPTIONS", "/{report_name}/flapping/endpoints", Options},
	{"trends.options", "OPTIONS", "/{report_name}/flapping/services", Options},
	{"trends.options", "OPTIONS", "/{report_name}/flapping/groups", Options},
	{"trends.options", "OPTIONS", "/{report_name}/status/metrics", Options},
}
