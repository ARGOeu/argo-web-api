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
	"github.com/ARGOeu/argo-web-api/app/availabilityProfiles"
	"github.com/ARGOeu/argo-web-api/app/endpointGroupAvailability"
	"github.com/ARGOeu/argo-web-api/app/factors"
	"github.com/ARGOeu/argo-web-api/app/groupGroupsAvailability"
	"github.com/ARGOeu/argo-web-api/app/metricProfiles"
	"github.com/ARGOeu/argo-web-api/app/metric_result"
	"github.com/ARGOeu/argo-web-api/app/recomputations"
	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/app/results"
	"github.com/ARGOeu/argo-web-api/app/serviceFlavorAvailability"
	"github.com/ARGOeu/argo-web-api/app/statusDetail"
	"github.com/ARGOeu/argo-web-api/app/statusEndpointGroups"
	"github.com/ARGOeu/argo-web-api/app/statusEndpoints"
	"github.com/ARGOeu/argo-web-api/app/statusMetrics"
	"github.com/ARGOeu/argo-web-api/app/statusMsg"
	"github.com/ARGOeu/argo-web-api/app/statusServices"
	"github.com/ARGOeu/argo-web-api/app/tenants"
)

var subroutes = []SubRouter{
	{"Results", "/results", results.HandleSubrouter},
	{"Metric Result", "/metric_result", metric_result.HandleSubrouter},
	{"Status metric timelines", "/status", statusMetrics.HandleSubrouter},
	{"Status service timelines", "/status", statusServices.HandleSubrouter},
	{"Status endpoint group timelines", "/status", statusEndpointGroups.HandleSubrouter},
	{"Status endpoint timelines", "/status", statusEndpoints.HandleSubrouter},
}

var routes = []Route{

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

	//Poem Profiles compatibility
	{"List poems", "GET", "/poems", metricProfiles.ListPoems},

	//Metric Profiles
	{"list metric profile", "GET", "/metric_profiles", metricProfiles.List},
	{"metric profile create", "POST", "/metric_profiles", metricProfiles.Create},
	{"metric profile delete", "DELETE", "/metric_profiles/{id}", metricProfiles.Delete},
	{"metric profile update", "PUT", "/metric_profiles/{id}", metricProfiles.Update},

	//Recalculations
	{"recomputation create", "POST", "/recomputations", recomputations.Create},
	{"recomputation list", "GET", "/recomputations", recomputations.List},

	{"factors list", "GET", "/factors", factors.List},

	//Status
	{"status detail list", "GET", "/status/metrics/timeline/{group}", statusDetail.List},

	//Status Raw Msg
	{"status message list", "GET", "/status/metrics/msg/{hostname}/{service}/{metric}", statusMsg.List},
}
