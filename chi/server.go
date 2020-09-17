package chi

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Handler struct {
	UserHandler
	AuthHandler
	r *chi.Mux
}

func NewHandler() *Handler {
	h := Handler{}
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Route("/auth", h.AuthHandler.getRoutes())

	r.Group(func(r chi.Router) {
		r.Use(h.Authenticate)
		r.Route("/users", h.UserHandler.getRoutes())
	})

	h.r = r
	return &h
}

func (h *Handler) ListenAndServe(addr string) {
	fmt.Printf("running server at address %s\n", addr)
	http.ListenAndServe(addr, h.r)
}

type Response struct {
	Status int `json:"status"`
	Result interface{} `json:"result,omitempty"`
	Message string `json:"message,omitempty"`
}
