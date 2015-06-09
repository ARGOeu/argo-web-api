package tenants

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/argoeu/argo-web-api/utils/authentication"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"
)

// Create function used to implement create tenant request
func Create(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	message := ""

	//Authentication procedure
	if authentication.Authenticate(r.Header, cfg) {

		session, err := mongo.OpenSession(cfg.MongoDB)

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		//Reading the json input
		reqBody, err := ioutil.ReadAll(r.Body)

		input := Tenant{}
		results := []Tenant{}
		//Unmarshalling the json input into byte form
		err = json.Unmarshal(reqBody, &input)

		if err != nil {
			if err != nil {
				message = "Malformated json input data" // User provided malformed json input data
				output, err := messageXML(message)

				if err != nil {
					code = http.StatusInternalServerError
					return code, h, output, err
				}

				code = http.StatusBadRequest
				h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
				return code, h, output, err
			}
		}

		//Making sure that no profile with the requested name and namespace combination already exists in the DB

		query := searchName(input.Name)
		err = mongo.Find(session, cfg.MongoDB.Db, "tenants", query, "name", &results)

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		// If no name is found then add the new tenant
		if len(results) <= 0 {
			//If name-namespace combination is unique we insert the new record into mongo
			query := createTenant(input)
			err = mongo.Insert(session, cfg.MongoDB.Db, "tenants", query)

			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}

			mongo.CloseSession(session)

			//Providing with the appropriate user response
			message = "Tenant inforsmation successfully added"
			output, err := messageXML(message) //Render the response into XML

			if err != nil {
				code = http.StatusInternalServerError
				return code, h, output, err
			}

			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err

		}

		// Name was found so print the error message
		message = "Tenant with the same name already exists"
		output, err := messageXML(message) //Render the response into XML

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		code = http.StatusBadRequest
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err

	}
	// Unauthorized user
	output = []byte(http.StatusText(http.StatusUnauthorized))
	code = http.StatusUnauthorized //If wrong api key is passed we return UNAUTHORIZED http status
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err

}
