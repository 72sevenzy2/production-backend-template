package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/72sevenzy2/grouter/grouter"
	repo "github.com/72sevenzy2/production-backend-template/repository"
)

func HandleCreateUser(u *repo.Users) grouter.HandlerFunc {
	return func(w http.ResponseWriter, r *grouter.Request) {
		limitedReader := io.LimitReader(r.Body, 1024*1024) // 1 MiB
		var user repo.User
		if err := json.NewDecoder(limitedReader).Decode(&user); err != nil {
			grouter.DisplayInternalServerError(w)
			return
		}

		constructedUser := repo.NewUser(user.ID, user.Following, user.Followers, user.Username, user.Password, "", 0)

		if ok := u.AddUser(*constructedUser); !ok {
			grouter.DisplayStatus(w, grouter.DisplayStatusConfig{
				Code:       http.StatusInternalServerError,
				Title:      http.StatusText(http.StatusInternalServerError),
				Heading:    fmt.Sprintf("status %d", http.StatusInternalServerError),
				Descripion: "sorry, this user has already been registered!",
			})
			return
		}
		w.WriteHeader(http.StatusCreated)

	}
}
