package main

import (
	"log"

	"github.com/quentin-fox/fox/auth"
	"github.com/quentin-fox/fox/chi"
	"github.com/quentin-fox/fox/config"
	"github.com/quentin-fox/fox/sqlx"
)

func main() {
	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	s := sqlx.NewStorageService(&c)
	a := auth.NewAuthService(&c)

	h := chi.NewHandler()
	h.AuthHandler.UserService = &s.UserService
	h.UserHandler.UserService = &s.UserService
	h.AuthService = a

	h.ListenAndServe(c.Addr)
}
