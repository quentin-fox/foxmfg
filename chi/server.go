package chi

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Handler struct {
	UserHandler
	r *chi.Mux
}

func NewHandler() *Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	h := Handler{r: r}
	h.r.Route("/users", h.UserHandler.getRoutes())
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

func NewRouteContext() *chi.Context {
	return chi.NewRouteContext()
}

var (
	RouteCtxKey = chi.RouteCtxKey
)

