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

package operationsProfiles

import (
	"github.com/ARGOeu/argo-web-api/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/ARGOeu/argo-web-api/respond"
)

// HandleSubrouter uses the subrouter for a specific calls and creates a tree of sorts
// handling each route with a different subrouter
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {
	s.Methods("GET").
		Path("/operations_profiles").
		Name("List Operations Profiles").
		Handler(confhandler.Respond(List))

	s.Methods("GET").
		Path("/operations_profiles/{ID}").
		Name("List One Operations Profile").
		Handler(confhandler.Respond(ListOne))

	s.Methods("POST").
		Path("/operations_profiles").
		Name("Create Operations Profile").
		Handler(confhandler.Respond(Create))

	s.Methods("PUT").
		Path("/operations_profiles/{ID}").
		Name("Update Operations Profile").
		Handler(confhandler.Respond(Update))

	s.Methods("DELETE").
		Path("/operations_profiles/{ID}").
		Name("Delete Operations Profile").
		Handler(confhandler.Respond(Delete))
}
