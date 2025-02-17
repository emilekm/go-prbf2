package prism

import (
	"context"
	"errors"
)

type serverService struct {
	c       *Client
	started bool
}

// Using `serverdetailsalways` instead of `serverdetails`
// because the latter might not return all the fields.
func (s *serverService) Details(ctx context.Context) (*ServerDetails, error) {
	resp, err := s.c.Send(ctx, &Request{
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

	s.started = true

	return &serverDetails, nil
}

func (s *serverService) DetailsUpdates(ctx context.Context) (Subscriber, error) {
	if !s.started {
		_, err := s.Details(ctx)
		if err != nil {
			return nil, err
		}
	}

	sub := s.c.Subscribe(SubjectUpdateServerDetails)
	return sub, nil
}

type gameplayService struct {
	c       *Client
	started bool
}

func (s *gameplayService) Details(ctx context.Context) (*GameplayDetails, error) {
	resp, err := s.c.Send(ctx, &Request{
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

	s.started = true

	return &gameplayDetails, nil
}

func (s *gameplayService) DetailsUpdates(ctx context.Context) (Subscriber, error) {
	if !s.started {
		_, err := s.Details(ctx)
		if err != nil {
			return nil, err
		}
	}

	sub := s.c.Subscribe(SubjectUpdateServerDetails)
	return sub, nil
}

type playersService struct {
	c       *Client
	started bool
}

func (s *playersService) List(ctx context.Context) (Players, error) {
	resp, err := s.c.Send(ctx, &Request{
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

	s.started = true

	return players, nil
}

func (s *playersService) ListUpdates(ctx context.Context) (Subscriber, error) {
	if !s.started {
		_, err := s.List(ctx)
		if err != nil {
			return nil, err
		}
	}

	sub := s.c.Subscribe(SubjectUpdatePlayers)
	return sub, nil
}

func (s *playersService) PlayerLeaveUpdates(ctx context.Context) (Subscriber, error) {
	if !s.started {
		_, err := s.List(ctx)
		if err != nil {
			return nil, err
		}
	}

	sub := s.c.Subscribe(SubjectPlayerLeave)
	return sub, nil
}

type adminService struct {
	c *Client
}

func (s *adminService) APIAdmin(ctx context.Context, command string) (string, error) {
	resp, err := s.c.Send(ctx, &Request{
		Message:         NewMessage(CommandAPIAdmin, []byte(command)),
		ExpectedSubject: SubjectAPIAdminResult,
	})
	if err != nil {
		return "", err
	}

	return string(resp.Message.Body()), nil
}

func (s *adminService) RACommand(ctx context.Context, command string) (*RACommandOutcome, error) {
	resp, err := s.c.Send(ctx, &Request{
		Message:         NewMessage(CommandRACommand, []byte(command)),
		ExpectedSubject: SubjectSuccess,
	})
	if err != nil {
		var msgErr Error
		if errors.As(err, &msgErr) {
			return nil, err
		}
		var msg RACommandOutcome
		err2 := Unmarshal(resp.Message.Body(), &msg)
		if err2 != nil {
			return nil, err
		}
		return &msg, nil
	}

	var msg RACommandOutcome
	err = Unmarshal(resp.Message.Body(), &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}
