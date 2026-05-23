# SDD Design: CLI Skeleton — TUI Base + Avatar Animado

**Change ID:** cli-skeleton  
**Author:** yorch  
**Date:** 2026-05-23  
**Status:** draft  
**Phase:** design

---

## 1. Overview

Este design documenta las decisiones de implementación para el CLI skeleton de Elena MVP. Detalla archivo por archivo, dependencias, y orden de creación para los 4 PRs encadenados.

---

## 2. PRs Planificados

| PR | Contenido | Líneas estimadas |
|----|-----------|------------------|
| **PR-1** | Estructura proyecto + go.mod + domain base | ~150 |
| **PR-2** | Use cases + Output ports | ~200 |
| **PR-3** | TUI adapter (BubbleTea model + view + layout) | ~350 |
| **PR-4** | Avatar animado + commands | ~300 |
| **Total** | | **~1000** |

---

## 3. PR-1: Estructura proyecto + Domain base

### 3.1 Archivos a crear

#### `go.mod`

```go
module elena

go 1.21

require (
    github.com/charmbracelet/bubbletea v0.27.0
    github.com/charmbracelet/lipgloss v0.11.0
)
```

#### `cmd/elena/main.go`

```go
package main

import (
    "elena/internal/app"
    "elena/internal/infrastructure/entrypoints"
)

func main() {
    a := app.Wire()
    if err := entrypoints.StartTUI(a); err != nil {
        // handle error
    }
}
```

**Líneas:** ~15

#### `internal/core/domain/mood.go`

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

**Líneas:** ~20

#### `internal/core/domain/message.go`

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

func (m *Message) ID() string             { return m.id }
func (m *Message) Author() Author        { return m.author }
func (m *Message) Content() string       { return m.content }
func (m *Message) Timestamp() time.Time  { return m.timestamp }
```

**Líneas:** ~30

#### `internal/core/domain/session.go`

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

**Líneas:** ~35

### 3.2 Dependencias PR-1

```
domain/ (puro, sin deps externos)
```

### 3.3 Decisiones de PR-1

| Decisión | Justificación |
|-----------|---------------|
| Go 1.21+ | Versión estable con generics |
| Módulos Go con `elena` | Namespace estándar |
| Campos privados en domain | Encapsulamiento real — lenguaje fuerza getters/setters |

---

## 4. PR-2: Use cases + Output ports

### 4.1 Archivos a crear

#### `internal/core/ports/output/display_port.go`

```go
package output

type DisplayPort interface {
    Show(text string) error
}
```

**Líneas:** ~8

#### `internal/core/ports/output/response_port.go`

```go
package output

type ResponsePort interface {
    Generate(text string) string
}
```

**Líneas:** ~8

#### `internal/core/usecases/chat/chat.go`

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
    // Simple UUID v4 placeholder — futuro: IDGenerator interface
    return "placeholder-id"
}
```

**Líneas:** ~45

#### `internal/core/usecases/identity/identity.go`

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

**Líneas:** ~30

#### `internal/infrastructure/adapters/mock/response_generator/service.go`

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

**Líneas:** ~30

#### `internal/app/application.go`

```go
package app

type App struct {
    session *domain.Session
}

func Wire() *App {
    return &App{
        session: domain.NewSession("session-1"),
    }
}

func (a *App) Start() error {
    // Stub — se implementa en PR-3 con BubbleTea
    return nil
}
```

**Líneas:** ~20

### 4.2 Dependencias PR-2

```
domain/ (puro)
    ↑
ports/output/ (DisplayPort, ResponsePort)
    ↑
use cases/ (ChatUseCase, IdentityUseCase)
    ↑
adapters/mock/ (implementa ResponsePort)
```

### 4.3 Decisiones de PR-2

| Decisión | Justificación |
|-----------|---------------|
| Output ports en `core/ports/output/` | Límites canónicos del sistema |
| Interfaz de servicio en consumer | Principio del consumidor |
| generateID() simple | Placeholder — future hotspot |
| Mock en adapters/ | Infraestructura, no core |

---

## 5. PR-3: TUI adapter

### 5.1 Archivos a crear

#### `internal/infrastructure/entrypoints/tui.go`

