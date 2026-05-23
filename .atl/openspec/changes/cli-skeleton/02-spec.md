# SDD Spec: CLI Skeleton — TUI Base + Avatar Animado

**Change ID:** cli-skeleton  
**Author:** yorch  
**Date:** 2026-05-23  
**Status:** draft  
**Phase:** spec

---

## 1. Overview

Este spec documenta los contratos formales del CLI skeleton de Elena MVP.

**Scope real:**
- Estados simples del avatar (idle, processing)
- Enviar mensajes a Elena y recibir respuestas mock rotativas
- 3 comandos: `/exit`, `/reset`, `/mood`

**No incluye:**
- Memoria, persistencia, inferencia real
- Sistema de eventos, command routing formal
- Concurrency model más allá de single-threaded

---

## 2. Scope Note

Elena MVP es un **conversational runtime skeleton**:

- UI reactiva con avatar animado
- Flujo de mensajes simple (sin memoria real)
- Respuestas mock rotativas
- Estados simples del avatar

**No es:**
- Un agente cognitivo
- Un sistema con memoria persistente
- Un sistema con inferencia real
- Un sistema con personalidad emergente

Esto es intencional. Primero runtime funcional, después sofisticación.

---

## 3. Runtime Model

### 3.1 Runtime Owner

**El runtime owner del MVP es BubbleTea Model.**

Todas las mutaciones de Session ocurren en el BubbleTea runtime thread. No hay goroutines, no hay concurrencia real.

```
BubbleTea Model (runtime owner)
    ├── Session (aggregate root — encapsulado, campos privados)
    ├── ChatUseCase
    ├── IdentityUseCase
    └── DisplayPort → TUI Renderer
```

### 3.2 Concurrency Assumption

**Single-threaded runtime.** Bubble Tea ya es single-threaded en su event loop.

- Session, Mood, Messages son accedidos solo desde el runtime thread
- No se requieren mutexes o channels por ahora

**Cuando se agregue streaming LLM real, se reevaluará:**
- Goroutines para inference
- Channels para comunicación UI ← → Inference
- Posible necesidad de sincronización

---

## 4. Domain Entities

### 4.1 Mood

**File:** `internal/core/domain/mood.go`

```go
package domain

type Mood int

const (
    MoodIdle        Mood = iota // Estado pasivo, esperando interacción
    MoodProcessing              // Procesando respuesta
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

### 4.2 Message

**File:** `internal/core/domain/message.go`

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

// NewMessage recibe el ID desde afuera — domain NO genera IDs.
func NewMessage(id string, author Author, content string) *Message {
    return &Message{
        id:        id,
        author:    author,
        content:   content,
        timestamp: time.Now(),
    }
}

// Getters
func (m *Message) ID() string        { return m.id }
func (m *Message) Author() Author   { return m.author }
func (m *Message) Content() string  { return m.content }
func (m *Message) Timestamp() time.Time { return m.timestamp }
```

### 4.3 Session — Campos privados, getters/setters obligatorios

**File:** `internal/core/domain/session.go`

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

// ─── Getters ───

func (s *Session) ID() string           { return s.id }
func (s *Session) Mood() Mood           { return s.mood }
func (s *Session) Messages() []*Message {
    out := make([]*Message, len(s.messages))
    copy(out, s.messages)
    return out
}
func (s *Session) StartedAt() time.Time { return s.startedAt }
func (s *Session) MessageCount() int   { return len(s.messages) }

// ─── Mutaciones centralizadas ───
// SIEMPRE usar estos métodos. Campos son PRIVADOS — no acceder directamente.

func (s *Session) AddMessage(msg *Message) {
    s.messages = append(s.messages, msg)
}

func (s *Session) SetMood(m Mood) {
    s.mood = m
}

