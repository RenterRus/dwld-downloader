package http

import (
	"fmt"
	"net/http"
	"time"
)

func NewRoute() *http.ServeMux {
	srv := http.NewServeMux()
	srv.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Available at %s", time.Now().Format(time.DateTime))))
	})

	return srv
}
