package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	utils "github.com/72sevenzy2/production-backend-template/internal"

	"github.com/72sevenzy2/grouter/grouter"
	handlers "github.com/72sevenzy2/production-backend-template/handlers"
	repo "github.com/72sevenzy2/production-backend-template/repository"
)

func TestEditUser(t *testing.T) {
	r := grouter.NewGrouter()

	users := repo.NewUsers()
	user := repo.NewUser(123, 0, 0, "test", "")
	users.AddUser(*user)

	cache := repo.NewCache[string](func(s string) (string, bool) {
		id, err := strconv.Atoi(s)
		if err != nil {
			return "", false
		}
		val, ok := users.GetUser(id)
		return val.GetName(), ok
	}, utils.NewSingleFlight[string, string](), time.Minute*5)

	r.Post("EditUser/:id", handlers.HandleEditUser(users, cache))

	data := []byte(`
	{
		"username": "test",
		"followers": 0,
		"following": 0,
		"id": 0
	}
	`)

	req := httptest.NewRequest(http.MethodPost, "/EditUser/123", bytes.NewReader(data))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
}
