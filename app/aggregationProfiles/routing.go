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

package aggregationProfiles

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
	{"aggregationProfiles.list", "GET", "/aggregation_profiles", List},
	{"aggregationProfiles.get", "GET", "/aggregation_profiles/{ID}", ListOne},
	{"aggregationProfiles.create", "POST", "/aggregation_profiles", Create},
	{"aggregationProfiles.update", "PUT", "/aggregation_profiles/{ID}", Update},
	{"aggregationProfiles.delete", "DELETE", "/aggregation_profiles/{ID}", Delete},
	{"aggregationProfiles.options", "OPTIONS", "/aggregation_profiles", Options},
	{"aggregationProfiles.options", "OPTIONS", "/aggregation_profiles/{ID}", Options},
}
