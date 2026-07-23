package statusV5

import (
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/gorilla/mux"
)

func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {
	respond.PrepAppRoutes(s, confhandler, statusV5Routes)
}

var statusV5Routes = []respond.AppRoutes{
	{Name: "status.list", Verb: "GET", Path: "/{report-name}/groups", SubrouterHandler: GetGroupStatus},
	{Name: "status.get", Verb: "GET", Path: "/{report-name}/groups/{group-name}", SubrouterHandler: GetGroupStatus},
	{Name: "status.list", Verb: "GET", Path: "/{report-name}/groups/{group-name}/service-types", SubrouterHandler: GetServiceStatus},
	{Name: "status.get", Verb: "GET", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}", SubrouterHandler: GetServiceStatus},
	{Name: "status.list", Verb: "GET", Path: "/{report-name}/groups/{group-name}/endpoints", SubrouterHandler: GetEndpointStatus},
	{Name: "status.get", Verb: "GET", Path: "/{report-name}/groups/{group-name}/endpoints/{endpoint-name}", SubrouterHandler: GetEndpointStatus},
	{Name: "status.list", Verb: "GET", Path: "/{report-name}/groups/{group-name}/endpoints/{endpoint-name}/metrics", SubrouterHandler: GetMetricStatus},
	{Name: "status.get", Verb: "GET", Path: "/{report-name}/groups/{group-name}/endpoints/{endpoint-name}/metrics/{metric-name}", SubrouterHandler: GetMetricStatus},
	{Name: "status.list", Verb: "GET", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints", SubrouterHandler: GetEndpointStatus},
	{Name: "status.get", Verb: "GET", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints/{endpoint-name}", SubrouterHandler: GetEndpointStatus},
	{Name: "status.list", Verb: "GET", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints/{endpoint-name}/metrics", SubrouterHandler: GetMetricStatus},
	{Name: "status.get", Verb: "GET", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints/{endpoint-name}/metrics/{metric-name}", SubrouterHandler: GetMetricStatus},

	{Name: "status.options", Verb: "OPTIONS", Path: "/{report-name}/groups", SubrouterHandler: Options},
	{Name: "status.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}", SubrouterHandler: Options},
	{Name: "status.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/service-types", SubrouterHandler: Options},
	{Name: "status.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}", SubrouterHandler: Options},
	{Name: "status.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/endpoints", SubrouterHandler: Options},
	{Name: "status.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/endpoints/{endpoint-name}", SubrouterHandler: Options},
	{Name: "status.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/endpoints/{endpoint-name}/metrics", SubrouterHandler: Options},
	{Name: "status.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/endpoints/{endpoint-name}/metrics/{metric-name}", SubrouterHandler: Options},
	{Name: "status.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints", SubrouterHandler: Options},
	{Name: "status.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints/{endpoint-name}", SubrouterHandler: Options},
	{Name: "status.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints/{endpoint-name}/metrics", SubrouterHandler: Options},
	{Name: "status.options", Verb: "OPTIONS", Path: "/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints/{endpoint-name}/metrics/{metric-name}", SubrouterHandler: Options},
}
