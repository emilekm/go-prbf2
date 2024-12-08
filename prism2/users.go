package prism2

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

type Users struct {
	c *Client
}

func New(c *Client) *Users {
	return &Users{c}
}

func (u *Users) List(ctx context.Context) (UserList, error) {
	resp, err := u.c.Send(ctx, &Request{
		Message:         NewMessage(CommandGetUsers, nil),
		ExpectedSubject: SubjectGetUsers,
	})
	if err != nil {
		return nil, err
	}

	return usersList(resp.Message)
}

func (u *Users) Add(ctx context.Context, newUser AddUser) (UserList, error) {
	payload, err := Marshal(newUser)
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

func (u *Users) Change(ctx context.Context, user ChangeUser) (UserList, error) {
	payload, err := Marshal(user)
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

func (u *Users) Delete(ctx context.Context, name string) (UserList, error) {
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
