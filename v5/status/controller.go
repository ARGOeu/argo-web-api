package statusV5

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/utils/config"
	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// parseZuluDate is used to parse a zulu formatted date to integer
func parseZuluDate(dateStr string) (int, error) {
	parsedTime, _ := time.Parse(zuluForm, dateStr)
	return strconv.Atoi(parsedTime.Format(ymdForm))
}

// getPrevDay returns the previous day
func getPrevDay(dateStr string) (int, error) {
	parsedTime, _ := time.Parse(zuluForm, dateStr)
	prevTime := parsedTime.AddDate(0, 0, -1)
	return strconv.Atoi(prevTime.Format(ymdForm))
}

// GetGroupStatus returns a list of group status timelines
func GetGroupStatus(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("List Metric Timelines")
	err := error(nil)
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	parsedStart, _ := parseZuluDate(urlValues.Get("start-time"))

	parsedEnd, _ := parseZuluDate(urlValues.Get("end-time"))

	input := GroupInputParams{
		parsedStart,
		parsedEnd,
		vars["report-name"],
		vars["group-name"],
		"application/json",
	}

	// This is going to be used to determine a detailed view or not of the results
	view := urlValues.Get("view")
	details := false
	if view == "details" {
		details = true
	}

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	// Mongo Session
	results := []GroupDataOutput{}
	statusGroupCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("status_endpoint_groups")

	report := reports.MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	err = rCol.FindOne(context.TODO(), bson.M{"info.name": vars["report-name"]}).Decode(&report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report-name"] + " does not exist"
		output, err := createErrorMessage(message, code, contentType)
		return code, h, output, err
	}
	opts := options.Find().SetSort(bson.D{
		{Key: "endpoint_group", Value: 1}, // Primary sort: endpoint (A-Z)
		{Key: "date_integer", Value: 1},   // Secondary sort: timestamp (Oldest first)
	})
	cursor, err := statusGroupCol.Find(context.TODO(), groupQuery(input, report.ID), opts)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	parsedPrev, _ := getPrevDay(urlValues.Get("start-time"))

	//if no status results yet show previous days results
	if len(results) == 0 {
		// Zero query results
		input.startTime = parsedPrev
		cursor, err := statusGroupCol.Find(context.TODO(), groupQuery(input, report.ID), opts)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &results)

	}

	output, err = createGroupView(report, results, input, urlValues.Get("end-time"), details) //Render the results into JSON

	return code, h, output, err
}

