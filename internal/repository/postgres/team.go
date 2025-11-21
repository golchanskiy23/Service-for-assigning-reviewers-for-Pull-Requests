package postgres

import (
	"database/sql"
)

type TeamRepo struct {
	db *sql.DB
}

func NewTeamRepo(db *sql.DB) *TeamRepo {
	return &TeamRepo{db: db}
}

/*
func (r *TeamRepo) CreateTeam(name string, users []entity.User) (*entity.Team, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO teams (team_name) VALUES ($1)`, name)
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		_, err = tx.Exec(`
            INSERT INTO users (user_id, username, team_name, is_active)
            VALUES ($1, $2, $3, TRUE)`,
			u.ID, u.Username, name,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &entity.Team{
		Name: name,
	}, nil
}

func (r *TeamRepo) GetTeam(name string) (*entity.Team, error) {
	t := entity.Team{}
	err := r.db.QueryRow(`
        SELECT team_name, created_at
        FROM teams
        WHERE team_name = $1`,
		name,
	).Scan(&t.Name, &t.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &t, nil
}*/
