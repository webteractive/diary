# Diary Agent Instructions

Diary is a Go CLI for preserving local AI implementation context across runs. It records compact handoffs as Markdown, retrieves prompt-ready context, lists prior records, migrates old stores, installs harness instructions/skills, and self-updates from GitHub releases.

## Project Shape

- Binary entrypoint: `cmd/diary/`.
- CLI commands: `internal/cli/`, one focused file per command where practical.
- Project/root resolution: `internal/project/`.
- Storage, records, indexes, and project maps: `internal/storage/`.
- Content hashing: `internal/hash/`.
- Markdown/JSON rendering: `internal/render/`.
- Harness instruction and skill installation: `internal/setup/` and `internal/install/`.
- GitHub release self-update logic: `internal/update/`.

Do not reintroduce planning docs under `docs/`. Durable maintainer guidance belongs in this file or `README.md`.

## Core Invariants

- Diary is local-first. Do not add network behavior except where it already exists for `self-update`.
- Markdown records are the source of truth. `index.json` is a rebuildable cache.
- New projects use the user-level Diary store by default under `~/.diary/`.
- Existing project-local `.diary/` stores remain supported for backward compatibility.
- `projects.json` maps project roots to stable project ids. Preserve that root-based mapping when changing project resolution, rename, list, get, record, or migrate behavior.
- Project ids must avoid collisions between checkouts with the same directory name.
- Record ids should remain harness-agnostic.
- Do not read, record, or commit secrets, credentials, tokens, private keys, or `.env` contents.

## Command Expectations

Keep command names short and intentional:

- `record`
- `get`
- `list`
- `rename`
- `init`
- `migrate`
- `install-skills`
- `self-update`

When adding or changing commands:

- Wire the command in `internal/cli/root.go`.
- Add behavior tests in `internal/cli/`.
- Keep flag names explicit and consistent with existing commands.
- Preserve `--root` behavior for commands that operate on the user-level store.
- Avoid creating or mutating `projects.json` from read-only commands such as `get` and `list` unless the command's purpose is to write.

## Storage Rules

Use structured storage helpers instead of ad hoc path manipulation when possible:

- `NewPaths`
- `NewDiaryRootPaths`
- `ResolveStore`
- `ResolveStoreForRoot`
- `ResolveNamedStore`
- `ReadProjectMap`
- `ReadRecords`
- `RebuildIndex`

When moving or copying records:

- Preserve existing Markdown records.
- Rebuild destination indexes after writes.
- Do not overwrite existing destination records unless an explicit force option exists.
- Do not delete source records unless an explicit delete option exists.

For renamed projects, preserve records by root mapping. A rename should update `projects.json` and move the mapped project directory rather than creating a second unrelated project.

## Coding Style

- Use `gofmt` on every touched Go file.
- Keep package boundaries narrow and boring.
- Prefer behavior-oriented helpers over clever abstractions.
- Keep Cobra command files readable: parse flags in the command file, delegate storage behavior to `internal/storage`.
- Use clear error messages that name the relevant project, root, id, hash, or path.

## Tests

Add or update tests for every behavior change. Current coverage focuses on:

- CLI parsing and command behavior.
- Project resolution.
- Storage root selection and project maps.
- Record creation, parsing, lookup, and indexing.
- Migration behavior.
- Skill installation.
- Self-update behavior.

Prefer behavior-oriented test names such as `TestCreateRecordWritesRecordLatestAndIndex`.

Run the full suite before handing off implementation work:

```bash
go test ./...
```

If the sandbox blocks Go's build cache, rerun the same command with the required approval rather than changing cache paths.

## Documentation

After implementing a behavior change, update `README.md` in the same work session so command examples, storage notes, and usage guidance stay current.

README updates should cover:

- New or changed commands and flags.
- Storage behavior that affects where records live.
- Migration, rename, or compatibility notes.
- Install, update, and release instructions when relevant.

## Harness Skills And Instructions

When working on `install-skills`, keep generated skills harness-agnostic unless target-specific behavior is strictly required.

Current installed skills:

- `diary-init`
- `diary-get`
- `diary-record`
- `diary-list`

`diary-init` should install a short `<diary>...</diary>` instruction block. `diary-record` should preserve the compaction workflow, including files in scope, files out of scope, verification, blockers, and next steps.

## Release Guidance

Releases are tag-driven. Pushing a `v*` tag triggers `.github/workflows/release.yml`, which builds binaries for macOS, Linux, and Windows and publishes a GitHub Release.

Before tagging:

- Run `go test ./...`.
- Confirm README install and command examples are accurate.
- Avoid retagging an existing version; use the next semver tag.

## Git And PRs

- Never commit or push unless the user explicitly asks.
- Never add `Co-Authored-By` to commits.
- Use concise imperative commit messages, for example `Improve README usage guide`.
- Pull requests should include a summary, testing notes, and any release impact.
