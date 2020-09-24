package chi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

const host = "localhost:3000"

func makePostRequest(t *testing.T, route string, body map[string]interface{}, handler http.HandlerFunc) *http.Response {
	data, err := json.Marshal(body)

	if err != nil {
		t.Fatalf("could not marshall body to json for route %s; error: %v", route, err)
	}

	reader := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPost, host+route, reader)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		t.Fatalf("could not create request for route %s: %v", route, err)
	}

	rec := httptest.NewRecorder()
	handler(rec, req)

	return rec.Result()
}

func makeGetRequest(t *testing.T, route string, handler http.HandlerFunc, params *map[string]string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, host+route, nil)
	if err != nil {
		t.Fatalf("could not create request for route %s: %v", route, err)
	}

	if params != nil {
		ctx := chi.NewRouteContext()
		for k, v := range *params {
			ctx.URLParams.Add(k, v)
		}
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	}

	rec := httptest.NewRecorder()
	handler(rec, req)

	return rec.Result()
}

func decodeRequest(t *testing.T, res *http.Response, v interface{}) {
	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		t.Errorf("could not decode request body: %v", err)
	}
}

func testStatus(t *testing.T, status int, expected int) {
	if status != expected {
		t.Errorf("status should be %d; got %d", expected, status)
	}
}
