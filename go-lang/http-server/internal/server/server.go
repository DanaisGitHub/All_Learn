package server

import (
	"fmt"
	"http-server/internal/request"
	"http-server/internal/response"
	"net"
	"strconv"
	"sync/atomic"
)

const (
	STATE_INIT int = iota
	STATE_LISTENING
	STATE_CLOSED
)

type Server struct {
	state     int
	Listener  net.Listener
	closeBool atomic.Bool // because atomic no need for locking
	port      int
}

type IServer interface {
	Close() error
	listen() error
	handle(conn net.Conn)
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, fmt.Errorf("failure to create tcp server on port %d due to error: %w", port, err)
	}
	s := &Server{
		state:     STATE_INIT,
		Listener:  listener,
		closeBool: atomic.Bool{},
		port:      port,
	}
	go func() {
		err := s.listen()
		if err != nil {
			panic(err)
		}
	}()
	return s, nil
}

func (s *Server) Close() error {
	s.state = STATE_CLOSED
	s.closeBool.Store(true)
	return s.Listener.Close()
}

func (s *Server) listen() error {
	handle := func(conn net.Conn) {
		err := s.handle(conn)
		if err != nil {
			fmt.Printf("couldn't handle request because: %v", err.Error())
		}
	}

	fmt.Printf("server open and listening on port %d\n", s.port)
	s.state = STATE_LISTENING
	for !s.closeBool.Load() { // whilst false keep going
		conn, err := s.Listener.Accept()
		if err != nil {
			return err
		}
		go handle(conn)
	}
	return nil
}

func (s *Server) handle(conn net.Conn) error {
	// Read
	r, err := request.RequestFromReader(conn)
	if err != nil {
		return fmt.Errorf("error handling singular conn: %v\nerror:=%w", conn, err)
	}
	fmt.Println(r)

	// Time
	statusCode64, err := strconv.ParseInt("200", 10, 64)
	if err != nil {
		return fmt.Errorf("when handling the connection couldn't read the method: %w", err)
	}
	err = response.WriteStatusLine(conn, response.StatusCode(statusCode64))
	if err != nil {
		return fmt.Errorf("couldn't write response headline: %w", err)
	}
	response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	return nil
}
