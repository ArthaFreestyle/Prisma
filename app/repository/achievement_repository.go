package repository

import (
	"context"
	"database/sql"
	"prisma/app/model"
	"prisma/utils"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository interface {
	Create(ctx context.Context, Achievement model.AchievementMongo) (*model.AchievementMongo, error)
	Update(ctx context.Context, Achievement model.AchievementMongo) (*model.AchievementMongo, error)
	FindAll(ctx context.Context, Id []string) ([]model.AchievementMongo, error)
	FindById(ctx context.Context, id string) (*model.AchievementMongo, error)
}

type AchievementRepositoryImpl struct {
	collection *mongo.Collection
	DB         *sql.DB
	Log        *logrus.Logger
}

func NewAchievementRepository(DB *mongo.Database, Log *logrus.Logger) AchievementRepository {
	return &AchievementRepositoryImpl{
		collection: DB.Collection("student_achievements"),
		Log:        Log,
	}
}

func (repo *AchievementRepositoryImpl) Create(ctx context.Context, Achievement model.AchievementMongo) (*model.AchievementMongo, error) {
	Achievement.CreatedAt = time.Now()
	result, err := repo.collection.InsertOne(ctx, Achievement)
	if err != nil {
		return nil, err
	}
	Achievement.ID = result.InsertedID.(primitive.ObjectID)
	return &Achievement, nil
}

func (repo *AchievementRepositoryImpl) Update(ctx context.Context, Achievement model.AchievementMongo) (*model.AchievementMongo, error) {
	Achievement.UpdatedAt = time.Now()
	update := bson.M{
		"$set": Achievement,
	}
	_, err := repo.collection.UpdateOne(ctx, bson.M{"_id": Achievement.ID}, update)
	if err != nil {
		return nil, err
	}
	return &Achievement, nil
}

func (repo *AchievementRepositoryImpl) FindAll(ctx context.Context, Id []string) ([]model.AchievementMongo, error) {
	ObjectID, err := utils.ToObjectsId(Id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": bson.M{
			"$in": ObjectID,
		},
	}
	res, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer res.Close(ctx)
	achievements := []model.AchievementMongo{}
	for res.Next(ctx) {
		var achievement model.AchievementMongo
		err := res.Decode(&achievement)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, achievement)
	}
	return achievements, nil
}

func (repo *AchievementRepositoryImpl) FindById(ctx context.Context, id string) (*model.AchievementMongo, error) {
	oid, err := utils.ToObjectId(id)
	cursor, err := repo.collection.Find(ctx, bson.M{"_id": oid})
	if err != nil {
		return nil, err
	}
	achievement := &model.AchievementMongo{}
	err = cursor.Decode(achievement)
	if err != nil {
		return nil, err
	}
	return achievement, nil
}
