package service

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"brainbash_backend/internal/model/entity"
	"brainbash_backend/internal/repository"
)

// UserService handles business logic for user operations.
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// FindByEmail looks up a user by email.
func (s *UserService) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// UpsertFromGoogleLogin creates or updates a user from Google login info.
// Returns the persisted user document.
func (s *UserService) UpsertFromGoogleLogin(ctx context.Context, googleUser *GoogleUserInfo) (*entity.User, error) {
	user := &entity.User{
		GaID:    googleUser.Sub,
		Email:   googleUser.Email,
		Name:    googleUser.Name,
		Picture: googleUser.Picture,
	}
	return s.userRepo.UpsertByGaID(ctx, user)
}

// FindByUserID looks up a user by their MongoDB ObjectID.
func (s *UserService) FindByUserID(ctx context.Context, userID string) (*entity.User, error) {
	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	return s.userRepo.FindByUserID(ctx, objID)
}
