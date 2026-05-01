# Diary

Diary is a local-first CLI for preserving AI implementation context across runs.

It gives agents and developers a small shared memory layer: record compact handoffs into `.diary/`, retrieve the latest context before work, list prior records, and install optional skills for supported harnesses.

## Why

AI coding sessions often lose the details that matter most:

- what changed and why
- which files were intentionally touched
- which dirty files were unrelated
- what tests or commands were run
- what remains blocked or unfinished

Diary stores that context locally in inspectable Markdown files so the next run can continue with less guesswork.

## Install

Install the latest release:

```bash
curl -fsSL https://raw.githubusercontent.com/webteractive/diary/main/scripts/install.sh | sh
```

Install a specific version:

```bash
curl -fsSL https://raw.githubusercontent.com/webteractive/diary/main/scripts/install.sh | env DIARY_VERSION=v0.0.1 sh
```

Install to a user-writable directory:

```bash
curl -fsSL https://raw.githubusercontent.com/webteractive/diary/main/scripts/install.sh | env DIARY_INSTALL_DIR="$HOME/.local/bin" sh
```

Build from source:

```bash
go build -o bin/diary ./cmd/diary
```

## Quick Start

Record a compact handoff:

```bash
diary record "Implemented record storage, added hashing tests, and left release docs pending."
```

Retrieve context at the start of the next run:

```bash
diary get
```

List available records:

```bash
diary list
```

List known projects:

```bash
diary list --projects
```

## Commands

### `diary record`

Stores implementation context for the current project.

```bash
diary record "Updated install-skills to install diary-get, diary-record, and diary-list."
diary record --project campaign-builder "Fixed importer validation and added tests."
diary record --file internal/storage/record.go "Changed record hashing behavior."
```

Messages can also come from stdin:

```bash
cat handoff.md | diary record
```

If `--project` is not supplied, Diary resolves the project from `.diary/config.yml`, the Git root directory, or the current directory name.

### `diary get`

Retrieves prompt-ready context for the resolved project.

```bash
diary get
diary get --project campaign-builder
diary get --id 2026-05-01T103000Z-codex-a7f3c9
diary get --hash abc123
diary get --json
```

By default, `get` returns the latest recorded context. Exact lookup by id or hash prefix is available when a harness needs a specific record.

### `diary list`

Lists stored Diary inventory without returning full context.

```bash
diary list
diary list --projects
diary list --project campaign-builder --json
```

Use this to discover projects, record ids, and hash prefixes before calling `diary get`.

### `diary install-skills`

Installs intent-specific skills for supported harness targets.

```bash
diary install-skills --target codex
diary install-skills --target claude
diary install-skills --target all
diary install-skills --target codex --dry-run
```

Installed skills:

- `diary-get` - retrieve prior context before work
- `diary-record` - compact the latest implementation and record it
- `diary-list` - list available projects and records

Default install locations:

- Codex: `~/.codex/skills/`
- Claude: `~/.claude/skills/`

Existing skill files are not overwritten unless `--force` is supplied.

### `diary self-update`

Updates the installed binary from GitHub Releases.

```bash
diary self-update
diary self-update --version v0.0.1
diary self-update --dry-run
```

Self-update currently supports macOS and Linux. Windows builds are published, but Windows self-replacement is intentionally not enabled yet.

## Storage

Diary writes local project data under `.diary/`:

```text
.diary/
  projects/
    <project>/
      latest.md
      index.json
      records/
        <record-id>.md
```

Records are Markdown files with YAML frontmatter. Markdown remains the source of truth; `index.json` is a rebuildable lookup cache.

## Security Notes

- Do not record secrets, credentials, tokens, private keys, or `.env` contents.
- `.diary/` is ignored by the repository `.gitignore` by default.
- Diary is local-first; it does not sync data unless you choose to version or share `.diary/`.

## Development

Run tests:

```bash
go test ./...
```

Run locally:

```bash
go run ./cmd/diary --help
```

Build:

```bash
go build -o bin/diary ./cmd/diary
```

Release builds are created by GitHub Actions when a `v*` tag is pushed.
