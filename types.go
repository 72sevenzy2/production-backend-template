package main

type User struct {
	ID        int
	Username  string
	Followers int
	Following int
}

type Users struct {
	UsersDB []User
}

func NewUsers() *Users {
	return &Users{
		UsersDB: []User{},
	}
}

func (u Users) GetName(ID int) (string, bool) {
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
		}
	}
}

type EditUserConfig struct {
	Username  string
	Followers int
	Following int
}

func (u *Users) EditUser(ID int, conf *EditUserConfig) bool {
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

func (u *Users) AddUser(name string, ID, followers, following int) *User {
	user := User{
		ID:        ID,
		Username:  name,
		Followers: followers,
		Following: following,
	}

	u.UsersDB = append(u.UsersDB, user)

	return &user
}
