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
	"github.com/ARGOeu/argo-web-api/app/availabilityProfiles"
	"github.com/ARGOeu/argo-web-api/app/endpointGroupAvailability"
	"github.com/ARGOeu/argo-web-api/app/factors"
	"github.com/ARGOeu/argo-web-api/app/groupGroupsAvailability"
	"github.com/ARGOeu/argo-web-api/app/metricProfiles"
	"github.com/ARGOeu/argo-web-api/app/metricResult"
	"github.com/ARGOeu/argo-web-api/app/operationsProfiles"
	"github.com/ARGOeu/argo-web-api/app/recomputations"
	"github.com/ARGOeu/argo-web-api/app/recomputations2"
	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/app/reportsv2"
	"github.com/ARGOeu/argo-web-api/app/results"
	"github.com/ARGOeu/argo-web-api/app/serviceFlavorAvailability"
	"github.com/ARGOeu/argo-web-api/app/statusEndpointGroups"
	"github.com/ARGOeu/argo-web-api/app/statusEndpoints"
	"github.com/ARGOeu/argo-web-api/app/statusMetrics"
	"github.com/ARGOeu/argo-web-api/app/statusServices"
	"github.com/ARGOeu/argo-web-api/app/tenants"
)

var routesV2 = []RouteV2{
	{"Results", "/results", results.HandleSubrouter},
	{"Metric Result", "/metric_result", metricResult.HandleSubrouter},
	{"Status metric timelines", "/status", statusMetrics.HandleSubrouter},
	{"Status service timelines", "/status", statusServices.HandleSubrouter},
	{"Status endpoint group timelines", "/status", statusEndpointGroups.HandleSubrouter},
	{"Status endpoint timelines", "/status", statusEndpoints.HandleSubrouter},
	{"Recomputations", "", recomputations2.HandleSubrouter},
	{"Reports", "/reports", reportsv2.HandleSubrouter},
	{"Metric Profiles", "", metricProfiles.HandleSubrouter},
	{"Aggregation Profiles", "", aggregationProfiles.HandleSubrouter},
	{"Operations Profiles", "", operationsProfiles.HandleSubrouter},
}

var routesV1 = []RouteV1{

	//-----------------------------------Old requests for here on down -------------------------------------------------
	{"Group Availability", "GET", "/group_availability", endpointGroupAvailability.List},
	{"Group Groups Availability", "GET", "/group_groups_availability", groupGroupsAvailability.List},
	{"Endpoint Group Availability", "GET", "/endpoint_group_availability", endpointGroupAvailability.List},
	{"Service flavor Availability", "GET", "/service_flavor_availability", serviceFlavorAvailability.List},
	{"AP List", "GET", "/AP", availabilityProfiles.List},
	{"AP Create", "POST", "/AP", availabilityProfiles.Create},
	{"AP update", "PUT", "/AP/{id}", availabilityProfiles.Update},
	{"AP delete", "DELETE", "/AP/{id}", availabilityProfiles.Delete},
	{"Service Falvor Availability", "GET", "/service_flavor_availability", serviceFlavorAvailability.List},
	{"tenant create", "POST", "/tenants", tenants.Create},
	{"tenant update", "PUT", "/tenants/{name}", tenants.Update},
	{"tenant delete", "DELETE", "/tenants/{name}", tenants.Delete},
	{"tenant list", "GET", "/tenants", tenants.List},
	{"tenant list one", "GET", "/tenants/{name}", tenants.ListOne},

	//reports
	{"reports create", "POST", "/reports", reports.Create},
	{"reports update", "PUT", "/reports/{name}", reports.Update},
	{"reports delete", "DELETE", "/reports/{name}", reports.Delete},
	{"reports list", "GET", "/reports", reports.List},
	{"reports list one", "GET", "/reports/{name}", reports.ListOne},

	//Recalculations
	{"recomputation create", "POST", "/recomputations", recomputations.Create},
	{"recomputation list", "GET", "/recomputations", recomputations.List},

	{"factors list", "GET", "/factors", factors.List},
}
