package repository

import (
	"context"
	"prisma/app/model"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AnalyticsRepository interface {
	Statistics(ctx context.Context) ([]model.Statistics, error)
	Reporting(ctx context.Context, id string) ([]*model.Statistics, error)
}

type AnalyticsRepositoryImpl struct {
	DB  *mongo.Collection
	Log *logrus.Logger
}

func NewAnalyticsRepository(Log *logrus.Logger, db *mongo.Database) *AnalyticsRepositoryImpl {
	return &AnalyticsRepositoryImpl{
		DB:  db.Collection("student_achievements"),
		Log: Log,
	}
}

func (repo *AnalyticsRepositoryImpl) Statistics(ctx context.Context) ([]model.Statistics, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "$year", Value: "$createdAt"}}},
			{Key: "international", Value: bson.D{
				{Key: "$sum", Value: bson.D{
					{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$eq", Value: bson.A{"$details.competitionLevel", "international"}}},
						1,
						0,
					}},
				}},
			}},
			{Key: "national", Value: bson.D{
				{Key: "$sum", Value: bson.D{
					{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$eq", Value: bson.A{"$details.competitionLevel", "national"}}},
						1,
						0,
					}},
				}},
			}},
			{Key: "regional", Value: bson.D{
				{Key: "$sum", Value: bson.D{
					{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$eq", Value: bson.A{"$details.competitionLevel", "regional"}}},
						1,
						0,
					}},
				}},
			}},
			{Key: "local", Value: bson.D{
				{Key: "$sum", Value: bson.D{
					{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$eq", Value: bson.A{"$details.competitionLevel", "local"}}},
						1,
						0,
					}},
				}},
			}},
		}}},

		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "tahun", Value: bson.D{{Key: "$toString", Value: "$_id"}}}, // Convert Year Int ke String
			{Key: "data", Value: bson.D{
				{Key: "international", Value: "$international"},
				{Key: "national", Value: "$national"},
				{Key: "regional", Value: "$regional"},
				{Key: "local", Value: "$local"},
			}},
		}}},

		{{Key: "$sort", Value: bson.D{{Key: "tahun", Value: -1}}}},
	}

	cursor, err := repo.DB.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var stats []model.Statistics
	if err = cursor.All(ctx, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

func (repo *AnalyticsRepositoryImpl) Reporting(ctx context.Context, id string) ([]*model.Statistics, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "studentId", Value: id},
		}}},

		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "$year", Value: "$createdAt"}}},
			{Key: "international", Value: bson.D{
				{Key: "$sum", Value: bson.D{
					{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$eq", Value: bson.A{"$details.competitionLevel", "international"}}},
						1, 0,
					}},
				}},
			}},
			{Key: "national", Value: bson.D{
				{Key: "$sum", Value: bson.D{
					{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$eq", Value: bson.A{"$details.competitionLevel", "national"}}},
						1, 0,
					}},
				}},
			}},
			{Key: "regional", Value: bson.D{
				{Key: "$sum", Value: bson.D{
					{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$eq", Value: bson.A{"$details.competitionLevel", "regional"}}},
						1, 0,
					}},
				}},
			}},
			{Key: "local", Value: bson.D{
				{Key: "$sum", Value: bson.D{
					{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$eq", Value: bson.A{"$details.competitionLevel", "local"}}},
						1, 0,
					}},
				}},
			}},
		}}},

		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "tahun", Value: bson.D{{Key: "$toString", Value: "$_id"}}},
			{Key: "data", Value: bson.D{
				{Key: "international", Value: "$international"},
				{Key: "national", Value: "$national"},
				{Key: "regional", Value: "$regional"},
				{Key: "local", Value: "$local"},
			}},
		}}},

		{{Key: "$sort", Value: bson.D{{Key: "tahun", Value: -1}}}},
	}

	cursor, err := repo.DB.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var stats []*model.Statistics
	if err := cursor.All(ctx, &stats); err != nil {
		return nil, err
	}

	if stats == nil {
		stats = []*model.Statistics{}
	}

	return stats, nil
}
