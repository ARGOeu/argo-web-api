package factors

import (
	"fmt"

	"github.com/argoeu/argo-web-api/utils/authentication"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"
	"net/http"
)

// List returns a list of factors (weights) per endpoint group (i.e. site)
func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	tenantDbConfig, err := authentication.AuthenticateTenant(r.Header, cfg)

	if err != nil {
		output = []byte(http.StatusText(http.StatusUnauthorized))
		code = http.StatusUnauthorized //If wrong api key is passed we return UNAUTHORIZED http status
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	session, err := mongo.OpenSession(tenantDbConfig)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	defer mongo.CloseSession(session)

	results := []FactorsOutput{}
	err = mongo.Find(session, tenantDbConfig.Db, "weights", nil, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createView(results) //Render the results into XML format

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	mongo.CloseSession(session)
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}
