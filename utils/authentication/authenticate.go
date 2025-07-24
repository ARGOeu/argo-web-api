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
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/ARGOeu/argo-web-api/utils/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Auth struct {
	ApiKey       string `bson:"api_key"`
	Restricted   bool   `bson:"restricted"`
	SuperAdminUI bool   `bson:"super_admin_ui"`
}

type Tenant struct {
	Info InfoStruct `bson:"info"`
}

type InfoStruct struct {
	Name string `bson:"name"`
}

type Info struct {
	Name string `bson:"name"`
}

type DbInfoUsers struct {
	Info   Info         `bson:"info"`
	DbConf []DbConfItem `bson:"db_conf"`
	Users  []UserItem   `bson:"users"`
}

type DbConfItem struct {
	Store    string `bson:"store"`
	Server   string `bson:"server"`
	Port     int    `bson:"port"`
	Database string `bson:"database"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}

type UserItem struct {
	Id     string   `bson:"id"`
	Name   string   `bson:"name"`
	Email  string   `bson:"email"`
	ApiKey string   `bson:"api_key"`
	Roles  []string `bson:"roles"`
}

func Authenticate(h http.Header, cfg config.Config) bool {

	authCol := cfg.MongoClient.Database(cfg.MongoDB.Db).Collection("authentication")
	apiKey := h.Get("x-api-key")
	query := bson.M{
		"api_key": apiKey,
	}
	result := Auth{}
	err := authCol.FindOne(context.TODO(), query).Decode(&result)
	if err == nil {
		if result.ApiKey == apiKey {
			return true
		}
	}
	return false
}

// AuthenticateAdmin is used to authenticate and administrator of ARGO
// and allow further CRUD ops wrt the argo_core database (i.e. add a new
// tenant, modify another tenant's configuration etc)
func AuthenticateAdmin(h http.Header, cfg config.Config) bool {
	return Authenticate(h, cfg)
}

func queryAdminRoles(h http.Header, cfg config.Config) Auth {

	authCol := cfg.MongoClient.Database(cfg.MongoDB.Db).Collection("authentication")

	query := bson.M{
		"api_key": h.Get("x-api-key"),
	}

	result := Auth{}
	err := authCol.FindOne(context.TODO(), query).Decode(&result)

	if err == nil {
		return result
	}

	return Auth{Restricted: false, SuperAdminUI: false}
}

// IsAdminRestricted resturns a boolean value if an admin user is in restricted mode or not.
// Admin user in restricted mode has read only access to certain calls
func IsAdminRestricted(h http.Header, cfg config.Config) bool {
	auth := queryAdminRoles(h, cfg)
	return auth.Restricted
}

// IsSuperAdminUI resturns a boolean value if the user is a dedicated super admin ui service user
func IsSuperAdminUI(h http.Header, cfg config.Config) bool {
	auth := queryAdminRoles(h, cfg)
	return auth.SuperAdminUI
}

// AuthenticateTenant is used to find which tenant the user making the requests
// belongs to and return the database configuration for that specific tenant.
// If the api-key in the request is not found in any tenant an empty configuration is
// returned along with an error
func AuthenticateTenant(h http.Header, cfg config.Config) (config.MongoConfig, string, error) {

	tenantsCol := cfg.MongoClient.Database(cfg.MongoDB.Db).Collection("tenants")

	apiKey := h.Get("x-api-key")
	query := bson.M{"users.api_key": apiKey}
	projection := bson.M{"_id": 0, "info.name": 1, "db_conf": 1, "users": 1}

	var result DbInfoUsers

	err := tenantsCol.FindOne(context.TODO(), query, options.FindOne().SetProjection(projection)).Decode(&result)

	if err == nil {

		mongoConf := config.MongoConfig{}

		for _, user := range result.Users {
			if user.ApiKey == apiKey {
				mongoConf.User = user.Name
				mongoConf.Email = user.Email
				mongoConf.Roles = user.Roles
			}
		}

		mongoConf.Db = result.DbConf[0].Database

		log.Printf("ACCESS User: %s", mongoConf.User)
		log.Printf("ACESSS Tenant: %s", result.Info.Name)
		return mongoConf, result.Info.Name, nil

	}

	return config.MongoConfig{}, "", errors.New("Unauthorized")

}
