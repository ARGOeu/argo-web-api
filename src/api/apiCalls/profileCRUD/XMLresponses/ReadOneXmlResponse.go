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

package XMLresponses

import (
	"encoding/xml"
)

type MongoProfile struct {
	Name  string     "p"
	Group [][]string "g"
}

func ReadOneXmlResponse(results []MongoProfile) ([]byte, error) {

	type Groups struct {
		XMLName         xml.Name `xml:"Group"`
		Service_flavors []string `xml:"service_flavor"`
	}

	type Profile struct {
		XMLName xml.Name `xml:"Profile"`
		Name    string   `xml:"name,attr"`
		Groups  []Groups
	}

	type Root struct {
		XMLName  xml.Name `xml:"root"`
		Profiles []Profile
	}

	v := &Root{}

	for key, result := range results {
		v.Profiles = append(v.Profiles,
			Profile{
				Name: result.Name,
			})
		for _, result2 := range result.Group {
			v.Profiles[key].Groups = append(v.Profiles[key].Groups,
				Groups{
					Service_flavors: result2,
				})
		}
	}

	output, err := xml.MarshalIndent(v, " ", "  ")
	return output, err
}
