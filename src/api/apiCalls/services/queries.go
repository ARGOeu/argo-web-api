/*
 * Copyright (c) 2013 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
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

package services

import (
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

type list []interface{}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

func prepareFilter(input ApiServiceAvailabilityInProfileInput) bson.M {

	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	filter := bson.M{
		"d":  bson.M{"$gte": tsYMD, "$lte": teYMD},
		"vo": bson.M{"$in": input.vo_name},
		"p":  bson.M{"$in": input.profile_name},
	}

	if len(input.namespace) > 0 {
		filter["ns"] = bson.M{"$in": input.namespace}
	}

	if len(input.group_name) > 0 {
		// TODO: We do not have the site name in the timeline
	}

	if len(input.service_flavour) > 0 {
		filter["sf"] = bson.M{"$in": input.service_flavour}
	}

	if len(input.service_hostname) > 0 {
		filter["h"] = bson.M{"$in": input.service_hostname}
	}

	return filter
}

func Timeline(input ApiServiceAvailabilityInProfileInput) []bson.M {
	filter := prepareFilter(input)

	query := []bson.M{{"$match": filter}, {"$sort": bson.D{{"p", 1}, {"h", 1}, {"sf", 1}, {"d", 1}}}}

	return query

}
