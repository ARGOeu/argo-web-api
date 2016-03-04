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

	"github.com/ARGOeu/argo-web-api/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/ARGOeu/argo-web-api/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/authentication"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
)

// HandleSubrouter uses the subrouter for a specific calls and creates a tree of sorts
// handling each route with a different subrouter
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {

	serviceSubrouter := s.PathPrefix("/{report_name}").Subrouter()

	serviceSubrouter.Path("/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services/{service_type}").
		Methods("GET").
		Name("Service Flavor").
		Handler(confhandler.Respond(ListServiceFlavorResults))

	serviceSubrouter.Path("/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services").
		Methods("GET").
		Name("Service Flavor").
		Handler(confhandler.Respond(ListServiceFlavorResults))

	serviceSubrouter.Path("/{lgroup_type}/{lgroup_name}/services/{service_type}").
		Methods("GET").
		Name("Service Flavor").
		Handler(confhandler.Respond(ListServiceFlavorResults))

	serviceSubrouter.Path("/{lgroup_type}/{lgroup_name}/services").
		Methods("GET").
		Name("Service Flavor").
		Handler(confhandler.Respond(ListServiceFlavorResults))

	groupSubrouter := s.PathPrefix("/{report_name}/{group_type}").Subrouter()
	groupSubrouter.
		Path("/{group_name}/{lgroup_type}/{lgroup_name}").
		Methods("GET").
		Name("Group name").
		Handler(confhandler.Respond(ListEndpointGroupResults))
	groupSubrouter.
		Path("/{group_name}/{lgroup_type}").
		Methods("GET").
		Name("Group name").
		Handler(confhandler.Respond(ListEndpointGroupResults))
	groupSubrouter.
		Path("/{group_name}").
		Methods("GET").
		Name("Group name").
		Handler(confhandler.Respond(routeGroup))
	groupSubrouter.
		Methods("GET").
		Name("Group Type").
		Handler(confhandler.Respond(routeGroup))

	s.Methods("OPTIONS").
		Path("/{report_name}/{group_type}").
		Name("List Options of Resource").
		Handler(confhandler.Respond(Options))

	s.Methods("OPTIONS").
		Path("/{report_name}/{group_type}/{group_name}/{lgroup_type}").
		Name("List Options of Resource").
		Handler(confhandler.Respond(Options))

	s.Methods("OPTIONS").
		Path("/{report_name}/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}").
		Name("List Options of Resource").
		Handler(confhandler.Respond(Options))

	s.Methods("OPTIONS").
		Path("/{report_name}/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services").
		Name("List Options of Resource").
		Handler(confhandler.Respond(Options))

	s.Methods("OPTIONS").
		Path("/{report_name}/{group_type}/{group_name}/{lgroup_type}/{lgroup_name}/services/{service_name}").
		Name("List Options of Resource").
		Handler(confhandler.Respond(Options))

	s.Methods("OPTIONS").
		Path("/{report_name}/{lgroup_type}/{lgroup_name}").
		Name("List Options of Resource").
		Handler(confhandler.Respond(Options))

	s.Methods("OPTIONS").
		Path("/{report_name}/{lgroup_type}/{lgroup_name}/services").
		Name("List Options of Resource").
		Handler(confhandler.Respond(Options))

	s.Methods("OPTIONS").
		Path("/{report_name}/{lgroup_type}/{lgroup_name}/services/{service_name}").
		Name("List Options of Resource").
		Handler(confhandler.Respond(Options))

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
	tenantcfg, err := authentication.AuthenticateTenant(r.Header, cfg)
	if err != nil {
		if err.Error() == "Unauthorized" {
			code = http.StatusUnauthorized
			message := err.Error()
			output, err = createErrorMessage(message, code, contentType)
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}
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
