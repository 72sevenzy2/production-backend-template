package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/72sevenzy2/grouter/grouter"
	"github.com/72sevenzy2/production-backend-template/handlers"
	repo "github.com/72sevenzy2/production-backend-template/repository"
)

func TestUserLogin(t *testing.T) {
	r := grouter.NewGrouter()

	users := repo.NewUsers()
	users.AddUser("test", "pass123", 10, 0, 0) // pre-initialisation for test user credentials
	r.Post("loginuser", handlers.HandleUserLogin(users))

	data := []byte(`
	{
		"id": 10,
		"password": "pass123"
	}
	`)

	req := httptest.NewRequest(http.MethodPost, "/loginuser", bytes.NewReader(data))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}
