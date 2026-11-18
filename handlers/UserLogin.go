package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/72sevenzy2/grouter/grouter"
	utils "github.com/72sevenzy2/production-backend-template/internal"
	repo "github.com/72sevenzy2/production-backend-template/repository"
)

func HandleUserLogin(u *repo.Users) grouter.HandlerFunc {
	return func(w http.ResponseWriter, r *grouter.Request) {
		limit_body := io.LimitReader(r.Body, 1024*1024) // 1 MiB

		var LoginReq utils.LoginRequest
		if err := json.NewDecoder(limit_body).Decode(&LoginReq); err != nil {
			grouter.DisplayInternalServerError(w)
			return
		}

		user, ok := u.GetUser(LoginReq.ID)
		if !ok {
			grouter.DisplayStatus(w, grouter.DisplayStatusConfig{
				Code:       http.StatusFailedDependency,
				Title:      fmt.Sprintf("status %d", http.StatusFailedDependency),
				Heading:    fmt.Sprintf("status %d", http.StatusFailedDependency),
				Descripion: "Sorry, your request failed as the user you were looking for was not found.",
			})
			return
		}

		if user.Password != LoginReq.Password {
			grouter.DisplayStatus(w, grouter.DisplayStatusConfig{
				Code:       http.StatusForbidden,
				Title:      fmt.Sprintf("status %d", http.StatusForbidden),
				Heading:    fmt.Sprintf("status %d", http.StatusForbidden),
				Descripion: "Sorry, the password you entered for the account was wrong.",
			})
			return
		}

		returning := struct {
			Username string `json:"username"`
			Id       int    `json:"id"`
		}{
			Username: user.Username,
			Id:       user.ID,
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(returning); err != nil {
			grouter.DisplayInternalServerError(w)
			return
		}
	}

}
