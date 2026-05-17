
# Elena MVP

## Purpose

The initial MVP of Elena is not focused on intelligence, autonomy, or advanced agent capabilities.

The purpose of the MVP is to validate a simpler and more important idea:

> Can a persistent terminal presence begin to feel continuous through state, memory, and behavior alone?

The first milestone is not AGI.

The first milestone is presence.

---

# Technology Decision

## Language

Elena will initially be developed in Go.

### Why Go?

The decision prioritizes:
- fast iteration,
- architectural clarity,
- concurrency,
- maintainability,
- and rapid experimentation.

At this stage, discovering the correct cognitive architecture is more important than low-level optimization.

---

# Architectural Direction

The project will follow:
- incremental development,
- clean architecture principles,
- event-driven thinking,
- and spec-driven development.

The system should remain modular from the beginning, even if many components are initially implemented as placeholders or dummy behaviors.

---

# Initial MVP Scope

The first MVP focuses on:

- terminal presence,
- visual identity,
- simple internal states,
- and persistence.

The MVP does NOT require:
- LLM providers,
- tool execution,
- autonomy,
- embeddings,
- or complex reasoning.

---

# Initial Features

## CLI/TUI Runtime

A terminal-based interface will act as Elena’s primary environment.

The interface itself is considered part of the agent’s identity.

---

## ASCII Avatar

Elena will include a minimal ASCII avatar capable of expressing simple states such as:
- idle,
- attentive,
- reflective,
- sleeping,
- or processing.

The avatar is not decorative.

It is part of the perception of presence.

---

## Emotional/Behavioral States

The MVP will experiment with lightweight internal states such as:
- curiosity,
- attentiveness,
- reflection,
- calmness,
- or inactivity.

These are not intended to simulate human emotions.

Their purpose is to create continuity and behavioral variation.

---

## Dummy Conversation Runtime

Initial responses may be fully mocked or scripted.

The first phase does not require real inference.

The objective is to validate:
- interaction flow,
- state transitions,
- persistence,
- and interface behavior.

---

## Persistence

Even in the earliest version, Elena should preserve:
- session state,
- mood/state,
- and basic interaction history.

Persistence is considered fundamental from day one.

---

# Proposed CLI/TUI Stack

The current preferred stack is based on the Charm ecosystem:

- Bubble Tea
- Lip Gloss
- Bubbles

These libraries provide:
- reactive terminal interfaces,
- animation support,
- stateful rendering,
- and flexible TUI composition.

This ecosystem is considered highly aligned with Elena’s goals.

---

# Early Development Philosophy

The project should avoid premature complexity.

The focus is:
- continuity before autonomy,
- memory before intelligence,
- architecture before optimization,
- and presence before capability.

---

# First Milestone

The first successful milestone for Elena is simple:

> A terminal-based entity that feels persistent, reactive, and present — even without a real LLM behind it. Think of it as a Tamagotchi with steroids.
