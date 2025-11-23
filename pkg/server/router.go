package server

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(h *handlers.Services, r *chi.Mux) {
	r.Route("/team", func(r chi.Router) {
		r.Post("/add", h.TeamAddHandler)
		r.Get("/get", h.TeamGetHandler)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", h.UserSetIsActiveHandler)
		r.Get("/getReview", h.UserGetReviewHandler)
		r.Post("/deactivate", h.UsersMassDeactivateHandler)
	})

	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", h.PRCreateHandler)
		r.Post("/merge", h.PRMergeHandler)
		r.Post("/reassign", h.PRReassignHandler)
	})
}
