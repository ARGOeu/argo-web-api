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

package routing

import (
	"github.com/argoeu/argo-web-api/app/availabilityProfiles"
	"github.com/argoeu/argo-web-api/app/endpointGroupAvailability"
	"github.com/argoeu/argo-web-api/app/factors"
	"github.com/argoeu/argo-web-api/app/groupGroupsAvailability"
	"github.com/argoeu/argo-web-api/app/jobs"
	"github.com/argoeu/argo-web-api/app/metricProfiles"
	"github.com/argoeu/argo-web-api/app/recomputations"
	"github.com/argoeu/argo-web-api/app/results"
	"github.com/argoeu/argo-web-api/app/serviceFlavorAvailability"
	"github.com/argoeu/argo-web-api/app/statusDetail"
	"github.com/argoeu/argo-web-api/app/statusEndpointGroups"
	"github.com/argoeu/argo-web-api/app/statusEndpoints"
	"github.com/argoeu/argo-web-api/app/statusMsg"
	"github.com/argoeu/argo-web-api/app/statusServices"
	"github.com/argoeu/argo-web-api/app/tenants"
)

var subroutes = SubRouters{
	{"results", "/results", results.HandleSubrouter},
}

var routes = Routes{

	//-----------------------------------Old requests for here on down -------------------------------------------------
	{"group_availability", "GET", "/group_availability", endpointGroupAvailability.List},
	{"group_groups_availability", "GET", "/group_groups_availability", groupGroupsAvailability.List},
	{"endpoint_group_availability", "GET", "/endpoint_group_availability", endpointGroupAvailability.List},
	{"service_flavor_availability", "GET", "/service_flavor_availability", serviceFlavorAvailability.List},
	{"AP List", "GET", "/AP", availabilityProfiles.List},
	{"AP Create", "POST", "/AP", availabilityProfiles.Create},
	{"AP update", "PUT", "/AP/{id}", availabilityProfiles.Update},
	{"AP delete", "DELETE", "/AP/{id}", availabilityProfiles.Delete},
	{"PLACEHOLDER", "GET", "/service_flavor_availability", serviceFlavorAvailability.List},
	{"tenant create", "GET", "/tenants", tenants.Create},
	{"tenant update", "PUT", "/tenants/{name}", tenants.Update},
	{"tenant delete", "DELETE", "/tenants/{name}", tenants.Delete},
	{"tenant list", "GET", "/tenants", tenants.List},
	{"tenant list one", "GET", "/tenants/{name}", tenants.ListOne},

	//jobs
	{"jobs create", "POST", "/jobs", jobs.Create},
	{"job update", "PUT", "/jobs/{name}", jobs.Update},
	{"job delete", "DELETE", "/jobs/{name}", jobs.Delete},
	{"job list", "GET", "/jobs", jobs.List},
	{"job list one", "GET", "/jobs/{name}", jobs.ListOne},

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

	//Status Endpoints
	{"status endpoint list", "GET", "/status/endpoints/timeline/{hostname}/{service_type}", statusEndpoints.List},

	//Status Services
	{"status service list", "GET", "/status/services/timeline/{group}", statusServices.List},

	//Status Sites
	{"status endpoint group list", "GET", "/status/sites/timeline/{group}", statusEndpointGroups.List},
}
