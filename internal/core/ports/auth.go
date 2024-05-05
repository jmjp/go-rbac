package ports

import "github.com/jmjp/go-rbac/internal/core/entities"

type AuthUsecase interface {
	Login(email string, avatar, username *string) (*string, error)
	Verify(code, email, ip, ua string) (*entities.User, *entities.Session, error)
	Refresh(hash string) (*entities.User, *entities.Session, error)
	Sessions(userId string) ([]*entities.Session, error)
	Logout(userId string, hash string) error
}
