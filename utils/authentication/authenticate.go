/*
 * Copyright (c) 2014 GRNET S.A.
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

package authentication

import (
	"github.com/argoeu/ar-web-api/utils/config"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

type Auth struct {
	apiKey string `bson:"apiKey"`
}

func Authenticate(h http.Header, cfg config.Config) bool {

	var result []Auth

	session, err := mgo.Dial(cfg.MongoDB.Host) //conect to mongo server
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)                //set mongo to monotonic behavioer
	c := session.DB(cfg.MongoDB.Db).C("authentication") //connect to collection
	//define the query to be retrieved from mongo
	retrieve := bson.M{
		"apiKey": h.Get("x-api-key"),
	}
	err = c.Find(retrieve).All(&result)
	if err != nil {
		panic(err)
	}
	//if password is found we return true
	if len(result) > 0 {
		return true
	}
	return false //else we return false
}
