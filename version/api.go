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

package version

import (
	"fmt"

	"net/http"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/gorilla/mux"
)

// HandleSubrouter for api access to version infomation
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {
	s = respond.PrepAppRoutes(s, confhandler, appRoutesV2)
}

var appRoutesV2 = []respond.AppRoutes{
	{"version.list", "GET", "/version", ListVersion},
	{"version.options", "OPTIONS", "/version", Options},
}

// Version struct holds version information about the binary build
type Version struct {
	Release   string `xml:"release" json:"release"`
	Commit    string `xml:"commit" json:"commit"`
	BuildTime string `xml:"build_time" json:"build_time"`
	GO        string `xml:"golang" json:"golang"`
	Compiler  string `xml:"compiler" json:"compiler"`
	OS        string `xml:"os" json:"os"`
	Arch      string `xml:"architecture" json:"architecture"`
}

// ListVersion displays version information about the service
func ListVersion(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	contentType, _ = respond.ParseAcceptHeader(r)
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	v := Version{
		Release:   Release,
		Commit:    Commit,
		BuildTime: BuildTime,
		GO:        GO,
		Compiler:  Compiler,
		OS:        OS,
		Arch:      Arch,
	}

	output, err = respond.MarshalContent(v, contentType, "", " ")

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	return code, h, output, err
}

func Options(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/plain"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	h.Set("Allow", fmt.Sprintf("GET, OPTIONS"))
	return code, h, output, err

}
