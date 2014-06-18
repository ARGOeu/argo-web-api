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

package main

import (
	"github.com/argoeu/ar-web-api/app/availabilityProfiles"
	"github.com/argoeu/ar-web-api/app/ngiAvailability"
	"github.com/argoeu/ar-web-api/app/poemProfiles"
	"github.com/argoeu/ar-web-api/app/recomputations"
	"github.com/argoeu/ar-web-api/app/serviceFlavorAvailability"
	"github.com/argoeu/ar-web-api/app/siteAvailability"
	"github.com/argoeu/ar-web-api/app/voAvailability"
	"github.com/argoeu/ar-web-api/app/factors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func main() {

	//Create the server router
	mainRouter := mux.NewRouter()
	//SUBROUTER DEFINITIONS
	getSubrouter := mainRouter.Methods("GET").Subrouter() //Routes only GET requests
	postSubrouter := mainRouter.Methods("POST").Headers("x-api-key", "").Subrouter() //Routes only POST requests
	deleteSubrouter := mainRouter.Methods("DELETE").Headers("x-api-key", "").Subrouter() //Routes only DELETE requests
	putSubrouter := mainRouter.Methods("PUT").Headers("x-api-key", "").Subrouter() //Routes only PUT requests
	//All requests that modify data must provide with authentication credentials

	// Grouping calls.
	// Groups are routed depending on the value of the parameter group type.
	// 2) Provide with a default call informing the user of an invalid parameter
	getSubrouter.HandleFunc("/api/v1/group_availability", Respond(voAvailability.List)).
		Queries("group_type", "vo")
	getSubrouter.HandleFunc("/api/v1/group_availability", Respond(siteAvailability.List)).
		Queries("group_type", "site")
	getSubrouter.HandleFunc("/api/v1/group_availability", Respond(ngiAvailability.List)).
		Queries("group_type", "ngi")

	// Service Flavor Availability
	getSubrouter.HandleFunc("/api/v1/service_flavor_availability", Respond(serviceFlavorAvailability.List))

	//Availability Profiles
	postSubrouter.HandleFunc("/api/v1/AP", Respond(availabilityProfiles.Create))
	getSubrouter.HandleFunc("/api/v1/AP", Respond(availabilityProfiles.List))
	putSubrouter.HandleFunc("/api/v1/AP/{id}", Respond(availabilityProfiles.Update))
	deleteSubrouter.HandleFunc("/api/v1/AP/{id}", Respond(availabilityProfiles.Delete))

	//POEM Profiles
	getSubrouter.HandleFunc("/api/v1/poems", Respond(poemProfiles.List))

	//Recalculations
	postSubrouter.HandleFunc("/api/v1/recomputations", Respond(recomputations.Create))
	getSubrouter.HandleFunc("/api/v1/recomputations", Respond(recomputations.List))
	
	getSubrouter.HandleFunc("/api/v1/factors", Respond(factors.List))
	

	http.Handle("/", mainRouter)

	//Cache
	//get_subrouter.HandleFunc("/api/v1/reset_cache", Respond("text/xml", "utf-8", ResetCache))

	//Web service binds to server. Requests served over HTTPS.
	err := http.ListenAndServeTLS(cfg.Server.Bindip+":"+strconv.Itoa(cfg.Server.Port), "/etc/pki/tls/certs/localhost.crt", "/etc/pki/tls/private/localhost.key", nil)

	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
