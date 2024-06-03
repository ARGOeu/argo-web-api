/*
 * Copyright (c) 2014 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
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
 * or implied, of either GRNET S.A., SRCE or IN2P3 CNRS Computing
 * Centre
 *
 * The work represented by this source file is partially funded by
 * the EGI-InSPIRE project through the European Commission's 7th
 * Framework Programme (contract # INFSO-RI-261323)
 */

package reports

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
	{Name: "reports.list", Verb: "GET", Path: "/reports", SubrouterHandler: List},
	{Name: "reports.get", Verb: "GET", Path: "/reports/{id}", SubrouterHandler: ListOne},
	{Name: "reports.create", Verb: "POST", Path: "/reports", SubrouterHandler: Create},
	{Name: "reports.update", Verb: "PUT", Path: "/reports/{id}", SubrouterHandler: Update},
	{Name: "reports.delete", Verb: "DELETE", Path: "/reports/{id}", SubrouterHandler: Delete},
	{Name: "reports.options", Verb: "OPTIONS", Path: "/reports", SubrouterHandler: Options},
	{Name: "reports.options", Verb: "OPTIONS", Path: "/reports/{id}", SubrouterHandler: Options},
}
