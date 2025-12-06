package headers

import (
	"bytes"
	"fmt"
	"maps"
	"regexp"
	"strings"
)

var NEWLINE = []byte("\r\n")

type Headers map[string]string

func (h Headers) Get(field string) string {
	return h[strings.ToLower(field)]
}

func (h Headers) Set(field, value string) {
	field = strings.ToLower(field)
	keys := maps.Keys(h)

	found := false
	for key := range keys {
		if key == field {
			found = true
			h[field] = h[field] + fmt.Sprintf(", %s", value)
			break
		}
	}
	if !found {
		h[field] = value
	}
}

func isToken(fieldName string) bool {
	var tokenRegex = regexp.MustCompile(`^[!#$%&'*+\-.^_` + "`" + `|~0-9A-Za-z]+$`)
	return tokenRegex.MatchString(fieldName)
}

// this needs to be able to parse multiple lines
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// need to get a full line of the data

	n = 0
	done = false
	err = nil
	for {
		idx := bytes.Index(data[n:], NEWLINE)
		if idx == -1 {
			break
		}
		if idx == 0 { // no header
			done = true
			n += len(NEWLINE)
			break
		}
		currData := data[n : n+idx]
		n += idx + len(NEWLINE)
		parts := bytes.SplitN(currData, []byte(":"), 2)
		if len(parts) != 2 {
			return 0, false, fmt.Errorf("malformed header")
		}
		field := parts[0]
		if bytes.HasSuffix(field, []byte(" ")) {
			return 0, false, fmt.Errorf("field of header malformed, space(s) between field & :, this is not allowed")
		}
		field = bytes.TrimSpace(field)
		value := bytes.TrimSpace(parts[1])
		if !isToken(string(field)) {
			return 0, false, fmt.Errorf("field name %s is not valid token", field)
		}
		h.Set(string(field), string(value))
	}
	return n, done, err
}

func NewHeaders() Headers {
	return Headers{}
}
