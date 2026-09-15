package prism

import (
	"context"
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
	defer c.pipeline.EndResponse(id)

	var subjects []Subject
	if req.ExpectedSubject != Subject("") {
		subjects = append(subjects, req.ExpectedSubject)
	}
	if !req.IgnoreError {
		subjects = append(subjects, SubjectError, SubjectCriticalError)
	}

	if len(subjects) == 0 {
		return nil, nil
	}

	waiter := c.beginWait(subjects...)
	defer c.endWait(waiter)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-waiter.result:
		if result.err != nil {
			return nil, result.err
		}

		switch result.subject {
		case SubjectError, SubjectCriticalError:
			var errMsg Error
			if err := Unmarshal(result.message.body, &errMsg); err != nil {
				return &Response{Message: result.message}, err
			}
			return nil, errMsg
		default:
			return &Response{Message: result.message}, nil
		}
	}
}
