# AI Harness Journal CLI PRD

Date: 2026-05-01
Status: Draft

## Summary

Build a local-first CLI tool that lets AI coding harnesses journal implementation progress, decisions, blockers, and compacted context into durable local files. The journal should act like a structured memory layer between runs: readable by humans, easy for AI agents to append to, and retrievable by future runs as concise context.

The product is not a general note-taking app. Its primary job is to preserve implementation continuity when an AI harness session ends, compacts, crashes, or resumes later with limited context.

## Problem

AI harnesses frequently lose useful implementation context across runs:

- Why a design or technical decision was made.
- What files were touched and what remains unfinished.
- Which commands were run and what failed.
- What assumptions should be carried forward.
- What a future agent should read first before continuing.

Existing solutions are either too informal, such as scattered notes, or too heavyweight, such as databases and project management tools. Agents need a predictable local protocol for writing and reading compact implementation memory.

## Goals

- Provide a CLI that any AI harness can call during or at the end of a run.
- Store journal data locally in human-readable Markdown with structured frontmatter.
- Support compact summaries that can be injected into a future AI prompt.
- Make journals project-scoped by default.
- Preserve implementation continuity without requiring cloud services, accounts, or background daemons.
- Keep v1 simple enough to be reliable under automated agent use.

## Non-Goals

- Cloud sync.
- Multi-user collaboration.
- Hosted dashboards.
- Semantic vector search.
- Automatic code analysis beyond metadata explicitly provided by the caller.
- Replacing Git history, issue trackers, or project documentation.

## Target Users

- AI coding harnesses that need persistent local context.
- Developers using multiple AI agents on the same project.
- Developers who want inspectable records of AI implementation work.
- Future automation that needs to retrieve a compact "what happened last time" summary.

## Product Approach

Use a structured-local-files approach.

Each project gets a journal directory named `.diary/` by default. Entries are Markdown files with YAML frontmatter and predictable naming. The initial CLI provides five commands: `record`, `get`, `list`, `install-skills`, and `self-update`.

Records are project-aware. Callers may pass `--project <name>`, but the CLI should infer a project name when the flag is omitted.

This approach is recommended because it is portable, debuggable, easy to version or ignore, and friendly to both humans and AI harnesses. SQLite can be considered later if retrieval requirements become more complex.

The CLI should be implemented in Go. Go fits the product because it can ship as a single binary, has strong filesystem support, avoids runtime dependency setup, and is practical for agent-driven installation in many environments.

Use Cobra for CLI command structure. Cobra fits the product because `diary` is subcommand-oriented from the start, needs consistent help output, and will likely grow additional maintenance commands later. Avoid Viper in v1; configuration needs are small enough to handle directly.

## Core Concepts

### Run

A run is one contiguous AI harness session in a project. It has:

- A unique run id.
- Start and optional end timestamps.
- Harness name and version when available.
- Working directory.
- Branch or commit metadata when available.
- A summary generated during or at the end of the session.

### Event

An event is an append-only journal record inside a run. Examples:

- `started`
- `intent`
- `decision`
- `changed`
- `command`
- `blocked`
- `handoff`
- `completed`

Events should be concise and structured enough to summarize later.

### Compact Summary

A compact summary is the handoff artifact intended for the next run. It should include:

- Current task.
- Decisions made.
- Files changed.
- Commands run and important results.
- Known blockers.
- Next recommended steps.
- References to source journal entries.

### Reference

A reference points to local context that future runs should inspect:

- File paths.
- Line numbers when available.
- Commands.
- Journal event ids.
- Related docs.

### Hash

A hash is a stable content fingerprint for a record or run file. Hashes support integrity checks, exact references, future deduplication, and optional retrieval by hash prefix. They should not be the primary day-to-day retrieval mechanism.

Each record should have both:

- A human-debuggable `id`.
- A content-derived `hash`.

The hash should be computed from normalized frontmatter and body content, excluding the `hash` field itself.

### Project

A project is the logical workspace the record belongs to. It is stored on every run and entry, even though v1 stores data locally inside the current project's `.diary/` directory.

Project resolution order:

1. Use `--project <name>` when supplied.
2. Use a configured project name from `.diary/config.yml` when present.
3. Use the Git repository root directory name when inside a Git repo.
4. Use the current working directory basename.

## Functional Requirements

### Record Context

The CLI must support recording structured implementation context:

