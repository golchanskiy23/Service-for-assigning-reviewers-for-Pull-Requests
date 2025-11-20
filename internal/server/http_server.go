package server

import (
	handlers "Service-for-assigning-reviewers-for-Pull-Requests/internal/handlers"
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/service"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type HTTPServer struct {
	addr string
	r    *chi.Mux
}

func NewHTTPServer(addr string) *HTTPServer {
	execution := &handlers.ServiceExecution{
		TeamService: &service.TeamService{},
		UserService: &service.UserService{},
		PrService:   &service.PRService{},
	}
	r := chi.NewRouter()
	RegisterRoutes(execution, r)
	return &HTTPServer{
		addr: addr,
		r:    r,
	}
}

func (s *HTTPServer) Run() {
	log.Println("Starting server on", s.addr)
	if err := http.ListenAndServe(s.addr, s.r); err != nil {
		log.Fatal(err)
	}
}