// ClearMessages limpia el historial.
// Nota: el nombre puede evolucionar cuando haya memoria persistente real.
func (s *Session) ClearMessages() {
    s.messages = make([]*Message, 0)
}
```

**Decisiones:**
- Campos son PRIVADOS — el lenguaje fuerza a usar getters/setters
- `AddMessage()` y `SetMood()` centralizan mutaciones
- Permite agregar hooks, validación, logging sin cambiar callers
- `ClearMessages()` por ahora — puede renombrarse a `ResetHistory()` cuando haya contexto/sesión/memoria separados

---

## 5. Output Ports del Sistema

### 5.1 DisplayPort

**File:** `internal/core/ports/output/display_port.go`

```go
package output

// DisplayPort es el puerto de salida para mostrar contenido al usuario.
type DisplayPort interface {
    Show(text string) error
}
```

**Nota de diseño futuro:** DisplayPort opera con strings. El conversation model real es `*Message`. Cuando se necesite timestamps, tool traces, streaming, o metadata, el contrato puede evolucionar a `ShowMessage(*domain.Message)` o `RenderConversation(session)`.

### 5.2 ResponsePort

**File:** `internal/core/ports/output/response_port.go`

```go
package output

// ResponsePort es el output port para generación de respuestas.
// En MVP: mock implementation en adapters/mock/response_generator/
// En producción: llamada a provider de LLM
type ResponsePort interface {
    Generate(text string) string
}
```

---

## 6. Use Cases

### 6.1 Chat Use Case

**File:** `internal/core/usecases/chat/chat.go`

```go
package chat

import (
    "context"
    "elena/internal/core/domain"
    "elena/internal/core/ports/output"
)

// ChatUseCaseInput configura el use case con sus dependencias.
type ChatUseCaseInput struct {
    ResponsePort output.ResponsePort // output port
    DisplayPort  output.DisplayPort  // output port
}

// ChatUseCase implementa el input port de chat.
type ChatUseCase struct {
    input ChatUseCaseInput
}

func NewChatUseCase(input ChatUseCaseInput) *ChatUseCase {
    return &ChatUseCase{input: input}
}

// Execute es el input port.
// Recibe session para gobernarla — los mensajes se agregan a session.
func (c *ChatUseCase) Execute(ctx context.Context, session *domain.Session, userText string) (string, error) {
    // 1. Mensaje del usuario → session (usar AddMessage — campo privado)
    userMsg := domain.NewMessage(generateID(), domain.AuthorUser, userText)
    session.AddMessage(userMsg)
    
    // 2. Transición a processing (usar SetMood — campo privado)
    session.SetMood(domain.MoodProcessing)
    
    // 3. Generar respuesta via ResponsePort (output port)
    response := c.input.ResponsePort.Generate(userText)
    
    // 4. Mensaje de Elena → session
    elenaMsg := domain.NewMessage(generateID(), domain.AuthorElena, response)
    session.AddMessage(elenaMsg)
    
    // 5. Transición a idle
    session.SetMood(domain.MoodIdle)
    
    // 6. Mostrar respuesta via DisplayPort (output port)
    if err := c.input.DisplayPort.Show(response); err != nil {
        return "", err
    }
    
    return response, nil
}
```

**Decisiones:**
- Campos de Session son privados — se usan `AddMessage()` y `SetMood()`
- `generateID()` es helper en el package — **future architectural hotspot** (ver sección 8)

---

### 6.2 Identity Use Case

**File:** `internal/core/usecases/identity/identity.go`

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
    return session.Mood() // getter — campo privado
}

func (c *IdentityUseCase) ChangeMood(session *domain.Session, newMood domain.Mood) {
    session.SetMood(newMood) // setter — campo privado
}

func (c *IdentityUseCase) ShowMood(session *domain.Session) error {
    text := "Estado actual: " + session.Mood().String()
    return c.input.DisplayPort.Show(text)
}
```

---

## 7. Mock Response Implementation

**File:** `internal/infrastructure/adapters/mock/response_generator/service.go`

