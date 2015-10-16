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

package statusServices

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/authentication"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// HandleSubrouter contains the different paths to follow during subrouting
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {

	// Goes up to /report/REPORT_NAME/group_type
	groupSubrouter := s.PathPrefix("/{report_name}/{group_type}").Subrouter()

	// eg. timelines/critical/SITE/mysite/services/apache
	groupSubrouter.
		Path("/{group_name}/services/{service_name}").
		Methods("GET").
		Name("service name").
		Handler(confhandler.Respond(routeCheckGroup))

	// eg. timelines/critical/SITE/mysite/services
	groupSubrouter.
		Path("/{group_name}/services/").
		Methods("GET").
		Name("all services").
		Handler(confhandler.Respond(routeCheckGroup))

}

func routeCheckGroup(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("group check")
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
		return code, h, output, err
	}
	session, err := mongo.OpenSession(tenantcfg)
	defer mongo.CloseSession(session)
	if err != nil {
		return code, h, output, err
	}
	result := reports.MongoInterface{}
	err = mongo.FindOne(session, tenantcfg.Db, "reports", bson.M{"info.name": vars["report_name"]}, &result)

	if err != nil {
		message := "The report with the name " + vars["report_name"] + " does not exist"
		output, err := createMessageOUT(message, format) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	if vars["group_type"] != result.GetEndpointGroupType() {
		message := "The report " + vars["report_name"] + " does not define endpoint group type: " + vars["group_type"]
		output, err := createMessageOUT(message, format) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	return ListServiceTimelines(r, cfg)

}
