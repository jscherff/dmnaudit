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

import (
	"encoding/xml"
	"fmt"
)

// ==============================================================================
// See https://docs.camunda.org/manual/7.4/reference/dmn11/decision-table/
// ==============================================================================


// -----------------------
// Dmn Object and Methods.
// -----------------------

// Dmn contains Decision Model and Notation definitions defining the 
// Decision and DecisionTable.
type Dmn struct {
	XMLName           xml.Name            `json:"xmlName"`
	Xmlns             string              `xml:"xmlns,attr" json:"xmlns"`
	Id                string              `xml:"id,attr" json:"id"`
	Name              string              `xml:"name,attr" json:"name"`
	Namespace         string              `xml:"namespace,attr" json:"namespace"`
	Decision          *Decision           `xml:"decision,child" json:"decision"`
}

// NewDmn creates and loads a new Dmn object from a JSON source.
func NewDmn(src interface{}) (*Dmn, error) {
	this := new(Dmn)
	err := this.Load(src)
	return this, err
}

// Load unmarshals JSON from Reader, url, file or string into an object.
func (this *Dmn) Load(src interface{}) error {
	return load(this, src, `xml`)
}

// Json marshals an object into a JSON byte array.
func (this *Dmn) Json() ([]byte, error) {
	return toJson(this)
}

// A DecisionTable is decision logic which can be depicted as a table in
// DMN 1.1. It consists of inputs, outputs and rules and is represented
// by a <decisionTable> element inside a <decision> element.

// The name describes the decision for which the decision table provides
// the decision logic. It is set as the name attribute on the decision
// element. The id is the technical identifier of the decision. It is
// set in the id attribute on the decision element.

type Decision struct {
	XMLName           xml.Name            `json:"xmlName"`
	Id                string              `xml:"id,attr" json:"id"`
	Name              string              `xml:"name,attr" json:"name"`
	DecisionTable     *DecisionTable      `xml:"decisionTable,child" json:"decisionTable"`
}

type DecisionTable struct {
	XMLName           xml.Name            `json:"xmlName"`
	Id                string              `xml:"id,attr" json:"id"`
	HitPolicy         string              `xml:"hitPolicy,attr" json:"hitPolicy"`
	Inputs            []*Input            `xml:"input,child" json:"input"`
	Outputs           []*Output           `xml:"output,child" json:"output"`
	Rules             []*Rule             `xml:"rule,child" json:"rule"`
}

// A decision table can have one or more inputs, also called input
// clauses. An input clause defines the id, label, expression and type
// of a decision table input. An input clause is represented by an input
// element inside a decisionTable XML element.

// The input id is an unique identifier of the decision table input.
// It is used by the Camunda BPMN platform to reference the input in the
// history of evaluated decisions. Therefore it is required by the Camunda
// DMN engine. It is set as the id attribute of the input XML element.

// An input label is a short description of the input. It is set on the
// input XML element in the label attribute. Note that the label is not
// required but recommended since it helps to understand the decision.

type Input struct {
	XMLName           xml.Name            `json:"xmlName"`
	Id                string              `xml:"id,attr" json:"id"`
	Label             string              `xml:"label,attr" json:"label"`
	InputExpressions  []*InputExpression  `xml:"inputExpression,child" json:"inputExpression"`
}

// An input expression specifies how the value of the input clause is
// generated. It is an expression which will be evaluated by the DMN
// engine. It is usually simple and references a variable which is
// available during the evaluation. The expression is set inside a text
// element that is a child of the inputExpression XML element.

// The type of the input clause can be specified by the typeRef attribute
// on the inputExpression XML element. After the input expression is
// evaluated by the DMN engine it converts the result to the specified
// type.

// The expression language of the input expression can be specified by
// the expressionLanguage attribute on the inputExpression XML element.
// If no expression language is set then the global expression language
// is used which is set on the definitions XML element. In case no global
// expression language is set, the default expression language is used
// instead. The default expression language for input expressions is JUEL.

type InputExpression struct {
	XMLName           xml.Name            `json:"xmlName"`
	Id                string              `xml:"id,attr" json:"id"`
	TypeRef           string              `xml:"typeRef,attr" json:"typeRef"`
	Text              string              `xml:"text" json:"text"`
}

// A decision table can have one or more output, also called output clauses.
// An output clause defines the id, label, name and type of a decision table
// output. An output clause is represented by an output element inside a
// decisionTable XML element.

// The output id is an unique identifier of the decision table output. It
// is used by the Camunda BPMN platform to reference the output in the
// history of evaluated decisions. Therefore it is required by the Camunda
// DMN engine. It is set as the id attribute of the output XML element.

// An output label is a short description of the output. It is set on the
// output XML element in the label attribute. Note that the label is not
// required but recommended since it helps to understand the decision.

// The name of the output is used to reference the value of the output in
// the decision table result. It is specified by the name attribute on the
// output XML element. If the decision table has more than one output then
// all outputs must have an unique name.

