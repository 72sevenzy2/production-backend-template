package main

import (
	"encoding/json"
	"io"
	"log/slog"
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
		w.Write([]byte("fully operational."))
	}, grouter.Logger(0))

	users := NewUsers()
	r.Post("api/users/createuser", func(w http.ResponseWriter, r *grouter.Request) {
		limitedReader := io.LimitReader(r.Body, 1024*1024) // 1 MiB
		var u User
		if err := json.NewDecoder(limitedReader).Decode(&u); err != nil {
			grouter.DisplayInternalServerError(w)
		}

		user := users.AddUser(u.Username, u.ID, u.Followers, u.Following)
		slog.Info("POST_REQUEST", "CREATE_USER", "USER", user)

	}, grouter.Timeout(time.Second*5)) // 5 second safe-gaurd

	r.Get("api/users/getuser/:id", HandleGetUsername(users))

	http.ListenAndServe(":9000", r)
}
