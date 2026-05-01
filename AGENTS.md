# Repository Guidelines

## Project Structure & Module Organization

This repository is currently documentation-first. Product planning documents live in `docs/plans/`, including the PRD for the Diary CLI. Keep future planning documents in that directory using the pattern:

```text
docs/plans/YYYY-MM-DD-short-topic.md
```

When implementation begins, place the Go CLI entrypoint in `cmd/diary/` and internal packages under `internal/`. Generated build output should go in a directory that can be ignored by Git.

## Build, Test, and Development Commands

This project uses Go and builds a `diary` CLI from `cmd/diary`.

Useful commands:

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

Prefer clear, small files with explicit names. Use lowercase, hyphen-separated names for Markdown documents, for example `2026-05-01-ai-harness-journal-cli-prd.md`.

For future CLI code, keep command names short and stable. The initial command vocabulary is `record`, `get`, and `list`, and storage should use `.diary/`.
Harness integration uses `install-skills` to install `diary-get`, `diary-record`, and `diary-list` into supported targets.

## Testing Guidelines

There is no test framework yet. When implementation starts, add tests for command parsing, project name resolution, `.diary/` storage behavior, JSON output, and error cases.

Name tests after behavior rather than implementation details, for example `records_message_with_inferred_project`.

## Commit & Pull Request Guidelines

Use concise imperative commit messages such as:

```text
Add initial record command PRD
```

Pull requests should include a short summary, relevant issue links, testing notes, and screenshots only when UI output is involved.

## Security & Configuration Tips

Never read or commit `.env` files. The Diary CLI design should avoid recording secrets by default and should support ignore patterns for sensitive files. Keep `.diary/` local unless the user explicitly decides to version or sync it.

## Previous Agent Context

```text
<claude-mem-context>
# Memory Context

# [diary] recent context, 2026-05-01 1:40pm GMT+8

Legend: 🎯session 🔴bugfix 🟣feature 🔄refactor ✅change 🔵discovery ⚖️decision
Format: ID TIME TYPE TITLE
Fetch details: get_observations([IDs]) | Search: mem-search skill

Stats: 20 obs (7,247t read) | 477,963t work | 98% savings

### May 1, 2026
3239 11:22a ⚖️ CLI Tool for AI Harness Implementation Journaling — PRD Initiated
3241 " 🔵 diary Project Directory — Nearly Empty, No Git Repo
3242 11:23a 🟣 AI Harness Journal CLI — Full PRD Written
3256 11:34a ⚖️ diary CLI — Storage Renamed to `.diary/`, Command Set Narrowed to Three
3278 11:55a 🟣 AI Harness Journal CLI PRD — Project Resolution Rules and Storage Model Added
3279 11:57a ✅ diary Repo — AGENTS.md Created With Repository Guidelines for AI Agents
3281 12:03p ⚖️ Existing Component Preserved — To Be Replaced by Diary Feature
3282 12:06p 🟣 Diary CLI PRD — Hash Model Added to Data Spec
3283 " 🟣 Diary CLI PRD — `get` Retrieval Behavior Formally Specified
3284 " 🟣 Diary CLI PRD — Storage Layout Updated With `records/` Dir and `index.json`
3285 " ⚖️ Diary CLI PRD — `list` Is an Inventory Command, Not a Context Retriever
3315 12:22p 🔵 User Preference — Go Over Node.js
3331 12:37p ⚖️ diary CLI — Cobra Chosen as CLI Framework
3332 " 🟣 diary — Full Go Module Scaffolded With All Initial Source Files
3333 " 🔵 diary — `go mod tidy` Fails With Build Cache Permission Error in Sandboxed Environment
3335 12:40p ⚖️ AI Harness Journal CLI — PRD Initiated for Local-First Session Continuity Tool
3336 12:41p 🟣 diary CLI — Full Go Implementation Shipped and Verified
3349 1:06p 🟣 install-skills Command — Planned for AI Harness CLI
3350 1:09p 🔵 diary CLI — Current Command Structure Confirmed Before install-skills Work
3353 1:12p 🟣 diary CLI — install-skills Command Implemented

Access 478k tokens of past work via get_observations([IDs]) or mem-search skill.
</claude-mem-context>
```


<claude-mem-context>
# Memory Context

# [diary] recent context, 2026-05-01 11:55am GMT+8

Legend: 🎯session 🔴bugfix 🟣feature 🔄refactor ✅change 🔵discovery ⚖️decision
Format: ID TIME TYPE TITLE
Fetch details: get_observations([IDs]) | Search: mem-search skill

Stats: 5 obs (2,071t read) | 42,403t work | 95% savings

### May 1, 2026
3239 11:22a ⚖️ CLI Tool for AI Harness Implementation Journaling — PRD Initiated
3241 " 🔵 diary Project Directory — Nearly Empty, No Git Repo
3242 11:23a 🟣 AI Harness Journal CLI — Full PRD Written
3256 11:34a ⚖️ diary CLI — Storage Renamed to `.diary/`, Command Set Narrowed to Three
3278 11:55a 🟣 AI Harness Journal CLI PRD — Project Resolution Rules and Storage Model Added

Access 42k tokens of past work via get_observations([IDs]) or mem-search skill.
</claude-mem-context>
