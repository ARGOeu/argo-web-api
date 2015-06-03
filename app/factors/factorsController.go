package factors

import (
	"fmt"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"
	"net/http"
)

func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	session, err := mongo.OpenSession(cfg.MongoDB)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	results := []FactorsOutput{}
	err = mongo.Find(session, "AR", "hepspec", nil, "p", &results)

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
