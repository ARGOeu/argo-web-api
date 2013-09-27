package main

import (
	"github.com/dpapathanasiou/go-api"
	"net/http"
)

func main() {
	handlers := map[string]func(http.ResponseWriter, *http.Request){}
	handlers["/service_availability_in_profile"] = func(w http.ResponseWriter, r *http.Request) {
		api.Respond("text/xml", "utf-8", ServiceAvailabilityInProfile)(w, r)
	}
	handlers["/sites_availability_in_profile"] = func(w http.ResponseWriter, r *http.Request) {
                api.Respond("text/xml", "utf-8", SitesAvailabilityInProfile)(w, r)
	}
	api.NewServer(8080, api.DefaultServerReadTimeout, handlers)
}
