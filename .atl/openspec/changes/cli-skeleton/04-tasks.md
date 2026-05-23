# SDD Tasks: CLI Skeleton — TUI Base + Avatar Animado

**Change ID:** cli-skeleton  
**Author:** yorch  
**Date:** 2026-05-23  
**Status:** draft  
**Phase:** tasks

---

## 1. Overview

Este documento desglosa el diseño en tareas ejecutables para cada PR encadenado.

---

## 2. PR-1: Estructura proyecto + Domain base

### 2.1 Estructura inicial

**Tarea 1.1:** Crear estructura de directorios

```bash
mkdir -p cmd/elena
mkdir -p internal/core/domain
mkdir -p internal/app
mkdir -p internal/core/ports/output
mkdir -p internal/core/usecases/chat
mkdir -p internal/core/usecases/identity
mkdir -p internal/infrastructure/adapters/tui/input
mkdir -p internal/infrastructure/adapters/tui/output
mkdir -p internal/infrastructure/adapters/tui/bubbletea
mkdir -p internal/infrastructure/adapters/mock/response_generator
mkdir -p internal/infrastructure/entrypoints
```

### 2.2 go.mod

**Tarea 1.2:** Crear `go.mod`

```bash
touch go.mod
```

Contenido:
```
module elena

go 1.21

require (
    github.com/charmbracelet/bubbletea v0.27.0
    github.com/charmbracelet/lipgloss v0.11.0
)
```

### 2.3 main.go

**Tarea 1.3:** Crear `cmd/elena/main.go`

```go
package main

import (
    "elena/internal/app"
    "elena/internal/infrastructure/entrypoints"
)

func main() {
    a := app.Wire()
    if err := entrypoints.StartTUI(a); err != nil {
        panic(err)
    }
}
```

### 2.4 Domain: Mood

**Tarea 1.4:** Crear `internal/core/domain/mood.go`

```go
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
```

### 2.5 Domain: Message

**Tarea 1.5:** Crear `internal/core/domain/message.go`

```go
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

func (m *Message) ID() string            { return m.id }
func (m *Message) Author() Author         { return m.author }
func (m *Message) Content() string        { return m.content }
func (m *Message) Timestamp() time.Time  { return m.timestamp }
```

### 2.6 Domain: Session

**Tarea 1.6:** Crear `internal/core/domain/session.go`

```go
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

func (s *Session) ID() string             { return s.id }
func (s *Session) Mood() Mood            { return s.mood }
func (s *Session) Messages() []*Message  {
    out := make([]*Message, len(s.messages))
    copy(out, s.messages)
    return out
}
func (s *Session) StartedAt() time.Time   { return s.startedAt }
func (s *Session) MessageCount() int     { return len(s.messages) }

func (s *Session) AddMessage(msg *Message) { s.messages = append(s.messages, msg) }
func (s *Session) SetMood(m Mood)          { s.mood = m }
func (s *Session) ClearMessages()          { s.messages = make([]*Message, 0) }
```

### 2.7 Verificación PR-1

**Tarea 1.7:** Verificar compilación

```bash
go build ./...
```

---

## 3. PR-2: Use cases + Output ports

### 3.1 Output ports

**Tarea 2.1:** Crear `internal/core/ports/output/display_port.go`

```go
package output

type DisplayPort interface {
    Show(text string) error
}
```

**Tarea 2.2:** Crear `internal/core/ports/output/response_port.go`

```go
package output

type ResponsePort interface {
    Generate(text string) string
}
```

### 3.2 Chat use case

**Tarea 2.3:** Crear `internal/core/usecases/chat/chat.go`

```go
package chat

import (
    "context"
    "elena/internal/core/domain"
    "elena/internal/core/ports/output"
)

type ChatUseCaseInput struct {
    ResponsePort output.ResponsePort
    DisplayPort  output.DisplayPort
}

type ChatUseCase struct {
    input ChatUseCaseInput
}

func NewChatUseCase(input ChatUseCaseInput) *ChatUseCase {
    return &ChatUseCase{input: input}
}

func (c *ChatUseCase) Execute(ctx context.Context, session *domain.Session, userText string) (string, error) {
    userMsg := domain.NewMessage(generateID(), domain.AuthorUser, userText)
    session.AddMessage(userMsg)
    
    session.SetMood(domain.MoodProcessing)
    
    response := c.input.ResponsePort.Generate(userText)
    
    elenaMsg := domain.NewMessage(generateID(), domain.AuthorElena, response)
    session.AddMessage(elenaMsg)
    
    session.SetMood(domain.MoodIdle)
    
    if err := c.input.DisplayPort.Show(response); err != nil {
        return "", err
    }
    
    return response, nil
}

func generateID() string {
    return "placeholder-id"
}
```