// GetServiceStatus
func GetServiceStatus(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("List Service Timelines")
	err := error(nil)
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	parsedStart, _ := parseZuluDate(urlValues.Get("start-time"))
	parsedEnd, _ := parseZuluDate(urlValues.Get("end-time"))

	input := ServiceInputParams{
		parsedStart,
		parsedEnd,
		vars["report-name"],
		vars["group-name"],
		vars["service-type-name"],
		"application/json",
	}

	// This is going to be used to determine a detailed view or not of the results
	view := urlValues.Get("view")
	details := false
	if view == "details" {
		details = true
	}

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	report := reports.MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	err = rCol.FindOne(context.TODO(), bson.M{"info.name": vars["report-name"]}).Decode(&report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report-name"] + " does not exist"
		output, err := createErrorMessage(message, code, contentType)
		return code, h, output, err
	}

	// Mongo Session
	results := []ServiceDataOutput{}

	statusServiceCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("status_services")

	opts := options.Find().SetSort(bson.D{
		{Key: "endpoint_group", Value: 1},
		{Key: "service", Value: 1},
		{Key: "date_integer", Value: 1},
	})

	cursor, err := statusServiceCol.Find(context.TODO(), serviceQuery(input, report.ID), opts)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	parsedPrev, _ := getPrevDay(urlValues.Get("start-time"))

	//if no status results yet show previous days results
	if len(results) == 0 {
		// Zero query results
		input.startTime = parsedPrev
		cursor, err := statusServiceCol.Find(context.TODO(), serviceQuery(input, report.ID), opts)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &results)
	}

	output, err = createServiceView(report, results, input, urlValues.Get("end-time"), details) //Render the results into JSON

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// GetEndpointStatus returns a list of status timelines
func GetEndpointStatus(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("List Endpoint Timelines")
	err := error(nil)
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	parsedStart, _ := parseZuluDate(urlValues.Get("start-time"))
	parsedEnd, _ := parseZuluDate(urlValues.Get("end-time"))

	input := EndpointInputParams{
		parsedStart,
		parsedEnd,
		vars["report-name"],
		vars["group-name"],
		vars["service-type-name"],
		vars["endpoint-name"],
		contentType,
	}

	// This is going to be used to determine a detailed view or not of the results
	view := urlValues.Get("view")
	details := false
	if view == "details" {
		details = true
	}

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	report := reports.MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	err = rCol.FindOne(context.TODO(), bson.M{"info.name": vars["report-name"]}).Decode(&report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report-name"] + " does not exist"
		output, err := createErrorMessage(message, code, contentType)
		return code, h, output, err
	}

	// Mongo Session
	results := []EndpointDataOutput{}

	statusEndpointCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("status_endpoints")

	opts := options.Find().SetSort(bson.D{
		{Key: "endpoint_group", Value: 1},
		{Key: "service", Value: 1},
		{Key: "host", Value: 1},
		{Key: "date_integer", Value: 1},
	})

	cursor, err := statusEndpointCol.Find(context.TODO(), endpointQuery(input, report.ID), opts)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	parsedPrev, _ := getPrevDay(urlValues.Get("start-time"))

	//if no status results yet show previous days results
	if len(results) == 0 {
		// Zero query results
		input.startTime = parsedPrev
		cursor, err := statusEndpointCol.Find(context.TODO(), endpointQuery(input, report.ID))
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &results)

	}

	output, err = createEndpointView(report, results, input, urlValues.Get("end-time"), details) //Render the results into JSON

	return code, h, output, err
}

// GetMetricStatus returns a list of status timelines
func GetMetricStatus(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("List Metric Timelines")
	err := error(nil)
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Set Content-Type response Header value
	contentType := r.Header.Get("Accept")
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	parsedStart, _ := parseZuluDate(urlValues.Get("start-time"))
	parsedEnd, _ := parseZuluDate(urlValues.Get("end-time"))

	view := urlValues.Get("view")
	details := false

	if view == "details" {
		details = true
	}

	input := MetricInputParams{
		parsedStart,
		parsedEnd,
		vars["report-name"],
		vars["group-name"],
		vars["service-type-name"],
		vars["endpoint-name"],
		vars["metric-name"],
		contentType,
	}

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	// Mongo Session
	results := []MetricDataOutput{}

	metricCollection := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("status_metrics")

	report := reports.MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")
	err = rCol.FindOne(context.TODO(), bson.M{"info.name": vars["report-name"]}).Decode(&report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the name " + vars["report-name"] + " does not exist"
		output, err := createErrorMessage(message, code, contentType)
		return code, h, output, err
	}

	opts := options.Find().SetSort(bson.D{
		{Key: "endpoint_group", Value: 1},
		{Key: "service", Value: 1},
		{Key: "host", Value: 1},
		{Key: "metric", Value: 1},
		{Key: "date_integer", Value: 1},
	})

	cursor, err := metricCollection.Find(context.TODO(), metricQuery(input, report.ID), opts)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &results)

	parsedPrev, _ := getPrevDay(urlValues.Get("start_time"))

	//if no status results yet show previous days results
	if len(results) == 0 {
		// Zero query results
		input.startTime = parsedPrev
		cursor, err := metricCollection.Find(context.TODO(), metricQuery(input, report.ID), opts)
		if err != nil {
			code = http.StatusInternalServerError
			return code, h, output, err
		}

		defer cursor.Close(context.TODO())
		cursor.All(context.TODO(), &results)
	}

	output, err = createMetricView(report, results, input, details) //Render the results into XML format

	return code, h, output, err
}

