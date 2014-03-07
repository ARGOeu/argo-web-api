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

package profileCRUD

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	//"strings"
	"api/apiCalls/profileCRUD/XMLresponses"
	"api/utils/config"
	//"api/apiCalls/profileCRUD/JSONresponses"
)

func ReadOneProfile(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte {
	session, err := mgo.Dial(cfg.MongoDB.Host + ":" + fmt.Sprint(cfg.MongoDB.Port))
	if err != nil {
		return []byte("ERROR") //TODO
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	c := session.DB(cfg.MongoDB.Db).C("Profiles")
	result := XMLresponses.MongoProfile{}
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
	results := []XMLresponses.MongoProfile{result}

	err3 := error(nil)
	output := []byte(nil)

	output, err3 = XMLresponses.ReadOneXmlResponse(results)

	//FIX THIS !!!!
	// if strings.ToLower(input.Output) == "json" {
	// 		output, err3 = JSONresponses.ReadOneJsonResponse(results)
	// 	} else {
	// 		output, err3 = XMLresponses.ReadOneXmlResponse(results)
	// 	}
	if err3 != nil {
		return []byte(fmt.Sprint(err3))
	}
	return output
}
