package handlers

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/util"
	"encoding/json"
	"net/http"
)

type UserSetIsActiveRequest struct {
	UserID   int  `json:"userId"`
	IsActive bool `json:"isActive"`
}

func (service *Services) UserSetIsActiveHandler(w http.ResponseWriter, r *http.Request) {
	var req UserSetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err := service.UserService.SetUserActive(req.UserID, req.IsActive)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (service *Services) UserGetReviewHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		util.RespondError(w, http.StatusBadRequest, "userId required")
		return
	}

	reviews, err := service.UserService.GetPRsAssignedTo(userIDStr)
	if err != nil {
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusOK, reviews)
}
