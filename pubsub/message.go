package main

type Message struct {
	topic string
	body  string
}

func NewMessage(topic string, body string) *Message {
	return &Message{
		topic: topic,
		body:  body,
	}
}

func (m *Message) GetTopic() string {
	return m.topic
}

func (m *Message) GetBody() string {
	return m.body
}
