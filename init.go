package main

import (
	"go-api-samp/application"
	"go-api-samp/repository"
	repositoryif "go-api-samp/repository/interface"
	"go-api-samp/util/config"
	"go-api-samp/util/log"
	"os"
	"time"
)

func Init(provider repository.InitProvider) error {
	// todo docker-compose 応急処置
	env := os.Getenv("env")
	if env == "local" {
		time.Sleep(10* time.Second)
	}

	if err := config.LoadConfig(); err != nil {
		return err
	}

	if err := application.NewDBOpen(config.DB); err != nil {
		return err
	}

	if err := loadLocations(provider.GetInitManager()); err != nil {
		return err
	}

	log.NewLogger(config.Log)

	return nil
}

func loadLocations(manager repositoryif.InitManager) error {
	application.NewLocationsMap()

	locations, err := manager.FindAllLocations()
	if err != nil {
		return err
	}

	m := application.GetLocationsMap()

	for _, v := range locations {
		m[v.Id] = v.City
	}

	return nil
}