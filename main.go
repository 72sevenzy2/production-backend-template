package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/72sevenzy2/grouter/grouter"

	handlers "github.com/72sevenzy2/production-backend-template/handlers"

	utils "github.com/72sevenzy2/production-backend-template/internal"
	repo "github.com/72sevenzy2/production-backend-template/repository"
)

func main() {
	r := grouter.NewGrouter()

	healthLimiter := grouter.NewLimiter(50, 1)
	r.Get("health", func(w http.ResponseWriter, r *grouter.Request) {
		w.Write([]byte("server fully operational."))
	}, grouter.Recoverer(), grouter.Logger(0), healthLimiter.RateLimiter())

	// API

	users := repo.NewUsers()
	cache := repo.NewCache[string](func(s string) (string, bool) {
		id, err := strconv.Atoi(s)
		if err != nil {
			return "", false
		}
		val, ok := users.GetName(id)
		return val, ok
	}, utils.NewSingleFlight[string, string](), time.Minute*5)

	api := r.Group("api")

	api.Use(grouter.Recoverer())
	globalLimiter := grouter.NewLimiter(20, 1)
	r.Use(globalLimiter.RateLimiter())       // global rate limiting
	r.Use(grouter.Timeout(time.Second * 10)) // global 10 second timeout.
	api.Post("users/createuser",
		handlers.HandleCreateUser(users),
		grouter.Timeout(time.Second*5), grouter.Logger(0)) // 5 second safe-gaurd and logger.

	sf := utils.NewSingleFlight[string, string]()
	GetUserRateLimiter := grouter.NewLimiter(50, 1) // allow more requests for read-only route.
	api.Get("users/getuser/:id", handlers.HandleGetUsername(users, cache, sf), grouter.Logger(0), GetUserRateLimiter.RateLimiter())

	secretToken := "abcdefg" // super secret bearer auth token
	api.Del("users/deleteuser/:id",
		handlers.HandleDeleteUser(users),
		grouter.BearerAuth(secretToken), grouter.Timeout(time.Second*3))

	secretToken2 := "123456"
	api.Put("users/editUser/:id",
		handlers.HandleEditUser(users, cache),
		grouter.BearerAuth(secretToken2),
		grouter.Logger(0))

	// few authenticated endpoints (creating users and logging in).

	// authLimiter := grouter.NewLimiter(5, 1)
	//
	// auth := api.Group("auth") // this now becomes api/auth/..
	// auth.Post("login",.., authLimiter.RateLimiter())

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
