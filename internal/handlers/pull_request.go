package handlers

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/util"
	"encoding/json"
	"net/http"
)

type PRCreateRequest struct {
	Title  string `json:"title"`
	Author int    `json:"author"`
}

type PRReassignRequest struct {
	PRID       int `json:"prId"`
	ReviewerID int `json:"reviewerId"`
}

type PRHandlers struct {
	svc PRService
}

type PRMergeRequest struct {
	PRID int `json:"prId"`
}

func (service *PRHandlers) PRCreateHandler(w http.ResponseWriter, r *http.Request) {
	var req PRCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	pr, err := service.CreatePR(req.Title, req.Author)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusCreated, pr)
}

func (service *PRHandlers) PRMergeHandler(w http.ResponseWriter, r *http.Request) {
	var req PRMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	pr, err := service.MergePR(req.PRID)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusOK, pr)
}

func (service *PRHandlers) PRReassignHandler(w http.ResponseWriter, r *http.Request) {
	var req PRReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	pr, err := service.ReassignReviewer(req.PRID, req.ReviewerID)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusOK, pr)
}
