package auth

import (
	"context"

	"github.com/CollabTed/CollabTed-Backend/internal/auth/repository"
	"github.com/CollabTed/CollabTed-Backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoAuthRepository struct {
	collection *mongo.Collection
}

func NewMongoAuthRepository(db *mongo.Database) repository.AuthRepository {
	return &MongoAuthRepository{
		collection: db.Collection("users"),
	}
}

func (r *MongoAuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *MongoAuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *MongoAuthRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return &user, err
}
