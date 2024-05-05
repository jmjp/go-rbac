package entities

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	Hash      string             `json:"hash"`
	IP        string             `json:"ip"`
	UserAgent string             `json:"user_agent" bson:"user_agent"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type SessionBuilder struct {
	session *Session
}

// NewSessionBuilder initializes a new SessionBuilder.
//
// Returns a pointer to a SessionBuilder.
func NewSessionBuilder() *SessionBuilder {
	return &SessionBuilder{
		session: &Session{
			IP:        "127.0.0.1",
			UserAgent: "unknow",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func (b *SessionBuilder) WithID(id primitive.ObjectID) *SessionBuilder {
	b.session.ID = id
	return b
}

func (b *SessionBuilder) WithUserID(userID primitive.ObjectID) *SessionBuilder {
	b.session.UserID = userID
	return b
}

func (b *SessionBuilder) WithHash(hash string) *SessionBuilder {
	b.session.Hash = hash
	return b
}

func (b *SessionBuilder) WithExpiresAt(expiresAt time.Time) *SessionBuilder {
	b.session.ExpiresAt = expiresAt
	return b
}

func (b *SessionBuilder) WithIP(ip string) *SessionBuilder {
	b.session.IP = ip
	return b
}

func (b *SessionBuilder) WithUserAgent(userAgent string) *SessionBuilder {
	b.session.UserAgent = userAgent
	return b
}

func (b *SessionBuilder) Build() *Session {
	return b.session
}

func (s *Session) IsValid() error {
	if s.ExpiresAt.Before(time.Now()) {
		return errors.New("session expired")
	}
	if s.UserID.IsZero() {
		return errors.New("user id is required")
	}
	return nil
}

func (s *Session) IsValidYet() error {
	if s.ExpiresAt.Before(time.Now()) {
		return errors.New("session expired")
	}
	return nil
}
