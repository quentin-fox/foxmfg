package sqlx

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/quentin-fox/fox"
)

type UserService struct {
	DB *sqlx.DB
}

func (s *UserService) GetSchema() string {
	return `
CREATE TABLE users (
	"id" serial primary key,
	"firstName" varchar not null,
	"lastName" varchar not null,
	"email" varchar not null,
	"isVerified" bool not null default false,
	"hash" varchar not null
)
	`
}

func (s *UserService) Create(u *fox.User) error {
	q := `
INSERT INTO users ("firstName", "lastName", "email", "isVerified", "hash")
VALUES ($1, $2, $3, $4, $5)
RETURNING id
	`
	row := s.DB.QueryRowx(q, u.FirstName, u.LastName, u.Email, u.IsVerified, u.Hash)
	err := row.Scan(&u.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(u)
	return nil
}

func (s *UserService) Update(u *fox.User) error {
	if u.ID == 0 {
		return errors.New("User ID cannot be 0")
	}

	q := `
UPDATE users
SET "firstName" = $1, "lastName" = $2, email = $3, "isVerified" = $4
WHERE id = $5
	`

	_, err := s.DB.Exec(q, u.FirstName, u.LastName, u.Email, u.IsVerified, u.ID)

	return err
}

func (s *UserService) List() ([]fox.User, error) {
	q := "SELECT * FROM users"
	var users []fox.User
	err := s.DB.Select(&users, q)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return users, nil
}

func (s *UserService) ListOne(id int) (fox.User, error) {
	q := `
SELECT "id", "firstName", "lastName", "email", "isVerified", "hash" FROM users
WHERE users.id = $1
	`
	var user fox.User
	err := s.DB.Get(&user, q, id)
	if err != nil {
		fmt.Println(err)
		return user, err
	}

	return user, nil
}

func (s *UserService) Verify(id int) error {
	q := `
UPDATE users
SET "isVerified" = true
WHERE id = $1 AND "isVerified" = false
`
	_, err := s.DB.Exec(q, id)
	return err
}
