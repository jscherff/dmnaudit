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

import (
	`flag`
	`fmt`
	`log`
	//`io`
	`os`
	`reflect`
	`github.com/jscherff/dmnsdk/api`
	`github.com/jscherff/dmnsdk/model`
)

const (
	esbeapPRD = `http://esbeap.24hourfit.com:8180`
	esbeapQA = `http://esbeap-qa.24hourfit.com:8180`
	esbeapDEV= `http://esbeap-dev.24hourfit.com:8180`

	dmnFailFmt = `[%s] could not get %s DMN key %s ver %d: %v`
	dmnCompFmt = `[%s] Env1 and Env2 DMN key %s ver %d: %s`
	elmFailFmt = `[%s] could not get %s DMN key %s ver %d elements: %v`
)

var (
	fSvcUrl1 = flag.String(`url1`, ``, "Use service at `http[s]://<hostname>[:<port>]`")
	fSvcUrl2 = flag.String(`url2`, ``, "Use service at `http[s]://<hostname>[:<port>]`")
	fOutFile = flag.String(`file`, ``, "Store results in file `<file>`")
	fSuccess = flag.Bool(`success`, false, "Show success message when DMNs match")
	fWarning = flag.Bool(`warning`, false, "Show warning message when DMNs do not match")
	fFailure = flag.Bool(`failure`, false, "Show failure message when DMN not found")
	fDetails = flag.Bool(`details`, false, "Show detailed differences between DMN elements")
	fVerbose = flag.Bool(`verbose`, false, "Show matching DMN elements along with differences")
)

func init() {
	log.SetFlags(0)
	flag.Parse()
}

func main() {

	var err error
	set := make(map[string]bool)

	flag.Visit(func(f *flag.Flag) {
		set[f.Name] = true
	})

	switch {
	case !set[`url1`] || !set[`url2`]:
		err = fmt.Errorf(`both -url1 and -url2 required`)
	case !set[`success`] && !set[`warning`] && !set[`failure`]:
		err = fmt.Errorf(`at least one of -success, -warning, or -failure requried`)
	case !set[`warning`] && set[`details`]:
		*fWarning = true
	}

	if err != nil {
		log.Println(err)
		flag.Usage()
		os.Exit(2)
	}

	// Get the API for both environments.

	api1 := api.NewDmnApi(*fSvcUrl1)
	api2 := api.NewDmnApi(*fSvcUrl2)

	var dmnList *model.DmnList

	// Get the DmnList for the first environment and sort it.

	if dmnList, err = api1.DmnList(); err != nil {
		log.Fatal(err)
	}

	dmnList.Sort()

	var succ, warn, fail []string

	// Iterate through keys and versions of first environment.

	for _, di := range *dmnList {

		var dmn1, dmn2 *model.Dmn

		// Retrieve the DMN from the first environment.

		if dmn1, err = api1.DmnByKeyVer(di.Key, di.Version); err != nil {
			fail = append(fail, fmt.Sprintf(dmnFailFmt, `FAILURE`, `Env1`, di.Key, di.Version, err))
			continue
		} else if dmn2, err = api2.DmnByKeyVer(di.Key, di.Version); err != nil {
			fail = append(fail, fmt.Sprintf(dmnFailFmt, `FAILURE`, `Env2`, di.Key, di.Version, err))
			continue
		} else if reflect.DeepEqual(dmn1, dmn2) {
			succ = append(succ, fmt.Sprintf(dmnCompFmt, `SUCCESS`, di.Key, di.Version, `identical`))
			continue
		} else {
			warn = append(warn, fmt.Sprintf(dmnCompFmt, `WARNING`, di.Key, di.Version, `different`))
		}

		if !*fDetails {
			continue
		}

		// Process differences.

		if de, err := model.NewDmnElements(dmn1); err != nil {
			fail = append(fail, fmt.Sprintf(elmFailFmt, `FAILURE`, `Env1`, di.Key, di.Version, err))
		} else if err := de.Compare(dmn2); err != nil {
			fail = append(fail, fmt.Sprintf(elmFailFmt, `FAILURE`, `Env2`, di.Key, di.Version, err))
		} else {
			for key, val := range de {
				switch val {
				case 1:
					warn = append(warn, fmt.Sprintf("\tEnv1 <---      : %s", key))
				case -1:
					warn = append(warn, fmt.Sprintf("\t     ---> Env2 : %s", key))
				case 0:
					if !*fVerbose {break}
					warn = append(warn, fmt.Sprintf("\tEnv1 <--> Env2 : %s", key))
				}
			}
		}
	}

	if *fSuccess {
		for _, line := range succ {
			fmt.Println(line)
		}
	}
	if *fFailure {
		for _, line := range fail {
			fmt.Println(line)
		}
	}
	if *fWarning {
		for _, line := range warn {
			fmt.Println(line)
		}
	}
}
