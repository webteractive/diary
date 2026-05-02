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

Rename the current project's Diary mapping after a project rename:

```bash
diary rename campaign-builder-api
```

Initialize a harness reminder:

```bash
diary init --target codex --install-skills
```

Migrate old project-local records to the user-level store:

```bash
diary migrate --from project --to user --dry-run
diary migrate --from project --to user
```

## Commands

### `diary record`

Stores implementation context for the current project.

```bash
diary record "Updated install-skills to install diary-get, diary-record, and diary-list."
diary record --project campaign-builder "Fixed importer validation and added tests."
diary record --file internal/storage/record.go "Changed record hashing behavior."
diary record --root ~/Documents/work-diary "Recorded context in a private Diary repo."
```

Messages can also come from stdin:

```bash
cat handoff.md | diary record
```

If `--project` is not supplied, Diary resolves the project from `.diary/config.yml`, the Git root directory, or the current directory name.

After writing the record, `diary record` prints the record id, content hash, and Markdown file path.

### `diary get`

Retrieves prompt-ready context for the resolved project.

```bash
diary get
diary get --project campaign-builder
diary get --id 2026-05-01T103000Z-a7f3c9
diary get --hash abc123
diary get --root ~/Documents/work-diary
diary get --json
```

By default, `get` returns the latest recorded context. Exact lookup by id or hash prefix is available when a harness needs a specific record.

### `diary list`

Lists stored Diary inventory without returning full context.

```bash
diary list
diary list --projects
diary list --project campaign-builder --json
diary list --root ~/Documents/work-diary
```

Use this to discover projects, record ids, and hash prefixes before calling `diary get`.

### `diary rename`

Renames the mapped Diary project for the current checkout. Use this after renaming a repository or changing its configured Diary project name so `diary list --projects` shows the new project id while preserving existing records.

```bash
diary rename campaign-builder-api
diary rename campaign-builder-api --root ~/Documents/work-diary
```

`rename` finds the project by the current checkout root, moves the existing records directory from the old mapped project id to the new one, and updates `projects.json`. It does not rewrite the frontmatter in existing record Markdown files.

### `diary init`

Installs a small Diary instruction for supported harness targets. The instruction reminds the harness to ask whether it should run the `diary-record` skill before ending a meaningful coding session, and to use the Diary CLI from the project directory so Diary can resolve the current project and configured storage location.

```bash
diary init --target codex
diary init --target claude
diary init --target all
diary init --target all --install-skills
diary init --target codex --scope project
diary init --target codex --dry-run
```

By default, `init` writes global harness instructions:

- Codex: `~/.codex/AGENTS.md`
- Claude: `~/.claude/CLAUDE.md`

Project-scoped initialization writes to the current project root:

- Codex: `AGENTS.md`
- Claude: `CLAUDE.md`

Diary-managed instructions are wrapped in a `<diary>...</diary>` block. Existing blocks are not overwritten unless `--force` is supplied. When `--install-skills` is supplied, the Diary skills are installed globally for the selected target.

### `diary migrate`

Migrates records between storage locations for the resolved project.

```bash
diary migrate --from project --to user --dry-run
diary migrate --from project --to user
diary migrate --from user --to project
diary migrate --from project --to "$HOME/Documents/work-diary"
diary migrate --from "$HOME/.diary" --to "$HOME/Documents/work-diary"
```

Storage names:

- `project` - the current project's `.diary/`
- `user` - the default user-level Diary store
- any path - a custom Diary root, such as a private Git repository

Migration copies records and rebuilds the destination index. It does not delete the source unless `--delete-source` is supplied. Existing destination records are not overwritten unless `--force` is supplied.

### `diary install-skills`

Installs intent-specific skills for supported harness targets.

```bash
diary install-skills --target codex
diary install-skills --target claude
diary install-skills --target all
diary install-skills --target codex --dry-run
```

Installed skills:

- `diary-init` - install Diary reminder instructions and optional usage skills
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

Diary writes records to a user-level Diary store by default:

```text
~/.diary/
  projects.json
  projects/
    <project-id>/
      latest.md
      index.json
      records/
        <record-id>.md
```

`projects.json` maps project paths to stable project ids so multiple projects with the same directory name do not collide. The id includes the sanitized project name plus a short hash of the root path. If the project is renamed, run `diary rename <new-project-name>` from that checkout to move the stored records to the new id.

If a project already has `.diary/`, Diary keeps using that project-local store for backward compatibility. To write records somewhere else, pass `--root <path>` or set `DIARY_ROOT`. This is useful when you want to keep Diary records in a separate private Git repository:

```bash
DIARY_ROOT="$HOME/Documents/work-diary" diary record "..."
diary record --root "$HOME/Documents/work-diary" "..."
```

Records are Markdown files with YAML frontmatter. Markdown remains the source of truth; `index.json` is a rebuildable lookup cache.

## Security Notes

- Do not record secrets, credentials, tokens, private keys, or `.env` contents.
- New projects use the user-level Diary store by default to avoid adding `.diary/` noise to public repositories.
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

Show the version:

```bash
diary --version
diary -v
```

Build:

```bash
go build -o bin/diary ./cmd/diary
```

Release builds are created by GitHub Actions when a `v*` tag is pushed. Local builds report `dev`; release builds report the tag.
