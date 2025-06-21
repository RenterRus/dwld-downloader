package http

import (
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	Host   string
	Port   int
	Enable bool
}

func NewHttpServer(s *Server) error {
	if s.Enable {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("Available at %s", time.Now().Format(time.DateTime))))
		})

		fmt.Printf("\nHTTP starting in %s:%d (/health)\n", s.Host, s.Port)

		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.Host, s.Port), nil); err != nil {
			return fmt.Errorf("http server is stopped: %w", err)
		}
	}

	return fmt.Errorf("http server not started: disabled in config")
}
