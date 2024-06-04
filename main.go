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

package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ARGOeu/argo-web-api/routing"
	"github.com/ARGOeu/argo-web-api/utils/store"
	"github.com/ARGOeu/argo-web-api/version"
	"github.com/gorilla/handlers"
)

func main() {

	log.SetOutput(os.Stdout)
	log.Println("starting argo-web-api " + version.Release)
	//Create a connection to mongo and store client in the configuration object

	cfg.MongoClient = store.GetMongoClient(cfg.MongoDB)
	// close connection at the end
	defer func() {
		if err := cfg.MongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	//Create the server router and add the middleware
	var mainRouter http.Handler
	mainRouter = routing.NewRouter(cfg)
	mainRouter = handlers.CombinedLoggingHandler(os.Stdout, mainRouter)

	// mainRouter = handlers.CompressHandler(mainRouter)

	http.Handle("/", mainRouter)

	//Cache
	//get_subrouter.HandleFunc("/api/v1/reset_cache", Respond("text/xml", "utf-8", ResetCache))

	//TLS support only
	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
	}

	server := &http.Server{Addr: cfg.Server.Bindip + ":" + strconv.Itoa(cfg.Server.Port), Handler: nil, TLSConfig: tlsConfig}
	//Web service binds to server. Requests served over HTTPS.
	log.Println("listening for requests at: " + server.Addr)
	err := server.ListenAndServeTLS(cfg.Server.Cert, cfg.Server.Privkey)

	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
