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

package tenants

// Tenant structure holds information about tenant information
// including db conf and users. Used in
type Tenant struct {
	ID     string         `bson:"id" json:"id"`
	Info   TenantInfo     `bson:"info" json:"info"`
	DbConf []TenantDbConf `bson:"db_conf" json:"db_conf,omitempty"`
	Users  []TenantUser   `bson:"users" json:"users,omitempty"`
}

// TenantInfo struct holds information about tenant name, contact details
type TenantInfo struct {
	Name    string `bson:"name" json:"name"`
	Email   string `bson:"email" json:"email"`
	Website string `bson:"website" json:"website"`
	Created string `bson:"created" json:"created"`
	Updated string `bson:"updated" json:"updated"`
}

// TenantDbConf structure holds information about tenant's
// database configuration
type TenantDbConf struct {
	Store    string `bson:"store" json:"store"`
	Server   string `bson:"server" json:"server"`
	Port     int    `bson:"port" json:"port"`
	Database string `bson:"database" json:"database"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

// TenantUser structure holds information about tenant's
// database configuration
type TenantUser struct {
	Name   string   `bson:"name"       json:"name"`
	Email  string   `bson:"email"      json:"email"`
	APIkey string   `bson:"api_key"    json:"api_key"`
	Roles  []string `bson:"roles,omitempty"      json:"roles,omitempty"`
}

// SelfReference to hold links and id
type SelfReference struct {
	ID    string `json:"id" bson:"id,omitempty"`
	Links Links  `json:"links"`
}

// Links struct to hold links
type Links struct {
	Self string `json:"self"`
}
