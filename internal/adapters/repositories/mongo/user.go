package repositories

import (
	"context"
	"time"

	"github.com/jmjp/go-rbac/internal/core/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserMongoRepository struct {
	db *mongo.Database
}

// NewUserMongoRepository creates a new instance of UserMongoRepository.
//
// It takes a pointer to a mongo.Database as a parameter and returns a pointer to UserMongoRepository.
func NewUserMongoRepository(db *mongo.Database) *UserMongoRepository {
	col := db.Collection("user")
	col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"username": 1},
			Options: options.Index().SetUnique(true),
		},
	})
	return &UserMongoRepository{
		db: db,
	}
}

func (r *UserMongoRepository) Create(ctx context.Context, user *entities.User) (out *entities.User, err error) {
	result, err := r.db.Collection("user").InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	err = r.db.Collection("user").FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&out)
	if err != nil {
		return nil, err
	}
	return
}

func (r *UserMongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (user *entities.User, err error) {
	if err := r.db.Collection("user").FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return
}

func (r *UserMongoRepository) GetByEmail(ctx context.Context, email string) (user *entities.User, err error) {
	if err := r.db.Collection("user").FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return
}

func (r *UserMongoRepository) Update(ctx context.Context, user *entities.User) error {
	_, err := r.db.Collection("user").ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	return err
}

func (r *UserMongoRepository) GetByValidOTP(ctx context.Context, code string) (user *entities.User, err error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"code": code,
				"expires_at": bson.M{
					"$gte": time.Now(),
				},
			},
		},
	}
	//lookup
	lookup := bson.M{
		"$lookup": bson.M{
			"from":         "user",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user",
		},
	}
	pipeline = append(pipeline, lookup)
	// unwind
	unwind := bson.M{
		"$unwind": "$user",
	}
	pipeline = append(pipeline, unwind)

	// replace root
	root := bson.M{
		"$replaceRoot": bson.M{
			"newRoot": "$user",
		},
	}
	pipeline = append(pipeline, root)

	cur, err := r.db.Collection("otp").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if cur.Next(ctx) {
		err = cur.Decode(&user)
		if err != nil {
			return nil, err
		}
	}
	return
}
