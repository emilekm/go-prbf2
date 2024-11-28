package users

import (
	"context"

	"github.com/emilekm/go-prbf2/prism"
)

type Users struct {
	c *prism.Client
}

func New(c *prism.Client) *Users {
	return &Users{c}
}

func (u *Users) List(ctx context.Context) (UserList, error) {
	rawMsg, err := u.c.Command(ctx, CommandGetUsers, nil, SubjectGetUsers)
	if err != nil {
		return nil, err
	}

	return usersList(rawMsg)
}

func (u *Users) Add(ctx context.Context, newUser AddUser) (UserList, error) {
	rawMsg, err := u.c.Command(ctx, CommandAddUser, &newUser, SubjectGetUsers)
	if err != nil {
		return nil, err
	}

	return usersList(rawMsg)
}

func (u *Users) Change(ctx context.Context, user ChangeUser) (UserList, error) {
	rawMsg, err := u.c.Command(ctx, CommandChangeUser, &user, SubjectGetUsers)
	if err != nil {
		return nil, err
	}

	return usersList(rawMsg)
}

func (u *Users) Delete(ctx context.Context, name string) (UserList, error) {
	rawMsg, err := u.c.Command(ctx, CommandDeleteUser, &name, SubjectGetUsers)
	if err != nil {
		return nil, err
	}

	return usersList(rawMsg)
}

func usersList(rawMsg *prism.RawMessage) (UserList, error) {
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
