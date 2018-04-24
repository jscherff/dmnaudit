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
	`encoding/json`
	`encoding/xml`
	`fmt`
	`io`
	`net/http`
	`os`
	`strings`
)

const (
	jsonPrefix = ``
	jsonIndent = `    `
)

// load unmarshals JSON or XML from an io.Reader.
func load(dst interface{}, src interface{}, enc string) (err error) {

	switch obj := src.(type) {

	case io.Reader:

		switch strings.ToLower(enc) {

		case `json`:
			return json.NewDecoder(obj).Decode(&dst)

		case `xml`:
			return xml.NewDecoder(obj).Decode(&dst)

		default:
			return fmt.Errorf(`unsupported encoding: %s`, enc)
		}

	case string:

		buf := new(bytes.Buffer)

		if _, err = read(buf, obj); err != nil {
			return err
		} else {
			return load(dst, buf, enc)
		}

	default:
		return fmt.Errorf(`unsupported source: %T`, obj)
	}

	return nil
}

// read returns an io.Reader ready for unmarshalling.
func read(w io.Writer, s string) (int64, error) {

	switch true {

	case strings.HasPrefix(strings.ToLower(s), `http:`):
		return readUrl(w, s)

	case strings.HasPrefix(strings.ToLower(s), `https:`):
		return readUrl(w, s)

	case fileExists(s):
		return readFile(w, s)

	default:
		return readString(w, s)
	}
}

// readUrl returns a buffer filled from a URL.
func readUrl(w io.Writer, u string) (int64, error) {

	if resp, err := http.Get(u); err != nil {
		return 0, err
	} else {
		defer resp.Body.Close()
		return io.Copy(w, resp.Body)
	}
}

// readFile returns a byte buffer filled from a file.
func readFile(w io.Writer, f string) (int64, error) {

	if fh, err := os.Open(f); err != nil {
		return 0, err
	} else {
		defer fh.Close()
		return io.Copy(w, fh)
	}
}

// readString returns a byte buffer filled from a string.
func readString(w io.Writer, s string) (int64, error) {
	n, err := io.WriteString(w, s)
	return int64(n), err
}

// fileExists returns true if a file exists, false if it does not.
func fileExists(f string) (bool) {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return false
	} else if _, err := os.Stat(f); err == nil {
		return true
	}
	return false
}

// toJson marshals an object into a JSON byte array.
func toJson(t interface{}) ([]byte, error) {
	return json.MarshalIndent(t, jsonPrefix, jsonIndent)
}


// toMap converts a DMN into a hierarchy of map[string]interface{} and
// []interface{} objects.
func toMap(t interface{}) (map[string]interface{}, error) {

	b := new(bytes.Buffer)
	m := make(map[string]interface{})

	if err := json.NewEncoder(b).Encode(&t); err != nil {
		return nil, err
	} else if err := json.NewDecoder(b).Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}
