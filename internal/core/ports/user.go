package ports

import (
	"context"

	"github.com/jmjp/go-rbac/internal/core/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) (*entities.User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByValidOTP(ctx context.Context, otp string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
}
