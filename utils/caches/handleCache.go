/*
 * Copyright (c) 2015 GRNET S.A.
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

package caches

import (
	"fmt"

	"github.com/ARGOeu/argo-web-api/Godeps/_workspace/src/github.com/ARGOeu/go-lru-cache"
	"github.com/ARGOeu/argo-web-api/utils/config"
)

var httpcache *cache.LRUCache

type mystring string

func (s mystring) Size() int {
	return len(s)
}

//TODO Add multitenant support on cache util.
func HitCache(name string, input interface{}, cfg config.Config) (bool, []byte) {
	output := []byte(nil)
	found := false
	if cfg.Server.Cache == true {
		out, found := httpcache.Get(name + fmt.Sprint(input))
		if found {
			output = []byte(fmt.Sprint(out))
		}
	}
	return found, output
}

//TODO Add multitenant support on cache util.
func WriteCache(name string, input interface{}, output interface{}, cfg config.Config) bool {

	if cfg.Server.Cache == true {
		httpcache.Set(name+fmt.Sprint(input), output.(mystring))
	}
	return true
}

func ResetCache() {
	httpcache.Clear()
}
