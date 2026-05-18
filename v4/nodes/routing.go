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
		Name:             "v4.results.groups.item",
		Verb:             "GET",
		Path:             "/results/groups/{item}",
		SubrouterHandler: GetGroupResults,
	},
	{
		Name:             "v4.results.groups",
		Verb:             "GET",
		Path:             "/results/groups",
		SubrouterHandler: GetGroupResults,
	},
	{
		Name:             "v4.results.groups.options",
		Verb:             "OPTIONS",
		Path:             "/results/groups",
		SubrouterHandler: Options,
	},
	{
		Name:             "v4.results.groups.item.options",
		Verb:             "OPTIONS",
		Path:             "/results/groups/{item}",
		SubrouterHandler: Options,
	},
	{
		Name:             "v4.status.groups.item",
		Verb:             "GET",
		Path:             "/status/groups/{item}",
		SubrouterHandler: GetGroupStatus,
	},
	{
		Name:             "v4.status.groups",
		Verb:             "GET",
		Path:             "/status/groups",
		SubrouterHandler: GetGroupStatus,
	},
	{
		Name:             "v4.status.groups.options",
		Verb:             "OPTIONS",
		Path:             "/status/groups",
		SubrouterHandler: Options,
	},
	{
		Name:             "v4.status.groups.item.options",
		Verb:             "OPTIONS",
		Path:             "/status/groups/{item}",
		SubrouterHandler: Options,
	},
	{
		Name:             "v4.nodes.summary.item",
		Verb:             "GET",
		Path:             "/nodes/{node_name}/capabilities/summary/{item}",
		SubrouterHandler: GetSummary,
	},
	{
		Name:             "v4.nodes.summary",
		Verb:             "GET",
		Path:             "/nodes/{node_name}/capabilities/summary",
		SubrouterHandler: GetSummary,
	},
	{
		Name:             "v4.nodes.availability",
		Verb:             "GET",
		Path:             "/nodes/{node_name}/capabilities/availability",
		SubrouterHandler: GetAvailability,
	},
	{
		Name:             "v4.nodes.availability.item",
		Verb:             "GET",
		Path:             "/nodes/{node_name}/capabilities/availability/{item}",
		SubrouterHandler: GetAvailability,
	},
	{
		Name:             "v4.nodes.uptime",
		Verb:             "GET",
		Path:             "/nodes/{node_name}/capabilities/uptime",
		SubrouterHandler: GetUptime,
	},
	{
		Name:             "v4.nodes.uptime.item",
		Verb:             "GET",
		Path:             "/nodes/{node_name}/capabilities/uptime/{item}",
		SubrouterHandler: GetUptime,
	},
	{
		Name:             "v4.nodes.status",
		Verb:             "GET",
		Path:             "/nodes/{node_name}/capabilities/status",
		SubrouterHandler: GetStatus,
	},
	{
		Name:             "v4.nodes.status.item",
		Verb:             "GET",
		Path:             "/nodes/{node_name}/capabilities/status/{item}",
		SubrouterHandler: GetStatus,
	},
	{
		Name:             "v4.nodes.status.options",
		Verb:             "OPTIONS",
		Path:             "/nodes/{node_name}/capabilities/status",
		SubrouterHandler: Options,
	},
	{
		Name:             "v4.nodes.status.item.options",
		Verb:             "OPTIONS",
		Path:             "/nodes/{node_name}/capabilities/status/{item}",
		SubrouterHandler: Options,
	},
	{
		Name:             "v4.nodes.availability.item.options",
		Verb:             "OPTIONS",
		Path:             "/nodes/{node_name}/capabilities/availability/{item}",
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
	{
		Name:             "v4.nodes.uptime.item.options",
		Verb:             "OPTIONS",
		Path:             "/nodes/{node_name}/capabilities/uptime/{item}",
		SubrouterHandler: Options,
	},
}