### 3.3 Identity use case

**Tarea 2.4:** Crear `internal/core/usecases/identity/identity.go`

```go
package identity

import (
    "elena/internal/core/domain"
    "elena/internal/core/ports/output"
)

type IdentityUseCaseInput struct {
    DisplayPort output.DisplayPort
}

type IdentityUseCase struct {
    input IdentityUseCaseInput
}

func NewIdentityUseCase(input IdentityUseCaseInput) *IdentityUseCase {
    return &IdentityUseCase{input: input}
}

func (c *IdentityUseCase) GetCurrentMood(session *domain.Session) domain.Mood {
    return session.Mood()
}

func (c *IdentityUseCase) ChangeMood(session *domain.Session, newMood domain.Mood) {
    session.SetMood(newMood)
}

func (c *IdentityUseCase) ShowMood(session *domain.Session) error {
    text := "Estado actual: " + session.Mood().String()
    return c.input.DisplayPort.Show(text)
}
```

### 3.4 Mock response implementation

**Tarea 2.5:** Crear `internal/infrastructure/adapters/mock/response_generator/service.go`

```go
package response_generator

type MockResponseService struct {
    responses []string
    index     int
}

func NewMockResponseService() *MockResponseService {
    return &MockResponseService{
        responses: []string{
            "¡Hola! ¿Cómo estás?",
            "Interesante...",
            "Cuéntame más.",
            "Hmm, déjame pensar.",
            "¡Qué buena pregunta!",
            "No estoy segura, pero...",
            "¿Y eso por qué?",
            "Ah, ya veo.",
        },
    }
}

func (s *MockResponseService) Generate(text string) string {
    response := s.responses[s.index]
    s.index = (s.index + 1) % len(s.responses)
    return response
}
```

### 3.5 Application bootstrap

**Tarea 2.6:** Crear `internal/app/application.go`

```go
package app

import (
    "elena/internal/core/domain"
    "elena/internal/core/usecases/chat"
    "elena/internal/core/usecases/identity"
    "elena/internal/infrastructure/adapters/mock/response_generator"
    "elena/internal/core/ports/output"
)

type App struct {
    session         *domain.Session
    chatUseCase     *chat.ChatUseCase
    identityUseCase *identity.IdentityUseCase
}

func Wire() *App {
    responsePort := response_generator.NewMockResponseService()
    
    return &App{
        session: domain.NewSession("session-1"),
        chatUseCase: chat.NewChatUseCase(chat.ChatUseCaseInput{
            ResponsePort: responsePort,
            DisplayPort:  &stubDisplayPort{},
        }),
        identityUseCase: identity.NewIdentityUseCase(identity.IdentityUseCaseInput{
            DisplayPort: &stubDisplayPort{},
        }),
    }
}

func (a *App) Session() *domain.Session    { return a.session }
func (a *App) ChatUseCase() *chat.ChatUseCase     { return a.chatUseCase }
func (a *App) IdentityUseCase() *identity.IdentityUseCase { return a.identityUseCase }

func (a *App) Start() error {
    return nil
}

type stubDisplayPort struct{}

func (s *stubDisplayPort) Show(text string) error {
    return nil
}
```

### 3.6 Verificación PR-2

**Tarea 2.7:** Verificar compilación

```bash
go build ./...
```

---

## 4. PR-3: TUI adapter

### 4.1 Entrypoint

**Tarea 3.1:** Crear `internal/infrastructure/entrypoints/tui.go`

```go
package entrypoints

import (
    "elena/internal/app"
    "elena/internal/infrastructure/adapters/tui/bubbletea"
)

func StartTUI(a *app.App) error {
    model := bubbletea.NewModel(
        a.Session(),
        a.ChatUseCase(),
        a.IdentityUseCase(),
    )
    return bubbletea.NewProgram(model).Start()
}
```

### 4.2 BubbleTea model

**Tarea 3.2:** Crear `internal/infrastructure/adapters/tui/bubbletea/model.go`

