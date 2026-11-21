package repository

import (
	"sync"
	"time"
)

// an mimic database for storing users, with O(n) lookup times to mimic an heavy db query.

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
}

func NewUser(id, following, followers int, username, password, sID string, expiresAt time.Duration) *User {
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

func (u *Users) GetUser(ID int) (*User, bool) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	for _, user := range u.UsersDB {
		if user.ID == ID {
			// returning a snapshot of User to avoid data races.
			copiedUser := user
			return &copiedUser, true
		}
	}
	return nil, false
}

func (u *Users) getName(ID int) (string, bool) {
	for _, user := range u.UsersDB {
		if user.ID == ID {
			return user.GetName(), true
		}
	}
	return "", false
}

func (u *Users) DeleteUser(ID int) {
	u.mu.Lock()
	defer u.mu.Unlock()

	for i, user := range u.UsersDB {
		if user.ID == ID {
			u.UsersDB = append(
				u.UsersDB[:i],
				u.UsersDB[i+1:]...,
			)
			return
		}
	}
}

type EditUserConfig struct {
	Username  string `json:"username"`
	ID        int    `json:"id"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
}

func (u *Users) EditUser(conf *EditUserConfig) bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, user := range u.UsersDB {
		if user.ID == conf.ID {
			user = User{
				ID:        conf.ID, // preseve user id
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

func (u *User) SetID(id int) {
	u.ID = id
}
func (u *User) SetUsername(name string) {
	u.Username = name
}

func (u *User) SetFollowers(n int) {
	u.Followers = n
}
func (u *User) SetFollowing(n int) {
	u.Following = n
}

// all the following getters implement the UserData interface.
func (u *User) GetName() string {
	return u.Username
}

func (u *User) GetID() int {
	return u.ID
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
