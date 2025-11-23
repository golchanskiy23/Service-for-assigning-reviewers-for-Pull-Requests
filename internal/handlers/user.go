package handlers

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/util"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
)

const (
	Zero        = 0
	userIDField = "user_id"
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

func validateUserSetIsActiveRequest(req *UserSetIsActiveRequest) error {
	if strings.TrimSpace(req.UserID) == "" {
		return errors.New("user_id is required")
	}
	return nil
}

func validateUserID(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return errors.New("user_id is required")
	}
	return nil
}

func (s *Services) UserSetIsActiveHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req UserSetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.Log.Warn("failed to decode user set active request", ERROR, err)
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			"invalid json",
		)

		return
	}

	if err := validateUserSetIsActiveRequest(&req); err != nil {
		s.Log.Warn("invalid user set active request", ERROR, err)
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			err.Error(),
		)

		return
	}

	ctx := r.Context()

	user, err := s.UserService.ChangeStatus(ctx, req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			s.Log.Warn("user not found for status change",
				userIDField,
				req.UserID)

			util.SendError(
				w,
				http.StatusNotFound,
				entity.CodeNotFound,
				"user not found",
			)

			return
		}

		s.Log.Error("failed to change user status",
			errFieldName,
			err,
			userIDField,
			req.UserID)

		util.SendError(w,
			http.StatusInternalServerError,
			entity.CodeInternalError,
			"internal server error",
		)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(
		UserSetIsActiveResponse{User: *user},
	); err != nil {
		s.Log.Error("failed to encode user set active response", ERROR, err)
		util.SendError(
			w,
			http.StatusInternalServerError,
			entity.CodeInternalError,
			"failed to encode response",
		)
	}
}

func (s *Services) UserGetReviewHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID := r.URL.Query().Get(userIDField)
	if err := validateUserID(userID); err != nil {
		s.Log.Warn("invalid user get review request", ERROR, err)
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			err.Error(),
		)

		return
	}

	ctx := r.Context()

	id, prs, err := s.UserService.GetPRsAssignedTo(ctx, userID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			s.Log.Warn("user not found for PR retrieval", "user_id", userID)
			util.SendError(
				w,
				http.StatusNotFound,
				entity.CodeNotFound,
				"user not found",
			)

			return
		}

		s.Log.Error("failed to get PRs for user",
			errFieldName,
			err,
			userIDField,
			userID)

		util.SendError(w,
			http.StatusInternalServerError,
			entity.CodeInternalError,
			"internal server error",
		)

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
		s.Log.Error("failed to encode user get review response", ERROR, err)
		util.SendError(w,
			http.StatusInternalServerError,
			entity.CodeInternalError,
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
		s.Log.Warn("failed to decode mass deactivate request", ERROR, err)
		util.SendError(w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			"invalid json")
		return
	}

	if len(req.Users) == Zero {
		s.Log.Warn("empty users list in mass deactivate request")
		util.SendError(w,
			http.StatusBadRequest,
			entity.CodeEmptyRequest,
			"empty request")
		return
	}

	if req.Flag {
		s.Log.Warn("invalid flag for mass deactivate (must be false)")
		util.SendError(w,
			http.StatusBadRequest,
			entity.CodeOnlyDeactivate,
			"only deactivation supported (flag must be false)")
		return
	}

	users := make([]entity.User, 0, len(req.Users))
	for _, u := range req.Users {
		users = append(users, entity.User(u))
	}

	ctx := r.Context()
	if err := s.UserService.MassDeactivate(ctx, users, req.Flag); err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			s.Log.Warn("user not found during mass deactivate")
			util.SendError(w,
				http.StatusNotFound,
				entity.CodeNotFound,
				"user not found")
			return
		case errors.Is(err, entity.ErrUsersFromDifferentTeams):
			s.Log.Warn("users from different teams in mass deactivate")
			util.SendError(w,
				http.StatusBadRequest,
				entity.CodeUsersFromDifferentTeams,
				"users belong to different teams")
			return
		case errors.Is(err, entity.ErrEmptyRequest):
			s.Log.Warn("empty request in mass deactivate")
			util.SendError(w,
				http.StatusBadRequest,
				entity.CodeEmptyRequest,
				"empty request")
			return
		case errors.Is(err, entity.ErrOnlyDeactivate):
			s.Log.Warn("only deactivation supported")
			util.SendError(w,
				http.StatusBadRequest,
				entity.CodeOnlyDeactivate,
				"only deactivation supported")
			return
		default:
			s.Log.Error("failed to mass deactivate users", ERROR, err)
			util.SendError(w,
				http.StatusInternalServerError,
				entity.CodeInternalError,
				"internal server error")
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
		s.Log.Error("failed to encode mass deactivate response", ERROR, err)
		util.SendError(w, http.StatusInternalServerError, entity.CodeInternalError, "failed to encode response")
	}
}
