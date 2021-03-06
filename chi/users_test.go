package chi_test

import (
	"testing"

	"github.com/quentin-fox/fox"
	"github.com/quentin-fox/fox/chi"
	"github.com/quentin-fox/fox/mock"
)

func TestUserList(t *testing.T) {
	h := chi.UserHandler{}
	us := mock.UserService{}

	us.ListFn = func() ([]fox.User, error) {
		var list []fox.User
		list = append(list, fox.User{
			ID:         1,
			FirstName:  "Quentin",
			LastName:   "Fox",
			Email:      "qfox@test.ca",
			IsVerified: true,
		})
		return list, nil
	}

	h.UserService = &us

	var response struct {
		Status int
		Result []fox.User
	}

	res := makeGetRequest(t, "/users/list", h.List, nil)
	decodeRequest(t, res, &response)
	testStatus(t, response.Status, 200)

	if length := len(response.Result); length != 1 {
		t.Errorf("result length should be 1; got %d", length)
	}

	if response.Result[0].ID != 1 {
		t.Errorf("first user id should be 1; got %d", response.Result[0].ID)
	}

	if !us.ListInvoked {
		t.Error("List method was not invoked")
	}
}

func TestUserListOne(t *testing.T) {
	h := chi.UserHandler{}
	us := mock.UserService{}

	us.ListOneFn = func(id int) (fox.User, error) {
		return fox.User{
			ID:         id,
			FirstName:  "Quentin",
			LastName:   "Fox",
			Email:      "qfox@test.ca",
			IsVerified: true,
		}, nil
	}

	h.UserService = &us

	var response struct {
		Status int
		Result fox.User
	}

	params := map[string]string{
		"id": "1",
	}
	res := makeGetRequest(t, "/users/1", h.ListOne, &params)
	decodeRequest(t, res, &response)
	testStatus(t, response.Status, 200)

	if response.Result.ID != 1 {
		t.Errorf("user should have id 1; got %d", response.Result.ID)
	}
}
