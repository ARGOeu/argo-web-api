package nodes

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetSummary lists availability, reliability and uptime for a specific group
func GetSummary(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	item := vars["item"]

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	// get node report id
	nodeReport := NodeReport{}

	report := reports.MongoInterface{}
	nrCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(nodeReportColName)
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(reportsColName)

	err = nrCol.FindOne(context.TODO(), bson.M{"_id": "node-report"}).Decode(&nodeReport)

	if err != nil {
		code = http.StatusNotFound
		message := "Node report not set"
		output, err := createErrorMessage(message, code, contentType)
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	err = rCol.FindOne(context.TODO(), bson.M{"id": nodeReport.ReportId}).Decode(&report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the id " + nodeReport.ReportId + " does not exist"
		output, err := createErrorMessage(message, code, contentType)
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	input :=
		basicQuery{
			Name:        item,
			Granularity: urlValues.Get("granularity"),
			Format:      contentType,
			StartTime:   urlValues.Get("start_time"),
			EndTime:     urlValues.Get("end_time"),
			Report:      report,
			Vars:        vars,
			Date:        urlValues.Get("date"),
			StartDate:   urlValues.Get("start_date"),
			EndDate:     urlValues.Get("end_date"),
		}
	errs := input.Validate()
	if len(errs) > 0 {
		out := respond.BadRequestSimple
		out.Errors = errs
		output = out.MarshalTo(contentType)
		code = 400
		return code, h, output, err
	}

	resultsGroups := []GroupInterface{}

	filter := bson.M{
		"date":   bson.M{"$gte": input.StartTimeInt, "$lte": input.EndTimeInt},
		"report": report.ID,
	}
	var query []primitive.M

	if input.Name != "" {
		filter["name"] = input.Name
	}

	// Prepare the group results
	// Select the granularity of the search daily/monthly

	if input.Granularity == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		query = DailyEndpoint(filter)

	} else if input.Granularity == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = MonthlyGroup(filter)

	}

	arCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(groupColName)
	cursor, err := arCol.Aggregate(context.TODO(), query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &resultsGroups)

	output, err = createSummaryView(resultsGroups)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

// GetAvailability lists group availabilities according to the http request
func GetAvailability(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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

	item := vars["item"]

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	// get node report id
	nodeReport := NodeReport{}

	report := reports.MongoInterface{}
	nrCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(nodeReportColName)
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(reportsColName)

	err = nrCol.FindOne(context.TODO(), bson.M{"_id": "node-report"}).Decode(&nodeReport)

	if err != nil {
		code = http.StatusNotFound
		message := "Node report not set"
		output, err := createErrorMessage(message, code, contentType)
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	err = rCol.FindOne(context.TODO(), bson.M{"id": nodeReport.ReportId}).Decode(&report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the id " + nodeReport.ReportId + " does not exist"
		output, err := createErrorMessage(message, code, contentType)
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	input :=
		basicQuery{
			Name:        item,
			Granularity: urlValues.Get("granularity"),
			Format:      contentType,
			StartTime:   urlValues.Get("start_time"),
			EndTime:     urlValues.Get("end_time"),
			Report:      report,
			Vars:        vars,
			Date:        urlValues.Get("date"),
			StartDate:   urlValues.Get("start_date"),
			EndDate:     urlValues.Get("end_date"),
		}
	errs := input.Validate()
	if len(errs) > 0 {
		out := respond.BadRequestSimple
		out.Errors = errs
		output = out.MarshalTo(contentType)
		code = 400
		return code, h, output, err
	}

	resultsGroups := []GroupInterface{}

	filter := bson.M{
		"date":   bson.M{"$gte": input.StartTimeInt, "$lte": input.EndTimeInt},
		"report": report.ID,
	}

	if input.Name != "" {
		filter["name"] = item
	}

	var query []primitive.M

	// Prepare the group results
	// Select the granularity of the search daily/monthly

	if input.Granularity == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		query = DailyEndpoint(filter)

	} else if input.Granularity == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = MonthlyGroup(filter)

	}

	arCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(groupColName)
	cursor, err := arCol.Aggregate(context.TODO(), query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &resultsGroups)

	output, err = createAvailabilityView(resultsGroups)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

// GetUptime lists group availabilities according to the http request
func GetUptime(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	item := vars["item"]

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	// get node report id
	nodeReport := NodeReport{}

	report := reports.MongoInterface{}
	nrCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(nodeReportColName)
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(reportsColName)

	err = nrCol.FindOne(context.TODO(), bson.M{"_id": "node-report"}).Decode(&nodeReport)

	if err != nil {
		code = http.StatusNotFound
		message := "Node report not set"
		output, err := createErrorMessage(message, code, contentType)
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	err = rCol.FindOne(context.TODO(), bson.M{"id": nodeReport.ReportId}).Decode(&report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the id " + nodeReport.ReportId + " does not exist"
		output, err := createErrorMessage(message, code, contentType)
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	input :=
		basicQuery{
			Name:        item,
			Granularity: urlValues.Get("granularity"),
			Format:      contentType,
			StartTime:   urlValues.Get("start_time"),
			EndTime:     urlValues.Get("end_time"),
			Report:      report,
			Vars:        vars,
			Date:        urlValues.Get("date"),
			StartDate:   urlValues.Get("start_date"),
			EndDate:     urlValues.Get("end_date"),
		}
	errs := input.Validate()
	if len(errs) > 0 {
		out := respond.BadRequestSimple
		out.Errors = errs
		output = out.MarshalTo(contentType)
		code = 400
		return code, h, output, err
	}

	resultsGroups := []GroupInterface{}

	filter := bson.M{
		"date":   bson.M{"$gte": input.StartTimeInt, "$lte": input.EndTimeInt},
		"report": report.ID,
	}
	var query []primitive.M

	if input.Name != "" {
		filter["name"] = input.Name
	}

	// Prepare the group results
	// Select the granularity of the search daily/monthly

	if input.Granularity == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		query = DailyEndpoint(filter)

	} else if input.Granularity == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = MonthlyGroup(filter)

	}

	arCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(groupColName)
	cursor, err := arCol.Aggregate(context.TODO(), query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &resultsGroups)

	output, err = createUptimeView(resultsGroups)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

// GetGroupResults lists availability, reliability and uptime for a specific group
func GetGroupResults(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	reportItem := urlValues.Get("report")
	period := urlValues.Get("period")

	item := vars["item"]

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	report := reports.MongoInterface{}

	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(reportsColName)

	// if report param is given
	if reportItem != "" {
		err = rCol.FindOne(context.TODO(), bson.M{"$or": bson.A{bson.M{"id": reportItem}, bson.M{"info.name": reportItem}}}).Decode(&report)

		if err != nil {
			code = http.StatusNotFound
			message := "The report: " + reportItem + " does not exist"
			output, err := createErrorMessage(message, code, contentType)
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}
	} else {
		count, err := rCol.CountDocuments(context.TODO(), bson.M{})
		if err != nil {
			code = http.StatusInternalServerError
			message := "Could not access reports"
			output, err := createErrorMessage(message, code, contentType)
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}
		switch {
		case count == 0:
			code = http.StatusNotFound
			message := "No reports found"
			output, err := createErrorMessage(message, code, contentType)
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		case count > 1:
			code = http.StatusBadRequest
			message := "Multiple reports found please specify one by using report url param"
			output, err := createErrorMessage(message, code, contentType)
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}
		err = rCol.FindOne(context.TODO(), bson.M{}).Decode(&report)
	}

	input :=
		basicQuery{
			Name:        item,
			Granularity: urlValues.Get("granularity"),
			Format:      contentType,
			StartTime:   urlValues.Get("start_time"),
			EndTime:     urlValues.Get("end_time"),
			Report:      report,
			Vars:        vars,
			Date:        urlValues.Get("date"),
			StartDate:   urlValues.Get("start_date"),
			EndDate:     urlValues.Get("end_date"),
		}

	// check if period declared
	if period != "" {
		days, err := periodToDays(period)
		if err != nil {
			code = http.StatusBadRequest
			output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(err.Error()), contentType, "", " ")
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}
		input.EndTime = time.Now().UTC().Format(time.RFC3339)
		input.StartTime = time.Now().UTC().AddDate(0, 0, -(days - 1)).Format(time.RFC3339)

	}

	errs := input.Validate()
	if len(errs) > 0 {
		out := respond.BadRequestSimple
		out.Errors = errs
		output = out.MarshalTo(contentType)
		code = 400
		return code, h, output, err
	}

	resultsGroups := []GroupInterface{}

	filter := bson.M{
		"date":   bson.M{"$gte": input.StartTimeInt, "$lte": input.EndTimeInt},
		"report": report.ID,
	}

	var query []primitive.M

	if input.Name != "" {
		filter["name"] = input.Name
	}

	// Prepare the group results
	// Select the granularity of the search daily/monthly

	if input.Granularity == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		query = DailyEndpoint(filter)

	} else if input.Granularity == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query = MonthlyGroup(filter)

	}

	arCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(groupColName)
	cursor, err := arCol.Aggregate(context.TODO(), query)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &resultsGroups)

	output, err = createSummaryView(resultsGroups)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

