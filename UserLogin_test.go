package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/72sevenzy2/grouter/grouter"
	"github.com/72sevenzy2/production-backend-template/handlers"
	utils "github.com/72sevenzy2/production-backend-template/internal"
	repo "github.com/72sevenzy2/production-backend-template/repository"
)

func TestUserLogin(t *testing.T) {
	r := grouter.NewGrouter()

	users := repo.NewUsers()
	sessions := repo.NewSessions()
	hashed, err := utils.HashPassword("test123")
	if err != nil {
		t.Fatal(err)
	}

	u := repo.NewUser(10, 0, 0, "test", hashed, "", 0)

	users.AddUser(*u) // pre-initialisation for test user credentials
	r.Post("loginuser", handlers.HandleUserLogin(users, sessions))

	data := []byte(`
	{
		"id": 10,
		"password": "test123"
	}
	`)

	req := httptest.NewRequest(http.MethodPost, "/loginuser", bytes.NewReader(data))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}
