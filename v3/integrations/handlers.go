package integrations

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"github.com/ARGOeu/argo-web-api/app/tenants"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ROLE_DICT = map[string]string{
	"engine":      "component_engine",
	"monbox":      "component_monbox",
	"probe":       "component_monbox",
	"poem-admin":  "component_poem",
	"poem-viewer": "component_poem",
}

func ComponentRefreshKey(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	vars := mux.Vars(r)
	roles := gcontext.Get(r, "roles").([]string)

	// check if the user has the role for the component
	compRole, exists := ROLE_DICT[vars["COMPONENT"]]
	if !exists || !hasRole(roles, compRole) {
		output, _ = respond.MarshalContent(respond.Forbidden, contentType, "", " ")
		code = http.StatusNotFound
		return code, h, output, err
	}

	// Create structure to hold query results
	result := tenants.Tenant{}

	// Try to get mongo client and target tenant collection
	tenantCol := cfg.MongoClient.Database(cfg.MongoDB.Db).Collection("tenants")

	// Create a simple query object to query tenant by name
	query := bson.M{"info.name": vars["TENANT"]}

	// Query collection tenants for the specific tenant name
	err = tenantCol.FindOne(context.TODO(), query).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			output, _ = respond.QuickResponse(fmt.Sprintf("tenant %s not found", vars["TENANT"]), 404)
			code = http.StatusNotFound
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	tenant := result

	token, err := genAccessKey(32)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	found := false
	for indx, user := range tenant.Users {
		if user.Component != "" && user.Component == vars["COMPONENT"] {
			found = true
			tenant.Users[indx].APIkey = token
		}
	}

	if !found {
		output, _ = respond.QuickResponse(fmt.Sprintf("Component %s not configured for tenant %s", vars["COMPONENT"], vars["TENANT"]), 404)
		code = http.StatusNotFound
		return code, h, output, err
	}

	if errMsg, errCode := tenants.ValidateTenantUsers(tenant, tenantCol); errMsg != "" && errCode != 0 {
		output, _ = respond.MarshalContent(respond.ErrConflict(errMsg), contentType, "", " ")
		code = errCode
		return code, h, output, err
	}

	// run the update query

	tenant.Info.Updated = time.Now().Format("2006-01-02 15:04:05")
	query = bson.M{"info.name": vars["TENANT"]}
	replaceResult, err := tenantCol.ReplaceOne(context.TODO(), query, tenant)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if replaceResult.MatchedCount == 0 {
		output, _ = respond.QuickResponse(fmt.Sprintf("Component %s not configured for tenant %s", vars["COMPONENT"], vars["TENANT"]), 400)
		code = http.StatusNotFound
		return code, h, output, err
	}

	// Create view for response message
	output, err = tenants.CreateRenewedToken(token, "component api key succesfully renewed", 200) //Render the results into JSON

	code = http.StatusOK
	return code, h, output, err
}

func hasRole(roles []string, checkRole string) bool {
	return len(roles) > 0 && roles[0] == checkRole
}

func genAccessKey(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
