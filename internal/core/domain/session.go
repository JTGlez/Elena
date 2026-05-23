package domain

import "time"

type Session struct {
	id        string
	mood      Mood
	messages  []*Message
	startedAt time.Time
}

func NewSession(id string) *Session {
	return &Session{
		id:        id,
		mood:      MoodIdle,
		messages:  make([]*Message, 0),
		startedAt: time.Now(),
	}
}

func (s *Session) ID() string         { return s.id }
func (s *Session) Mood() Mood         { return s.mood }
func (s *Session) Messages() []*Message {
	out := make([]*Message, len(s.messages))
	copy(out, s.messages)
	return out
}
func (s *Session) StartedAt() time.Time  { return s.startedAt }
func (s *Session) MessageCount() int     { return len(s.messages) }

func (s *Session) AddMessage(msg *Message) { s.messages = append(s.messages, msg) }
func (s *Session) SetMood(m Mood)          { s.mood = m }
func (s *Session) ClearMessages()          { s.messages = make([]*Message, 0) }
