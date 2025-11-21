package entity

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type ErrorCode string

const (
	CodeTeamExists  ErrorCode = "TEAM_EXISTS"
	CodePRExists    ErrorCode = "PR_EXISTS"
	CodePRMerged    ErrorCode = "PR_MERGED"
	CodeNotAssigned ErrorCode = "NOT_ASSIGNED"
	CodeNoCandidate ErrorCode = "NO_CANDIDATE"
	CodeNotFound    ErrorCode = "NOT_FOUND"
)
