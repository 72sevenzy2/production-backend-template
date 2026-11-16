package benchmarks

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/72sevenzy2/grouter/grouter"
	handlers "github.com/72sevenzy2/production-backend-template/handlers"

	repo "github.com/72sevenzy2/production-backend-template/repository"
)

func BenchmarkDeleteUser(b *testing.B) {
	handler := grouter.NewGrouter()

	u := repo.NewUsers()

	u.AddUser("testName", 10, 10, 10)
	handler.Del("api/deleteuser/:id", handlers.HandleDeleteUser(u))

	b.ReportAllocs()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodDelete, "/api/deleteuser/10", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			b.Fatalf("got status %d, and expected %d", rr.Code, http.StatusOK)
			return
		}

		b.StopTimer()
		u.AddUser("test", 10, 10, 10)
		b.StartTimer()
	}
}
