package respond

import (
	"log"
	"net/http"
	"strings"

	"github.com/ARGOeu/argo-web-api/utils/authentication"
	"github.com/ARGOeu/argo-web-api/utils/authorization"
	"github.com/ARGOeu/argo-web-api/utils/config"
	gcontext "github.com/gorilla/context"
)

// // WrapAll Wraps all wrap handlers. Note: Precedence is inversed
// func WrapAll(handler http.Handler, cfg config.Config) http.Handler {
//
// 	handler = WrapValidate(handler)
// 	handler = WrapAuthorize(handler)
// 	handler = WrapAuthenticate(handler, cfg)
//
// 	return handler
// }

func isComponentRoute(routeName string) bool {
	return strings.HasPrefix(routeName, "v3.components")
}

func needsAPIAdmin(routeName string) bool {

	routePart := strings.Split(routeName, ".")[0]
	if routePart == "tenants" || routePart == "metrics_admin" {
		return true
	}

	return false
}

// WrapAuthenticate handle wrapper to apply authentication
func WrapAuthenticate(hfn http.Handler, cfg config.Config, routeName string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var errs []ErrorResponse
		// check if component route

		// check if has x-tenant-id header
		adminTenantId := r.Header.Get("x-tenant-id")

		if isComponentRoute(routeName) {

			compRole := authentication.GetComponentRole(r.Header, cfg)
			// check if user has a component role
			if compRole != "" {
				gcontext.Set(r, "roles", []string{compRole})
			} else {
				// Because user has no component role
				Error(w, r, ErrAuthen, cfg, errs)
				return
			}
			hfn.ServeHTTP(w, r)

		} else if needsAPIAdmin(routeName) {
			// check if api admin authentication is needed (for tenants etc...)
			if !(authentication.AuthenticateAdmin(r.Header, cfg)) {
				// Because not authenticated respond with error
				Error(w, r, ErrAuthen, cfg, errs)
				return
			}

			// admin api authenticated so continue serving
			gcontext.Set(r, "authen", true)
			// Add admin restricted or not information -- used in get tenants
			if authentication.IsAdminRestricted(r.Header, cfg) {
				gcontext.Set(r, "roles", []string{"super_admin_restricted"})
			} else if authentication.IsSuperAdminUI(r.Header, cfg) {
				gcontext.Set(r, "roles", []string{"super_admin_ui"})
			} else {
				gcontext.Set(r, "roles", []string{"super_admin"})
			}

			hfn.ServeHTTP(w, r)

		} else if adminTenantId != "" {
			// check if it is an admin trying to access a tenant route
			authen, username, email := authentication.AuthenticateAdminName(r.Header, cfg)
			if !(authen) {
				// Because not authenticated respond with error
				Error(w, r, ErrAuthen, cfg, errs)
				return
			}

			tenantConf, tenantName, tErr := authentication.AuthenticateAdminTenant(r.Header, cfg)
			// If error respond with error
			if tErr != nil {
				Error(w, r, ErrAuthen, cfg, errs)
				return
			}

			tenantConf.User = username
			tenantConf.Email = email
			if authentication.IsAdminRestricted(r.Header, cfg) {
				gcontext.Set(r, "roles", []string{"super_admin_restricted", "viewer"})
				tenantConf.Roles = []string{"viewer"}
			} else if authentication.IsSuperAdminUI(r.Header, cfg) {
				gcontext.Set(r, "roles", []string{"super_admin_ui", "admin_ui"})
				tenantConf.Roles = []string{"viewer"}
			} else {
				gcontext.Set(r, "roles", []string{"super_admin", "admin"})
				tenantConf.Roles = []string{"admin"}
			}

			gcontext.Set(r, "tenant_conf", tenantConf)
			gcontext.Set(r, "tenant_name", tenantName)
			gcontext.Set(r, "authen", authen)
			log.Printf("Admin User: %s Accessing Tenant: %s", username, tenantName)
			hfn.ServeHTTP(w, r)

		} else {

			// authenticate tenant user
			tenantConf, name, tErr := authentication.AuthenticateTenant(r.Header, cfg)
			// If tenant user not authenticated respond with  error
			if tErr != nil {
				Error(w, r, ErrAuthen, cfg, errs)
				return
			}

			gcontext.Set(r, "roles", tenantConf.Roles)
			gcontext.Set(r, "tenant_conf", tenantConf)
			gcontext.Set(r, "tenant_name", name)
			gcontext.Set(r, "authen", true)
			hfn.ServeHTTP(w, r)

		}

	})
}

// WrapAuthorize handle wrapper to apply authorization
func WrapAuthorize(hfn http.Handler, cfg config.Config, routeName string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var errs []ErrorResponse

		roles := gcontext.Get(r, "roles").([]string)

		if roles != nil {

			author := authorization.HasResourceRoles(cfg, routeName, roles)
			if author {
				hfn.ServeHTTP(w, r)
				return
			}
		}

		Error(w, r, ErrAuthor, cfg, errs)

	})
}

// WrapValidate handle wrapper to apply validation
func WrapValidate(hfn http.Handler, cfg config.Config, routeName string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var errs []ErrorResponse

		headers := r.Header
		queries := r.URL.Query()
		resource := strings.Split(routeName, ".")[0]

		// Validate Accept Header (globally unless OPTIONS verb is used)
		if r.Method != "OPTIONS" {
			err := ValidateAcceptHeader(headers.Get("Accept"))
			if err != (ErrorResponse{}) {
				Error(w, r, ErrValidHead, cfg, errs)
				return
			}
			if strings.Contains(resource, "status") {
				errs = ValidateStatusParams(queries)
				if len(errs) > 0 {
					Error(w, r, ErrValidQuery, cfg, errs)
					return
				}
			}
			if strings.Contains(resource, "results") {
				errs = ValidateResultsParams(queries)
				if len(errs) > 0 {
					Error(w, r, ErrValidQuery, cfg, errs)
					return
				}
			}
			if strings.Contains(resource, "metricResult") {
				errs = ValidateMetricParams(queries)
				if len(errs) > 0 {
					Error(w, r, ErrValidQuery, cfg, errs)
					return
				}
			}
		}

		hfn.ServeHTTP(w, r)
	})
}
