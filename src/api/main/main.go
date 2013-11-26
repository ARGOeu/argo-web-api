package main

import (
	"github.com/makistsan/go-api"
	"github.com/makistsan/go-lru-cache"
	"net/http"
	"runtime"
	"strconv"
	"fmt"
	"runtime/pprof"
	"os"
	"log"
	"os/signal"
)

var httpcache *cache.LRUCache

type mystring string

func (s mystring) Size() int {
	return len(s)
}

var cfg = LoadConfiguration()

func main() {
	httpcache = cache.NewLRUCache(uint64(cfg.Server.Lrucache))
	runtime.GOMAXPROCS(cfg.Server.Maxprocs)
	
	 if *flProfile != "" {
                f, err := os.Create(*flProfile)
                if err != nil {
                        log.Fatal(err)
                }
                pprof.StartCPUProfile(f)
                defer pprof.StopCPUProfile()
        }
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
    	for sig := range c {
    	    // sig is a ^C, handle it
		if *flProfile != "" {
		pprof.StopCPUProfile()
		}
		log.Printf("captured %v, stopping profiler and exiting..", sig)
		os.Exit(1)
    	}
	}()
	defer fmt.Println("FInished") 
	handlers := map[string]func(http.ResponseWriter, *http.Request){}
	
	//Basic api calls
	handlers["/api/v1/service_availability_in_profile"] = func(w http.ResponseWriter, r *http.Request) {
		api.Respond("text/xml", "utf-8", ServiceAvailabilityInProfile)(w, r)
	}
	handlers["/api/v1/sites_availability_in_profile"] = func(w http.ResponseWriter, r *http.Request) {
		api.Respond("text/xml", "utf-8", SitesAvailabilityInProfile)(w, r)
	}
	handlers["/api/v1/ngi_availability_in_profile"] = func(w http.ResponseWriter, r *http.Request) {
		api.Respond("text/xml", "utf-8", NgiAvailabilityInProfile)(w, r)
	}
	handlers["/api/v1/profiles"] = func(w http.ResponseWriter, r *http.Request) {
		api.Respond("text/xml", "utf-8", GetProfileNames)(w, r)
	}

	//CRUD functions for profiles
	handlers["/api/v1/profiles/create"] = func(w http.ResponseWriter, r *http.Request) {
		api.Respond("text/xml", "utf-8", AddProfile)(w, r)
	}
	handlers["/api/v1/profiles/remove"] = func(w http.ResponseWriter, r *http.Request) {
		api.Respond("text/xml", "utf-8", RemoveProfile)(w, r)
	}
	handlers["/api/v1/profiles/getone"] = func(w http.ResponseWriter, r *http.Request) {
		api.Respond("text/xml", "utf-8", GetProfile)(w, r)
	}

	//Miscallenious calls
	handlers["/reset_cache"] = func(w http.ResponseWriter, r *http.Request) {
		api.Respond("text/xml", "utf-8", ResetCache)(w, r)
	}
	api.NewServer(cfg.Server.Bindip+":"+strconv.Itoa(cfg.Server.Port), api.DefaultServerReadTimeout, handlers)
}

func ResetCache(w http.ResponseWriter, r *http.Request) string {
	httpcache.Clear()
	return "Cache Emptied"
}
