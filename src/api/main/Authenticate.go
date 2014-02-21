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

package main

import (
	"bufio"
	"os"
	"strings"
	"net/http"
)

var m=make(map[string]string,128)

func Authenticate(h http.Header) bool{
	user:=h.Get("x-api-requestor")//Suggested username: host FQDN passed with the http request headers
	key:=h.Get("x-api-key")//Api key shared with webUI
	file,err:=os.Open("/root/api_password_file")//Username-passwords stored into file. Future work: store password values hashed into mongoDB
	if err!=nil{
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
			line:=strings.Split(scanner.Text() ,":")//all values stored inside a map structure
			m[line[0]]=line[1]
		}	
	if m[user]==key{
		return true
	}else{
		return false
	}
}