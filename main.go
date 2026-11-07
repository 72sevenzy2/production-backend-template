package main

import (
	"net/http"
	"time"

	"github.com/72sevenzy2/grouter/grouter"
)

func main() {
	r := grouter.NewGrouter()

	r.Use(grouter.Recoverer())
	lim := grouter.NewLimiter(20, 1)
	r.Use(lim.RateLimiter()) // global rate limiting

	r.Get("health", func(w http.ResponseWriter, r *grouter.Request) {
		w.Write([]byte("server fully operational."))
	}, grouter.Logger(0))

	// API
	users := NewUsers()
	r.Post("api/users/createuser",
		HandleCreateUser(users),
		grouter.Timeout(time.Second*5), grouter.Logger(0)) // 5 second safe-gaurd and logger.

	r.Get("api/users/getuser/:id", HandleGetUsername(users), grouter.Logger(0))

	secretToken := "abcdefg" // super secret bearer auth token
	r.Put("api/users/deleteuser/:id",
		HandleDeleteUser(users),
		grouter.BearerAuth(secretToken))

	http.ListenAndServe(":9000", r)
}
