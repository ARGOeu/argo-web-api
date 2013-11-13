package main

import (
	"code.google.com/p/gcfg"
	"flag"
	"os"
)

var flConfig = flag.String("conf", "", "specify configuration file")
var flServerIp = flag.String("ip", "", "ip address the server will bind to")
var flServerPort = flag.Int("port", 0, "specify the port to listen on")
var flServerMaxProcs = flag.Int("maxprocs", 0, "specify the GOMAXPROCS")
var flMongoHost = flag.String("mongo-host", "", "specify the IP address of the MongoDB instance")
var flMongoPort = flag.Int("mongo-port", 0, "specify the port on which the MongoDB instance listens on")
var flMongoDatabase = flag.String("mongo-db", "", "specify the MongoDB database to connect to")
var flCache = flag.String("cache", "no", "specify weather to use cache or not [yes/no]")

type Config struct {
	Server struct {
		Bindip   string
		Port     int
		Maxprocs int
		Cache    bool
		Lrucache int
	}
	MongoDB struct {
		Host string
		Port int
		Db   string
	}
}

const defaultConfig = `
    [server]
    bindip = ""
    port = 8080
    maxprocs = 4
    cache = false
    lrucache = 700000000

    [mongodb]
    host = "127.0.0.1"
    port = 27017
    db = "AR"
`

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

	return cfg
}
