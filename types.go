package main

type User struct {
	ID        string
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

func (u Users) GetName(ID string) (bool, string) {
	for v := range u.UsersDB {
		if u.UsersDB[v].ID == ID {
			return true, u.UsersDB[v].Username
		}
	}
	return false, ""
}

func (u *Users) AddUser(name, ID string, followers, following int) *User {
	user := User{
		ID:        ID,
		Username:  name,
		Followers: followers,
		Following: following,
	}

	u.UsersDB = append(u.UsersDB, user)

	return &user
}
