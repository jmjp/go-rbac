package entities

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type TeamBuilder struct {
	team *Team
}

// NewTeamBuilder initializes a new TeamBuilder.
//
// No parameters.
// Returns a pointer to a TeamBuilder.
func NewTeamBuilder() *TeamBuilder {
	return &TeamBuilder{
		team: &Team{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func (b *TeamBuilder) WithID(id primitive.ObjectID) *TeamBuilder {
	b.team.ID = id
	return b
}

func (b *TeamBuilder) WithName(name string) *TeamBuilder {
	b.team.Name = name
	return b
}

func (b *TeamBuilder) Build() *Team {
	return b.team
}

func (t *Team) IsValid() error {
	if t.Name == "" {
		return errors.New("name is required")
	}
	return nil
}
