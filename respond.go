/*
 * Copyright (c) 2014 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
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
	"github.com/argoeu/argo-web-api/utils/config"
	"log"
	"net/http"
	"strings"
)

type list []interface{}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

// The respond function that will be called to answer to http requests to the PI
func Respond(fn func(r *http.Request, cfg config.Config) (int, http.Header, []byte, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		code, header, output, err := fn(r, cfg)

		if code == http.StatusInternalServerError {
			log.Panic("Internal Server Error:", err)
		}

		encoding := strings.Split(r.Header.Get("Accept-Encoding"), ",")[0] //get the first accepted encoding
		if (cfg.Server.Gzip) == true && r.Header.Get("Accept-Encoding") != "" {
			var b bytes.Buffer
			if encoding == "gzip" {
				writer := gzip.NewWriter(&b)
				writer.Write(output)
				writer.Close()
				w.Header().Set("Content-Encoding", "gzip")

			} else if encoding == "deflate" {
				writer := zlib.NewWriter(&b)
				writer.Write(output)
				writer.Close()
				w.Header().Set("Content-Encoding", "deflate")
			}
			output = b.Bytes()
		}
		//Add headers

		header.Set("Content-Length", fmt.Sprintf("%d", len(output)))

		for name, values := range header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		w.WriteHeader(code)
		w.Write(output)
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
