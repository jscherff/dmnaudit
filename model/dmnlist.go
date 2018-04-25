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

// ------------------------------------------------------------------------
// DmnList.
// ------------------------------------------------------------------------

// DmnList is a collection of DmnsInfo objects.
type DmnList []*DmnInfo

// ------------------------------------------------------------------------
// DmnList Methods.
// ------------------------------------------------------------------------

// NewDmnList creates and loads a new DmnList object from a JSON source.
func NewDmnList(src interface{}) (*DmnList, error) {
	this := new(DmnList)
	err := this.Load(src)
	return this, err
}

// Load unmarshals JSON from a Reader, url, file, or string to an object.
func (this *DmnList) Load(src interface{}) error {
	return load(this, src, `json`)
}

// Json marshals an object into a JSON byte array.
func (this *DmnList) Json() ([]byte, error) {
	return toJson(this)
}

// Map creates a map of Dmns indexed by key and version number.
func (this *DmnList) Map() (DmnMap, error) {

	dm := make(DmnMap)

	for _, di := range *this {

		if dm[di.Key] == nil {
			dm[di.Key] = make(map[int]*DmnInfo)
		}

		dm[di.Key][di.Version] = di
	}

	return dm, nil
}
