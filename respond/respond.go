/*
 * Copyright (c) 2015 GRNET S.A.
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

package respond

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ARGOeu/argo-web-api/utils/caches"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/logging"
	"github.com/gorilla/mux"
)

type list []interface{}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

// ConfHandler Keeps all the configuration/variables required by all the requests
type ConfHandler struct {
	Config config.Config
}

// Respond will be called to answer to http requests to the PI
func (confhandler *ConfHandler) Respond(fn func(r *http.Request, cfg config.Config) (int, http.Header, []byte, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if r := recover(); r != nil {
				logging.HandleError(r)
			}
		}()
		code, header, output, err := fn(r, confhandler.Config)

		if code == http.StatusInternalServerError {
			log.Panic("Internal Server Error:", err)
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
	})

}

func (confhandler *ConfHandler) walker(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	// route.Handler(route.GetHandler())
	return nil
}

// ResetCache resets the cache if it is set
func ResetCache(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte {
	answer := ""
	if cfg.Server.Cache == true {
		caches.ResetCache()
		answer = "Cache Emptied"
	}
	answer = "No Caching is active"
	return []byte(answer)
}
