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
)
