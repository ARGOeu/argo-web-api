package statusEndpoints

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ARGOeu/argo-web-api/utils/authentication"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/mux"
	"labix.org/v2/mgo/bson"
)

// ListMetricTimelines returns a list of metric timelines
func ListEndpointTimelines(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("List Endpoint Timelines")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Parse the request into the input
	urlValues := r.URL.Query()
	vars := mux.Vars(r)

	input := InputParams{
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		vars["report_name"],
		vars["group_type"],
		vars["group_name"],
		vars["service_name"],
		vars["endpoint_name"],
	}

	// Call authenticateTenant to check the api key and retrieve
	// the correct tenant db conf
	tenantDbConfig, err := authentication.AuthenticateTenant(r.Header, cfg)

	if err != nil {
		output = []byte(http.StatusText(http.StatusUnauthorized))
		code = http.StatusUnauthorized //If wrong api key is passed we return UNAUTHORIZED http status
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	// Mongo Session
	results := []DataOutput{}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	metricCollection := session.DB(tenantDbConfig.Db).C("status_endpoints")

	// Query the detailed metric results
	err = metricCollection.Find(prepareQuery(input)).All(&results)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createView(results, input) //Render the results into XML format

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

func prepareQuery(input InputParams) bson.M {

	//Time Related
	const zuluForm = "2006-01-02T15:04:05Z"
	const ymdForm = "20060102"

	ts, _ := time.Parse(zuluForm, input.startTime)
	te, _ := time.Parse(zuluForm, input.endTime)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	// prepare the match filter
	filter := bson.M{
		"date_int":       bson.M{"$gte": tsYMD, "$lte": teYMD},
		"report":         input.report,
		"endpoint_group": input.group,
		"service":        input.service,
	}

	if len(input.hostname) > 0 {
		filter["hostname"] = input.hostname
	}

	return filter
}
