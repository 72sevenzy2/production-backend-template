package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/72sevenzy2/grouter/grouter"
)

func BenchmarkCreateUser(b *testing.B) {
	handler := grouter.NewGrouter()

	u := NewUsers()
	handler.Post("api/createuser/:id", HandleCreateUser(u))

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
