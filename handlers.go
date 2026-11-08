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
		ok2 := u.EditUser(param2, &UserEditConf)
		if !ok2 {
			grouter.DisplayInternalServerError(w)
			return
		}
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

func HandleGetUsername(u *Users, cache *Cache[string], sf *SingleFlight[string]) grouter.HandlerFunc {
	return func(w http.ResponseWriter, r *grouter.Request) {
		param, ok := r.Get("id")
		if !ok {
			grouter.DisplayInternalServerError(w)
			return
		}
		param2, err := strconv.Atoi(param)
		if err != nil {
			grouter.DisplayInternalServerError(w)
			return
		}
		val, err := sf.Do(param, func() (string, error) {
			if val, ok := cache.Get(param); ok {
				return val, nil
			}

			time.Sleep(time.Second * 1) // simulating heavy db query.
			val, ok := u.GetName(param2)
			if !ok {
				grouter.DisplayInternalServerError(w)
				return "", ErrNonExistentUser
			}
			cache.Put(param, val)
			return val, nil
		})
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
