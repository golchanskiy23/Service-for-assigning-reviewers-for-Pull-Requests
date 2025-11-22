package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/util"
)

const (
	teamNameField = "team_name"
	ERROR         = "error"
)

type TeamAddResponse struct {
	Team entity.Team `json:"team"`
}

func validateTeamAddRequest(team *entity.Team) error {
	if strings.TrimSpace(team.TeamName) == "" {
		return errors.New("team_name is required")
	}
	return nil
}

func validateTeamName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("team_name is required")
	}
	return nil
}

func (s *Services) TeamAddHandler(w http.ResponseWriter, r *http.Request) {
	var req entity.Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.Log.Warn("failed to decode team add request", ERROR, err)
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			"invalid json",
		)

		return
	}

	if err := validateTeamAddRequest(&req); err != nil {
		s.Log.Warn("invalid team add request", ERROR, err)
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			err.Error(),
		)

		return
	}

	ctx := r.Context()

	team, err := s.TeamService.AddTeam(ctx, &req)
	if err != nil {
		if errors.Is(err, entity.ErrTeamExists) {
			s.Log.Info("attempt to create team with existing name",
				teamNameField,
				req.TeamName)

			util.SendError(
				w,
				http.StatusBadRequest,
				entity.CodeTeamExists,
				"team_name already exists",
			)

			return
		}

		s.Log.Error("failed to add team",
			errFieldName, err,
			teamNameField,
			req.TeamName)

		util.SendError(
			w,
			http.StatusInternalServerError,
			entity.CodeInternalError,
			"internal server error",
		)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(TeamAddResponse{Team: *team})
	if err != nil {
		s.Log.Error("failed to encode team add response", ERROR, err)
		util.SendError(
			w,
			http.StatusInternalServerError,
			entity.CodeInternalError,
			"failed to encode response",
		)
	}
}

func (s *Services) TeamGetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get(teamNameField)
	if err := validateTeamName(name); err != nil {
		s.Log.Warn("invalid team get request", ERROR, err)
		util.SendError(
			w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			err.Error(),
		)

		return
	}

	ctx := r.Context()

	team, err := s.TeamService.GetTeam(ctx, name)
	if err != nil {
		s.Log.Warn("team not found", "team_name", name)
		util.SendError(
			w,
			http.StatusNotFound,
			entity.CodeNotFound,
			"team not found",
		)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(team); err != nil {
		s.Log.Error("failed to encode team get response", ERROR, err)
		util.SendError(
			w,
			http.StatusInternalServerError,
			entity.CodeInternalError,
			"failed to encode response",
		)
	}
}
