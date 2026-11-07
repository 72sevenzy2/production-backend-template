package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/72sevenzy2/grouter/grouter"
)

func HandleCreateUser(users *Users) grouter.HandlerFunc {
	return func(w http.ResponseWriter, r *grouter.Request) {
		limitedReader := io.LimitReader(r.Body, 1024*1024) // 1 MiB
		var u User
		if err := json.NewDecoder(limitedReader).Decode(&u); err != nil {
			grouter.DisplayInternalServerError(w)
		}

		w.WriteHeader(http.StatusCreated)
		users.AddUser(u.Username, u.ID, u.Followers, u.Following)
	}
}

func HandleGetUsername(u *Users) grouter.HandlerFunc {
	time.Sleep(time.Second * 1) // simulating a heavy db query
	return func(w http.ResponseWriter, r *grouter.Request) {
		param, ok := r.Get("id")
		if !ok {
			grouter.DisplayInternalServerError(w)
		}
		param2, err := strconv.Atoi(param)
		if err != nil {
			grouter.DisplayInternalServerError(w)
		}

		val, ok := u.GetName(param2)

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(val); err != nil {
			grouter.DisplayInternalServerError(w)
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
