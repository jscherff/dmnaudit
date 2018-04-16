package main

import (
	//`fmt`
	`log`
	`reflect`
	`github.com/jscherff/dmnsdk/api`
	`github.com/jscherff/dmnsdk/model`
)

func init() {
	log.SetFlags(0)
}

func main() {

	dmnApi1 := api.NewDmnApi(`http://esbeap.24hourfit.com:8180`)
	dmnApi2 := api.NewDmnApi(`http://esbeap-qa.24hourfit.com:8180`)

	var (
		diMap1, diMap2 model.DefinitionInfoMap
		err error
	)

	if diMap1, err = dmnApi1.GetDefinitionInfoMap(); err != nil {
		log.Fatal(err)
	}

	if diMap2, err = dmnApi2.GetDefinitionInfoMap(); err != nil {
		log.Fatal(err)
	}


	for key, verMap := range diMap1 {
		for ver, di1  := range verMap {

			var (
				d1, d2 *model.Definition
				di2 *model.DefinitionInfo
			)

			if d1, err = dmnApi1.GetDefinitionById(di1.Id); err != nil {
				log.Printf(`WARNING: could not get Prd Definition for key '%s' version '%s': %v`, key, ver, err)
				continue
			}


			if di2, err = diMap2.Get(key, ver); err != nil {
				log.Printf(`WARNING: could not get QA Definition Info: %v`, err)
				continue
			}

			if d2, err = dmnApi2.GetDefinitionById(di2.Id); err != nil {
				log.Printf(`WARNING: could not get QA Definition for key '%s' version '%s': %v`, key, ver, err)
				continue
			}

			if reflect.DeepEqual(d1, d2) {
				log.Printf(`Prd and QA Definition for Key '%s' version '%d': identical`, key, ver)
			} else {
				log.Printf(`WARNING: Prd and QA Definition for Key '%s' version '%d': different`, key, ver)
			}
		}
	}
}
