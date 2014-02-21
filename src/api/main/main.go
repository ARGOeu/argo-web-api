/*
 * Copyright (c) 2013 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
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
	"github.com/gorilla/mux"
	"github.com/makistsan/go-lru-cache"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strconv"
)

var httpcache *cache.LRUCache

type mystring string

func (s mystring) Size() int {
	return len(s)
}

// Load the configurations that we have set through flags and through the configuration file
var cfg = LoadConfiguration()

func init() {

	//Create a recover function to log the case of a failure
	defer func() {
		if err := recover(); err != nil {
			log.Println("work failed:", err)
		}
	}()

	//Initialize the cache
	httpcache = cache.NewLRUCache(uint64(cfg.Server.Lrucache))

	//Set GOMAXPROCS
	runtime.GOMAXPROCS(cfg.Server.Maxprocs)

	//Start the profiler if the flag flProfile is set to a filename, where profile data will be writter
	if *flProfile != "" {
		f, err := os.Create(*flProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	//Catch an terminate signal and write all profiling data before exiting
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			if *flProfile != "" {
				pprof.StopCPUProfile()
			}
			log.Printf("captured %v, stopping profiler and exiting..", sig)
			os.Exit(1)
		}
	}()
}

func main() {

	//Create the server router
	main_router := mux.NewRouter()
	//first_subrouter := main_router.Headers("x-api-key","").Subrouter()//routes only the requets that provide an api key
	get_subrouter := main_router.Methods("GET").Subrouter()                          //routes only GET requests
	post_subrouter := main_router.Methods("POST").Subrouter()                        //routes only POST requests
	auth_subrouter := post_subrouter.Headers("x-api-key", "", "x-api-requestor", "") //calls requested with POST must provide authentication credentials otherwise will not be routed

	//Basic api calls
	get_subrouter.HandleFunc("/api/v1/service_availability_in_profile", Respond("text/xml", "utf-8", ServiceAvailabilityInProfile))
	get_subrouter.HandleFunc("/api/v1/sites_availability_in_profile", Respond("text/xml", "utf-8", SitesAvailabilityInProfile))
	get_subrouter.HandleFunc("/api/v1/ngi_availability_in_profile", Respond("text/xml", "utf-8", NgiAvailabilityInProfile))
	//get_subrouter.HandleFunc("/api/v1/service_flavor_availability_in_profile", Respond("text/xml", "utf-8", ServiceFlavorAvailabilityInProfile))
	//CRUD functions for profiles
	post_subrouter.HandleFunc("/api/v1/profiles/create", Respond("text/xml", "utf-8", AddProfile))
	get_subrouter.HandleFunc("/api/v1/profiles", Respond("text/xml", "utf-8", GetProfileNames))
	get_subrouter.HandleFunc("/api/v1/profiles/getone", Respond("text/xml", "utf-8", GetProfile))
	//SOME UPDATE METHOD MISSING
	post_subrouter.HandleFunc("/api/v1/profiles/remove", Respond("text/xml", "utf-8", RemoveProfile))
	//Miscallenious calls
	get_subrouter.HandleFunc("/api/v1/reset_cache", Respond("text/xml", "utf-8", ResetCache))
	post_subrouter.HandleFunc("/api/v1/recalculate", Respond("text/xml", "utf-8", Recalculate))
	http.Handle("/", main_router)
	//Web service binds to server.
	//plain http bidning to be removed
	// err := http.ListenAndServe(cfg.Server.Bindip+":"+strconv.Itoa(cfg.Server.Port), nil)
	// 		if err != nil {
	// 			log.Fatal("ListenAndServe:", err)
	// 		}
	//HTTPS support for the API server. We have to issue a valid certificate for our production server and replace the parameters with the actual path where the certificate will be placed
	err := http.ListenAndServeTLS(cfg.Server.Bindip+":"+strconv.Itoa(cfg.Server.Port), "/etc/pki/tls/certs/localhost.crt", "/etc/pki/tls/private/localhost.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
