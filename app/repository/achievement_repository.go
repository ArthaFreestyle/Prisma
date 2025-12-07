package repository

import (
	"context"
	"prisma/app/model"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository interface {
	Create(ctx context.Context, Achievement model.AchievementMongo) (*model.AchievementMongo, error)
	Update(ctx context.Context, Achievement model.AchievementMongo) (*model.AchievementMongo, error)
	Delete(ctx context.Context, Achievement model.AchievementMongo) error
	FindAll(ctx context.Context) ([]model.AchievementMongo, error)
	FindById(ctx context.Context, id string) (model.AchievementMongo, error)
	FindByUser(ctx context.Context, UserId string) ([]model.AchievementMongo, error)
}

type AchievementRepositoryImpl struct {
	collection *mongo.Collection
	Log        *logrus.Logger
}

func NewAchievementRepository(DB *mongo.Database, Log *logrus.Logger) AchievementRepository {
	return &AchievementRepositoryImpl{
		collection: DB.Collection("student_achievements"),
		Log:        Log,
	}
}

func (repo AchievementRepositoryImpl) Create(ctx context.Context, Achievement model.AchievementMongo) (*model.AchievementMongo, error) {
	Achievement.CreatedAt = time.Now()
	result, err := repo.collection.InsertOne(ctx, Achievement)
	if err != nil {
		return nil, err
	}
	Achievement.ID = result.InsertedID.(primitive.ObjectID)
	return &Achievement, nil
}

func (repo AchievementRepositoryImpl) Update(ctx context.Context, Achievement model.AchievementMongo) (*model.AchievementMongo, error) {

	panic("implement me")
}

func (repo AchievementRepositoryImpl) Delete(ctx context.Context, Achievement model.AchievementMongo) error {
	//TODO implement me
	panic("implement me")
}

func (repo AchievementRepositoryImpl) FindAll(ctx context.Context) ([]model.AchievementMongo, error) {
	//TODO implement me
	panic("implement me")
}

func (repo AchievementRepositoryImpl) FindById(ctx context.Context, id string) (model.AchievementMongo, error) {
	//TODO implement me
	panic("implement me")
}

func (repo AchievementRepositoryImpl) FindByUser(ctx context.Context, UserId string) ([]model.AchievementMongo, error) {
	//TODO implement me
	panic("implement me")
}
