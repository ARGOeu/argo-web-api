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
	"fmt"
	"net/http"
	"strings"

	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// HandleSubrouter uses the subrouter for a specific calls and creates a tree of sorts
// handling each route with a different subrouter
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {

	serviceSubrouter := s.StrictSlash(false).PathPrefix("/{report_name}").Subrouter()
	serviceSubrouter = respond.PrepAppRoutes(serviceSubrouter, confhandler, appServiceRoutes)

	groupSubrouter := s.StrictSlash(false).PathPrefix("/{report_name}").Subrouter()
	groupSubrouter = respond.PrepAppRoutes(groupSubrouter, confhandler, appGroupRoutes)

	s = respond.PrepAppRoutes(s, confhandler, appRoutesV2)
}

var appServiceRoutes = []respond.AppRoutes{
	{"results.get", "GET", "/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services/{service_type}", ListServiceFlavorResults},
	{"results.list", "GET", "/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services", ListServiceFlavorResults},
	{"results.get", "GET", "/{lgroup_type}/{lgroup_name}/services/{service_type}", ListServiceFlavorResults},
	{"results.list", "GET", "/{lgroup_type}/{lgroup_name}/services", ListServiceFlavorResults},
}

var appGroupRoutes = []respond.AppRoutes{
	{"results.get", "GET", "/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}", ListEndpointGroupResults},
	{"results.list", "GET", "/{group_type}/{group_name}/{lgroup_type}", ListEndpointGroupResults},
	{"results.get", "GET", "/{group_type}/{group_name}", routeGroup},
	{"results.list", "GET", "/{group_type}", routeGroup},
}

var appRoutesV2 = []respond.AppRoutes{
	{"results.options", "OPTIONS", "/{report_name}/{group_type}", Options},
	{"results.options", "OPTIONS", "/{report_name}/{group_type}/{group_name}/{lgroup_type}", Options},
	{"results.options", "OPTIONS", "/{report_name}/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}", Options},
	{"results.options", "OPTIONS", "/{report_name}/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services", Options},
	{"results.options", "OPTIONS", "/{report_name}/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services/{service_name}", Options},
	{"results.options", "OPTIONS", "/{report_name}/{lgroup_type}/{lgroup_name}", Options},
	{"results.options", "OPTIONS", "/{report_name}/{lgroup_type}/{lgroup_name}/services", Options},
	{"results.options", "OPTIONS", "/{report_name}/{lgroup_type}/{lgroup_name}/services/{service_name}", Options},
}

func routeGroup(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "application/xml"
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Handle response format based on Accept Header
	// Default is application/xml
	format := r.Header.Get("Accept")
	if strings.EqualFold(format, "application/json") {
		contentType = "application/json"
	}

	vars := mux.Vars(r)
	// Grab Tenant DB configuration from context
	tenantcfg := context.Get(r, "tenant_conf").(config.MongoConfig)

	session, err := mongo.OpenSession(tenantcfg)
	defer mongo.CloseSession(session)

	if err != nil {
		return code, h, output, err
	}

	requestedReport := reports.MongoInterface{}
	err = mongo.FindOne(session, tenantcfg.Db, "reports", bson.M{"info.name": vars["report_name"]}, &requestedReport)

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
	output, err = createErrorMessage(message, code, format) //Render the response into XML or JSON
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", format, charset))
	return code, h, output, err

}
