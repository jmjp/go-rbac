package ports

import (
	"context"

	"github.com/jmjp/go-rbac/internal/core/entities"
)

type OTPRepository interface {
	Create(ctx context.Context, otp *entities.OTP) error
	DeleteByCode(ctx context.Context, code string) error
}
