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
	dmnFailFmt = `could not get %s DMN for key '%s' v%d: %v`
	dmnCompFmt = `Env1 and Env2 DMN Key %s version %d: %s`
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

	report := make(map[string][]string)

	// Iterate through keys and versions of first environment.

	for _, di := range *dmnList {

		var dmn1, dmn2 *model.Dmn

		// Retrieve the DMN from the first environment.

		if dmn1, err = api1.DmnByKeyVer(di.Key, di.Version); err != nil {
			msg := fmt.Sprintf(dmnFailFmt, `Env1`, di.Key, di.Version, err)
			report[`FAILURE`] = append(report[`FAILURE`], msg)
		}

		// Retrieve the DMN from the second environment.

		if dmn2, err = api2.DmnByKeyVer(di.Key, di.Version); err != nil {
			msg := fmt.Sprintf(dmnFailFmt, `Env2`, di.Key, di.Version, err)
			report[`FAILURE`] = append(report[`FAILURE`], msg)
		}

		// Deeply compare the two DMNs and show the results.

		if reflect.DeepEqual(dmn1, dmn2) {
			msg := fmt.Sprintf(dmnCompFmt, di.Key, di.Version, `identical`)
			report[`SUCCESS`] = append(report[`SUCCESS`], msg)
			println(msg)
			continue
		} else {
			msg := fmt.Sprintf(dmnCompFmt, di.Key, di.Version, `different`)
			report[`WARNING`] = append(report[`WARNING`], msg)
			println(msg)
		}

		// Process differences.

		if de, err := model.NewDmnElements(dmn1); err != nil {
			log.Println(err)
		} else if err := de.Compare(dmn2); err != nil {
			log.Println(err)
		} else {

			var msg string

			for _, key := range de.SortedKeys() {
				if de[key] == 1 {
					msg = fmt.Sprintf("%s:\t%s", `Env1`, key)
				} else if de[key] == -1 {
					msg = fmt.Sprintf("%s:\t%s", `Env2`, key)
				} else if *fVerbose {
					msg = fmt.Sprintf("%s:\t%s", `Both`, key)
				} else {
					continue
				}
			}

			report[`DETAILS`] = append(report[`DETAILS`], msg)
		}
	}

	for category, lines := range report {
		for _, line := range lines {
			fmt.Printf("[%s] %s\n", category, line)
		}
	}
}



/*
	if set[`file`] {
		if out, err = os.Create(*fCsvFile); err != nil {
			log.Fatal(err)
		}
		defer out.Close()
	} else {
		out = os.Stdout
	}

	if _, err := rules.Write(out); err != nil {
		log.Fatal(err)
	}
*/
