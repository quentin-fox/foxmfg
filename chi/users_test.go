package chi_test

import (
	"encoding/json"
	"testing"

	"github.com/quentin-fox/fox"
	"github.com/quentin-fox/fox/chi"
	"github.com/quentin-fox/fox/mock"
)

func TestUserCreation(t *testing.T) {
	h := chi.UserHandler{}
	us := mock.UserService{}

	us.CreateFn = func(u *fox.User) error {
		u.ID = 1
		return nil
	}

	h.UserService = &us
	
	body := map[string]interface{}{
		"firstName": "Quentin",
		"lastName": "Fox",
		"email": "qfox@test.ca",
		"isVerified": false,
	}

	res := makePostRequest(t, "/users/create", body, h.Create)

	var response struct {
		Status int
		Result fox.User
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Error("Could not decode response")
	}

	if response.Status != 200 {
		t.Errorf("status should be 200; got %d", response.Status)
	}

	if response.Result.ID != 1 {
		t.Errorf("user id should be 1; got %d", response.Result.ID)
	}

	if !us.CreateInvoked {
		t.Error("Create method was not invoked")
	}
}

func TestUserList(t *testing.T) {
	h := chi.UserHandler{}
	us := mock.UserService{}

	us.ListFn = func() ([]fox.User, error) {
		var list []fox.User
		list = append(list, fox.User{
			ID: 1,
			FirstName: "Quentin",
			LastName: "Fox",
			Email: "qfox@test.ca",
			IsVerified: true,
		})
		return list, nil
	}

	h.UserService = &us

	var response struct {
		Status int
		Result []fox.User
	}

	res := makeGetRequest(t, "/users/list", h.List)

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Error("could not decode response")
	}

	if response.Status != 200 {
		t.Errorf("status should be 200; got %d", response.Status)
	}

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
