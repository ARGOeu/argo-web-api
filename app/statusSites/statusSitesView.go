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

package statusSites

import "encoding/xml"

func createView(results []StatusSitesOutput, input StatusSitesInput) ([]byte, error) {

	docRoot := &ReadRoot{}

	if len(results) == 0 {
		output, err := xml.MarshalIndent(docRoot, " ", "  ")
		return output, err
	}

	profile := &Profile{}
	profile.Name = input.profile
	vo := &Group{}
	vo.Type = "vo"
	vo.Name = input.vo

	prevSite := ""
	prevRoc := ""

	var pp_Site *Site
	var pp_Roc *Group

	for _, row := range results {

		// filter by profile
		if row.Profile != input.profile {
			continue
		}

		if row.Roc != prevRoc && row.Roc != "" {
			roc := &Group{}
			roc.Name = row.Roc
			roc.Type = "ngi"
			vo.Groups = append(vo.Groups, roc)
			prevRoc = roc.Name
			pp_Roc = roc
		}

		if row.Site != prevSite {

			site := &Site{}
			//Add the prev status as the firstone

			site.Name = row.Site

			pp_Roc.Sites = append(pp_Roc.Sites, site)
			prevSite = row.Site
			pp_Site = site

			status := &Status{}
			status.Timestamp = input.start_time
			status.Status = row.P_status
			pp_Site.Timeline = append(pp_Site.Timeline, status)

		} else {
			status := &Status{}
			status.Timestamp = row.Timestamp
			status.Status = row.Status
			pp_Site.Timeline = append(pp_Site.Timeline, status)
		}

	}

	profile.Groups = append(profile.Groups, vo)
	docRoot.Profile = profile

	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err

}
