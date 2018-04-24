package main

import (
	`fmt`
	`log`
	`reflect`
	`github.com/jscherff/dmnsdk/api`
	`github.com/jscherff/dmnsdk/model`
)

const (
	dmnRetrieve = `[%s] could not get %s %s for key '%s' version '%d': %v`
	defCompare = `[%s] %s and %s %s for Key '%s' version '%d': %s`
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
		dmnList *model.DmnList
	)

	// Get the [key][ver]->DMNinfo map for first environment.
	if dmnList, err = dmnApi1.DmnList(); err != nil {
		log.Fatal(err)
	}

	// Iterate through keys and versions of first environment.
	for _, di := range *dmnList {

		var dmn1, dmn2 *model.Dmn

		// Retrieve the DMN for key/ver for the first environment.
		if dmn1, err = dmnApi1.DmnByKeyVer(di.Key, di.Version); err != nil {
			log.Printf(dmnRetrieve, `FAILURE`, `PRD`, `DMN`, di.Key, di.Version, err)
		}

		// Retrieve the DMN for key/ver for the first environment.
		if dmn2, err = dmnApi2.DmnByKeyVer(di.Key, di.Version); err != nil {
			log.Printf(dmnRetrieve, `FAILURE`, `QA`, `DMN`, di.Key, di.Version, err)
		}

		// Deeply compare the two DMNs and show the results.
		if reflect.DeepEqual(dmn1, dmn2) {
			log.Printf(defCompare, `SUCCESS`, `PRD`, `QA`, `DMN`, di.Key, di.Version, `identical`)
		} else {
			log.Printf(defCompare, `WARNING`, `PRD`, `QA`, `DMN`, di.Key, di.Version, `different`)
			showDiff(dmn1, dmn2)
		}
	}
}

func showDiff(d1, d2 *model.Dmn) {

	if de, err := model.NewDmnElements(d1); err != nil {
		log.Println(err)
	} else if err := de.Compare(d2); err != nil {
		log.Println(err)
	} else {
		for _, key := range de.SortedKeys() {
			if de[key] == 1 {
				fmt.Printf("\t<---\t%s\n", key)
			} else if de[key] == -1 {
				fmt.Printf("\t--->\t%s\n", key)
			} else {
				fmt.Printf("\t####\t%s\n", key)
			}
		}
	}
}
