package util

import (
	"encoding/json"
	"net/http"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
)

func SendError(
	w http.ResponseWriter,
	status int,
	code entity.ErrorCode,
	message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := entity.ErrorResponse{
		Error: entity.ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Log error but can't send another response.
		_ = err
	}
}
