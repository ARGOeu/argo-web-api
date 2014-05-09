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

package poems

import (
	"api/utils/config"
	"api/utils/mongo"
	"encoding/xml"
	"net/http"

)

type ApiPOEM struct{
	Poem string `bson:"p"`
}

type Poem struct{
	Poem string `xml:"profile,attr"`
}

type POEMxml struct{
	XMLName      xml.Name `xml:"POEM_List"`
	Poem	  	 []*Poem 
}

func ReadPoems(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte {

	results := []ApiPOEM{}
	output  := &POEMxml{}

	session := mongo.OpenSession(cfg)

	err := mongo.Find(session, "AR", "poem_list", nil, "p", &results)
	
	for _, row := range results{
		p := &Poem{}
		p.Poem = row.Poem
		output.Poem =append(output.Poem,p)
	}
	
	answer, err := xml.MarshalIndent(output, "", " ")
	
	if err != nil {
		panic(err)
	}

	mongo.CloseSession(session)

	return []byte("<root>" + string(answer) + "</root>")
}