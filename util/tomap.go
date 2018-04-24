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

package util

import (
	`fmt`
	`reflect`
	`strings`
)

// ToMap creates a map from a struct, including only fields that match a tag.
func ToMap(i interface{}, m map[string]interface{}, tid string) (error) {

	v := reflect.ValueOf(i).Elem()

	if v.Type().Kind() != reflect.Struct {
		return fmt.Errorf(`kind %q is not struct`, v.Type().Kind())
	}

	for i := 0; i < v.NumField(); i++ {

		f := v.Field(i)
		t := v.Type().Field(i)

		if !f.IsValid() || !f.CanAddr() || !f.CanInterface() {
			continue
		}

		if tag, ok := t.Tag.Lookup(tid); !ok {
			continue
		} else {

			tval := strings.Split(tag, `,`)[0]

			switch tval {
			case `-`:
				continue
			case ``:
				m[t.Name] = f.Interface()
			default:
				m[tval] = f.Interface()
			}
		}
	}

	return nil
}
