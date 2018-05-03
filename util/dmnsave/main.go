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
	`os`
	`github.com/jscherff/dmnsdk/api`
)

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	flag.Parse()
}

func main() {

	api := api.NewDmnApi(`http://esbeap.24hourfit.com:8180`)

	dmnList, err := api.DmnList()

	if err != nil {
		log.Fatal(err)
	}

	dmnList.Sort()

	for _, di := range *dmnList {

		id, key, ver := di.Id, di.Key, di.Version

		if xml, err := api.DmnXmlByKeyVer(key, ver); err != nil {
			log.Println(err)
		} else if b, err := xml.Xml(); err != nil {
			log.Println(err)
		} else {
			save(fmt.Sprintf(`dmnxml_key_%s_ver_%d.xml`, key, ver), b)
			save(fmt.Sprintf(`dmnxml_id_%s.xml`, id), b)
		}

		if dmn, err := api.DmnByKeyVer(key, ver); err != nil {
			log.Println(err)
		} else if b, err := dmn.Json(); err != nil {
			log.Println(err)
		} else {
			save(fmt.Sprintf(`dmn_key_%s_ver_%d.json`, key, ver), b)
			save(fmt.Sprintf(`dmn_id_%s.json`, id), b)
		}
	}
}

func save(f string, b []byte) {
	if fh, err := os.Create(f); err != nil {
		log.Println(err)
	} else {
		defer fh.Close()
		fh.Write(b)
	}
}
