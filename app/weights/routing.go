package weights

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
	{Name: "weights.create", Verb: "POST", Path: "/weights", SubrouterHandler: Create},
	{Name: "weights.update", Verb: "PUT", Path: "/weights/{ID}", SubrouterHandler: Update},
	{Name: "weights.list", Verb: "GET", Path: "/weights", SubrouterHandler: List},
	{Name: "weights.get", Verb: "GET", Path: "/weights/{ID}", SubrouterHandler: ListOne},
	{Name: "weights.delete", Verb: "DELETE", Path: "/weights/{ID}", SubrouterHandler: Delete},
	{Name: "weights.options", Verb: "OPTIONS", Path: "/weights", SubrouterHandler: Options},
	{Name: "weights.options", Verb: "OPTIONS", Path: "/weights/{ID}", SubrouterHandler: Options},
}
