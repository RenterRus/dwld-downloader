// Package grpcserver implements HTTP server.
package grpcserver

import (
	"fmt"
	"net"

	pbgrpc "google.golang.org/grpc"
)

const (
	_defaultAddr = ":80"
)

// Server -.
type Server struct {
	App     *pbgrpc.Server
	notify  chan error
	address string
}

// New -.
func New(opts ...Option) *Server {
	s := &Server{
		App:     pbgrpc.NewServer(),
		notify:  make(chan error, 1),
		address: _defaultAddr,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start() {
	go func() {
		ln, err := net.Listen("tcp", s.address)
		if err != nil {
			s.notify <- fmt.Errorf("failed to listen: %w", err)
			close(s.notify)

			return
		}
		fmt.Println("grpc listen in: ", s.address)
		s.notify <- s.App.Serve(ln)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	s.App.GracefulStop()

	return nil
}