// GetGroupStatus lists group timelines per report
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
	item := vars["item"]
	reportItem := urlValues.Get("report")

	// This is going to be used to determine a detailed/latest view of the results
	view := urlValues.Get("view")
	history := urlValues.Get("history")
	details := false
	if view == "details" {
		details = true
	}

	var parsedStart, parsedEnd int

	// check if user provided start and end time correctly
	urlStartTime := urlValues.Get("start_time")
	urlEndTime := urlValues.Get("end_time")

	endDate := ""

	if urlStartTime == "" && urlEndTime == "" {
		isoTimeNow := time.Now().UTC().Format(time.RFC3339)
		startDate := strings.Split(isoTimeNow, "T")[0]
		startTime := startDate + "T00:00:00Z"
		endDate = startDate
		parsedStart, _ = parseZuluDate(startTime)
		parsedEnd, _ = parseZuluDate(isoTimeNow)
	} else {
		if parsedStart, err = parseZuluDate(urlStartTime); err != nil {
			code = http.StatusBadRequest
			message := fmt.Sprintf("Error parsing start_time=%s - please use zulu format like %s", urlStartTime, zuluForm)
			output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(message), contentType, "", " ")
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}

		if parsedEnd, err = parseZuluDate(urlEndTime); err != nil {
			code = http.StatusBadRequest
			message := fmt.Sprintf("Error parsing end_time=%s - please use zulu format like %s", urlEndTime, zuluForm)
			output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(message), contentType, "", " ")
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}

		endDate = strings.Split(urlEndTime, "T")[0]
		history = "true"
	}

	input := InputStatus{
		parsedStart,
		parsedEnd,
		urlValues.Get("report_name"),
		urlValues.Get("group_type"),
		item,
		contentType,
		"",
	}

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	report := reports.MongoInterface{}
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(reportsColName)

	// if report param is given
	if reportItem != "" {
		err = rCol.FindOne(context.TODO(), bson.M{"$or": bson.A{bson.M{"id": reportItem}, bson.M{"info.name": reportItem}}}).Decode(&report)

		if err != nil {
			code = http.StatusNotFound
			message := "The report: " + reportItem + " does not exist"
			output, err := createErrorMessage(message, code, contentType)
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}
	} else {
		count, err := rCol.CountDocuments(context.TODO(), bson.M{})
		if err != nil {
			code = http.StatusInternalServerError
			message := "Could not access reports"
			output, err := createErrorMessage(message, code, contentType)
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}
		switch {
		case count == 0:
			code = http.StatusNotFound
			message := "No reports found"
			output, err := createErrorMessage(message, code, contentType)
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		case count > 1:
			code = http.StatusBadRequest
			message := "Multiple reports found please specify one by using report url param"
			output, err := createErrorMessage(message, code, contentType)
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}
		err = rCol.FindOne(context.TODO(), bson.M{}).Decode(&report)
	}

	input.groupType = report.Topology.Group.Group.Type

	groupCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(statusGroupColName)
	groupResults := []GroupStatusData{}

	cursor, err := groupCol.Find(context.TODO(), queryStatusGroups(input, report.ID))
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &groupResults)

	output, err = createStatusView(groupResults, input, endDate, details, history == "") //Render the results into JSON

	return code, h, output, err
}

