package util

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"encoding/json"
	"net/http"
)

func SendError(w http.ResponseWriter, status int, code entity.ErrorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := entity.ErrorResponse{
		Error: entity.ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
	json.NewEncoder(w).Encode(resp)
}
