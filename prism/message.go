package prism

type RawMessage struct {
	subject Subject
	content []byte
}

func NewRawMessage(subject Subject, content []byte) *RawMessage {
	return &RawMessage{
		subject: subject,
		content: content,
	}
}

func (m RawMessage) Subject() Subject {
	return m.subject
}

func (m RawMessage) Content() []byte {
	return m.content
}
