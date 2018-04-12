package main

import (
	`fmt`
	`log`
	`github.com/jscherff/dmnaudit/model`
)

const (
	//dmnService =	`http://esbeap-qa.24hourfit.com:8180`
	dmnService =	`http://esb-qa-13.24hourfit.com:8080`
	dmnListPath =	`/engine-rest/decision-definition`
	dmnXmlPath =	`/engine-rest/decision-definition/key/%s/xml`
)

func init() {
	log.SetFlags(0)
}

func main() {

	dmns := new(model.Dmns)

	url := dmnService + dmnListPath

	if err := dmns.ReadUrl(url); err != nil {
		log.Fatal(err)
	}

	for _, dmn := range *dmns {

		fmt.Printf("DMN Id: %s\n", dmn.Id)
		fmt.Printf("DMN Key: %s\n", dmn.Key)
		fmt.Printf("DMN Category: %s\n", dmn.Category)
		fmt.Printf("DMN Name: %s\n", dmn.Name)
		fmt.Printf("DMN Version: %d\n", dmn.Version)

		url := fmt.Sprintf(dmnService + dmnXmlPath, dmn.Key)

		if err := dmn.ReadUrl(url); err != nil {
			log.Fatal(err)
		}

		if err := dmn.Definitions.ReadString(dmn.Xml); err != nil {
			log.Fatal(err)
		}

		dmn.Definitions = new(model.Definitions)
		def := dmn.Definitions
		def.ReadString(dmn.Xml)

		fmt.Printf("DMN Def XMLName: %#v\n", def.XMLName)
		fmt.Printf("DMN Def ID: %s\n", def.Id)
		fmt.Printf("DMN Def Xmlns: %s\n", def.Xmlns)
		fmt.Printf("DMN Def Name: %s\n", def.Name)
		fmt.Printf("DMN Def Namespace: %s\n", def.Namespace)

		dec := def.Decision

		fmt.Printf("Decision XMLName: %#v\n", dec.XMLName)
		fmt.Printf("Decision ID: %s\n", dec.Id)
		fmt.Printf("Decision Name: %s\n", dec.Name)

		decTab := dec.DecisionTable

		fmt.Printf("Decision Table XMLName: %#v\n", decTab.XMLName)
		fmt.Printf("Decision Table ID: %s\n", decTab.Id)
		fmt.Printf("Decision Table HitPolicy: %s\n", decTab.HitPolicy)

		for _, input := range decTab.Inputs {

			fmt.Printf("\tInput ID: %s\n", input.Id)
			fmt.Printf("\tInput Label: %s\n", input.Label)

			for _, inpExp := range input.InputExpressions {

				fmt.Printf("\t\tInput Expression ID: %s\n", inpExp.Id)
				fmt.Printf("\t\tInput Expression TypeRef: %s\n", inpExp.TypeRef)
				fmt.Printf("\t\tInput Expression Text: %s\n", inpExp.Text)
			}
		}

		for _, output := range decTab.Outputs {

			fmt.Printf("\tOutput ID: %s\n", output.Id)
			fmt.Printf("\tOutput Label: %s\n", output.Label)
			fmt.Printf("\tOutput Name: %s\n", output.Name)
			fmt.Printf("\tOutput TypeRef: %s\n", output.TypeRef)
		}

		for _, rule := range decTab.Rules {

			fmt.Printf("\tRule ID: %s\n", rule.Id)

			for _, inpEnt := range rule.InputEntries {

				fmt.Printf("\t\tInput Entry ID: %s\n", inpEnt.Id)
				fmt.Printf("\t\tInput Entry Text: %s\n", inpEnt.Text)
			}

			for _, outEnt := range rule.OutputEntries {

				fmt.Printf("\t\tOutput Entry ID: %s\n", outEnt.Id)
				fmt.Printf("\t\tOutput Entry Text: %s\n", outEnt.Text)
			}
		}
	}
}
