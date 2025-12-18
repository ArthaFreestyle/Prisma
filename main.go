package main

import (
	"fmt"
	"prisma/config"
	_ "prisma/docs" // Import generated swagger docs
)

// @title Prisma API
// @version 1.0
// @description API Documentation for Prisma Application
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@prisma.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	viperConfig := config.NewViper()
	log := config.NewLog(viperConfig)
	app := config.NewFiber(viperConfig)
	postgres := config.PostgresConnect(viperConfig, log)
	mongo := config.MongoConnect(viperConfig, log)
	redis := config.NewRedisClient(viperConfig, log)
	validate := config.NewValidator()

	config.Bootstrap(&config.BootstrapConfig{
		Postgres: postgres,
		App:      app,
		Log:      log,
		Config:   viperConfig,
		MongoDB:  mongo,
		Redis:    redis,
		Validate: validate,
	})

	// Swagger route

	port := viperConfig.GetInt("app.port")
	log.Infof("Server running on http://localhost:%d", port)
	log.Infof("Swagger UI available at http://localhost:%d/swagger/index.html", port)

	err := app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to start app: %v", err)
	}
}
