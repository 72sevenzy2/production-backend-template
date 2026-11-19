package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/72sevenzy2/grouter/grouter"
	handlers "github.com/72sevenzy2/production-backend-template/handlers"
	repo "github.com/72sevenzy2/production-backend-template/repository"
)

func TestCreateUserWithAuth(t *testing.T) {
	r := grouter.NewGrouter()

	users := repo.NewUsers()
	r.Post("createUser", handlers.HandleAuthCreateUser(users))

	data := []byte(`
	{
		"username": "test",
		"id": 123,
		"password": "password12345"
	}
	`)

	req := httptest.NewRequest(http.MethodPost, "/createUser", bytes.NewReader(data))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
}
