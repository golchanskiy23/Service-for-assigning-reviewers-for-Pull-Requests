package handlers

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/util"
	"encoding/json"
	"net/http"
	"strconv"
)

type UserSetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserGetReviewResponse struct {
	UserID       string                    `json:"user_id"`
	PullRequests []entity.PullRequestShort `json:"pull_requests"`
}

func (service *Services) UserSetIsActiveHandler(w http.ResponseWriter, r *http.Request) {
	var req UserSetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusNotFound, entity.CodeNotFound, "user not found")
		return
	}

	user, err := service.UserService.ChangeActivateStatus(req.UserID, req.IsActive)
	if err != nil {
		util.SendError(w, http.StatusNotFound, entity.CodeNotFound, "user not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (service *Services) UserGetReviewHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		util.SendError(w, http.StatusNotFound, entity.CodeNotFound, "empty param")
		return
	}
	// string, []pull_request(pull_request_short)
	id, arr := service.UserService.GetPRsAssignedTo(userIDStr)

	resp := UserGetReviewResponse{
		UserID:       id,
		PullRequests: transform(arr),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
