package feeds

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2/bson"
)

//Create a new feeds resource
const feedsTopoCol = "feeds_topology"
const feedsWeightsCol = "feeds_weights"
const feedsDataCol = "feeds_data"
const tenantsColName = "tenants"

// Update request handler creates a new feed topo resource
func UpdateData(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
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

	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	// open a master session to the argo core database
	coreSession, err := mongo.OpenSession(cfg.MongoDB)
	defer mongo.CloseSession(coreSession)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	tenantsCol := coreSession.DB(cfg.MongoDB.Db).C(tenantsColName)
	incoming := Data{}

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

	// prep query to find info based on all tenant names given in the incoming list
	aggr := []bson.M{
		bson.M{"$match": bson.M{"info.name": bson.M{"$in": incoming.Tenants}}},
		bson.M{"$project": bson.M{"id": "$id", "name": "$info.name"}},
	}

	// open a normal session to the specific tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	tInfo := []TenantInfo{}
	err = tenantsCol.Pipe(aggr).All(&tInfo)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	tenantIDs := Data{}

	// check if incoming tenant names were found
	for _, name := range incoming.Tenants {
		found := false
		for _, tItem := range tInfo {

			if name == tItem.Name {
				tenantIDs.Tenants = append(tenantIDs.Tenants, tItem.ID)
				found = true
			}
		}
		if found == false {
			output, err = createMsgView(fmt.Sprintf("Tenant %s not found", name), 404)
			code = http.StatusNotFound
			return code, h, output, err
		}

	}

	// all incoming tenants were found

	// update the new information

	ftCol := session.DB(tenantDbConfig.Db).C(feedsDataCol)
	_, err = ftCol.Upsert(bson.M{}, tenantIDs)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createDataListView([]Data{incoming}, "Feeds resource succesfully updated", 200) //Render the results into JSON
	code = 200
	return code, h, output, err
}

// Update request handler creates a new feed topo resource
func UpdateTopo(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
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

	incoming := Topo{}

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

	_, err = mongo.Remove(session, tenantDbConfig.Db, feedsTopoCol, bson.M{})
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	err = mongo.Insert(session, tenantDbConfig.Db, feedsTopoCol, incoming)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createTopoListView([]Topo{incoming}, "Feeds resource succesfully updated", 200) //Render the results into JSON
	code = 200
	return code, h, output, err
}

// UpdateWeights request handler creates a new weights feed resource
func UpdateWeights(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
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

	_, err = mongo.Remove(session, tenantDbConfig.Db, feedsWeightsCol, bson.M{})
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	err = mongo.Insert(session, tenantDbConfig.Db, feedsWeightsCol, incoming)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createWeightsListView([]Weights{incoming}, "Feeds resource succesfully updated", 200) //Render the results into JSON
	code = 200
	return code, h, output, err
}

// ListData lists data feeds Results
func ListData(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	// Open session to tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Retrieve Results from database
	result := Data{}

	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	ftCol := session.DB(tenantDbConfig.Db).C(feedsDataCol)
	err = ftCol.Find(bson.M{}).One(&result)

	if err != nil {
		if err.Error() == "not found" {

			output, err = createMsgView("No tenant data feeds were defined. Please specify new ones!", 404)
			code = http.StatusNotFound
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	tenantNames := Data{}

	if len(result.Tenants) > 0 {

		// open a master session to the argo core database
		coreSession, err := mongo.OpenSession(cfg.MongoDB)
		defer mongo.CloseSession(coreSession)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		tenantsCol := coreSession.DB(cfg.MongoDB.Db).C(tenantsColName)

		// prep query to find info based on all tenant ids given
		aggr := []bson.M{
			bson.M{"$match": bson.M{"id": bson.M{"$in": result.Tenants}}},
			bson.M{"$project": bson.M{"id": "$id", "name": "$info.name"}},
		}

		tInfo := []TenantInfo{}

		err = tenantsCol.Pipe(aggr).All(&tInfo)

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		// check if incoming tenant IDs were found
		for _, ID := range result.Tenants {
			found := false
			for _, tItem := range tInfo {
				if ID == tItem.ID {
					tenantNames.Tenants = append(tenantNames.Tenants, tItem.Name)
					found = true
				}
			}
			if found == false {
				output, err = createMsgView(fmt.Sprintf("Tenant with ID: %s not found. Please update the feed correctly!", ID), 404)
				code = http.StatusNotFound
				return code, h, output, err
			}

		}

		// all incoming tenants were found
	}

	// Create view of the results
	output, err = createDataListView([]Data{tenantNames}, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// ListTopo lists topology Results
func ListTopo(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	// Open session to tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Retrieve Results from database
	result := Topo{}

	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	ftCol := session.DB(tenantDbConfig.Db).C(feedsTopoCol)
	err = ftCol.Find(bson.M{}).One(&result)
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
	output, err = createTopoListView([]Topo{result}, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// ListWeight list weights feeds results
func ListWeights(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	// Open session to tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Retrieve Results from database
	result := Weights{}

	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	ftCol := session.DB(tenantDbConfig.Db).C(feedsWeightsCol)
	err = ftCol.Find(bson.M{}).One(&result)
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
	output, err = createWeightsListView([]Weights{result}, "Success", code) //Render the results into JSON

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