```go
package entrypoints

import "elena/internal/app"

func StartTUI(a *app.App) error {
    // Crear TUI model con dependencias injectadas
    model := bubbletea.NewModel(
        a.Session(),
        a.ChatUseCase(),
        a.IdentityUseCase(),
    )
    
    p := tea.NewProgram(model)
    _, err := p.Run()
    return err
}
```

**Líneas:** ~20

#### `internal/infrastructure/adapters/tui/bubbletea/model.go`

```go
package bubbletea

import (
    "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
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
}

func NewModel(session *domain.Session, chatUC *chat.ChatUseCase, identityUC *identity.IdentityUseCase) *Model {
    return &Model{
        session:         session,
        chatUseCase:     chatUC,
        identityUseCase: identityUC,
        currentInput:    "",
        renderer:        output.NewRenderer(),
    }
}

func (m *Model) Init() tea.Cmd {
    // Mostrar saludo inicial
    return func() tea.Msg {
        return tea.Msg("")
    }
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.Type == tea.KeyEnter {
            if len(m.currentInput) == 0 {
                return m, nil
            }
            
            // Handle command or message
            if cmd, ok := parseCommand(m.currentInput); ok {
                return m, m.handleCommand(cmd)
            }
            
            // Process message
            return m, func() tea.Msg {
                m.chatUseCase.Execute(nil, m.session, m.currentInput)
                m.currentInput = ""
                return renderMsg{}
            }
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
    return m.renderer.Render(m.session, m.currentInput)
}

func (m *Model) handleCommand(cmd string) tea.Cmd {
    return func() tea.Msg {
        switch cmd {
        case "/exit":
            // handled in View
        case "/reset":
            m.session.ClearMessages()
        case "/mood":
            m.identityUseCase.ShowMood(m.session)
        }
        return renderMsg{}
    }
}

type renderMsg struct{}
```

**Líneas:** ~80

#### `internal/infrastructure/adapters/tui/output/renderer.go`

```go
package output

import (
    "github.com/charmbracelet/lipgloss"
    "elena/internal/core/domain"
)

type Renderer struct {
    baseStyle   lipgloss.Style
    messageStyle lipgloss.Style
}

func NewRenderer() *Renderer {
    return &Renderer{
        baseStyle:    lipgloss.NewStyle().Width(80).Height(24),
        messageStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("green")),
    }
}

func (r *Renderer) Render(session *domain.Session, currentInput string) string {
    // Layout: avatar panel (left) + chat panel (right)
    avatar := r.renderAvatar(session.Mood())
    chat := r.renderChat(session, currentInput)
    
    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        avatar,
        chat,
    )
}

func (r *Renderer) renderAvatar(mood domain.Mood) string {
    // TODO: implementar con frames ASCII
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
    // Implementa DisplayPort — por ahora solo log
    return nil
}
```

**Líneas:** ~60

### 5.2 Dependencias PR-3

```
domain/ (puro)
    ↑
use cases/ (ChatUseCase, IdentityUseCase)
    ↑
ports/output/ (DisplayPort)
    ↑
adapters/tui/ (implementa DisplayPort)
```

### 5.3 Decisiones de PR-3

| Decisión | Justificación |
|-----------|---------------|
| Bubble Tea model como runtime owner | Documentado en spec |
| Renderer como componente separado | Separa lógica de rendering |
| Input handling en Update() | Estándar Bubble Tea |
| No duplicar estado | IsProcessing() deriva de session.Mood() |

---

## 6. PR-4: Avatar animado + commands

### 6.1 Archivos a crear

#### `internal/infrastructure/adapters/tui/avatar.go`

```go
package tui

import (
    "elena/internal/core/domain"
    "github.com/charmbracelet/lipgloss"
)

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

**Líneas:** ~45

#### `internal/infrastructure/adapters/tui/input/command.go`

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

**Líneas:** ~40

#### Actualización de `internal/infrastructure/adapters/tui/bubbletea/model.go`

Agregar ticker para animación:

```go
func (m *Model) Init() tea.Cmd {
    // Ticker para animación de avatar (500ms por frame)
    return tea.Every(time.Millisecond*500, func(t time.Time) tea.Msg {
        return tickMsg{}
    })
}

