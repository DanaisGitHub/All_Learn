package request

import (
	"bytes"
	"fmt"
	"http-server/internal/headers"
	rl "http-server/internal/requestLine"
	"io"
	"strconv"
	"strings"
)

const NEWLINEs string = "\r\n"

var NEWLINEb = []byte(NEWLINEs)

// Request States
const (
	REQ_REQUESTLINE int = iota
	REQ_HEADER
	REQ_BODY
	REQ_DONE
)

// eg GET /coffee HTTP/1.1

type Request struct {
	RequestLine          *rl.RequestLine
	state                int
	Headers              headers.Headers // doesn't need to be a pointer
	Body                 []byte
	contentLength        int
	accumulator          []byte
	remainingAccumulator []byte
}

func (r *Request) String() string {
	rl := fmt.Sprintf("Request line:\n - Method: %s\n - Target: %s \n - Version: %s\n",
		r.RequestLine.Method,
		r.RequestLine.RequestTarget,
		r.RequestLine.HttpVersion)

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

func NewRequest() *Request {
	return &Request{
		RequestLine:   rl.NewRequestLine(),
		state:         REQ_REQUESTLINE,
		Headers:       headers.NewHeaders(),
		contentLength: 0,
		accumulator:   make([]byte, 0),
		Body:          make([]byte, 0),
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

// takes a given chunk, from running memorey and
func newLineParser(r *Request, chunk []byte, bytesRead int, maxRead int) (int, bool) {
	if bytesRead == 0 {
		return -1, false
	}

	r.accumulator = append(r.accumulator, cleanChunk(chunk[:min(maxRead, bytesRead)])...)
	idx := bytes.Index(r.accumulator, NEWLINEb)
	if idx == -1 { // No new line found
		return idx, false
	}
	// newline found withing chunk, need both before and after newline
	r.remainingAccumulator = append(r.remainingAccumulator, cleanChunk(r.accumulator[idx+len(NEWLINEb):])...) // may contain \x00... // maybe x,y,\x00,z // i don't care not
	r.accumulator = r.accumulator[:idx+len(NEWLINEb)]
	return idx, true
}

func moveState(r *Request, newState int) {
	r.state = newState
	r.accumulator = r.remainingAccumulator
	r.remainingAccumulator = make([]byte, 0)

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	const MAXREAD = 8
	r := NewRequest()
	chunk := make([]byte, MAXREAD) // at maximum 8 bytes will be read
	bytesRead := 0
	var isFullLine bool
	var tempContentLength int

	for { // assumed to always be reading, it should stop reading when stopped

		n, err := reader.Read(chunk)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read chunk:=> %w", err)
		}
		bytesRead += n
		_, isFullLine = newLineParser(r, chunk, n, MAXREAD)

		switch r.state {
		case REQ_REQUESTLINE:
			if !isFullLine {
				break
			}
			reqLineDone := false
			r.RequestLine, reqLineDone, err = rl.ParseRequestLine(r.accumulator, NEWLINEs)
			if err != nil {
				return nil, fmt.Errorf("couldn't create requestLine structure:\n %w", err)
			}
			if !reqLineDone {
				return nil, fmt.Errorf("request line malformed, first line couldn't be parsed as a line: %w", err)
			}
			moveState(r, REQ_HEADER)
			break

		case REQ_HEADER:
			if idx := bytes.Index(r.accumulator, NEWLINEb); !isFullLine && idx != 0 {
				break
			}
			_, done, err := r.parseHeaderLines(r.accumulator) // end of headers not being noticed
			if err != nil {
				return nil, fmt.Errorf("couldn't parse chunk: %w", err)
			}
			if cntLenStr, ok := r.Headers.Get("Content-Length"); ok {
				cntLen, err := strconv.ParseInt(cntLenStr, 10, 32)
				if err != nil {
					return nil, fmt.Errorf("content-length not properly formatted, it should be a number: %w", err)
				}
				r.contentLength = int(cntLen)
				tempContentLength = r.contentLength
			}
			if done {
				if r.contentLength > 0 {
					tempContentLength = r.contentLength
					moveState(r, REQ_BODY)
				} else { // when there's no body for some reason this isn't being called
					fmt.Println("no content length found moving to done")
					moveState(r, REQ_DONE)
					return r, nil
				}

			} else {
				r.accumulator = r.remainingAccumulator
				r.remainingAccumulator = make([]byte, 0)
			}
			break

		case REQ_BODY:
			// NEED finish body based of content length
			tempContentLength -= len(r.accumulator)
			r.Body = append(r.Body, r.accumulator...)

			if tempContentLength <= 0 || n <= 0 {

				fmt.Println("reached REQ_DONE")
				bodyLen := len(r.Body)
				moveState(r, REQ_DONE)
				if bodyLen > r.contentLength {
					return r, fmt.Errorf("content length is too small for actual content length\n acclaimedLength = %d, real length = %d \n%w", r.contentLength, bodyLen, err)
				} else if bodyLen < r.contentLength {
					return r, fmt.Errorf("content length is too large for actual content length\n acclaimedLength = %d, real length = %d \n%w", r.contentLength, bodyLen, err)
				}
				return r, nil
			}
			r.accumulator = make([]byte, 0)
			break

		default:
			return nil, fmt.Errorf("mismatched request reading states")
			break
		}
	}

}
