package entity

import "go.mongodb.org/mongo-driver/v2/bson"

// User represents a user document in the "users" MongoDB collection.
type User struct {
	UserID  bson.ObjectID `bson:"_id,omitempty" json:"user_id"`
	GaID    string        `bson:"ga_id"         json:"ga_id"`
	Email   string        `bson:"email"         json:"email"`
	Name    string        `bson:"name"          json:"name"`
	Picture string        `bson:"picture"       json:"picture"`
}
