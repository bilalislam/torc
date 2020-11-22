package helpers

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func V4UUIDFromString(s string) (*primitive.Binary, error) {
	u, err := uuid.Parse(s)
	if err == nil {
		b, err := u.MarshalBinary()
		if err == nil {
			return &primitive.Binary{
				Subtype: bsontype.BinaryUUID,
				Data:    b,
			}, nil
		}
		return nil, err
	}
	return nil, err
}
