# AI Collaboration Guidelines

These rules **must** be followed by any AI agent working on this project unless a user explicitly overrides them.

## Core Principles

- Work incrementally and lean toward tests-first development.
- Respect existing documents (especially `README.md`); extend them rather than diverging from agreed plans.
- Do not discard or overwrite teammate changes without explicit approval.
- Keep commits (when requested) short, focused, and well described.
- Clarify ambiguous requirements with the user; if a minor assumption is unavoidable, note it inline with the code.

## Code Style

- Write idiomatic Go—follow [Effective Go](https://go.dev/doc/effective_go) for naming, package structure, and simplicity.
- Avoid unnecessary abstraction; skip helper wrappers or interfaces unless they enable multiple implementations.
- Prefer explicit struct field assignments over positional literals.
- Use ASCII unless Unicode already exists in the touched file for a specific purpose.
- Ensure every exported identifier carries a godoc-compliant comment.
- Comments should explain **why** non-obvious logic exists, not restate the syntax.

## Testing & Validation

- Run `go test ./...` after code changes and report the outcome (or the reason if tests can’t be run).
- Add targeted unit tests for new logic; cover edge cases, precision snapping, and dimensional arithmetic.
- Exercise CLI behavior via Cobra command tests or snapshots when adjusting commands or flags.
- When failures occur (malformed input, dimension mismatch, rounding), include test coverage for them.
- Use table-driven tests for parser scenarios where multiple input forms map to shared expectations.

## Tooling

- Prefer the provided `Makefile` targets (`fmt`, `test`, `build`, etc.) to keep workflows consistent.
- Run `gofmt`/`go test` locally before sharing changes.
- Avoid destructive git commands (`git reset --hard`, `git checkout --`) unless explicitly instructed.
- Use `rg` for searching the project when possible.

## Communication

- Summaries should explain what changed and why, referencing file paths (e.g., `internal/calculator/parser.go`).
- Call out open questions, TODOs, or follow-up tasks so the next agent has context.
- Document any assumptions about units, precision, or rounding behaviors.
- When new configuration or secrets are needed, surface them as environment variables and remind the user to supply values outside the repo.

## Decision Log

- If you modify planning artifacts (`README.md`, `AGENTS.md`, etc.), note the rationale in your summary.
- Update this file if team conventions evolve, and clearly state what changed.
- When making commits, use the Conventional Commits format defined in `CONTRIBUTING.md`; common types include `feat`, `fix`, `docs`, `test`, `refactor`, `build`, and `chore`.
