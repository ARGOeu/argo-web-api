package main

import (
	"github.com/makistsan/go-lru-cache"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"api/utils"
)

var httpcache *cache.LRUCache

type mystring string

func (s mystring) Size() int {
	return len(s)
}

// Load the configurations that we have set through flags and through the configuration file
var cfg = utils.LoadConfiguration()

func init() {

	//Create a recover function to log the case of a failure
	defer func() {
		if err := recover(); err != nil {
			log.Println("work failed:", err)
		}
	}()

	//Initialize the cache
	httpcache = cache.NewLRUCache(uint64(cfg.Server.Lrucache))

	//Set GOMAXPROCS
	runtime.GOMAXPROCS(cfg.Server.Maxprocs)

	//Start the profiler if the flag flProfile is set to a filename, where profile data will be writter
	if cfg.Profile != "" {
		f, err := os.Create(cfg.Profile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	//Catch an terminate signal and write all profiling data before exiting
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			if cfg.Profile != "" {
				pprof.StopCPUProfile()
			}
			log.Printf("captured %v, stopping profiler and exiting..", sig)
			os.Exit(1)
		}
	}()
}