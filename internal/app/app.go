package app

import (
	"log"

	"github.com/aintsashqa/go-video-service/internal/config"
	"github.com/aintsashqa/go-video-service/internal/repository"
)

type App struct {
	Config     *config.Config
	Repository *repository.Container
}

func Initialize(configFilename string) *App {
	if len(configFilename) == 0 {
		configFilename = "default-config"
	}

	log.Print("Initialize config")
	conf, err := config.Init(configFilename)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Initialize repository container")
	repository, err := repository.NewContainer(conf)
	if err != nil {
		log.Fatal(err)
	}

	app := App{
		Config:     &conf,
		Repository: repository,
	}

	log.Print("Application successfully initialized")
	return &app
}
