package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

var NEWLINE = []byte("\r\n")

type Headers map[string]string

// returns the value of specific header field
//   - if none will return empty string
func (h Headers) Get(field string) (string,bool) {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	value, ok := h[strings.ToLower(field)]
	if !ok {
		return "",false
	}
	return value,true
}

// will create a new header with a value or will append the key with new value and comma
func (h Headers) Set(field, value string) {
	_, ok := h[field]
	if !ok {
		h[field] = value
		return
	}
	h[field] += ", " + value
}

func isToken(fieldName string) bool {
	var tokenRegex = regexp.MustCompile(`^[!#$%&'*+\-.^_` + "`" + `|~0-9A-Za-z]+$`)
	return tokenRegex.MatchString(fieldName)
}

// this needs to be able to parse multiple lines
//   - n = number of bytes parsed
//   - done = finishing parsing headers
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// need to get a full line of the data

	n = 0
	done = false
	err = nil
	for {
		// idx = index of new line
		// data[n:] = everything left to be parsed
		idx := bytes.Index(data[n:], NEWLINE)
		if idx == -1 { // no new lines to be parsed
			break
		}
		if idx == 0 { // assumed the end of the header
			done = true
			n += len(NEWLINE) // that new line at 0 counts as parts of the header
			break
		}

		// currData = what is currently being parsed
		currData := data[n : n+idx]
		// n += --everything-- NEWLINE + NEWLINE
		n += idx + len(NEWLINE)

		field, value, err := h.checkValidKeyValuePair(currData)
		if err != nil {
			return 0, false, err
		}

		if !isToken(string(field)) {
			return 0, false, fmt.Errorf("field name %s is not valid token", field)
		}
		h.Set(string(field), string(value))
	}
	return n, done, err
}

// Checks if data passed in can be converted into a key-value pair
func (h Headers) checkValidKeyValuePair(data []byte) ([]byte, []byte, error) {
	field := []byte("")
	value := []byte("")
	parts := bytes.SplitN(data, []byte(":"), 2)

	if len(parts) != 2 {
		return field, value, fmt.Errorf("malformed header")
	}
	field = parts[0]
	if bytes.HasSuffix(field, []byte(" ")) {
		return field, value, fmt.Errorf("field of header malformed, space(s) between field & :, this is not allowed")
	}
	field = bytes.ToLower(bytes.TrimSpace(field))
	value = bytes.TrimSpace(parts[1])

	return field, value, nil
}

func NewHeaders() Headers {
	// since headers is of type map and maps are already pointers, no pointers made
	return Headers{}
}
