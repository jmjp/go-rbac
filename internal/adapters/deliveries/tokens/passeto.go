package tokens

import (
	"errors"
	"os"
	"time"

	"github.com/jmjp/go-rbac/internal/core/entities"
	"github.com/jmjp/go-rbac/pkg/random"

	"github.com/o1egl/paseto"
)

func GeneratePasetoToken(userId, email string, expires time.Duration, teams []entities.Member) (*string, error) {
	key := []byte(os.Getenv("TOKEN_SECRET"))

	now := time.Now()
	exp := now.Add(expires)
	nbt := now

	jsonToken := paseto.JSONToken{
		Audience:   os.Getenv("HOST"),
		Issuer:     os.Getenv("HOST"),
		Jti:        random.String(32),
		Subject:    userId,
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  nbt,
	}
	jsonToken.Set("user", userId)
	jsonToken.Set("email", email)

	var jsonTeams []PayloadTeams
	for _, team := range teams {
		jsonTeams = append(jsonTeams, PayloadTeams{
			TeamID: team.Team.ID.Hex(),
			Role:   string(team.Role),
		})
	}

	token, err := paseto.NewV2().Encrypt(key, jsonToken, jsonTeams)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

type Payload struct {
	UserId    string         `json:"user_id"`
	Email     string         `json:"email"`
	Teams     []PayloadTeams `json:"teams"`
	ExpiresIn time.Time      `json:"expires_in"`
}

type PayloadTeams struct {
	TeamID string `json:"team_id"`
	Role   string `json:"role"`
}

func ValidatePasseto(token string) (*Payload, error) {
	key := []byte(os.Getenv("TOKEN_SECRET"))
	var jsonToken paseto.JSONToken
	var teams []PayloadTeams
	if err := paseto.NewV2().Decrypt(token, key, &jsonToken, &teams); err != nil {
		return nil, err
	}
	if jsonToken.Expiration.Before(time.Now()) {
		return nil, errors.New("token expired or invalid")
	}
	if jsonToken.Issuer != os.Getenv("HOST") {
		return nil, errors.New("invalid token claims")
	}
	usr := jsonToken.Get("user")
	email := jsonToken.Get("email")
	if usr == "" {
		return nil, errors.New("invalid token claims")
	}
	payload := &Payload{
		UserId:    usr,
		Email:     email,
		ExpiresIn: jsonToken.Expiration,
		Teams:     teams,
	}
	return payload, nil
}
