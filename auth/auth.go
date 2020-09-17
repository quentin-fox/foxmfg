package auth

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/quentin-fox/fox"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{
	privateKey string // filename
	publicKey string // filename
	GetTime func() time.Time
}

func NewAuthService(privateKey string, publicKey string) *AuthService {
	return &AuthService{
		privateKey: privateKey,
		publicKey: publicKey,
		GetTime: time.Now, // injected to ease testing 
	}
}

func (s *AuthService) GenerateHash(password string) (hash string, err error) {
	if len(password) == 0 {
		return "", errors.New("password cannot have zero length")
	}
	bhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bhash), nil
}

func (s *AuthService) ValidatePassword(hash string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	success := err == nil
	return success, nil
}

func (s *AuthService) IssueJWT(u fox.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": u.ID,
		"name": u.FirstName + u.LastName,
		"iat": s.GetTime().Unix(),
		"nbf": s.GetTime().Unix(),
		"exp": s.GetTime().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	privateKeyB, err := ioutil.ReadFile(s.privateKey)
	if err != nil {
		return "", err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyB)
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}

func (s *AuthService) VerifyJWT(tokenStr string) (bool, error) {
	isVerified := false

	publicKeyB, err := ioutil.ReadFile(s.publicKey)

	if err != nil {
		return isVerified, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyB)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil 
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		notExpired := claims.VerifyExpiresAt(s.GetTime().Unix(), true)
		notBefore := claims.VerifyNotBefore(s.GetTime().Unix(), true)

		isVerified = notExpired && notBefore
	}

	return isVerified, nil
}
