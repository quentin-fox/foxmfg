package chi_test

import (
	"net/http"
	"testing"

	"github.com/quentin-fox/fox"
	"github.com/quentin-fox/fox/chi"
	"github.com/quentin-fox/fox/mock"
)

func TestSignup(t *testing.T) {
	h := chi.AuthHandler{}
	us := mock.UserService{}
	as := mock.AuthService{}

	var userHash string
	us.CreateFn = func(u *fox.User) error {
		userHash = u.Hash
		u.ID = 1
		return nil
	}

	as.GenerateHashFn = func(password string) (hash string, err error) {
		return "hashstring", nil
	}

	h.UserService = &us
	h.AuthService = &as

	body := map[string]interface{}{
		"firstName":  "Quentin",
		"lastName":   "Fox",
		"email":      "qfox@test.ca",
		"password":   "testtest",
		"isVerified": false,
	}

	res := makePostRequest(t, "/auth/signup", body, h.Signup)

	if res.Body == nil {
		t.Fatal("response had no body")
	}

	var response struct {
		Status int
		Result struct {
			ID int
		}
	}

	decodeRequest(t, res, &response)
	testStatus(t, response.Status, 200)

	if response.Result.ID != 1 {
		t.Errorf("created user id should be 1; got %d", response.Result)
	}

	if userHash != "hashstring" {
		t.Errorf(`created user hash should be "hashstring"; got "%s"`, userHash)
	}

	if !us.CreateInvoked {
		t.Error("create function not invoked")
	}

	if !as.GenerateHashInvoked {
		t.Error("generate hash funtion not invoked")
	}
}

func TestLogin(t *testing.T) {
	tt := []struct {
		name            string
		hash            string
		isValidPassword bool
		passwordErr     error
		tokenStr        string
		tokenErr        error
		status          int
		result          interface{}
		message         string
	}{
		{
			name:            "logs in successfully",
			hash:            "passhash",
			isValidPassword: true,
			tokenStr:        "jwtToken",
			status:          http.StatusOK,
			result: struct {
				Token string
			}{
				Token: "jwtToken",
			},
		},
		// {
		// 	name:            "incorrect password",
		// 	hash:            "passhash",
		// 	isValidPassword: false,
		// 	status:          http.StatusUnauthorized,
		// 	message:         "incorrect password",
		// },
	}

	for _, tc := range tt {

		t.Run(tc.name, func (t *testing.T) {
			h := chi.AuthHandler{}
			us := mock.UserService{}
			as := mock.AuthService{}

			us.ListWithHashFn = func(id int) (fox.User, error) {
				return fox.User{
					ID:   id,
					Hash: tc.hash,
				}, nil
			}

			as.ValidatePasswordFn = func(hash string, password string) (bool, error) {
				return tc.isValidPassword, tc.passwordErr
			}

			as.IssueJWTFn = func(u fox.User) (string, error) {
				return tc.tokenStr, tc.tokenErr
			}

			h.UserService = &us
			h.AuthService = &as

			body := map[string]interface{}{
				"id": 1,
				"password": "testtest",
			}

			res := makePostRequest(t, "/auth/login", body, h.Login)

			var response struct {
				Status int
				Result struct {
					Token string
				}
				Message string
			}

			decodeRequest(t, res, &response)
			testStatus(t, response.Status, tc.status)
		})
	}
}
