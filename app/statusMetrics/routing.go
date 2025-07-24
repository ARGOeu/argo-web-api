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

package statusMetrics

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

// HandleSubrouter contains the different paths to follow during subrouting
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {

	respond.PrepAppRoutes(s, confhandler, appRoutesV2)

}

var appRoutesV2 = []respond.AppRoutes{
	{Name: "status.get", Verb: "GET", Path: "/{report_name}/{group_type}/{group_name}/services/{service_name}/endpoints/{endpoint_name}/metrics/{metric_name}", SubrouterHandler: routeCheckGroup},
	{Name: "status.list", Verb: "GET", Path: "/{report_name}/{group_type}/{group_name}/services/{service_name}/endpoints/{endpoint_name}/metrics", SubrouterHandler: routeCheckGroup},
	{Name: "status.options", Verb: "OPTIONS", Path: "/{report_name}/{group_type}/{group_name}/services/{service_name}/endpoints/{endpoint_name}/metrics/{metric_name}", SubrouterHandler: Options},
	{Name: "status.options", Verb: "OPTIONS", Path: "/{report_name}/{group_type}/{group_name}/services/{service_name}/endpoints/{endpoint_name}/metrics", SubrouterHandler: Options},
}

func routeCheckGroup(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Handle response format based on Accept Header
	contentType := r.Header.Get("Accept")

	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantcfg := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	result := reports.MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantcfg.Db).Collection("reports")
	err = rCol.FindOne(context.TODO(), bson.M{"info.name": vars["report_name"]}).Decode(&result)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report_name"] + " does not exist"
		output, err := createMessageOUT(message, code, contentType) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	if vars["group_type"] != result.GetEndpointGroupType() {
		code = http.StatusNotFound
		message := "The report " + vars["report_name"] + " does not define endpoint group type: " + vars["group_type"]
		output, err := createMessageOUT(message, code, contentType) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	return ListMetricTimelines(r, cfg)

}
