package main

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"github.com/makistsan/go-api"
	"github.com/makistsan/go-lru-cache"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strconv"
)

var httpcache *cache.LRUCache

type mystring string

func (s mystring) Size() int {
	return len(s)
}

var cfg = LoadConfiguration()

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

func main() {
	httpcache = cache.NewLRUCache(uint64(cfg.Server.Lrucache))
	runtime.GOMAXPROCS(cfg.Server.Maxprocs)

	if *flProfile != "" {
		f, err := os.Create(*flProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			if *flProfile != "" {
				pprof.StopCPUProfile()
			}
			log.Printf("captured %v, stopping profiler and exiting..", sig)
			os.Exit(1)
		}
	}()
	handlers := map[string]func(http.ResponseWriter, *http.Request){}

	//Basic api calls
	handlers["/api/v1/service_availability_in_profile"] = func(w http.ResponseWriter, r *http.Request) {
		Respond("text/xml", "utf-8", ServiceAvailabilityInProfile)(w, r)
	}
	handlers["/api/v1/sites_availability_in_profile"] = func(w http.ResponseWriter, r *http.Request) {
		Respond("text/xml", "utf-8", SitesAvailabilityInProfile)(w, r)
	}
	handlers["/api/v1/ngi_availability_in_profile"] = func(w http.ResponseWriter, r *http.Request) {
		Respond("text/xml", "utf-8", NgiAvailabilityInProfile)(w, r)
	}
	handlers["/api/v1/profiles"] = func(w http.ResponseWriter, r *http.Request) {
		Respond("text/xml", "utf-8", GetProfileNames)(w, r)
	}

	//CRUD functions for profiles
	handlers["/api/v1/profiles/create"] = func(w http.ResponseWriter, r *http.Request) {
		Respond("text/xml", "utf-8", AddProfile)(w, r)
	}
	handlers["/api/v1/profiles/remove"] = func(w http.ResponseWriter, r *http.Request) {
		Respond("text/xml", "utf-8", RemoveProfile)(w, r)
	}
	handlers["/api/v1/profiles/getone"] = func(w http.ResponseWriter, r *http.Request) {
		Respond("text/xml", "utf-8", GetProfile)(w, r)
	}

	//Miscallenious calls
	handlers["/reset_cache"] = func(w http.ResponseWriter, r *http.Request) {
		Respond("text/xml", "utf-8", ResetCache)(w, r)
	}
	api.NewServer(cfg.Server.Bindip+":"+strconv.Itoa(cfg.Server.Port), api.DefaultServerReadTimeout, handlers)
}

func ResetCache(w http.ResponseWriter, r *http.Request) []byte {
	httpcache.Clear()
	return []byte("Cache Emptied")
}
