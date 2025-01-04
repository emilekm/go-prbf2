package prism

import (
	"context"
)

const (
	SubjectGetUsers Subject = "getusers"

	CommandGetUsers   Command = "getusers"
	CommandAddUser    Command = "adduser"
	CommandChangeUser Command = "changeuser"
	CommandDeleteUser Command = "deleteuser"
)

// User returned with `getusers` message
type User struct {
	Name  string
	Power int
}

// List of users returned with `getusers` message
type UserList []User

func (u *UserList) UnmarshalMessage(content []byte) error {
	users, err := UnmarshalMultipartBody[User](content)
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

type usersService struct {
	c *Client
}

func (u *usersService) List(ctx context.Context) (UserList, error) {
	resp, err := u.c.Send(ctx, &Request{
		Message:         NewMessage(CommandGetUsers, nil),
		ExpectedSubject: SubjectGetUsers,
	})
	if err != nil {
		return nil, err
	}

	return usersList(resp.Message)
}

func (u *usersService) Add(ctx context.Context, newUser *AddUser) (UserList, error) {
	user := *newUser
	user.Password = stringHash(newUser.Password)
	payload, err := Marshal(&user)
	if err != nil {
		return nil, err
	}

	resp, err := u.c.Send(ctx, &Request{
		Message:         NewMessage(CommandAddUser, payload),
		ExpectedSubject: SubjectGetUsers,
	})
	if err != nil {
		return nil, err
	}

	return usersList(resp.Message)
}

func (u *usersService) Change(ctx context.Context, changedUser *ChangeUser) (UserList, error) {
	user := *changedUser
	if changedUser.NewPassword != "" {
		user.NewPassword = stringHash(changedUser.NewPassword)
	}
	payload, err := Marshal(&user)
	if err != nil {
		return nil, err
	}

	resp, err := u.c.Send(ctx, &Request{
		Message:         NewMessage(CommandChangeUser, payload),
		ExpectedSubject: SubjectGetUsers,
	})
	if err != nil {
		return nil, err
	}

	return usersList(resp.Message)
}

func (u *usersService) Delete(ctx context.Context, name string) (UserList, error) {
	resp, err := u.c.Send(ctx, &Request{
		Message:         NewMessage(CommandDeleteUser, []byte(name)),
		ExpectedSubject: SubjectGetUsers,
	})
	if err != nil {
		return nil, err
	}

	return usersList(resp.Message)
}

func usersList(rawMsg *Message) (UserList, error) {
	var users UserList
	if len(rawMsg.Body()) == 0 {
		return users, nil
	}

	err := users.UnmarshalMessage(rawMsg.Body())
	if err != nil {
		return nil, err
	}

	return users, nil
}