```bash
diary record "Use Markdown files for v1 storage"
diary record --project campaign-builder "Updated importer validation flow" --file app/Import.php
echo "Tests fail because local database is unavailable" | diary record --project campaign-builder
```

Expected behavior:

- Creates `.diary/` automatically if it does not exist.
- Requires a compact/context message, either as an argument or through stdin.
- Resolves the project name before writing the record.
- Creates or appends to the active run by default.
- Allows explicit `--run <id>` when needed.
- Accepts message body, tags, files, and references.
- Keeps entries timestamped.
- Supports caller-provided compact/context summaries through stdin.
- Updates a latest context pointer or equivalent when recording handoff-style entries.

### Retrieve Context

The CLI must support printing context for the next harness run:

```bash
diary get
diary get --since 7d
diary get --task "import workflow"
diary get --id 2026-05-01T103000Z-codex-a7f3c9
diary get --hash abc123
```

Expected behavior:

- Acts as a context retriever for the next AI harness run.
- Resolves the project using the same rules as `diary record`.
- Prints the latest compact summary by default.
- Can include recent decisions and blockers.
- Can retrieve an exact record by id or hash prefix.
- Can emit Markdown or JSON.
- Keeps output bounded by a configurable character or token budget.

### List Runs and Records

The CLI must support listing stored Diary objects without returning full context:

```bash
diary list
diary list --projects
diary list --project campaign-builder
```

Expected behavior:

- Lists records and runs for the resolved project by default.
- Shows id, hash prefix, timestamp, harness, task or message preview, status, and tags.
- Can list known projects with `--projects`.
- Supports `--json`.
- Can filter by project, status, tag, harness, or date in later versions.
- Does not assemble prompt-ready context; callers should use `diary get` for that.

### Install Harness Skills

The CLI must support installing Diary usage instructions into supported AI harnesses:

```bash
diary install-skills --target codex
diary install-skills --target claude
diary install-skills --target all --dry-run
```

Expected behavior:

- Supports Codex and Claude in v1.
- Installs three intentional Diary skills for each supported target: `diary-get`, `diary-record`, and `diary-list`.
- `diary-get` retrieves previous context before work.
- `diary-record` compacts the latest implementation, classifies files in scope and out of scope, and records the handoff.
- `diary-list` inventories available projects, records, ids, and hash prefixes.
- The installed skills do not require harness-specific flags.
- Writes to default skill paths unless `--path <dir>` is supplied.
- Refuses to overwrite an existing skill unless `--force` is supplied.
- Supports `--dry-run` and `--json`.
- Does not record secrets or read `.env` files.

### Self Update

The CLI must support updating itself from GitHub Releases:

```bash
diary self-update
diary self-update --version v0.0.1
diary self-update --dry-run
```

Expected behavior:

- Resolves the latest release unless `--version` is supplied.
- Downloads the matching archive for the current OS and architecture.
- Replaces the running `diary` binary on macOS and Linux.
- Tries `sudo` when the binary directory is not writable.
- Supports `--dry-run` and `--json`.
- Returns a clear unsupported message on Windows in v1.

## CLI Requirements

- Commands must be safe for non-interactive harness use.
- Every command that changes state should have predictable stdout.
- Errors should be concise and actionable.
- `--json` should be supported for automation.
- The tool should work without a daemon.
- The active run should be discoverable without shell-specific state.
- The CLI should avoid destructive behavior unless explicitly requested.
- The CLI should build into a single `diary` binary.
- The codebase should use a small command structure that can grow without making the v1 commands hard to follow.
- The CLI should use Cobra for command parsing, subcommands, shared flags, and help text.

## Go Project Structure

Recommended implementation layout:

```text
cmd/diary/
  main.go
internal/cli/
  root.go
  record.go
  get.go
  list.go
  install_skills.go
internal/install/
  install.go
  templates.go
internal/project/
internal/storage/
internal/hash/
internal/render/
```

Responsibilities:

- `cmd/diary`: binary entrypoint.
- `internal/cli`: command parsing and output mode handling.
- `internal/install`: target skill installation and embedded Diary skill templates.
- `internal/project`: project name and root resolution.
- `internal/storage`: `.diary/` reads, writes, indexes, and Markdown frontmatter.
- `internal/hash`: content normalization and hashing.
- `internal/render`: Markdown and JSON output formatting.

## Storage Requirements

Default local storage:

