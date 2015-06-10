package tenants

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

	// if authentication procedure fails then
	// return unauthorized
	if authentication.Authenticate(r.Header, cfg) == false {

		output = []byte(http.StatusText(http.StatusUnauthorized))
		//If wrong api key is passed we return UNAUTHORIZED http status
		code = http.StatusUnauthorized
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
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

	// Try to open the mongo session
	session, err := mongo.OpenSession(cfg.MongoDB)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	//Making sure that no profile with the requested name and namespace combination already exists in the DB

	query := searchName(input.Name)
	err = mongo.Find(session, cfg.MongoDB.Db, "tenants", query, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if len(results) > 0 {
		// Name was found so print the error message
		message = "Tenant with the same name already exists"
		output, err = messageXML(message) //Render the response into XML

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		code = http.StatusBadRequest
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err

	}

	//If name-namespace combination is unique we insert the new record into mongo
	query = createTenant(input)
	err = mongo.Insert(session, cfg.MongoDB.Db, "tenants", query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	//Providing with the appropriate user response
	message = "Tenant information successfully added"
	output, err = messageXML(message) //Render the response into XML

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err

}

// List function
func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Try to open the mongo session
	session, err := mongo.OpenSession(cfg.MongoDB)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	results := []Tenant{}
	err = mongo.Find(session, cfg.MongoDB.Db, "tenants", nil, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createView(results) //Render the results into XML format

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// ListOne function
func ListOne(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	//Extracting record id from url
	urlValues := r.URL.Path
	nameFromURL := strings.Split(urlValues, "/")[4]

	// Try to open the mongo session
	session, err := mongo.OpenSession(cfg.MongoDB)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	query := searchName(nameFromURL)

	results := []Tenant{}
	err = mongo.Find(session, cfg.MongoDB.Db, "tenants", query, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if len(results) == 0 {
		message := "Tenant not found!"
		output, err := messageXML(message)

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		code = http.StatusBadRequest
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	output, err = createView(results) //Render the results into XML format

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// Update function used to implement update tenant request
func Update(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	message := ""

	// if authentication procedure fails then
	// return unauthorized
	if authentication.Authenticate(r.Header, cfg) == false {

		output = []byte(http.StatusText(http.StatusUnauthorized))
		//If wrong api key is passed we return UNAUTHORIZED http status
		code = http.StatusUnauthorized
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	//Extracting record id from url
	urlValues := r.URL.Path
	nameFromURL := strings.Split(urlValues, "/")[4]

	//Reading the json input
	reqBody, err := ioutil.ReadAll(r.Body)

	input := Tenant{}
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

	// Try to open the mongo session
	session, err := mongo.OpenSession(cfg.MongoDB)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// We search by name and update
	query := searchName(nameFromURL)
	err = mongo.Update(session, cfg.MongoDB.Db, "tenants", query, input)

	if err != nil {

		if err.Error() != "not found" {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		message = "Tenant not found"

	} else {
		message = "Tenant successfully updated"
	}

	output, err = messageXML(message) //Render the response into XML

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err

}

// Delete function used to implement remove tenant request
func Delete(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	message := ""

	// if authentication procedure fails then
	// return unauthorized
	if authentication.Authenticate(r.Header, cfg) == false {

		output = []byte(http.StatusText(http.StatusUnauthorized))
		//If wrong api key is passed we return UNAUTHORIZED http status
		code = http.StatusUnauthorized
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	//Extracting record id from url
	urlValues := r.URL.Path
	nameFromURL := strings.Split(urlValues, "/")[4]

	// Try to open the mongo session
	session, err := mongo.OpenSession(cfg.MongoDB)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// We search by name and delete
	query := searchName(nameFromURL)
	info, err := mongo.Remove(session, cfg.MongoDB.Db, "tenants", query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	//Providing with the appropriate user response
	if info.Removed > 0 {
		message = "Tenant information successfully deleted"
	} else {
		message = "Tenant not found"
	}
	output, err = messageXML(message) //Render the response into XML

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err

}
