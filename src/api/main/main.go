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
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"github.com/makistsan/go-lru-cache"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strconv"
	"github.com/gorilla/mux"
)

var httpcache *cache.LRUCache

type mystring string

func (s mystring) Size() int {
	return len(s)
}

// Load the configurations that we have set through flags and through the configuration file
var cfg = LoadConfiguration()


// The respond function that will be called to answer to http requests to the PI
func Respond(mediaType string, charset string, fn func(w http.ResponseWriter, r *http.Request) []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=%s", mediaType, charset))
		output := fn(w, r)
		var b bytes.Buffer
		var data []byte
		if (cfg.Server.Gzip) == true && r.Header.Get("Accept-Encoding") != "" {
			encodings := parseCSV(r.Header.Get("Accept-Encoding"))
			for _, val := range encodings {
				if val == "gzip" {
					writer := gzip.NewWriter(&b)
					writer.Write(output)
					writer.Close()
					w.Header().Set("Content-Encoding", "gzip")
					break
				} else if val == "deflate" {
					writer := zlib.NewWriter(&b)
					writer.Write(output)
					writer.Close()
					w.Header().Set("Content-Encoding", "deflate")
					break
				}
			}
			data = b.Bytes()
		} else {
			data = output
		}
		fmt.Println(len(data))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.Write(data)
	}
}



func main() {

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


	//Create the server router
	r := mux.NewRouter()

	
	//Basic api calls
	r.HandleFunc("/api/v1/service_availability_in_profile", Respond("text/xml", "utf-8", ServiceAvailabilityInProfile))	
	r.HandleFunc("/api/v1/sites_availability_in_profile", Respond("text/xml", "utf-8", SitesAvailabilityInProfile))	
	r.HandleFunc("/api/v1/ngi_availability_in_profile", Respond("text/xml", "utf-8", NgiAvailabilityInProfile))		
	//CRUD functions for profiles
	
	r.HandleFunc("/api/v1/profiles", Respond("text/xml", "utf-8", GetProfileNames))	
	r.HandleFunc("/api/v1/profiles/create", Respond("text/xml", "utf-8", AddProfile))	
	r.HandleFunc("/api/v1/profiles/remove", Respond("text/xml", "utf-8", RemoveProfile))	
	r.HandleFunc("/api/v1/profiles/getone", Respond("text/xml", "utf-8", GetProfile))	
	//Miscallenious calls
	r.HandleFunc("/reset_cache", Respond("text/xml", "utf-8", ResetCache))	

	http.Handle("/", r)
	err := http.ListenAndServe(cfg.Server.Bindip+":"+strconv.Itoa(cfg.Server.Port), nil)
    	if err != nil {
        	log.Fatal("ListenAndServe:", err)
    	}
}


//Reset the cache if it is set
func ResetCache(w http.ResponseWriter, r *http.Request) []byte {
	if cfg.Server.Cache == true {
	httpcache.Clear()
	return []byte("Cache Emptied")
	}
	return []byte("No Caching is active")
}
