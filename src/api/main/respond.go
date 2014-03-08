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
	"api/utils/config"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"net/http"
	"strings"
	//  "encoding/xml"
)

type list []interface{}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"


func parseCSV(data string) []string {
	splitted := strings.SplitN(data, ",", -1)

	data_tmp := make([]string, len(splitted))

	for i, val := range splitted {
		data_tmp[i] = strings.TrimSpace(val)
	}

	return data_tmp
}

// The respond function that will be called to answer to http requests to the PI
func Respond(mediaType string, charset string, fn func(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=%s", mediaType, charset))
		output := fn(w, r, cfg)
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
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.Write(data)
	}
}

//Reset the cache if it is set
func ResetCache(w http.ResponseWriter, r *http.Request) []byte {
	answer := ""
	if cfg.Server.Cache == true {
		httpcache.Clear()
		answer = "Cache Emptied"
	}
	answer = "No Caching is active"
	return []byte(answer)
}
