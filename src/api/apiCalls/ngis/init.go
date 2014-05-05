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

package ngis

var customForm []string

type ApiNgiAvailabilityInProfileInput struct {
	// mandatory values
	start_time           string // UTC time in W3C format
	end_time             string // UTC time in W3C format
	availability_profile string //availability profile
	// optional values
	granularity    string //availability period; possible values: `DAILY`, MONTHLY`
	infrastructure string //infrastructure name
	production     string //production or not
	monitored      string //yes or no
	certification  string //certification status
	//format    string   // default XML; possible values are: XML, JSON
	group_name []string // site name; may appear more than once
}

type ApiNgiAvailabilityInProfileOutput struct {
	Date         string  "dt"
	Namespace    string  "ns"
	Profile      string  "p"
	Ngi          string  "n"
	Availability float64 "a"
	Reliability  float64 "r"
}

func init() {
	customForm = []string{"20060102", "2006-01-02"} //{"Format that is returned by the database" , "Format that will be used in the generated report"}
}
