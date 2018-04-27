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

import	`sort`

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

// Json returns the DmnList object as a JSON byte array.
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

// Sort sorts the DmnList by rules defined in byDmnInfo.Less().
func (this *DmnList) Sort() {
	sort.Stable(byDmnInfo(*this))
}

// ------------------------------------------------------------------------
// byDmnInfo.
// ------------------------------------------------------------------------

// byDmnInfo is a derivative object used in sorting.
type byDmnInfo DmnList

// ------------------------------------------------------------------------
// byDmnInfo Methods.
// ------------------------------------------------------------------------

// Len returns the number of elements in the slice.
func (this byDmnInfo) Len() int {
	return len(this)
}

// Swap exchanges the values of two elements.
func (this byDmnInfo) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// Less defines the rules for soring the elements.
func (this byDmnInfo) Less(i, j int) bool {

	if this[i].Name != this[j].Name {
		return this[i].Name < this[j].Name
	} else if this[i].Key != this[j].Key {
		return this[i].Key < this[j].Key
	} else if this[i].Version != this[j].Version {
		return this[i].Version < this[j].Version
	} else {
		return this[i].Id < this[j].Id
	}
}
