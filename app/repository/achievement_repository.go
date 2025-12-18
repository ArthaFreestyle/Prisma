package repository

import (
	"context"
	"errors"
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

	updateData := bson.M{
		"student_id":      Achievement.StudentID,
		"achievment_type": Achievement.AchievementType,
		"title":           Achievement.Title,
		"description":     Achievement.Description,
		"details":         Achievement.Details,
		"attachments":     Achievement.Attachments,
		"tags":            Achievement.Tags,
		"updated_at":      time.Now(),
		"points":          Achievement.Points,
	}
	update := bson.M{
		"$set": updateData,
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
	oid, err := utils.ToObjectId(id) // Asumsi utils ini sudah benar handling errornya
	if err != nil {
		return nil, err
	}

	achievement := &model.AchievementMongo{}
	err = repo.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(achievement)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("achievement not found")
		}
		return nil, err
	}
	return achievement, nil
}