// GetStatus lists group timelines
func GetStatus(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

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
	item := vars["item"]

	// This is going to be used to determine a detailed/latest view of the results
	view := urlValues.Get("view")
	history := urlValues.Get("history")
	details := false
	if view == "details" {
		details = true
	}

	var parsedStart, parsedEnd int

	// check if user provided start and end time correctly
	urlStartTime := urlValues.Get("start_time")
	urlEndTime := urlValues.Get("end_time")

	endDate := ""

	if urlStartTime == "" && urlEndTime == "" {
		isoTimeNow := time.Now().UTC().Format(time.RFC3339)
		startDate := strings.Split(isoTimeNow, "T")[0]
		startTime := startDate + "T00:00:00Z"
		endDate = startDate
		parsedStart, _ = parseZuluDate(startTime)
		parsedEnd, _ = parseZuluDate(isoTimeNow)
	} else {
		if parsedStart, err = parseZuluDate(urlStartTime); err != nil {
			code = http.StatusBadRequest
			message := fmt.Sprintf("Error parsing start_time=%s - please use zulu format like %s", urlStartTime, zuluForm)
			output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(message), contentType, "", " ")
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}

		if parsedEnd, err = parseZuluDate(urlEndTime); err != nil {
			code = http.StatusBadRequest
			message := fmt.Sprintf("Error parsing end_time=%s - please use zulu format like %s", urlEndTime, zuluForm)
			output, _ = respond.MarshalContent(respond.ErrBadRequestDetails(message), contentType, "", " ")
			h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
			return code, h, output, err
		}

		endDate = strings.Split(urlEndTime, "T")[0]
		history = "true"
	}

	input := InputStatus{
		parsedStart,
		parsedEnd,
		urlValues.Get("report_name"),
		urlValues.Get("group_type"),
		item,
		contentType,
		"",
	}

	// Grab Tenant DB configuration from context
	tenantDbConfig := gcontext.Get(r, "tenant_conf").(config.MongoConfig)

	nodeReport := NodeReport{}

	report := reports.MongoInterface{}
	nrCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(nodeReportColName)
	rCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(reportsColName)

	err = nrCol.FindOne(context.TODO(), bson.M{"_id": "node-report"}).Decode(&nodeReport)

	if err != nil {
		code = http.StatusNotFound
		message := "Node report not set"
		output, err := createErrorMessage(message, code, contentType)
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	err = rCol.FindOne(context.TODO(), bson.M{"id": nodeReport.ReportId}).Decode(&report)

	if err != nil {
		code = http.StatusNotFound
		message := "The report with the id " + nodeReport.ReportId + " does not exist"
		output, err := createErrorMessage(message, code, contentType)
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	input.groupType = report.Topology.Group.Group.Type

	groupCol := cfg.MongoClient.Database(tenantDbConfig.Db).Collection(statusGroupColName)
	groupResults := []GroupStatusData{}

	cursor, err := groupCol.Find(context.TODO(), queryStatusGroups(input, report.ID))
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &groupResults)

	output, err = createStatusView(groupResults, input, endDate, details, history == "") //Render the results into JSON

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
