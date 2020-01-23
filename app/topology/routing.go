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
	"fmt"
	"net/http"

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

	s = respond.PrepAppRoutes(s, confhandler, appRoutesV2)

}

var appRoutesV2 = []respond.AppRoutes{
	{"topology_endpoints.insert", "POST", "/endpoints", CreateEndpoints},
	{"topology_endpoints.list", "GET", "/endpoints", ListEndpoints},
	{"topology_endpoints.delete", "DELETE", "/endpoints", DeleteEndpoints},
	{"topology_endpoints.options", "OPTIONS", "/endpoints", Options},
	{"topology_stats.list", "GET", "/stats/{report_name}", routeCheckGroup},
	{"topology.options", "OPTIONS", "/stats/{report_name}", Options},
}

func routeCheckGroup(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("group check")
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Handle response format based on Accept Header
	contentType := r.Header.Get("Accept")

	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantcfg := context.Get(r, "tenant_conf").(config.MongoConfig)

	session, err := mongo.OpenSession(tenantcfg)
	defer mongo.CloseSession(session)
	if err != nil {
		code = http.StatusInternalServerError
		output, _ = respond.MarshalContent(respond.InternalServerErrorMessage, contentType, "", " ")
		return code, h, output, err
	}
	result := reports.MongoInterface{}
	err = mongo.FindOne(session, tenantcfg.Db, "reports", bson.M{"info.name": vars["report_name"]}, &result)
	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report_name"] + " does not exist"
		output, err := createMessageOUT(message, code, contentType) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}
	// default generic group names
	group_type := "group"
	endpoint_group_type := "endpoint_group"

	// if specific group names are available in report use them
	if result.Topology.Group != nil {
		group_type = result.Topology.Group.Type
		if result.Topology.Group.Group != nil {
			endpoint_group_type = result.Topology.Group.Group.Type
		}
	}

	// set group names as part of the context
	context.Set(r, "group_type", group_type)
	context.Set(r, "endpoint_group_type", endpoint_group_type)
	return ListTopoStats(r, cfg)

}
