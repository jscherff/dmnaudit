// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
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

package main

import `github.com/jscherff/dmnsdk/model`

func diff(di *model.DmnInfo, dmn1, dmn2 *model.Dmn) {

	if de, err := model.NewDmnElements(dmn1); err != nil {
		failure.Add(elmFailure.Fmt(`Env1`, di.Key, di.Version, err))
	} else if err := de.Compare(dmn2); err != nil {
		failure.Add(elmFailure.Fmt(`Env2`, di.Key, di.Version, err))
	} else {
		for _, key := range de.SortedKeys() {
			switch de[key] {
			case 1:
				warning.Add(cmpDetails.Fmt(`Env1`, key))
			case -1:
				warning.Add(cmpDetails.Fmt(`Env2`, key))
			case 0:
				if !*fVerbose {break}
				warning.Add(cmpDetails.Fmt(`Both`, key))
			}
		}
	}
}

