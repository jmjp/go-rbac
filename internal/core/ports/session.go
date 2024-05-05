package ports

import (
	"context"

	"github.com/jmjp/go-rbac/internal/core/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionRepository interface {
	Create(ctx context.Context, session *entities.Session) (*entities.Session, error)
	GetByHash(ctx context.Context, hash string) (*entities.Session, error)
	GetByUserId(ctx context.Context, userId primitive.ObjectID) ([]*entities.Session, error)
	Update(ctx context.Context, session *entities.Session) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
