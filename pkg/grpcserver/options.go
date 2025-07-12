package grpcserver

import (
	"net"
)

type Option func(*Server)

func Port(host, port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort(host, port)
	}
}
