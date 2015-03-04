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

package config

import (
	"code.google.com/p/gcfg"
	"flag"
	"os"
)

//All the flags that can be added when starting the PI
var flConfig = flag.String("conf", "", "specify configuration file")
var flServerIp = flag.String("ip", "", "ip address the server will bind to")
var flServerPort = flag.Int("port", 0, "specify the port to listen on")
var flServerMaxProcs = flag.Int("maxprocs", 0, "specify the GOMAXPROCS")
var flMongoHost = flag.String("mongo-host", "", "specify the IP address of the MongoDB instance")
var flMongoPort = flag.Int("mongo-port", 0, "specify the port on which the MongoDB instance listens on")
var flMongoDatabase = flag.String("mongo-db", "", "specify the MongoDB database to connect to")
var flCache = flag.String("cache", "no", "specify weather to use cache or not [yes/no]")
var flGzip = flag.String("gzip", "yes", "specify weather to use compression or not [yes/no]")
var flProfile = flag.String("cpuprofile", "", "write cpu profile to file")
var flCert = flag.String("cert", "", "speficy path to the host certificate")
var flPrivKey = flag.String("privkey", "", "speficy path to the private key file")

type Config struct {
	Server struct {
		Bindip   string
		Port     int
		Maxprocs int
		Cache    bool
		Lrucache int
		Gzip     bool
		Cert     string
		Privkey  string
	}
	MongoDB struct {
		Host string
		Port int
		Db   string
	}
	Profile string
}

const defaultConfig = `
    [server]
    bindip = ""
    port = 443
    maxprocs = 4
    cache = false
    lrucache = 700000000
    gzip = true
    cert = /etc/pki/tls/certs/localhost.crt
    privkey = /etc/pki/tls/private/localhost.key

    [mongodb]
    host = "127.0.0.1"
    port = 27017
    db = "AR"
`

//Loads the configurations passed either by flags or by the configuration file
func LoadConfiguration() Config {
	flag.Parse()
	var cfg Config
	if *flConfig != "" {
		_ = gcfg.ReadFileInto(&cfg, *flConfig)
	} else {
		_ = gcfg.ReadStringInto(&cfg, defaultConfig)
	}

	var env = os.Getenv("EGI_AR_REST_API_ENV")
	switch env {
	default:
		os.Setenv("EGI_AR_REST_API_ENV", "development")
	case "test":
	case "production":
	}

	if *flServerIp != "" {
		cfg.Server.Bindip = *flServerIp
	}
	if *flServerPort != 0 {
		cfg.Server.Port = *flServerPort
	}
	if *flServerMaxProcs != 0 {
		cfg.Server.Maxprocs = *flServerMaxProcs
	}
	if *flMongoHost != "" {
		cfg.MongoDB.Host = *flMongoHost
	}
	if *flMongoPort != 0 {
		cfg.MongoDB.Port = *flMongoPort
	}
	if *flMongoDatabase != "" {
		cfg.MongoDB.Db = *flMongoDatabase
	}
	if *flCache == "yes" {
		cfg.Server.Cache = true
	}
	if *flGzip == "no" {
		cfg.Server.Gzip = false
	}
	if *flProfile != "" {
		cfg.Profile = *flProfile
	}

	if *flCert != "" {
		cfg.Server.Cert = *flCert
	}

	if *flPrivKey != "" {
		cfg.Server.Privkey = *flPrivKey
	}

	return cfg
}
