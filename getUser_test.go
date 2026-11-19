package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/72sevenzy2/grouter/grouter"
	handlers "github.com/72sevenzy2/production-backend-template/handlers"
	utils "github.com/72sevenzy2/production-backend-template/internal"
	repo "github.com/72sevenzy2/production-backend-template/repository"
)

func TestGetUser(t *testing.T) {
	r := grouter.NewGrouter()

	users := repo.NewUsers()
	user := repo.NewUser(123, 0, 0, "test", "", "", 0)
	users.AddUser(*user)
	cache := repo.NewCache[string](func(s string) (string, bool) {
		id, err := strconv.Atoi(s)
		if err != nil {
			return "", false
		}
		val, ok := users.GetUser(id)
		return val.GetName(), ok
	}, utils.NewSingleFlight[string, string](), time.Minute*5)
	r.Get("GetUser/:id", handlers.HandleGetUsername(users, cache, utils.NewSingleFlight[string, string]()))

	req := httptest.NewRequest(http.MethodGet, "/GetUser/123", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}
