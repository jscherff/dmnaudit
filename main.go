package main

import (
	//`fmt`
	`log`
	`reflect`
	`github.com/jscherff/dmnsdk/api`
	`github.com/jscherff/dmnsdk/model`
)

const (
	defRetrieve = `%s: could not get %s %s for key '%s' version '%d': %v`
	defCompare = `%s: %s and %s %s for Key '%s' version '%d': %s`
	esbeapPRD = `http://esbeap.24hourfit.com:8180`
	esbeapQA = `http://esbeap-qa.24hourfit.com:8180`
	esbeapDEV= `http://esbeap-dev.24hourfit.com:8180`
)

func init() {
	log.SetFlags(0)
}

func main() {

	// Get the API for both environments.
	dmnApi1 := api.NewDmnApi(esbeapPRD)
	dmnApi2 := api.NewDmnApi(esbeapQA)

	var (
		err error
		diMap1, diMap2 model.DefinitionInfoMap
	)

	// Get the [key][ver]->DMNinfo map for first environment.
	if diMap1, err = dmnApi1.GetDefinitionInfoMap(); err != nil {
		log.Fatal(err)
	}

	// Get the [key][ver]->DMNinfo map for second environment.
	if diMap2, err = dmnApi2.GetDefinitionInfoMap(); err != nil {
		log.Fatal(err)
	}


	// Iterate through keys and versions of first environment.
	for key, verMap := range diMap1 {
		for ver, di1  := range verMap {

			var (
				d1, d2 *model.Definition
				di2 *model.DefinitionInfo
			)

			// Retrieve the DMN for key/ver for the first environment.
			if d1, err = dmnApi1.GetDefinitionById(di1.Id); err != nil {
				log.Printf(defRetrieve, `FAILURE`, `PRD`, `DMN`, key, ver, err)
				continue
			}


			// Retrieve the DMN Info for key/ver for the second environment.
			if di2, err = diMap2.Get(key, ver); err != nil {
				log.Printf(defRetrieve, `FAILURE`, `QA`, `DMN Info`, key, ver, err)
				continue
			}

			// Retrieve the DMN for key/ver for the second environment.
			if d2, err = dmnApi2.GetDefinitionById(di2.Id); err != nil {
				log.Printf(defRetrieve, `FAILURE`, `QA`, `DMN`, key, ver, err)
				continue
			}

			// Deeply compare the two DMNs and show the results.
			if reflect.DeepEqual(d1, d2) {
				log.Printf(defCompare, `SUCCESS`, `PRD`, `QA`, `DMN`, key, ver, `identical`)
			} else {
				log.Printf(defCompare, `WARNING`, `PRD`, `QA`, `DMN`, key, ver, `different`)
			}
		}
	}
}
