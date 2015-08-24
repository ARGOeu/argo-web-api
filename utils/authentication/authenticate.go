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
	"errors"
	"net/http"

	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Auth struct {
	ApiKey string `bson:"api_key"`
}

func Authenticate(h http.Header, cfg config.Config) bool {

	session, err := mongo.OpenSession(cfg.MongoDB)

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

// AuthenticateAdmin is used to authenticate and administrator of ARGO
// and allow further CRUD ops wrt the argo_core database (i.e. add a new
// tenant, modify another tenant's configuration etc)
func AuthenticateAdmin(h http.Header, cfg config.Config) bool {

	session, err := mongo.OpenSession(cfg.MongoDB)
	if err != nil {
		panic(err)
	}
	defer mongo.CloseSession(session)

	query := bson.M{"api_key": h.Get("x-api-key")}
	projection := bson.M{"_id": 0, "name": 0, "email": 0}

	results := []Auth{}
	err = mongo.FindAndProject(session, cfg.MongoDB.Db, "authentication", query, projection, "api_key", &results)

	if err != nil {
		return false
	}

	if len(results) > 0 {
		return true
	}
	return false
}

// AuthenticateTenant is used to find which tenant the user making the requests
// belongs to and return the database configuration for that specific tenant.
// If the api-key in the request is not found in any tenant an empty configuration is
// returned along with an error
func AuthenticateTenant(h http.Header, cfg config.Config) (config.MongoConfig, error) {

	session, err := mongo.OpenSession(cfg.MongoDB)
	if err != nil {
		return config.MongoConfig{}, err
	}
	defer mongo.CloseSession(session)

	query := bson.M{"users.api_key": h.Get("x-api-key")}
	projection := bson.M{"_id": 0, "name": 1, "db_conf": 1}

	var results []map[string][]config.MongoConfig
	mongo.FindAndProject(session, cfg.MongoDB.Db, "tenants", query, projection, "server", &results)

	if len(results) == 0 {
		return config.MongoConfig{}, errors.New("Unauthorized")
	}
	mongoConf := results[0]["db_conf"][0]
	return mongoConf, nil
}
