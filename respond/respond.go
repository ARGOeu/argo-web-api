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

package respond

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ARGOeu/argo-web-api/utils/caches"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/logging"
	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// AppRoutes holds
type AppRoutes struct {
	Name             string
	Verb             string
	Path             string
	SubrouterHandler func(r *http.Request, cfg config.Config) (int, http.Header, []byte, error)
}

type list []interface{}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

// ErrEnum used as type for enumerations of errors
type ErrEnum int

const (
	//ErrAuthen is Error during authentication
	ErrAuthen ErrEnum = iota
	//ErrAuthor is Error during authorization
	ErrAuthor ErrEnum = iota
	//ErrValidHead is Error during validation
	ErrValidHead ErrEnum = iota
	//ErrValidQuery is Error during validation
	ErrValidQuery ErrEnum = iota
)

// ConfHandler Keeps all the configuration/variables required by all the requests
type ConfHandler struct {
	Config config.Config
}

// ResponseMessage is used to construct and marshal correctly response messages
type ResponseMessage struct {
	XMLName xml.Name       `xml:"root" json:"-"`
	Status  StatusResponse `xml:"status,omitempty" json:"status,omitempty"`
	Data    interface{}    `xml:"data>result,omitempty" json:"data,omitempty"`
	Errors  interface{}    `xml:"errors>error,omitempty" json:"errors,omitempty"`
}

// StatusResponse accompanies the ResponseMessage struct to construct a response
type StatusResponse struct {
	Message string `xml:"message,omitempty" json:"message,omitempty"`
	Code    string `xml:"code,omitempty" json:"code,omitempty"`
	Details string `xml:"details,omitempty" json:"details,omitempty"`
}

// ErrorResponse holds a list of error objects
type ErrorResponse struct {
	Message string `xml:"message,omitempty" json:"message,omitempty"`
	Code    string `xml:"code,omitempty" json:"code,omitempty"`
	Details string `xml:"details,omitempty" json:"details,omitempty"`
}

// PrepAppRoutes is used in apps to prepare app's routes
func PrepAppRoutes(s *mux.Router, confHandler *ConfHandler, routes []AppRoutes) *mux.Router {
	for _, route := range routes {
		// prepare handle wrappers
		var handler http.HandlerFunc

		handler = confHandler.Respond(route.SubrouterHandler)
		handler = WrapValidate(handler, confHandler.Config, route.Name)
		if (route.Verb != "OPTIONS") && (route.Name != "version.list") {
			handler = WrapAuthorize(handler, confHandler.Config, route.Name)
			handler = WrapAuthenticate(handler, confHandler.Config, route.Name)
		}
		s.Methods(route.Verb).
			Path(route.Path).
			Handler(gcontext.ClearHandler(handler))
	}

	return s

}

// Error responds immediately when errors arise in handler chain
func Error(w http.ResponseWriter, r *http.Request, errType ErrEnum, cfg config.Config, errs []ErrorResponse) {
	//Add headers

	var msg ResponseMessage
	var contentType string
	var output []byte
	var code int
	header := r.Header

	switch errType {
	case ErrAuthen:
		msg = UnauthorizedMessage
		code = http.StatusUnauthorized
		contentType = r.Header.Get("Accept")
		output, _ = MarshalContent(msg, contentType, "", " ")
	case ErrAuthor:
		msg = Forbidden
		code = http.StatusForbidden
		contentType = r.Header.Get("Accept")
		output, _ = MarshalContent(msg, contentType, "", " ")
	case ErrValidHead:
		msg = NotAcceptableContentType
		code = http.StatusNotAcceptable
		contentType = "application/json"
		output, _ = MarshalContent(msg, contentType, "", " ")
	case ErrValidQuery:
		msg = BadRequestSimple
		code = http.StatusBadRequest
		contentType = r.Header.Get("Accept")
		output = CreateFailureResponseMessage("Bad Request", strconv.Itoa(code), errs).MarshalTo(contentType)
	default:
		msg = InternalServerErrorMessage
		code = http.StatusInternalServerError
		contentType = "application/json"
		output, _ = MarshalContent(msg, contentType, "", " ")
	}

	header.Set("Content-Length", fmt.Sprintf("%d", len(output)))

	if cfg.Server.EnableCors {
		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Headers", "Content-Type, Accept, x-api-key")
		header.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
	}

	for name, values := range header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(code)
	w.Write(output)
}

