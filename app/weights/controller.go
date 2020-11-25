package weights

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

//Create a new weights resource
const weightCol = "weights"

func prepMultiQuery(dt int, name string) interface{} {

	matchQuery := bson.M{"date_integer": bson.M{"$lte": dt}}

	if name != "" {
		matchQuery["name"] = name
	}

	return []bson.M{
		{
			"$match": matchQuery,
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"id": "$id",
				},
				// weights collection is meant to have an index with date_integer:-1 and id:1 so
				// when searching by date the documents are sorted with the recent timestamp first
				// so we need the recent item available to our query timepoint which is specific date
				"id":          bson.M{"$first": "$id"},
				"date":        bson.M{"$first": "$date"},
				"name":        bson.M{"$first": "$name"},
				"weight_type": bson.M{"$first": "$weight_type"},
				"group_type":  bson.M{"$first": "$group_type"},
				"groups":      bson.M{"$first": "$groups"},
			},
		},
		{
			"$sort": bson.M{"id": 1},
		},
	}

}

func prepQuery(dt int, id string) interface{} {

	return bson.M{"date_integer": bson.M{"$lte": dt}, "id": id}

}

// Create request handler creates a new weight resource
func Create(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
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

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)
	urlValues := r.URL.Query()
	dateStr := urlValues.Get("date")
	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	incoming := Weights{}

	// Try ingest request body
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	if err := r.Body.Close(); err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Parse body json
	if err := json.Unmarshal(body, &incoming); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
		code = 400
		return code, h, output, err
	}

	incoming.DateInt = dt
	incoming.Date = dateStr

	// check if the weights resource name is unique
	results := []Weights{}
	query := bson.M{"name": incoming.Name}

	err = mongo.Find(session, tenantDbConfig.Db, weightCol, query, "", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// If results are returned for the specific name
	// then we already have an existing report and we must
	// abort creation notifying the user
	if len(results) > 0 {
		output, _ = respond.MarshalContent(respond.ErrConflict("Weights resource with the same name already exists"), contentType, "", " ")
		code = http.StatusConflict
		return code, h, output, err
	}

	// Generate new weights id
	incoming.ID = mongo.NewUUID()

	err = mongo.Insert(session, tenantDbConfig.Db, weightCol, incoming)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createRefView(incoming, "Weights resource succesfully created", 201, r) //Render the results into JSON
	code = 201
	return code, h, output, err
}

// ListOne handles the listing of one weight resource based on its given id
func ListOne(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	urlValues := r.URL.Query()
	dateStr := urlValues.Get("date")

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Open session to tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Retrieve Results from database
	result := Weights{}
	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}
	wQuery := prepQuery(dt, vars["ID"])

	wCol := session.DB(tenantDbConfig.Db).C(weightCol)
	err = wCol.Find(wQuery).One(&result)
	if err != nil {
		if err.Error() == "not found" {
			output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
			code = 404
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createListView([]Weights{result}, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// List the existing weight resources for the tenant making the request
func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	urlValues := r.URL.Query()
	dateStr := urlValues.Get("date")
	name := urlValues.Get("name")

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Open session to tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}
	wQuery := prepMultiQuery(dt, name)

	wCol := session.DB(tenantDbConfig.Db).C(weightCol)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	results := []Weights{}
	err = wCol.Pipe(wQuery).All(&results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createListView(results, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

//Update function to update contents of an existing weights resource
func Update(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	vars := mux.Vars(r)
	urlValues := r.URL.Query()
	dateStr := urlValues.Get("date")
	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	incoming := Weights{}

	// ingest body data
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	if err := r.Body.Close(); err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	// parse body json
	if err := json.Unmarshal(body, &incoming); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
		code = 400
		return code, h, output, err
	}

	incoming.DateInt = dt
	incoming.Date = dateStr

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	// create filter to retrieve specific profile with id
	filter := bson.M{"id": vars["ID"]}

	incoming.ID = vars["ID"]

	// Retrieve Results from database
	results := []Weights{}
	err = mongo.Find(session, tenantDbConfig.Db, weightCol, filter, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Check if nothing found
	if len(results) < 1 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	// check if the weights dataset name is unique
	if incoming.Name != results[0].Name {

		results = []Weights{}
		query := bson.M{"name": incoming.Name, "id": bson.M{"$ne": vars["ID"]}}

		err = mongo.Find(session, tenantDbConfig.Db, weightCol, query, "", &results)

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		// If results are returned for the specific name
		// then we already have an existing operations profile and we must
		// abort creation notifying the user
		if len(results) > 0 {
			output, _ = respond.MarshalContent(respond.ErrConflict("Weights resource with the same name already exists"), contentType, "", " ")
			code = http.StatusConflict
			return code, h, output, err
		}
	}
	// run the update query
	wCol := session.DB(tenantDbConfig.Db).C(weightCol)
	info, err := wCol.Upsert(bson.M{"id": vars["ID"], "date_integer": dt}, incoming)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	updMsg := "Weights resource successfully updated"

	if info.Updated <= 0 {
		updMsg = "Weights resource successfully updated (new history snapshot)"
	}

	// Create view for response message
	output, err = createMsgView(updMsg, 200) //Render the results into JSON
	code = 200
	return code, h, output, err
}

//Delete weights resource based on ID
func Delete(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	// Grab Tenant DB configuration from context
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Open session to tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	filter := bson.M{"id": vars["ID"]}

	// Retrieve Results from database
	results := []Weights{}
	err = mongo.Find(session, tenantDbConfig.Db, weightCol, filter, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Check if nothing found
	if len(results) < 1 {
		output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	// Remove also removes multiple documents with specific id, thus erasing all history
	mongo.Remove(session, tenantDbConfig.Db, weightCol, filter)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createMsgView("Weights resource successfully deleted", 200) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// Options request handler
func Options(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/plain"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	h.Set("Allow", fmt.Sprintf("GET, POST, DELETE, PUT, OPTIONS"))
	return code, h, output, err

}
