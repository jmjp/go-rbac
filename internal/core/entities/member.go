package entities

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberRole string

const (
	MemberRoleOwner    MemberRole = "owner"
	MemberRoleInternal MemberRole = "internal"
	MemberRoleExternal MemberRole = "external"
)

type Member struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Team     Team               `json:"team" bson:"team,omitempty"`
	Role     MemberRole         `json:"role" bson:"role"`
	JoinedAt time.Time          `json:"joined_at" bson:"joined_at"`
}

type MemberBuilder struct {
	member *Member
}

// NewMemberBuilder creates a new MemberBuilder with a default Member object.
//
// It sets the CreatedAt and UpdatedAt fields of the Member object to the current time.
//
// Returns a pointer to the newly created MemberBuilder.
func NewMemberBuilder() *MemberBuilder {
	return &MemberBuilder{
		member: &Member{
			JoinedAt: time.Now(),
		},
	}
}

func (b *MemberBuilder) WithID(id primitive.ObjectID) *MemberBuilder {
	b.member.ID = id
	return b
}

func (b *MemberBuilder) WithTeam(team Team) *MemberBuilder {
	b.member.Team = team
	return b
}

func (b *MemberBuilder) WithRole(role MemberRole) *MemberBuilder {
	b.member.Role = role
	return b
}

func (b *MemberBuilder) Build() *Member {
	return b.member
}

func (m *Member) IsValid() error {
	if m.Team.ID.IsZero() {
		return errors.New("team is required")
	}
	if m.Role == "" {
		return errors.New("role is required")
	}
	return nil
}
