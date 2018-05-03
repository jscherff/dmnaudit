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
	`bytes`
	`encoding/csv`
	`io`
)

// ------------------------------------------------------------------------
// DmnRules.
// ------------------------------------------------------------------------

type DmnRules interface {
	Bytes() ([]byte, error)
	String() (string)
	Write(io.Writer) (int, error)
	Headers() ([][]string)
	Rules() ([][]string)
}

type dmnRules struct {
	headers	[][]string
	rules	[][]string
}

func NewDmnRules(dmn *Dmn) (DmnRules, error) {

	// Create shorthand for nested objects to improve readability.

	inputs := dmn.Decision.DecisionTable.Inputs
	outputs := dmn.Decision.DecisionTable.Outputs
	rules := dmn.Decision.DecisionTable.Rules

	// Determine number of rows by counting rules. Determine number
	// of columns by counting input and output entries for the first
	// rule.

	rows := len(rules)
	cols := len(rules[0].InputEntries) + len(rules[0].OutputEntries)

	// Create the data structures: [row][col]string.

	table := &dmnRules{
		headers: make([][]string, 4),
		rules: make([][]string, rows),
	}

	for i := range table.headers {
		table.headers[i] = make([]string, cols)
	}

	for i := range table.rules {
		table.rules[i] = make([]string, cols)
	}

	// Populate the data structure.

	hcol := 0

	for _, input := range inputs {

		for _, inputExp := range input.InputExpressions {

			table.headers[0][hcol] = `Input`
			table.headers[1][hcol] = input.Label
			table.headers[2][hcol] = inputExp.Text
			table.headers[3][hcol] = inputExp.TypeRef

			hcol++
		}
	}

	for _, output := range outputs {

		table.headers[0][hcol] = `Output`
		table.headers[1][hcol] = output.Label
		table.headers[2][hcol] = output.Name
		table.headers[3][hcol] = output.TypeRef

		hcol++
	}


	for row, rule := range rules {

		ecol := 0

		for _, inputEntry := range rule.InputEntries {
			table.rules[row][ecol] = inputEntry.Text
			ecol++
		}

		for _, outputEntry := range rule.OutputEntries {
			table.rules[row][ecol] = outputEntry.Text
			ecol++
		}
	}

	table.headers[0] = append([]string{`Flow`}, table.headers[0]...)
	table.headers[1] = append([]string{`Label`}, table.headers[1]...)
	table.headers[2] = append([]string{`Name`}, table.headers[2]...)
	table.headers[3] = append([]string{`Type`}, table.headers[3]...)

	for i, rule := range table.rules {
		table.rules[i] = append([]string{`Rule`}, rule...)
	}

	return table, nil
}

func (this *dmnRules) Bytes() ([]byte, error) {

	buf := new(bytes.Buffer)
	cw := csv.NewWriter(buf)
	if err := cw.WriteAll(this.headers); err != nil {
		return nil, err
	}

	if err := cw.WriteAll(this.rules); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (this *dmnRules) String() (string) {

	if b, err := this.Bytes(); err != nil {
		return ``
	} else {
		return string(b)
	}
}

func (this *dmnRules) Write(w io.Writer) (int, error) {

	if b, err := this.Bytes(); err != nil {
		return 0, err
	} else {
		return w.Write(b)
	}
}

func (this *dmnRules) Headers() ([][]string) {
	return this.headers
}

func (this *dmnRules) Rules() ([][]string) {
	return this.rules
}
