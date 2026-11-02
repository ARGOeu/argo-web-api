package logging

import (
	"log"
)

// HandleError accepts errors and logs the appropriate information
func HandleError(reqErr interface{}) {

	log.Printf("%+v", reqErr)
}
