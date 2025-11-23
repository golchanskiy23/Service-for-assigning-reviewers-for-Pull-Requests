package entity

import "errors"

var (
	ErrPRExists                = errors.New("PR_EXISTS")
	ErrNotFound                = errors.New("NOT_FOUND")
	ErrPRMerged                = errors.New("PR_MERGED")
	ErrTeamExists              = errors.New("TEAM_EXISTS")
	ErrNotAssigned             = errors.New("NOT_ASSIGNED")
	ErrNoCandidate             = errors.New("NO_CANDIDATE")
	ErrEmptyRequest            = errors.New("EMPTY_REQUEST")
	ErrUsersFromDifferentTeams = errors.New("USERS_FROM_DIFFERENT_TEAMS")
	ErrOnlyDeactivate          = errors.New("ONLY_DEACTIVATE")
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type ErrorCode string

const (
	CodeOnlyDeactivate          ErrorCode = "ONLY_DEACTIVATE"
	CodeEmptyRequest            ErrorCode = "EMPTY_REQUEST"
	CodeUsersFromDifferentTeams ErrorCode = "USERS_FROM_DIFFERENT_TEAMS"
	CodeTeamExists              ErrorCode = "TEAM_EXISTS"
	CodePRExists                ErrorCode = "PR_EXISTS"
	CodePRMerged                ErrorCode = "PR_MERGED"
	CodeNotAssigned             ErrorCode = "NOT_ASSIGNED"
	CodeNoCandidate             ErrorCode = "NO_CANDIDATE"
	CodeNotFound                ErrorCode = "NOT_FOUND"
	CodeBadRequest              ErrorCode = "BAD_REQUEST"
	CodeInternalError           ErrorCode = "INTERNAL_ERROR"
)
