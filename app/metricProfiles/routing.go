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

package metricProfiles

import (
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/gorilla/mux"
)

// HandleSubrouter uses the subrouter for a specific calls and creates a tree of sorts
// handling each route with a different subrouter
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {

	s = respond.PrepAppRoutes(s, confhandler, appRoutesV2)

}

var appRoutesV2 = []respond.AppRoutes{
	{"metricProfiles.list", "GET", "/metric_profiles", List},
	{"metricProfiles.get", "GET", "/metric_profiles/{ID}", ListOne},
	{"metricProfiles.create", "POST", "/metric_profiles", Create},
	{"metricProfiles.update", "PUT", "/metric_profiles/{ID}", Update},
	{"metricProfiles.delete", "DELETE", "/metric_profiles/{ID}", Delete},
	{"metricProfiles.options", "OPTIONS", "/metric_profiles", Options},
	{"metricProfiles.options", "OPTIONS", "/metric_profiles/{ID}", Options},
}
