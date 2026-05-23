package tui

import "elena/internal/core/domain"

var avatarFrames = map[domain.Mood][]string{
	domain.MoodIdle: {
		"  /\\    /\\ \n ( ◕‿◕)\n  \\__/  \\__/",
		"  /\\    /\\ \n ( ◕ ‿ ◕)\n  \\__/  \\__/",
	},
	domain.MoodProcessing: {
		"  /\\    /\\ \n ( ◕ ○ ◕)\n  \\__/  \\__/",
		"  /\\    /\\ \n ( ◠ ‿ ◠)\n  \\__/  \\__/",
	},
}

// Avatar renders an animated ASCII face with mood-aware frames.
type Avatar struct {
	currentMood domain.Mood
	frameIndex  int
}

// NewAvatar creates an Avatar with the given mood.
func NewAvatar(mood domain.Mood) *Avatar {
	return &Avatar{
		currentMood: mood,
		frameIndex:  0,
	}
}

// CurrentFrame returns the current frame string for the active mood.
func (a *Avatar) CurrentFrame() string {
	frames, ok := avatarFrames[a.currentMood]
	if !ok {
		return ""
	}
	return frames[a.frameIndex%len(frames)]
}

// Tick advances the frame index by one.
func (a *Avatar) Tick() {
	a.frameIndex++
}

// SetMood changes the active mood and resets the frame index.
func (a *Avatar) SetMood(m domain.Mood) {
	a.currentMood = m
	a.frameIndex = 0
}
