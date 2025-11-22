package handlers

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/util"
	"encoding/json"
	"net/http"
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
		util.SendError(w, http.StatusBadRequest, entity.CodeNotFound, "invalid json")
		return
	}

	ctx := r.Context()
	pr, _, err := s.PRService.CreatePR(ctx, req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		if err.Error() == "PR_EXISTS" {
			util.SendError(w, http.StatusConflict, entity.CodePRExists, "PR id already exists")
		} else if err.Error() == "NOT_FOUND" {
			util.SendError(w, http.StatusNotFound, entity.CodeNotFound, "author/team not found")
		} else {
			util.SendError(w, http.StatusInternalServerError, entity.CodeNotFound, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(PRCreateResponse{PR: *pr})
}

func (s *Services) PRMergeHandler(w http.ResponseWriter, r *http.Request) {
	var req PRMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, entity.CodeNotFound, "invalid json")
		return
	}

	ctx := r.Context()
	pr, err := s.PRService.MergePR(ctx, req.PullRequestID)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			util.SendError(w, http.StatusNotFound, entity.CodeNotFound, "PR not found")
		} else {
			util.SendError(w, http.StatusInternalServerError, entity.CodeNotFound, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PRMergeResponse{PR: *pr})
}

func (s *Services) PRReassignHandler(w http.ResponseWriter, r *http.Request) {
	var req PRReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, entity.CodeNotFound, "invalid json")
		return
	}

	ctx := r.Context()
	pr, replacedBy, err := s.PRService.ReassignReviewer(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			util.SendError(w, http.StatusNotFound, entity.CodeNotFound, "PR or user not found")
			return
		}
		if err.Error() == "PR_MERGED" {
			util.SendError(w, http.StatusConflict, entity.CodePRMerged, "cannot reassign on merged PR")
			return
		}
		if err.Error() == "NOT_ASSIGNED" {
			util.SendError(w, http.StatusConflict, entity.CodeNotAssigned, "reviewer is not assigned to this PR")
			return
		}
		if err.Error() == "NO_CANDIDATE" {
			util.SendError(w, http.StatusConflict, entity.CodeNoCandidate, "no active replacement candidate in team")
			return
		}
		util.SendError(w, http.StatusInternalServerError, entity.CodeNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PRReassignResponse{
		PR:         *pr,
		ReplacedBy: replacedBy,
	})
}
