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

type PRReassignRequest struct {
	PullRequestID int `json:"pull_request_id"`
	OldReviewerID int `json:"reviewerId"`
}

type PRMergeRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

func (service *Services) PRCreateHandler(w http.ResponseWriter, r *http.Request) {
	var req PRCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, entity.CodeNotFound, "invalid json")
		return
	}

	// pull_request, code
	pr, code := service.PRService.CreatePR(req.Title, req.Author)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if code == entity.CodeNotFound {
		util.SendError(w, http.StatusNotFound, entity.CodeNotFound, "author/team not found")
	} else if code == entity.CodePRExists {
		util.SendError(w, http.StatusConflict, entity.CodeNotFound, "PR id already exists")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pr)
}

func (service *Services) PRMergeHandler(w http.ResponseWriter, r *http.Request) {
	var req PRMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, entity.CodeNotFound, "invalid json")
		return
	}

	pr, err := service.PRService.MergePR(req.PullRequestID)
	if err != nil {
		util.SendError(w, http.StatusNotFound, entity.CodeNotFound, "pr not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pr)
}

func (h *PullRequestHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req PRReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, entity.CodeNotFound, "invalid json")
		return
	}

	// 1. Получить PR
	pr, err := h.PRRepo.GetByID(req.PullRequestID)
	if err != nil || pr == nil {
		writeJSON(w, http.StatusNotFound, errorResp("NOT_FOUND", "pull request not found"))
		return
	}

	// 2. Проверить статус
	if pr.Status == "MERGED" {
		writeJSON(w, http.StatusConflict, errorResp("PR_MERGED", "cannot reassign on merged PR"))
		return
	}

	// 3. Проверить, что old_user_id был ревьювером
	index := -1
	for i, rid := range pr.AssignedReviewers {
		if rid == req.OldUserID {
			index = i
			break
		}
	}
	if index == -1 {
		writeJSON(w, http.StatusConflict, errorResp("NOT_ASSIGNED", "reviewer is not assigned to this PR"))
		return
	}

	// 4. Найти пользователя и его команду
	oldUser, err := h.UserRepo.GetByID(req.OldUserID)
	if err != nil || oldUser == nil {
		writeJSON(w, http.StatusNotFound, errorResp("NOT_FOUND", "old reviewer not found"))
		return
	}

	// 5. Найти замену в его команде
	candidates, err := h.UserRepo.GetTeamActiveUsers(oldUser.TeamID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResp("INTERNAL", "internal error"))
		return
	}

	var replacement *User
	for _, u := range candidates {
		if u.ID != req.OldUserID && u.Active {
			replacement = u
			break
		}
	}

	if replacement == nil {
		writeJSON(w, http.StatusConflict, errorResp("NO_CANDIDATE", "no active replacement candidate in team"))
		return
	}

	// 6. Заменить ревьювера
	pr.AssignedReviewers[index] = replacement.ID

	// 7. Сохранить
	if err := h.PRRepo.Save(pr); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResp("INTERNAL", "failed to save PR"))
		return
	}

	// 8. Ответ
	resp := ReassignResponse{
		PR:         *pr,
		ReplacedBy: replacement.ID,
	}
	writeJSON(w, http.StatusOK, resp)
}
