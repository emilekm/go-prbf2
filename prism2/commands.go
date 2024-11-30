package prism2

type Request struct {
	Message         *Message
	ExpectedSubject Subject
}

type Response struct {
	Message *Message
	Error   error
}

func (c *Client) Send(req *Request) <-chan Response {
	resp := make(chan Response, 1)

	go func() {
		defer close(resp)
	}()

	return resp
}
