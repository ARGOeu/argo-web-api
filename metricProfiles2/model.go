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

package metricProfiles2

// MongoInterface to retrieve and insert metricProfiles in mongo
type MongoInterface struct {
	UUID     string    `bson:"uuid" json:"uuid"`
	Name     string    `bson:"name" json:"name"`
	Services []Service `bson:"services" json:"services"`
}

// Service struct to represent services with their metrics
type Service struct {
	Service string   `bson:"service" json:"service"`
	Metrics []string `bson:"metrics" json:"metrics"`
}

type SelfReference struct {
	UUID  string `json:"uuid" bson:"uuid,omitempty"`
	Links Links  `json:"links"`
}

type Links struct {
	Self string `json:"self"`
}
