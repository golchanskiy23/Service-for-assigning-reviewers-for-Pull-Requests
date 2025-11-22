//nolint:revive // max-public-structs: Multiple request/response structs are needed for API
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/util"
)

const (
	invalidJSONMsg    = "invalid json"
	contentTypeHeader = "Content-Type"
	encodeErrorMsg    = "failed to encode response"
	applicationJSON   = "application/json"
)

type PRCreateRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type PRCreateResponse struct {
	PR entity.PullRequest `json:"pr"`
}

type PRReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type PRReassignResponse struct {
	ReplacedBy string             `json:"replaced_by"`
	PR         entity.PullRequest `json:"pr"`
}

type PRMergeRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type PRMergeResponse struct {
	PR entity.PullRequest `json:"pr"`
}

func validatePRCreateRequest(req *PRCreateRequest) error {
	if strings.TrimSpace(req.PullRequestID) == "" {
		return errors.New("pull_request_id is required")
	}
	if strings.TrimSpace(req.PullRequestName) == "" {
		return errors.New("pull_request_name is required")
	}
	if strings.TrimSpace(req.AuthorID) == "" {
		return errors.New("author_id is required")
	}
	return nil
}

func validatePRMergeRequest(req *PRMergeRequest) error {
	if strings.TrimSpace(req.PullRequestID) == "" {
		return errors.New("pull_request_id is required")
	}
	return nil
}

func validatePRReassignRequest(req *PRReassignRequest) error {
	if strings.TrimSpace(req.PullRequestID) == "" {
		return errors.New("pull_request_id is required")
	}
	if strings.TrimSpace(req.OldUserID) == "" {
		return errors.New("old_user_id is required")
	}
	return nil
}

func (s *Services) PRCreateHandler(w http.ResponseWriter, r *http.Request) {
	var req PRCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.Log.Warn("failed to decode PR create request", "error", err)
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			invalidJSONMsg,
		)

		return
	}

	if err := validatePRCreateRequest(&req); err != nil {
		s.Log.Warn("invalid PR create request", "error", err)
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			err.Error(),
		)

		return
	}

	ctx := r.Context()

	pr, _, err := s.PRService.CreatePR(ctx, req.PullRequestID,
		req.PullRequestName,
		req.AuthorID)

	if err != nil {
		switch {
		case errors.Is(err, entity.ErrPRExists):
			s.Log.Info("attempt to create PR with existing ID", "pr_id",
				req.PullRequestID)

			util.SendError(
				w,
				http.StatusConflict,
				entity.CodePRExists,
				"PR id already exists",
			)

		case errors.Is(err, entity.ErrNotFound):
			s.Log.Warn("author or team not found for PR creation",
				"author_id", req.AuthorID)

			util.SendError(
				w,
				http.StatusNotFound,
				entity.CodeNotFound,
				"author/team not found",
			)
		default:
			s.Log.Error("failed to create PR", "error", err, "pr_id", req.PullRequestID)
			util.SendError(
				w,
				http.StatusInternalServerError,
				entity.CodeInternalError,
				"internal server error",
			)
		}

		return
	}

	w.Header().Set(contentTypeHeader, applicationJSON)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(PRCreateResponse{PR: *pr}); err != nil {
		s.Log.Error("failed to encode PR create response", "error", err)
		util.SendError(
			w,
			http.StatusInternalServerError,
			entity.CodeInternalError,
			encodeErrorMsg,
		)
	}
}

func (s *Services) PRMergeHandler(w http.ResponseWriter, r *http.Request) {
	var req PRMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.Log.Warn("failed to decode PR merge request", "error", err)
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			invalidJSONMsg,
		)

		return
	}

	if err := validatePRMergeRequest(&req); err != nil {
		s.Log.Warn("invalid PR merge request", "error", err)
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			err.Error(),
		)

		return
	}

	ctx := r.Context()

	pr, err := s.PRService.MergePR(ctx, req.PullRequestID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			s.Log.Warn("PR not found for merge", "pr_id", req.PullRequestID)
			util.SendError(w, http.StatusNotFound, entity.CodeNotFound, "PR not found")
		} else {
			s.Log.Error("failed to merge PR", "error", err, "pr_id", req.PullRequestID)
			util.SendError(w, http.StatusInternalServerError, entity.CodeInternalError, "internal server error")
		}

		return
	}

	w.Header().Set(contentTypeHeader, applicationJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(PRMergeResponse{PR: *pr}); err != nil {
		s.Log.Error("failed to encode PR merge response", "error", err)
		util.SendError(
			w,
			http.StatusInternalServerError,
			entity.CodeInternalError,
			encodeErrorMsg,
		)
	}
}

//nolint:revive // Complex business logic for PR reassignment
func (s *Services) PRReassignHandler(w http.ResponseWriter, r *http.Request) {
	var req PRReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.Log.Warn("failed to decode PR reassign request", "error", err)
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			invalidJSONMsg,
		)

		return
	}

	if err := validatePRReassignRequest(&req); err != nil {
		s.Log.Warn("invalid PR reassign request", "error", err)
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			err.Error(),
		)

		return
	}

	ctx := r.Context()

	pr, replacedBy, err := s.PRService.ReassignReviewer(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			s.Log.Warn("PR or user not found for reassign", "pr_id", req.PullRequestID, "old_user_id", req.OldUserID)
			util.SendError(
				w,
				http.StatusNotFound,
				entity.CodeNotFound,
				"PR or user not found",
			)

		case errors.Is(err, entity.ErrPRMerged):
			s.Log.Info("attempt to reassign on merged PR", "pr_id", req.PullRequestID)
			util.SendError(
				w,
				http.StatusConflict,
				entity.CodePRMerged,
				"cannot reassign on merged PR",
			)

		case errors.Is(err, entity.ErrNotAssigned):
			s.Log.Info("reviewer not assigned to PR", "pr_id", req.PullRequestID, "user_id", req.OldUserID)
			util.SendError(
				w,
				http.StatusConflict,
				entity.CodeNotAssigned,
				"reviewer is not assigned to this PR",
			)

		case errors.Is(err, entity.ErrNoCandidate):
			s.Log.Warn("no candidate for PR reassignment", "pr_id", req.PullRequestID)
			util.SendError(
				w,
				http.StatusConflict,
				entity.CodeNoCandidate,
				"no active replacement candidate in team",
			)

		default:
			s.Log.Error("failed to reassign PR reviewer", "error", err, "pr_id", req.PullRequestID)
			util.SendError(w, http.StatusInternalServerError, entity.CodeInternalError, "internal server error")
		}

		return
	}

	w.Header().Set(contentTypeHeader, applicationJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(PRReassignResponse{
		PR:         *pr,
		ReplacedBy: replacedBy,
	}); err != nil {
		s.Log.Error("failed to encode PR reassign response", "error", err)
		util.SendError(
			w,
			http.StatusInternalServerError,
			entity.CodeInternalError,
			encodeErrorMsg,
		)
	}
}
