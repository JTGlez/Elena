package domain

type Mood int

const (
	MoodIdle        Mood = iota
	MoodProcessing
)

func (m Mood) String() string {
	switch m {
	case MoodIdle:
		return "idle"
	case MoodProcessing:
		return "processing"
	default:
		return "unknown"
	}
}
