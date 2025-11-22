package handlers

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/util"
	"encoding/json"
	"net/http"
)

type TeamAddResponse struct {
	Team entity.Team `json:"team"`
}

func (s *Services) TeamAddHandler(w http.ResponseWriter, r *http.Request) {
	var req entity.Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, entity.CodeNotFound, "invalid json")
		return
	}
	if req.TeamName == "" {
		util.SendError(w, http.StatusBadRequest, entity.CodeNotFound, "team name is required")
		return
	}

	ctx := r.Context()
	team, err := s.TeamService.AddTeam(ctx, &req)
	if err != nil {
		if err.Error() == "TEAM_EXISTS" {
			util.SendError(w, http.StatusBadRequest, entity.CodeTeamExists, "team_name already exists")
			return
		}
		util.SendError(w, http.StatusInternalServerError, entity.CodeNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(TeamAddResponse{Team: *team})
}

func (s *Services) TeamGetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("team_name")
	if name == "" {
		util.SendError(w, http.StatusNotFound, entity.CodeNotFound, "team not found")
		return
	}

	ctx := r.Context()
	team, err := s.TeamService.GetTeam(ctx, name)
	if err != nil {
		util.SendError(w, http.StatusNotFound, entity.CodeNotFound, "team not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(team)
}
