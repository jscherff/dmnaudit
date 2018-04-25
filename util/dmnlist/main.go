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
	`encoding/csv`
	`flag`
	`log`
	`io`
	`os`
	`strconv`
	`github.com/jscherff/dmnsdk/api`
	`github.com/jscherff/dmnsdk/model`
)

const (
	esbeapPRD = `http://esbeap.24hourfit.com:8180`
	esbeapQA = `http://esbeap-qa.24hourfit.com:8180`
	esbeapDEV= `http://esbeap-dev.24hourfit.com:8180`
)

var (
	fSvcUrl = flag.String(`url`, ``, "Use service at `http[s]://<hostname>[:<port>]`")
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

	if !set[`url`] {
		log.Println(`-service flag is required`)
		flag.Usage()
		os.Exit(2)
	}

	var (
		dmns *model.DmnList
		rows [][]string
		out io.WriteCloser
	)

	api := api.NewDmnApi(*fSvcUrl)

	if dmns, err = api.DmnList(); err != nil {
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

	dmns.Sort()
	rows = append(rows, []string{`name`, `key`, `version`, `id`})

	for _, dmn := range *dmns {
		row := []string{dmn.Name, dmn.Key, strconv.Itoa(dmn.Version), dmn.Id}
		rows = append(rows, row)
	}

	if err := csv.NewWriter(out).WriteAll(rows); err != nil {
		log.Fatal(err)
	}
}
