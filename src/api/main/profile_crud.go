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

package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strings"
)

type MongoProfile struct {
	Name  string     "p"
	Group [][]string "g"
}

func AddProfile(w http.ResponseWriter, r *http.Request) []byte {
	answer:=""
	if Authenticate(r.Header) {
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
		err = r.ParseForm()
		urlValues := r.Form
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
			return []byte(fmt.Sprint(err3))

		} else if len(input.Json) > 0 {
			return ErrorXML("Not implemented yet")
		} else {
			return ErrorXML("Could not find data to save")
		}

	} else {
		answer = http.StatusText(403)
	}
	return []byte(answer)
}

func RemoveProfile(w http.ResponseWriter, r *http.Request) []byte {
	answer :=""
	if Authenticate(r.Header) {
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
		urlValues := r.URL.Query() //CONVERT TO POST
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
		return []byte(fmt.Sprint(err3))
	} else {
		answer = http.StatusText(403)
	}
	return []byte(answer)
}

func GetProfile(w http.ResponseWriter, r *http.Request) []byte {
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
		Name   string
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
		return []byte(fmt.Sprint(err2))
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
		return []byte(fmt.Sprint(err3))
	}
	return output

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

	output, err := json.MarshalIndent(v, " ", "  ")
	return output, err

}

func CreateProfileXmlResponse(results []MongoProfile) ([]byte, error) {

	type Groups struct {
		XMLName         xml.Name `xml:"Group"`
		Service_flavors []string `xml:"service_flavor"`
	}

	type Profile struct {
		XMLName xml.Name `xml:"Profile"`
		Name    string   `xml:"name,attr"`
		Groups  []Groups
	}

	type Root struct {
		XMLName  xml.Name `xml:"root"`
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
