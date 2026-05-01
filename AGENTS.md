# Repository Guidelines

## Project Structure & Module Organization

Diary is a Go CLI. The binary entrypoint lives in `cmd/diary/`, and implementation packages live under `internal/`.

Key packages:

- `internal/cli` - Cobra commands and command output handling.
- `internal/project` - project/root resolution.
- `internal/storage` - `.diary/` paths, Markdown records, indexes, and lookups.
- `internal/hash` - content hashing helpers.
- `internal/install` - harness skill installation and templates.
- `internal/update` - GitHub release self-update logic.
- `internal/render` - Markdown and JSON rendering.

Do not reintroduce planning docs under `docs/`; durable maintainer guidance belongs here or in `README.md`.

## Build, Test, and Development Commands

```bash
go test ./...
```

Runs the full test suite.

```bash
go run ./cmd/diary --help
```

Shows CLI help.

```bash
go build -o bin/diary ./cmd/diary
```

Builds the local CLI binary.

## Coding Style & Naming Conventions

Use `gofmt` for all Go files. Keep command files focused by command, for example `internal/cli/record.go` and `internal/cli/self_update.go`.

Command names should stay short and intentional:

- `record`
- `get`
- `list`
- `install-skills`
- `self-update`

Diary storage must remain local-first under `.diary/`. Markdown records are the source of truth; `index.json` is a rebuildable cache.

## Testing Guidelines

Add or update tests for every behavior change. Current coverage focuses on CLI parsing, project resolution, storage behavior, hashing, skill installation, and self-update behavior.

Prefer behavior-oriented test names such as `TestCreateRecordWritesRecordLatestAndIndex`.

Run `go test ./...` before committing.

## Release Guidelines

Releases are tag-driven. Pushing a `v*` tag triggers `.github/workflows/release.yml`, which builds binaries for macOS, Linux, and Windows and publishes a GitHub Release.

Before tagging:

- Run `go test ./...`.
- Confirm README install and command examples are still accurate.
- Avoid retagging an existing version; use the next semver tag.

## Security & Configuration Tips

Never read or commit `.env` files. Diary should avoid recording secrets by default. Do not commit `.diary/`; it is local runtime data and ignored by `.gitignore`.

When working on `install-skills`, keep generated skills harness-agnostic unless target-specific behavior is strictly required. Current installed skills are:

- `diary-get`
- `diary-record`
- `diary-list`

`diary-record` should preserve the compaction workflow, including files in scope, files out of scope, verification, blockers, and next steps.

## Commit & Pull Request Guidelines

Use concise imperative commit messages, for example:

```text
Improve README usage guide
```

Pull requests should include a summary, testing notes, and any release impact.
