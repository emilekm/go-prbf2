package users

import "github.com/emilekm/go-prbf2/prism"

const (
	SubjectGetUsers prism.Subject = "getusers"

	CommandGetUsers   prism.Command = "getusers"
	CommandAddUser    prism.Command = "adduser"
	CommandChangeUser prism.Command = "changeuser"
	CommandDeleteUser prism.Command = "deleteuser"
)

// User returned with `getusers` message
type User struct {
	Name  string
	Power int
}

// List of users returned with `getusers` message
type UserList []User

func (u *UserList) UnmarshalMessage(content []byte) error {
	users, err := prism.UnmarshalMultipartBody[User](content)
	if err != nil {
		return err
	}

	*u = users
	return nil
}

type AddUser struct {
	Name     string
	Password string
	Power    int
}

type ChangeUser struct {
	Name        string
	NewName     string
	NewPassword string
	NewPower    int
}
