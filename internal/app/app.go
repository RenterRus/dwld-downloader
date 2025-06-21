package app

import (
	"dwld-downloader/internal/controller/http"
	"fmt"
	"sync"
)

func NewApp() error {
	conf, err := ReadConfig(".", "config_example")
	if err != nil {
		return fmt.Errorf("ReadConfig: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		http.NewHttpServer(&http.Server{
			Host:   conf.HTTP.Host,
			Port:   conf.HTTP.Port,
			Enable: conf.HTTP.Enable,
		})
	}()

	wg.Wait()

	return nil
}
