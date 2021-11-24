package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ARGOeu/argo-web-api/app/metricProfiles"
	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

//Create a new metrics resource
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

	incoming := []Metric{}

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

	_, err = mongo.Remove(coreSession, cfg.MongoDB.Db, metricsColName, bson.M{})
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	err = mongo.MultiInsert(coreSession, cfg.MongoDB.Db, metricsColName, incoming)

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

	// open a master session to the argo core database
	coreSession, err := mongo.OpenSession(cfg.MongoDB)
	defer mongo.CloseSession(coreSession)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Retrieve Results from database
	result := []Metric{}

	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	mCol := coreSession.DB(cfg.MongoDB.Db).C(metricsColName)
	err = mCol.Find(bson.M{}).All(&result)
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
	output, err = createMetricsListView(result, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

//ListMetricsByReport list all metrics available in the metric profile used by a report
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
	tenantDbConfig := context.Get(r, "tenant_conf").(config.MongoConfig)

	// Open session to tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	colReports := session.DB(tenantDbConfig.Db).C(reportsColName)
	//get the report
	report := reports.MongoInterface{}
	err = colReports.Find(bson.M{"info.name": reportName}).One(&report)

	if err != nil {
		if err.Error() == "not found" {
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

	colMProfiles := session.DB(tenantDbConfig.Db).C(mprofilesColName)
	//get the report
	mprofile := metricProfiles.MetricProfile{}
	err = colMProfiles.Find(mpQuery).One(&mprofile)

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

	// open a master session to the argo core database
	coreSession, err := mongo.OpenSession(cfg.MongoDB)
	defer mongo.CloseSession(coreSession)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Retrieve Results from database
	result := []Metric{}

	if err != nil {
		code = http.StatusBadRequest
		output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
		return code, h, output, err
	}

	mCol := coreSession.DB(cfg.MongoDB.Db).C(metricsColName)
	err = mCol.Find(bson.M{"name": bson.M{"$in": metrics}}).All(&result)

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
	output, err = createMetricsListView(result, "Success", code) //Render the results into JSON

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
