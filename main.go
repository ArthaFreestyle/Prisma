package main

import (
	"prisma/config"

	"github.com/spf13/viper"
)

func main() {
	viperConfig := viper.New()
	log := config.Log()
	Mongo := config.MongoConnect(viperConfig, log)
	Postgres := config.PostgresConnect(viperConfig, log)
}
