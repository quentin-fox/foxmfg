package auth

import "golang.org/x/crypto/bcrypt"

type AuthService struct{}

func (s *AuthService) GenerateHash(password string) (hash string, err error) {
	bhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bhash), nil
}

func (s *AuthService) ValidatePassword(hash string, password string) (bool, error) {
	bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return true, nil
}
