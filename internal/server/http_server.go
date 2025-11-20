package server

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type HTTPServer struct {
	addr string
	r    *chi.Mux
}

/*
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
}*/

func StartServer(srv *http.Server) error {
	//log.Println("Starting server on", s.addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	return nil
}
