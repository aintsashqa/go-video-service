package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/aintsashqa/go-video-service/internal/app"
	"github.com/aintsashqa/go-video-service/internal/module/video"
)

func main() {
	config := flag.String("config", "default-config", "Choose configuration filename")
	filename := flag.String("filename", "", "Choose file to upload")
	flag.Parse()

	*filename = strings.TrimSpace(*filename)
	if len(*filename) == 0 {
		log.Fatal("Option filename cannot be empty value")
	}

	f, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("Cannot open file: %s", err)
	}
	defer f.Close()

	fInfo, err := f.Stat()
	if err != nil {
		log.Fatalf("Cannot get file info: %s", err)
	}

	app := app.Initialize(*config)

	log.Print("Creating file")

	deps := video.CreateCommandDeps{
		Context:    context.Background(),
		Repository: app.Repository.VideoRepository,
	}

	args := video.CreateCommandArgs{
		OriginalName: fInfo.Name(),
		Size:         fInfo.Size(),
		Reader:       f,
	}

	resultId, err := video.CreateCommand(deps, args)
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}

	log.Printf("File created with name %s", resultId)
}
