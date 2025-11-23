package handlers

import (
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

type PRMergeRequest struct {
	PRID int `json:"prId"`
}

func (service *Services) PRCreateHandler(w http.ResponseWriter, r *http.Request) {
	/*var req PRCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	pr, err := service.PRService.CreatePR(req.Title, req.Author)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusCreated, pr)*/
}

func (service *Services) PRMergeHandler(w http.ResponseWriter, r *http.Request) {
	/*var req PRMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	pr, err := service.PRService.MergePR(req.PRID)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusOK, pr)*/
}

func (service *Services) PRReassignHandler(w http.ResponseWriter, r *http.Request) {
	/*var req PRReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	pr, err := service.PRService.ReassignReviewer(req.PRID, req.ReviewerID)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusOK, pr)*/
}