type tickMsg struct{}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tickMsg:
        // Rotar frame de avatar
        // ...
    }
    // ... resto del Update
}
```

**Líneas adicionales:** ~15

#### Actualización de `internal/infrastructure/adapters/tui/output/renderer.go`

Integrar Avatar en el render:

```go
func (r *Renderer) Render(session *domain.Session, currentInput string, avatar *Avatar) string {
    avatarPanel := r.renderAvatarWithFrame(avatar.CurrentFrame())
    chatPanel := r.renderChat(session, currentInput)
    
    return lipgloss.JoinHorizontal(lipgloss.Top, avatarPanel, chatPanel)
}
```

**Líneas adicionales:** ~10

### 6.2 Dependencias PR-4

```
domain/ (Mood, Session)
    ↑
adapters/tui/avatar.go (frames)
adapters/tui/input/command.go (commands)
```

### 6.3 Decisiones de PR-4

| Decisión | Justificación |
|-----------|---------------|
| Frames como strings en el package | Simple, no requiere assets externos |
| Ticker en Init() | Bubble Tea subscription pattern |
| Commands en adapter-local | Documentado como future hotspot |

---

## 7. Resumen de archivos por PR

### PR-1: Estructura + Domain

| Archivo | Líneas |
|---------|--------|
| `go.mod` | ~10 |
| `cmd/elena/main.go` | ~15 |
| `internal/core/domain/mood.go` | ~20 |
| `internal/core/domain/message.go` | ~30 |
| `internal/core/domain/session.go` | ~35 |
| **Subtotal** | **~110** |

### PR-2: Use cases + Output ports

| Archivo | Líneas |
|---------|--------|
| `internal/core/ports/output/display_port.go` | ~8 |
| `internal/core/ports/output/response_port.go` | ~8 |
| `internal/core/usecases/chat/chat.go` | ~45 |
| `internal/core/usecases/identity/identity.go` | ~30 |
| `internal/infrastructure/adapters/mock/response_generator/service.go` | ~30 |
| `internal/app/application.go` | ~20 |
| **Subtotal** | **~141** |

### PR-3: TUI adapter

| Archivo | Líneas |
|---------|--------|
| `internal/infrastructure/entrypoints/tui.go` | ~20 |
| `internal/infrastructure/adapters/tui/bubbletea/model.go` | ~80 |
| `internal/infrastructure/adapters/tui/output/renderer.go` | ~60 |
| **Subtotal** | **~160** |

### PR-4: Avatar + Commands

| Archivo | Líneas |
|---------|--------|
| `internal/infrastructure/adapters/tui/avatar.go` | ~45 |
| `internal/infrastructure/adapters/tui/input/command.go` | ~40 |
| Updates a model.go (+ticker) | ~15 |
| Updates a renderer.go (+avatar) | ~10 |
| **Subtotal** | **~110** |

### Total: ~521 líneas

**Nota:** Estimación conservadora. Los archivos reales pueden variar ±20%.

---

## 8. Orden de implementación

```
PR-1 (Domain)
    │
    ▼
PR-2 (Use cases + Output ports)
    │
    ▼
PR-3 (TUI adapter — stub funciona)
    │
    ▼
PR-4 (Avatar + Commands — funcionalidad completa)
```

---

## 9. Decisiones técnicas resumen

| Tema | Decisión |
|------|----------|
| Módulos Go | `elena` como module name |
| Encapsulamiento | Campos privados con getters/setters |
| Output ports | `DisplayPort`, `ResponsePort` en `core/ports/output/` |
| Use case servicios | Interfaz declarada donde el consumer la necesita |
| Avatar | Frames como strings, ticker en BubbleTea Init() |
| Commands | Adapter-local por MVP, documentado como hotspot |
| Concurrency | Single-threaded, Bubble Tea event loop |
| generateID() | Placeholder, documentado como hotspot |

---

## 10. Validación

- [x] Cada PR tiene archivos, dependencias, y líneas estimadas
- [x] Dependencias siempre apuntan hacia domain (puro)
- [x] Orden de implementación es válido (PR-1 sin deps externos)
- [x] Decisiones técnicas documentadas
- [x] Future hotspots reconocidos en cada PR

**Status:** ⏳ Pending approval