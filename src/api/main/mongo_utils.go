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
	"encoding/xml"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strings"
)

func ErrorXML(s string) []byte {
	return []byte("<root><error>" + s + "</error></root>")
}

type PoemProfile struct {
	Name           string "p"
	Namespace      string "ns"
	Group          string "g"
	Service_flavor string "sf"
}

func parseCSV(data string) []string {
	splitted := strings.SplitN(data, ",", -1)

	data_tmp := make([]string, len(splitted))

	for i, val := range splitted {
		data_tmp[i] = strings.TrimSpace(val)
	}

	return data_tmp
}

func CheckAndCompress(w http.ResponseWriter, r *http.Request, output *[]byte) []byte {

	if (cfg.Server.Gzip) == true && r.Header.Get("Accept-Encoding") != "" {
		encodings := parseCSV(r.Header.Get("Accept-Encoding"))

		var b bytes.Buffer
		fmt.Println(encodings)
		for _, val := range encodings {
			if val == "gzip" {
				fmt.Println("gzipping")
				//w.Header().Set("Accept-Encoding", "gzip")
				writer, _ := gzip.NewWriterLevel(&b, gzip.BestSpeed)
				writer.Write(*output)
				writer.Close()
				w.Header().Set("Content-Encoding", "gzip") //http.DetectContentType(b.Bytes()))
				return b.Bytes()
			} else if val == "deflate" {
				fmt.Println("zlib")
				//w.Header().Set("Accept-Encoding", "deflate")
				writer, _ := zlib.NewWriterLevel(&b, zlib.BestSpeed)
				writer.Write(*output)
				writer.Close()
				w.Header().Set("Content-Encoding", "deflate") //http.DetectContentType(b.Bytes()))
				return b.Bytes()
			}
		}
	}
	fmt.Println("No compression")
	return *output
}

func CreatePoemProfileNameXmlResponse(results []PoemProfile) ([]byte, error) {
	type Profile struct {
		XMLName        xml.Name `xml:"Profile"`
		Name           string   `xml:"name,attr"`
		Namespace      string   `xml:"namespace,attr"`
		Group          string   `xml:"group,attr"`
		Service_flavor string   `xml:"service_flavor,attr"`
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Profile []Profile
	}

	v := &Root{}

	for _, result := range results {
		v.Profile = append(v.Profile,
			Profile{
				Name:           result.Name,
				Namespace:      result.Namespace,
				Group:          result.Group,
				Service_flavor: result.Service_flavor,
			})
	}

	output, err := xml.MarshalIndent(v, " ", "  ")
	return output, err

}

func GetProfileNames(w http.ResponseWriter, r *http.Request) []byte {
	var results []PoemProfile
	session, err := mgo.Dial(cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(cfg.MongoDB.Db).C("sites")
	err = c.Pipe([]bson.M{{"$group": bson.M{"_id": bson.M{"ns": "$ns", "p": "$p"}}}, {"$project": bson.M{"ns": "$_id.ns", "p": "$_id.p"}}}).All(&results)
	if err != nil {
		return []byte("<root><error>" + err.Error() + "</error></root>")
	}

	fmt.Println(results)
	output, err := CreatePoemProfileNameXmlResponse(results)
	return output

}