```text
.diary/
  config.yml
  index.md
  projects/
    diary/
      latest.md
      index.json
      runs/
        2026-05-01T103000Z-codex-abc123.md
      records/
        2026-05-01T103000Z-codex-a7f3c9.md
```

Run and record files should use Markdown with YAML frontmatter:

```yaml
---
id: 2026-05-01T103000Z-codex-abc123
project: diary
hash: sha256:abc123...
parent_hash:
harness: codex
task: Build import workflow
status: active
started_at: 2026-05-01T10:30:00Z
ended_at:
files: []
tags: []
---
```

Run file bodies should contain sections for events, compact summary, references, and handoff notes. Record file bodies may be shorter and contain the recorded context message plus references.

## Data Model

Minimum run fields:

- `id`
- `project`
- `hash`
- `parent_hash`
- `harness`
- `task`
- `status`
- `started_at`
- `ended_at`
- `cwd`
- `git_branch`
- `git_commit`
- `tags`

Minimum record fields:

- `id`
- `hash`
- `type`
- `timestamp`
- `project`
- `body`
- `files`
- `refs`
- `tags`

## Retrieval Behavior

The `get` command should prioritize:

1. Resolve the project.
2. Return an exact record when `--id` or `--hash` is supplied.
3. Read the project's latest compact summary.
4. Include active blockers.
5. Include recent decisions.
6. Include files changed in the latest run.
7. Include next steps.

When output must be shortened, it should remove lower-priority historical detail before removing blockers or next steps.

Hashes should be used for exact references and integrity checks. Normal context retrieval should rely on project, latest summary, ids, dates, tags, and task filters.

The `list` command should be treated as an inventory command. It helps humans and harnesses discover available projects, records, runs, ids, and hash prefixes before calling `get`.

## Error Handling

- If no journal exists, `diary record` should create `.diary/` automatically, while `diary get` and `diary list` should explain that no records exist yet.
- If no active run exists, `diary record` should create one automatically unless `--run <id>` is provided.
- If the project cannot be inferred, `diary record` should ask the caller to pass `--project <name>`.
- If a run id is ambiguous, the CLI should list matching ids.
- If a hash prefix matches multiple records, the CLI should list the matching records and ask for a longer prefix.
- If a Markdown file cannot be parsed, the CLI should preserve it and report the parse issue.
- If Git metadata is unavailable, the CLI should continue without it.

## Security and Privacy

- The CLI must not read `.env` files.
- The CLI should avoid recording secrets by default.
- The CLI should support a configurable ignore list for files and patterns.
- The journal should be local-only unless the user explicitly versions or syncs it.
- The tool should make it easy to add `.diary/` to `.gitignore`, but should not do so automatically in v1.

## Success Metrics

- A new harness run can retrieve useful context with one command.
- A developer can inspect the journal files without the CLI.
- The initial command set covers recording, retrieving, and listing journal context.
- Supported targets can install `diary-get`, `diary-record`, and `diary-list` with one command.
- Records can be referenced exactly by id or hash.
- The tool can be used repeatedly in the same project without manual cleanup.
- Context output remains concise enough to paste into an AI prompt.

## Open Questions

- Should the binary name be `diary`, or should the package expose another name while storing data in `.diary/`?
- Should active run state live in `.diary/active` or be inferred from the latest run with `status: active`?
- Should compaction be fully caller-provided in v1, or should the CLI include a basic local summarizer template?
- Should `.diary/` be ignored by default, or should the CLI ask interactively when used by a human?
- Should there be a global journal mode for cross-project memories, or only project-scoped journals in v1?

## V1 Scope

V1 should include:

- `record`
- `get`
- `list`
- `install-skills`
- `self-update`
- Markdown plus YAML frontmatter storage.
- Content hashes for records and run files.
- Rebuildable `index.json` files for faster lookup.
- JSON output mode.
- Basic file references.
- Configurable output budget.

V1 should exclude:

- Cloud sync.
- Daemon mode.
- Vector search.
- Web UI.
- Team sharing.
- Automatic secret scanning beyond simple ignore patterns.

## Future Enhancements

- SQLite storage backend.
- Semantic retrieval.
- Built-in prompt templates for different harnesses.
- Integration adapters for Codex, Claude Code, Cursor, and other agents.
- Optional Git hooks.
- Secret detection.
- Global memory layer.
- Import/export between projects.
