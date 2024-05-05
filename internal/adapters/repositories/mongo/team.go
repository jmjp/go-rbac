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

type TeamMongoRepository struct {
	db *mongo.Database
}

func NewTeamMongoRepository(db *mongo.Database) *TeamMongoRepository {
	col := db.Collection("team")
	col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.M{"name": 1},
			Options: options.Index().SetUnique(true),
		},
	})
	return &TeamMongoRepository{
		db: db,
	}
}

func (r *TeamMongoRepository) Create(ctx context.Context, team *entities.Team) (out *entities.Team, err error) {
	result, err := r.db.Collection("team").InsertOne(ctx, team)
	if err != nil {
		return nil, err
	}
	err = r.db.Collection("team").FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&out)
	return
}

func (r *TeamMongoRepository) GetAll(ctx context.Context, id primitive.ObjectID) (out []*entities.Team, err error) {
	cursor, err := r.db.Collection("team").Find(ctx, bson.M{"user_id": id})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &out); err != nil {
		return nil, err
	}
	return
}

func (r *TeamMongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (team *entities.Team, err error) {
	if err := r.db.Collection("team").FindOne(ctx, bson.M{"_id": id}).Decode(&team); err != nil {
		return nil, err
	}
	return
}

func (r *TeamMongoRepository) Update(ctx context.Context, team *entities.Team) error {
	_, err := r.db.Collection("team").ReplaceOne(ctx, bson.M{"_id": team.ID}, team)
	return err
}

func (r *TeamMongoRepository) UpdateName(ctx context.Context, id primitive.ObjectID, name string) (team *entities.Team, err error) {
	// por conta da estruturação de dados (extended references), é necessario fazer a atualização atomica.
	updatedAt := time.Now()
	_, err = r.db.Collection("team").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"name": name, "updated_at": updatedAt}})
	if err != nil {
		return nil, err
	}
	_, err = r.db.Collection("user").UpdateMany(ctx, bson.M{"teams.team._id": id}, bson.M{"$set": bson.M{"teams.$.team.name": name, "teams.$.team.updated_at": updatedAt}})
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *TeamMongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.Collection("team").DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	_, err = r.db.Collection("user").UpdateMany(ctx, bson.M{"teams.team._id": id}, bson.M{"$pull": bson.M{"teams": bson.M{"team._id": id}}})
	if err != nil {
		return err
	}
	return err
}