// Respond will be called to answer to http requests to the PI
func (confhandler *ConfHandler) Respond(fn func(r *http.Request, cfg config.Config) (int, http.Header, []byte, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if r := recover(); r != nil {
				logging.HandleError(r)
			}
		}()

		code, header, output, err := fn(r, confhandler.Config)

		if code == http.StatusInternalServerError {
			log.Println("Internal Server Error:", fmt.Sprintf("%+v", err))

			errResponse := ErrorResponse{Message: err.Error(), Code: "500", Details: err.Error()}
			errMsg := CreateFailureResponseMessage(err.Error(), "500", []ErrorResponse{errResponse})
			output = errMsg.MarshalTo("application/json")
		}

		//Add headers
		header.Set("Content-Length", fmt.Sprintf("%d", len(output)))

		if confhandler.Config.Server.EnableCors {
			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Headers", "Content-Type, Accept, x-api-key")
			header.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
		}

		for name, values := range header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		w.WriteHeader(code)
		w.Write(output)
	})

}

var acceptedContentTypes = []string{
	"application/xml",
	"application/json",
}

var defaultContentType = "application/json"

// ParseAcceptHeader parses the accept header to determine the content type
func ParseAcceptHeader(r *http.Request) (string, error) {
	contentType := r.Header.Get("Accept")
	if r.Header.Get("Accept") == "" {
		return defaultContentType, nil
	}
	// contentType := httputil.NegotiateContentType(r, acceptedContentTypes, "notvalid")
	if strings.Contains(contentType, "application/json") {
		return "application/json", nil
	} else if strings.Contains(contentType, "application/xml") {
		return "application/xml", nil
	} else if strings.Contains(contentType, "*/*") {
		return "application/json", nil
	}

	return defaultContentType, errors.New("Not Acceptable ContentType")

}

// MarshalContent marshals content using the marshaler that corresponds to the contentType parameter
func MarshalContent(doc interface{}, contentType string, prefix string, indent string) ([]byte, error) {
	var output []byte
	var err error

	if contentType == "application/xml" {
		output, err = xml.MarshalIndent(doc, prefix, indent)
	} else {
		output, err = json.MarshalIndent(doc, prefix, indent)
	}

	return output, err
}

func (confhandler *ConfHandler) walker(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	// route.Handler(route.GetHandler())
	return nil
}

