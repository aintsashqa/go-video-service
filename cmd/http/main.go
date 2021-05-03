package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/aintsashqa/go-video-service/internal/app"
	"github.com/aintsashqa/go-video-service/internal/delivery/http"
)

func main() {
	ctx := context.Background()
	config := flag.String("config", "default-config", "Choose configuration filename")
	flag.Parse()

	app := app.Initialize(*config)

	log.Print("Initialize handler")
	handler := http.NewHandler(app.Repository)

	log.Print("Initialize server")
	server := http.NewServer(app, handler)

	log.Printf("Starting server on %s:%d", app.Config.HttpConfig.Host, app.Config.HttpConfig.Port)
	go func() {
		if err := server.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Print("Initialize graceful shutdown")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Kill, os.Interrupt)
	<-signalChan

	log.Print("Starting graceful shutdown")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Print("Graceful shutdown was successfully done")
}
