# Changelog

All notable changes to **Crona** will be documented in this file.

This project follows **semantic-ish versioning**, but during early development versions may change rapidly as core concepts solidify.

---

## [Unreleased]

### Added
- Local-first **kernel architecture** with HTTP + SSE APIs
- Event-driven **TimerService** with authoritative session segments
- Pomodoro-style **boundary engine** (work → break → work)
- `active_context` model (repo → stream → issue)
- Git-inspired **stash** system for suspending context
- Filesystem-backed **scratchpads**
  - Variable path support:
    - `[[date]]`
    - `[[time]]`
    - `[[datetime]]`
    - `[[timestamp]]`
    - `[[random]]`
- Commit-style **session end messages**
- Session notes parser with structured sections:
  - `::commit`
  - `::context`
  - `::work`
  - `::notes`
- Kernel auto-discovery and auto-launch
- SSE event stream (`/events`) for real-time UI updates
- End-to-end test harness with isolated SQLite databases
- Ink-based **Terminal UI (TUI)** scaffold
- Graceful kernel shutdown handling
- Context change propagation events
- Scratchpad pinning support

---

## [0.1.0] — Initial Architecture Cut

### Core
- Repository, Stream, Issue domain models
- SQLite persistence using Kysely
- Command-based mutation layer
- Central event bus
- Operation log for replay/debugging

### Kernel
- Local kernel bootstrap
- Token-based local auth
- HTTP APIs for:
  - repos
  - streams
  - issues
  - sessions
  - timer
  - stash
  - scratchpads
- Kernel info endpoint (port, token, scratch dir)
- In-memory timer boundary scheduler

### Timer
- Sessions composed of segments
- Pause treated as `rest` segment
- Boundary-driven transitions
- Idempotent start / pause / resume / end
- Derived timer state (no duplicated storage)

### Stash
- Context snapshotting
- Optional session suspension
- Apply / drop semantics
- Event emission on create/apply/drop

### Scratchpads
- Metadata registry (path, name, pinned)
- Filesystem as source of truth
- Kernel-exposed scratch directory

### Testing
- Full kernel E2E coverage
- Session lifecycle tests
- Boundary transition tests
- Stash behavior tests
- Scratchpad API tests

---

## Design Decisions

- Timer state is **never stored**, only derived
- Sessions end with commit-style messages
- Scratchpads are **not scoped** (no issue/session binding)
- Kernel owns persistence; UIs are stateless clients
- No background daemons beyond the kernel

---

## Upcoming

### Planned
- TUI command palette (vim/lazygit-style)
- Session amend commands (git-like)
- Session history browser
- Daily planning workflow
- Exportable work logs (Markdown / JSON)
- Embedded editor support for scratchpads
- Ops/debug panel in TUI
- Multi-client kernel attachment (future)

---
