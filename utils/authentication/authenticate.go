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

package authentication

import (
	"github.com/argoeu/ar-web-api/utils/config"
	"github.com/argoeu/ar-web-api/utils/mongo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

type Auth struct {
	ApiKey string `bson:"apiKey"`
}

func Authenticate(h http.Header, cfg config.Config) bool {

	session, err := mongo.OpenSession(cfg)

	query := bson.M{
		"apiKey": h.Get("x-api-key"),
	}

	results := []Auth{}
	err = mongo.Find(session, cfg.MongoDB.Db, "authentication", query, "apiKey", &results)

	if err != nil {
		return false
	}

	if len(results) > 0 {
		return true
	}
	return false
}
