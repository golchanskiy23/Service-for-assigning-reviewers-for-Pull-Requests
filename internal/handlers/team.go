package handlers

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/util"
	"encoding/json"
	"net/http"
)

type TeamAddRequest struct {
	Team entity.Team `json:"team"`
}

type UserGetReviewRequest struct {
	Query entity.TeamNameQuery `json:"query"`
}

func (service *Services) TeamAddHandler(w http.ResponseWriter, r *http.Request) {
	var req TeamAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, entity.CodeNotFound, "invalid json")
		return
	}
	if req.Team.TeamName == "" {
		util.SendError(w, http.StatusBadRequest, entity.CodeNotFound, "team name is required")
		return
	}

	team, err := service.TeamService.AddTeam(req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, entity.CodeTeamExists, "team_name already exists")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(team)
}

func (service *Services) TeamGetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("team_name")
	if name == "" {
		util.SendError(w, http.StatusNotFound, entity.CodeNotFound, "team_not_found")
		return
	}

	team, err := service.TeamService.GetTeam(name)
	if err != nil {
		util.SendError(w, http.StatusNotFound, entity.CodeNotFound, "team_not_found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(team)
}
