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

	respond.PrepAppRoutes(s, confhandler, appRoutesV2)
}

var appRoutesV2 = []respond.AppRoutes{

	{Name: "trends.flapping_metrics_tags", Verb: "GET", Path: "/{report_name}/flapping/metrics/tags", SubrouterHandler: ListFlappingMetricsTags},
	{Name: "trends.flapping_metrics", Verb: "GET", Path: "/{report_name}/flapping/metrics", SubrouterHandler: ListFlappingMetrics},
	{Name: "trends.flapping_endpoints", Verb: "GET", Path: "/{report_name}/flapping/endpoints", SubrouterHandler: ListFlappingEndpoints},
	{Name: "trends.flapping_services", Verb: "GET", Path: "/{report_name}/flapping/services", SubrouterHandler: ListFlappingServices},
	{Name: `trends.flapping_endpoint_groups`, Verb: "GET", Path: "/{report_name}/flapping/groups", SubrouterHandler: ListFlappingEndpointGroups},
	{Name: "trends.status_metrics_tags", Verb: "GET", Path: "/{report_name}/status/metrics/tags", SubrouterHandler: ListStatusMetricsTags},
	{Name: "trends.status_metrics", Verb: "GET", Path: "/{report_name}/status/metrics", SubrouterHandler: ListStatusMetrics},
	{Name: "trends.status_endpoints", Verb: "GET", Path: "/{report_name}/status/endpoints", SubrouterHandler: ListStatusEndpoints},
	{Name: "trends.status_services", Verb: "GET", Path: "/{report_name}/status/services", SubrouterHandler: ListStatusServices},
	{Name: "trends.status_endpoint_groups", Verb: "GET", Path: "/{report_name}/status/groups", SubrouterHandler: ListStatusEgroups},
	{Name: "trends.options", Verb: "OPTIONS", Path: "/{report_name}/flapping/metrics", SubrouterHandler: Options},
	{Name: "trends.options", Verb: "OPTIONS", Path: "/{report_name}/flapping/endpoints", SubrouterHandler: Options},
	{Name: "trends.options", Verb: "OPTIONS", Path: "/{report_name}/flapping/services", SubrouterHandler: Options},
	{Name: "trends.options", Verb: "OPTIONS", Path: "/{report_name}/flapping/groups", SubrouterHandler: Options},
	{Name: "trends.options", Verb: "OPTIONS", Path: "/{report_name}/status/metrics/tags", SubrouterHandler: Options},
	{Name: "trends.options", Verb: "OPTIONS", Path: "/{report_name}/status/metrics", SubrouterHandler: Options},
	{Name: "trends.options", Verb: "OPTIONS", Path: "/{report_name}/status/endpoints", SubrouterHandler: Options},
	{Name: "trends.options", Verb: "OPTIONS", Path: "/{report_name}/status/services", SubrouterHandler: Options},
	{Name: "trends.options", Verb: "OPTIONS", Path: "/{report_name}/status/groups", SubrouterHandler: Options},
}
