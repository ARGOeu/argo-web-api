package downtimes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils"
	"github.com/ARGOeu/argo-web-api/utils/config"
	gcontext "github.com/gorilla/context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create a new downtimes resource
const downtimeColName = "downtimes"

// prepFilterQuery produces an aggregation query in downtimes collection that allows to filter endpoints by severity and/or classification
// the modeled query will resemble the following mongo aggregation
// db.downtimes.aggregate([
//
//	{  "$match": {"date_integer": 20220510 } },
//	{  "$project": {
//	      "date" : "$date",
//	      "date_integer": "$date_integer",
//	      "endpoints": {
//	        "$filter": {
//	          "input": "$endpoints",
//	          "as": "endpoint",
//	          "cond": { "$and": [ { "$eq": ["$$endpoint.severity","warning"] },
//	                              { "$eq": ["$$endpoint.classification", "outage" ] } ] }
//	      }
//	    }
//	  }
//	}
//
// ])
func prepFilterQuery(dt int, sev string, cl string) interface{} {

	// prepare the match part of the aggergation
	matchQuery := bson.M{"date_integer": dt}

	// if both severity and classification filter strings are empty proceed with the simple query
	if sev == "" && cl == "" {
		return []bson.M{{
			"$match": matchQuery,
		}}
	}

	// prepare the filter condition
	andParts := []bson.M{}
	// if severity filter has value add it to the condition
	if sev != "" {
		andParts = append(andParts, bson.M{"$eq": []string{"$$endpoint.severity", sev}})
	}
	// if classification filter has value add it to the condition
	if cl != "" {
		andParts = append(andParts, bson.M{"$eq": []string{"$$endpoint.classification", cl}})
	}
	// prepare the condition
	cond := bson.M{"$and": andParts}

	// prepare the projection query
	projectQuery := bson.M{
		"date":         "$date",
		"date_integer": "$date_integer",
		"endpoints": bson.M{
			"$filter": bson.M{
				"input": "$endpoints",
				"as":    "endpoint",
				"cond":  cond,
			},
		},
	}

	// return the whole aggregation with match and project parts
	return []bson.M{
		{"$match": matchQuery},
		{"$project": projectQuery},
	}

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
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)
	urlValues := r.URL.Query()
	dateStr := urlValues.Get("date")
	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	incoming := Downtimes{}
	incoming.DateInt = dt
	incoming.Date = dateStr

	// Try ingest request body
	body, err := io.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// check if topology already exists for current day

	existing := Downtimes{}
	downtimeCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(downtimeColName)
	err = downtimeCol.FindOne(context.TODO(), bson.M{"date_integer": dt}).Decode(&existing)
	if err != nil {
		// Stop at any error except not found. We want to have not found
		if err != mongo.ErrNoDocuments {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
		// else continue correctly -
	} else {
		// If found we need to inform user that the topology is already created for this date
		output, err = createMsgView(fmt.Sprintf("Downtimes already exists for date: %s, please either update it or delete it first!", dateStr), 409)
		code = 409
		return code, h, output, err
	}

	// Parse body json
	if err := json.Unmarshal(body, &incoming); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestInvalidJSON, contentType, "", " ")
		code = 400
		return code, h, output, err
	}

	_, err = downtimeCol.InsertOne(context.TODO(), incoming)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createMsgView(fmt.Sprintf("Downtimes set created for date: %s", dateStr), 201) //Render the results into JSON
	code = 201
	return code, h, output, err
}

// List the existing downtime resource for the exact date requested
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
	sev := urlValues.Get("severity")
	cl := urlValues.Get("classification")

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	dCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(downtimeColName)

	results := []Downtimes{}
	query := prepFilterQuery(dt, sev, cl)

	cursor, err := dCol.Aggregate(context.TODO(), query)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &results)

	if !(len(results) > 0) {
		downtimes := Downtimes{Date: dateStr, Endpoints: []Downtime{}}
		results = append(results, downtimes)
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

// delete downtimes resource based on ID
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

	urlValues := r.URL.Query()

	dateStr := urlValues.Get("date")
	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	if err != nil {
		output, _ = respond.MarshalContent(respond.UnauthorizedMessage, contentType, "", " ")
		code = http.StatusUnauthorized //If wrong api key is passed we return UNAUTHORIZED http status
		return code, h, output, err
	}

	dCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(downtimeColName)
	deleteResult, err := dCol.DeleteOne(context.TODO(), bson.M{"date_integer": dt})
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if deleteResult.DeletedCount == 0 {
		output, err = createMsgView(fmt.Sprintf("Downtimes dataset not found for date: %s", dateStr), 404)
		code = http.StatusNotFound
		return code, h, output, err
	}

	// Create view of the results
	output, err = createMsgView(fmt.Sprintf("Downtimes set deleted for date: %s", dateStr), 200)
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
