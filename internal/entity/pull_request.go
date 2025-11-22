package entity

import "time"

const (
	OPEN   PRStatus = "OPEN"
	MERGED PRStatus = "MERGED"
)

type PRStatus string

type PullRequest struct {
	CreatedAt         *time.Time `json:"created_at"`
	MergedAt          *time.Time `json:"merged_at"`
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            PRStatus   `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
}

type PullRequestShort struct {
	PullRequestID   string   `json:"pull_request_id"`
	PullRequestName string   `json:"pull_request_name"`
	AuthorID        string   `json:"author_id"`
	Status          PRStatus `json:"status"`
}
