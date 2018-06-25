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

package routing

import (
	"github.com/ARGOeu/argo-web-api/app/aggregationProfiles"
	"github.com/ARGOeu/argo-web-api/app/factors"
	"github.com/ARGOeu/argo-web-api/app/metricProfiles"
	"github.com/ARGOeu/argo-web-api/app/metricResult"
	"github.com/ARGOeu/argo-web-api/app/operationsProfiles"
	"github.com/ARGOeu/argo-web-api/app/recomputations2"
	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/app/results"
	"github.com/ARGOeu/argo-web-api/app/statusEndpointGroups"
	"github.com/ARGOeu/argo-web-api/app/statusEndpoints"
	"github.com/ARGOeu/argo-web-api/app/statusMetrics"
	"github.com/ARGOeu/argo-web-api/app/statusServices"
	"github.com/ARGOeu/argo-web-api/app/tenants"
	"github.com/ARGOeu/argo-web-api/app/thresholdsProfiles"
)

var routesV2 = []RouteV2{
	{"Results", "/results", results.HandleSubrouter},
	{"Metric Result", "/metric_result", metricResult.HandleSubrouter},
	{"Status metric timelines", "/status", statusMetrics.HandleSubrouter},
	{"Status service timelines", "/status", statusServices.HandleSubrouter},
	{"Status endpoint group timelines", "/status", statusEndpointGroups.HandleSubrouter},
	{"Status endpoint timelines", "/status", statusEndpoints.HandleSubrouter},
	{"Recomputations", "", recomputations2.HandleSubrouter},
	{"Metric Profiles", "", metricProfiles.HandleSubrouter},
	{"Reports", "", reports.HandleSubrouter},
	{"Aggregation Profiles", "", aggregationProfiles.HandleSubrouter},
	{"Operations Profiles", "", operationsProfiles.HandleSubrouter},
	{"Thresholds Profiles", "", thresholdsProfiles.HandleSubrouter},
	{"Tenants", "/admin", tenants.HandleSubrouter},
	{"Factors", "", factors.HandleSubrouter},
}
