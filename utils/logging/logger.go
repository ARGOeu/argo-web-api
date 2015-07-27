package logging

import (
	"log"
	"net/http"
	"time"
)

type RequestError struct {
	code   int
	header http.Header
	output []byte
	err    error
}

func Logger(inner http.Handler, name string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func HandleError(reqErr interface{}) {

}
