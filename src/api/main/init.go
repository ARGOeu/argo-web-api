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
	"github.com/makistsan/go-lru-cache"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"api/utils/config"
)

var httpcache *cache.LRUCache

type mystring string

func (s mystring) Size() int {
	return len(s)
}

// Load the configurations that we have set through flags and through the configuration file
var cfg = config.LoadConfiguration()

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