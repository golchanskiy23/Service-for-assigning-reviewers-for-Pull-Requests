package handlers

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/util"
	"encoding/json"
	"net/http"
)

type TeamAddRequest struct {
	Name string `json:"name"`
}

func (service *Services) TeamAddHandler(w http.ResponseWriter, r *http.Request) {
	var req TeamAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Name == "" {
		util.RespondError(w, http.StatusBadRequest, "team name is required")
		return
	}

	team, err := service.TeamService.AddTeam(req.Name)
	if err != nil {
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusCreated, team)
}

func (service *Services) TeamGetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		util.RespondError(w, http.StatusBadRequest, "name is required")
		return
	}

	team, err := service.TeamService.GetTeam(name)
	if err != nil {
		util.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusOK, team)
}
