package feeds

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	gcontext "github.com/gorilla/context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create a new feeds resource
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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	tenantsCol := cfg.MongoClient.Database(cfg.MongoDB.Db).Collection(tenantsColName)
	incoming := Data{}

	// Try ingest request body
	body, err := io.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
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
		{"$match": bson.M{"info.name": bson.M{"$in": incoming.Tenants}}},
		{"$project": bson.M{"id": "$id", "name": "$info.name"}},
	}

	tInfo := []TenantInfo{}
	cursor, err := tenantsCol.Aggregate(context.TODO(), aggr)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &tInfo)

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
		if !found {
			output, err = createMsgView(fmt.Sprintf("Tenant %s not found", name), 404)
			code = http.StatusNotFound
			return code, h, output, err
		}

	}

	// all incoming tenants were found

	// update the new information
	ftCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(feedsDataCol)

	_, err = ftCol.ReplaceOne(context.TODO(), bson.M{}, tenantIDs, options.Replace().SetUpsert(true))

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	incoming := Topo{}

	// Try ingest request body
	body, err := io.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
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

	ftCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(feedsTopoCol)

	_, err = ftCol.ReplaceOne(context.TODO(), bson.M{}, incoming, options.Replace().SetUpsert(true))

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	incoming := Weights{}

	// Try ingest request body
	body, err := io.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
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

	ftCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(feedsWeightsCol)

	_, err = ftCol.ReplaceOne(context.TODO(), bson.M{}, incoming, options.Replace().SetUpsert(true))

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	// Retrieve Results from database
	result := Data{}

	ftCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(feedsDataCol)
	err = ftCol.FindOne(context.TODO(), bson.M{}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {

			output, err = createMsgView("No tenant data feeds were defined. Please specify new ones!", 404)
			code = http.StatusNotFound
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	tenantNames := Data{}

	if len(result.Tenants) > 0 {

		tenantsCol := cfg.MongoClient.Database(cfg.MongoDB.Db).Collection(tenantsColName)

		// prep query to find info based on all tenant ids given
		aggr := []bson.M{
			{"$match": bson.M{"id": bson.M{"$in": result.Tenants}}},
			{"$project": bson.M{"id": "$id", "name": "$info.name"}},
		}

		tInfo := []TenantInfo{}

		cursor, err := tenantsCol.Aggregate(context.TODO(), aggr)

		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())

		cursor.All(context.TODO(), &tInfo)

		// check if incoming tenant IDs were found
		for _, ID := range result.Tenants {
			found := false
			for _, tItem := range tInfo {
				if ID == tItem.ID {
					tenantNames.Tenants = append(tenantNames.Tenants, tItem.Name)
					found = true
				}
			}
			if !found {
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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	// Retrieve Results from database
	result := Topo{}

	ftCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(feedsTopoCol)
	err = ftCol.FindOne(context.TODO(), bson.M{}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	// Retrieve Results from database
	result := Weights{}

	ftCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(feedsWeightsCol)
	err = ftCol.FindOne(context.TODO(), bson.M{}).Decode(&result)
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
	h.Set("Allow", "GET, POST, DELETE, PUT, OPTIONS")
	return code, h, output, err

}
