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
	"api/utils/authentication"
	"api/utils/config"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func DeleteProfile(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte {
	answer := ""
	if authentication.Authenticate(r.Header, cfg) {
		session, err := mgo.Dial(cfg.MongoDB.Host + ":" + fmt.Sprint(cfg.MongoDB.Port))
		if err != nil {
			return []byte("ERROR") //TODO
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
			return []byte("ERROR") //TODO
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
