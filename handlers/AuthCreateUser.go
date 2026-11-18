package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/72sevenzy2/grouter/grouter"
	repo "github.com/72sevenzy2/production-backend-template/repository"

	utils "github.com/72sevenzy2/production-backend-template/internal"
)

// create users but with passwords (to further experiment grouters capabilites with auth).

func HandleAuthCreateUser(u *repo.Users) grouter.HandlerFunc {
	return func(w http.ResponseWriter, r *grouter.Request) {
		limit_body := io.LimitReader(r.Body, 1024*1204) // 1 MiB

		var regReq utils.RegisterRequest
		if err := json.NewDecoder(limit_body).Decode(&regReq); err != nil {
			grouter.DisplayInternalServerError(w)
			return
		}

		if ok := u.AddUser(regReq.Username, regReq.Password, regReq.ID, 0, 0); !ok {
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
