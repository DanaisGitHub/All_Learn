package response

import (
	"fmt"
	"http-server/internal/headers"
	"io"
)

type StatusCode int

const (
	STATUS_OK           StatusCode = 200
	STATUS_BAD_REQUEST  StatusCode = 400
	STATUS_SERVER_ERROR StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	rl := make([]byte, 0)
	switch statusCode {
	case STATUS_OK:
		rl = fmt.Appendf(nil, "HTTP/1.1 %d OK", STATUS_OK)
	case STATUS_BAD_REQUEST:
		rl = fmt.Appendf(nil, "HTTP/1.1 %d Bad Request", STATUS_BAD_REQUEST)
	case STATUS_SERVER_ERROR:
		rl = fmt.Appendf(nil, "HTTP/1.1 %d Internal Server Error", STATUS_SERVER_ERROR)
	default:
		rl = fmt.Appendf(nil, "HTTP/1.1 %d", statusCode)
	}
	fmt.Println(string(rl))
	w.Write(rl)
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", fmt.Sprint(contentLen))
	headers.Set("connection", "close")
	headers.Set("content-type", "text/plain")
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) {
	for key, val := range headers {
		header := fmt.Sprintf("%s: %s", key, val)
		w.Write(fmt.Appendf(nil, "%s", header))
		fmt.Println(header)
	}
}
