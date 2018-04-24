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
	`fmt`
	`sort`
)

// DmnElement captures DMN XML elements.
type DmnElement struct {
	Tag			string
	Id			string
	Property		string
	Value			string
}

func (this DmnElement) String() (string) {
	return fmt.Sprintf(`<%s id="%s"><%s>%s<%[3]s/><%[1]s/>`,
		this.Tag, this.Id, this.Property, this.Value,
	)
}

type byDmnElement []DmnElement

func (this byDmnElement) Len() int {
	return len(this)
}

func (this byDmnElement) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this byDmnElement) Less(i, j int) bool {

	if this[i].Tag != this[j].Tag {
		return this[i].Tag < this[j].Tag
	} else if this[i].Id != this[j].Id {
		return this[i].Id < this[j].Id
	} else if this[i].Property != this[j].Property {
		return this[i].Property < this[j].Property
	} else {
		return this[i].Value < this[j].Value
	}
}

// DmnElements is a collection of DmnElement objects.
type DmnElements map[DmnElement]int

// NewDmnElements creates a collection of DmnElement objects from a DMN.
func NewDmnElements(t interface{}) (DmnElements, error) {

	this := make(DmnElements)

	if err := this.Load(t); err != nil {
		return nil, err
	}

	return this, nil
}

func (this DmnElements) Keys() (keys []DmnElement) {
	for key := range this {
		keys = append(keys, key)
	}
	return keys
}

func (this DmnElements) SortedKeys() (keys []DmnElement) {
	keys = this.Keys()
	sort.Sort(byDmnElement(keys))
	return keys
}

func (this DmnElements) Compare(t interface{}) (error) {

	if err := this.load(t, -1); err != nil {
		return err
	}

	return nil
}

// Load imports a DMN into a data structure that allows comparison
// of DMN elements between objects.
func (this DmnElements) Load(t interface{}) (error) {
	return this.load(t, 1)
}

// load imports a DMN into a data structure that allows comparison
// of DMN elements between objects.
func (this DmnElements) load(t interface{}, inc int) (error) {

	var tag, id string

	switch obj := t.(type) {

	case *Dmn:

		if em, err := toMap(t); err != nil {
			return err
		} else {
			return this.load(em, inc)
		}

	case map[string]interface{}:

		tm := make(map[string]interface{})

		for name, value := range obj {

			switch obj := value.(type) {

			case string:

				if name == `id` {
					id = obj
				} else {
					tm[name] = value
				}

			case map[string]interface{}:

				if name == `xmlName` {
					tag = obj[`Local`].(string)
				} else if err := this.load(obj, inc); err != nil {
					return err
				}

			default:

				if err := this.load(obj, inc); err != nil {
					return err
				}
			}
		}

		if tag == `` || id == `` {
			return fmt.Errorf(`missing Tag and Id`)
		}

		for name, value := range tm {
			el := DmnElement{tag, id, name, value.(string)}
			this[el] += inc
		}

	case []interface{}:

		for _, value := range obj {
			if err := this.load(value, inc); err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf(`unsupported type %T`, t)
	}

	return nil
}
