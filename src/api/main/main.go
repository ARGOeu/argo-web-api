package main

import (
	"github.com/makistsan/go-api"
	"github.com/makistsan/go-lru-cache"
	"net/http"
	"runtime"
)

var httpcache *cache.LRUCache

type mystring string

func (s mystring) Size() int {
	return len(s)
}

func main() {
	httpcache = cache.NewLRUCache(uint64(700000000))
	runtime.GOMAXPROCS(4)
	handlers := map[string]func(http.ResponseWriter, *http.Request){}
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

	handlers["/reset_cache"] = func(w http.ResponseWriter, r *http.Request) {
		api.Respond("text/xml", "utf-8", ResetCache)(w, r)
	}
	api.NewServer(":8080", api.DefaultServerReadTimeout, handlers)
}

func ResetCache(w http.ResponseWriter, r *http.Request) string {
	httpcache.Clear()
	return "Cache Emptied"
}
