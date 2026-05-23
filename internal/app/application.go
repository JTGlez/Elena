package app

import "elena/internal/core/domain"

type App struct {
	session *domain.Session
}

func Wire() *App {
	return &App{
		session: domain.NewSession("session-1"),
	}
}

func (a *App) Session() *domain.Session { return a.session }
