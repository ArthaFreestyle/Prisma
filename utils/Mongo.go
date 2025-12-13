package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToObjectsId(Ids []string) ([]primitive.ObjectID, error) {
	ObjectIds := make([]primitive.ObjectID, 0, len(Ids))
	for _, Id := range Ids {
		oid, err := primitive.ObjectIDFromHex(Id)
		if err != nil {
			return nil, err
		}
		ObjectIds = append(ObjectIds, oid)
	}
	return ObjectIds, nil
}

func ToObjectId(Id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return oid, nil
}
