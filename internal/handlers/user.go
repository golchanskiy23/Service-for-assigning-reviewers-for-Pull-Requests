package handlers

import (
	"encoding/json"
	"net/http"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/util"
)

const (
	Zero = 0
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

type UserMassChangeRequest struct {
	Users []entity.UserItem `json:"users"`
	Flag  bool              `json:"flag"`
}

type UserMassChangeResponse struct {
	Deactivated []string `json:"deactivated_user_ids"`
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

//nolint:revive // easy for reading function
func (s *Services) UsersMassDeactivateHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req UserMassChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w,
			http.StatusBadRequest,
			entity.CodeNotFound,
			"invalid json")
		return
	}

	if len(req.Users) == Zero {
		util.SendError(w,
			http.StatusBadRequest,
			entity.CodeNotFound,
			"empty request")
		return
	}

	if req.Flag {
		util.SendError(w,
			http.StatusBadRequest,
			entity.CodeNotFound,
			"only deactivation supported (flag must be false)")
		return
	}

	users := make([]entity.User, 0, len(req.Users))
	for _, u := range req.Users {
		users = append(users, entity.User(u))
	}

	ctx := r.Context()
	if err := s.UserService.MassDeactivate(ctx, users, req.Flag); err != nil {
		switch err.Error() {
		case string(entity.CodeNotFound):
			util.SendError(w,
				http.StatusNotFound,
				entity.CodeNotFound,
				"user not found")
			return
		case string(entity.CodeUsersFromDifferentTeams):
			util.SendError(w,
				http.StatusBadRequest,
				entity.CodeNotFound,
				"users belong to different teams")
			return
		case string(entity.CodeEmptyRequest):
			util.SendError(w,
				http.StatusBadRequest,
				entity.CodeNotFound,
				"empty request")
			return
		case string(entity.CodeOnlyDeactivate):
			util.SendError(w,
				http.StatusBadRequest,
				entity.CodeNotFound,
				"only deactivation supported")
			return
		default:
			util.SendError(w,
				http.StatusInternalServerError,
				entity.CodeNotFound,
				err.Error())
			return
		}
	}

	userIDs := make([]string, 0, len(req.Users))
	for _, u := range req.Users {
		userIDs = append(userIDs, u.UserID)
	}

	resp := UserMassChangeResponse{Deactivated: userIDs}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		util.SendError(w, http.StatusInternalServerError, entity.CodeNotFound, "failed to encode response")
	}
}
