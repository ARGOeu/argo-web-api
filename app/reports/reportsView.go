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

package reports

import "github.com/ARGOeu/argo-web-api/respond"

// SubmitSuccesful marshals a response struct when a report is successfully inserted in the
// database
func SubmitSuccesful(inserted MongoInterface, contentType string, link string) ([]byte, error) {
	docRoot := respond.ResponseMessage{
		Status: respond.StatusResponse{
			Message: "Successfully Created Report",
			Code:    "201",
		},
		Data: respond.SelfReference{
			UUID:  inserted.UUID,
			Links: respond.SelfLinks{Self: link},
		},
	}
	output, err := respond.MarshalContent(docRoot, contentType, "", " ")
	return output, err
}

// ReportNotFound consructs marshals a response struct when the requested
// report is not found
func ReportNotFound(contentType string) ([]byte, error) {
	docRoot := respond.ResponseMessage{
		Status: respond.StatusResponse{
			Message: "Report Not Found",
			Code:    "404",
		},
	}
	output, err := respond.MarshalContent(docRoot, contentType, "", " ")
	return output, err
}

func createView(results interface{}, format string) ([]byte, error) {
	docRoot := &respond.ResponseMessage{
		Status: respond.StatusResponse{
			Message: "Success",
			Code:    "200",
		},
	}

	docRoot.Data = &results
	output, err := respond.MarshalContent(docRoot, format, "", " ")
	// output, err := xml.MarshalIndent(docRoot, "", " ")
	return output, err
}