func metricQuery(input MetricInputParams, reportID string) bson.M {

	// prepare the match filter
	filter := bson.M{
		"report":         reportID,
		"date_integer":   bson.M{"$gte": input.startTime, "$lte": input.endTime},
		"endpoint_group": input.group,
	}

	if len(input.service) > 0 {
		filter["service"] = input.service
	}

	filter["host"] = input.hostname

	if len(input.metric) > 0 {
		filter["metric"] = input.metric
	}

	return filter
}

func endpointQuery(input EndpointInputParams, reportID string) bson.M {

	// prepare the match filter
	filter := bson.M{
		"date_integer": bson.M{"$gte": input.startTime, "$lte": input.endTime},
		"report":       reportID,
	}

	if len(input.group) > 0 {
		filter["endpoint_group"] = input.group
	}

	if len(input.service) > 0 {
		filter["service"] = input.service
	}

	if len(input.hostname) > 0 {
		filter["host"] = input.hostname
	}

	return filter
}

func serviceQuery(input ServiceInputParams, reportID string) bson.M {

	// prepare the match filter
	filter := bson.M{
		"date_integer":   bson.M{"$gte": input.startTime, "$lte": input.endTime},
		"report":         reportID,
		"endpoint_group": input.group,
	}

	if len(input.service) > 0 {
		filter["service"] = input.service
	}

	return filter
}

func groupQuery(input GroupInputParams, reportID string) bson.M {
	filter := bson.M{
		"date_integer": bson.M{"$gte": input.startTime, "$lte": input.endTime},
		"report":       reportID,
	}

	if len(input.group) > 0 {
		filter["endpoint_group"] = input.group
	}

	return filter
}

// Options responds to an OPTIONS request
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

// GetMetricDetails returns the detailed message from a probe
func GetMetricDetails(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	reportName := urlValues.Get("report")
	reportID := ""

	reportCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("reports")

	if reportName != "" {

		queryResult := reportCol.FindOne(context.TODO(), bson.M{"info.name": reportName})

		if queryResult.Err() != nil {
			if queryResult.Err() == mongo.ErrNoDocuments {
				code = http.StatusNotFound
				message := "The report with the name " + reportName + " does not exist"
				output, err := createErrorMessage(message, code, contentType) //Render the response into XML or JSON
				h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
				return code, h, output, err
			}
			code = http.StatusInternalServerError
			return code, h, output, err
		}

	}

	input := metricDetailsQuery{
		GroupName:       vars["group-name"],
		ServiceTypeName: vars["service-type-name"],
		EndpointName:    vars["endpoint-name"],
		MetricName:      vars["metric-name"],
		ExecTime:        urlValues.Get("timestamp"),
	}

	var result metricResultOutput

	metricCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection("status_metrics")
	q := prepQuery(input, reportID)
	fmt.Println(q)
	err = metricCol.FindOne(context.TODO(), q).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			code = http.StatusNotFound
			message := "Metric not found!"
			output, err := createErrorMessage(message, code, contentType)
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = json.MarshalIndent(result, " ", "  ")

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

func prepQuery(input metricDetailsQuery, reportID string) bson.M {

	//Time Related
	const zuluForm = "2006-01-02T15:04:05Z"
	const ymdForm = "20060102"

	ts, _ := time.Parse(zuluForm, input.ExecTime)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))

	// parse time as integer
	tsInt := (ts.Hour() * 10000) + (ts.Minute() * 100) + ts.Second()

	query := bson.M{
		"date_integer":   tsYMD,
		"endpoint_group": input.GroupName,
		"service":        input.ServiceTypeName,
		"host":           input.EndpointName,
		"metric":         input.MetricName,
		"time_integer":   tsInt,
	}

	if reportID != "" {
		query["report"] = reportID
	}

	return query

}
