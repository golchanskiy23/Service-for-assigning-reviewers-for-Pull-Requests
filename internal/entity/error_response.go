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
	CodeTeamExists  ErrorCode = "TEAM_EXISTS"  // team
	CodePRExists    ErrorCode = "PR_EXISTS"    // pr
	CodePRMerged    ErrorCode = "PR_MERGED"    // pr
	CodeNotAssigned ErrorCode = "NOT_ASSIGNED" // pr
	CodeNoCandidate ErrorCode = "NO_CANDIDATE" // pr
	CodeNotFound    ErrorCode = "NOT_FOUND"    // all
)
