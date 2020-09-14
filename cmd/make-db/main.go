package main

import (
	"log"

	"github.com/quentin-fox/fox/config"
	"github.com/quentin-fox/fox/sqlx"
)

func main() {
	c, err := config.GetConfig()

	if err != nil {
		log.Fatal(err)
	}

	s := sqlx.NewStorageService(&c)

	if err := s.CreateTables(); err != nil {
		log.Fatal(err)
	}
}
