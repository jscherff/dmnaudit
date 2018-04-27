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
	`io`
	`os`
	`reflect`
	`github.com/jscherff/dmnsdk/api`
)

var success, warning, failure Report

type Message string

func (this Message) Fmt(args ...interface{}) string {
	return fmt.Sprintf(string(this), args...)
}

type Report []string

func (this *Report) Add(args ...string) {
	*this = append(*this, args...)
}

func (this *Report) Print(w io.Writer) {
	for _, line := range *this {
		fmt.Fprintln(w, line)
	}
}

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	flag.Parse()
}

func main() {

	var (
		out io.WriteCloser
		err error
	)

	set := make(map[string]bool)

	flag.Visit(func(f *flag.Flag) {
		set[f.Name] = true
	})

	if !set[`url1`] || !set[`url2`] {
		log.Println(`both -url1 and -url2 required`)
		flag.Usage()
		os.Exit(2)
	}

	if !*fSuccess && !*fWarning && !*fFailure {
		*fWarning = true
		*fFailure = true
	}

	if !*fWarning && *fDetails {
		*fWarning = true
	}

	if !set[`file`] {
		out = os.Stdout
	} else if out, err = os.Create(*fOutFile); err != nil {
		log.Fatal(err)
	} else {
		defer out.Close()
	}

	// Get the API for both environments.

	api1 := api.NewDmnApi(*fSvcUrl1)
	api2 := api.NewDmnApi(*fSvcUrl2)

	// Get the DmnList for the first environment and sort it.

	dmnList, err := api1.DmnList()

	if err != nil {
		log.Fatal(err)
	}

	dmnList.Sort()

	// Iterate through keys and versions of first environment.

	for _, di := range *dmnList {

		if dmn1, err := api1.DmnByKeyVer(di.Key, di.Version); err != nil {
			failure.Add(dmnFailure.Fmt(`Env1`, di.Key, di.Version, err))
		} else if dmn2, err := api2.DmnByKeyVer(di.Key, di.Version); err != nil {
			failure.Add(dmnFailure.Fmt(`Env2`, di.Key, di.Version, err))
		} else if reflect.DeepEqual(dmn1, dmn2) {
			success.Add(cmpSuccess.Fmt(di.Key, di.Version))
		} else {
			warning.Add(cmpWarning.Fmt(di.Key, di.Version))
			if *fDetails {
				diff(di, dmn1, dmn2)
			}
		}
	}

	if *fSuccess {
		success.Print(out)
	}
	if *fFailure {
		failure.Print(out)
	}
	if *fWarning {
		warning.Print(out)
	}
}
