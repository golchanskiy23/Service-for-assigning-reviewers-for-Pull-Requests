package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/database"
)

const (
	initialReviewerCap = 0
)

type UserRepository interface {
	GetUser(ctx context.Context, userID string) (*entity.User, error)
	SetIsActive(ctx context.Context, userID string, active bool) error
	GetActiveUsersByTeam(
		ctx context.Context,
		teamName string,
		exclude []string,
	) ([]*entity.User, error)
	//nolint:revive // monolith func
	GetPRsForReviewer(ctx context.Context, userID string) ([]*entity.PullRequestShort, error)
	//nolint:revive // monolith func
	MassDeactivateAndReassign(ctx context.Context, teamName string, userIDs []string) error
}

type userPGRepository struct {
	db *database.DatabaseSource
}

func NewUserPGRepository(db *database.DatabaseSource) UserRepository {
	return &userPGRepository{db: db}
}

func (r *userPGRepository) GetUser(
	ctx context.Context,
	userID string,
) (*entity.User, error) {
	var user entity.User

	err := r.db.Pool.QueryRow(ctx,
		`SELECT user_id, username, team_name, is_active
		 FROM users WHERE user_id = $1`,
		userID,
	).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		return nil, errors.New(string(entity.CodeNotFound))
	}

	return &user, nil
}

func (r *userPGRepository) SetIsActive(
	ctx context.Context,
	userID string,
	active bool,
) error {
	result, err := r.db.Pool.Exec(ctx,
		`UPDATE users SET is_active = $1 WHERE user_id = $2`, active, userID,
	)
	if err != nil {
		return err
	}

	const noRowsAffected = 0
	if result.RowsAffected() == noRowsAffected {
		return errors.New(string(entity.CodeNotFound))
	}

	return nil
}

func (r *userPGRepository) GetActiveUsersByTeam(
	ctx context.Context,
	teamName string,
	exclude []string,
) ([]*entity.User, error) {
	query := `
    	SELECT
        	user_id,
        	username,
        	team_name,
        	is_active
    	FROM users
    	WHERE team_name = $1
      		AND is_active = TRUE
	`
	args := []interface{}{teamName}

	const emptySlice = 0
	if len(exclude) > emptySlice {
		query += ` AND NOT (user_id = ANY($2::text[]))`
		args = append(args, exclude)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User

	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
			&user.UserID,
			&user.Username,
			&user.TeamName,
			&user.IsActive,
		); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (r *userPGRepository) GetPRsForReviewer(
	ctx context.Context,
	userID string,
) ([]*entity.PullRequestShort, error) {
	rows, err := r.db.Pool.Query(ctx,
		//nolint:revive // sql query
		`SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
		 FROM pr_reviewers prr
		 JOIN pull_requests pr ON pr.pull_request_id = prr.pull_request_id
		 WHERE prr.reviewer_id = $1
		 ORDER BY pr.created_at DESC`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*entity.PullRequestShort

	for rows.Next() {
		var pr entity.PullRequestShort
		if err := rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&pr.Status,
		); err != nil {
			return nil, err
		}

		prs = append(prs, &pr)
	}

	return prs, nil
}

//nolint:revive // to heavy for coordination func
func (r *userPGRepository) MassDeactivateAndReassign(
	ctx context.Context,
	teamName string,
	userIDs []string,
) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	//nolint:errcheck // best-effort
	defer tx.Rollback(ctx)

	if e := r.deactivateUsers(ctx, tx, teamName, userIDs); e != nil {
		return e
	}

	prs, err := r.fetchAffectedPRs(ctx, tx, userIDs)
	if err != nil {
		return err
	}

	for _, pr := range prs {
		remaining := filterRemainingReviewers(pr.Reviewers, userIDs)

		if err := r.removeReviewersFromPR(ctx, tx, pr.ID, userIDs); err != nil {
			return err
		}

		if len(remaining) == 0 {
			if err := r.assignFallbackReviewer(ctx,
				tx,
				pr.ID,
				teamName,
				userIDs); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

//nolint:revive // useless linter here
func (r *userPGRepository) deactivateUsers(
	ctx context.Context,
	tx pgx.Tx,
	teamName string,
	userIDs []string,
) error {
	_, err := tx.Exec(ctx,
		`UPDATE users
         SET is_active = FALSE
         WHERE user_id = ANY($1::text[]) AND team_name = $2`,
		userIDs, teamName,
	)
	return err
}

// PRInfo stores PR ID and current reviewers.
type PRInfo struct {
	ID        string
	Reviewers []string
}

//nolint:revive // useless linter here
func (r *userPGRepository) fetchAffectedPRs(
	ctx context.Context,
	tx pgx.Tx,
	userIDs []string,
) ([]PRInfo, error) {
	rows, err := tx.Query(ctx,
		`SELECT pr.pull_request_id,
                array_agg(prr.reviewer_id ORDER BY prr.assigned_at)
         FROM pr_reviewers prr
         JOIN pull_requests pr ON pr.pull_request_id = prr.pull_request_id
         WHERE pr.status = 'OPEN'
           AND prr.reviewer_id = ANY($1::text[])
         GROUP BY pr.pull_request_id`,
		userIDs,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []PRInfo
	for rows.Next() {
		var pr PRInfo
		if err := rows.Scan(&pr.ID, &pr.Reviewers); err != nil {
			return nil, err
		}
		result = append(result, pr)
	}
	return result, rows.Err()
}

func filterRemainingReviewers(reviewers, disabled []string) []string {
	disabledSet := make(map[string]struct{}, len(disabled))
	for _, id := range disabled {
		disabledSet[id] = struct{}{}
	}

	res := make([]string, initialReviewerCap, len(reviewers))
	for _, rID := range reviewers {
		if _, skip := disabledSet[rID]; !skip {
			res = append(res, rID)
		}
	}
	return res
}

//nolint:revive // useless linter here
func (r *userPGRepository) removeReviewersFromPR(
	ctx context.Context,
	tx pgx.Tx,
	prID string,
	userIDs []string,
) error {
	_, err := tx.Exec(ctx,
		`DELETE FROM pr_reviewers
         WHERE pull_request_id = $1
           AND reviewer_id = ANY($2::text[])`,
		prID, userIDs,
	)
	return err
}

//nolint:revive // useless linter here
func (r *userPGRepository) assignFallbackReviewer(
	ctx context.Context,
	tx pgx.Tx,
	prID string,
	teamName string,
	excluded []string,
) error {
	var candidate string
	err := tx.QueryRow(ctx,
		`SELECT user_id
         FROM users
         WHERE team_name = $1
           AND is_active = TRUE
           AND NOT (user_id = ANY($2::text[]))
         LIMIT 1`,
		teamName, excluded,
	).Scan(&candidate)

	if err != nil || candidate == "" {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
         VALUES ($1, $2)`,
		prID, candidate,
	)
	return err
}
