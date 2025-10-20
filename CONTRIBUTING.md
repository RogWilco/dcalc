# Contributing Guidelines

Thanks for helping improve `dcalc`. These guidelines apply to both human contributors and AI assistants.

## Commit Messages (Conventional Commits)

- Follow the [Conventional Commits specification](https://www.conventionalcommits.org/en/v1.0.0/): `<type>[optional scope]: <description>`.
  - Common types: `feat`, `fix`, `docs`, `test`, `refactor`, `build`, `chore`.
  - Example: `feat(calculator): add mixed-number parser`.
- Keep the subject (header line) ≤ 72 characters; description is imperative and lowercase.
- Avoid ending the header with punctuation.
- Reference issues when relevant using parentheses or footer (e.g., `fix: adjust precision (#12)` or `Closes #12` in the body).
- When more detail is needed, add a blank line and wrap body text at ≤ 80 characters per line.

## Workflow Expectations

- Create focused changes and run `go test ./...` before proposing a merge.
- Format code with `gofmt` (or `make fmt`) prior to submission.
- Update documentation (`README.md`, `AGENTS.md`, etc.) when behavior or process changes.
- When adding dependencies, explain why they are needed in the pull request or summary.
