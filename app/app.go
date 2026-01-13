package app

import (
	"log"
	"todogin/internal/api"
	"todogin/internal/config"
	"todogin/internal/database"
)

func Run() {
	conf, err := config.ConfigInit()
	log.Printf("(config.ConfigInit): Err: %v\n", err)

	db, err := database.DatabaseInit(conf)
	log.Printf("(database.DatabaseInit): Err: %v\n", err)

	api := api.ApiInit(db, conf)
	api.Run()
}
