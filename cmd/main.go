package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RenterRus/dwld-downloader/internal/app"
)

func main() {
	path := flag.String("config", "../config.yaml", "path to config. Example: ../config.yaml")
	flag.Parse()
	if path == nil || len(*path) < 6 {
		log.Fatal("config flag not found")
		os.Exit(1)
	}

	time.Sleep(time.Second * 17)

	if err := app.NewApp(*path); err != nil {
		fmt.Println(err)
	}
}
