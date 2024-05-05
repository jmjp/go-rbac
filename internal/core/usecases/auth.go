package usecases

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/jmjp/go-rbac/internal/core/entities"
	"github.com/jmjp/go-rbac/internal/core/ports"
	"github.com/jmjp/go-rbac/pkg/random"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthUseCase struct {
	user    ports.UserRepository
	otp     ports.OTPRepository
	session ports.SessionRepository
}

var (
	errInvalidCodeOrEmail = errors.New("invalid code or email")
	errUserIsBlocked      = errors.New("user is blocked")
	errInvalidUserOrBlock = errors.New("invalid user or block")
)

func NewAuthUseCase(user ports.UserRepository, otp ports.OTPRepository, session ports.SessionRepository) *AuthUseCase {
	return &AuthUseCase{user: user, otp: otp, session: session}
}

func (repos *AuthUseCase) Login(email string, avatar, username *string) (*string, error) {
	ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer timeout()
	user, err := repos.user.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if user != nil && user.Blocked {
		return nil, errUserIsBlocked
	}
	if user == nil {
		builder := entities.NewUserBuilder().WithEmail(email).WithAvatar(avatar).WithUsername(username).Build()
		if err := builder.IsValid(); err != nil {
			return nil, err
		}
		user, err = repos.user.Create(ctx, builder)
		if err != nil {
			return nil, err
		}
	}
	otp := "1234"
	if os.Getenv("env") != "dev" {
		otp = random.Int(4)
	}
	builder := entities.NewOTPBuilder().WithCode(otp).WithUserID(user.ID).WithExpiresAt(time.Now().Add(5 * time.Minute)).Build()
	if err := repos.otp.Create(ctx, builder); err != nil {
		return nil, err
	}
	message := "an otp has been sent to your email"
	return &message, nil
}

func (repos *AuthUseCase) Verify(code, email, ip, ua string) (*entities.User, *entities.Session, error) {
	ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer timeout()
	user, err := repos.user.GetByValidOTP(ctx, code)
	if err != nil {
		return nil, nil, errInvalidCodeOrEmail
	}
	if user == nil || user.Blocked || user.Email != email {
		return nil, nil, errInvalidUserOrBlock
	}
	builder := entities.NewSessionBuilder().WithUserID(user.ID).WithIP(ip).WithExpiresAt(time.Now().Add(7 * 24 * time.Hour)).WithUserAgent(ua).WithHash(random.String(64)).Build()
	session, err := repos.session.Create(ctx, builder)
	if err != nil {
		return nil, nil, err
	}
	go func() {
		_ = repos.otp.DeleteByCode(context.Background(), code)
	}()
	return user, session, nil
}

func (repos *AuthUseCase) Refresh(hash string) (*entities.User, *entities.Session, error) {
	ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer timeout()
	session, err := repos.session.GetByHash(ctx, hash)
	if err != nil {
		return nil, nil, err
	}
	if err := session.IsValidYet(); err != nil {
		return nil, nil, err
	}
	user, err := repos.user.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, nil, err
	}
	if user.Blocked {
		return nil, nil, errInvalidUserOrBlock
	}

	nextThreeDays := time.Now().Add(3 * 24 * time.Hour)
	if nextThreeDays.After(session.ExpiresAt) {
		session.Hash = random.String(64)
		session.ExpiresAt = time.Now().Add(7 * 24 * time.Hour) // 7 days
		go func() {
			if err := repos.session.Update(ctx, session); err != nil {
				return
			}
		}()
	}
	return user, session, nil
}

func (repos *AuthUseCase) Sessions(userId string) ([]*entities.Session, error) {
	ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer timeout()
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	sessions, err := repos.session.GetByUserId(ctx, oid)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (repos *AuthUseCase) Logout(userId string, hash string) error {
	ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer timeout()
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	sessions, err := repos.session.GetByUserId(ctx, oid)
	if err != nil {
		return err
	}
	if len(sessions) == 0 {
		return errors.New("no sessions found")
	}
	for _, session := range sessions {
		if session.Hash == hash {
			err = repos.session.Delete(context.Background(), session.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
