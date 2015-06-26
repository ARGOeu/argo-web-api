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
	"crypto/tls"
	"log"
	"net/http"
	"strconv"

	"github.com/argoeu/argo-web-api/app/availabilityProfiles"
	"github.com/argoeu/argo-web-api/app/endpointGroupAvailability"
	"github.com/argoeu/argo-web-api/app/factors"
	"github.com/argoeu/argo-web-api/app/groupGroupsAvailability"
	"github.com/argoeu/argo-web-api/app/jobs"
	"github.com/argoeu/argo-web-api/app/metricProfiles"
	"github.com/argoeu/argo-web-api/app/recomputations"
	"github.com/argoeu/argo-web-api/app/serviceFlavorAvailability"
	"github.com/argoeu/argo-web-api/app/statusDetail"
	"github.com/argoeu/argo-web-api/app/statusEndpointGroups"
	"github.com/argoeu/argo-web-api/app/statusEndpoints"
	"github.com/argoeu/argo-web-api/app/statusMsg"
	"github.com/argoeu/argo-web-api/app/statusServices"
	"github.com/argoeu/argo-web-api/app/tenants"
	"github.com/argoeu/argo-web-api/app/voAvailability"
	"github.com/gorilla/mux"
)

func main() {

	//Create the server router
	mainRouter := mux.NewRouter()
	//SUBROUTER DEFINITIONS
	getSubrouter := mainRouter.Methods("GET").Subrouter()                                //Routes GET requests
	authGetSubrouter := mainRouter.Methods("GET").Headers("x-api-key", "").Subrouter()   //Routes GET requests with auth
	postSubrouter := mainRouter.Methods("POST").Headers("x-api-key", "").Subrouter()     //Routes only POST requests
	deleteSubrouter := mainRouter.Methods("DELETE").Headers("x-api-key", "").Subrouter() //Routes only DELETE requests
	putSubrouter := mainRouter.Methods("PUT").Headers("x-api-key", "").Subrouter()       //Routes only PUT requests
	//All requests that modify data must provide with authentication credentials

	// Grouping calls.
	// Groups are routed depending on the value of the parameter group type.
	// 2) Provide with a default call informing the user of an invalid parameter
	getSubrouter.HandleFunc("/api/v1/group_availability", Respond(voAvailability.List)).
		Queries("group_type", "vo")

	// Group of Groups availability
	getSubrouter.HandleFunc("/api/v1/group_groups_availability", Respond(groupGroupsAvailability.List))
	getSubrouter.HandleFunc("/api/v1/endpoint_group_availability", Respond(endpointGroupAvailability.List))

	// Service Flavor Availability
	getSubrouter.HandleFunc("/api/v1/service_flavor_availability", Respond(serviceFlavorAvailability.List))

	//Availability Profiles
	postSubrouter.HandleFunc("/api/v1/AP", Respond(availabilityProfiles.Create))
	getSubrouter.HandleFunc("/api/v1/AP", Respond(availabilityProfiles.List))
	putSubrouter.HandleFunc("/api/v1/AP/{id}", Respond(availabilityProfiles.Update))
	deleteSubrouter.HandleFunc("/api/v1/AP/{id}", Respond(availabilityProfiles.Delete))

	//tenants
	postSubrouter.HandleFunc("/api/v1/tenants", Respond(tenants.Create))
	putSubrouter.HandleFunc("/api/v1/tenants/{name}", Respond(tenants.Update))
	deleteSubrouter.HandleFunc("/api/v1/tenants/{name}", Respond(tenants.Delete))
	getSubrouter.HandleFunc("/api/v1/tenants", Respond(tenants.List))
	getSubrouter.HandleFunc("/api/v1/tenants/{name}", Respond(tenants.ListOne))

	//jobs
	postSubrouter.HandleFunc("/api/v1/jobs", Respond(jobs.Create))
	putSubrouter.HandleFunc("/api/v1/jobs/{name}", Respond(jobs.Update))
	deleteSubrouter.HandleFunc("/api/v1/jobs/{name}", Respond(jobs.Delete))
	getSubrouter.HandleFunc("/api/v1/jobs", Respond(jobs.List))
	getSubrouter.HandleFunc("/api/v1/jobs/{name}", Respond(jobs.ListOne))

	//Poem Profiles compatibility
	getSubrouter.HandleFunc("/api/v1/poems", Respond(metricProfiles.ListPoems))

	//Metric Profiles
	getSubrouter.HandleFunc("/api/v1/metric_profiles", Respond(metricProfiles.List))
	postSubrouter.HandleFunc("/api/v1/metric_profiles", Respond(metricProfiles.Create))
	deleteSubrouter.HandleFunc("/api/v1/metric_profiles/{id}", Respond(metricProfiles.Delete))
	putSubrouter.HandleFunc("/api/v1/metric_profiles/{id}", Respond(metricProfiles.Update))

	//Recalculations
	postSubrouter.HandleFunc("/api/v1/recomputations", Respond(recomputations.Create))
	getSubrouter.HandleFunc("/api/v1/recomputations", Respond(recomputations.List))

	authGetSubrouter.HandleFunc("/api/v1/factors", Respond(factors.List))

	//Status
	getSubrouter.HandleFunc("/api/v1/status/metrics/timeline/{group}", Respond(statusDetail.List))

	//Status Raw Msg
	getSubrouter.HandleFunc("/api/v1/status/metrics/msg/{hostname}/{service}/{metric}", Respond(statusMsg.List))

	//Status Endpoints
	getSubrouter.HandleFunc("/api/v1/status/endpoints/timeline/{hostname}/{service_type}", Respond(statusEndpoints.List))

	//Status Services
	getSubrouter.HandleFunc("/api/v1/status/services/timeline/{group}", Respond(statusServices.List))

	//Status Sites
	getSubrouter.HandleFunc("/api/v1/status/sites/timeline/{group}", Respond(statusEndpointGroups.List))

	//Tenants Operations
	authGetSubrouter.HandleFunc("/api/v1/tenants", Respond(tenants.List))

	http.Handle("/", mainRouter)

	//Cache
	//get_subrouter.HandleFunc("/api/v1/reset_cache", Respond("text/xml", "utf-8", ResetCache))

	//TLS support only
	config := &tls.Config{
		MinVersion: tls.VersionTLS10,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		},
		PreferServerCipherSuites: true,
	}
	server := &http.Server{Addr: cfg.Server.Bindip + ":" + strconv.Itoa(cfg.Server.Port), Handler: nil, TLSConfig: config}
	//Web service binds to server. Requests served over HTTPS.

	err := server.ListenAndServeTLS(cfg.Server.Cert, cfg.Server.Privkey)

	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
