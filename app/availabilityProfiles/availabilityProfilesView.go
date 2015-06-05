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

package availabilityProfiles

import "encoding/xml"

func createView(results []AvailabilityProfileOutput) ([]byte, error) {

	metricProfileName := ""

	docRoot := &ReadRoot{}
	for _, row := range results {

		// If metricprofiles array doesn't contain items add empty string to metricprofiles attribute
		if len(row.MetricProfiles) > 0 {
			metricProfileName = row.MetricProfiles[0]
		} else {
			metricProfileName = ""
		}

		profile := &Profile{
			ID:            row.ID.Hex(),
			Name:          row.Name,
			Namespace:     row.Namespace,
			MetricProfile: metricProfileName,
		}
		and := &And{}
		docRoot.Profile = append(docRoot.Profile, profile)
		for _, group := range row.Groups {
			or := &Or{}
			for sf, op := range group.Services {
				group := &Group{
					ServiceFlavor:    sf,
					ServiceOperation: op,
				}
				or.Group = append(or.Group, group)
			}
			and.Or = append(and.Or, or)
		}
		profile.And = and
	}
	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err

}

func messageXML(answer string) ([]byte, error) {
	docRoot := &Message{}
	docRoot.Message = answer
	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err
}
