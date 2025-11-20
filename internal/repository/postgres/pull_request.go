package postgres

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/database"
	"context"
)

type PullRequest struct {
	ID         int
	Title      string
	AuthorID   int
	ReviewerID int
	Status     string
}

type PullRequestRepository interface {
	Create(ctx context.Context, pr *PullRequest) error
	Merge(ctx context.Context, prID int) error
	Reassign(ctx context.Context, prID int, newReviewerID int) error
}

type prPGRepository struct {
	db *database.DatabaseSource
}

func NewPullRequestPGRepository(db *database.DatabaseSource) PullRequestRepository {
	return &prPGRepository{db: db}
}

func (r *prPGRepository) Create(ctx context.Context, pr *PullRequest) error {
	_, err := r.db.Pool.Exec(ctx,
		`INSERT INTO pull_requests (title, author_id, reviewer_id, status)
         VALUES ($1, $2, $3, $4)`,
		pr.Title, pr.AuthorID, pr.ReviewerID, pr.Status)
	return err
}

func (r *prPGRepository) Merge(ctx context.Context, prID int) error {
	_, err := r.db.Pool.Exec(ctx,
		`UPDATE pull_requests SET status='merged' WHERE id=$1`, prID)
	return err
}

func (r *prPGRepository) Reassign(ctx context.Context, prID int, newReviewerID int) error {
	_, err := r.db.Pool.Exec(ctx,
		`UPDATE pull_requests SET reviewer_id=$1 WHERE id=$2`, newReviewerID, prID)
	return err
}

/*func (r *PRRepo) CreatePR(pr *entity.PullRequest) error {
	tx, err := r.db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		_, err = tx.Exec(`
	        INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
	        VALUES ($1, $2, $3, $4)`,
			pr.ID, pr.Name, pr.AuthorID, pr.Status,
		)
		if err != nil {
			return err
		}

		for _, reviewerID := range pr.Reviewers {
			_, err := tx.Exec(`
	            INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
	            VALUES ($1, $2)`,
				pr.ID, reviewerID,
			)
			if err != nil {
				return err
			}
		}

		return tx.Commit()
	return nil
}

func (r *PRRepo) GetPR(id string) (*entity.PullRequest, error) {
	pr := entity.PullRequest{}

		err := r.db.QueryRow(`
	        SELECT pull_request_id, pull_request_name, author_id,
	               status, created_at, merged_at
	        FROM pull_requests
	        WHERE pull_request_id = $1`,
			id,
		).Scan(
			&pr.ID, &pr.Name, &pr.AuthorID,
			&pr.Status, &pr.CreatedAt, &pr.MergedAt,
		)
		if err != nil {
			return nil, err
		}

		// Load reviewers
		rows, err := r.db.Query(`
	        SELECT reviewer_id
	        FROM pr_reviewers
	        WHERE pull_request_id = $1`,
			id,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		pr.Reviewers = []string{}
		for rows.Next() {
			var reviewer string
			rows.Scan(&reviewer)
			pr.Reviewers = append(pr.Reviewers, reviewer)
		}

	return nil, nil
}

func (r *PRRepo) UpdatePR(pr *entity.PullRequest) error {
	/*tx, err := r.db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		_, err = tx.Exec(`
	        UPDATE pull_requests
	        SET status = $1,
	            merged_at = CASE WHEN $1 = 'MERGED' THEN now() ELSE NULL END
	        WHERE pull_request_id = $2`,
			pr.Status, pr.ID,
		)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`DELETE FROM pr_reviewers WHERE pull_request_id = $1`, pr.ID)
		if err != nil {
			return err
		}

		for _, reviewer := range pr.Reviewers {
			_, err := tx.Exec(`
	            INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
	            VALUES ($1, $2)`,
				pr.ID, reviewer,
			)
			if err != nil {
				return err
			}
		}

	return nil
}

func (r *PRRepo) GetActiveReviewers(teamName string, exclude []string) ([]entity.User, error) {
	/*rows, err := r.db.Query(`
	        SELECT user_id, username, team_name, is_active
	        FROM users
	        WHERE team_name = $1
	          AND is_active = TRUE
	          AND NOT (user_id = ANY($2))`,
			teamName, exclude,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		list := []entity.User{}
		for rows.Next() {
			var u entity.User
			rows.Scan(&u.ID, &u.Username, &u.TeamName, &u.Active)
			list = append(list, u)
		}

	return nil, nil
}
*/
