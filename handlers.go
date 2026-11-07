package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/72sevenzy2/grouter/grouter"
)

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

		if err := json.NewEncoder(w).Encode(val); err != nil {
			grouter.DisplayInternalServerError(w)
		}
	}
}
