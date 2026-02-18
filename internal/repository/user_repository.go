package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"brainbash_backend/internal/model/entity"
)

const usersCollection = "users"

// UserRepository handles MongoDB operations for the users collection.
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new UserRepository with the given database.
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection(usersCollection),
	}
}

// UpsertByGaID inserts a new user or updates an existing one matched by ga_id.
// Returns the upserted/found user.
func (r *UserRepository) UpsertByGaID(ctx context.Context, user *entity.User) (*entity.User, error) {
	filter := bson.M{"ga_id": user.GaID}
	update := bson.M{
		"$set": bson.M{
			"email":   user.Email,
			"name":    user.Name,
			"picture": user.Picture,
		},
		"$setOnInsert": bson.M{
			"ga_id": user.GaID,
		},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result entity.User
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}

	return &result, nil
}

// FindByEmail finds a user by email.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}

// FindByUserID finds a user by their MongoDB ObjectID.
func (r *UserRepository) FindByUserID(ctx context.Context, userID bson.ObjectID) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by user_id: %w", err)
	}
	return &user, nil
}
