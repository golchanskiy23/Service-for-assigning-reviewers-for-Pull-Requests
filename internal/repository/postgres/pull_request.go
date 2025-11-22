package postgres

import (
	"context"
	"errors"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/database"
)

type PullRequestRepository interface {
	//nolint:revive // interface func
	CreatePR(ctx context.Context, pr *entity.PullRequest, reviewerIDs []string) error
	GetPR(ctx context.Context, prID string) (*entity.PullRequest, error)
	PRExists(ctx context.Context, prID string) (bool, error)
	UpdatePR(ctx context.Context, pr *entity.PullRequest) error
	UpdateReviewers(
		ctx context.Context,
		prID string,
		reviewerIDs []string,
	) error
	//nolint:revive // interface func
	GetOpenPRsByReviewer(ctx context.Context, reviewerID string) ([]string, error)
}

type prPGRepository struct {
	db *database.DatabaseSource
}

//nolint:revive // idiomatic constructor sight
func NewPullRequestPGRepository(db *database.DatabaseSource) PullRequestRepository {
	return &prPGRepository{db: db}
}

func (r *prPGRepository) CreatePR(
	ctx context.Context,
	pr *entity.PullRequest,
	reviewerIDs []string,
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
		`SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`,
		pr.PullRequestID,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("PR_EXISTS")
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO pull_requests (pull_request_id, 
                           pull_request_name, 
                           author_id, status)
		 VALUES ($1, $2, $3, $4)`,
		pr.PullRequestID,
		pr.PullRequestName,
		pr.AuthorID,
		pr.Status,
	)
	if err != nil {
		return err
	}

	for _, reviewerID := range reviewerIDs {
		_, err = tx.Exec(ctx,
			`INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
			 VALUES ($1, $2)`,
			pr.PullRequestID, reviewerID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *prPGRepository) GetPR(
	ctx context.Context,
	prID string,
) (*entity.PullRequest, error) {
	var pr entity.PullRequest

	err := r.db.Pool.QueryRow(ctx,
		`SELECT pull_request_id, pull_request_name, author_id, status,
		 created_at, merged_at
		 FROM pull_requests
		 WHERE pull_request_id = $1`,
		prID,
	).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		return nil, errors.New(string(entity.CodeNotFound))
	}

	rows, err := r.db.Pool.Query(ctx,
		`SELECT reviewer_id
		 FROM pr_reviewers
		 WHERE pull_request_id = $1
		 ORDER BY assigned_at`,
		prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pr.AssignedReviewers = []string{}

	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, err
		}

		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewerID)
	}

	return &pr, nil
}

//nolint:revive // func
func (r *prPGRepository) PRExists(ctx context.Context, prID string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx,
		//nolint:revive // monolit sql query
		`SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`, prID).Scan(&exists)

	return exists, err
}

//nolint:revive // func
func (r *prPGRepository) UpdatePR(ctx context.Context, pr *entity.PullRequest) error {
	_, err := r.db.Pool.Exec(ctx,
		`UPDATE pull_requests
		 SET status = $1,
		     merged_at = CASE WHEN $1 = 'MERGED' AND merged_at 
			 IS NULL THEN now() ELSE merged_at END
		 WHERE pull_request_id = $2`,
		pr.Status, pr.PullRequestID)

	return err
}

//nolint:revive // sql query
func (r *prPGRepository) UpdateReviewers(ctx context.Context, prID string, reviewerIDs []string) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	//nolint:errcheck // Rollback in defer is best-effort cleanup
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`DELETE FROM pr_reviewers WHERE pull_request_id = $1`,
		prID)

	if err != nil {
		return err
	}

	for _, reviewerID := range reviewerIDs {
		_, err = tx.Exec(ctx,
			`INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
			 VALUES ($1, $2)`,
			prID, reviewerID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

//nolint:revive // func
func (r *prPGRepository) GetOpenPRsByReviewer(ctx context.Context, reviewerID string) ([]string, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT pr.pull_request_id
		 FROM pr_reviewers prr
		 JOIN pull_requests pr ON pr.pull_request_id = prr.pull_request_id
		 WHERE prr.reviewer_id = $1 AND pr.status = 'OPEN'`,
		reviewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prIDs []string

	for rows.Next() {
		var prID string
		if err := rows.Scan(&prID); err != nil {
			return nil, err
		}

		prIDs = append(prIDs, prID)
	}

	return prIDs, nil
}
