package repositories

import (
	"context"

	"github.com/jmjp/go-rbac/internal/core/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SessionMongoRepository struct {
	db *mongo.Database
}

// NewSessionMongoRepository creates a new instance of SessionMongoRepository.
//
// It takes a pointer to a mongo.Database as a parameter and returns a pointer to SessionMongoRepository.
func NewSessionMongoRepository(db *mongo.Database) *SessionMongoRepository {
	col := db.Collection("session")
	col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.M{"hash": 1},
			Options: options.Index().SetUnique(true),
		},
	})
	return &SessionMongoRepository{
		db: db,
	}
}

func (r *SessionMongoRepository) Create(ctx context.Context, session *entities.Session) (out *entities.Session, err error) {
	result, err := r.db.Collection("session").InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}
	err = r.db.Collection("session").FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&out)
	return
}

func (r *SessionMongoRepository) GetByHash(ctx context.Context, hash string) (out *entities.Session, err error) {
	if err := r.db.Collection("session").FindOne(ctx, bson.M{"hash": hash}).Decode(&out); err != nil {
		return nil, err
	}
	return
}

func (r *SessionMongoRepository) GetByUserId(ctx context.Context, userID primitive.ObjectID) (out []*entities.Session, err error) {
	cursor, err := r.db.Collection("session").Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &out); err != nil {
		return nil, err
	}
	return
}

func (r *SessionMongoRepository) Update(ctx context.Context, session *entities.Session) error {
	_, err := r.db.Collection("session").ReplaceOne(ctx, bson.M{"_id": session.ID}, session)
	return err
}

func (r *SessionMongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.Collection("session").DeleteOne(ctx, bson.M{"_id": id})
	return err
}
