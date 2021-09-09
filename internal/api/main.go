package api

import (
	"genesis/internal/api/handlers"
	"genesis/internal/server"
	"github.com/go-chi/chi/v5"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()
	s := server.Server{}
	s.NewServer()
	r.Get("/btc", handlers.GetBTC)
	r.Get("/user/create", s.CreateUser)
	r.Get("/user/auth", s.AuthenticateUser)
	return r
}
