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
// DmnInfo.
// ------------------------------------------------------------------------

// DmnInfo contains DMN metadata which can be used to retrieve other data,
// such as the DMN XML describing the DMN.
type DmnInfo struct {
	Id                string              `json:"id"`
	Key               string              `json:"key"`
	Category          string              `json:"category"`
	Name              string              `json:"name"`
	Version           int                 `json:"version"`
	Resource          string              `json:"resource"`
	DeploymentId      string              `json:"deploymentId"`
	TenantId          string              `json:"tenantId"`
	DecisionReqDefId  string              `json:"decisionRequirementsDmnId"`
	DecisionReqDefKey string              `json:"decisionRequirementsDmnKey"`
	HistoryTtl        string              `json:"historyTimeToLive"`
	DmnXml            string              `json:"dmnXml"`
}

// ------------------------------------------------------------------------
// DmnInfo Methods.
// ------------------------------------------------------------------------

// NewDmnInfo creates and loads a new DmnInfo object from a JSON source.
func NewDmnInfo(src interface{}) (*DmnInfo, error) {
	this := new(DmnInfo)
	err := this.Load(src)
	return this, err
}

// Load unmarshals JSON from a Reader, url, file, or string to an object.
func (this *DmnInfo) Load(src interface{}) error {
	return load(this, src, `json`)
}
