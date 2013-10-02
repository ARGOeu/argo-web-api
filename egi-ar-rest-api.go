package main

import (
	"github.com/makis192/go-api"
	"net/http"
	"runtime"
)

func main() {
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
	api.NewServer(":8080", api.DefaultServerReadTimeout, handlers)
}
