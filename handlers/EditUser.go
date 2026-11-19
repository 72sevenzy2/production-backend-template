package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/72sevenzy2/grouter/grouter"
	repo "github.com/72sevenzy2/production-backend-template/repository"
)

func HandleEditUser(u *repo.Users, cache *repo.Cache[string]) grouter.HandlerFunc {
	return func(w http.ResponseWriter, r *grouter.Request) {
		param, ok := r.Get("id")
		if !ok {
			grouter.DisplayStatus(w, grouter.DisplayStatusConfig{
				Code:       http.StatusFailedDependency,
				Title:      fmt.Sprintf("status %d", http.StatusFailedDependency),
				Heading:    fmt.Sprintf("status %d", http.StatusFailedDependency),
				Descripion: "Sorry, your request failed as an route parameter was not found, please try again.",
			})

			return
		}
		param2, err := strconv.Atoi(param)
		if err != nil {
			grouter.DisplayInternalServerError(w)
			return
		}

		var UserEditConf repo.EditUserConfig
		limitedReader := io.LimitReader(r.Body, 1024*1024) // 1 MiB
		if err := json.NewDecoder(limitedReader).Decode(&UserEditConf); err != nil {
			grouter.DisplayInternalServerError(w)
			return
		}
		UserEditConf.ID = param2

		ok2 := u.EditUser(&UserEditConf)
		if !ok2 {
			grouter.DisplayStatus(w, grouter.DisplayStatusConfig{
				Code:       http.StatusFailedDependency,
				Title:      fmt.Sprintf("status %d", http.StatusFailedDependency),
				Heading:    fmt.Sprintf("status %d", http.StatusFailedDependency),
				Descripion: "Sorry, your request failed as the user you were looking to edit for was not found.",
			})

			return
		}
		// invalidate cache.
		_ = cache.Edit(param, UserEditConf.Username)
		w.WriteHeader(http.StatusCreated)
	}
}
