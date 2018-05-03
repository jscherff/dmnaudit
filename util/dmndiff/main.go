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
	`reflect`
	`strconv`
	`github.com/jscherff/dmnsdk/api`
	`github.com/jscherff/dmnsdk/model`
)


const (
	cmpSuccess = `DMNs are Identical`
	cmpWarning = `DMNs are Different`
	dmnFailure = `Could not get DMN`
	elmWarning = `Elements are Different`
	elmSuccess = `Elements are Identical`
	elmFailure = `Could not process DMN`
)

type Results interface {
	Success(*model.DmnInfo, ...string)
	Warning(*model.DmnInfo, ...string)
	Failure(*model.DmnInfo, ...string)
	Print(io.Writer) error
}

func NewResults() Results {

	this := new(results)

	this.rows = append(this.rows, []string{
		`Result`,
		`DMN Key`,
		`DMN Name`,
		`DMN ID`,
		`DMN Version`,
		`Reason`,
		`Service`,
		`Details`,
	})

	return this
}

type results struct {
	rows [][]string
}

func (this *results) add(result string, di *model.DmnInfo, cols ...string) {
	row := []string{result, di.Key, di.Name, di.Id, strconv.Itoa(di.Version)}
	row = append(row, cols...)
	this.rows = append(this.rows, row)
}

func (this *results) Success(di *model.DmnInfo, cols ...string) {
	this.add(`SUCCESS`, di, cols...)
}

func (this *results) Warning(di *model.DmnInfo, cols ...string) {
	this.add(`WARNING`, di, cols...)
}

func (this *results) Failure(di *model.DmnInfo, cols ...string) {
	this.add(`FAILURE`, di, cols...)
}

func (this *results) Print(w io.Writer) error {

	if err := csv.NewWriter(w).WriteAll(this.rows); err != nil {
		return err
	}

	return nil
}

var report = NewResults()

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

	if !*fSuccess && !*fWarning && !*fFailure && !*fDetails {
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
			if *fFailure {
				report.Failure(di, dmnFailure, *fSvcUrl1, err.Error())
			}
		} else if dmn2, err := api2.DmnByKeyVer(di.Key, di.Version); err != nil {
			if *fFailure {
				report.Failure(di, dmnFailure, *fSvcUrl2, err.Error())
			}
		} else if reflect.DeepEqual(dmn1, dmn2) {
			if *fSuccess {
				report.Success(di, cmpSuccess)
			}
		} else {
			if *fWarning {
				report.Warning(di, cmpWarning)
			}
			if *fDetails {
				diff(di, dmn1, dmn2)
			}
		}
	}

	if err := report.Print(out); err != nil {
		log.Fatal(err)
	}
}

func diff(di *model.DmnInfo, dmn1, dmn2 *model.Dmn) {

	if de, err := model.NewDmnElements(dmn1); err != nil {
		if *fFailure {
			report.Failure(di, elmFailure, *fSvcUrl1, err.Error())
		}
	} else if err := de.Compare(dmn2); err != nil {
		if *fFailure {
			report.Failure(di, elmFailure, *fSvcUrl2, err.Error())
		}
	} else {
		for _, key := range de.SortedKeys() {
			switch de[key] {
			case 1:
				report.Warning(di, elmWarning, *fSvcUrl1, key.String())
			case -1:
				report.Warning(di, elmWarning, *fSvcUrl2, key.String())
			case 0:
				if *fVerbose {
					report.Success(di, elmSuccess, `Both Services`, key.String())
				}
			}
		}
	}
}
