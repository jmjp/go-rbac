package entities

import (
	"time"

	"github.com/jmjp/go-rbac/pkg/random"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OTP struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	Code      string             `json:"code"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type OTPBuilder struct {
	OTP *OTP
}

// NewOTPBuilder creates a new OTPBuilder with a default OTP object.
//
// It generates a random 4-digit code and sets the CreatedAt and UpdatedAt fields to the current time.
//
// Returns a pointer to the newly created OTPBuilder.
func NewOTPBuilder() *OTPBuilder {
	return &OTPBuilder{
		OTP: &OTP{
			Code:      random.Int(4),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func (b *OTPBuilder) WithID(id primitive.ObjectID) *OTPBuilder {
	b.OTP.ID = id
	return b
}

func (b *OTPBuilder) WithCode(code string) *OTPBuilder {
	b.OTP.Code = code
	return b
}

func (b *OTPBuilder) WithExpiresAt(expiresAt time.Time) *OTPBuilder {
	b.OTP.ExpiresAt = expiresAt
	return b
}

func (b *OTPBuilder) WithUserID(userID primitive.ObjectID) *OTPBuilder {
	b.OTP.UserID = userID
	return b
}

func (b *OTPBuilder) Build() *OTP {
	return b.OTP
}