```go
package bubbletea

import (
    "github.com/charmbracelet/bubbletea"
    "elena/internal/core/domain"
    "elena/internal/core/usecases/chat"
    "elena/internal/core/usecases/identity"
    "elena/internal/infrastructure/adapters/tui/output"
)

type Model struct {
    session         *domain.Session
    chatUseCase     *chat.ChatUseCase
    identityUseCase *identity.IdentityUseCase
    currentInput    string
    renderer        *output.Renderer
    quit            bool
}

func NewModel(session *domain.Session, chatUC *chat.ChatUseCase, identityUC *identity.IdentityUseCase) *Model {
    return &Model{
        session:         session,
        chatUseCase:     chatUC,
        identityUseCase: identityUC,
        currentInput:    "",
        renderer:        output.NewRenderer(),
        quit:            false,
    }
}

func (m *Model) Init() tea.Cmd {
    return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.Type == tea.KeyEnter {
            if len(m.currentInput) == 0 {
                return m, nil
            }
            
            // Check if command
            if cmd, ok := parseCommand(m.currentInput); ok {
                return m, m.handleCommand(cmd)
            }
            
            // Process message
            _, _ = m.chatUseCase.Execute(nil, m.session, m.currentInput)
            m.currentInput = ""
            return m, nil
        }
        
        if msg.Type == tea.KeyBackspace {
            if len(m.currentInput) > 0 {
                m.currentInput = m.currentInput[:len(m.currentInput)-1]
            }
        }
        
        if msg.Type == tea.KeyRune {
            m.currentInput += msg.String()
        }
    }
    
    return m, nil
}

func (m *Model) View() string {
    if m.quit {
        return "¡Hasta luego!\n"
    }
    return m.renderer.Render(m.session, m.currentInput)
}

func parseCommand(input string) (string, bool) {
    if len(input) > 0 && input[0] == '/' {
        return input, true
    }
    return "", false
}

func (m *Model) handleCommand(cmd string) tea.Cmd {
    return func() tea.Msg {
        switch cmd {
        case "/exit":
            m.quit = true
        case "/reset":
            m.session.ClearMessages()
        case "/mood":
            m.identityUseCase.ShowMood(m.session)
        }
        return nil
    }
}
```

### 4.3 Renderer

**Tarea 3.3:** Crear `internal/infrastructure/adapters/tui/output/renderer.go`

```go
package output

import (
    "github.com/charmbracelet/lipgloss"
    "elena/internal/core/domain"
)

type Renderer struct {
    width int
    height int
}

func NewRenderer() *Renderer {
    return &Renderer{
        width:  80,
        height: 24,
    }
}

func (r *Renderer) Render(session *domain.Session, currentInput string) string {
    avatar := r.renderAvatar(session.Mood())
    chat := r.renderChat(session, currentInput)
    
    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        lipgloss.NewStyle().Width(25).Render(avatar),
        lipgloss.NewStyle().Width(55).Render(chat),
    )
}

func (r *Renderer) renderAvatar(mood domain.Mood) string {
    return "  /\\    /\\\n ( ◕‿◕)\n  \\__/  \\__/\n\nEstado: " + mood.String()
}

func (r *Renderer) renderChat(session *domain.Session, currentInput string) string {
    var lines []string
    
    for _, msg := range session.Messages() {
        prefix := "Tú: "
        if msg.Author() == domain.AuthorElena {
            prefix = "Elena: "
        }
        lines = append(lines, prefix+msg.Content())
    }
    
    lines = append(lines, "")
    lines = append(lines, "[Escribe algo...] "+currentInput)
    
    result := ""
    for _, line := range lines {
        result += line + "\n"
    }
    return result
}

func (r *Renderer) Show(text string) error {
    // DisplayPort implementation - for logging/debugging
    return nil
}
```

### 4.4 Verificación PR-3

**Tarea 3.4:** Verificar compilación

```bash
go build ./...
```

**Tarea 3.5:** Test manual

```bash
go run ./cmd/elena
```

---

## 5. PR-4: Avatar animado + Commands

### 5.1 Avatar

**Tarea 4.1:** Crear `internal/infrastructure/adapters/tui/avatar.go`

```go
package tui

import "elena/internal/core/domain"

var avatarFrames = map[domain.Mood][]string{
    domain.MoodIdle: {
        `  /\    /\
 ( ◕‿◕)
  \__/  \__/`,
        `  /\    /\
 ( ◕ ‿ ◕ )
  \__/  \__/`,
    },
    domain.MoodProcessing: {
        `  /\    /\
 ( ◕ ○ ◕)
  \__/  \__/`,
        `  /\    /\
 ( ◠ ‿ ◠)
  \__/  \__/`,
    },
}

type Avatar struct {
    currentMood domain.Mood
    frameIndex  int
}

func NewAvatar(mood domain.Mood) *Avatar {
    return &Avatar{
        currentMood: mood,
        frameIndex:  0,
    }
}

func (a *Avatar) CurrentFrame() string {
    frames, ok := avatarFrames[a.currentMood]
    if !ok {
        return ""
    }
    return frames[a.frameIndex % len(frames)]
}

func (a *Avatar) Tick() {
    a.frameIndex++
}

func (a *Avatar) SetMood(m domain.Mood) {
    a.currentMood = m
    a.frameIndex = 0
}
```

