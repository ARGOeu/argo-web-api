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

// Create function is used to implement the create tenant request.
// The request is an http POST request with the tenant description
// provided as json structure in the request body
func Create(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// if authentication procedure fails then
	// return unauthorized http status
	if authentication.Authenticate(r.Header, cfg) == false {

		output = []byte(http.StatusText(http.StatusUnauthorized))
		//If wrong api key is passed we return UNAUTHORIZED http status
		code = http.StatusUnauthorized
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	//Reading the json input from the request body
	reqBody, err := ioutil.ReadAll(r.Body)
	input := Tenant{}
	//Unmarshalling the json input into byte form
	err = json.Unmarshal(reqBody, &input)

	// Check if json body is malformed
	if err != nil {
		if err != nil {
			// Msg in xml style, to notify for malformed json
			output, err := messageXML("Malformated json input data")

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

	// Prepare structure for storing query results
	results := []Tenant{}

	//Making sure that no profile with the requested name and namespace combination already exists in the DB
	query := searchName(input.Name)
	err = mongo.Find(session, cfg.MongoDB.Db, "tenants", query, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// If results are returned for the specific name
	// then we already have an existing tenant and we must
	// abort creation notifing the user
	if len(results) > 0 {
		// Name was found so print the error message in xml
		output, err = messageXML("Tenant with the same name already exists")

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		code = http.StatusBadRequest
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err

	}

	// If no tenant exists with this name create a new one
	query = createTenant(input)
	err = mongo.Insert(session, cfg.MongoDB.Db, "tenants", query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Notify user that the tenant has been created. In xml style
	output, err = messageXML("Tenant information successfully added")

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err

}

// List function that implements the http GET request that retrieves
// all avaiable tenant information
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

	// Create structure for storing query results
	results := []Tenant{}
	// Query tenant collection for all available documents.
	// nil query param == match everything
	err = mongo.Find(session, cfg.MongoDB.Db, "tenants", nil, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	// After successfully retrieving the db results
	// call the createView function to render them into idented xml
	output, err = createView(results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// ListOne function implement an http GET request that accepts
// a name parameter urlvar and retrieves information only for the
// specific tenant
func ListOne(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	//Extracting urlvar "name" from url path
	urlValues := r.URL.Path
	nameFromURL := strings.Split(urlValues, "/")[4]

	// Try to open the mongo session
	session, err := mongo.OpenSession(cfg.MongoDB)
	defer session.Close()

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create structure to hold query results
	results := []Tenant{}

	// Create a simple query object to query by name
	query := searchName(nameFromURL)
	// Query collection tenants for the specific tenant name
	err = mongo.Find(session, cfg.MongoDB.Db, "tenants", query, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// If query returned zero result then no tenant matched this name,
	// abort and notify user accordingly
	if len(results) == 0 {

		output, err := messageXML("Tenant not found!")

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		code = http.StatusBadRequest
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	// After successfully retrieving the db results
	// call the createView function to render them into idented xml
	output, err = createView(results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// Update function used to implement update tenant request.
// This is an http PUT request that gets a specific tenant's name
// as a urlvar parameter input and a json structure in the request
// body in order to update the datastore document for the specific
// tenant
func Update(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

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
			// User provided malformed json input data
			output, err := messageXML("Malformated json input data")

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
		//Render the response into XML
		output, err = messageXML("Tenant not found")

	} else {
		//Render the response into XML
		output, err = messageXML("Tenant successfully Added")
	}

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

	// We search by name and delete the document in db
	query := searchName(nameFromURL)
	info, err := mongo.Remove(session, cfg.MongoDB.Db, "tenants", query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// info.Removed > 0 means that many documents have been removed
	// If deletion took place we notify user accordingly.
	// Else we notify that no tenant matched the specific name
	if info.Removed > 0 {
		output, err = messageXML("Tenant information successfully deleted")
	} else {
		output, err = messageXML("Tenant not found")
	}
	//Render the response into XML
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err

}
