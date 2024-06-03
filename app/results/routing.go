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

package results

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
)

// HandleSubrouter uses the subrouter for a specific calls and creates a tree of sorts
// handling each route with a different subrouter
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {

	// Routes for serving endpoint a/r results under services
	endpointSubrouter := s.StrictSlash(false).PathPrefix("/{report_name}").Subrouter()
	respond.PrepAppRoutes(endpointSubrouter, confhandler, appEndpointRoutes)

	serviceSubrouter := s.StrictSlash(false).PathPrefix("/{report_name}").Subrouter()
	respond.PrepAppRoutes(serviceSubrouter, confhandler, appServiceRoutes)

	groupSubrouter := s.StrictSlash(false).PathPrefix("/{report_name}").Subrouter()
	respond.PrepAppRoutes(groupSubrouter, confhandler, appGroupRoutes)

	respond.PrepAppRoutes(s, confhandler, appRoutesV2)
}

var appEndpointRoutes = []respond.AppRoutes{

	{Name: "results.get", Verb: "GET", Path: "/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services/{service_type}/endpoints/{endpoint_name}", SubrouterHandler: ListEndpointResults},
	{Name: "results.get", Verb: "GET", Path: "/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services/{service_type}/endpoints", SubrouterHandler: ListEndpointResults},
	{Name: "results.get", Verb: "GET", Path: "/{lgroup_type}/{lgroup_name}/services/{service_type}/endpoints/{endpoint_name}", SubrouterHandler: ListEndpointResults},
	{Name: "results.get", Verb: "GET", Path: "/{lgroup_type}/{lgroup_name}/services/{service_type}/endpoints", SubrouterHandler: ListEndpointResults},
	//routes to quickly get endpoints included in an endpoint group without specifing service
	{Name: "results.get", Verb: "GET", Path: "/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/endpoints/{endpoint_name}", SubrouterHandler: ListEndpointResults},
	{Name: "results.get", Verb: "GET", Path: "/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/endpoints", SubrouterHandler: ListEndpointResults},
	{Name: "results.get", Verb: "GET", Path: "/{lgroup_type}/{lgroup_name}/endpoints/{endpoint_name}", SubrouterHandler: ListEndpointResults},
	{Name: "results.get", Verb: "GET", Path: "/{lgroup_type}/{lgroup_name}/endpoints", SubrouterHandler: ListEndpointResults},
	{Name: "results.get", Verb: "GET", Path: "/endpoints", SubrouterHandler: FlatListEndpointResults},
	// normal routes to get endpoints included in a service

}

var appServiceRoutes = []respond.AppRoutes{

	{Name: "results.get", Verb: "GET", Path: "/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services/{service_type}", SubrouterHandler: ListServiceFlavorResults},
	{Name: "results.list", Verb: "GET", Path: "/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services", SubrouterHandler: ListServiceFlavorResults},
	{Name: "results.get", Verb: "GET", Path: "/{lgroup_type}/{lgroup_name}/services/{service_type}", SubrouterHandler: ListServiceFlavorResults},
	{Name: "results.list", Verb: "GET", Path: "/{lgroup_type}/{lgroup_name}/services", SubrouterHandler: ListServiceFlavorResults},
}

var appGroupRoutes = []respond.AppRoutes{
	{Name: "results.get", Verb: "GET", Path: "/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}", SubrouterHandler: ListEndpointGroupResults},
	{Name: "results.list", Verb: "GET", Path: "/{group_type}/{group_name}/{lgroup_type}", SubrouterHandler: ListEndpointGroupResults},
	{Name: "results.get", Verb: "GET", Path: "/{group_type}/{group_name}", SubrouterHandler: routeGroup},
	{Name: "results.list", Verb: "GET", Path: "/{group_type}", SubrouterHandler: routeGroup},
}

var appRoutesV2 = []respond.AppRoutes{
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/endpoints", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{group_type}", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{group_type}/{group_name}/{lgroup_type}", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/endpoints", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/endpoints/{endpoint_name}", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services/{service_name}", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services/{service_name}/endpoints", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services/{service_name}/endpoints/{endpoint_name}", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{lgroup_type}/{lgroup_name}", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{lgroup_type}/{lgroup_name}/endpoints", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{lgroup_type}/{lgroup_name}/endpoints/{endpoint_name}", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{lgroup_type}/{lgroup_name}/services", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{lgroup_type}/{lgroup_name}/services/{service_name}", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{lgroup_type}/{lgroup_name}/services/{service_name}/endpoints", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report_name}/{lgroup_type}/{lgroup_name}/services/{service_name}/endpoints/{endpoint_name}", SubrouterHandler: Options},
}

func routeGroup(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	requestedReport := reports.MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantcfg.Db).Collection("reports")
	err = rCol.FindOne(context.TODO(), bson.M{"info.name": vars["report_name"]}).Decode(&requestedReport)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report_name"] + " does not exist"
		output, err := createErrorMessage(message, code, contentType) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	selectedGroupType := requestedReport.DetermineGroupType(vars["group_type"])

	if selectedGroupType == "endpoint" {
		if vars["lgroup_type"] == "" {
			vars["lgroup_type"] = vars["group_type"]
			vars["lgroup_name"] = vars["group_name"]
			vars["group_type"] = ""
			vars["group_name"] = ""
		}
		return ListEndpointGroupResults(r, cfg)
	} else if selectedGroupType == "group" {
		return ListSuperGroupResults(r, cfg)
	}

	code = http.StatusNotFound
	message := "The report " + vars["report_name"] + " does not define any group type: " + vars["group_type"]
	output, err = createErrorMessage(message, code, contentType) //Render the response into XML or JSON
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err

}