### 5.2 Command handler

**Tarea 4.2:** Crear `internal/infrastructure/adapters/tui/input/command.go`

```go
package tui

import (
    "strings"
    "elena/internal/core/domain"
    "elena/internal/core/usecases/identity"
)

type CommandHandler struct {
    identityUseCase *identity.IdentityUseCase
}

func NewCommandHandler(identityUC *identity.IdentityUseCase) *CommandHandler {
    return &CommandHandler{
        identityUseCase: identityUC,
    }
}

func (h *CommandHandler) Parse(input string) (cmd string, isCommand bool) {
    if len(input) > 0 && input[0] == '/' {
        parts := strings.SplitN(input, " ", 2)
        return parts[0], true
    }
    return "", false
}

func (h *CommandHandler) Execute(cmd string, session *domain.Session) string {
    switch cmd {
    case "/exit":
        return "¡Hasta luego! Fue agradable hablar contigo."
    case "/reset":
        session.ClearMessages()
        return "Conversación reiniciada."
    case "/mood":
        mood := h.identityUseCase.GetCurrentMood(session)
        return "Estado actual: " + mood.String()
    default:
        return "Comando no reconocido."
    }
}
```

### 5.3 Actualizar BubbleTea model con avatar y ticker

**Tarea 4.3:** Actualizar `internal/infrastructure/adapters/tui/bubbletea/model.go`

```go
package bubbletea

import (
    "time"
    
    "github.com/charmbracelet/bubbletea"
    "elena/internal/core/domain"
    "elena/internal/core/usecases/chat"
    "elena/internal/core/usecases/identity"
    "elena/internal/infrastructure/adapters/tui"
    "elena/internal/infrastructure/adapters/tui/output"
)

type tickMsg struct{}

type Model struct {
    session         *domain.Session
    chatUseCase     *chat.ChatUseCase
    identityUseCase *identity.IdentityUseCase
    avatar          *tui.Avatar
    currentInput    string
    renderer        *output.Renderer
    cmdHandler      *tui.CommandHandler
    quit            bool
}

func NewModel(session *domain.Session, chatUC *chat.ChatUseCase, identityUC *identity.IdentityUseCase) *Model {
    return &Model{
        session:         session,
        chatUseCase:     chatUC,
        identityUseCase: identityUC,
        avatar:          tui.NewAvatar(session.Mood()),
        currentInput:    "",
        renderer:        output.NewRenderer(),
        cmdHandler:      tui.NewCommandHandler(identityUC),
        quit:            false,
    }
}

func (m *Model) Init() tea.Cmd {
    // Ticker para animación de avatar (500ms)
    return tea.Every(time.Millisecond*500, func(t time.Time) tea.Msg {
        return tickMsg{}
    })
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tickMsg:
        m.avatar.Tick()
        return m, nil
        
    case tea.KeyMsg:
        if msg.Type == tea.KeyEnter {
            if len(m.currentInput) == 0 {
                return m, nil
            }
            
            if cmd, ok := m.cmdHandler.Parse(m.currentInput); ok {
                return m, m.handleCommand(cmd)
            }
            
            _, _ = m.chatUseCase.Execute(nil, m.session, m.currentInput)
            m.avatar.SetMood(m.session.Mood())
            m.currentInput = ""
            return m, nil
        }
        
        if msg.Type == tea.KeyBackspace {
            if len(m.currentInput) > 0 {
                m.currentInput = m.currentInput[:len(m.currentInput)-1]
            }
        }
        
        if msg.Type == tea.KeyRune {
            m.currentInput += msg.String()
        }
    }
    
    return m, nil
}

func (m *Model) View() string {
    if m.quit {
        return "¡Hasta luego!\n"
    }
    return m.renderer.RenderWithAvatar(m.session, m.currentInput, m.avatar.CurrentFrame())
}

func (m *Model) handleCommand(cmd string) tea.Cmd {
    return func() tea.Msg {
        result := m.cmdHandler.Execute(cmd, m.session)
        m.quit = (cmd == "/exit")
        return nil
    }
}
```

### 5.4 Actualizar Renderer con avatar

**Tarea 4.4:** Actualizar `internal/infrastructure/adapters/tui/output/renderer.go`

