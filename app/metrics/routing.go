package metrics

import (
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/gorilla/mux"
)

// HandleAdminSubrouter includes the admin specific calls (put/get)
func HandleAdminSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {

	respond.PrepAppRoutes(s, confhandler, appAdminRoutesV2)

}

// HandleSubrouter includes the tenant only specific calls (get only)
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {

	respond.PrepAppRoutes(s, confhandler, appRoutesV2)

}

var appAdminRoutesV2 = []respond.AppRoutes{
	{Name: "metrics_admin.update", Verb: "PUT", Path: "/metrics", SubrouterHandler: UpdateMetrics},
	{Name: "metrics_admin.get", Verb: "GET", Path: "/metrics", SubrouterHandler: ListMetrics},
	{Name: "metrics_admin.options", Verb: "OPTIONS", Path: "/metrics", SubrouterHandler: Options},
}

var appRoutesV2 = []respond.AppRoutes{
	{Name: "metrics_report.get", Verb: "GET", Path: "/metrics/by_report/{report}", SubrouterHandler: ListMetricsByReport},
	{Name: "metrics.get", Verb: "GET", Path: "/metrics", SubrouterHandler: ListMetrics},
	{Name: "metrics.options", Verb: "OPTIONS", Path: "/metrics", SubrouterHandler: Options},
}
