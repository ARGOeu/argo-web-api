package resultsV5

import (
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/gorilla/mux"
)

func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {
	respond.PrepAppRoutes(s, confhandler, resultsV5Routes)
}

var resultsV5Routes = []respond.AppRoutes{
	{Name: "results.list", Verb: "GET", Path: "/{report-name}/supergroups", SubrouterHandler: GetSupergroups},
	{Name: "results.get", Verb: "GET", Path: "/{report-name}/supergroups/{supergroup-name}", SubrouterHandler: GetSupergroups},
	{Name: "results.list", Verb: "GET", Path: "/{report-name}/groups", SubrouterHandler: GetGroups},
	{Name: "results.get", Verb: "GET", Path: "/{report-name}/groups/{group-name}", SubrouterHandler: GetGroups},
	{Name: "results.list", Verb: "GET", Path: "/{report-name}/endpoints", SubrouterHandler: GetEndpoints},
	{Name: "results.get", Verb: "GET", Path: "/{report-name}/endpoints/{endpoint-name}", SubrouterHandler: GetEndpoints},
	{Name: "results.list", Verb: "GET", Path: "/{report-name}/groups/{group-name}/service-types", SubrouterHandler: GetServiceTypes},
	{Name: "results.list", Verb: "GET", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}", SubrouterHandler: GetServiceTypes},
	{Name: "results.list", Verb: "GET", Path: "/{report-name}/groups/{group-name}/endpoints", SubrouterHandler: GetEndpoints},
	{Name: "results.list", Verb: "GET", Path: "/{report-name}/groups/{group-name}/endpoints/{endpoint-name}", SubrouterHandler: GetEndpoints},
	{Name: "results.list", Verb: "GET", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints", SubrouterHandler: GetEndpoints},
	{Name: "results.list", Verb: "GET", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints/{endpoint-name}", SubrouterHandler: GetEndpoints},

	{Name: "results.options", Verb: "OPTIONS", Path: "/{report-name}/groups", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report-name}/supergroups", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report-name}/supergroups/{supergroup-name}", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report-name}/endpoints", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report-name}/endpoints/{endpoint-name}", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/endpoints", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/endpoints/{endpoint-name}", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/service-types", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints", SubrouterHandler: Options},
	{Name: "results.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints/{endpoint-name}", SubrouterHandler: Options},
}
