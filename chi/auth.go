package chi

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi"
	"github.com/quentin-fox/fox"
)

type AuthHandler struct {
	UserService fox.UserService
	AuthService fox.AuthService
}

func (h *AuthHandler) getRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/signup", h.Signup)
		r.Post("/login", h.Login)
	}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var user struct{
		fox.User
		Password string `json:"password"`
	}
	if err := decode(r, &user); err != nil {
		fail(w, err, http.StatusBadRequest)
		return
	}

	hash, err := h.AuthService.GenerateHash(user.Password)

	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}

	user.Hash = hash
	h.UserService.Create(&user.User)
	response := struct{
		ID int `json:"id"`
	}{
		ID: user.ID,
	}

	ok(w, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct{
		ID int `json:"id"`
		Password string `json:"password"`
	}
	if err := decode(r, &body); err != nil {
		fail(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.UserService.ListWithHash(body.ID)

	if err != nil {
		fail(w, err, http.StatusUnauthorized)
		return
	}

	isValid, err := h.AuthService.ValidatePassword(user.Hash, body.Password)

	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}

	if !isValid {
		err = errors.New("incorrect password")
		fail(w, err, http.StatusUnauthorized)
		return
	}

	tokenStr, err := h.AuthService.IssueJWT(user)

	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}

	result := struct{
		Token string `json:"token"`
	}{
		Token: tokenStr,
	}

	ok(w, result)
}

func (h *AuthHandler) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeaders := r.Header["Authorization"]
		fmt.Println(authHeaders)

		if len(authHeaders) == 0 {
			err := errors.New("Authorization header is empty")
			fail(w, err, http.StatusUnauthorized)
			return
		}

		pattern := "^Bearer .*$"
		matched, err := regexp.MatchString(pattern, authHeaders[0])

		if !matched || err != nil {
			err := errors.New("Authorization header in incorrect format")
			fail(w, err, http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeaders[0], "Bearer ")
		isValid, err := h.AuthService.VerifyJWT(tokenStr)

		if !isValid {
			err := errors.New("Authorization token is not valid")
			fail(w, err, http.StatusUnauthorized)
			return
		}

		if err != nil {
			fail(w, err, http.StatusUnauthorized)
			return
		}


		next.ServeHTTP(w, r)
	})
}
