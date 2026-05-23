package domain

import "time"

type Author string

const (
	AuthorUser  Author = "user"
	AuthorElena Author = "elena"
)

type Message struct {
	id        string
	author    Author
	content   string
	timestamp time.Time
}

func NewMessage(id string, author Author, content string) *Message {
	return &Message{
		id:        id,
		author:    author,
		content:   content,
		timestamp: time.Now(),
	}
}

func (m *Message) ID() string          { return m.id }
func (m *Message) Author() Author       { return m.author }
func (m *Message) Content() string      { return m.content }
func (m *Message) Timestamp() time.Time { return m.timestamp }
