package benchmarks

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

func BenchmarkGetUser(b *testing.B) {
	handler := grouter.NewGrouter()

	u := repo.NewUsers()
	cache := repo.NewCache[string](func(s string) (string, bool) {
		id, err := strconv.Atoi(s)
		if err != nil {
			return "", false
		}
		val, ok := u.GetUser(id)
		return val.GetName(), ok

	}, utils.NewSingleFlight[string, string](), time.Minute*5)
	constructedUser := repo.NewUser(10, 10, 10, "testName", "", "", 0)
	u.AddUser(*constructedUser)
	handler.Get("api/getuser/:id", handlers.HandleGetUsername(u, cache, utils.NewSingleFlight[string, string]()))

	b.ReportAllocs()
	for b.Loop() {

		req := httptest.NewRequest(http.MethodGet, "/api/getuser/10", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			b.Fatalf("got status %d, and expected %d", rr.Code, http.StatusOK)
			return
		}
	}
}
