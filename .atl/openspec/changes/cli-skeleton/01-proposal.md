# SDD Proposal: CLI Skeleton — TUI Base + Avatar Animado

**Change ID:** cli-skeleton  
**Author:** yorch  
**Date:** 2026-05-23  
**Status:** approved  
**Phase:** proposal → spec approved

---

## 1. Summary

Implementar la base del CLI de Elena MVP: una TUI interactiva en Go usando Bubble Tea, con un avatar ASCII animado estilo Tamagotchi en layout split (avatar izq + chat der), respuestas dummy rotativas, y comandos básicos (`/exit`, `/reset`, `/mood`).

Este change establece el stack tecnológico y la infraestructura visual mínima para futuras features de memoria, personalidad y persistencia.

---

## 2. Motivation

Elena necesita una presencia visual que la diferencie de un chatbot genérico. La referencia es Meta Soulmate AI (sin el romantic angle): una presencia conversacional persistente con avatar animado que transmite estado y continuidad.

El CLI skeleton es el primer paso porque:
- Define el stack tecnológico (Go + Bubble Tea)
- Establece el look & feel visual de la aplicación
- Proporciona el framework para interacciones futuras
- Permite validar el flujo interactivo antes de agregar complejidad

---

## 3. Scope

### 3.1 In Scope

- Proyecto Go limpio con estructura modular
- TUI basada en Bubble Tea con layout split (avatar | chat)
- Avatar ASCII animado (estilo Tamagotchi)
- Animación de mínimo 2 frames: idle + processing
- Sesión interactiva: type + Enter para enviar mensajes
- Respuestas dummy rotativas (array de frases predeterminadas)
- Comandos: `/exit`, `/reset`, `/mood`
- Separación clara entre componentes

### 3.2 Out of Scope

- Memoria o persistencia (sin StorePort, sin MemoryService)
- LLM real o integración con providers (sin InferPort)
- CommandRouter o sistema de routing formal
- Concurrency model o runtime threads
- AppState centralizado
- Historial navegable
- Soporte para voz

---

## 4. Technical Approach

### 4.1 Stack Tecnológico

- **Lenguaje:** Go 1.21+
- **Framework UI:** Bubble Tea (charm.sh/talk-to-tools)
- **Estilos:** Lip Gloss (charm.sh/lip-gloss)
- **Layout:** Split 30% avatar / 70% chat

### 4.2 Arquitectura Final (Simplificada)

```
elena/
├── cmd/
│   └── elena/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   └── application.go       # Bootstrap simple
│   │
│   ├── core/
│   │   ├── domain/
│   │   │   ├── mood.go
│   │   │   ├── message.go
│   │   │   └── session.go
│   │   │
│   │   ├── usecases/
│   │   │   ├── chat/
│   │   │   │   ├── chat.go      # Input port — mock rotativo
│   │   │   │   └── services/
│   │   │   │       └── dummy_generator/
│   │   │   │           └── service.go
│   │   │   └── identity/
│   │   │       └── identity.go  # Muestra mood
│   │   │
│   │   └── ports/
│   │       └── output/
│   │           └── display_port.go
│   │
│   └── infrastructure/
│       ├── adapters/
│       │   ├── tui/
│       │   │   ├── input/
│       │   │   │   ├── chat.go
│       │   │   │   └── command.go
│       │   │   ├── output/
│       │   │   │   └── renderer.go
│       │   │   └── bubbletea/
│       │   │       ├── model.go
│       │   │       └── view.go
│       │   └── dummy/
│       │       └── responses.go
│       └── entrypoints/
│           └── tui.go
│
├── go.mod
└── go.sum
```

### 4.3 Convenciones de Interfaces

**Principio del consumidor:**
- Interfaz declarada ARRIBA del struct del use case que la consume
- Implementación vive en `use_case/services/myservice/service.go`

**Puertos de salida:**
- Solo DisplayPort necesario en MVP
- No hay StorePort ni InferPort

### 4.4 Decisiones de Diseño

- **Bubble Tea nativamente:** No se forkeara pi-agent. Control total del stack.
- **ASCII animation simple:** Array de strings con frames que rotan via ticker.
- **Dummy responses:** Array de strings que rota en round-robin.
- **Session como aggregate root:** Recibe mensajes, controla mood.
- **DisplayPort.Show():** La única abstracción de output necesaria.

---

## 5. User Experience

### 5.1 Arranque

```
$ elena

╔══════════════════════════════════════╗
║     ┌────────────────────────┐      ║
║     │  /\    /\              │      ║
║     │ ( ◕‿◕) │              │      ║
║     │  \__/  \__/            │      ║
║     └────────────────────────┘      ║
║                                      ║
║     ¡Hola! Soy Elena.               ║
║     Estoy aquí.                      ║
╚══════════════════════════════════════╝

[Escribe algo y presiona Enter...]
```

### 5.2 Flujo Interactivo

```
┌────────────────────────────────────┬──────────────────────────────────┐
│  ┌──────────────────────────┐     │ Elena: ¡Hola! Soy Elena.       │
│  │       /\    /\            │     │                                  │
│  │      ( ◕‿◕)              │     │ Tú: Hola Elena                  │
│  │       \__/  \__/           │     │                                  │
│  └──────────────────────────┘     │ Elena: Me alegra verte.        │
│                                    │                                  │
│  Estado: idle                      │ [Escribe algo...]               │
└────────────────────────────────────┴──────────────────────────────────┘
```

### 5.3 Comandos

| Comando | Descripción | Output |
|---------|-------------|--------|
| `/exit` | Cierra la sesión | Mensaje de despedida + exit code 0 |
| `/reset` | Limpia el historial | Chat vacío + confirmación |
| `/mood` | Muestra el estado emocional | "Estado actual: idle/processing" |

---

## 6. Milestones / Chained PRs

| PR | Contenido | Líneas estimadas |
|----|-----------|------------------|
| PR-1 | Estructura proyecto + go.mod + domain base | ~150 |
| PR-2 | Use cases (ChatUseCase + IdentityUseCase) + DisplayPort | ~200 |
| PR-3 | TUI adapter (BubbleTea model + view + layout split) | ~350 |
| PR-4 | Avatar animado (frames ASCII + state machine) + commands | ~300 |
| **Total** | | **~1000** |

---

## 7. Risks and Mitigations

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|-------------|---------|------------|
| Bubble Tea curva de aprendizaje | Media | Medio | Documentación oficial excelente |
| Frames ASCII en fonts proporcionales | Baja | Alto | Verificar en múltiples terminales |
| Flicker en re-renders | Media | Medio | Doble buffering de Bubble Tea |

---

## 8. Approval Gate

**Status:** ✅ APPROVED

- Scope de features definido y acordado
- Stack tecnológico aprobado (Go + Bubble Tea + Lip Gloss)
- Arquitectura simplificada (sin especulación)
- Sin persistencia, sin LLM, sin CommandRouter
- Solo lo implementable en MVP
- PRs encadenados aceptables (~1000 líneas en 4 PRs)

---

## 9. Next Steps

1. ~~Proposal~~ → ✅ Done
2. ~~Spec~~ → ✅ Approved
3. **Design** → In progress
4. **Tasks** → Pending
5. **Apply** → Pending
6. **Verify** → Pending