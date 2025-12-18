package config

import (
	"database/sql"
	"prisma/app/repository"
	"prisma/app/service"
	"prisma/middleware"
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
	StudentRepository := repository.NewStudentRepositoryImpl(config.Log, config.Postgres)
	LecturerRepository := repository.NewLecturerRepositoryImpl(config.Log, config.Postgres)
	AnalyticsRepository := repository.NewAnalyticsRepository(config.Log, config.MongoDB)
	LogoutRepository := repository.NewLogoutRepository(config.Redis, config.Log)
	AchievementRepository := repository.NewAchievementRepository(config.MongoDB, config.Log)
	AchievementRepositoryReference := repository.NewAchievementReferenceRepository(config.Log, config.Postgres)

	secret := []byte(config.Config.GetString("app.jwt-secret"))
	//Setup Service
	AchievementService := service.NewAchievementService(AchievementRepository, StudentRepository, AchievementRepositoryReference, config.Validate, config.Log)
	AuthService := service.NewAuthService(UserRepository, LogoutRepository, config.Log, secret)
	UserService := service.NewUserService(UserRepository, StudentRepository, LecturerRepository, config.Postgres, config.Validate, config.Log)
	StudentService := service.NewStudentService(StudentRepository, AchievementRepositoryReference)
	LecturerService := service.NewLecturerService(LecturerRepository, StudentRepository)
	AnalyticsService := service.NewAnalyticsService(AnalyticsRepository)

	RouteConfig := routes.RouteConfig{
		App:                config.App,
		UserService:        UserService,
		AuthService:        AuthService,
		AchievementService: AchievementService,
		LecturerService:    LecturerService,
		AnalyticsService:   AnalyticsService,
		StudentService:     StudentService,
		AuthMiddleware:     middleware.AuthRequired(secret),
	}

	RouteConfig.Setup()

}
