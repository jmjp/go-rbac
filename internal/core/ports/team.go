package ports

import (
	"context"

	"github.com/jmjp/go-rbac/internal/core/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TeamRepository defines the methods to interact with team entity
type TeamRepository interface {
	Create(ctx context.Context, team *entities.Team) (*entities.Team, error)
	GetAll(ctx context.Context, id primitive.ObjectID) ([]*entities.Team, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Team, error)
	Update(ctx context.Context, team *entities.Team) error
	UpdateName(ctx context.Context, id primitive.ObjectID, name string) (*entities.Team, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// TeamUsecase defines the methods to interact with team entity
type TeamUsecase interface {
	Create(userId, name string) (*entities.Team, error)
	GetAll(userId string) ([]*entities.Team, error)
	GetByID(id string) (*entities.Team, error)
	Update(userId, teamId, name string) (*entities.Team, error)
	Delete(userId, id string) error
}
