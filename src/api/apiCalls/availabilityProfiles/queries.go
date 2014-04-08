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

package availabilityProfiles

import (
	"labix.org/v2/mgo/bson"
)

func prepareFilter(input ApiAPSearch) bson.M {

	filter := bson.M{
		"name":      bson.M{"$in": input.Name},
		"namespace": bson.M{"$in": input.Namespace},
	}

	return filter
}

func createOne(input ApiAPInput) bson.M {
	query := bson.M{
		"name":      input.Name,
		"namespace": input.Namespace,
		"groups":    input.Groups,
		"poem":      input.Poem,
	}
	return query
}

func readOne(input ApiAPSearch) bson.M {
	filter := prepareFilter(input)
	return filter
}

// func deleteOne(input ApiAPInput) bson.M {
// 	filter := prepareFilter(input)
// 	return filter
// }
