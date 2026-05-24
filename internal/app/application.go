package app

import (
	"elena/internal/core/domain"
	"elena/internal/core/usecases/chat"
	"elena/internal/core/usecases/identity"
	"elena/internal/infrastructure/adapters/mock/response_generator"
	"elena/internal/infrastructure/adapters/tui/output"
)

// App holds the application dependencies.
type App struct {
	session         *domain.Session
	chatUseCase     *chat.ChatUseCase
	identityUseCase *identity.IdentityUseCase
	displayPort     *output.TUIDisplay
}

// Wire assembles the application dependencies and returns a wired App.
// The displayPort is created without a notification callback — set it after
// the BubbleTea model is created: app.DisplayPort().Notify = model.SetNotification
func Wire() *App {
	mockService := mock.NewService()

	displayPort := output.NewTUIDisplay()

	chatUC := chat.New(chat.Input{
		ResponsePort: mockService,
		DisplayPort:  displayPort,
	})

	identityUC := identity.New(identity.Input{
		DisplayPort: displayPort,
	})

	return &App{
		session:         domain.NewSession("session-1"),
		chatUseCase:     chatUC,
		identityUseCase: identityUC,
		displayPort:     displayPort,
	}
}

// Session returns the current domain Session.
func (a *App) Session() *domain.Session { return a.session }

// ChatUseCase returns the chat use case.
func (a *App) ChatUseCase() *chat.ChatUseCase { return a.chatUseCase }

// IdentityUseCase returns the identity use case.
func (a *App) IdentityUseCase() *identity.IdentityUseCase { return a.identityUseCase }

// DisplayPort returns the TUI display port.
// Wire the notification callback after the model is created:
//   app.DisplayPort().Notify = model.SetNotification
func (a *App) DisplayPort() *output.TUIDisplay { return a.displayPort }