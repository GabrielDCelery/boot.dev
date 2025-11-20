package server

import (
	"fmt"
	"net"
)

type Server struct {
	listener net.Listener
	errChan  chan error
	quitChan chan struct{}
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", port))

	if err != nil {
		return &Server{}, fmt.Errorf("failed to create listener on port %d, reason: %v\n", port, err)
	}

	server := &Server{
		listener: listener,
		errChan:  make(chan error, 1),
		quitChan: make(chan struct{}),
	}
	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	close(s.quitChan)
	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("failed to close listener: %v", err)
	}
	return nil
}

func (s *Server) Err() <-chan error {
	return s.errChan
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			select {
			case <-s.quitChan:
				// shutdown requested, exit gracefully
				return
			case s.errChan <- err:
				continue
			default:
				fmt.Printf("dropping error, %v", err)
				continue
			}
		}

		go func(conn net.Conn) {
			defer conn.Close()
			responseMsg := "HTTP/1.1 200 OK\r\n" + "Content-Type: text/plain\r\n" + "Content-Length: 14\r\n" + "\r\n" + "Hello World!\r\n"
			conn.Write([]byte(responseMsg))
		}(conn)
	}
}
