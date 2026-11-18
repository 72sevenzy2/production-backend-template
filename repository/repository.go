package repository

import (
	"sync"
)

// an mimic database for storing users, with O(n) lookup times to mimic an heavy db query.

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
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

func (u *Users) GetUser(ID int) (*User, bool) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	for v := range u.UsersDB {
		if u.UsersDB[v].ID == ID {
			return &u.UsersDB[v], true
		}
	}
	return nil, false
}

func (u *Users) getName(ID int) (string, bool) {
	for v := range u.UsersDB {
		if u.UsersDB[v].ID == ID {
			return u.UsersDB[v].Username, true
		}
	}
	return "", false
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
	Username  string `json:"username"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
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

func (u *Users) AddUser(name, pass string, ID, followers, following int) bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, ok := u.getName(ID); ok {
		return false
	}

	user := User{
		ID:        ID,
		Password:  pass,
		Username:  name,
		Followers: followers,
		Following: following,
	}

	u.UsersDB = append(u.UsersDB, user)

	return true
}
