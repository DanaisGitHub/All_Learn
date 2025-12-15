package request

import (
	"bytes"
	"fmt"
	"http-server/internal/headers"
	"io"
	"strings"
)

const NEWLINEs string = "\r\n"

var NEWLINEb = []byte(NEWLINEs)

// Request States
const (
	initialized int = iota
	requestLine
	header
	body
	done
)

// eg GET /coffee HTTP/1.1
type RequestLine struct {
	HttpVersion   string // 1.1 vs 2 vs ..
	RequestTarget string // /coffee
	Method        string // GET, POST ....
}

type Request struct {
	RequestLine *RequestLine
	state       int
	Headers     headers.Headers // doesn't need to be a pointer
}

// Parse and create the RequestLine which is just the first line of the whole request
func parseRequestLine(s string) (*RequestLine, error) {
	idx := strings.Index(s, NEWLINEs)
	if idx == -1 {
		return nil, fmt.Errorf("no SEPARATOR = %v found", NEWLINEs)
	}

	requestLine := s[:idx]
	// Seperate the Request line into 3 parts
	splitReqLine := strings.Split(requestLine, " ")
	if len(splitReqLine) != 3 {
		return nil, fmt.Errorf("malformed request line spaces")
	}
	httpVersion := strings.Split(splitReqLine[2], "/")
	if len(httpVersion) != 2 {
		return nil, fmt.Errorf("malformed http version")
	}
	if httpVersion[1] != "1.1" {
		return nil, fmt.Errorf("wrong http version")

	}
	rl := &RequestLine{
		HttpVersion:   httpVersion[1],
		Method:        splitReqLine[0],
		RequestTarget: splitReqLine[1],
	}
	return rl, nil
}

func NewRequest() *Request {
	return &Request{
		RequestLine: &RequestLine{},
		state:       initialized,
		Headers:     headers.NewHeaders(),
	}
}

// finds the position of sepator and return in i index including length of SEPARATOR
func (r *Request) parseHeaderLines(line []byte) (int, bool, error) {
	idx := bytes.Index(line, NEWLINEb)
	if idx == -1 {
		return 0, false, nil
	}

	n, requestDone, err := r.Headers.Parse(line)
	if err != nil {
		return -1, false, fmt.Errorf("couldn't parse header: %w\nn=%d\ndone=%v", err, n, requestDone)
	}
	return idx + len(NEWLINEs), requestDone, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := new(Request)
	chunk := make([]byte, 8) // at maximum 8 bytes will be read
	bytesRead := 0
	r.state = requestLine

	r.Headers = headers.NewHeaders()
	tempStr := ""

	for {
		n, err := reader.Read(chunk)
		if err == io.EOF {
			return nil, fmt.Errorf("EOF: %w", err)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read chunk:=> %w", err)
		}
		bytesRead += n

		switch r.state {
		case requestLine:
			i := bytes.Index(chunk, []byte(NEWLINEs))
			if i == -1 { //not found
				tempStr += string(chunk)
				break
			}
			tempStr += string(chunk)[:i+len([]byte(NEWLINEs))]
			r.RequestLine, err = parseRequestLine(tempStr)
			if err != nil {
				return nil, fmt.Errorf("couldn't create RequestLine structure %w", err)
			}
			tempStr = ""
			r.state = header

		case header:
			i := bytes.Index(chunk, NEWLINEb)
			if i == -1 { //not found
				tempStr += string(chunk)
				break
			}
			tempStr += string(chunk[:i+len(NEWLINEb)])
			_, done, err := r.parseHeaderLines([]byte(tempStr)) // ERROR: Cannot parse header malformed
			if err != nil {
				return nil, fmt.Errorf("couldn't parse chunk: %w", err)
			}
			if done {
				r.state = body
			}

		case body:
			fmt.Println("You are parsing the body now")
			r.state = done

		case done:
			return r, nil
		default:
			return nil, fmt.Errorf("mismatched Request reading states")
		}

	}

}