// ResetCache resets the cache if it is set
func ResetCache(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte {
	answer := ""
	if cfg.Server.Cache == true {
		caches.ResetCache()
		answer = "Cache Emptied"
	}
	answer = "No Caching is active"
	return []byte(answer)
}

// SelfReference struct for self referencing resource after they are created
type SelfReference struct {
	ID    string    `xml:"id" json:"id" bson:"id,omitempty"`
	Links SelfLinks `xml:"links" json:"links"`
}

// SelfLinks struct to acoomodate link inside the SelfReference
type SelfLinks struct {
	Self string `xml:"self" json:"self"`
}

// CreateResponseMessage creates an output using the parameters given and the correct marshaller
// according to the contetnType
func CreateResponseMessage(message string, code string, contentType string) ([]byte, error) {
	out := ResponseMessage{
		Status: StatusResponse{
			Message: message,
			Code:    code,
		},
	}

	output, err := MarshalContent(out, contentType, "", " ")
	return output, err
}

// CreateFailureResponseMessage creates a response message struct intance to represent all
// the errors that occured during the request
func CreateFailureResponseMessage(message string, code string, errs []ErrorResponse) ResponseMessage {

	responseMessage := ResponseMessage{
		Status: StatusResponse{
			Message: message,
			Code:    code,
		},
		Errors: errs,
	}

	return responseMessage

}

// MarshalTo takes a reponse message and uses the requested marshaller to render it into
// the desired content type
func (resp ResponseMessage) MarshalTo(contentType string) []byte {
	output, _ := MarshalContent(resp, contentType, "", " ")
	return output
}

// BadRequest is used to inform the user about malformed json body
var BadRequestSimple = ResponseMessage{
	Status: StatusResponse{
		Message: "Bad Request",
		Code:    "400",
	}}

// ErrBadRequestDetails provides an error messsage with specific details
var ErrBadRequestDetails = func(details string) ResponseMessage {
	return ResponseMessage{
		Status: StatusResponse{
			Message: "Bad Request",
			Code:    "400",
		},
		Errors: []StatusResponse{
			{Message: "Bad Request", Code: "400", Details: details},
		},
	}
}

// BadRequestBadJson is used to inform the user about malformed json body
var BadRequestBadJSON = ResponseMessage{
	Status: StatusResponse{
		Message: "Bad Request",
		Code:    "400",
		Details: "Request Body contains malformed JSON, thus rendering the Request Bad",
	}}

var BadRequestInvalidJSON = ResponseMessage{
	Status: StatusResponse{
		Message: "Bad Request",
		Code:    "400",
	},
	Errors: []StatusResponse{
		{Message: "Bad Request", Code: "400", Details: "Request Body contains malformed JSON, thus rendering the Request Bad"},
	},
}

// NotFound is used to inform the user about not found item
var NotFound = ResponseMessage{
	Status: StatusResponse{
		Message: "Not Found",
		Code:    "404",
		Details: "item with the specific ID was not found on the server",
	},
}

var ErrNotFoundQuery = ResponseMessage{
	Status: StatusResponse{
		Message: "Not Found",
		Code:    "404",
	},
	Errors: []StatusResponse{
		{Message: "Not Found", Code: "404", Details: "Specific query returned no items"},
	},
}

// NotFound is used to inform the user about not found item
var ErrNotFound = ResponseMessage{
	Status: StatusResponse{
		Message: "Not Found",
		Code:    "404",
	},
	Errors: []StatusResponse{
		{Message: "Not Found", Code: "404", Details: "item with the specific ID was not found on the server"},
	},
}

// UnauthorizedMessage is used to inform the user about incorrect api key and can be marshaled to xml and json
var UnauthorizedMessage = ResponseMessage{
	Status: StatusResponse{
		Message: "Unauthorized",
		Code:    "401",
		Details: "You need to provide a correct authentication token using the header 'x-api-key'",
	}}

// NotAcceptableContentType is used to inform the user about incorrect Accept header and can be marshaled to xml and json
var NotAcceptableContentType = ResponseMessage{
	Status: StatusResponse{
		Message: "Not Acceptable Content Type",
		Code:    "406",
		Details: "Accept header provided did not contain any valid content types. Acceptable content types are 'application/xml' and 'application/json'",
	}}

var ErrNotAcceptableContentType = func() ResponseMessage {
	return ResponseMessage{
		Status: StatusResponse{
			Message: "Not Acceptable Content Type",
			Code:    "406",
		},
		Errors: []StatusResponse{
			{Message: "Conflict", Code: "406", Details: "Accept header provided did not contain any valid content types. Acceptable content types are 'application/xml' and 'application/json'"},
		},
	}
}

// MalformedJSONInput is used to marshal a response when user json input is malformed
var MalformedJSONInput = ResponseMessage{
	Status: StatusResponse{
		Code:    "400",
		Message: "Malformated json input data",
		Details: "Check that your json input is valid",
	},
}

var ErrConflict = func(details string) ResponseMessage {
	return ResponseMessage{
		Status: StatusResponse{
			Message: "Conflict",
			Code:    "409",
		},
		Errors: []StatusResponse{
			{Message: "Conflict", Code: "409", Details: details},
		},
	}
}

// UnprocessableEntity is used to marshal a response
var UnprocessableEntity = ResponseMessage{
	Status: StatusResponse{
		Code:    "422",
		Message: "Unprocessable Entity",
	},
}

var ErrUnprocessableEntity = func(details string) ResponseMessage {
	return ResponseMessage{
		Status: StatusResponse{
			Message: "Unprocessable Entity",
			Code:    "422",
		},
		Errors: []StatusResponse{
			{Message: "Unprocessable Entity", Code: "422", Details: details},
		},
	}
}

// InternalServerErrorMessage is used to marshal a response
var InternalServerErrorMessage = ResponseMessage{
	Status: StatusResponse{
		Code:    "500",
		Message: "Internal Server Error",
	},
}

// Forbidden is used to marshal a response
var Forbidden = ResponseMessage{
	Status: StatusResponse{
		Code:    "403",
		Message: "Access to the resource is Forbidden",
	},
}
