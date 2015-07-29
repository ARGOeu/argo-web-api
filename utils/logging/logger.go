package logging

import (
	"log"
	"net/http"
)

type RequestError struct {
	code   int
	header http.Header
	output []byte
	err    error
}

// HandleError accepts errors and logs the appropriate information
func HandleError(reqErr interface{}) {
	log.Printf("%s", reqErr.(error).Error())
}
