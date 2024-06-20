package prism

import "context"

func (c *Client) ListUsers(ctx context.Context) (Users, error) {
	rawMsg, err := c.Command(ctx, CommandGetUsers, nil, SubjectGetUsers)
	if err != nil {
		return nil, err
	}

	return usersList(rawMsg)
}

func (c *Client) AddUser(ctx context.Context, newUser AddUser) (Users, error) {
	rawMsg, err := c.Command(ctx, CommandAddUser, &newUser, SubjectGetUsers)
	if err != nil {
		return nil, err
	}

	return usersList(rawMsg)
}

func (c *Client) ChangeUser(ctx context.Context, user ChangeUser) (Users, error) {
	rawMsg, err := c.Command(ctx, CommandChangeUser, &user, SubjectGetUsers)
	if err != nil {
		return nil, err
	}

	return usersList(rawMsg)
}

func (c *Client) DeleteUser(ctx context.Context, name string) (Users, error) {
	rawMsg, err := c.Command(ctx, CommandDeleteUser, name, SubjectGetUsers)
	if err != nil {
		return nil, err
	}

	return usersList(rawMsg)
}

func usersList(rawMsg *RawMessage) (Users, error) {
	var users Users
	if len(rawMsg.Body()) == 0 {
		return users, nil
	}

	err := UnmarshalMessage(rawMsg.Body(), &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
