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
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"net/http"
	//"encoding/json"
	"time"
	"labix.org/v2/mgo"	
	"labix.org/v2/mgo/bson"	
	
)

type list []interface{}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

// The respond function that will be called to answer to http requests to the PI
func Respond(mediaType string, charset string, fn func(w http.ResponseWriter, r *http.Request) []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=%s", mediaType, charset))
		output := fn(w, r)
		var b bytes.Buffer
		var data []byte
		if (cfg.Server.Gzip) == true && r.Header.Get("Accept-Encoding") != "" {
			encodings := parseCSV(r.Header.Get("Accept-Encoding"))
			for _, val := range encodings {
				if val == "gzip" {
					writer := gzip.NewWriter(&b)
					writer.Write(output)
					writer.Close()
					w.Header().Set("Content-Encoding", "gzip")
					break
				} else if val == "deflate" {
					writer := zlib.NewWriter(&b)
					writer.Write(output)
					writer.Close()
					w.Header().Set("Content-Encoding", "deflate")
					break
				}
			}
			data = b.Bytes()
		} else {
			data = output
		}
		fmt.Println(len(data))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.Write(data)
	}
}

//Reset the cache if it is set
func ResetCache(w http.ResponseWriter, r *http.Request) []byte {
	if cfg.Server.Cache == true {
		httpcache.Clear()
		return []byte("Cache Emptied")
	}
	return []byte("No Caching is active")
}
//Scedule a recalculation
func Recalculate(w http.ResponseWriter, r *http.Request) []byte{
	type ApiRecalculationInput struct {
		Start_time          string   
		End_time            string   
		Reason 				string
		Vo_name 			string
		Ngi_name 			string
		Exclude_site		[]string
		Status				string
		Timestamp			int64			
		//Exclude_sf			[]string
		//Exclude_end_point   []string		
	}
	session, err := mgo.Dial(cfg.MongoDB.Host + ":" + fmt.Sprint(cfg.MongoDB.Port))
	if err != nil {
		return ErrorXML("Error while connecting to MongoDB")
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	c := session.DB(cfg.MongoDB.Db).C("Recalculations")
	urlValues := r.URL.Query()
	now:=time.Now()
	input:=ApiRecalculationInput{
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues.Get("reason"),
		urlValues.Get("vo_name"),
		urlValues.Get("ngi_name"),
		urlValues["exclude_site"],
		"pending",
		now.Unix(),
		//urlValues["exclude_sf"],
		//urlValues["exclude_end_point"],		
	}
	toMongo:=bson.M{
		"start_time" : input.Start_time,
		"end_time" : input.End_time,
		"reason" : input.Reason,
		"vo" : input.Vo_name,
		"ngi" : input.Ngi_name,
		"status" : input.Status,
		"timestamp": input.Timestamp,
	}
	
	answer:="An appropriate output to the WebUI"//Provide the webUI with an appropriate xml/json response 
	err=c.Insert(toMongo)
	if err!=nil{
		return ErrorXML("MongoDB write error")
	}
	return []byte(answer)
}