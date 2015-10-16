package logging

import (
	"log"
	"net/http"
)

// RequestError is used to pass the errors from controllers to repsonse handlers]
type RequestError struct {
	code   int
	header http.Header
	output []byte
	err    error
}

// HandleError accepts errors and logs the appropriate information
func HandleError(reqErr interface{}) {
	log.Printf("%+v", reqErr.(error).Error())
}
