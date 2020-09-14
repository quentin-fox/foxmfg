package main

import (
	"log"

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

	h := chi.NewHandler()
	h.UserService = &s.UserService

	h.ListenAndServe(c.Addr)
}
