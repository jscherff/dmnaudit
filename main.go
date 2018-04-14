package main

import (
	`fmt`
	`log`
	`github.com/jscherff/dmnsdk/api`
	//`github.com/jscherff/dmnsdk/model`
)

func main() {

	dmnApi := api.NewDmnApi(`http://esbeap-qa.24hourfit.com:8180`)

	/*
	if dl, err := dmnApi.GetDefinitionList(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%#v\n", dl)
	}
	*/

	if d, err := dmnApi.GetDefinitionByKey(`mi9-user-provisioning-rules-roles`); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%#v\n", d)
	}

	/*
	if d, err := model.NewDefinition(`tmp`); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%#v\n", d)
	}
	*/
}
