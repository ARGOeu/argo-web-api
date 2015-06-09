package tenants

import "labix.org/v2/mgo/bson"

// Tenant structure holds information about tenant information
// including db conf and users
type Tenant struct {
	Name   string         `bson:"name" json:"name"`
	DbConf []TenantDbConf `bson:"db_conf" json:"db_conf"`
	Users  []TenantUser   `bson:"users" json:"users"`
}

// TenantDbConf structure holds information about tenant's
// database configuration
type TenantDbConf struct {
	Store    string `bson:"store" json:"store"`
	Server   string `bson:"server" json:"server"`
	Port     int    `bson:"port" json:"port"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

// TenantUser structure holds information about tenant's
// database configuration
type TenantUser struct {
	Name   string `bson:"name" json:"name"`
	Email  string `bson:"email" json:"email"`
	APIkey string `bson:"api_key" json:"api_key"`
}

// createTenant is used to create a new
func createTenant(input Tenant) bson.M {
	query := bson.M{
		"name":    input.Name,
		"db_conf": input.DbConf,
		"users":   input.Users,
	}
	return query
}

func searchName(name String) bson.M {
	query := bson.M{
		"name": name,
	}
}
