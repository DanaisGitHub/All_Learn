package requestline

import (
	"fmt"
	"strings"
)

type RequestLine struct {
	HttpVersion   string // 1.1 vs 2 vs ..
	RequestTarget string // /coffee
	Method        string // GET, POST ....
}

func NewRequestLine()*RequestLine{
	return &RequestLine{}
}

func ParseRequestLine(b []byte, NEWLINEs string) (*RequestLine, bool, error) {
	s := string(b)
	idx := strings.Index(s, NEWLINEs)
	if idx == -1 {
		return nil, false, nil
	}

	requestLine := s[:idx]
	// Seperate the Request line into 3 parts
	splitReqLine := strings.Split(requestLine, " ")
	if len(splitReqLine) != 3 {
		return nil, false, fmt.Errorf("malformed request line spaces")
	}
	httpVersion := strings.Split(splitReqLine[2], "/")
	if len(httpVersion) != 2 {
		return nil, false, fmt.Errorf("malformed http version")
	}
	if httpVersion[1] != "1.1" {
		return nil, false, fmt.Errorf("wrong http version")

	}
	rl := &RequestLine{
		HttpVersion:   httpVersion[1],
		Method:        splitReqLine[0],
		RequestTarget: splitReqLine[1],
	}
	return rl, true, nil
}
