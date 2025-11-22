package entity

import "time"

type PRReviewer struct {
	AssignedAt    time.Time `db:"assigned_at"`
	PullRequestID string    `db:"pull_request_id"`
	ReviewerID    string    `db:"reviewer_id"`
}
