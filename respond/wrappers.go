package respond

import (
	"net/http"
	"strings"

	"github.com/ARGOeu/argo-web-api/utils/authentication"
	"github.com/ARGOeu/argo-web-api/utils/authorization"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/gorilla/context"
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

func needsAPIAdmin(routeName string) bool {
	if strings.Split(routeName, ".")[0] == "tenants" {
		return true
	}

	return false
}

// WrapAuthenticate handle wrapper to apply authentication
func WrapAuthenticate(hfn http.Handler, cfg config.Config, routeName string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// check if api admin authentication is needed (for tenants etc...)
		if needsAPIAdmin(routeName) {

			// Authenticate admin api and check
			if (authentication.AuthenticateAdmin(r.Header, cfg)) == false {
				// Because not authenticated respond with error
				Error(w, r, ErrAuthen, cfg)
				return
			}
			// admin api authenticated so continue serving
			context.Set(r, "authen", true)
			context.Set(r, "roles", []string{"super_admin"})

			hfn.ServeHTTP(w, r)

		} else {

			// authenticate tenant user
			tenantConf, tErr := authentication.AuthenticateTenant(r.Header, cfg)
			// If tenant user not authenticated respond with  error
			if tErr != nil {
				Error(w, r, ErrAuthen, cfg)
				return
			}

			context.Set(r, "roles", tenantConf.Roles)
			context.Set(r, "tenant_conf", tenantConf)
			context.Set(r, "authen", true)
			hfn.ServeHTTP(w, r)

		}

	})
}

// WrapAuthorize handle wrapper to apply authorization
func WrapAuthorize(hfn http.Handler, cfg config.Config, routeName string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// log.Printf(" >> Authorization takes place here...")
		var roles []string

		roles = context.Get(r, "roles").([]string)

		if roles != nil {
			author := authorization.HasResourceRoles(cfg, routeName, roles)
			if author != false {
				hfn.ServeHTTP(w, r)
				return
			}
		}

		Error(w, r, ErrAuthor, cfg)

	})
}

// WrapValidate handle wrapper to apply validation
func WrapValidate(hfn http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// log.Printf(" >> Validation takes place here...")
		hfn.ServeHTTP(w, r)
	})
}
