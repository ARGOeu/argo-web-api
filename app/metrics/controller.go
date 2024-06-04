package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ARGOeu/argo-web-api/app/metricProfiles"
	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils"
	"github.com/ARGOeu/argo-web-api/utils/config"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create a new metrics resource
const metricsColName = "monitoring_metrics"
const reportsColName = "reports"
const mprofilesColName = "metric_profiles"

// gets a report result and retuns the metric profile ID
func getReportMetricProfileID(r reports.MongoInterface) string {
	for _, item := range r.Profiles {
		if item.Type == "metric" {
			return item.ID
		}
	}
	return ""
}

// gets a list with metric names from a metric profile
func getMetricsFromProfile(mp metricProfiles.MetricProfile) []string {
	set := make(map[string]bool)
	result := []string{}

	for _, service := range mp.Services {
		for _, metric := range service.Metrics {
			set[metric] = true
		}
	}

	for key := range set {
		result = append(result, key)
	}

	return result
}

func prepQueryMprofile(dt int, id string) interface{} {

	return bson.M{"date_integer": bson.M{"$lte": dt}, "id": id}

}

// Update request handler creates a new list of metrics
func UpdateMetrics(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
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

	incoming := []Metric{}

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

	coreCol := cfg.MongoClient.Database(cfg.MongoDB.Db).Collection(metricsColName)

	_, err = coreCol.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	incomingInf := make([]interface{}, len(incoming))
	for i, value := range incoming {
		incomingInf[i] = value
	}

	_, err = coreCol.InsertMany(context.TODO(), incomingInf)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createMetricsListView(incoming, "Metrics resource succesfully updated", 200) //Render the results into JSON
	code = 200
	return code, h, output, err
}

// ListMetrics actually list monitoring metrics
func ListMetrics(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	// Retrieve Results from database
	results := []Metric{}

	coreCol := cfg.MongoClient.Database(cfg.MongoDB.Db).Collection(metricsColName)
	cursor, err := coreCol.Find(context.TODO(), bson.M{})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
			code = 404
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	// Create view of the results
	output, err = createMetricsListView(results, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// ListMetricsByReport list all metrics available in the metric profile used by a report
func ListMetricsByReport(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	reportName := vars["report"]

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	colReports := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(reportsColName)
	//get the report
	report := reports.MongoInterface{}
	err = colReports.FindOne(context.TODO(), bson.M{"info.name": reportName}).Decode(&report)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			output, err = createMessageOUT(fmt.Sprintf("No report with name: %s exists!", reportName), 404, "json")
			code = 404
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	mprofileID := getReportMetricProfileID(report)

	if mprofileID == "" {
		output, err = createMessageOUT("Report doesn't contain a metric profile", 404, "json")
		code = 404
		return code, h, output, err
	}

	dt, dateStr, err := utils.ParseZuluDate(dateStr)
	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	mpQuery := prepQueryMprofile(dt, mprofileID)

	colMProfiles := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(mprofilesColName)
	//get the report
	mprofile := metricProfiles.MetricProfile{}
	err = colMProfiles.FindOne(context.TODO(), mpQuery).Decode(&mprofile)

	if err != nil {
		if err.Error() == "not found" {
			output, err = createMessageOUT(fmt.Sprintf("No metric profiles with ID: %s exists!", mprofileID), 404, "json")
			code = 404
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	metrics := getMetricsFromProfile(mprofile)

	// Retrieve Results from database
	results := []Metric{}

	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	mCol := cfg.MongoClient.Database(cfg.MongoDB.Db).Collection(metricsColName)
	cursor, err := mCol.Find(context.TODO(), bson.M{"name": bson.M{"$in": metrics}})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			output, _ = respond.MarshalContent(respond.ErrNotFound, contentType, "", " ")
			code = 404
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	// Create view of the results
	output, err = createMetricsListView(results, "Success", code) //Render the results into JSON

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
