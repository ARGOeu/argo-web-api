package nodes

import (
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/gorilla/mux"
)

// HandleSubrouter uses the subrouter for a specific calls and creates a tree of sorts
// handling each route with a different subrouter
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {
	respond.PrepAppRoutes(s, confhandler, arRoutes)
}

var arRoutes = []respond.AppRoutes{
	{
		Name:             "v4.nodes.availability",
		Verb:             "GET",
		Path:             "/nodes/{node_name}/capabilities/availability",
		SubrouterHandler: GetAvailability,
	},
	{
		Name:             "v4.nodes.uptime",
		Verb:             "GET",
		Path:             "/nodes/{node_name}/capabilities/uptime",
		SubrouterHandler: GetUptime,
	},
	{
		Name:             "v4.nodes.status",
		Verb:             "GET",
		Path:             "/nodes/{node_name}/capabilities/status",
		SubrouterHandler: GetStatus,
	},
	{
		Name:             "v4.nodes.status.options",
		Verb:             "OPTIONS",
		Path:             "/nodes/{node_name}/capabilities/status",
		SubrouterHandler: Options,
	},
	{
		Name:             "v4.nodes.availability.options",
		Verb:             "OPTIONS",
		Path:             "/nodes/{node_name}/capabilities/availability",
		SubrouterHandler: Options,
	},
	{
		Name:             "v4.nodes.uptime.options",
		Verb:             "OPTIONS",
		Path:             "/nodes/{node_name}/capabilities/uptime",
		SubrouterHandler: Options,
	},
}
