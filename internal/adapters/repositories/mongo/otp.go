package repositories

import (
	"context"

	"github.com/jmjp/go-rbac/internal/core/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OTPMongoRepository struct {
	db *mongo.Database
}

// NewOTPMongoRepository creates a new instance of OTPMongoRepository.
//
// It takes a pointer to a mongo.Database as a parameter and returns a pointer to OTPMongoRepository.
// The function initializes a collection named "otp" in the given database and creates a unique index on the "code" field.
// The function then returns a pointer to the newly created OTPMongoRepository instance.
func NewOTPMongoRepository(db *mongo.Database) *OTPMongoRepository {
	col := db.Collection("otp")
	col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.M{"code": 1},
			Options: options.Index().SetUnique(true),
		},
	})
	return &OTPMongoRepository{
		db: db,
	}
}

func (r *OTPMongoRepository) Create(ctx context.Context, otp *entities.OTP) error {
	_, err := r.db.Collection("otp").InsertOne(ctx, otp)
	return err
}

func (r *OTPMongoRepository) DeleteByCode(ctx context.Context, code string) error {
	_, err := r.db.Collection("otp").DeleteOne(ctx, bson.M{"code": code})
	return err
}