```go
package output

import (
    "github.com/charmbracelet/lipgloss"
    "elena/internal/core/domain"
)

type Renderer struct {
    width  int
    height int
}

func NewRenderer() *Renderer {
    return &Renderer{
        width:  80,
        height: 24,
    }
}

func (r *Renderer) Render(session *domain.Session, currentInput string) string {
    return r.RenderWithAvatar(session, currentInput, r.renderAvatar(session.Mood()))
}

func (r *Renderer) RenderWithAvatar(session *domain.Session, currentInput string, avatarFrame string) string {
    avatar := lipgloss.NewStyle().
        Width(25).
        Align(lipgloss.Center).
        Render(avatarFrame + "\n\nEstado: " + session.Mood().String())
    
    chat := r.renderChat(session, currentInput)
    
    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        avatar,
        lipgloss.NewStyle().Width(55).Render(chat),
    )
}

func (r *Renderer) renderAvatar(mood domain.Mood) string {
    return "  /\\    /\\\n ( ◕‿◕)\n  \\__/  \\__/"
}

func (r *Renderer) renderChat(session *domain.Session, currentInput string) string {
    var lines []string
    
    for _, msg := range session.Messages() {
        prefix := "Tú: "
        if msg.Author() == domain.AuthorElena {
            prefix = "Elena: "
        }
        lines = append(lines, prefix+msg.Content())
    }
    
    lines = append(lines, "")
    lines = append(lines, "[Escribe algo...] "+currentInput)
    
    result := ""
    for _, line := range lines {
        result += line + "\n"
    }
    return result
}

func (r *Renderer) Show(text string) error {
    return nil
}
```

### 5.5 Verificación PR-4

**Tarea 4.5:** Verificar compilación

```bash
go build ./...
```

**Tarea 4.6:** Test manual

```bash
go run ./cmd/elena
# Escribir "Hola" y Enter
# Verificar respuesta mock
# Escribir "/mood" y Enter
# Verificar estado
# Escribir "/reset" y Enter
# Verificar limpieza
# Escribir "/exit" y Enter
# Verificar salida
```

---

## 6. Checklist de verificación final

### PR-1: Domain
- [ ] `go.mod` creado con dependencias
- [ ] `cmd/elena/main.go` creado
- [ ] `domain/mood.go` con MoodIdle y MoodProcessing
- [ ] `domain/message.go` con campos privados y getters
- [ ] `domain/session.go` con campos privados, getters/setters
- [ ] `go build ./...` pasa

### PR-2: Use cases + Output ports
- [ ] `ports/output/display_port.go` creado
- [ ] `ports/output/response_port.go` creado
- [ ] `usecases/chat/chat.go` creado
- [ ] `usecases/identity/identity.go` creado
- [ ] `adapters/mock/response_generator/service.go` creado
- [ ] `app/application.go` creado con Wire()
- [ ] `go build ./...` pasa

### PR-3: TUI adapter
- [ ] `entrypoints/tui.go` creado
- [ ] `tui/bubbletea/model.go` creado
- [ ] `tui/output/renderer.go` creado
- [ ] `go build ./...` pasa
- [ ] `go run ./cmd/elena` funciona (stub)

### PR-4: Avatar + Commands
- [ ] `tui/avatar.go` creado con frames
- [ ] `tui/input/command.go` creado
- [ ] model.go actualizado con avatar y ticker
- [ ] renderer.go actualizado con avatar
- [ ] `go build ./...` pasa
- [ ] `go run ./cmd/elena` funciona completo
- [ ] Avatar animado (frames rotan)
- [ ] `/mood` muestra estado
- [ ] `/reset` limpia chat
- [ ] `/exit` termina programa

---

## 7. Líneas estimadas por tarea

| Tarea | Descripción | Líneas |
|-------|-------------|--------|
| 1.1 | Estructura dirs | ~10 |
| 1.2 | go.mod | ~10 |
| 1.3 | main.go | ~15 |
| 1.4 | mood.go | ~20 |
| 1.5 | message.go | ~30 |
| 1.6 | session.go | ~35 |
| 2.1 | display_port.go | ~8 |
| 2.2 | response_port.go | ~8 |
| 2.3 | chat.go | ~45 |
| 2.4 | identity.go | ~30 |
| 2.5 | mock/service.go | ~30 |
| 2.6 | application.go | ~45 |
| 3.1 | entrypoints/tui.go | ~15 |
| 3.2 | model.go | ~85 |
| 3.3 | renderer.go | ~60 |
| 4.1 | avatar.go | ~45 |
| 4.2 | command.go | ~40 |
| 4.3 | model.go updates | ~25 |
| 4.4 | renderer.go updates | ~15 |
| **Total** | | **~521** |

---

**Status:** ⏳ Pending approval