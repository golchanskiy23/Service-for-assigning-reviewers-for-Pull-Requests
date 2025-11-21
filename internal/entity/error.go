package entity

import "errors"

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type ErrorCode string

var (
	ErrUsersFromDifferentTeams = errors.New("users_from_different_teams")
	ErrNoActiveUsersLeft       = errors.New("no_active_users_left")
	ErrUserNotFound            = errors.New("user_not_found")
	ErrIncorrectRequest        = errors.New("incorrect_request")
	ErrUnfamous                = errors.New("unfamous_error")
	ErrTeamMismatch            = errors.New("team_mismatch")
)

const (
	CodeUsersFromDifferentTeams ErrorCode = "USERS_FROM_DIFFERENT_TEAMS"
	CodeInvalidFileFormat       ErrorCode = "INVALID_FILE_FORMAT"
	CodeTeamExists              ErrorCode = "TEAM_EXISTS"
	CodePRExists                ErrorCode = "PR_EXISTS"
	CodePRMerged                ErrorCode = "PR_MERGED"
	CodeNotAssigned             ErrorCode = "NOT_ASSIGNED"
	CodeNoCandidate             ErrorCode = "NO_CANDIDATE"
	CodeNotFound                ErrorCode = "NOT_FOUND"
)
