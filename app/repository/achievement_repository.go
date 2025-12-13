package repository

import (
	"context"
	"database/sql"
	"errors"
	"prisma/app/model"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository interface {
	Create(ctx context.Context, Achievement model.AchievementMongo) (*model.AchievementMongo, error)
	Update(ctx context.Context, Achievement model.AchievementMongo) (*model.AchievementMongo, error)
	Delete(ctx context.Context, Id string) error
	FindAll(ctx context.Context) ([]model.AchievementMongo, error)
	FindById(ctx context.Context, id string) (*model.AchievementMongo, error)
	FindByUser(ctx context.Context, UserId string) ([]model.AchievementMongo, error)
	FindByLecturer(ctx context.Context, UserId string) ([]model.AchievementMongo, error)
}

type AchievementRepositoryImpl struct {
	collection *mongo.Collection
	DB         *sql.DB
	Log        *logrus.Logger
}

func NewAchievementRepository(DB *mongo.Database, SQL *sql.DB, Log *logrus.Logger) AchievementRepository {
	return &AchievementRepositoryImpl{
		collection: DB.Collection("student_achievements"),
		Log:        Log,
		DB:         SQL,
	}
}

func (repo *AchievementRepositoryImpl) FindByLecturer(ctx context.Context, UserId string) ([]model.AchievementMongo, error) {
	//TODO implement me
	SQL := `SELECT a.mongo_achievement_id FROM lecturers as l 
       LEFT JOIN students s ON s.advisor_id = l.id
       JOIN achievement_references a ON a.student_id = s.id
       WHERE l.user_id = $1`
	rows, err := repo.DB.Query(SQL, UserId)
	if err != nil {
		return nil, err
	}
	panic("implement me")
}

func (repo *AchievementRepositoryImpl) Create(ctx context.Context, Achievement model.AchievementMongo) (*model.AchievementMongo, error) {
	ts := time.Now()
	Achievement.CreatedAt = ts
	result, err := repo.collection.InsertOne(ctx, Achievement)
	if err != nil {
		return nil, err
	}
	Achievement.ID = result.InsertedID.(primitive.ObjectID)
	SQL := "INSERT INTO achievement_references(student_id, mongo_achievement_id,status,created_at) VALUES ($1, $2, $3, $4)"
	res, err := repo.DB.ExecContext(ctx, SQL, Achievement.StudentID, Achievement.ID, "draft", ts)
	if err != nil {
		_, _ = repo.collection.DeleteOne(ctx, bson.M{
			"_id": Achievement.ID,
		})
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		_, _ = repo.collection.DeleteOne(ctx, bson.M{
			"_id": Achievement.ID,
		})
		return nil, err
	}
	if rows != 1 {
		_, _ = repo.collection.DeleteOne(ctx, bson.M{"_id": Achievement.ID})
		repo.Log.Fatalf("insert affected %d rows", rows)
	}
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

func (repo *AchievementRepositoryImpl) Delete(ctx context.Context, Id string) error {
	res, err := repo.DB.ExecContext(ctx, "UPDATE achievement_references SET status = 'DELETED' WHERE mongo_achievement_id = $1", Id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("no affected rows")
	}
	return nil
}

func (repo *AchievementRepositoryImpl) FindAll(ctx context.Context) ([]model.AchievementMongo, error) {
	res, err := repo.collection.Find(ctx, bson.D{})
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
	cursor, err := repo.collection.Find(ctx, bson.M{"_id": id})
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

func (repo *AchievementRepositoryImpl) FindByUser(ctx context.Context, UserId string) ([]model.AchievementMongo, error) {
	cursor, err := repo.collection.Find(ctx, bson.M{"studentId": UserId})
	if err != nil {
		return nil, err
	}
	achievements := []model.AchievementMongo{}

	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}
	return achievements, nil
}
