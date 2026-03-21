package prism

import (
	"context"
	"io"
)

type Request struct {
	Message         *Message
	ExpectedSubject Subject
	IgnoreError     bool
}

type Response struct {
	Message *Message
}

func (c *Client) Send(ctx context.Context, req *Request) (*Response, error) {
	id := c.pipeline.Next()
	c.pipeline.StartRequest(id)

	err := c.WriteMessage(req.Message)
	if err != nil {
		c.pipeline.EndRequest(id)
		return nil, err
	}

	c.pipeline.StartResponse(id)
	c.pipeline.EndRequest(id)

	// TODO: check if case with empty channel is valid
	channels := make(map[Subject]Subscriber)

	if req.ExpectedSubject != Subject("") {
		channels[req.ExpectedSubject] = c.Subscribe(req.ExpectedSubject)
		defer c.Unsubscribe(channels[req.ExpectedSubject])
	}

	if !req.IgnoreError {
		channels[SubjectError] = c.Subscribe(SubjectError)
		channels[SubjectCriticalError] = c.Subscribe(SubjectCriticalError)
		defer c.Unsubscribe(channels[SubjectError])
		defer c.Unsubscribe(channels[SubjectCriticalError])
	}

	defer c.pipeline.EndResponse(id)

	if len(channels) == 0 {
		return nil, nil
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case message, ok := <-channels[SubjectCriticalError]:
		if !ok {
			return nil, io.EOF
		}
		var errMsg Error
		err = Unmarshal(message.body, &errMsg)
		if err != nil {
			return &Response{
				Message: message,
			}, err
		}

		return nil, errMsg
	case message, ok := <-channels[SubjectError]:
		if !ok {
			return nil, io.EOF
		}
		var errMsg Error
		err = Unmarshal(message.body, &errMsg)
		if err != nil {
			return &Response{
				Message: message,
			}, err
		}

		return nil, errMsg
	case message, ok := <-channels[req.ExpectedSubject]:
		if !ok {
			return nil, io.EOF
		}
		return &Response{
			Message: message,
		}, nil
	}
}
