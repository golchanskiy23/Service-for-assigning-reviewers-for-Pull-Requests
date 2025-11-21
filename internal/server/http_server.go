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

func NewHTTPServer(addr string) *HTTPServer {
	r := chi.NewRouter()
	RegisterRoutes(r)
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
