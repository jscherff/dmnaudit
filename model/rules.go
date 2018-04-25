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
)

// ------------------------------------------------------------------------
// DmnRules.
// ------------------------------------------------------------------------

type DmnRules interface {
	NoOp()
	ToCsv() ([]byte, error)
	String() (string)
}

type dmnRules struct {
	Headers	[][]string
	Rules	[][]string
}

func (this *dmnRules) NoOp() {

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
		Headers: make([][]string, 4),
		Rules: make([][]string, rows),
	}

	for i := range table.Headers {
		table.Headers[i] = make([]string, cols)
	}

	for i := range table.Rules {
		table.Rules[i] = make([]string, cols)
	}

	// Populate the data structure.

	hcol := 0

	for _, input := range inputs {

		for _, inputExp := range input.InputExpressions {

			table.Headers[0][hcol] = `Input`
			table.Headers[1][hcol] = input.Label
			table.Headers[2][hcol] = inputExp.Text
			table.Headers[3][hcol] = inputExp.TypeRef

			hcol++
		}
	}

	for _, output := range outputs {

		table.Headers[0][hcol] = `Output`
		table.Headers[1][hcol] = output.Label
		table.Headers[2][hcol] = output.Name
		table.Headers[3][hcol] = output.TypeRef

		hcol++
	}

	ecol := 0

	for i, rule := range rules {

		for _, inputEntry := range rule.InputEntries {
			table.Rules[i][ecol] = inputEntry.Text
			ecol++
		}

		for _, outputEntry := range rule.OutputEntries {
			table.Rules[i][ecol] = outputEntry.Text
			ecol++
		}
	}

	table.Headers[0] = append([]string{`Entry Type`}, table.Headers[0]...)
	table.Headers[1] = append([]string{`Entry Label`}, table.Headers[1]...)
	table.Headers[2] = append([]string{`Entry Name`}, table.Headers[2]...)
	table.Headers[3] = append([]string{`Data Type`}, table.Headers[3]...)

	for i, rule := range table.Rules {
		table.Rules[i] = append([]string{`Rule`}, rule...)
	}

	return table, nil
}

func (this *dmnRules) ToCsv() ([]byte, error) {

	buf := new(bytes.Buffer)
	cw := csv.NewWriter(buf)
	if err := cw.WriteAll(this.Headers); err != nil {
		return nil, err
	}

	if err := cw.WriteAll(this.Rules); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (this *dmnRules) String() (string) {

	if b, err := this.ToCsv(); err != nil {
		return ``
	} else {
		return string(b)
	}
}


