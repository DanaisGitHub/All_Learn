package request

import (
	"bytes"
	"fmt"
	"http-server/internal/headers"
	"io"
	"strconv"
	"strings"

	"golang.org/x/text/number"
)

const NEWLINEs string = "\r\n"

var NEWLINEb = []byte(NEWLINEs)

// Request States
const (
	REQ_INIT int = iota
	REQ_REQUESTLINE
	REQ_HEADER
	REQ_BODY
	REQ_DONE
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
	Body        []byte
}

func (r *Request) String() string {
	rl := fmt.Sprintf("Request line:\n - Method: %s\n - Target: %s \n - Version: %s\n", r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
	var h strings.Builder
	h.WriteString("Headers:\n")
	for key, value := range r.Headers {
		fmt.Fprintf(&h, "- %s: %s\n", key, value)
	}
	return rl + h.String()
}

// Parse and create the RequestLine which is just the first line of the whole request
func parseRequestLine(b []byte) (*RequestLine, bool, error) {
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

func NewRequest() *Request {
	return &Request{
		RequestLine: &RequestLine{},
		state:       REQ_INIT,
		Headers:     headers.NewHeaders(),
	}
}

// finds the position of sepator and return in i index including length of SEPARATOR
func (r *Request) parseHeaderLines(line []byte) (int, bool, error) {
	idx := bytes.Index(line, NEWLINEb)
	if idx == -1 {
		return -1, false, nil
	}
	n, requestDone, err := r.Headers.Parse(line) // what are we doing n
	if err != nil {
		return 0, false, fmt.Errorf("couldn't parse header: %w\nn = %d\ndone = %v", err, n, requestDone)
	}
	return n, requestDone, nil

}

func appendAccumulator(chunk []byte, accumulator []byte, bytesRead int, maxRead int) ([]byte, []byte, bool) {
	remainingChunk := make([]byte, 0)
	if bytesRead == 0 {
		return accumulator, remainingChunk, false
	}
	// add to accumulator
	idx := bytes.Index(chunk, NEWLINEb)
	if idx == -1 {
		// need to clean chunk here
		accumulator = append(accumulator, chunk[:min(maxRead, bytesRead)]...)
		return accumulator, remainingChunk, false
	}
	newIdk := idx + len(NEWLINEb)
	newChunk := chunk[:newIdk]
	accumulator = append(accumulator, newChunk...)
	remainingChunk = append(remainingChunk, chunk[newIdk:]...)
	return accumulator, remainingChunk, true
}

func moveState(newState int, accumulator []byte, maxRead int) (int, []byte) {
	return newState, make([]byte, 0)
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("-------------------------------------- PANIC --------------------------------")
			fmt.Printf("Recovery: %s", r)
		}
	}()
	const MAXREAD = 8

	r := new(Request)
	chunk := make([]byte, MAXREAD) // at maximum 8 bytes will be read
	bytesRead := 0
	r.state = REQ_REQUESTLINE
	r.Headers = headers.NewHeaders()
	accumulator := make([]byte, 0)
	remainingChunk := accumulator
	var isFullLine bool

	for {
		n, err := reader.Read(chunk)
		// ERROR: finished reading malformed request

		// if err == io.EOF && len(accumulator) == 0 && r.state != done { // last chunk + EOF
		// 	return nil, fmt.Errorf("EOF: %w", err) // if you get EOF before r.state == done, then there has to be an error
		// }
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read chunk:=> %w", err)
		}
		bytesRead += n

		accumulator, remainingChunk, isFullLine = appendAccumulator(chunk, accumulator, n, MAXREAD)

		switch r.state {

		case REQ_REQUESTLINE:
			reqLineDone := false
			// how will this act when no newline present
			r.RequestLine, reqLineDone, err = parseRequestLine(accumulator)
			if err != nil {
				return nil, fmt.Errorf("couldn't create requestLine structure:\n %w", err)
			}
			if reqLineDone {
				r.state, accumulator = moveState(REQ_HEADER, accumulator, MAXREAD)
			}

		case REQ_HEADER:

			_, done, err := r.parseHeaderLines(accumulator) // end of headers not being noticed
			if err != nil {
				return nil, fmt.Errorf("couldn't parse chunk: %w", err)
			}
			if done {
				r.state, accumulator = moveState(REQ_BODY, accumulator, MAXREAD)
			}
			if isFullLine {
				// clear acc
				accumulator = remainingChunk
			}

		case REQ_BODY:
			fmt.Println("You are parsing the body now")
			// based on content length we need to add that to accumalor then add to body bytes

			bodyLenStr, ok := r.Headers["content-length"]
			if !ok {
				fmt.Println("no content length found moving to done")
				moveState(REQ_DONE, accumulator, MAXREAD)
				break // break out of case
			}
			bodyLen,err :=  strconv.ParseInt(bodyLenStr,10,32) // numStr,base10,int32
			if err!=nil{
				
			}

		case REQ_DONE:
			fmt.Println(r)
			return r, nil

		default:
			return nil, fmt.Errorf("mismatched Request reading states")
		}
	}

}
