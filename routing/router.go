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
 * or implied, of either GRNET S.A., SRCE or IN2P3 CNRS Computing
 * Centre
 *
 * The work represented by this source file is partially funded by
 * the EGI-InSPIRE project through the European Commission's 7th
 * Framework Programme (contract # INFSO-RI-261323)
 */

package routing

import (
	"fmt"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/gorilla/mux"
)

// RouteV2 represents the new style of routes that are handled by each package respectively
type RouteV2 struct {
	Name             string
	Pattern          string
	SubrouterHandler func(*mux.Router, *respond.ConfHandler)
}

// Same Route model as V2
type RouteV3 RouteV2

// NewRouter creates the main router that will be used by the api (contains both v2 and v3 routes)
func NewRouter(cfg config.Config) *mux.Router {

	confhandler := respond.ConfHandler{Config: cfg}

	router := mux.NewRouter().StrictSlash(false)

	// Add v2 subroutes
	for _, subroute := range routesV2 {
		subrouter := router.
			PathPrefix("/api/v2" + subroute.Pattern).
			Name(subroute.Name).
			Subrouter()
		subroute.SubrouterHandler(subrouter, &confhandler)
	}

	// Add v3 subroutes
	for _, subroute := range routesV3 {
		subrouter := router.
			PathPrefix("/api/v3" + subroute.Pattern).
			Name(subroute.Name).
			Subrouter()
		subroute.SubrouterHandler(subrouter, &confhandler)
	}
	return router
}

// PrintRoutes Attempts to print all register routes when called using mux.Router.Walk(PrintRoutes)
func PrintRoutes(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	url, _ := route.URL()
	fmt.Println(route.GetName(), url)
	return nil
}
