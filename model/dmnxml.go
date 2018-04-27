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
// DmnXml.
// ------------------------------------------------------------------------

// DmnXml is the DMN Decision Definition in XML format.
type DmnXml struct {
	Id                string              `json:"id"`
	DmnXml            string              `json:"dmnXml"`
}

// ------------------------------------------------------------------------
// DmnXml Methods.
// ------------------------------------------------------------------------

// NewDmnXml creates and loads a new DmnXml object from a JSON source.
func NewDmnXml(src interface{}) (*DmnXml, error) {
	this := new(DmnXml)
	err := this.Load(src)
	return this, err
}

// Load unmarshals JSON from a Reader, url, file, or string to an object.
func (this *DmnXml) Load(src interface{}) error {
	return load(this, src, `json`)
}

// Json returns the DmnXml object as a JSON byte array.
func (this *DmnXml) Json() ([]byte, error) {
	return toJson(this)
}

// Xml returns DmnXml as an XML byte array.
func (this *DmnXml) Xml() ([]byte, error) {
	return []byte(this.DmnXml), nil
}

// String implements the Stringer interface for DmnXml.
func (this *DmnXml) String() (string) {
	return this.DmnXml
}
