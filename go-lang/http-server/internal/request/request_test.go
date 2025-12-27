package request

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

func TestBodyParsingo(t *testing.T) {
	str := "POST /coffee HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/8.4.0\r\n" +
		"Accept: */*\r\n" +
		"Content-Type: application/json\r\n" +
		"Content-Length: 39\r\n" +
		"\r\n" +
		`{"type": "dark mode", "size": "medium"}`
	// Test: Standard Body
	reader := &chunkReader{
		data:            str,
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, `{"type": "dark mode", "size": "medium"}`, string(r.Body))
}

func TestBodyParsingTo(t *testing.T) {
	str := "POST /submit HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"hello world!\n"
	// Test: Standard Body
	reader := &chunkReader{
		data:            str,
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "hello world!\n", string(r.Body))

	// Test: Body shorter than reported content length
	reader = &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 20\r\n" +
			"\r\n" +
			"partial content",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.Error(t, err)
}

func TestPrintingOfRequest(t *testing.T) {
	// Test: Good GET Request line
	str := "GET / HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"
	r, _ := RequestFromReader(strings.NewReader(str))
	fmt.Print(r)
}

func TestRequestLineParse0(t *testing.T) {
	// Test: Good GET Request line
	str := "GET / HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"
	r, err := RequestFromReader(strings.NewReader(str))
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)

}

func TestRequestLineParse1(t *testing.T) {
	str := "GET /coffee HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"

	// Test: Good GET Request line with path
	r, err := RequestFromReader(strings.NewReader(str))
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestRequestLineParse2(t *testing.T) {

	str := "/coffee HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"
	//Test: Invalid number of parts in request line
	_, err := RequestFromReader(strings.NewReader(str))
	require.Error(t, err)
}

func TestRequestLineParse3(t *testing.T) {
	str := "GET / HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"

	// Test: Good GET Request line
	reader := &chunkReader{
		data:            str,
		numBytesPerRead: 50,
	}

	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)
}
func TestRequestLineParsex(t *testing.T) {
	str := "GET / HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"

	// Test: Good GET Request line
	reader := &chunkReader{
		data:            str,
		numBytesPerRead: 50,
	}

	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)
}
func TestRequestLineParse4(t *testing.T) {

	str := "GET /coffee HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"
	// Test: Good GET Request line with path
	reader := &chunkReader{
		data:            str,
		numBytesPerRead: 1,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)

}

func TestHeaders(t *testing.T) {
	str := "GET / HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"
	// Test: Standard Headers
	reader := &chunkReader{
		data:            str,
		numBytesPerRead: 50,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "localhost:42069", r.Headers.Get("host"))
	assert.Equal(t, "curl/7.81.0", r.Headers.Get("user-agent"))
	assert.Equal(t, "*/*", r.Headers["accept"])

	// Test: Malformed Header
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.Error(t, err)
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}
