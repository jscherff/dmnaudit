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
	`github.com/jscherff/dmnsdk/api`
	`github.com/jscherff/dmnsdk/model`
)

const (
	esbeapPRD = `http://esbeap.24hourfit.com:8180`
	esbeapQA = `http://esbeap-qa.24hourfit.com:8180`
	esbeapDEV= `http://esbeap-dev.24hourfit.com:8180`
)

var (
	fSvcUrl = flag.String(`url`, ``, "SvcUrl `<url>` and port")
	fDmnId = flag.String(`id`, ``, "Retrieve DMN with ID `<id>`")
	fDmnKey = flag.String(`key`, ``, "Retrieve DMN with key `<key>`")
	fDmnVer = flag.Int(`ver`, 0, "Retrieve DMN version `<ver>` (requires -key)")
	fCsvFile = flag.String(`file`, ``, "Store results in file `<file>`")
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
	case !set[`url`]:
		err = fmt.Errorf(`-service flag is required`)
	case !set[`key`] && !set[`id`]:
		err = fmt.Errorf(`-id or -key must be set`)
	case !set[`key`] && set[`ver`]:
		err = fmt.Errorf(`-ver requires -key`)
	}

	if err != nil {
		log.Printf("%v\n\n", err)
		flag.Usage()
		os.Exit(2)
	}

	var (
		dmn *model.Dmn
		rules model.DmnRules
		out io.WriteCloser
	)

	api := api.NewDmnApi(*fSvcUrl)

	switch {
	case set[`ver`]:
		dmn, err = api.DmnByKeyVer(*fDmnKey, *fDmnVer)
	case set[`id`]:
		dmn, err = api.DmnById(*fDmnId)
	case set [`key`]:
		dmn, err = api.DmnByKey(*fDmnKey)
	}

	if err != nil {
		log.Fatal(err)
	} else if rules, err = dmn.Rules(); err != nil {
		log.Fatal(err)
	}

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
}