```go
package response_generator

import "math/rand"

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

---

## 8. Design Notes — Future Hotspots

### 8.1 Commands son adapter-local

Los comandos (`/exit`, `/reset`, `/mood`) están implementados directamente en el TUI adapter.

**Esto está bien para MVP**, pero no es arquitectura final.

Cuando se agreguen:
- Discord adapter
- Voice runtime
- WebSocket UI
- API endpoints

El parsing y dispatch de comandos debería vivir en una capa compartida (`core/commands/`).

Por ahora: "MVP commands are adapter-local".

### 8.2 generateID() — Future architectural hotspot

`generateID()` es un helper en el package que usa el use case.

Hoy: single runtime, sin persistencia — no hay problema.

Futuro: cuando haya:
- Persistencia
- Replay
- Import/export
- Sync

El ownership de identidad necesitará formalizarse.

Por ahora: helper simple en el package que lo necesita.

### 8.3 ClearMessages() — Nombre puede evolucionar

Semánticamente:
- `reset UI` ≠ `reset context` ≠ `reset session` ≠ `clear history`

Cuando aparezcan memoria persistente, analytics, replay, el nombre puede evolucionar a `ResetHistory()`, `ClearContext()`, etc.

Por ahora: `ClearMessages()` está bien.

### 8.4 DisplayPort — Gap semántico

DisplayPort opera con strings, pero el conversation model real es `*Message`.

Cuando se necesite timestamps, tool traces, streaming, o metadata, el contrato puede evolucionar a:
- `ShowMessage(*domain.Message)`
- `RenderConversation(session)`

No cambiar ahora —准备好 para el salto.

---

## 9. TUI Adapter

### 9.1 Model — Sin estado duplicado

**File:** `internal/infrastructure/adapters/tui/bubbletea/model.go`

```go
package bubbletea

import (
    "elena/internal/core/domain"
    "elena/internal/core/usecases/chat"
    "elena/internal/core/usecases/identity"
)

type Model struct {
    session         *domain.Session
    chatUseCase     *chat.ChatUseCase
    identityUseCase *identity.IdentityUseCase
    currentInput    string
    // NO hay isProcessing duplicado — se deriva de session.Mood()
}

// IsProcessing retorna true si el mood actual es MoodProcessing.
// Deriva de Session — una sola source of truth.
func (m *Model) IsProcessing() bool {
    return m.session.Mood() == domain.MoodProcessing
}
```

**Decisiones:**
- No duplicar estado. `IsProcessing()` deriva de `session.Mood()`
- Session es la única fuente de verdad

### 9.2 Avatar States

**File:** `internal/infrastructure/adapters/tui/avatar.go`

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
```

### 9.3 Command Handling

**File:** `internal/infrastructure/adapters/tui/input/command.go`

```go
package tui

import "strings"

func ParseCommand(input string) (cmd string, isCommand bool) {
    if len(input) > 0 && input[0] == '/' {
        parts := strings.SplitN(input, " ", 2)
        return parts[0], true
    }
    return "", false
}

func (m *Model) handleCommand(cmd string) {
    switch cmd {
    case "/exit":
        m.displayFarewell()
        m.quit()
    case "/reset":
        m.session.ClearMessages()
        m.displayPort.Show("Conversación reiniciada.")
    case "/mood":
        m.identityUseCase.ShowMood(m.session)
    }
}
```

---

## 10. Application Events (reconocidos, no implementados)

Aunque no hay event bus, estos son los eventos naturales del sistema:

| Evento | Cuándo ocurre |
|--------|---------------|
| `MessageReceived` | Usuario envía mensaje |
| `ResponseGenerated` | Elena genera respuesta |
| `MoodChanged` | Session.SetMood() es llamado |
| `CommandExecuted` | Comando (/exit, /reset, /mood) es ejecutado |

**Futuro:** Cuando se necesite telemetry, analytics, replay, o animation sync, estos eventos se formalizarán en un event bus.

