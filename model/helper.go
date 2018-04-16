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
	`bufio`
	`bytes`
	`encoding/xml`
	`encoding/json`
	`fmt`
	`io`
	`net/http`
	`os`
	`reflect`
	`strings`
)

const (
	jsonPrefix = ``
	jsonIndent = `    `
)

// load unmarshals JSON or XML from an io.Reader.
func load(dst interface{}, src interface{}, enc string) (error) {

	switch t := src.(type) {

	case io.Reader:

		switch strings.ToLower(enc) {

		case `json`:
			return json.NewDecoder(t).Decode(&dst)

		case `xml`:
			return xml.NewDecoder(t).Decode(&dst)

		default:
			return fmt.Errorf(`unsupported encoding: %s`, enc)
		}

	case string:

		var (
			b bytes.Buffer
			w = bufio.NewWriter(&b)
		)

		if _, err := read(w, t); err != nil {
			return err
		} else {
			w.Flush()
			return load(dst, bufio.NewReader(&b), enc)
		}

	default:
		return fmt.Errorf(`unsupported source: %T`, t)
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

// json marshals an object into a JSON byte array.
func toJson(t interface{}) ([]byte, error) {
	return json.MarshalIndent(t, jsonPrefix, jsonIndent)
}

// toMap creates a map from a struct, including only fields that match a tag.
func toMap(t interface{}, m map[string]interface{}, tid string) (error) {

	v := reflect.ValueOf(t).Elem()

	if v.Type().Kind() != reflect.Struct {
		return fmt.Errorf(`kind %q is not struct`, v.Type().Kind().String())
	}

	for i := 0; i < v.NumField(); i++ {

		f := v.Field(i)
		t := v.Type().Field(i)

		if !f.IsValid() || !f.CanAddr() || !f.CanInterface() {
			continue
		}

		if tag, ok := t.Tag.Lookup(tid); !ok {
			continue
		} else {

			tval := strings.Split(tag, `,`)[0]

			switch tval {
			case `-`:
				continue
			case ``:
				m[t.Name] = f.Interface()
			default:
				m[tval] = f.Interface()
			}
		}
	}

	return nil
}
