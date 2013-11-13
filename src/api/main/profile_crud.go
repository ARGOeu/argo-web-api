package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strings"
	"encoding/xml"
	"encoding/json"
)

type MongoProfile struct {
	Name  string     "p"
	Group [][]string "g"
}

func AddProfile(w http.ResponseWriter, r *http.Request) string {
	session, err := mgo.Dial(cfg.MongoDB.Host + ":" + fmt.Sprint(cfg.MongoDB.Port))
	if err != nil {
		return ErrorXML("Error while connecting to mongodb")
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	c := session.DB(cfg.MongoDB.Db).C("Profiles")
	result := MongoProfile{}
	type ProfileInput struct {
		// mandatory values
		Name  string
		Group []string
		Json  string
	}
	urlValues := r.URL.Query()
	input := ProfileInput{
		urlValues.Get("name"),
		urlValues["group"],
		urlValues.Get("json"),
	}
	q := bson.M{
		"p": input.Name,
	}
	err2 := c.Find(q).One(&result)
	if fmt.Sprint(err2) != "not found" {
		return ErrorXML("Already exists")
	}
	if len(input.Group) > 0 {

		doc := bson.M{
			"p": input.Name,
		}
		groups := make(list, 0)
		for _, value := range input.Group {
			//doc["g"] = append(doc["g"], value)
			groups = append(groups, strings.Split(value, ","))
		}
		doc["g"] = groups
		err3 := c.Insert(doc)
		return fmt.Sprint(err3)

	} else if len(input.Json) > 0 {
		return ErrorXML("Not implemented yet")
	} else {
		return ErrorXML("Could not find data to save")
	}

}

func RemoveProfile(w http.ResponseWriter, r *http.Request) string {
	session, err := mgo.Dial(cfg.MongoDB.Host + ":" + fmt.Sprint(cfg.MongoDB.Port))
	if err != nil {
		return ErrorXML("Error while connecting to mongodb")
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	c := session.DB(cfg.MongoDB.Db).C("Profiles")
	result := MongoProfile{}
	type ProfileInput struct {
		// mandatory values
		Name string
	}
	urlValues := r.URL.Query()
	input := ProfileInput{
		urlValues.Get("name"),
	}
	q := bson.M{
		"p": input.Name,
	}
	err2 := c.Find(q).One(&result)
	if fmt.Sprint(err2) == "not found" {
		return ErrorXML("Doesn't exists")
	}
	doc := bson.M{
		"p": input.Name,
	}
	err3 := c.Remove(doc)
	return fmt.Sprint(err3)
}

func GetProfile(w http.ResponseWriter, r *http.Request) string {
	session, err := mgo.Dial(cfg.MongoDB.Host + ":" + fmt.Sprint(cfg.MongoDB.Port))
	if err != nil {
		return ErrorXML("Error while connecting to mongodb")
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	c := session.DB(cfg.MongoDB.Db).C("Profiles")
	result := MongoProfile{}
	type ProfileInput struct {
		// mandatory values
		Name string
		Output string
	}
	urlValues := r.URL.Query()
	input := ProfileInput{
		urlValues.Get("name"),
		urlValues.Get("output"),
	}
	q := bson.M{
		"p": input.Name,
	}
	err2 := c.Find(q).One(&result)
	if err2 != nil {
		return fmt.Sprint(err2)
	}
	results := []MongoProfile{result}
	
	err3 := error(nil)
	output := []byte(nil)
	
	if strings.ToLower(input.Output) == "json" {
		output, err3 = CreateProfileJsonResponse(results)
	} else {
		output, err3 = CreateProfileXmlResponse(results)
	}
	if err3 != nil {
		return fmt.Sprint(err3)
	}
	return string(output)

}

func CreateProfileJsonResponse(results []MongoProfile) ([]byte, error) {

	type Groups struct {
		Service_flavors []string
	}

	type Profile struct {
		Name   string
		Groups []Groups
	}

	type Root struct {
		Profiles []Profile
	}
	v := &Root{}

	for key, result := range results {
		v.Profiles = append(v.Profiles,
			Profile{
				Name: result.Name,
			})
		for _, result2 := range result.Group {
			v.Profiles[key].Groups = append(v.Profiles[key].Groups,
				Groups{
					Service_flavors: result2,
				})
		}
	}

	output, err := json.MarshalIndent(v," ", "  ")
	return output, err

}

func CreateProfileXmlResponse(results []MongoProfile) ([]byte, error) {

	type Groups struct {
		XMLName xml.Name `xml:"Group"`
		Service_flavors []string `xml:"service_flavor"`   
	}

	type Profile struct {
		XMLName        xml.Name `xml:"Profile"`
		Name           string   `xml:"name,attr"`
		Groups         []Groups
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Profiles []Profile
	}

	v := &Root{}

	for key, result := range results {
		v.Profiles = append(v.Profiles,
			Profile{
				Name: result.Name,
			})
		for _, result2 := range result.Group {
			v.Profiles[key].Groups = append(v.Profiles[key].Groups,
				Groups{
					Service_flavors: result2,
				})
		}
	}
	
	output, err := xml.MarshalIndent(v, " ", "  ")
	return output, err
}
