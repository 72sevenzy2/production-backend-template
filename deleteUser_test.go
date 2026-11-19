package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/72sevenzy2/grouter/grouter"
	handlers "github.com/72sevenzy2/production-backend-template/handlers"
	repo "github.com/72sevenzy2/production-backend-template/repository"
)

func TestDeleteUser(t *testing.T) {
	r := grouter.NewGrouter()

	users := repo.NewUsers()
	user := repo.NewUser(123, 0, 0, "test", "", "", 0)
	users.AddUser(*user)
	r.Del("DeleteUser/:id", handlers.HandleDeleteUser(users))

	req := httptest.NewRequest(http.MethodDelete, "/DeleteUser/123", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}
