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
	"context"
	"fmt"
	"net/http"

	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// HandleSubrouter uses the subrouter for a specific calls and creates a tree of sorts
// handling each route with a different subrouter
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {

	respond.PrepAppRoutes(s, confhandler, appRoutesV2)

}

var appRoutesV2 = []respond.AppRoutes{
	{Name: "topology_endpoints_report.list", Verb: "GET", Path: "/endpoints/by_report/{report}", SubrouterHandler: ListEndpointsByReport},
	{Name: "topology_groups_report.list", Verb: "GET", Path: "/groups/by_report/{report}", SubrouterHandler: ListGroupsByReport},
	{Name: "topology_groups.delete", Verb: "DELETE", Path: "/groups", SubrouterHandler: DeleteGroups},
	{Name: "topology_groups.list", Verb: "GET", Path: "/groups", SubrouterHandler: ListGroups},
	{Name: "topology_groups.insert", Verb: "POST", Path: "/groups", SubrouterHandler: CreateGroups},
	{Name: "topology_groups.options", Verb: "OPTIONS", Path: "/groups", SubrouterHandler: Options},
	{Name: "topology_endpoints.insert", Verb: "POST", Path: "/endpoints", SubrouterHandler: CreateEndpoints},
	{Name: "topology_endpoints.list", Verb: "GET", Path: "/endpoints", SubrouterHandler: ListEndpoints},
	{Name: "topology_endpoints.delete", Verb: "DELETE", Path: "/endpoints", SubrouterHandler: DeleteEndpoints},
	{Name: "topology_endpoints.options", Verb: "OPTIONS", Path: "/endpoints", SubrouterHandler: Options},
	{Name: "topology_service_types.insert", Verb: "POST", Path: "/service-types", SubrouterHandler: CreateServiceTypes},
	{Name: "topology_service_types.list", Verb: "GET", Path: "/service-types", SubrouterHandler: ListServiceTypes},
	{Name: "topology_service_types.delete", Verb: "DELETE", Path: "/service-types", SubrouterHandler: DeleteServiceTypes},
	{Name: "topology_service_types.options", Verb: "OPTIONS", Path: "/service-types", SubrouterHandler: Options},
	{Name: "topology_tags.list", Verb: "GET", Path: "/tags", SubrouterHandler: ListTopoTags},
	{Name: "topology_tags.options", Verb: "OPTIONS", Path: "/tags", SubrouterHandler: Options},
	{Name: "topology_stats.list", Verb: "GET", Path: "/stats/{report_name}", SubrouterHandler: routeCheckGroup},
	{Name: "topology.options", Verb: "OPTIONS", Path: "/stats/{report_name}", SubrouterHandler: Options},
}

func routeCheckGroup(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Handle response format based on Accept Header
	contentType := r.Header.Get("Accept")

	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantcfg := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	result := reports.MongoInterface{}
	reportsCol := cfg.MongoClient.Database(tenantcfg.Db).Collection("reports")
	err = reportsCol.FindOne(context.TODO(), bson.M{"info.name": vars["report_name"]}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			code = http.StatusNotFound
			message := "The report with the name " + vars["report_name"] + " does not exist"
			output, err := createMessageOUT(message, code, contentType) //Render the response into XML or JSON
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err

	}
	// default generic group names
	groupType := "group"
	endpointGroupType := "endpoint_group"

	// if specific group names are available in report use them
	if result.Topology.Group != nil {
		groupType = result.Topology.Group.Type
		if result.Topology.Group.Group != nil {
			endpointGroupType = result.Topology.Group.Group.Type
		}
	}

	// set group names as part of the context
	gcontext.Set(r, "group_type", groupType)
	gcontext.Set(r, "endpoint_group_type", endpointGroupType)
	return ListTopoStats(r, cfg)

}
