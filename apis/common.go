package apis

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func primitiveToString(p primitive.ObjectID) string {
	return p.Hex()
}

func stringToPrimitive(s string) primitive.ObjectID {
	a, _ := primitive.ObjectIDFromHex(s)
	return a
}
