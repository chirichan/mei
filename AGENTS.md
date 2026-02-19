# AGENTS.md

## Project Overview
- Language: Go (`go 1.26.0`)
- Module: `github.com/chirichan/mei`
- Type: CLI toolkit with multiple commands
- Main binaries:
  - `cmd/pwdgen`
  - `cmd/bing15`

## Repository Layout
- `cmd/`: executable entrypoints
  - `cmd/pwdgen/`: password and utility CLI
  - `cmd/bing15/`: secondary CLI
- `internal/entities/`: shared data models
- `version/`: version metadata
- `testdata/`: local testing fixtures
- `release/`: packaging outputs
- `build.sh`: project build script

## Working Rules For Agents
- Keep changes minimal and focused on the requested task.
- Prefer editing existing files over introducing new abstractions.
- Do not rename public commands/flags unless explicitly requested.
- Preserve CLI behavior and backward compatibility for existing flags.
- Never commit secrets from `.env`.

## Build And Run
- Install dependencies:
  - `go mod tidy`
- Build all:
  - `go build ./...`
- Run pwdgen locally:
  - `go run ./cmd/pwdgen --help`
- Run bing15 locally:
  - `go run ./cmd/bing15 --help`

## Testing And Validation
- Run tests:
  - `go test ./...`
- For CLI changes, validate:
  - help output
  - changed flags/arguments
  - error paths for invalid input

## Coding Conventions
- Follow standard Go formatting:
  - `gofmt` on changed files
- Keep package boundaries clear:
  - command parsing in `cmd/...`
  - data structures in `internal/entities`
- Prefer explicit error propagation with context.

## Change Checklist
- Code compiles with `go build ./...`
- Tests pass with `go test ./...` (if tests exist)
- No accidental binary or generated artifact changes unless required
- Documentation/README updated when behavior changes

## Notes Specific To This Repo
- `cmd/pwdgen/cmd.go` is feature-heavy; avoid broad refactors in small tasks.
- Some output text appears to contain legacy encoding artifacts; do not normalize text unless requested.
