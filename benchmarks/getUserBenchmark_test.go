package benchmarks

import (
	"bytes"
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

func BenchmarkGetUser(b *testing.B) {
	handler := grouter.NewGrouter()

	u := repo.NewUsers()
	cache := repo.NewCache[string](func(s string) (string, bool) {
		id, err := strconv.Atoi(s)
		if err != nil {
			return "", false
		}
		val, ok := u.GetName(id)
		return val, ok

	}, utils.NewSingleFlight[string, string](), time.Minute*5)
	u.AddUser("testName", 10, 10, 10)
	handler.Post("api/getuser/:id", handlers.HandleGetUsername(u, cache, utils.NewSingleFlight[string, string]()))

	b.ReportAllocs()
	for b.Loop() {
		data := []byte(`
		{
			"id": 10,
			"followers": 31,
			"following": 31,
			"username": "test"
		}
		`)

		req := httptest.NewRequest(http.MethodPost, "/api/createuser/12", bytes.NewReader(data))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			b.Fatalf("got status %d, and expected %d", rr.Code, http.StatusCreated)
			return
		}
	}
}
