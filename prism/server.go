package prism

import "context"

// Using `serverdetailsalways` instead of `serverdetails`
// because the latter might not return all the fields.
func (c *Client) ServerDetails(ctx context.Context) (*ServerDetails, error) {
	rawMsg, err := c.Command(ctx, CommandServerDetailsAlways, nil, SubjectServerDetails)
	if err != nil {
		return nil, err
	}

	var serverDetails ServerDetails
	err = UnmarshalMessage(rawMsg.Body(), &serverDetails)
	if err != nil {
		return nil, err
	}

	return &serverDetails, nil
}

func (c *Client) ListPlayers(ctx context.Context) (Players, error) {
	rawMsg, err := c.Command(ctx, CommandListPlayers, nil, SubjectListPlayers)
	if err != nil {
		return nil, err
	}

	var players Players
	if len(rawMsg.Body()) == 0 {
		return players, nil
	}

	err = UnmarshalMessage(rawMsg.Body(), &players)
	if err != nil {
		return nil, err
	}

	return players, nil
}
