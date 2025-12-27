package request

import (
	"bytes"
	"fmt"
	"http-server/internal/headers"
	"io"
	"strconv"
	"strings"
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
	RequestLine   *RequestLine
	state         int
	Headers       headers.Headers // doesn't need to be a pointer
	Body          []byte
	contentLength int
	accumulator   []byte
}

func (r *Request) String() string {
	rl := fmt.Sprintf("Request line:\n - Method: %s\n - Target: %s \n - Version: %s\n", r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
	var h strings.Builder
	h.WriteString("Headers:\n")
	for key, value := range r.Headers {
		fmt.Fprintf(&h, "- %s: %s\n", key, value)
	}
	b := "Body:\n"
	b += string(r.Body)
	return rl + h.String() + b
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
		RequestLine:   &RequestLine{},
		state:         REQ_INIT,
		Headers:       headers.NewHeaders(),
		contentLength: -1,
		accumulator:   make([]byte, 0),
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

func cleanChunk(c []byte) []byte {
	return bytes.TrimRight(c, "\x00")
}

func appendAccumulator(chunk []byte, accumulator []byte, bytesRead int, maxRead int) ([]byte, []byte, int, bool) {
	remainingChunk := make([]byte, 0)
	if bytesRead == 0 {
		return accumulator, remainingChunk, -1, false
	}
	accumulator = append(accumulator, cleanChunk(chunk[:min(maxRead, bytesRead)])...)
	idx := bytes.Index(accumulator, NEWLINEb)
	if idx == -1 { // No new line found
		return accumulator, remainingChunk, idx, false
	}
	// newline found withing chunk, need both before and after newline
	newIdk := idx + len(NEWLINEb)
	remainingChunk = append(remainingChunk, cleanChunk(accumulator[newIdk:])...) // may contain \x00... // maybe x,y,\x00,z // i don't care not
	accumulator = accumulator[:newIdk]

	return accumulator, remainingChunk, idx, true
}

func moveState(r *Request, newState int, accumulator []byte, maxRead int) []byte {
	r.state = newState
	return make([]byte, 0)
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	const MAXREAD = 8

	r := new(Request)
	chunk := make([]byte, MAXREAD) // at maximum 8 bytes will be read
	bytesRead := 0
	r.state = REQ_REQUESTLINE
	r.Headers = headers.NewHeaders()
	accumulator := make([]byte, 0)
	remainingChunk := accumulator
	var isFullLine bool
	//idx := -1

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

		accumulator, remainingChunk, _, isFullLine = appendAccumulator(chunk, accumulator, n, MAXREAD)

		switch r.state {

		case REQ_REQUESTLINE:
			if !isFullLine {
				break
			}
			reqLineDone := false
			r.RequestLine, reqLineDone, err = parseRequestLine(accumulator)
			if err != nil {
				return nil, fmt.Errorf("couldn't create requestLine structure:\n %w", err)
			}
			if !reqLineDone {
				return nil, fmt.Errorf("request line malformed, first line couldn't be parsed as a line: %w", err)
			}
			accumulator = moveState(r, REQ_HEADER, accumulator, MAXREAD)
			accumulator = remainingChunk

		case REQ_HEADER:
			if idxAcc := bytes.Index(accumulator, NEWLINEb); idxAcc == 0 {
				accumulator = moveState(r, REQ_BODY, accumulator, MAXREAD)
				accumulator = remainingChunk
				break
			}
			if !isFullLine {
				break
			}
			_, done, err := r.parseHeaderLines(accumulator) // end of headers not being noticed
			if err != nil {
				return nil, fmt.Errorf("couldn't parse chunk: %w", err)
			}
			if cntLenStr, ok := r.Headers["content-length"]; ok {
				cntLen, err := strconv.ParseInt(cntLenStr, 10, 32)
				if err != nil {
					return nil, fmt.Errorf("content-length not properly formatted, it should be a number: %w", err)
				}
				r.contentLength = int(cntLen)
			}
			if done && r.contentLength != -1 {
				accumulator = moveState(r, REQ_BODY, accumulator, MAXREAD)
			}
			if done && r.contentLength == -1 {
				fmt.Println("no content length found moving to done")
				accumulator = moveState(r, REQ_DONE, accumulator, MAXREAD)
			}

			accumulator = remainingChunk

		case REQ_BODY:
			if err == io.EOF {
				cntLen := r.contentLength
				accumulator = append(accumulator, cleanChunk(remainingChunk)...)
				r.Body = append(r.Body, accumulator...)
				bodyLen := len(r.Body)
				if bodyLen > cntLen {
					return r, fmt.Errorf("content length is too small for actual content length\n acclaimedLength = %d, real length = %d \n%w", cntLen, bodyLen, err)
				} else if cntLen > bodyLen {
					return r, fmt.Errorf("content length is too large for actual content length\n acclaimedLength = %d, real length = %d \n%w", cntLen, bodyLen, err)
				}
				accumulator = moveState(r, REQ_DONE, accumulator, MAXREAD)
			}

		case REQ_DONE:
			fmt.Println(r)
			return r, nil

		default:
			return nil, fmt.Errorf("mismatched Request reading states")
		}
	}

}
