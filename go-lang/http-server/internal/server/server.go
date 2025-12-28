package server

import (
	"fmt"
	"net"
)

const (
	STATE_INIT int = iota
	STATE_LISTENING
	STATE_CLOSED
)

type Server struct {
	state    int
	Listener net.Listener
}

type IServer interface {
	Close() error
	listen() error
	handle(conn net.Conn)
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprint(port))
	if err != nil {
		return nil, fmt.Errorf("failure to create tcp server on port %d due to error: %w", port, err)
	}

	s := &Server{
		state:    STATE_INIT,
		Listener: listener,
	}

	s.listen()

	return s, nil
}

func (s *Server) Close() error {

	s.state = STATE_CLOSED
	return s.Listener.Close()
}

func (s *Server) listen() error {
	s.Listener.Accept()
	return nil
}

func (s *Server) handle(conn net.Conn) error {
	return nil
}
