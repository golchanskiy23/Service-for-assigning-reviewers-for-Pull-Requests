package server

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router) {
	r.Route("/team", func(r chi.Router) {
		r.Post("/add", TeamAddHandler)
		r.Get("/get", TeamGetHandler)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", UserSetIsActiveHandler)
		r.Get("/getReview", UserGetReviewHandler)
	})

	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", PRCreateHandler)
		r.Post("/merge", PRMergeHandler)
		r.Post("/reassign", PRReassignHandler)
	})

}
