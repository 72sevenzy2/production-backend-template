package repository

import (
	"sync"
)

// an mimic database for storing users, with O(n) lookup times to mimic an heavy db query.

type UserData interface {
	GetName() string
	GetHashedPassword() string
	GetFollowers() int
	GetFollowing() int
}

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
}

// GetName, GetFollowers, GetFollowing. and GetHashedPassword implement the UserData interface.
func (u *User) GetName() string {
	return u.Username
}

func (u *User) GetHashedPassword() string {
	return u.Password
}

func (u *User) GetFollowers() int {
	return u.Followers
}

func (u *User) GetFollowing() int {
	return u.Following
}

func NewUser(id, following, followers int, username, password string) *User {
	return &User{
		ID:        id,
		Username:  username,
		Password:  password,
		Followers: followers,
		Following: following,
	}
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

func (u *Users) GetUser(ID int) (UserData, bool) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	for v := range u.UsersDB {
		if u.UsersDB[v].ID == ID {
			// returning a snapshot of User to avoid data races.
			user := u.UsersDB[v]
			return &user, true
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

func (u *Users) AddUser(user User) bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, ok := u.getName(user.ID); ok {
		return false
	}

	u.UsersDB = append(u.UsersDB, user)

	return true
}
