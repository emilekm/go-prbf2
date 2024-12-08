package prism2

import "context"

type Request struct {
	Message         *Message
	ExpectedSubject Subject
	IgnoreError     bool
}

type Response struct {
	Message *Message
}

func (c *Client) Send(ctx context.Context, req *Request) (*Response, error) {
	id := c.Next()
	c.StartRequest(id)

	err := c.WriteMessage(req.Message)
	if err != nil {
		c.EndRequest(id)
		return nil, err
	}

	c.StartResponse(id)
	c.EndRequest(id)

	// TODO: check if case with empty channel is valid
	channels := make(map[Subject]Subscriber)

	if req.ExpectedSubject != Subject("") {
		channels[req.ExpectedSubject] = c.Subscribe(req.ExpectedSubject)
	}

	if !req.IgnoreError {
		channels[SubjectError] = c.Subscribe(SubjectError)
		channels[SubjectCriticalError] = c.Subscribe(SubjectCriticalError)
	}

	defer c.EndResponse(id)

	if len(channels) == 0 {
		return nil, nil
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case message := <-channels[SubjectCriticalError]:
		var errMsg Error
		err = Unmarshal(message.body, &errMsg)
		if err != nil {
			return nil, err
		}

		return nil, errMsg
	case message := <-channels[SubjectError]:
		var errMsg Error
		err = Unmarshal(message.body, &errMsg)
		if err != nil {
			return nil, err
		}

		return nil, errMsg
	case message := <-channels[req.ExpectedSubject]:
		return &Response{
			Message: message,
		}, nil
	}
}