---

## 11. Escenarios de Uso

### 11.1 Usuario envía un mensaje

**Given:** Elena corriendo con mood idle  
**When:** Usuario escribe "Hola" y presiona Enter  
**Then:**
1. Mensaje aparece en chat como "user"
2. Mood cambia a `MoodProcessing` (via `session.SetMood()` — campo privado)
3. Avatar muestra frame processing
4. Respuesta mock aparece como "elena"
5. Mood vuelve a `MoodIdle`
6. Avatar muestra frame idle

### 11.2 Usuario ejecuta /mood

**When:** Usuario escribe "/mood" y presiona Enter  
**Then:**
1. `IdentityUseCase.ShowMood()` es llamado
2. `DisplayPort.Show("Estado actual: idle")`
3. Mensaje aparece en chat

### 11.3 Usuario ejecuta /reset

**When:** Usuario escribe "/reset" y presiona Enter  
**Then:**
1. `Session.ClearMessages()` limpia historial
2. Mensaje de confirmación aparece
3. Chat panel queda vacío

### 11.4 Usuario ejecuta /exit

**When:** Usuario escribe "/exit" y presiona Enter  
**Then:**
1. Mensaje de despedida aparece
2. Programa termina con exit code 0

---

## 12. Contracts — Resumen

### 12.1 ChatUseCase

```go
type ChatUseCaseInput struct {
    ResponsePort output.ResponsePort
    DisplayPort  output.DisplayPort
}

func (c *ChatUseCase) Execute(ctx context.Context, session *Session, userText string) (string, error)
```

### 12.2 IdentityUseCase

```go
type IdentityUseCaseInput struct {
    DisplayPort output.DisplayPort
}

func (c *IdentityUseCase) GetCurrentMood(session *Session) Mood
func (c *IdentityUseCase) ChangeMood(session *Session, newMood Mood)
func (c *IdentityUseCase) ShowMood(session *Session) error
```

### 12.3 Session

```go
// Campos PRIVADOS
func (s *Session) ID() string
func (s *Session) Mood() Mood
func (s *Session) Messages() []*Message
func (s *Session) StartedAt() time.Time
func (s *Session) MessageCount() int

// Mutaciones centralizadas
func (s *Session) AddMessage(msg *Message)
func (s *Session) SetMood(m Mood)
func (s *Session) ClearMessages()
```

### 12.4 Output Ports

```go
type DisplayPort interface {
    Show(text string) error
}

type ResponsePort interface {
    Generate(text string) string
}
```

---

## 13. Dependencias

```
domain/          ← 0 deps (puro) — campos privados, getters/setters
    ↑
ports/output/   ← DisplayPort, ResponsePort (interfaces)
    ↑
use cases/      ← dependen de domain + ports/output
    ↑
adapters/       ← implementan ports/output
    ├── tui/         → DisplayPort
    └── mock/        → ResponsePort
```

---

## 14. Validación

- [x] Domain puro (sin generateID en domain, sin eventos implementados)
- [x] Session con campos PRIVADOS — lenguaje fuerza encapsulamiento
- [x] Session centraliza mutaciones (AddMessage, SetMood) — getters/setters obligatorios
- [x] No hay estado duplicado (isProcessing eliminado — deriva de session.Mood())
- [x] ResponsePort es output port (implementado en adapters/mock/)
- [x] DisplayPort simple con string — gap reconocido para futuro
- [x] Use cases reciben Session y la gobiernan
- [x] Runtime owner documentado (BubbleTea Model)
- [x] Concurrency assumption documentada (single-threaded)
- [x] Scope MVP documentado (no es agente todavía)
- [x] Application events reconocidos (no implementados)
- [x] Future hotspots documentados (commands, generateID, ClearMessages, DisplayPort)
- [x] Escenarios cubren todos los flujos
- [x] Spec es implementable sin ambigüedad

**Status:** ⏳ Pending approval