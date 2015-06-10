package tenants

import (
	"encoding/xml"

	"labix.org/v2/mgo/bson"
)

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

// Message struct for xml message response
type Message struct {
	XMLName xml.Name `xml:"root"`
	Message string
}

// TenantXML used for xml response
type TenantXML struct {
	XMLName xml.Name           `xml:"tenant" json:"-"`
	Name    string             `xml:"name,attr" json:"name"`
	DbConf  []*TenantDBConfXML `xml:"db_conf" json:"db_conf"`
	Users   []*TenantUserXML   `xml:"users" json:"users"`
}

// TenantDBConfXML used for XML response
type TenantDBConfXML struct {
	XMLName  xml.Name `xml:"db_conf" json:"-"`
	Store    string   `xml:"store,attr" json:"store"`
	Server   string   `xml:"server,attr" json:"server"`
	Port     int      `xml:"port,attr" json:"port"`
	Database string   `xml:"database,attr" json:"database"`
	Username string   `xml:"username,attr" json:"username"`
	Password string   `xml:"password,attr" json:"password"`
}

// Root struct to represent the root of the xml/json document
type Root struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Tenants []*TenantXML
}

// TenantUserXML used for XML response
type TenantUserXML struct {
	XMLName xml.Name `xml:"user" json:"-"`
	Name    string   `xml:"name,attr" json:"name"`
	Email   string   `xml:"email,attr" json:"email"`
	APIkey  string   `xml:"api_key,attr" json:"api_key"`
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

func searchName(name string) bson.M {
	query := bson.M{
		"name": name,
	}

	return query
}
