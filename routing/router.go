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
	"net/http"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"

	"github.com/gorilla/mux"
)

// Route represents the old style routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(*http.Request, config.Config) (int, http.Header, []byte, error)
}

// Subrouter represents the new style of routes that are handled by each package respectively
type SubRouter struct {
	Name             string
	Pattern          string
	SubrouterHandler func(*mux.Router, *respond.ConfHandler)
}

// NewRouter creates the main router that will be used by the api
func NewRouter(cfg config.Config) *mux.Router {

	confhandler := respond.ConfHandler{Config: cfg}
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		handler := confhandler.Respond(route.HandlerFunc)
		router.
			PathPrefix("/api/v1").
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	for _, subroute := range subroutes {
		subrouter := router.
			PathPrefix("/api/v2").
			PathPrefix(subroute.Pattern).
			Subrouter()
		subroute.SubrouterHandler(subrouter, &confhandler)
	}
	// router.Walk(PrintRoutes)
	return router
}

// PrintRoutes Attempts to print all register routes when called using mux.Router.Walk(PrintRoutes)
func PrintRoutes(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	url, _ := route.URL()
	fmt.Println(route.GetName(), url)
	return nil
}
