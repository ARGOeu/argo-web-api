/*
 * Copyright (c) 2014 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the
 * License. You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an "AS
 * IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language
 * governing permissions and limitations under the License.
 *
 * The views and conclusions contained in the software and
 * documentation are those of the authors and should not be
 * interpreted as representing official policies, either expressed
 * or implied, of either GRNET S.A., SRCE or IN2P3 CNRS Computing
 * Centre
 *
 * The work represented by this source file is partially funded by
 * the EGI-InSPIRE project through the European Commission's 7th
 * Framework Programme (contract # INFSO-RI-261323)
 */

package endpointGroupAvailability

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ARGOeu/argo-web-api/utils/authentication"
	"github.com/ARGOeu/argo-web-api/utils/caches"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
)

// List endpoint group availabilities according to the http request
func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END
	tenantDbConfig, err := authentication.AuthenticateTenant(r.Header, cfg)
	if err != nil {
		if err.Error() == "Unauthorized" {
			code = http.StatusUnauthorized
			return code, h, output, err
		}
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Parse the request into the input
	urlValues := r.URL.Query()

	input := EndpointGroupAvailabilityInput{
		StartTime:   urlValues.Get("start_time"),
		EndTime:     urlValues.Get("end_time"),
		Granularity: urlValues.Get("granularity"),
		Format:      urlValues.Get("format"),
		Job:         urlValues.Get("job"),
		GroupName:   urlValues["group_name"],
		SuperGroup:  urlValues["supergroup_name"],
	}

	if strings.ToLower(input.Format) == "json" {
		contentType = "application/json"
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	found, output := caches.HitCache("endpoint_group_ar", input, cfg)

	if found {
		return code, h, output, err
	}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	results := []MongoInterface{}

	// Select the granularity of the search daily/monthly
	if len(input.Granularity) == 0 || strings.ToLower(input.Granularity) == "daily" {
		customForm[0] = "20060102"
		customForm[1] = "2006-01-02"
		query := Daily(input)
		err = mongo.Pipe(session, tenantDbConfig.Db, "endpoint_group_ar", query, &results)

	} else if strings.ToLower(input.Granularity) == "monthly" {
		customForm[0] = "200601"
		customForm[1] = "2006-01"
		query := Monthly(input)
		err = mongo.Pipe(session, tenantDbConfig.Db, "endpoint_group_ar", query, &results)
	}
	// mongo.Find(session, tenantDbConfig.Db, "endpoint_group_ar", bson.M{}, "_id", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createView(results, input.Format)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	if len(results) > 0 {
		caches.WriteCache("endpointGroup", input, output, cfg)
	}

	return code, h, output, err
}