// The type of the output clause can be specified by the typeRef attribute
// on the output XML element. After an output entry is evaluated by the DMN
// engine it converts the result to the specified type. Note that the type
// is not required but recommended since it provides a type safety of the
// output values. Additionally, the type can be used to transform the output
// value into another type. For example, transform the output value 80% of
// type String into a Double using a custom data type.

type Output struct {
	XMLName           xml.Name            `json:"xmlName"`
	Id                string              `xml:"id,attr" json:"id"`
	Label             string              `xml:"label,attr" json:"label"`
	Name              string              `xml:"name,attr" json:"name"`
	TypeRef           string              `xml:"typeRef,attr" json:"typeRef"`
}

// A decision table can have one or more rules. Each rule contains input
// and output entries. The input entries are the condition and the output
// entries the conclusion of the rule. If each input entry (condition) is
// satisfied then the rule is satisfied and the decision result contains
// the output entries (conclusion) of this rule. A rule is represented by
// a rule element inside a decisionTable XML element.

type Rule struct {
	XMLName           xml.Name            `json:"xmlName"`
	Id                string              `xml:"id,attr" json:"id"`
	InputEntries      []*InputEntry       `xml:"inputEntry,child" json:"inputEntry"`
	OutputEntries     []*OutputEntry      `xml:"outputEntry,child" json:"outputEntry"`
}

// A rule can have one or more input entries which are the conditions of
// the rule. Each input entry contains an expression in a text element as
// child of an inputEntry XML element. The input entry is satisfied when
// the evaluated expression returns true. In case an input entry is
// irrelevant for a rule, the expression is empty which is always satisfied.

// The expression language of the input entry can be specified by the
// expressionLanguage attribute on the inputEntry XML element. If no
// expression language is set then the global expression language is
// used which is set on the definitions XML element. In case no global
// expression language is set, the default expression language is used
// instead.

type InputEntry struct {
	XMLName           xml.Name            `json:"xmlName"`
	Id                string              `xml:"id,attr" json:"id"`
	Text               string             `xml:"text" json:"text"`
}

// A rule can have one or more output entries which are the conclusions
// of the rule. Each output entry contains an expression in a text element
// as child of an outputEntry XML element. If the output entry is empty then
// the output is ignored and not part of the decision table result.

// The expression language of the expression can be specified by the
// expressionLanguage attribute on the outputEntry XML element. If no
// expression language is set then the global expression language is
// used which is set on the definitions XML element. In case no global
// expression language is set, the default expression language is used
// instead.

// A rule can be annotated with a description that provides additional
// information. The description text is set inside the description XML
// element.

type OutputEntry struct {
	XMLName           xml.Name            `json:"xmlName"`
	Id                string              `xml:"id,attr" json:"id"`
	Description       string              `xml:"description" json:"description"`
	Text              string              `xml:"text" json:"text"`
}

// -----------------------------------------------------------------------
// DmnInfo Object and Methods.
// -----------------------------------------------------------------------

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


// -----------------------------------------------------------------------
// DmnList Object and Methods.
// -----------------------------------------------------------------------

// DmnList is a collection of DmnsInfo objects.
type DmnList []*DmnInfo

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

// --------------------------
// DmnXml Object and Methods.
// --------------------------

// NewDmnXml creates and loads a new DmnXml object from a JSON source.
func NewDmnXml(src interface{}) (*DmnXml, error) {
	this := new(DmnXml)
	err := this.Load(src)
	return this, err
}

// DmnXml is the DMN Decision Definition in XML format.
type DmnXml struct {
	Id                string              `json:"id"`
	DmnXml            string              `json:"dmnXml"`
}

// Load unmarshals JSON from a Reader, url, file, or string to an object.
func (this *DmnXml) Load(src interface{}) error {
	return load(this, src, `json`)
}

// Json marshals an object into a JSON byte array.
func (this *DmnXml) Json() ([]byte, error) {
	return toJson(this)
}

// String implements the Stringer interface for DmnXml.
func (this *DmnXml) String() (string) {
	return this.DmnXml
}

// --------------------------
// DmnMap Object and Methods.
// --------------------------

// DmnMap is a collection of DmnInfo objects indexed by DMN Key and Version.
type DmnMap map[string]map[int]*DmnInfo

// Get returns a DmnInfo object given its key and version.
func (this DmnMap) Info(key string, ver int) (*DmnInfo, error) {

	if di, ok := this[key][ver]; !ok {
		return nil, fmt.Errorf(`key %s version %d not found`, key, ver)
	} else {
		return di, nil
	}
}

// GetId returns the DMN ID given the DMN Key and Version.
func (this DmnMap) Id(key string, ver int) (string, error) {

	if di, err := this.Info(key, ver); err != nil {
		return ``, err
	} else {
		return di.Id, nil
	}
}
