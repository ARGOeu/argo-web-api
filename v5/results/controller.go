package resultsV5

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ListGroupResults list monitoring results (availability/reliability)
func GetGroups(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	report := reports.MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	err = rCol.FindOne(context.TODO(), bson.M{"info.name": vars["report-name"]}).Decode(&report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report-name"] + " does not exist"
		output, err := createErrorMessage(message, code, contentType) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	input := endpointGroupResultQuery{
		basicQuery{
			Name:        vars["group-name"],
			Granularity: urlValues.Get("granularity"),
			Format:      contentType,
			StartTime:   urlValues.Get("start-time"),
			EndTime:     urlValues.Get("end-time"),
			Report:      report,
			Vars:        vars,
		}, "",
	}

	errs := input.Validate()
	if len(errs) > 0 {
		out := respond.BadRequestSimple
		out.Errors = errs
		output = out.MarshalTo(contentType)
		code = 400
		return code, h, output, err
	}

	results := []EndpointGroupInterface{}

	// Construct the query to mongodb based on the input
	filter := bson.M{
		"date":   bson.M{"$gte": input.StartTimeInt, "$lte": input.EndTimeInt},
		"report": report.ID,
	}

	if input.Name != "" {
		filter["name"] = input.Name
	}

	// Select the granularity of the search daily/monthly
	custom := false
	arCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("endpoint_group_ar")
	var query []primitive.M

	if input.Granularity == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		query = DailyGroup(filter)
	} else if input.Granularity == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = MonthlyGroup(filter)
	} else if input.Granularity == "custom" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = CustomGroup(filter)
		custom = true
	}

	cursor, err := arCol.Aggregate(context.TODO(), query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	if len(results) == 0 {
		code = http.StatusNotFound
		message := "No results found for given query"
		output, err = createErrorMessage(message, code, contentType)
		return code, h, output, err
	}

	output, err = createGroupResultView(results, report, input.Format, custom)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

// ListSuperGroupResults supergroup availabilities according to the http request
func GetSupergroups(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	report := reports.MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	err = rCol.FindOne(context.TODO(), bson.M{"info.name": vars["report-name"]}).Decode(&report)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			code = http.StatusNotFound
			message := "The report with the name " + vars["report-name"] + " does not exist"
			output, err := createErrorMessage(message, code, contentType) //Render the response into XML or JSON
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		} else {
			code = http.StatusInternalServerError
			return code, h, output, err
		}
	}

	input := endpointGroupResultQuery{
		basicQuery{
			Name:        vars["supergroup-name"],
			Granularity: urlValues.Get("granularity"),
			Format:      contentType,
			StartTime:   urlValues.Get("start-time"),
			EndTime:     urlValues.Get("end-time"),
			Report:      report,
			Vars:        vars,
		}, "",
	}

	errs := input.Validate()
	if len(errs) > 0 {
		out := respond.BadRequestSimple
		out.Errors = errs
		output = out.MarshalTo(contentType)
		code = 400
		return code, h, output, err
	}

	results := []SuperGroupInterface{}

	// Construct the query to mongodb based on the input
	filter := bson.M{
		"date":   bson.M{"$gte": input.StartTimeInt, "$lte": input.EndTimeInt},
		"report": report.ID,
	}

	if input.Name != "" {
		filter["supergroup"] = input.Name
	}

	custom := false

	arCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("endpoint_group_ar")
	var query []primitive.M
	// Select the granularity of the search daily/monthly
	if input.Granularity == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		query = DailySuperGroup(filter)
	} else if input.Granularity == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = MonthlySuperGroup(filter)
	} else if input.Granularity == "custom" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = CustomSuperGroup(filter)
		custom = true
	}
	cursor, err := arCol.Aggregate(context.TODO(), query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	if len(results) == 0 {
		code = http.StatusNotFound
		message := "No results found for given query"
		output, err = createErrorMessage(message, code, contentType)
		return code, h, output, err
	}

	output, err = createSuperGroupView(results, report, input.Format, custom)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

// GetEndpoints is responsible to return results for endpoints
func GetEndpoints(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
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

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	report := reports.MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	err = rCol.FindOne(context.TODO(), bson.M{"info.name": vars["report-name"]}).Decode(&report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report_name"] + " does not exist"
		output, err := createErrorMessage(message, code, contentType) //Render the response into XML or JSON
		return code, h, output, err
	}

	input := endpointResultQuery{
		basicQuery: basicQuery{
			Name: vars["endpoint-name"],

			Granularity: urlValues.Get("granularity"),
			Format:      contentType,
			StartTime:   urlValues.Get("start-time"),
			EndTime:     urlValues.Get("end-time"),
			Report:      report,
			Vars:        vars,
		},
		EndpointGroup: vars["group-name"],
		Service:       vars["service-type"],
	}

	errs := input.Validate()
	if len(errs) > 0 {
		out := respond.BadRequestSimple
		out.Errors = errs
		output = out.MarshalTo(contentType)
		code = 400
		return code, h, output, err
	}

	results := []EndpointInterface{}

	// Construct the query to mongodb based on the input
	filter := bson.M{
		"date":   bson.M{"$gte": input.StartTimeInt, "$lte": input.EndTimeInt},
		"report": report.ID,
	}

	if input.Name != "" {
		filter["name"] = input.Name
	}

	if input.EndpointGroup != "" {
		filter["supergroup"] = input.EndpointGroup
	}

	if input.Service != "" {
		filter["service"] = input.Service
	}

	// Select the granularity of the search daily/monthly
	custom := false
	arCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("endpoint_ar")
	var query []primitive.M
	if input.Granularity == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		query = DailyEndpoint(filter)

	} else if input.Granularity == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = MonthlyEndpoint(filter)

	} else if input.Granularity == "custom" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = CustomEndpoint(filter)
		custom = true
	}

	cursor, err := arCol.Aggregate(context.TODO(), query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	if len(results) == 0 {
		code = http.StatusNotFound
		message := "No results found for given query"
		output, err = createErrorMessage(message, code, contentType)
		return code, h, output, err
	}

	output, err = createEndpointResultView(results, report, input.Format, custom)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err

}

