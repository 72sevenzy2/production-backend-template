package main

import (
	"sync"
)

type User struct {
	ID        int
	Username  string
	Followers int
	Following int
}

type Users struct {
	mu      sync.RWMutex
	UsersDB []User
}

func NewUsers() *Users {
	return &Users{
		UsersDB: []User{},
	}
}

func (u *Users) GetName(ID int) (string, bool) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	for v := range u.UsersDB {
		if u.UsersDB[v].ID == ID {
			return u.UsersDB[v].Username, true
		}
	}
	return "", false
}

func (u *Users) DeleteUser(ID int) {
	for v := range u.UsersDB {
		if u.UsersDB[v].ID == ID {
			u.UsersDB = append(
				u.UsersDB[:v],
				u.UsersDB[v+1:]...,
			)
			return
		}
	}
}

type EditUserConfig struct {
	Username  string
	Followers int
	Following int
}

func (u *Users) EditUser(ID int, conf *EditUserConfig) bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	for v := range u.UsersDB {
		if u.UsersDB[v].ID == ID {
			u.UsersDB[v] = User{
				ID:        ID, // preseve user id
				Username:  conf.Username,
				Followers: conf.Followers,
				Following: conf.Following,
			}
			return true
		}
	}

	return false
}

func (u *Users) AddUser(name string, ID, followers, following int) bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	user := User{
		ID:        ID,
		Username:  name,
		Followers: followers,
		Following: following,
	}

	u.UsersDB = append(u.UsersDB, user)

	return true
}
