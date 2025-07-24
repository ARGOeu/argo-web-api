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
	respond.PrepAppRoutes(s, confhandler, appRoutesV2)
}

var appRoutesV2 = []respond.AppRoutes{
	{Name: "tenants.user_by_id", Verb: "GET", Path: "/users:byID/{ID}", SubrouterHandler: GetUserByID},
	{Name: "tenants.list", Verb: "GET", Path: "/tenants", SubrouterHandler: List},
	{Name: "tenants.get_status", Verb: "GET", Path: "/tenants/{ID}/status", SubrouterHandler: ListStatus},
	{Name: "tenants.get", Verb: "GET", Path: "/tenants/{ID}", SubrouterHandler: ListOne},
	{Name: "tenants.create", Verb: "POST", Path: "/tenants", SubrouterHandler: Create},
	{Name: "tenants.update_status", Verb: "PUT", Path: "/tenants/{ID}/status", SubrouterHandler: UpdateStatus},
	{Name: "tenants.create_user", Verb: "POST", Path: "/tenants/{ID}/users", SubrouterHandler: CreateUser},
	{Name: "tenants.list_users", Verb: "GET", Path: "/tenants/{ID}/users", SubrouterHandler: ListUsers},
	{Name: "tenants.update_user", Verb: "PUT", Path: "/tenants/{ID}/users/{USER_ID}", SubrouterHandler: UpdateUser},
	{Name: "tenants.delete_user", Verb: "DELETE", Path: "/tenants/{ID}/users/{USER_ID}", SubrouterHandler: DeleteUser},
	{Name: "tenants.get_user", Verb: "GET", Path: "/tenants/{ID}/users/{USER_ID}", SubrouterHandler: GetUser},
	{Name: "tenants.user_refresh_token", Verb: "POST", Path: "/tenants/{ID}/users/{USER_ID}/renew_api_key", SubrouterHandler: RefreshToken},
	{Name: "tenants.update", Verb: "PUT", Path: "/tenants/{ID}", SubrouterHandler: Update},
	{Name: "tenants.delete", Verb: "DELETE", Path: "/tenants/{ID}", SubrouterHandler: Delete},
	{Name: "tenants.options", Verb: "OPTIONS", Path: "/tenants", SubrouterHandler: Options},
	{Name: "tenants.options", Verb: "OPTIONS", Path: "/tenants/{ID}", SubrouterHandler: Options},
	{Name: "tenants.options", Verb: "OPTIONS", Path: "/tenants/{ID}/users", SubrouterHandler: Options},
	{Name: "tenants.options", Verb: "OPTIONS", Path: "/tenants/{ID}/users/{USER_ID}", SubrouterHandler: Options},
	{Name: "tenants.options", Verb: "OPTIONS", Path: "/tenants/{ID}/users/{USER_ID}/renew_api_key", SubrouterHandler: Options},
}
