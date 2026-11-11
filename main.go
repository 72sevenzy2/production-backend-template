package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/72sevenzy2/grouter/grouter"
)

func main() {
	r := grouter.NewGrouter()

	r.Use(grouter.Recoverer())
	globalLimiter := grouter.NewLimiter(20, 1)
	r.Use(globalLimiter.RateLimiter())       // global rate limiting
	r.Use(grouter.Timeout(time.Second * 10)) // global 10 second timeout.

	r.Get("health", func(w http.ResponseWriter, r *grouter.Request) {
		w.Write([]byte("server fully operational."))
	}, grouter.Logger(0))

	// API

	users := NewUsers()
	cache := NewCache[string](func(s string) (string, bool) {
		id, err := strconv.Atoi(s)
		if err != nil {
			return "", false
		}
		val, ok := users.GetName(id)
		return val, ok
	}, NewSingleFlight[string, string](), time.Minute*5)

	r.Post("api/users/createuser",
		HandleCreateUser(users),
		grouter.Timeout(time.Second*5), grouter.Logger(0)) // 5 second safe-gaurd and logger.

	sf := NewSingleFlight[string, string]()
	GetUserRateLimiter := grouter.NewLimiter(50, 1) // allow more requests for read-only route.
	r.Get("api/users/getuser/:id", HandleGetUsername(users, cache, sf), grouter.Logger(0), GetUserRateLimiter.RateLimiter())

	secretToken := "abcdefg" // super secret bearer auth token
	r.Del("api/users/deleteuser/:id",
		HandleDeleteUser(users),
		grouter.BearerAuth(secretToken), grouter.Timeout(time.Second*3))

	secretToken2 := "123456"
	r.Put("api/users/editUser/:id",
		HandleEditUser(users),
		grouter.BearerAuth(secretToken2),
		grouter.Logger(0))

	srv := &http.Server{
		Addr:              ":9000",
		IdleTimeout:       60 * time.Second,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
