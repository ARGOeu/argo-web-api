package integrations

import (
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/gorilla/mux"
)

// HandleSubrouter uses the subrouter for a specific calls and creates a tree of sorts
// handling each route with a different subrouter
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {
	respond.PrepAppRoutes(s, confhandler, intRoutes)
}

var intRoutes = []respond.AppRoutes{
	{
		Name:             "v3.components.access_refresh",
		Verb:             "POST",
		Path:             "/components/{COMPONENT}/by-tenant-name/{TENANT}/refresh",
		SubrouterHandler: ComponentRefreshKey,
	},
}
