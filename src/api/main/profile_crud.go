package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"fmt"
	)

type MongoProfile struct{
	Name string "p"
	Namespace string "ns"
	Group string "g"
	Service_flavor string "sf"
}


func AddProfile(w http.ResponseWriter, r *http.Request) string {
	session, err := mgo.Dial(cfg.MongoDB.Host + ":" + fmt.Sprint(cfg.MongoDB.Port))
	if err != nil { return "Error while connecting to mongodb" }
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	c := session.DB(cfg.MongoDB.Db).C("Profiles")
	result := MongoProfile{}
	type ProfileInput struct {
                // mandatory values
                Name  string   
		Namespace string
        	Group string
		Service_flavor string
	}
        urlValues := r.URL.Query()
        input := ProfileInput{
                urlValues.Get("name"),
        	urlValues.Get("namespace"),
		urlValues.Get("group"),
		urlValues.Get("service_flavor"),
	}
	q := bson.M{
                "p":  input.Name,
        }
	err2 := c.Find(q).One(&result)
	if fmt.Sprint(err2) != "not found" {
	return "<root><error> Already exists </error></root>"
	}
	doc := bson.M{ 
			"p"  : input.Name,
			"ns" : input.Namespace,
			"g"  : input.Group,
			"sf" : input.Service_flavor,
			}
	err3 := c.Insert(doc)	
	return fmt.Sprint(err3)
}

func RemoveProfile(w http.ResponseWriter, r *http.Request) string {
	session, err := mgo.Dial(cfg.MongoDB.Host + ":" + fmt.Sprint(cfg.MongoDB.Port))
	if err != nil { return "Error while connecting to mongodb" }
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	c := session.DB(cfg.MongoDB.Db).C("Profiles")
	result := MongoProfile{}
	type ProfileInput struct {
                // mandatory values
                Name  string   
		Namespace string
        	Group string
		Service_flavor string
	}
        urlValues := r.URL.Query()
        input := ProfileInput{
                urlValues.Get("name"),
        	urlValues.Get("namespace"),
		urlValues.Get("group"),
		urlValues.Get("service_flavor"),
	}
	q := bson.M{
                "p":  input.Name,
        }
	err2 := c.Find(q).One(&result)
	if fmt.Sprint(err2) == "not found" {
	return "<root><error> Doesn't exists </error></root>"
	}
	doc := bson.M{ 
			"p"  : input.Name,
			"ns" : input.Namespace,
			"g"  : input.Group,
			"sf" : input.Service_flavor,
			}
	err3 := c.Remove(doc)	
	return fmt.Sprint(err3)
}


func GetProfile(w http.ResponseWriter, r *http.Request) string {
	session, err := mgo.Dial(cfg.MongoDB.Host + ":" + fmt.Sprint(cfg.MongoDB.Port))
	if err != nil { return "Error while connecting to mongodb" }
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	c := session.DB(cfg.MongoDB.Db).C("Profiles")
	result := MongoProfile{}
	type ProfileInput struct {
                // mandatory values
                Name  string   
		Namespace string
        	Group string
		Service_flavor string
	}
        urlValues := r.URL.Query()
        input := ProfileInput{
                urlValues.Get("name"),
        	urlValues.Get("namespace"),
		urlValues.Get("group"),
		urlValues.Get("service_flavor"),
	}
	q := bson.M{
                "p":  input.Name,
        }
	err2 := c.Find(q).One(&result)
	if err2 != nil { return fmt.Sprint(err2)}
	results := []MongoProfile{result}	
	output,err3 := CreateProfileNameXmlResponse(results)
	if err3 != nil {return fmt.Sprint(err3)}	
	return string(output)

}
