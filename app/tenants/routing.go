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

package tenants

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
	{"tenants.user_by_id", "GET", "/users:byID/{ID}", GetUserByID},
	{"tenants.list", "GET", "/tenants", List},
	{"tenants.get_status", "GET", "/tenants/{ID}/status", ListStatus},
	{"tenants.get", "GET", "/tenants/{ID}", ListOne},
	{"tenants.create", "POST", "/tenants", Create},
	{"tenants.update_status", "PUT", "/tenants/{ID}/status", UpdateStatus},
	{"tenants.create_user", "POST", "/tenants/{ID}/users", CreateUser},
	{"tenants.update_user", "PUT", "/tenants/{ID}/users/{USER_ID}", UpdateUser},
	{"tenants.update", "PUT", "/tenants/{ID}", Update},
	{"tenants.delete", "DELETE", "/tenants/{ID}", Delete},
	{"tenants.options", "OPTIONS", "/tenants", Options},
	{"tenants.options", "OPTIONS", "/tenants/{ID}", Options},
	{"tenants.options", "OPTIONS", "/tenants/{ID}/users", Options},
	{"tenants.options", "OPTIONS", "/tenants/{ID}/users/{USER_ID}", Options},
}
