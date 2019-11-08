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

package operationsProfiles

import "sort"

// OpsProfile to retrieve and insert operationsProfiles in mongo
type OpsProfile struct {
	ID          string        `bson:"id" json:"id"`
	DateInt     int           `bson:"date_integer" json:"-"`
	Date        string        `bson:"date" json:"date"`
	Name        string        `bson:"name" json:"name"`
	AvailStates []string      `bson:"available_states" json:"available_states"`
	Defaults    DefaultStates `bson:"defaults" json:"defaults"`
	Operations  []Operation   `bson:"operations" json:"operations"`
}

// DefaultStates struct to represent defaults states
type DefaultStates struct {
	Down    string `bson:"down" json:"down"`
	Missing string `bson:"missing" json:"missing"`
	Unknown string `bson:"unknown" json:"unknown"`
}

// Operation struct to represent an operation
type Operation struct {
	Name       string      `bson:"name" json:"name"`
	TruthTable []Statement `bson:"truth_table" json:"truth_table"`
}

// Statement holds an operation statement expressed in the form of A {op} B -> X
type Statement struct {
	A string `bson:"a" json:"a"`
	B string `bson:"b" json:"b"`
	X string `bson:"x" json:"x"`
}

// SelfReference to hold links and id
type SelfReference struct {
	ID    string `json:"id" bson:"id,omitempty"`
	Links Links  `json:"links"`
}

// Links struct to hold links
type Links struct {
	Self string `json:"self"`
}

func (oprof *OpsProfile) hasState(state string) bool {
	for _, item := range oprof.AvailStates {
		if item == state {
			return true
		}
	}
	return false
}

// validateStates checks all state references for undeclared states
func (oprof *OpsProfile) validateStates() []string {
	// check default states
	var errList []string

	if !(oprof.hasState(oprof.Defaults.Down)) {
		errList = append(errList, "Default Down State: "+oprof.Defaults.Down+" not in available States")
	}
	if !(oprof.hasState(oprof.Defaults.Missing)) {
		errList = append(errList, "Default Missing State: "+oprof.Defaults.Missing+" not in available States")
	}
	if !(oprof.hasState(oprof.Defaults.Unknown)) {
		errList = append(errList, "Default Unknown State: "+oprof.Defaults.Unknown+" not in available States")
	}
	// check operations
	for _, op := range oprof.Operations {
		for _, st := range op.TruthTable {
			if !(oprof.hasState(st.A)) {
				errList = append(errList, "In Operation: "+op.Name+", statement member a: "+st.A+" contains undeclared state")
			}
			if !(oprof.hasState(st.B)) {
				errList = append(errList, "In Operation: "+op.Name+", statement member b: "+st.B+" contains undeclared state")
			}
			if !(oprof.hasState(st.X)) {
				errList = append(errList, "In Operation: "+op.Name+", statement member x: "+st.X+" contains undeclared state")
			}
		}
	}

	return errList
}

// validateMentions checks if we have enough state mentions in the truth table to accomodate all cases
func (oprof *OpsProfile) validateMentions() []string {

	var errList []string

	counters := make(map[string]int)
	// threshold of mentions for each element = number of elements + 1
	thold := len(oprof.AvailStates) - 1
	// for each operation
	for _, op := range oprof.Operations {
		// init counter map
		for _, state := range oprof.AvailStates {
			counters[state] = 0
		}

		// for all statements in truth table
		for _, st := range op.TruthTable {
			if oprof.hasState(st.A) {
				counters[st.A]++
			}
			if oprof.hasState(st.B) {
				counters[st.B]++
			}

		}

		// Check counters if contain mentions >= threshold
		var sorted []string
		for key := range counters {
			sorted = append(sorted, key)
		}
		// sort keys
		sort.Strings(sorted)

		for _, key := range sorted {
			if counters[key] < thold {
				errList = append(errList, "Not enough mentions of state:"+key+" in operation: "+op.Name)
			}
		}
	}

	return errList
}

// validateDuplicates checks if we have duplicate states
func (oprof *OpsProfile) validateDuplicates() []string {

	var errList []string

	// check duplicate in states
	for i := 0; i < len(oprof.AvailStates); i++ {
		for j := i + 1; j < len(oprof.AvailStates); j++ {
			if oprof.AvailStates[i] == oprof.AvailStates[j] {
				errList = append(errList, "State:"+oprof.AvailStates[i]+" is duplicated")
			}
		}
	}

	// check duplicate in operations
	for i := 0; i < len(oprof.Operations); i++ {
		for j := i + 1; j < len(oprof.Operations); j++ {
			if oprof.Operations[i].Name == oprof.Operations[j].Name {
				errList = append(errList, "Operation:"+oprof.Operations[i].Name+" is duplicated")
			}
		}
	}

	return errList
}
