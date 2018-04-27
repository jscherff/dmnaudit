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

import `flag`

var (
	fSvcUrl1 = flag.String(`url1`, ``, "Use service at `http[s]://<hostname>[:<port>]`")
	fSvcUrl2 = flag.String(`url2`, ``, "Use service at `http[s]://<hostname>[:<port>]`")
	fOutFile = flag.String(`file`, ``, "Store results in file `<file>`")
	fWarning = flag.Bool(`warning`, true, "Show warning message when DMNs do not match")
	fFailure = flag.Bool(`failure`, true, "Show failure message when DMN not found")
	fSuccess = flag.Bool(`success`, false, "Show success message when DMNs match")
	fDetails = flag.Bool(`details`, false, "Show detailed differences between DMN elements")
	fVerbose = flag.Bool(`verbose`, false, "Show matching DMN elements along with differences")
)
