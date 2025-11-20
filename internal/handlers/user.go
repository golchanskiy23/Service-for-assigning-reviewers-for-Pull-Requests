package handlers

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/util"
	"encoding/json"
	"net/http"
)

type UserSetIsActiveRequest struct {
	UserID   int  `json:"userId"`
	IsActive bool `json:"isActive"`
}

type UserHandlers struct {
	//svc UserService
}

func (service *UserHandlers) UserSetIsActiveHandler(w http.ResponseWriter, r *http.Request) {
	var req UserSetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err := service.SetUserActive(req.UserID, req.IsActive)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (service *UserHandlers) UserGetReviewHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		util.RespondError(w, http.StatusBadRequest, "userId required")
		return
	}

	reviews, err := service.GetPRsAssignedTo(userIDStr)
	if err != nil {
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusOK, reviews)
}
