package config

import (
	"database/sql"
	"prisma/app/repository"
	"prisma/app/service"
	"prisma/routes"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type BootstrapConfig struct {
	App      *fiber.App
	Postgres *sql.DB
	MongoDB  *mongo.Database
	Redis    *redis.Client
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {

	//Setup Repository
	UserRepository := repository.NewUserRepository(config.Postgres, config.Log)

	secret := []byte(config.Config.GetString("app.jwt-secret"))
	//Setup Service
	UserService := service.NewUserService(UserRepository, config.Log, secret)

	RouteConfig := routes.RouteConfig{
		App:            config.App,
		UserService:    UserService,
		AuthMiddleware: nil,
	}

	RouteConfig.Setup()

}
