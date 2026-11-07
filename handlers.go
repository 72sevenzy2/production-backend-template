package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/72sevenzy2/grouter/grouter"
)

func HandleEditUser(u *Users) grouter.HandlerFunc {
	return func(w http.ResponseWriter, r *grouter.Request) {
		param, ok := r.Get("id")
		if !ok {
			grouter.DisplayInternalServerError(w)
		}
		param2, err := strconv.Atoi(param)
		if err != nil {
			grouter.DisplayInternalServerError(w)
		}

		var UserEditConf EditUserConfig
		limitedReader := io.LimitReader(r.Body, 1024*1023) // 1 MiB
		if err := json.NewDecoder(limitedReader).Decode(&UserEditConf); err != nil {
			grouter.DisplayInternalServerError(w)
		}

		w.WriteHeader(http.StatusOK)
		ok2 := u.EditUser(param2, &UserEditConf)
		if !ok2 {
			grouter.DisplayInternalServerError(w)
		}
	}
}

func HandleCreateUser(u *Users) grouter.HandlerFunc {
	return func(w http.ResponseWriter, r *grouter.Request) {
		limitedReader := io.LimitReader(r.Body, 1024*1024) // 1 MiB
		var user User
		if err := json.NewDecoder(limitedReader).Decode(&user); err != nil {
			grouter.DisplayInternalServerError(w)
		}

		w.WriteHeader(http.StatusCreated)
		u.AddUser(user.Username, user.ID, user.Followers, user.Following)
	}
}

func HandleGetUsername(u *Users, cache *Cache[string]) grouter.HandlerFunc {
	return func(w http.ResponseWriter, r *grouter.Request) {
		param, ok := r.Get("id")
		if !ok {
			grouter.DisplayInternalServerError(w)
		}
		param2, err := strconv.Atoi(param)
		if err != nil {
			grouter.DisplayInternalServerError(w)
		}
		val, ok := cache.Get(param)
		// todo: fix cache stampede
		if !ok {
			time.Sleep(time.Second * 5) // simulating heavy db query
			val, ok := u.GetName(param2)
			if !ok {
				grouter.DisplayInternalServerError(w)
			}
			cache.Put(param, val)
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
			grouter.DisplayInternalServerError(w)
		}
		param2, err := strconv.Atoi(param)
		if err != nil {
			grouter.DisplayInternalServerError(w)
		}

		w.WriteHeader(http.StatusOK)
		u.DeleteUser(param2)
	}
}
