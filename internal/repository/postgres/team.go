package postgres

import (
	"context"
	"errors"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/database"
)

type TeamRepository interface {
	AddTeam(ctx context.Context, team *entity.Team) error
	GetTeam(ctx context.Context, teamName string) (*entity.Team, error)
	TeamExists(ctx context.Context, teamName string) (bool, error)
}

type teamPGRepository struct {
	db *database.DatabaseSource
}

func NewTeamPGRepository(db *database.DatabaseSource) TeamRepository {
	return &teamPGRepository{db: db}
}

func (r *teamPGRepository) AddTeam(
	ctx context.Context,
	team *entity.Team,
) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	//nolint:errcheck // Rollback in defer is best-effort cleanup
	defer tx.Rollback(ctx)

	var exists bool

	err = tx.QueryRow(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`,
		team.TeamName,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return errors.New(string(entity.CodeTeamExists))
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO teams (team_name) VALUES ($1)`,
		team.TeamName,
	)

	if err != nil {
		return err
	}

	for _, member := range team.Members {
		_, err = tx.Exec(ctx,
			`INSERT INTO users (user_id, username, team_name, is_active)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (user_id) DO UPDATE SET
			 username = EXCLUDED.username,
			 team_name = EXCLUDED.team_name,
			 is_active = EXCLUDED.is_active`,
			member.UserID, member.Username, team.TeamName, member.IsActive)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *teamPGRepository) GetTeam(
	ctx context.Context,
	teamName string,
) (*entity.Team, error) {
	var exists bool

	err := r.db.Pool.QueryRow(ctx,

		`SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`, teamName).
		Scan(&exists)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New(string(entity.CodeNotFound))
	}

	rows, err := r.db.Pool.Query(ctx,
		`SELECT user_id, username, is_active
		 FROM users
		 WHERE team_name = $1
		 ORDER BY user_id`,
		teamName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []entity.TeamMember

	for rows.Next() {
		var member entity.TeamMember
		if err := rows.Scan(
			&member.UserID,
			&member.Username,
			&member.IsActive,
		); err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return &entity.Team{
		TeamName: teamName,
		Members:  members,
	}, nil
}

func (r *teamPGRepository) TeamExists(
	ctx context.Context,
	teamName string,
) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`,
		teamName,
	).Scan(&exists)

	return exists, err
}
