package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/72sevenzy2/grouter/grouter"
)

func HandleEditUser(u *Users, cache *Cache[string]) grouter.HandlerFunc {
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

		var UserEditConf EditUserConfig
		limitedReader := io.LimitReader(r.Body, 1024*1023) // 1 MiB
		if err := json.NewDecoder(limitedReader).Decode(&UserEditConf); err != nil {
			grouter.DisplayInternalServerError(w)
			return
		}

		w.WriteHeader(http.StatusOK)
		// todo: add cache invalidation.
		ok2 := u.EditUser(param2, &UserEditConf)
		_ = cache.Edit(param, UserEditConf.Username)
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
	}
}

func HandleCreateUser(u *Users) grouter.HandlerFunc {
	return func(w http.ResponseWriter, r *grouter.Request) {
		limitedReader := io.LimitReader(r.Body, 1024*1024) // 1 MiB
		var user User
		if err := json.NewDecoder(limitedReader).Decode(&user); err != nil {
			grouter.DisplayInternalServerError(w)
			return
		}

		w.WriteHeader(http.StatusCreated)
		u.AddUser(user.Username, user.ID, user.Followers, user.Following)
	}
}

func HandleGetUsername(u *Users, cache *Cache[string], sf *SingleFlight[string, string]) grouter.HandlerFunc {
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

		val, ok := cache.Get(param)
		if !ok {
			grouter.DisplayStatus(w, grouter.DisplayStatusConfig{
				Code:       http.StatusFailedDependency,
				Title:      fmt.Sprintf("status %d", http.StatusFailedDependency),
				Heading:    fmt.Sprintf("status %d", http.StatusFailedDependency),
				Descripion: "Sorry, your request failed as the user you were looking for was not found.",
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(val); err != nil {
			return
		}
	}
}

func HandleDeleteUser(u *Users) grouter.HandlerFunc {
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

		w.WriteHeader(http.StatusOK)
		u.DeleteUser(param2)
	}
}
