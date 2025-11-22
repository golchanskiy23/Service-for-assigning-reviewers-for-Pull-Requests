package entity

import "errors"

var (
	ErrPRExists                = errors.New("PR_EXISTS")
	ErrPRNotFound              = errors.New("PR_NOT_FOUND")
	ErrNotFound                = errors.New("NOT_FOUND")
	ErrPRMerged                = errors.New("PR_MERGED")
	ErrUserNotFound            = errors.New("USER_NOT_FOUND")
	ErrTeamNotFound            = errors.New("TEAM_NOT_FOUND")
	ErrTeamExists              = errors.New("TEAM_EXISTS")
	ErrUserInactive            = errors.New("USER_INACTIVE")
	ErrNotAssigned             = errors.New("NOT_ASSIGNED")
	ErrNoCandidate             = errors.New("NO_CANDIDATE")
	ErrInvalidRequest          = errors.New("INVALID_REQUEST")
	ErrEmptyRequest            = errors.New("EMPTY_REQUEST")
	ErrUsersFromDifferentTeams = errors.New("USERS_FROM_DIFFERENT_TEAMS")
	ErrOnlyDeactivate          = errors.New("ONLY_DEACTIVATE")
	ErrPRCountError            = errors.New("PR_COUNT_ERROR")
	ErrInternalError           = errors.New("INTERNAL_ERROR")
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
	CodePRCount                 ErrorCode = "PR_COUNT_ERROR"
	CodeUsersFromDifferentTeams ErrorCode = "USERS_FROM_DIFFERENT_TEAMS"
	CodeInvalidFileFormat       ErrorCode = "INVALID_FILE_FORMAT"
	CodeTeamExists              ErrorCode = "TEAM_EXISTS"
	CodePRExists                ErrorCode = "PR_EXISTS"
	CodePRMerged                ErrorCode = "PR_MERGED"
	CodeNotAssigned             ErrorCode = "NOT_ASSIGNED"
	CodeNoCandidate             ErrorCode = "NO_CANDIDATE"
	CodeNotFound                ErrorCode = "NOT_FOUND"
	CodeBadRequest              ErrorCode = "BAD_REQUEST"
	CodeConflict                ErrorCode = "CONFLICT"
	CodeInternalError           ErrorCode = "INTERNAL_ERROR"
)
