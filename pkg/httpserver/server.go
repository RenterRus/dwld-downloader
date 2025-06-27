package httpserver

import (
	"fmt"
	"net/http"
)

type Server struct {
	Host   string
	Port   int
	Enable bool
	Mux    *http.ServeMux
}

func NewHttpServer(s *Server) error {
	if s.Enable {
		fmt.Printf("\nHTTP starting in %s:%d (/health)\n", s.Host, s.Port)

		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.Host, s.Port), s.Mux); err != nil {
			return fmt.Errorf("http server is stopped: %w", err)
		}
	}

	return nil
}
