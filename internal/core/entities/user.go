package entities

import (
	"errors"
	"regexp"
	"time"

	"github.com/jmjp/go-rbac/pkg/random"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Avatar    *string            `json:"avatar"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	Blocked   bool               `json:"blocked"`
	Teams     []Member           `json:"teams" bson:"teams"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserBuilder struct {
	user *User
}

// NewUserBuilder initializes a new UserBuilder.
//
// It returns a pointer to a UserBuilder with a new User instance.
// The User instance has a randomly generated username, set to false for blocked,
// and the current time for created and updated.
// The function returns a pointer to the newly created UserBuilder.
func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		user: &User{
			Username:  "rbac_" + random.String(6),
			Blocked:   false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func (b *UserBuilder) WithID(id primitive.ObjectID) *UserBuilder {
	b.user.ID = id
	return b
}

func (b *UserBuilder) WithUsername(username *string) *UserBuilder {
	if username == nil {
		return b
	}
	b.user.Username = *username
	return b
}

func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.user.Email = email
	return b
}

func (b *UserBuilder) WithAvatar(avatar *string) *UserBuilder {
	if avatar == nil {
		return b
	}
	b.user.Avatar = avatar
	return b
}

func (b *UserBuilder) WithBlocked(blocked bool) *UserBuilder {
	b.user.Blocked = blocked
	return b
}

func (b *UserBuilder) WithTeams(teams []Member) *UserBuilder {
	b.user.Teams = teams
	return b
}

func (b *UserBuilder) Build() *User {
	return b.user
}

func (u *User) IsValid() error {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(u.Username) {
		return errors.New("invalid username")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email")
	}
	return nil
}

func (u *User) HasTeam(teamId string) bool {
	for _, t := range u.Teams {
		if t.Team.ID.Hex() == teamId {
			return true
		}
	}
	return false
}

func (u *User) HasOwner(teamId string) bool {
	for _, t := range u.Teams {
		if t.Team.ID.Hex() == teamId && t.Role == MemberRoleOwner {
			return true
		}
	}
	return false
}
