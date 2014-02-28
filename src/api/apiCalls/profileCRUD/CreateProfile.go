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
	"strings"
	"api/utils/config"
	"api/utils/authentication"
)

type MongoProfile struct {
	Name  string     "p"
	Group [][]string "g"
}

type list []interface{}

func CreateProfile(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte {
	answer:=""
	if authentication.Authenticate(r.Header,cfg){
		session, err := mgo.Dial(cfg.MongoDB.Host + ":" + fmt.Sprint(cfg.MongoDB.Port))
		if err != nil {
			return []byte("ERROR")//TODO
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
		err=r.ParseForm()
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
			return []byte("ERROR")//TODO
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
			return []byte("ERROR")//TODO
		} else {
			return []byte("ERROR")//TODO
		}
	} else {
		answer = http.StatusText(403)
	}
	return []byte (answer)
}	