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
}

func (h *UserHandler) getRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/list", h.List)
		r.Get("/{id}", h.ListOne)
		r.Post("/create", h.Create)
		r.Post("/update", h.Update)
		r.Post("/verify", h.Verify)
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
	var user fox.User
	if err := decode(r, &user); err != nil {
		fail(w, err, http.StatusBadRequest)
		return
	}

	h.UserService.Create(&user)
	ok(w, user)
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
