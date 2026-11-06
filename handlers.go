package main

import (
	"time"
)

func HandleGetUsername(u *Users, ID string) string {
	time.Sleep(time.Second * 1) // simulating a heavy db query

	ok, val := u.GetName(ID)
	if !ok {
		return ""
	}

	return val
}