// GetServiceTypes is responsible for handling request to list service flavor results
func GetServiceTypes(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
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

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	report := reports.MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	err = rCol.FindOne(context.TODO(), bson.M{"info.name": vars["report-name"]}).Decode(&report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report_name"] + " does not exist"
		output, err := createErrorMessage(message, code, contentType) //Render the response into XML or JSON
		return code, h, output, err
	}

	input := serviceFlavorResultQuery{
		basicQuery: basicQuery{
			Name:        vars["service-type-name"],
			Granularity: urlValues.Get("granularity"),
			Format:      contentType,
			StartTime:   urlValues.Get("start-time"),
			EndTime:     urlValues.Get("end-time"),
			Report:      report,
			Vars:        vars,
		},
		EndpointGroup: vars["group-name"],
	}

	errs := input.Validate()
	if len(errs) > 0 {
		out := respond.BadRequestSimple
		out.Errors = errs
		output = out.MarshalTo(contentType)
		code = 400
		return code, h, output, err
	}

	results := []ServiceFlavorInterface{}

	// Construct the query to mongodb based on the input
	filter := bson.M{
		"date":   bson.M{"$gte": input.StartTimeInt, "$lte": input.EndTimeInt},
		"report": report.ID,
	}

	if input.Name != "" {
		filter["name"] = input.Name
	}

	if input.EndpointGroup != "" {
		filter["supergroup"] = input.EndpointGroup
	}

	// Select the granularity of the search daily/monthly
	custom := false
	arCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("service_ar")
	var query []primitive.M

	if input.Granularity == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		query = DailyServiceFlavor(filter)
	} else if input.Granularity == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = MonthlyServiceFlavor(filter)
	} else if input.Granularity == "custom" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = CustomServiceFlavor(filter)
		custom = true
	}

	cursor, err := arCol.Aggregate(context.TODO(), query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	if len(results) == 0 {
		code = http.StatusNotFound
		message := "No results found for given query"
		output, err = createErrorMessage(message, code, contentType)
		return code, h, output, err
	}

	output, err = createServiceFlavorResultView(results, report, input.Format, custom)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err

}

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
	h.Set("Allow", "GET, OPTIONS")
	return code, h, output, err

}
