package handlers

import (
	"encoding/json"
	"net/http"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/util"
)

type UserSetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserSetIsActiveResponse struct {
	User entity.User `json:"user"`
}

type UserGetReviewResponse struct {
	UserID       string                    `json:"user_id"`
	PullRequests []entity.PullRequestShort `json:"pull_requests"`
}

func (s *Services) UserSetIsActiveHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req UserSetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeNotFound,
			"invalid json",
		)

		return
	}

	ctx := r.Context()

	user, err := s.UserService.ChangeStatus(ctx, req.UserID, req.IsActive)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			util.SendError(
				w,
				http.StatusNotFound,
				entity.CodeNotFound,
				"user not found",
			)

			return
		}

		util.SendError(w,
			http.StatusInternalServerError,
			entity.CodeNotFound,
			err.Error())

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(
		UserSetIsActiveResponse{User: *user},
	); err != nil {
		util.SendError(
			w,
			http.StatusInternalServerError,
			entity.CodeNotFound,
			"failed to encode response",
		)
	}
}

func (s *Services) UserGetReviewHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeNotFound,
			"user_id parameter is required",
		)

		return
	}

	ctx := r.Context()

	id, prs, err := s.UserService.GetPRsAssignedTo(ctx, userID)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			util.SendError(
				w,
				http.StatusNotFound,
				entity.CodeNotFound,
				"user not found",
			)

			return
		}

		util.SendError(w,
			http.StatusInternalServerError,
			entity.CodeNotFound,
			err.Error())

		return
	}

	pullRequests := make([]entity.PullRequestShort, len(prs))
	for i, pr := range prs {
		pullRequests[i] = *pr
	}

	resp := UserGetReviewResponse{
		UserID:       id,
		PullRequests: pullRequests,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		util.SendError(w,
			http.StatusInternalServerError,
			entity.CodeNotFound,
			"failed to encode response")
	}
}
