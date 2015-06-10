package tenants

import (
	"encoding/xml"

	"labix.org/v2/mgo/bson"
)

// Tenant structure holds information about tenant information
// including db conf and users. Used in
type Tenant struct {
	XMLName xml.Name       `bson:",omitempty" json:"-"       xml:"tenant" `
	Name    string         `bson:"name"       json:"name"    xml:"name,attr" `
	DbConf  []TenantDbConf `bson:"db_conf"    json:"db_conf" xml:"db_confs>db_conf"`
	Users   []TenantUser   `bson:"users"      json:"users"   xml:"users>user"`
}

// TenantDbConf structure holds information about tenant's
// database configuration
type TenantDbConf struct {
	XMLName  xml.Name `bson:",omitempty" json:"-"        xml:"db_conf"`
	Store    string   `bson:"store"      json:"store"    xml:"store,attr"`
	Server   string   `bson:"server"     json:"server"   xml:"server,attr"`
	Port     int      `bson:"port"       json:"port"     xml:"port,attr"`
	Username string   `bson:"username"   json:"username" xml:"username,attr"`
	Password string   `bson:"password"   json:"password" xml:"password,attr"`
}

// TenantUser structure holds information about tenant's
// database configuration
type TenantUser struct {
	XMLName xml.Name `bson:",omitempty" json:"-"       xml:"user"`
	Name    string   `bson:"name"       json:"name"    xml:"name,attr"`
	Email   string   `bson:"email"      json:"email"   xml:"email,attr"`
	APIkey  string   `bson:"api_key"    json:"api_key" xml:"api_key,attr"`
}

// Message struct for xml message response
type Message struct {
	XMLName xml.Name `xml:"root"`
	Message string
}

// RootXML struct to represent the root of the xml/json document
type RootXML struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Tenants *[]Tenant
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

// searchName is used to create a simple query object based on name
func searchName(name string) bson.M {
	query := bson.M{
		"name": name,
	}

	return query
}
