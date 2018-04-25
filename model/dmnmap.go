// Copyright 2018 John Scherff
//
// Licensed under the Apache License, version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import `fmt`

// ------------------------------------------------------------------------
// DmnMap.
// ------------------------------------------------------------------------

// DmnMap is a collection of DmnInfo objects indexed by DMN Key and Version.
type DmnMap map[string]map[int]*DmnInfo

// ------------------------------------------------------------------------
// DmnMap Methods.
// ------------------------------------------------------------------------

// DmnInfo returns a DmnInfo object given its key and version.
func (this DmnMap) DmnInfo(key string, ver int) (*DmnInfo, error) {

	if di, ok := this[key][ver]; !ok {
		return nil, fmt.Errorf(`key %s version %d not found`, key, ver)
	} else {
		return di, nil
	}
}

// DmnId returns the DMN ID given the DMN Key and Version.
func (this DmnMap) DmnId(key string, ver int) (string, error) {

	if di, err := this.DmnInfo(key, ver); err != nil {
		return ``, err
	} else {
		return di.Id, nil
	}
}
