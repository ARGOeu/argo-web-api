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
	"net/http"

	"github.com/argoeu/argo-web-api/respond"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/gorilla/mux"
)

// HandleSubrouter contains the different paths to follow during subrouting
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {

	// Goes up to /report/REPORT_NAME/group_type
	groupSubrouter := s.PathPrefix("/{report_name}/{group_type}").Subrouter()

	// eg. timelines/critical/SITE/mysite/service/apache/endpoints/apache01.host/metrics/memory_used
	groupSubrouter.
		Path("/{group_name}/services/{service_name}/endpoints/{endpoint_name}/metrics/{metric_name}").
		Methods("GET").
		Name("metric name").
		Handler(confhandler.Respond(ListMetricTimeline))

}

func routeGroup(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)

	return code, h, output, err

}
