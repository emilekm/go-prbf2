package prism

import "context"

// Using `serverdetailsalways` instead of `serverdetails`
// because the latter might not return all the fields.
func (c *Client) ServerDetails(ctx context.Context) (*ServerDetails, error) {
	resp, err := c.Send(ctx, &Request{
		Message:         NewMessage(CommandServerDetailsAlways, nil),
		ExpectedSubject: SubjectServerDetails,
	})
	if err != nil {
		return nil, err
	}

	var serverDetails ServerDetails
	err = Unmarshal(resp.Message.Body(), &serverDetails)
	if err != nil {
		return nil, err
	}

	return &serverDetails, nil
}

func (c *Client) GameplayeDetails(ctx context.Context) (*GameplayDetails, error) {
	resp, err := c.Send(ctx, &Request{
		Message:         NewMessage(CommandGameplayDetails, nil),
		ExpectedSubject: SubjectGameplayDetails,
	})
	if err != nil {
		return nil, err
	}

	var gameplayDetails GameplayDetails
	err = Unmarshal(resp.Message.Body(), &gameplayDetails)
	if err != nil {
		return nil, err
	}

	return &gameplayDetails, nil
}

func (c *Client) ListPlayers(ctx context.Context) (Players, error) {
	resp, err := c.Send(ctx, &Request{
		Message:         NewMessage(CommandListPlayers, nil),
		ExpectedSubject: SubjectListPlayers,
	})
	if err != nil {
		return nil, err
	}

	var players Players
	if len(resp.Message.Body()) == 0 {
		return players, nil
	}

	err = Unmarshal(resp.Message.Body(), &players)
	if err != nil {
		return nil, err
	}

	return players, nil
}
