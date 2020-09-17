package chi

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/quentin-fox/fox"
)

type UserHandler struct {
	UserService fox.UserService
	AuthService fox.AuthService
}

func (h *UserHandler) getRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/list", h.List)
		r.Get("/{id}", h.ListOne)
		r.Post("/create", h.Create)
		r.Post("/update", h.Update)
		r.Post("/verify", h.Verify)
		r.Post("/auth", h.Authenticate)
	}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserService.List()

	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}

	ok(w, users)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
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
	ok(w, user)
}

func (h *UserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var body struct{
		ID int `json:"id"`
		Password string `json:"password"`
	}
	if err := decode(r, &body); err != nil {
		fail(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.UserService.ListOne(body.ID)

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

func (h *UserHandler) ListOne(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.UserService.ListOne(id)

	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}

	ok(w, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user fox.User
	if err := decode(r, &user); err != nil {
		fail(w, err, http.StatusBadRequest)
		return
	}

	if user.ID == 0 {
		err := errors.New("user id cannot be 0")
		fail(w, err, http.StatusBadRequest)
		return
	}

	if err := h.UserService.Update(&user); err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}

	ok(w, user)
}

func (h *UserHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var user fox.User
	if err := decode(r, &user); err != nil {
		fail(w, err, http.StatusBadRequest)
		return
	}

	if user.ID == 0 {
		err := errors.New("user id cannot be 0")
		fail(w, err, http.StatusBadRequest)
	}

	dbUser, err := h.UserService.ListOne(user.ID)

	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}

	if dbUser.IsVerified {
		err = errors.New("user is already verified")
		fail(w, err, http.StatusBadRequest)
	}

	if err := h.UserService.Verify(user.ID); err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}

	ok(w, user.ID)
}
