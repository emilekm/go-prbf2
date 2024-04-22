package prism

import "bytes"

type RawMessage struct {
	subject Subject
	content []byte
}

func NewRawMessage(subject Subject, content []byte) RawMessage {
	return RawMessage{
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

func (m RawMessage) Encode() []byte {
	return bytes.Join([][]byte{
		SeparatorStart,
		stringToBytes(string(m.subject)),
		SeparatorSubject,
		m.content,
		SeparatorEnd,
	}, []byte{})
}

func (m *RawMessage) Decode(data []byte) error {
	subject, content, err := decodeData(data)
	if err != nil {
		return err
	}

	m.subject = subject
	m.content = content

	return nil
}
