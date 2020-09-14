package sqlx

import (
	"fmt"
	"log"

	"github.com/quentin-fox/fox"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type StorageService struct {
	UserService
}

func NewStorageService(c *fox.Config) *StorageService {
	params := []interface{}{
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Database,
	}
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		params...,
	)

	db, err := sqlx.Connect("postgres", dataSourceName)

	if err != nil {
		log.Fatal(err)
	}

	s := StorageService{}
	s.DB = db
	return &s
}

func (s *StorageService) CreateTables() error {
	tx, err := s.DB.Beginx()
	if err != nil {
		return err
	}

	if err = s.applySchemas(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *StorageService) applySchemas() error {
	schemas := s.listSchemas()

	for _, schema := range schemas {
		if _, err := s.DB.Exec(schema); err != nil {
			return err
		}
	}

	return nil
}

func (s *StorageService) listSchemas() []string {
	schemas := make([]string, 1)
	schemas[0] = s.UserService.GetSchema()
	return schemas
}
