//nolint:revive // max-public-structs: Multiple request/response structs are needed for API
package handlers

import (
	"encoding/json"
	"net/http"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/util"
)

const (
	invalidJSONMsg    = "invalid json"
	notFoundErrorMsg  = "NOT_FOUND"
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
	PR         entity.PullRequest `json:"pr"`
	ReplacedBy string             `json:"replaced_by"`
}

type PRMergeRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type PRMergeResponse struct {
	PR entity.PullRequest `json:"pr"`
}

func (s *Services) PRCreateHandler(w http.ResponseWriter, r *http.Request) {
	var req PRCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeNotFound,
			invalidJSONMsg,
		)

		return
	}

	ctx := r.Context()

	pr, _, err := s.PRService.CreatePR(ctx, req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		switch err.Error() {
		case "PR_EXISTS":
			util.SendError(
				w,
				http.StatusConflict,
				entity.CodePRExists,
				"PR id already exists",
			)
		case notFoundErrorMsg:
			util.SendError(
				w,
				http.StatusNotFound,
				entity.CodeNotFound,
				"author/team not found",
			)
		default:
			util.SendError(
				w,
				http.StatusInternalServerError,
				entity.CodeNotFound,
				err.Error(),
			)
		}

		return
	}

	w.Header().Set(contentTypeHeader, applicationJSON)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(PRCreateResponse{PR: *pr}); err != nil {
		util.SendError(
			w,
			http.StatusInternalServerError,
			entity.CodeNotFound,
			encodeErrorMsg,
		)
	}
}

func (s *Services) PRMergeHandler(w http.ResponseWriter, r *http.Request) {
	var req PRMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeNotFound,
			invalidJSONMsg,
		)

		return
	}

	ctx := r.Context()

	pr, err := s.PRService.MergePR(ctx, req.PullRequestID)
	if err != nil {
		if err.Error() == notFoundErrorMsg {
			util.SendError(w, http.StatusNotFound, entity.CodeNotFound, "PR not found")
		} else {
			util.SendError(w, http.StatusInternalServerError, entity.CodeNotFound, err.Error())
		}

		return
	}

	w.Header().Set(contentTypeHeader, applicationJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(PRMergeResponse{PR: *pr}); err != nil {
		util.SendError(
			w,
			http.StatusInternalServerError,
			entity.CodeNotFound,
			encodeErrorMsg,
		)
	}
}

//nolint:revive // Complex business logic for PR reassignment
func (s *Services) PRReassignHandler(w http.ResponseWriter, r *http.Request) {
	var req PRReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeNotFound,
			invalidJSONMsg,
		)

		return
	}

	ctx := r.Context()

	pr, replacedBy, err := s.PRService.ReassignReviewer(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {
		if err.Error() == notFoundErrorMsg {
			util.SendError(
				w,
				http.StatusNotFound,
				entity.CodeNotFound,
				"PR or user not found",
			)

			return
		}

		if err.Error() == "PR_MERGED" {
			util.SendError(
				w,
				http.StatusConflict,
				entity.CodePRMerged,
				"cannot reassign on merged PR",
			)

			return
		}

		if err.Error() == "NOT_ASSIGNED" {
			util.SendError(
				w,
				http.StatusConflict,
				entity.CodeNotAssigned,
				"reviewer is not assigned to this PR",
			)

			return
		}

		if err.Error() == "NO_CANDIDATE" {
			util.SendError(
				w,
				http.StatusConflict,
				entity.CodeNoCandidate,
				"no active replacement candidate in team",
			)

			return
		}

		util.SendError(w, http.StatusInternalServerError, entity.CodeNotFound, err.Error())

		return
	}

	w.Header().Set(contentTypeHeader, applicationJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(PRReassignResponse{
		PR:         *pr,
		ReplacedBy: replacedBy,
	}); err != nil {
		util.SendError(
			w,
			http.StatusInternalServerError,
			entity.CodeNotFound,
			encodeErrorMsg,
		)
	}
}
