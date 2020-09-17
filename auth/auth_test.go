package auth_test

import (
	"testing"
	"time"

	"github.com/quentin-fox/fox"
	"github.com/quentin-fox/fox/auth"
)

const (
	privateKeyPath = "../id_rsa"
	publicKeyPath = "../id_rsa.pub"
)

func TestGenerateHash(t *testing.T) {
	a := auth.AuthService{} // no private/public key required
	tt := []struct{
		name string
		password string
		err bool
	}{
		{
			name: "hashes normal password",
			password: "secure_pass123",
			err: false,

		},
		{
			name: "won't hash zero length password",
			password: "",
			err: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := a.GenerateHash(tc.password)
			if (err != nil) != tc.err {
				t.Errorf("error from generating hash was unexpected; %v", err)
			}
		})
	}
}

func TestHashingUnhashing(t *testing.T) {
	a := auth.AuthService{}

	password := "secure_pass123"
	hash, err := a.GenerateHash(password)
	if err != nil {
		t.Errorf("error when generating hash: %v", err)
	}

	valid, err := a.ValidatePassword(hash, password)

	if err != nil {
		t.Errorf("error when validating password: %v", err)
	}

	if !valid {
		t.Errorf("did not validate password when it should have")
	}

	valid2, err := a.ValidatePassword(hash, "wrong_password")

	if err != nil {
		t.Errorf("error when validating incorrect password: %v", err)
	}

	if valid2 {
		t.Errorf("validated password when it should not have")
	}
}

func TestIssueJWT(t *testing.T) {
	a := auth.NewAuthService(privateKeyPath, publicKeyPath) // testing occurs relative to package
	u := fox.User{
		ID: 1,
		FirstName: "Quentin",
		LastName: "Fox",
		Email: "testemail@test.test",
		IsVerified: true,
	}

	tokenStr, err := a.IssueJWT(u)
	t.Log(tokenStr)

	if err != nil {
		t.Errorf("unexpected error when issuing jwt: %v", err)

	}

	ok, err := a.VerifyJWT(tokenStr)

	if err != nil {
		t.Errorf("unexpected error when verifying jwt: %v", err)
	}

	if !ok {
		t.Error("did not verify jwt when it should have")
	}
}

func TestInvalidJWTs(t *testing.T) {
	u := fox.User{
		ID: 1,
		FirstName: "Quentin",
		LastName: "Fox",
		Email: "testemail@test.test",
		IsVerified: true,
	}
	tt := []struct{
		name string
		verifyTime time.Time
	}{
		{
			name: "will not verify tokens before nbf claim",
			verifyTime: time.Now().Add(time.Hour * -1),
		},
		{
			name: "will not verify tokens after exp claim",
			verifyTime: time.Now().Add(time.Hour * -25), // expires after 24h
		},
		
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			a := auth.NewAuthService(privateKeyPath, publicKeyPath)
			tokenStr, err := a.IssueJWT(u)
			if err != nil {
				t.Errorf("unexpected error when issuing jwt: %v", err)
			}

			a.GetTime = func() time.Time { return tc.verifyTime }

			isValid, err := a.VerifyJWT(tokenStr)

			if err != nil {
				t.Errorf("unexpected error when verifying jwt: %v", err)
			}

			if isValid {
				t.Error("found jwt to be valid when it should not have been")
			}

		})
	}
}
