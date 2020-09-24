package mock

import "github.com/quentin-fox/fox"

type AuthService struct {
	GenerateHashFn func(password string) (hash string, err error)
	GenerateHashInvoked bool

	ValidatePasswordFn func(hash string, password string) (bool, error)
	ValidatePasswordInvoked bool

	IssueJWTFn func(u fox.User) (string, error)
	IssueJWTInvoked bool

	VerifyJWTFn func(tokenStr string) (bool, error)
	VerifyJWTInvoked bool
}

func (s *AuthService) GenerateHash(password string) (hash string, err error) {
	s.GenerateHashInvoked = true
	return s.GenerateHashFn(password)
}

func (s *AuthService) ValidatePassword(hash string, password string) (bool, error) {
	s.ValidatePasswordInvoked = true
	return s.ValidatePasswordFn(hash, password)
}

func (s *AuthService) IssueJWT(u fox.User) (string, error) {
	s.IssueJWTInvoked = true
	return s.IssueJWTFn(u)
}

func (s *AuthService) VerifyJWT(tokenStr string) (bool, error) {
	s.VerifyJWTInvoked = true
	return s.VerifyJWTFn(tokenStr)
}
