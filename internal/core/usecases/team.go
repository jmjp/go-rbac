package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/jmjp/go-rbac/internal/core/entities"
	"github.com/jmjp/go-rbac/internal/core/ports"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeamUsecase struct {
	user ports.UserRepository
	team ports.TeamRepository
}

func NewTeamUsecase(user ports.UserRepository, team ports.TeamRepository) *TeamUsecase {
	return &TeamUsecase{user: user, team: team}
}

func (repos *TeamUsecase) Create(userId, name string) (*entities.Team, error) {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer timeout()
	user, err := repos.user.GetByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Blocked {
		return nil, errors.New("user is blocked")
	}
	if len(user.Teams) >= 5 {
		return nil, errors.New("user already has 5 teams")
	}
	builder := entities.NewTeamBuilder().WithName(name).Build()
	if err := builder.IsValid(); err != nil {
		return nil, err
	}
	team, err := repos.team.Create(ctx, builder)
	if err != nil {
		return nil, err
	}
	go func() {
		builder := entities.NewMemberBuilder().WithTeam(*team).WithRole(entities.MemberRoleOwner).Build()
		user.Teams = append(user.Teams, *builder)
		err := repos.user.Update(context.Background(), user)
		if err != nil {
			return
		}
	}()
	return team, nil
}

func (repos *TeamUsecase) GetAll(userId string) ([]*entities.Team, error) {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer timeout()
	teams, err := repos.team.GetAll(ctx, oid)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (repos *TeamUsecase) GetByID(id string) (*entities.Team, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer timeout()
	team, err := repos.team.GetByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (repos *TeamUsecase) Update(userId, teamId, name string) (*entities.Team, error) {
	ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer timeout()
	_, err := repos.ownerPolicy(ctx, userId, teamId)
	if err != nil {
		return nil, err
	}
	// check if name is valid
	builder := entities.NewTeamBuilder().WithName(name).Build()
	if err := builder.IsValid(); err != nil {
		return nil, err
	}
	tid, err := primitive.ObjectIDFromHex(teamId)
	if err != nil {
		return nil, err
	}
	team, err := repos.team.UpdateName(ctx, tid, name)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (repos *TeamUsecase) Delete(userId, teamId string) error {
	ctx, timeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer timeout()
	_, err := repos.ownerPolicy(ctx, userId, teamId)
	if err != nil {
		return err
	}
	oid, err := primitive.ObjectIDFromHex(teamId)
	if err != nil {
		return err
	}
	return repos.team.Delete(ctx, oid)
}

func (repos *TeamUsecase) ownerPolicy(ctx context.Context, userId, teamId string) (*entities.User, error) {
	uid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	user, err := repos.user.GetByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Blocked {
		return nil, errors.New("user is blocked")
	}
	// if user not have teamId return error
	if !user.HasOwner(teamId) {
		return nil, errors.New("user does not have this team or not an owner")
	}
	return user, nil
}
