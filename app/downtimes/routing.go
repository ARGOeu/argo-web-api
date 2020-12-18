package downtimes

import (
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/gorilla/mux"
)

// HandleSubrouter uses the subrouter for a specific calls and creates a tree of sorts
// handling each route with a different subrouter
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {

	s = respond.PrepAppRoutes(s, confhandler, appRoutesV2)

}

var appRoutesV2 = []respond.AppRoutes{
	{Name: "downtimes.create", Verb: "POST", Path: "/downtimes", SubrouterHandler: Create},
	{Name: "downtimes.list", Verb: "GET", Path: "/downtimes", SubrouterHandler: List},
	{Name: "downtimes.delete", Verb: "DELETE", Path: "/downtimes", SubrouterHandler: Delete},
	{Name: "downtimes.options", Verb: "OPTIONS", Path: "/downtimes", SubrouterHandler: Options},
}
